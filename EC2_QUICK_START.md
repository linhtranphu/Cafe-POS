# 🚀 EC2 Quick Start - Deploy via Docker Hub

## 📋 Tóm tắt quy trình

```
Local Machine
  ↓ Build & Push images
Docker Hub
  ↓ Pull images
EC2 Server
  ↓ Configure & Deploy
✅ Application Running
```

---

## ⚡ Quick Deploy (5 phút)

### Bước 1: Build & Push Images (Local Machine)

```bash
# Chạy script tự động
./deploy-to-ec2.sh

# Nhập Docker Hub username khi được hỏi
# Script sẽ:
# - Build backend image
# - Build frontend image
# - Push to Docker Hub
# - Tạo files cho EC2
```

### Bước 2: Copy Files to EC2

```bash
# Thay your-key.pem và your-ec2-ip
scp -i your-key.pem docker-compose.hub.yml.ec2 ec2-user@your-ec2-ip:~/cafe-pos/
scp -i your-key.pem .env.ec2 ec2-user@your-ec2-ip:~/cafe-pos/
scp -i your-key.pem ec2-deploy.sh ec2-user@your-ec2-ip:~/cafe-pos/
```

### Bước 3: Deploy on EC2

```bash
# SSH to EC2
ssh -i your-key.pem ec2-user@your-ec2-ip

# Setup
cd ~/cafe-pos
mv docker-compose.hub.yml.ec2 docker-compose.hub.yml
mv .env.ec2 .env

# Edit .env with secure MongoDB password
nano .env

# Set permissions
chmod 600 .env
chmod +x ec2-deploy.sh

# Deploy
./ec2-deploy.sh
```

### Bước 4: Access Application

```
Frontend: http://your-ec2-public-ip
Backend:  http://your-ec2-public-ip:8080

Login: admin / admin123
⚠️ Change password immediately!
```

---

## 📝 Detailed Steps

### Step 1: Prepare Local Machine

```bash
# Ensure Docker is running
docker --version
docker info

# Login Docker Hub
docker login
```

### Step 2: Build & Push

```bash
# Run automated script
./deploy-to-ec2.sh

# Or build manually:
# Backend
cd backend
docker build -t your-username/cafe-pos-backend:latest .
docker push your-username/cafe-pos-backend:latest
cd ..

# Frontend
cd frontend
docker build -t your-username/cafe-pos-frontend:latest .
docker push your-username/cafe-pos-frontend:latest
cd ..
```

### Step 3: Prepare EC2

```bash
# SSH to EC2
ssh -i your-key.pem ec2-user@your-ec2-ip

# Create directory
mkdir -p ~/cafe-pos
cd ~/cafe-pos
```

### Step 4: Copy Files

```bash
# From local machine
scp -i your-key.pem docker-compose.hub.yml.ec2 ec2-user@your-ec2-ip:~/cafe-pos/
scp -i your-key.pem .env.ec2 ec2-user@your-ec2-ip:~/cafe-pos/
scp -i your-key.pem ec2-deploy.sh ec2-user@your-ec2-ip:~/cafe-pos/
```

### Step 5: Configure on EC2

```bash
# SSH to EC2
ssh -i your-key.pem ec2-user@your-ec2-ip
cd ~/cafe-pos

# Rename files
mv docker-compose.hub.yml.ec2 docker-compose.hub.yml
mv .env.ec2 .env

# Generate secure passwords
echo "MongoDB Password:"
openssl rand -base64 32

echo ""
echo "JWT Secret:"
openssl rand -base64 64

# Edit .env
nano .env

# Update:
# MONGO_INITDB_ROOT_PASSWORD=<generated-password>
# MONGODB_URI=mongodb://admin:<generated-password>@mongodb:27017
# JWT_SECRET=<generated-secret>

# Save and exit (Ctrl+X, Y, Enter)

# Set permissions
chmod 600 .env
chmod +x ec2-deploy.sh
```

### Step 6: Deploy

```bash
# Run deployment script
./ec2-deploy.sh

# Script will:
# - Install Docker
# - Install Docker Compose
# - Pull images
# - Start services
# - Verify MongoDB
# - Seed data
```

### Step 7: Verify

```bash
# Check services
docker-compose -f docker-compose.hub.yml ps

# Test backend
curl http://localhost:8080/api/health

# View logs
docker-compose -f docker-compose.hub.yml logs -f
```

---

## 🔐 MongoDB Password Setup

### Generate Secure Password

```bash
# On EC2
openssl rand -base64 32

# Example output:
# v0ys4k/cduilPyonooIa23XgBWtNv+f6qEmYyNEIzfM=
```

### Update .env

```bash
# Edit .env
nano .env

# Change these lines:
MONGO_INITDB_ROOT_PASSWORD=v0ys4k/cduilPyonooIa23XgBWtNv+f6qEmYyNEIzfM=
MONGODB_URI=mongodb://admin:v0ys4k/cduilPyonooIa23XgBWtNv+f6qEmYyNEIzfM=@mongodb:27017

# Save (Ctrl+X, Y, Enter)
```

### Verify Authentication

```bash
# Test MongoDB connection
docker exec -it cafe-pos-mongodb mongosh \
  -u admin \
  -p 'v0ys4k/cduilPyonooIa23XgBWtNv+f6qEmYyNEIzfM=' \
  --authenticationDatabase admin

# In mongosh:
use cafe_pos
db.users.find()
exit
```

---

## 📊 Credentials

### MongoDB
- **Username**: admin
- **Password**: (from .env - MONGO_INITDB_ROOT_PASSWORD)

### Application
- **Username**: admin
- **Password**: admin123
- **⚠️ Change after first login!**

### JWT Secret
- (from .env - JWT_SECRET)

---

## 🌐 Access Application

### Get EC2 Public IP

```bash
# On EC2
curl http://169.254.169.254/latest/meta-data/public-ipv4

# Or check AWS Console
```

### URLs

```
Frontend: http://your-ec2-public-ip
Backend:  http://your-ec2-public-ip:8080
```

### Login

```
Username: admin
Password: admin123
```

---

## 🛠️ Useful Commands

### View Logs

```bash
# All services
docker-compose -f docker-compose.hub.yml logs -f

# Specific service
docker-compose -f docker-compose.hub.yml logs -f backend
docker-compose -f docker-compose.hub.yml logs -f mongodb
```

### Restart Services

```bash
docker-compose -f docker-compose.hub.yml restart
```

### Stop Services

```bash
docker-compose -f docker-compose.hub.yml down
```

### Update Images

```bash
docker-compose -f docker-compose.hub.yml pull
docker-compose -f docker-compose.hub.yml up -d
```

### Backup MongoDB

```bash
docker exec cafe-pos-mongodb mongodump \
  --username=admin \
  --password=$(grep MONGO_INITDB_ROOT_PASSWORD .env | cut -d= -f2) \
  --authenticationDatabase=admin \
  --db=cafe_pos \
  --out=/backup

docker cp cafe-pos-mongodb:/backup ./mongodb-backup-$(date +%Y%m%d)
```

---

## ⚠️ Important Notes

1. **Replace placeholders:**
   - `your-key.pem` → Your EC2 key pair file
   - `your-ec2-ip` → Your EC2 public IP
   - `your-username` → Your Docker Hub username

2. **Security:**
   - Never expose MongoDB port (27017) to internet
   - Always use strong passwords
   - Keep .env file secure (chmod 600)
   - Enable SSL/TLS in production
   - Regular backups

3. **Maintenance:**
   - Monitor logs regularly
   - Backup MongoDB weekly
   - Update images periodically
   - Check disk space

---

## 🆘 Troubleshooting

### Docker not found

```bash
# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER
```

### Docker Compose not found

```bash
# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

### Images won't pull

```bash
# Check Docker Hub connection
docker pull hello-world

# Login again
docker login

# Pull specific image
docker pull your-username/cafe-pos-backend:latest
```

### MongoDB authentication failed

```bash
# Check .env
cat .env | grep MONGO_INITDB

# Check MongoDB logs
docker-compose -f docker-compose.hub.yml logs mongodb

# Restart MongoDB
docker-compose -f docker-compose.hub.yml restart mongodb
```

### Port already in use

```bash
# Check what's using port 80
sudo lsof -i :80

# Check what's using port 8080
sudo lsof -i :8080
```

---

## 📞 Next Steps

1. ✅ Build & push images: `./deploy-to-ec2.sh`
2. ✅ Copy files to EC2: `scp ...`
3. ✅ Configure .env on EC2: `nano .env`
4. ✅ Deploy: `./ec2-deploy.sh`
5. ✅ Access application: `http://your-ec2-ip`
6. ✅ Change admin password
7. ✅ Setup backups
8. ✅ Monitor logs

---

**Version**: 1.0.0  
**Last Updated**: January 2026  
**Status**: Ready for Production
