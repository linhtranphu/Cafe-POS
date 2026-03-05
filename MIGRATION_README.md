# Migration to Version 2.0 - Quick Reference

## Files Overview

| File | Purpose | Where to Run |
|------|---------|--------------|
| `migrate-v2.0.sh` | Local migration testing | Local machine |
| `migrate-v2.0-production.sh` | Production migration | EC2 server |
| `verify-migration.sh` | Check migration status | Both |
| `backend/cmd/migrate-v2.0/main.go` | Migration logic | Embedded in Docker |
| `PRODUCTION_DEPLOYMENT_STEPS.md` | Detailed deployment guide | Reference |
| `DEPLOYMENT_V2.0_GUIDE.md` | Complete deployment guide | Reference |

## Quick Start - Production Deployment

### 1. Backup (5 minutes)
```bash
ssh ubuntu@your-ec2-ip
sudo docker exec mongodb mongodump \
  --uri="mongodb://admin:password123@localhost:27017/cafe_pos?authSource=admin" \
  --out=/backup/cafe_pos_$(date +%Y%m%d_%H%M%S)
```

### 2. Build & Push (10-15 minutes)
```bash
# Local machine
./build_docker_hub.sh
```

### 3. Deploy (5 minutes)
```bash
# EC2
ssh ubuntu@your-ec2-ip
sudo docker-compose pull
sudo docker-compose down
sudo docker-compose up -d
```

### 4. Migrate (2 minutes)
```bash
# Copy script to EC2 first
scp migrate-v2.0-production.sh ubuntu@your-ec2-ip:/home/ubuntu/

# Then on EC2
sudo bash /home/ubuntu/migrate-v2.0-production.sh
```

### 5. Verify (1 minute)
```bash
# Copy verify script
scp verify-migration.sh ubuntu@your-ec2-ip:/home/ubuntu/

# Run verification
bash /home/ubuntu/verify-migration.sh
```

### 6. Restart (1 minute)
```bash
sudo docker-compose restart
```

**Total time: ~25-30 minutes**

## What the Migration Does

### ✅ Creates (if not exists)
- `shop_settings` collection with default values
- `print_jobs` collection
- `printer_configs` collection
- `print_templates` collection
- `print_notifications` collection

### ✅ Adds
- Indexes for print collections (performance)
- Default shop settings (name, address, phone)

### ✅ Preserves
- All existing orders
- All existing users
- All existing menu items
- All existing ingredients
- All existing shifts
- **NO DATA LOSS**

## Migration Safety

The migration is designed to be:
- **Idempotent**: Can run multiple times safely
- **Non-destructive**: Only adds, never deletes
- **Reversible**: Can rollback if needed

## Verification Checklist

After migration, verify:

```bash
# 1. Check shop settings exist
curl https://tacafe.store/api/manager/shop-settings \
  -H "Authorization: Bearer TOKEN"

# 2. Check backend logs
sudo docker-compose logs backend | grep "shop_settings"

# 3. Run verification script
bash verify-migration.sh
```

Expected results:
- ✅ shop_settings: 1 document
- ✅ print_jobs: collection exists
- ✅ printer_configs: collection exists
- ✅ print_templates: collection exists
- ✅ orders: preserved (same count as before)
- ✅ users: preserved (same count as before)

## Troubleshooting

### Migration fails with "connection refused"
**Cause:** MongoDB not accessible from container

**Solution:**
```bash
# Check MongoDB is running
sudo docker ps | grep mongodb

# Check network
sudo docker network ls
sudo docker network inspect cafe-pos_default
```

### Migration creates shop_settings but API returns 404
**Cause:** Backend not restarted after migration

**Solution:**
```bash
sudo docker-compose restart backend
```

### Want to re-run migration
**Safe to do:**
```bash
sudo bash migrate-v2.0-production.sh
```

Migration checks if data exists before creating, so it's safe to run multiple times.

## Rollback

If you need to rollback:

### Option 1: Rollback Docker images only
```bash
# Edit docker-compose.yml to use old version
sudo nano docker-compose.yml
# Change :2.0 to :1.0 or previous version

sudo docker-compose down
sudo docker-compose up -d
```

### Option 2: Rollback database (DESTRUCTIVE)
```bash
# Restore from backup
sudo docker exec mongodb mongorestore \
  --uri="mongodb://admin:password123@localhost:27017/cafe_pos?authSource=admin" \
  --drop \
  /backup/cafe_pos_YYYYMMDD_HHMMSS
```

⚠️ **Warning:** Option 2 will lose all data created after the backup!

## Support

### Check logs
```bash
# Backend logs
sudo docker-compose logs -f backend

# MongoDB logs
sudo docker logs mongodb

# All logs
sudo docker-compose logs -f
```

### Check database directly
```bash
# Connect to MongoDB
sudo docker exec -it mongodb mongo cafe_pos

# In mongo shell:
> db.shop_settings.find().pretty()
> db.print_jobs.countDocuments()
> db.orders.countDocuments()
```

### Common commands
```bash
# Restart everything
sudo docker-compose restart

# Stop everything
sudo docker-compose down

# Start everything
sudo docker-compose up -d

# View running containers
sudo docker-compose ps

# View images
sudo docker images | grep cafe-pos
```

## Version 2.0 Features

### New Features
- ✅ Logo support in bills
- ✅ Optimized chromedp rendering (80mm bills)
- ✅ ESC/POS GS v 0 command (better compatibility)
- ✅ Vietnamese font support
- ✅ Dynamic bill height (memory efficient)
- ✅ HTML template management
- ✅ Visual bill preview

### Technical Changes
- Go 1.24 (from 1.21)
- template.URL for logo rendering
- Manual GS v 0 raster conversion
- Font loading wait in chromedp
- Improved font stack

### Breaking Changes
- Requires shop_settings collection (created by migration)
- Requires Go 1.24+ (Docker image updated)

## Questions?

1. Read `PRODUCTION_DEPLOYMENT_STEPS.md` for detailed guide
2. Run `verify-migration.sh` to check status
3. Check logs: `sudo docker-compose logs -f`
4. Test locally first with `migrate-v2.0.sh`
