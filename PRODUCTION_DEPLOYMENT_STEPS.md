# Production Deployment Steps - Version 2.0

## Current Production Structure
```
/home/ubuntu/
├── docker-compose.yml
├── migrate-v2.0-production.sh
└── (Docker containers running)
```

## Quick Deployment Guide

### Step 1: Backup Database (IMPORTANT!)

```bash
# SSH to EC2
ssh ubuntu@your-ec2-ip

# Backup MongoDB
sudo docker exec mongodb mongodump \
  --uri="mongodb://admin:password123@localhost:27017/cafe_pos?authSource=admin" \
  --out=/backup/cafe_pos_$(date +%Y%m%d_%H%M%S)

# Verify backup
ls -lh /backup/
```

### Step 2: Build and Push Docker Images (Local Machine)

```bash
# On your local machine
cd /path/to/cafe-pos

# Build and push version 2.0
./build_docker_hub.sh
```

Wait for build to complete. Verify on Docker Hub:
- https://hub.docker.com/r/linhtranphu/cafe-pos-backend/tags
- https://hub.docker.com/r/linhtranphu/cafe-pos-frontend/tags

### Step 3: Update docker-compose.yml on Production

```bash
# SSH to EC2
ssh ubuntu@your-ec2-ip

# Edit docker-compose.yml
sudo nano docker-compose.yml
```

Update image tags to version 2.0:
```yaml
services:
  backend:
    image: linhtranphu/cafe-pos-backend:2.0  # Change from :latest to :2.0
    # ... rest of config

  frontend:
    image: linhtranphu/cafe-pos-frontend:2.0  # Change from :latest to :2.0
    # ... rest of config
```

Save and exit (Ctrl+X, Y, Enter)

### Step 4: Pull New Images

```bash
# Pull new images
sudo docker-compose pull

# Verify images downloaded
sudo docker images | grep cafe-pos
```

You should see:
```
linhtranphu/cafe-pos-backend   2.0      ...
linhtranphu/cafe-pos-frontend  2.0      ...
```

### Step 5: Copy Migration Script to EC2

```bash
# From local machine
scp migrate-v2.0-production.sh ubuntu@your-ec2-ip:/home/ubuntu/

# Make it executable
ssh ubuntu@your-ec2-ip
chmod +x /home/ubuntu/migrate-v2.0-production.sh
```

### Step 6: Stop Services

```bash
# On EC2
sudo docker-compose down
```

### Step 7: Start New Version

```bash
# Start services with new images
sudo docker-compose up -d

# Wait for containers to be healthy (30-60 seconds)
sudo docker-compose ps
```

### Step 8: Run Migration

```bash
# Run migration script
sudo bash /home/ubuntu/migrate-v2.0-production.sh
```

Type `yes` when prompted.

Expected output:
```
✅ Connected to MongoDB
✅ Shop settings created
✅ Indexes created
✅ Migration completed successfully!
```

### Step 9: Restart Services (to apply changes)

```bash
sudo docker-compose restart
```

### Step 10: Verify Deployment

#### Check Backend Logs
```bash
sudo docker-compose logs -f backend | head -50
```

Look for:
```
✅ MongoDB connected successfully
✅ HTML bill renderer initialized
✅ Chromedp print handler initialized
Server starting on :3000
```

#### Check Frontend
Open browser: https://tacafe.store

#### Test API
```bash
# Get shop settings (should return 200 OK)
curl -s https://tacafe.store/api/manager/shop-settings \
  -H "Authorization: Bearer YOUR_TOKEN" | jq
```

### Step 11: Configure Shop Settings

1. Login as admin
2. Go to: Print Management → Cài Đặt
3. Configure:
   - Shop name
   - Address
   - Phone
   - Upload logo (optional)
   - Enable auto-print

## Troubleshooting

### Issue: Migration script not found in container

**Solution:** Migration code is embedded in Docker image. If missing, rebuild:

```bash
# On local machine
./build_docker_hub.sh

# On EC2
sudo docker-compose pull
sudo docker-compose up -d
```

### Issue: 404 on shop-settings API

**Check 1:** Migration ran successfully?
```bash
# Check MongoDB
sudo docker exec mongodb mongo cafe_pos \
  --eval "db.shop_settings.find().pretty()"
```

**Check 2:** Backend running?
```bash
sudo docker-compose ps
sudo docker-compose logs backend | tail -50
```

### Issue: Backend container won't start

**Check logs:**
```bash
sudo docker-compose logs backend
```

Common causes:
- MongoDB not ready (wait 30 seconds)
- Port 3000 already in use
- Environment variables missing

**Solution:**
```bash
sudo docker-compose down
sudo docker-compose up -d
```

### Issue: Frontend shows old version

**Clear browser cache:**
- Chrome: Ctrl+Shift+Delete
- Or hard refresh: Ctrl+Shift+R

**Check frontend version:**
```bash
sudo docker-compose logs frontend | grep "version"
```

## Rollback Plan

If deployment fails:

### Option 1: Rollback Docker Images

```bash
# Edit docker-compose.yml
sudo nano docker-compose.yml

# Change back to previous version
# backend: image: linhtranphu/cafe-pos-backend:1.0
# frontend: image: linhtranphu/cafe-pos-frontend:1.0

# Restart
sudo docker-compose down
sudo docker-compose up -d
```

### Option 2: Rollback Database

```bash
# Restore from backup
sudo docker exec mongodb mongorestore \
  --uri="mongodb://admin:password123@localhost:27017/cafe_pos?authSource=admin" \
  --drop \
  /backup/cafe_pos_YYYYMMDD_HHMMSS
```

⚠️ **Warning:** This will lose data created after backup!

## Post-Deployment Checklist

- [ ] Backend logs show no errors
- [ ] Frontend loads correctly
- [ ] Shop settings API returns 200
- [ ] Can login as admin
- [ ] Print Management page loads
- [ ] Can create test order
- [ ] Auto-print works (if enabled)

## Version 2.0 New Features

✅ Logo support in bills
✅ Optimized chromedp rendering
✅ Better ESC/POS printer compatibility
✅ Vietnamese font support
✅ Dynamic bill height
✅ HTML template management

## Support

If issues persist:

1. **Check logs:**
   ```bash
   sudo docker-compose logs -f
   ```

2. **Check MongoDB:**
   ```bash
   sudo docker exec mongodb mongo cafe_pos
   ```

3. **Restart everything:**
   ```bash
   sudo docker-compose down
   sudo docker-compose up -d
   ```

4. **Contact support** with:
   - Error messages from logs
   - Steps that led to the issue
   - Screenshot if applicable
