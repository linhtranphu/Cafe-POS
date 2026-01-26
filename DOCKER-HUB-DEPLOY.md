# 🐳 Deploy Café POS từ Docker Hub

## 🚀 Quick Deploy (1 Command)

### Trên EC2 hoặc bất kỳ server Linux nào:

```bash
curl -fsSL https://raw.githubusercontent.com/linhtranphu/Cafe-POS/main/deploy-from-hub.sh | bash
```

Hoặc:

```bash
wget -qO- https://raw.githubusercontent.com/linhtranphu/Cafe-POS/main/deploy-from-hub.sh | bash
```

## 📦 Docker Hub Images

- **Backend**: `linhtranphu/cafe-pos-backend:latest`
- **Frontend**: `linhtranphu/cafe-pos-frontend:latest`
- **MongoDB**: `mongo:7.0` (official)

## 🔧 Manual Deploy

### Bước 1: Cài đặt Docker & Docker Compose

```bash
# Install Docker
curl -fsSL https://get.docker.com | sh

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

### Bước 2: Tạo docker-compose.yml

```bash
mkdir -p ~/cafe-pos
cd ~/cafe-pos

cat > docker-compose.yml << 'EOF'
version: '3.8'

services:
  mongodb:
    image: mongo:7.0
    container_name: cafe-pos-mongodb
    restart: always
    environment:
      MONGO_INITDB_ROOT_USERNAME: admin
      MONGO_INITDB_ROOT_PASSWORD: admin123
      MONGO_INITDB_DATABASE: cafe_pos
    volumes:
      - mongodb_data:/data/db
    networks:
      - cafe-pos-network

  backend:
    image: linhtranphu/cafe-pos-backend:latest
    container_name: cafe-pos-backend
    restart: always
    environment:
      - MONGODB_URI=mongodb://admin:admin123@mongodb:27017
      - MONGODB_DATABASE=cafe_pos
      - JWT_SECRET=your-secret-key-change-in-production
      - PORT=8080
    ports:
      - "8080:8080"
    depends_on:
      - mongodb
    networks:
      - cafe-pos-network

  frontend:
    image: linhtranphu/cafe-pos-frontend:latest
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

networks:
  cafe-pos-network:
    driver: bridge
EOF
```

### Bước 3: Deploy

```bash
# Pull images
docker-compose pull

# Start services
docker-compose up -d

# Check status
docker-compose ps

# View logs
docker-compose logs -f
```

## 🔐 Security Configuration

### Tạo file .env

```bash
cat > .env << 'EOF'
MONGO_INITDB_ROOT_USERNAME=admin
MONGO_INITDB_ROOT_PASSWORD=YOUR_STRONG_PASSWORD_HERE
MONGODB_URI=mongodb://admin:YOUR_STRONG_PASSWORD_HERE@mongodb:27017
MONGODB_DATABASE=cafe_pos
JWT_SECRET=YOUR_JWT_SECRET_HERE
EOF

chmod 600 .env
```

### Update docker-compose.yml để sử dụng .env

```yaml
services:
  mongodb:
    environment:
      MONGO_INITDB_ROOT_USERNAME: ${MONGO_INITDB_ROOT_USERNAME}
      MONGO_INITDB_ROOT_PASSWORD: ${MONGO_INITDB_ROOT_PASSWORD}
  
  backend:
    environment:
      - MONGODB_URI=${MONGODB_URI}
      - JWT_SECRET=${JWT_SECRET}
```

## 📊 Management Commands

```bash
# View logs
docker-compose logs -f

# Restart services
docker-compose restart

# Stop services
docker-compose down

# Update to latest version
docker-compose pull
docker-compose up -d

# Backup MongoDB
docker exec cafe-pos-mongodb mongodump --out /tmp/backup
docker cp cafe-pos-mongodb:/tmp/backup ./backup

# Restore MongoDB
docker cp ./backup cafe-pos-mongodb:/tmp/backup
docker exec cafe-pos-mongodb mongorestore /tmp/backup
```

## 🌐 Access Application

- **Frontend**: http://YOUR-SERVER-IP
- **Backend API**: http://YOUR-SERVER-IP:8080

## 👥 Default Users

- **Manager**: `admin/admin123`
- **Waiter**: `waiter1/waiter123`
- **Cashier**: `cashier1/cashier123`

⚠️ **Change these passwords in production!**

## 🔄 Update Application

```bash
cd ~/cafe-pos

# Pull latest images
docker-compose pull

# Restart with new images
docker-compose up -d

# Check logs
docker-compose logs -f
```

## 🛠️ Troubleshooting

### Check container status
```bash
docker-compose ps
```

### View logs
```bash
docker-compose logs backend
docker-compose logs frontend
docker-compose logs mongodb
```

### Restart specific service
```bash
docker-compose restart backend
```

### Clean restart
```bash
docker-compose down
docker-compose up -d
```

## 📈 Monitoring

### Check resource usage
```bash
docker stats
```

### Check disk space
```bash
df -h
docker system df
```

### Clean up unused images
```bash
docker system prune -a
```

## 🔒 Production Checklist

- [ ] Change MongoDB password
- [ ] Change JWT_SECRET
- [ ] Setup firewall (UFW)
- [ ] Setup SSL/HTTPS
- [ ] Setup backup cron job
- [ ] Change default user passwords
- [ ] Monitor logs regularly
- [ ] Setup alerts

## 📞 Support

For issues or questions, check:
- Container logs: `docker-compose logs -f`
- Container status: `docker-compose ps`
- System resources: `docker stats`
