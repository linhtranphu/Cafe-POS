# 🚀 Deploy to EC2 via Docker Hub with MongoDB Password

## 📋 Tổng quan quy trình

```
Local Machine (Build)
    ↓
    Build Docker images
    ↓
    Push to Docker Hub
    ↓
EC2 Server (Deploy)
    ↓
    Pull images from Docker Hub
    ↓
    Configure .env with MongoDB password
    ↓
    Start services with docker-compose
    ↓
    ✅ Application running with authentication
```

---

## 🔧 PHẦN 1: Chuẩn bị trên Local Machine

### Bước 1.1: Login Docker Hub

```bash
docker login
# Nhập Docker Hub username và password
```

### Bước 1.2: Build Backend Image

```bash
cd backend

# Build image
docker build -t your-dockerhub-username/cafe-pos-backend:v1.0.0 .
docker tag your-dockerhub-username/cafe-pos-backend:v1.0.0 \
           your-dockerhub-username/cafe-pos-backend:latest

# Push to Docker Hub
docker push your-dockerhub-username/cafe-pos-backend:v1.0.0
docker push your-dockerhub-username/cafe-pos-backend:latest

cd ..
```

### Bước 1.3: Build Frontend Image

```bash
cd frontend

# Build image
docker build -t your-dockerhub-username/cafe-pos-frontend:v1.0.0 .
docker tag your-dockerhub-username/cafe-pos-frontend:v1.0.0 \
           your-dockerhub-username/cafe-pos-frontend:latest

# Push to Docker Hub
docker push your-dockerhub-username/cafe-pos-frontend:v1.0.0
docker push your-dockerhub-username/cafe-pos-frontend:latest

cd ..
```

### Bước 1.4: Verify Images on Docker Hub

```bash
# Truy cập: https://hub.docker.com/u/your-dockerhub-username
# Kiểm tra 2 repositories:
# - cafe-pos-backend
# - cafe-pos-frontend
```

---

## 🖥️ PHẦN 2: Chuẩn bị EC2 Server

### Bước 2.1: SSH vào EC2

```bash
# Thay your-ec2-ip bằng IP của EC2
ssh -i your-key.pem ec2-user@your-ec2-ip

# Hoặc nếu dùng Ubuntu
ssh -i your-key.pem ubuntu@your-ec2-ip
```

### Bước 2.2: Cài Docker

```bash
# Cài Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Thêm user vào docker group
sudo usermod -aG docker $USER

# Logout và login lại để apply changes
exit
ssh -i your-key.pem ec2-user@your-ec2-ip
```

### Bước 2.3: Cài Docker Compose

```bash
# Download Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose

# Make executable
sudo chmod +x /usr/local/bin/docker-compose

# Verify
docker-compose --version
```

### Bước 2.4: Tạo thư mục deploy

```bash
# Tạo thư mục
mkdir -p ~/cafe-pos
cd ~/cafe-pos

# Verify
pwd
```

---

## 📁 PHẦN 3: Copy Files lên EC2

### Bước 3.1: Tạo docker-compose.hub.yml trên EC2

```bash
# Trên EC2, tạo file
cat > docker-compose.hub.yml << 'EOF'
version: '3.8'

services:
  mongodb:
    image: mongo:7.0
    container_name: cafe-pos-mongodb
    restart: always
    environment:
      MONGO_INITDB_ROOT_USERNAME: ${MONGO_INITDB_ROOT_USERNAME}
      MONGO_INITDB_ROOT_PASSWORD: ${MONGO_INITDB_ROOT_PASSWORD}
      MONGO_INITDB_DATABASE: ${MONGO_INITDB_DATABASE}
    ports:
      - "27017:27017"
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
    image: your-dockerhub-username/cafe-pos-backend:latest
    container_name: cafe-pos-backend
    restart: always
    environment:
      - MONGODB_URI=${MONGODB_URI}
      - MONGODB_DATABASE=${MONGODB_DATABASE}
      - JWT_SECRET=${JWT_SECRET}
      - PORT=8080
    ports:
      - "8080:8080"
    depends_on:
      mongodb:
        condition: service_healthy
    networks:
      - cafe-pos-network
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:8080/api/health"]
      interval: 30s
      timeout: 10s
      retries: 3

  frontend:
    image: your-dockerhub-username/cafe-pos-frontend:latest
    container_name: cafe-pos-frontend
    restart: always
    ports:
      - "80:80"
    depends_on:
      - backend
    networks:
      - cafe-pos-network
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost"]
      interval: 30s
      timeout: 10s
      retries: 3

volumes:
  mongodb_data:
    driver: local
  mongodb_config:
    driver: local

networks:
  cafe-pos-network:
    driver: bridge
EOF
```

**⚠️ QUAN TRỌNG**: Thay `your-dockerhub-username` bằng Docker Hub username của bạn!

### Bước 3.2: Tạo .env trên EC2

```bash
# Generate secure passwords
echo "MongoDB Password:"
MONGO_PASS=$(openssl rand -base64 32)
echo "$MONGO_PASS"

echo ""
echo "JWT Secret:"
JWT_SECRET=$(openssl rand -base64 64)
echo "$JWT_SECRET"

# Tạo file .env
cat > .env << EOF
# MongoDB Configuration
MONGO_INITDB_ROOT_USERNAME=admin
MONGO_INITDB_ROOT_PASSWORD=$MONGO_PASS
MONGO_INITDB_DATABASE=cafe_pos

# Backend Configuration
MONGODB_URI=mongodb://admin:$MONGO_PASS@mongodb:27017
MONGODB_DATABASE=cafe_pos
JWT_SECRET=$JWT_SECRET
PORT=8080
EOF

# Verify
cat .env
```

### Bước 3.3: Set permissions

```bash
# Secure .env file
chmod 600 .env

# Verify
ls -la .env
```

---

## 🚀 PHẦN 4: Deploy trên EC2

### Bước 4.1: Pull Images

```bash
# Pull images từ Docker Hub
docker-compose -f docker-compose.hub.yml pull

# Verify
docker images | grep cafe-pos
```

### Bước 4.2: Start Services

```bash
# Start services
docker-compose -f docker-compose.hub.yml up -d

# Check status
docker-compose -f docker-compose.hub.yml ps

# View logs
docker-compose -f docker-compose.hub.yml logs -f
```

### Bước 4.3: Wait for Services

```bash
# Đợi services khởi động (khoảng 30 giây)
sleep 30

# Check MongoDB
docker exec -it cafe-pos-mongodb mongosh \
  -u admin \
  -p $(grep MONGO_INITDB_ROOT_PASSWORD .env | cut -d= -f2) \
  --authenticationDatabase admin \
  --eval "db.adminCommand('ping')"
```

### Bước 4.4: Seed Data (Optional)

```bash
# Seed initial data
docker exec -it cafe-pos-backend ./cafe-pos-server seed

# Check logs
docker-compose -f docker-compose.hub.yml logs backend
```

### Bước 4.5: Verify Application

```bash
# Test backend health
curl http://localhost:8080/api/health

# Test frontend
curl http://localhost

# Check all services
docker-compose -f docker-compose.hub.yml ps
```

---

## 🌐 PHẦN 5: Access Application

### Bước 5.1: Get EC2 Public IP

```bash
# Trên EC2
curl http://169.254.169.254/latest/meta-data/public-ipv4

# Hoặc xem trong AWS Console
```

### Bước 5.2: Access Application

```
Frontend: http://your-ec2-public-ip
Backend:  http://your-ec2-public-ip:8080
```

### Bước 5.3: Login

```
Username: admin
Password: admin123

⚠️ CHANGE PASSWORD IMMEDIATELY!
```

---

## 🔐 PHẦN 6: Security Configuration

### Bước 6.1: Configure Security Group

```bash
# Trên AWS Console, edit Security Group:

Inbound Rules:
- Port 80 (HTTP):   0.0.0.0/0 (Allow from anywhere)
- Port 443 (HTTPS): 0.0.0.0/0 (Allow from anywhere)
- Port 8080:        0.0.0.0/0 (Optional, for API testing)
- Port 27017:       CLOSE (Never expose MongoDB!)

Outbound Rules:
- All traffic allowed
```

### Bước 6.2: Setup SSL/TLS (Recommended)

```bash
# Install Certbot
sudo yum install certbot python3-certbot-nginx -y

# Get certificate
sudo certbot certonly --standalone -d your-domain.com

# Update nginx config (if using nginx reverse proxy)
```

### Bước 6.3: Firewall Rules

```bash
# Check firewall status
sudo systemctl status firewalld

# Allow ports
sudo firewall-cmd --permanent --add-port=80/tcp
sudo firewall-cmd --permanent --add-port=443/tcp
sudo firewall-cmd --reload
```

---

## 📊 PHẦN 7: Monitoring & Maintenance

### Bước 7.1: View Logs

```bash
# All services
docker-compose -f docker-compose.hub.yml logs -f

# Specific service
docker-compose -f docker-compose.hub.yml logs -f backend
docker-compose -f docker-compose.hub.yml logs -f mongodb
docker-compose -f docker-compose.hub.yml logs -f frontend
```

### Bước 7.2: Check Status

```bash
# Services status
docker-compose -f docker-compose.hub.yml ps

# Resource usage
docker stats

# Disk usage
df -h
```

### Bước 7.3: Backup MongoDB

```bash
# Backup
docker exec cafe-pos-mongodb mongodump \
  --username=admin \
  --password=$(grep MONGO_INITDB_ROOT_PASSWORD .env | cut -d= -f2) \
  --authenticationDatabase=admin \
  --db=cafe_pos \
  --out=/backup

# Copy to host
docker cp cafe-pos-mongodb:/backup ./mongodb-backup-$(date +%Y%m%d)

# Compress
tar -czf mongodb-backup-$(date +%Y%m%d).tar.gz mongodb-backup-$(date +%Y%m%d)
```

### Bước 7.4: Update Images

```bash
# Pull latest images
docker-compose -f docker-compose.hub.yml pull

# Restart with new images
docker-compose -f docker-compose.hub.yml up -d

# Check logs
docker-compose -f docker-compose.hub.yml logs -f
```

---

## 🆘 Troubleshooting

### Backend không kết nối MongoDB

```bash
# Check .env
cat .env | grep MONGODB_URI

# Check backend logs
docker-compose -f docker-compose.hub.yml logs backend | grep -i mongo

# Test MongoDB connection
docker exec -it cafe-pos-mongodb mongosh \
  -u admin \
  -p $(grep MONGO_INITDB_ROOT_PASSWORD .env | cut -d= -f2) \
  --authenticationDatabase admin
```

### Port already in use

```bash
# Check what's using port 80
sudo lsof -i :80

# Check what's using port 8080
sudo lsof -i :8080

# Kill process if needed
sudo kill -9 <PID>
```

### Images won't pull

```bash
# Check Docker Hub connection
docker pull hello-world

# Login again
docker login

# Pull specific image
docker pull your-dockerhub-username/cafe-pos-backend:latest
```

### MongoDB authentication failed

```bash
# Check MongoDB logs
docker-compose -f docker-compose.hub.yml logs mongodb

# Check credentials in .env
cat .env | grep MONGO_INITDB

# Restart MongoDB
docker-compose -f docker-compose.hub.yml restart mongodb
```

---

## 📋 Checklist

### Local Machine
- [ ] Build backend image
- [ ] Build frontend image
- [ ] Push to Docker Hub
- [ ] Verify images on Docker Hub

### EC2 Server
- [ ] SSH vào EC2
- [ ] Cài Docker
- [ ] Cài Docker Compose
- [ ] Tạo thư mục ~/cafe-pos
- [ ] Tạo docker-compose.hub.yml
- [ ] Tạo .env với secure passwords
- [ ] Set permissions: chmod 600 .env
- [ ] Pull images
- [ ] Start services
- [ ] Verify MongoDB authentication
- [ ] Seed data (optional)
- [ ] Test application
- [ ] Configure Security Group
- [ ] Setup SSL/TLS (recommended)

### Post-Deployment
- [ ] Change admin password
- [ ] Create additional users
- [ ] Setup backups
- [ ] Monitor logs
- [ ] Test all features

---

## 🎯 Quick Reference

### On Local Machine
```bash
# Build and push
docker build -t your-username/cafe-pos-backend:latest ./backend
docker push your-username/cafe-pos-backend:latest

docker build -t your-username/cafe-pos-frontend:latest ./frontend
docker push your-username/cafe-pos-frontend:latest
```

### On EC2 Server
```bash
# SSH
ssh -i your-key.pem ec2-user@your-ec2-ip

# Deploy
cd ~/cafe-pos
docker-compose -f docker-compose.hub.yml pull
docker-compose -f docker-compose.hub.yml up -d

# Verify
docker-compose -f docker-compose.hub.yml ps
curl http://localhost:8080/api/health
```

---

## 📞 Important Notes

1. **Replace placeholders:**
   - `your-dockerhub-username` → Your Docker Hub username
   - `your-ec2-ip` → Your EC2 public IP
   - `your-key.pem` → Your EC2 key pair file
   - `your-domain.com` → Your domain (if using SSL)

2. **Security:**
   - Never expose MongoDB port (27017)
   - Always use strong passwords
   - Enable SSL/TLS in production
   - Regular backups
   - Monitor logs

3. **Maintenance:**
   - Update images regularly
   - Backup MongoDB weekly
   - Monitor disk space
   - Check logs daily

---

**Version**: 1.0.0  
**Last Updated**: January 2026  
**Status**: Ready for Production Deployment
