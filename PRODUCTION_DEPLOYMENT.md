# 🚀 Production Deployment Guide

## ⚠️ Security Checklist

### 1. Environment Variables Setup

Create a `.env` file in the root directory (DO NOT commit this file):

```bash
# MongoDB Configuration
MONGO_INITDB_ROOT_USERNAME=your_secure_username
MONGO_INITDB_ROOT_PASSWORD=your_very_secure_password_min_16_chars
MONGO_INITDB_DATABASE=cafe_pos

# Backend Configuration
MONGODB_URI=mongodb://your_secure_username:your_very_secure_password_min_16_chars@mongodb:27017
MONGODB_DATABASE=cafe_pos
JWT_SECRET=your_jwt_secret_minimum_32_characters_random_string
PORT=8080
```

### 2. Generate Secure Credentials

```bash
# Generate secure MongoDB password (32 characters)
openssl rand -base64 32

# Generate secure JWT secret (64 characters)
openssl rand -base64 64
```

### 3. Initial Admin Account

After deployment, the system will create ONE admin account:
- **Username**: `admin`
- **Password**: `admin123`

**⚠️ CRITICAL: Change this password immediately after first login!**

Steps:
1. Login with `admin/admin123`
2. Go to User Management
3. Click "Reset Password" for admin user
4. Set a strong password (min 12 characters, mix of letters, numbers, symbols)

### 4. Create Additional Users

After securing the admin account:
1. Login as admin
2. Go to User Management
3. Create users for each role:
   - **Manager**: Full system access
   - **Waiter**: Order management
   - **Barista**: Drink preparation
   - **Cashier**: Payment processing

**Never share passwords between users!**

## 🐳 Deployment Options

### Option 1: Build from Source

```bash
# 1. Create .env file with secure credentials
cp .env.example .env
nano .env  # Edit with your secure values

# 2. Build and start services
docker-compose up -d --build

# 3. Check services are running
docker-compose ps

# 4. View logs
docker-compose logs -f
```

### Option 2: Deploy from Docker Hub

```bash
# 1. Create .env file with secure credentials
cp .env.example .env
nano .env  # Edit with your secure values

# 2. Pull and start services
docker-compose -f docker-compose.hub.yml up -d

# 3. Check services are running
docker-compose -f docker-compose.hub.yml ps
```

## 📊 Post-Deployment Steps

### 1. Seed Initial Data

```bash
# Enter backend container
docker exec -it cafe-pos-backend sh

# Run seed command
./cafe-pos-server seed

# Exit container
exit
```

This will create:
- ✅ Admin user (username: admin, password: admin123)
- ✅ Sample menu items
- ✅ Sample ingredients
- ✅ Sample facilities

### 2. Verify Services

```bash
# Check backend health
curl http://localhost:8080/api/health

# Check frontend
curl http://localhost:80

# Check MongoDB connection
docker exec -it cafe-pos-mongodb mongosh -u your_username -p your_password
```

### 3. Access Application

Open browser: `http://your-server-ip`

Login with:
- Username: `admin`
- Password: `admin123`

**⚠️ Change password immediately!**

## 🔒 Security Best Practices

### 1. Firewall Configuration

```bash
# Allow only necessary ports
sudo ufw allow 80/tcp    # Frontend
sudo ufw allow 443/tcp   # HTTPS (if using SSL)
sudo ufw enable
```

### 2. MongoDB Security

- ✅ Use strong username/password
- ✅ Keep MongoDB port (27017) closed to external access
- ✅ Enable MongoDB authentication
- ✅ Regular backups

### 3. JWT Secret

- ✅ Minimum 32 characters
- ✅ Random alphanumeric + symbols
- ✅ Never commit to git
- ✅ Rotate periodically

### 4. SSL/TLS (Recommended)

Use nginx or Caddy as reverse proxy with Let's Encrypt:

```bash
# Install Caddy
sudo apt install caddy

# Configure Caddy
sudo nano /etc/caddy/Caddyfile
```

Example Caddyfile:
```
your-domain.com {
    reverse_proxy localhost:80
}
```

## 🔄 Backup & Restore

### Backup MongoDB

```bash
# Create backup
docker exec cafe-pos-mongodb mongodump \
  --username=your_username \
  --password=your_password \
  --authenticationDatabase=admin \
  --db=cafe_pos \
  --out=/backup

# Copy backup to host
docker cp cafe-pos-mongodb:/backup ./mongodb-backup-$(date +%Y%m%d)
```

### Restore MongoDB

```bash
# Copy backup to container
docker cp ./mongodb-backup cafe-pos-mongodb:/backup

# Restore
docker exec cafe-pos-mongodb mongorestore \
  --username=your_username \
  --password=your_password \
  --authenticationDatabase=admin \
  --db=cafe_pos \
  /backup/cafe_pos
```

## 📝 Monitoring

### View Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f backend
docker-compose logs -f frontend
docker-compose logs -f mongodb
```

### Check Resource Usage

```bash
# Container stats
docker stats

# Disk usage
docker system df
```

## 🆘 Troubleshooting

### Backend won't start

```bash
# Check logs
docker-compose logs backend

# Common issues:
# - MongoDB not ready: Wait 30 seconds and restart
# - Wrong credentials: Check .env file
# - Port conflict: Change PORT in .env
```

### Frontend shows connection error

```bash
# Check backend is running
curl http://localhost:8080/api/health

# Check frontend nginx config
docker exec cafe-pos-frontend cat /etc/nginx/conf.d/default.conf
```

### Cannot login

```bash
# Reset admin password
docker exec -it cafe-pos-mongodb mongosh \
  -u your_username \
  -p your_password \
  --authenticationDatabase admin

use cafe_pos
db.users.find({username: "admin"})

# If admin doesn't exist, run seed again
docker exec -it cafe-pos-backend ./cafe-pos-server seed
```

## 📞 Support

For issues or questions:
1. Check logs: `docker-compose logs -f`
2. Verify .env configuration
3. Ensure all services are running: `docker-compose ps`
4. Check MongoDB connection
5. Review this guide

## 🔐 Security Reminders

- ✅ Change default admin password immediately
- ✅ Use strong, unique passwords for all users
- ✅ Keep JWT_SECRET secure and random
- ✅ Never commit .env file to git
- ✅ Regular backups of MongoDB
- ✅ Monitor logs for suspicious activity
- ✅ Keep Docker images updated
- ✅ Use HTTPS in production
- ✅ Restrict MongoDB port access
- ✅ Regular security audits

---

**Last Updated**: January 2026
**Version**: 1.0.0
