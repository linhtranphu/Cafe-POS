# 🚀 Deploy Café POS lên AWS EC2

## 📋 Yêu Cầu

- AWS Account
- EC2 instance (Ubuntu 22.04 LTS)
- Domain name (optional, cho HTTPS)

## 🖥️ Bước 1: Tạo EC2 Instance

### 1.1 Launch EC2 Instance

1. Đăng nhập AWS Console → EC2 Dashboard
2. Click "Launch Instance"
3. Cấu hình:
   - **Name**: cafe-pos-server
   - **AMI**: Ubuntu Server 22.04 LTS
   - **Instance type**: t3.medium (2 vCPU, 4GB RAM) hoặc t3.small (2 vCPU, 2GB RAM)
   - **Key pair**: Tạo mới hoặc chọn existing key pair
   - **Network settings**:
     - Allow SSH (port 22)
     - Allow HTTP (port 80)
     - Allow HTTPS (port 443)
     - Allow Custom TCP (port 8080) - Backend API
   - **Storage**: 20GB gp3

4. Click "Launch Instance"

### 1.2 Security Group Rules

```
Type            Protocol    Port Range    Source
SSH             TCP         22            Your IP
HTTP            TCP         80            0.0.0.0/0
HTTPS           TCP         443           0.0.0.0/0
Custom TCP      TCP         8080          0.0.0.0/0
Custom TCP      TCP         27017         127.0.0.1/32 (localhost only)
```

## 🔧 Bước 2: Kết Nối và Cài Đặt

### 2.1 SSH vào EC2

```bash
# Thay YOUR-KEY.pem và EC2-PUBLIC-IP
chmod 400 YOUR-KEY.pem
ssh -i YOUR-KEY.pem ubuntu@EC2-PUBLIC-IP
```

### 2.2 Update System

```bash
sudo apt update && sudo apt upgrade -y
```

### 2.3 Cài Đặt Docker

```bash
# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Add user to docker group
sudo usermod -aG docker ubuntu

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Verify installation
docker --version
docker-compose --version

# Logout and login again for group changes
exit
```

### 2.4 Cài Đặt Git

```bash
sudo apt install git -y
```

## 📦 Bước 3: Deploy Application

### 3.1 Clone Repository

```bash
# SSH vào EC2 lại
ssh -i YOUR-KEY.pem ubuntu@EC2-PUBLIC-IP

# Clone project (thay YOUR-REPO-URL)
git clone YOUR-REPO-URL cafe-pos
cd cafe-pos
```

### 3.2 Cấu Hình Environment

```bash
# Tạo file .env cho production
cat > .env << 'EOF'
# MongoDB
MONGO_INITDB_ROOT_USERNAME=admin
MONGO_INITDB_ROOT_PASSWORD=CHANGE_THIS_PASSWORD_123

# Backend
JWT_SECRET=CHANGE_THIS_SECRET_KEY_IN_PRODUCTION
MONGODB_URI=mongodb://admin:CHANGE_THIS_PASSWORD_123@mongodb:27017
MONGODB_DATABASE=cafe_pos
PORT=8080

# Frontend
VITE_API_URL=http://YOUR-EC2-PUBLIC-IP:8080
EOF

# Thay đổi permissions
chmod 600 .env
```

### 3.3 Update docker-compose.yml cho Production

```bash
# Backup original
cp docker-compose.yml docker-compose.yml.backup

# Update docker-compose.yml
cat > docker-compose.yml << 'EOF'
version: '3.8'

services:
  mongodb:
    image: mongo:7.0
    container_name: cafe-pos-mongodb
    restart: always
    environment:
      MONGO_INITDB_ROOT_USERNAME: ${MONGO_INITDB_ROOT_USERNAME}
      MONGO_INITDB_ROOT_PASSWORD: ${MONGO_INITDB_ROOT_PASSWORD}
      MONGO_INITDB_DATABASE: ${MONGODB_DATABASE}
    volumes:
      - mongodb_data:/data/db
      - mongodb_config:/data/configdb
    networks:
      - cafe-pos-network
    healthcheck:
      test: echo 'db.runCommand("ping").ok' | mongosh localhost:27017/test --quiet
      interval: 10s
      timeout: 5s
      retries: 5

  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile
    container_name: cafe-pos-backend
    restart: always
    environment:
      - MONGODB_URI=${MONGODB_URI}
      - MONGODB_DATABASE=${MONGODB_DATABASE}
      - JWT_SECRET=${JWT_SECRET}
      - PORT=${PORT}
    ports:
      - "8080:8080"
    depends_on:
      mongodb:
        condition: service_healthy
    networks:
      - cafe-pos-network

  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile
      args:
        - VITE_API_URL=${VITE_API_URL}
    container_name: cafe-pos-frontend
    restart: always
    ports:
      - "80:80"
    depends_on:
      - backend
    networks:
      - cafe-pos-network

volumes:
  mongodb_data:
  mongodb_config:

networks:
  cafe-pos-network:
    driver: bridge
EOF
```

### 3.4 Build và Deploy

```bash
# Build images
docker-compose build

# Start services
docker-compose up -d

# Check status
docker-compose ps

# View logs
docker-compose logs -f
```

## 🔍 Bước 4: Verify Deployment

### 4.1 Check Services

```bash
# Check all containers are running
docker-compose ps

# Check backend health
curl http://localhost:8080/api/health

# Check frontend
curl http://localhost
```

### 4.2 Access Application

Mở browser và truy cập:
- **Frontend**: http://YOUR-EC2-PUBLIC-IP
- **Backend API**: http://YOUR-EC2-PUBLIC-IP:8080

## 🔐 Bước 5: Security Hardening (Recommended)

### 5.1 Setup Firewall

```bash
# Install UFW
sudo apt install ufw -y

# Configure firewall
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow ssh
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw allow 8080/tcp

# Enable firewall
sudo ufw enable
sudo ufw status
```

### 5.2 Setup SSL/HTTPS (Optional but Recommended)

```bash
# Install Certbot
sudo apt install certbot python3-certbot-nginx -y

# Install Nginx
sudo apt install nginx -y

# Configure Nginx as reverse proxy
sudo nano /etc/nginx/sites-available/cafe-pos

# Add this configuration:
```

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:80;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }

    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }
}
```

```bash
# Enable site
sudo ln -s /etc/nginx/sites-available/cafe-pos /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx

# Get SSL certificate
sudo certbot --nginx -d your-domain.com
```

## 📊 Bước 6: Monitoring & Maintenance

### 6.1 View Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f backend
docker-compose logs -f frontend
docker-compose logs -f mongodb
```

### 6.2 Restart Services

```bash
# Restart all
docker-compose restart

# Restart specific service
docker-compose restart backend
```

### 6.3 Update Application

```bash
# Pull latest code
git pull origin main

# Rebuild and restart
docker-compose down
docker-compose build
docker-compose up -d
```

### 6.4 Backup MongoDB

```bash
# Create backup script
cat > backup-mongodb.sh << 'EOF'
#!/bin/bash
BACKUP_DIR="/home/ubuntu/backups"
DATE=$(date +%Y%m%d_%H%M%S)
mkdir -p $BACKUP_DIR

docker exec cafe-pos-mongodb mongodump \
  --username admin \
  --password CHANGE_THIS_PASSWORD_123 \
  --authenticationDatabase admin \
  --out /tmp/backup

docker cp cafe-pos-mongodb:/tmp/backup $BACKUP_DIR/mongodb_$DATE
tar -czf $BACKUP_DIR/mongodb_$DATE.tar.gz -C $BACKUP_DIR mongodb_$DATE
rm -rf $BACKUP_DIR/mongodb_$DATE

# Keep only last 7 backups
ls -t $BACKUP_DIR/mongodb_*.tar.gz | tail -n +8 | xargs -r rm
EOF

chmod +x backup-mongodb.sh

# Setup cron job for daily backup
crontab -e
# Add this line:
# 0 2 * * * /home/ubuntu/cafe-pos/backup-mongodb.sh
```

## 🚨 Troubleshooting

### Issue: Container không start

```bash
# Check logs
docker-compose logs

# Check disk space
df -h

# Check memory
free -h
```

### Issue: Cannot connect to MongoDB

```bash
# Check MongoDB logs
docker-compose logs mongodb

# Restart MongoDB
docker-compose restart mongodb
```

### Issue: Port already in use

```bash
# Check what's using the port
sudo lsof -i :80
sudo lsof -i :8080

# Kill process if needed
sudo kill -9 PID
```

## 📈 Performance Optimization

### Enable Docker logging limits

```bash
# Edit docker-compose.yml, add to each service:
logging:
  driver: "json-file"
  options:
    max-size: "10m"
    max-file: "3"
```

### Setup Swap (if needed)

```bash
# Create 2GB swap
sudo fallocate -l 2G /swapfile
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile

# Make permanent
echo '/swapfile none swap sw 0 0' | sudo tee -a /etc/fstab
```

## 🎯 Quick Commands Reference

```bash
# Deploy
cd cafe-pos && docker-compose up -d

# Stop
docker-compose down

# Restart
docker-compose restart

# View logs
docker-compose logs -f

# Update
git pull && docker-compose down && docker-compose build && docker-compose up -d

# Backup
./backup-mongodb.sh

# Check status
docker-compose ps
```

## 📞 Support

**Access URLs:**
- Frontend: http://YOUR-EC2-PUBLIC-IP
- Backend: http://YOUR-EC2-PUBLIC-IP:8080
- MongoDB: localhost:27017 (internal only)

**Default Login:**
- Manager: `admin/admin123`
- Waiter: `waiter1/waiter123`
- Cashier: `cashier1/cashier123`

**Important:** Đổi tất cả passwords trong production!
