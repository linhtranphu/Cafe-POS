# Deployment Guide - Version 2.0

## Overview
Version 2.0 includes major improvements to bill printing with logo support, optimized chromedp rendering, and ESC/POS compatibility fixes.

## Pre-Deployment Checklist

- [ ] Backup production database
- [ ] Test migration script locally
- [ ] Build and push Docker images
- [ ] Verify Docker images on Docker Hub

## Deployment Steps

### 1. Backup Production Database

```bash
# SSH to EC2
ssh ubuntu@tacafe.store

# Backup MongoDB
sudo docker exec mongodb mongodump \
  --uri="mongodb://admin:password123@localhost:27017/cafe_pos?authSource=admin" \
  --out=/backup/cafe_pos_$(date +%Y%m%d_%H%M%S)
```

### 2. Build and Push Docker Images

```bash
# On local machine
./build_docker_hub.sh
```

This will build and push:
- `linhtranphu/cafe-pos-backend:2.0`
- `linhtranphu/cafe-pos-backend:latest`
- `linhtranphu/cafe-pos-frontend:2.0`
- `linhtranphu/cafe-pos-frontend:latest`

### 3. Update docker-compose.yml on Production

```bash
# SSH to EC2
ssh ubuntu@tacafe.store

# Navigate to project directory
cd /home/ubuntu/cafe-pos

# Pull latest images
sudo docker-compose pull

# Verify images
sudo docker images | grep cafe-pos
```

### 4. Run Database Migration

```bash
# Copy migration script to EC2 (from local machine)
scp migrate-v2.0-production.sh ubuntu@tacafe.store:/home/ubuntu/cafe-pos/
scp -r backend/cmd/migrate-v2.0 ubuntu@tacafe.store:/home/ubuntu/cafe-pos/backend/cmd/

# SSH to EC2
ssh ubuntu@tacafe.store

# Run migration
cd /home/ubuntu/cafe-pos
sudo bash migrate-v2.0-production.sh
```

**What the migration does:**
1. ✅ Creates `shop_settings` collection with default values
2. ✅ Creates indexes for print collections
3. ✅ Ensures print collections exist (print_jobs, printer_configs, print_templates)
4. ✅ Verifies existing orders (no data loss)

### 5. Restart Services

```bash
# Stop services
sudo docker-compose down

# Start services with new images
sudo docker-compose up -d

# Check logs
sudo docker-compose logs -f backend
sudo docker-compose logs -f frontend
```

### 6. Install Vietnamese Fonts (for chromedp)

```bash
# SSH to EC2
ssh ubuntu@tacafe.store

# Install fonts
sudo apt update
sudo apt install -y fonts-noto fonts-noto-cjk

# Verify fonts
fc-list :lang=vi | grep -i noto
```

### 7. Verify Deployment

#### Backend Health Check
```bash
curl https://tacafe.store/api/manager/shop-settings \
  -H "Authorization: Bearer YOUR_TOKEN"
```

Expected: 200 OK with shop settings JSON

#### Frontend Check
- Open: https://tacafe.store
- Login as admin
- Navigate to Print Management
- Verify shop settings load correctly

#### Test Bill Printing
1. Create a test order
2. Check print jobs are created
3. Verify logo renders correctly (if configured)

## Rollback Plan

If deployment fails:

```bash
# SSH to EC2
ssh ubuntu@tacafe.store

# Rollback to previous version
cd /home/ubuntu/cafe-pos
sudo docker-compose down

# Edit docker-compose.yml to use previous version tags
# Then restart
sudo docker-compose up -d
```

## Post-Deployment Tasks

### 1. Configure Shop Settings
- Login as manager
- Go to Print Management → Cài Đặt
- Upload logo (optional)
- Configure shop name, address, phone
- Enable auto-print if needed

### 2. Configure Printers
- Go to Print Management → Máy In
- Add bill printer (ESC/POS compatible)
- Add label printer
- Test connections

### 3. Configure Templates
- Go to Print Management → Templates
- Review default templates
- Customize if needed

## Troubleshooting

### Issue: 404 on /api/manager/shop-settings

**Cause:** Migration not run or failed

**Solution:**
```bash
# Re-run migration
sudo bash migrate-v2.0-production.sh
```

### Issue: Logo not displaying in bills

**Cause:** Vietnamese fonts not installed

**Solution:**
```bash
sudo apt install -y fonts-noto fonts-noto-cjk
sudo docker-compose restart backend
```

### Issue: Printer shows garbage characters

**Cause:** Printer doesn't support ESC/POS or wrong printer type

**Solution:**
1. Verify printer supports ESC/POS
2. Check printer IP and port (usually 9100)
3. Test with simple text print first

### Issue: Backend container won't start

**Cause:** Go version mismatch or missing dependencies

**Solution:**
```bash
# Check logs
sudo docker-compose logs backend

# Rebuild image
sudo docker-compose build --no-cache backend
sudo docker-compose up -d
```

## Version 2.0 Features

### New Features
- ✅ Logo support in bill printing
- ✅ Optimized chromedp rendering (576px width for 80mm bills)
- ✅ ESC/POS GS v 0 command (better printer compatibility)
- ✅ Vietnamese font support
- ✅ Dynamic bill height (saves memory)
- ✅ HTML template management
- ✅ Visual bill preview

### Technical Improvements
- ✅ Go 1.24 support
- ✅ template.URL for logo base64 (prevents HTML escaping)
- ✅ Manual GS v 0 raster conversion (no escpos library)
- ✅ Font loading wait in chromedp
- ✅ Improved font stack for Vietnamese

### Breaking Changes
- ⚠️ Requires Go 1.24+ (Docker image updated)
- ⚠️ Requires shop_settings collection (created by migration)
- ⚠️ Removed HTML Print tab from UI (chromedp used internally)

## Support

If you encounter issues:
1. Check logs: `sudo docker-compose logs -f`
2. Verify database: `sudo docker exec mongodb mongo cafe_pos`
3. Test API endpoints with curl
4. Review this guide's troubleshooting section

## Rollback Database (if needed)

```bash
# SSH to EC2
ssh ubuntu@tacafe.store

# Restore from backup
sudo docker exec mongodb mongorestore \
  --uri="mongodb://admin:password123@localhost:27017/cafe_pos?authSource=admin" \
  --drop \
  /backup/cafe_pos_YYYYMMDD_HHMMSS
```

**⚠️ Warning:** This will restore database to backup state, losing any data created after backup.
