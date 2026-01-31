# 🚀 Deployment Quick Reference

**Quick commands and URLs for Café POS deployment**

---

## 🎯 One-Command Deployment (Recommended)

### Deploy from GitHub to EC2

```bash
# SSH to EC2 first
ssh -i YOUR_KEY.pem ec2-user@13.229.74.162

# Then run this command
curl -fsSL https://raw.githubusercontent.com/linhtranphu/Cafe-POS/main/ec2-deploy-from-github.sh | bash
```

**That's it!** The script will:
- Install Docker & Docker Compose
- Clone repository from GitHub
- Generate secure passwords
- Deploy all services
- Seed initial data

---

## 🌐 Access Application

After deployment completes:

```
Frontend: http://13.229.74.162
Backend:  http://13.229.74.162:3000
```

### Login
```
Username: admin
Password: admin123
```

⚠️ **Change password immediately!**

---

## 📋 Manual Deployment (If Needed)

### Step 1: SSH to EC2
```bash
ssh -i YOUR_KEY.pem ec2-user@13.229.74.162
```

### Step 2: Clone Repository
```bash
git clone https://github.com/linhtranphu/Cafe-POS.git cafe-pos
cd cafe-pos
```

### Step 3: Create .env File
```bash
cat > .env << EOF
MONGO_INITDB_ROOT_USERNAME=admin
MONGO_INITDB_ROOT_PASSWORD=$(openssl rand -base64 32)
MONGO_INITDB_DATABASE=cafe_pos
MONGODB_URI=mongodb://admin:$(openssl rand -base64 32)@mongodb:27017
MONGODB_DATABASE=cafe_pos
JWT_SECRET=$(openssl rand -base64 64)
PORT=3000
EOF
```

**Or edit manually:**
```bash
nano .env
# Set your own passwords:
# MONGO_INITDB_ROOT_PASSWORD=your-secure-password
# MONGODB_URI=mongodb://admin:your-secure-password@mongodb:27017
# JWT_SECRET=your-secure-jwt-secret
```

### Step 4: Install Docker
```bash
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER
```

### Step 5: Install Docker Compose
```bash
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

### Step 6: Deploy
```bash
docker-compose -f docker-compose.hub.yml up -d
sleep 30
docker exec -it cafe-pos-backend ./cafe-pos-server seed
```

---

## 🔧 Common Commands

### View Logs
```bash
# All services
docker-compose -f docker-compose.hub.yml logs -f

# Specific service
docker-compose -f docker-compose.hub.yml logs -f backend
docker-compose -f docker-compose.hub.yml logs -f frontend
docker-compose -f docker-compose.hub.yml logs -f mongodb
```

### Restart Services
```bash
# All services
docker-compose -f docker-compose.hub.yml restart

# Specific service
docker-compose -f docker-compose.hub.yml restart backend
docker-compose -f docker-compose.hub.yml restart frontend
docker-compose -f docker-compose.hub.yml restart mongodb
```

### Stop Services
```bash
docker-compose -f docker-compose.hub.yml down
```

### Start Services
```bash
docker-compose -f docker-compose.hub.yml up -d
```

### Check Status
```bash
docker-compose -f docker-compose.hub.yml ps
```

---

## 🧪 Verify Deployment

### Test Frontend
```bash
curl http://localhost
```

### Test Backend
```bash
curl http://localhost:3000/api/health
```

### Test MongoDB
```bash
docker exec cafe-pos-mongodb mongosh \
  -u admin \
  -p YOUR_PASSWORD \
  --authenticationDatabase admin \
  --eval "db.adminCommand('ping')"
```

---

## 🔐 Security

### Change Admin Password
1. Login to application
2. Go to Profile
3. Change password
4. Use strong password (min 12 chars)

### Backup MongoDB
```bash
docker exec cafe-pos-mongodb mongodump \
  -u admin \
  -p YOUR_PASSWORD \
  --authenticationDatabase admin \
  --out /backup
```

### View Environment Variables
```bash
cat .env
```

### Update Environment Variables
```bash
nano .env
docker-compose -f docker-compose.hub.yml restart
```

---

## 📊 Configuration

| Item | Value |
|------|-------|
| Frontend Port | 80 |
| Backend Port | 3000 |
| MongoDB Port | 27017 (internal) |
| Docker Hub | linhtranphu |
| GitHub Repo | linhtranphu/Cafe-POS |
| EC2 IP | 13.229.74.162 |
| EC2 User | ec2-user |

---

## 🆘 Troubleshooting

### Services Won't Start
```bash
# Check logs
docker-compose -f docker-compose.hub.yml logs

# Restart all services
docker-compose -f docker-compose.hub.yml restart

# Rebuild and restart
docker-compose -f docker-compose.hub.yml pull
docker-compose -f docker-compose.hub.yml up -d
```

### MongoDB Connection Failed
```bash
# Check MongoDB is running
docker ps | grep mongodb

# Check MongoDB logs
docker-compose -f docker-compose.hub.yml logs mongodb

# Verify credentials in .env
cat .env | grep MONGO
```

### Frontend Shows Blank Page
```bash
# Check frontend logs
docker-compose -f docker-compose.hub.yml logs frontend

# Check backend is running
curl http://localhost:3000/api/health

# Restart frontend
docker-compose -f docker-compose.hub.yml restart frontend
```

### Port Already in Use
```bash
# Find what's using the port
sudo lsof -i :80
sudo lsof -i :3000

# Kill the process
sudo kill -9 <PID>

# Restart services
docker-compose -f docker-compose.hub.yml restart
```

---

## 📞 Support

### Documentation
- Full Guide: `DEPLOY_TO_EC2_WITH_DOCKER_HUB.md`
- Quick Start: `EC2_QUICK_START.md`
- Manual Steps: `MANUAL_EC2_DEPLOYMENT.md`
- Checklist: `DEPLOYMENT_READY_CHECKLIST.md`

### GitHub
- Repository: https://github.com/linhtranphu/Cafe-POS
- Issues: https://github.com/linhtranphu/Cafe-POS/issues

---

## ⚡ Quick Links

- **Frontend**: http://13.229.74.162
- **Backend**: http://13.229.74.162:3000
- **Docker Hub**: https://hub.docker.com/u/linhtranphu
- **GitHub**: https://github.com/linhtranphu/Cafe-POS

---

**Last Updated**: January 31, 2026  
**Status**: Ready for Deployment

