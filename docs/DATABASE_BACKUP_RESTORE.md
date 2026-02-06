# 📦 Database Backup & Restore Guide

## 🎯 Overview

Hướng dẫn backup và restore MongoDB database giữa EC2 và máy local.

## 📥 Backup Database từ EC2

### Cách 1: Sử dụng Script (Recommended)

```bash
# Set environment variables
export EC2_HOST=your-ec2-ip
export EC2_USER=ubuntu
export EC2_KEY=~/.ssh/your-key.pem

# Run backup script
./scripts/backup-db-from-ec2.sh
```

Hoặc inline:

```bash
EC2_HOST=your-ip EC2_USER=ubuntu EC2_KEY=~/.ssh/key.pem ./scripts/backup-db-from-ec2.sh
```

### Cách 2: Manual Steps

#### Step 1: SSH vào EC2
```bash
ssh -i ~/.ssh/your-key.pem ubuntu@your-ec2-ip
```

#### Step 2: Dump MongoDB từ container
```bash
# Create backup directory
mkdir -p ~/backups

# Dump database
docker exec mongodb mongodump \
  --db cafe_pos \
  --out /tmp/backup

# Copy from container to host
docker cp mongodb:/tmp/backup ~/backups/

# Create tarball
cd ~/backups
tar -czf cafe_pos_backup.tar.gz backup/

# Cleanup
rm -rf backup/
```

#### Step 3: Download về local
```bash
# On local machine
scp -i ~/.ssh/your-key.pem \
  ubuntu@your-ec2-ip:~/backups/cafe_pos_backup.tar.gz \
  ./backups/
```

#### Step 4: Extract
```bash
cd backups
tar -xzf cafe_pos_backup.tar.gz
```

## 📤 Restore Database vào Local

### Cách 1: Sử dụng Script (Recommended)

```bash
./scripts/restore-db-to-local.sh ./backups/cafe_pos_backup_20260205_150000
```

### Cách 2: Manual Steps

#### Option A: MongoDB trong Docker
```bash
# Copy backup vào container
docker cp ./backups/backup/cafe_pos mongodb:/tmp/

# Restore
docker exec mongodb mongorestore \
  --db cafe_pos \
  --drop \
  /tmp/cafe_pos

# Cleanup
docker exec mongodb rm -rf /tmp/cafe_pos
```

#### Option B: MongoDB local
```bash
mongorestore \
  --db cafe_pos \
  --drop \
  ./backups/backup/cafe_pos
```

## 🔍 Verify Restore

### Check database stats
```bash
# Docker
docker exec -it mongodb mongosh cafe_pos --eval 'db.stats()'

# Local
mongosh cafe_pos --eval 'db.stats()'
```

### Check collections
```bash
# Docker
docker exec -it mongodb mongosh cafe_pos --eval 'db.getCollectionNames()'

# Local
mongosh cafe_pos --eval 'db.getCollectionNames()'
```

### Check document count
```bash
# Docker
docker exec -it mongodb mongosh cafe_pos --eval '
  db.getCollectionNames().forEach(function(col) {
    print(col + ": " + db[col].countDocuments())
  })
'
```

## 📋 Backup Schedule

### Recommended Schedule
- **Daily**: Automatic backup at 2 AM
- **Weekly**: Full backup every Sunday
- **Before Deploy**: Manual backup before any deployment

### Setup Cron Job on EC2

```bash
# Edit crontab
crontab -e

# Add daily backup at 2 AM
0 2 * * * /home/ubuntu/backup-mongodb.sh

# Add weekly backup every Sunday at 3 AM
0 3 * * 0 /home/ubuntu/backup-mongodb-weekly.sh
```

### Backup Script on EC2

Create `/home/ubuntu/backup-mongodb.sh`:

```bash
#!/bin/bash

BACKUP_DIR="/home/ubuntu/backups"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_NAME="cafe_pos_${TIMESTAMP}"

# Create backup
docker exec mongodb mongodump \
  --db cafe_pos \
  --out /tmp/backup

# Copy from container
docker cp mongodb:/tmp/backup $BACKUP_DIR/$BACKUP_NAME

# Create tarball
cd $BACKUP_DIR
tar -czf ${BACKUP_NAME}.tar.gz $BACKUP_NAME/
rm -rf $BACKUP_NAME/

# Keep only last 7 days
find $BACKUP_DIR -name "cafe_pos_*.tar.gz" -mtime +7 -delete

# Cleanup container
docker exec mongodb rm -rf /tmp/backup

echo "Backup complete: ${BACKUP_NAME}.tar.gz"
```

Make it executable:
```bash
chmod +x /home/ubuntu/backup-mongodb.sh
```

## 🔄 Common Scenarios

### Scenario 1: Clone Production to Local for Testing

```bash
# 1. Backup from EC2
EC2_HOST=your-ip ./scripts/backup-db-from-ec2.sh

# 2. Restore to local
./scripts/restore-db-to-local.sh ./backups/cafe_pos_backup_20260205_150000

# 3. Test locally
npm run dev
```

### Scenario 2: Restore Specific Collection

```bash
# Restore only users collection
docker exec mongodb mongorestore \
  --db cafe_pos \
  --collection users \
  --drop \
  /tmp/cafe_pos/users.bson
```

### Scenario 3: Backup Before Risky Operation

```bash
# Quick backup
docker exec mongodb mongodump \
  --db cafe_pos \
  --out /tmp/backup_before_migration

# If something goes wrong, restore
docker exec mongodb mongorestore \
  --db cafe_pos \
  --drop \
  /tmp/backup_before_migration/cafe_pos
```

## 📊 Backup Size Estimation

Typical database sizes:
- **Empty DB**: ~1 MB
- **Development**: 10-50 MB
- **Production (1 month)**: 100-500 MB
- **Production (1 year)**: 1-5 GB

Compressed backup is typically 10-20% of original size.

## ⚠️ Important Notes

### Security
- ✅ Always use SSH keys, not passwords
- ✅ Store backups in secure location
- ✅ Encrypt sensitive backups
- ✅ Limit backup access to authorized users

### Performance
- ⚠️ Backup during low-traffic hours
- ⚠️ Large backups may take several minutes
- ⚠️ Restore will drop existing data

### Data Integrity
- ✅ Verify backup after creation
- ✅ Test restore periodically
- ✅ Keep multiple backup versions
- ✅ Document backup procedures

## 🆘 Troubleshooting

### Error: "Connection refused"
```bash
# Check if MongoDB container is running
docker ps | grep mongodb

# Check MongoDB logs
docker logs mongodb
```

### Error: "Permission denied"
```bash
# Fix SSH key permissions
chmod 600 ~/.ssh/your-key.pem

# Fix backup directory permissions
sudo chown -R $USER:$USER ./backups
```

### Error: "No space left on device"
```bash
# Check disk space on EC2
ssh -i ~/.ssh/key.pem ubuntu@your-ip 'df -h'

# Cleanup old backups
ssh -i ~/.ssh/key.pem ubuntu@your-ip 'rm -f ~/backups/cafe_pos_*.tar.gz'
```

### Error: "Database already exists"
```bash
# Use --drop flag to replace existing database
mongorestore --db cafe_pos --drop ./backups/backup/cafe_pos
```

## 📚 Related Scripts

- `scripts/backup-db-from-ec2.sh` - Backup from EC2 to local
- `scripts/restore-db-to-local.sh` - Restore backup to local
- `scripts/restore-db-to-ec2.sh` - Restore backup to EC2 (TODO)

## 🔗 References

- [MongoDB Backup Methods](https://www.mongodb.com/docs/manual/core/backups/)
- [mongodump Documentation](https://www.mongodb.com/docs/database-tools/mongodump/)
- [mongorestore Documentation](https://www.mongodb.com/docs/database-tools/mongorestore/)

---

**Last Updated:** February 5, 2026  
**Maintained By:** Development Team
