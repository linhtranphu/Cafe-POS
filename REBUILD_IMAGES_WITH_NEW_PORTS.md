# 🔨 Rebuild Docker Images with New Ports

## 📋 Tóm tắt thay đổi

- **Frontend**: Port 80 (unchanged)
- **Backend**: Port 8080 → **Port 3000** ✅
- **MongoDB**: Port 27017 (unchanged)

---

## 🚀 Rebuild & Push Images

### Bước 1: Build Backend Image

```bash
cd backend

# Build image
docker build -t your-username/cafe-pos-backend:v1.0.1 .

# Tag as latest
docker tag your-username/cafe-pos-backend:v1.0.1 \
           your-username/cafe-pos-backend:latest

# Push to Docker Hub
docker push your-username/cafe-pos-backend:v1.0.1
docker push your-username/cafe-pos-backend:latest

cd ..
```

### Bước 2: Build Frontend Image

```bash
cd frontend

# Build image
docker build -t your-username/cafe-pos-frontend:v1.0.1 .

# Tag as latest
docker tag your-username/cafe-pos-frontend:v1.0.1 \
           your-username/cafe-pos-frontend:latest

# Push to Docker Hub
docker push your-username/cafe-pos-frontend:v1.0.1
docker push your-username/cafe-pos-frontend:latest

cd ..
```

### Bước 3: Verify Images

```bash
# Check local images
docker images | grep cafe-pos

# Check on Docker Hub
# https://hub.docker.com/u/your-username
```

---

## 🧪 Test Locally

### Bước 1: Start Services

```bash
# Pull latest images
docker-compose -f docker-compose.hub.yml pull

# Start services
docker-compose -f docker-compose.hub.yml up -d

# Wait for services
sleep 30

# Check status
docker-compose -f docker-compose.hub.yml ps
```

### Bước 2: Test Frontend

```bash
# Test frontend (port 80)
curl http://localhost

# Should return HTML content
```

### Bước 3: Test Backend

```bash
# Test backend health (port 3000)
curl http://localhost:3000/api/health

# Should return: {"status":"ok"}
```

### Bước 4: Test MongoDB

```bash
# Test MongoDB connection
docker exec -it cafe-pos-mongodb mongosh \
  -u admin \
  -p $(grep MONGO_INITDB_ROOT_PASSWORD .env | cut -d= -f2) \
  --authenticationDatabase admin

# In mongosh:
use cafe_pos
db.users.find()
exit
```

### Bước 5: View Logs

```bash
# All services
docker-compose -f docker-compose.hub.yml logs -f

# Specific service
docker-compose -f docker-compose.hub.yml logs -f backend
docker-compose -f docker-compose.hub.yml logs -f frontend
docker-compose -f docker-compose.hub.yml logs -f mongodb
```

### Bước 6: Test Application

```bash
# Open browser
http://localhost

# Login
Username: admin
Password: admin123

# Test features
- Create order
- View dashboard
- Check API calls in browser console
```

### Bước 7: Stop Services

```bash
docker-compose -f docker-compose.hub.yml down
```

---

## 📊 Port Verification

### Check Ports in Use

```bash
# Check port 80 (frontend)
sudo lsof -i :80

# Check port 3000 (backend)
sudo lsof -i :3000

# Check port 27017 (mongodb)
sudo lsof -i :27017
```

### Expected Output

```
PORT    STATE    SERVICE
80      LISTEN   nginx (frontend)
3000    LISTEN   backend (go app)
27017   LISTEN   mongodb
```

---

## 🚀 Deploy to EC2

### Bước 1: Update deploy-to-ec2.sh

```bash
# The script already has port 3000 configured
# Just run it:
./deploy-to-ec2.sh
```

### Bước 2: Copy Files to EC2

```bash
scp -i your-key.pem docker-compose.hub.yml.ec2 ec2-user@your-ec2-ip:~/cafe-pos/
scp -i your-key.pem .env.ec2 ec2-user@your-ec2-ip:~/cafe-pos/
scp -i your-key.pem ec2-deploy.sh ec2-user@your-ec2-ip:~/cafe-pos/
```

### Bước 3: Deploy on EC2

```bash
ssh -i your-key.pem ec2-user@your-ec2-ip

cd ~/cafe-pos
mv docker-compose.hub.yml.ec2 docker-compose.hub.yml
mv .env.ec2 .env

# Edit .env with secure MongoDB password
nano .env

# Deploy
chmod 600 .env
chmod +x ec2-deploy.sh
./ec2-deploy.sh
```

### Bước 4: Verify on EC2

```bash
# Check services
docker-compose -f docker-compose.hub.yml ps

# Test frontend
curl http://localhost

# Test backend
curl http://localhost:3000/api/health

# Get public IP
curl http://169.254.169.254/latest/meta-data/public-ipv4

# Access application
# http://your-ec2-public-ip (frontend)
# http://your-ec2-public-ip:3000 (backend API)
```

---

## 📋 Checklist

### Build & Push
- [ ] Backend image built
- [ ] Backend image pushed to Docker Hub
- [ ] Frontend image built
- [ ] Frontend image pushed to Docker Hub
- [ ] Verify images on Docker Hub

### Local Testing
- [ ] Services started
- [ ] Frontend accessible on port 80
- [ ] Backend accessible on port 3000
- [ ] MongoDB authentication working
- [ ] Application login working
- [ ] API calls working
- [ ] Services stopped

### EC2 Deployment
- [ ] Files copied to EC2
- [ ] .env configured with secure password
- [ ] Services deployed
- [ ] Frontend accessible
- [ ] Backend accessible
- [ ] Application working
- [ ] Backups configured

---

## 🔐 Security Checklist

- [ ] MongoDB password is strong (min 16 chars)
- [ ] JWT secret is strong (min 32 chars)
- [ ] .env file has correct permissions (600)
- [ ] .env file NOT committed to git
- [ ] MongoDB port (27017) NOT exposed
- [ ] Backend port (3000) can be restricted via firewall
- [ ] Frontend port (80) accessible
- [ ] SSL/TLS enabled (recommended)

---

## 🆘 Troubleshooting

### Port Already in Use

```bash
# Find what's using the port
sudo lsof -i :3000

# Kill the process
sudo kill -9 <PID>

# Or change port in docker-compose.yml
# ports:
#   - "3001:3000"
```

### Backend Not Responding

```bash
# Check logs
docker-compose -f docker-compose.hub.yml logs backend

# Check if backend is running
docker-compose -f docker-compose.hub.yml ps backend

# Restart backend
docker-compose -f docker-compose.hub.yml restart backend
```

### Frontend Can't Connect to Backend

```bash
# Check nginx logs
docker-compose -f docker-compose.hub.yml logs frontend

# Check backend is running on port 3000
curl http://localhost:3000/api/health

# Check nginx proxy config
docker exec cafe-pos-frontend cat /etc/nginx/conf.d/default.conf
```

### MongoDB Connection Failed

```bash
# Check MongoDB logs
docker-compose -f docker-compose.hub.yml logs mongodb

# Check credentials in .env
cat .env | grep MONGO_INITDB

# Test connection
docker exec -it cafe-pos-mongodb mongosh \
  -u admin \
  -p 'your-password' \
  --authenticationDatabase admin
```

---

## 📞 Important Notes

1. **Version Tagging**
   - Use v1.0.1 for this update
   - Keep v1.0.0 for rollback if needed

2. **Port Changes**
   - Frontend: 80 (unchanged)
   - Backend: 8080 → 3000 (changed)
   - MongoDB: 27017 (unchanged)

3. **Testing**
   - Always test locally first
   - Verify all ports are accessible
   - Check logs for errors

4. **Deployment**
   - Update EC2 files before deploying
   - Backup MongoDB before deploying
   - Monitor logs after deployment

---

## 🎯 Quick Commands

```bash
# Build and push
docker build -t your-username/cafe-pos-backend:v1.0.1 ./backend
docker push your-username/cafe-pos-backend:v1.0.1

docker build -t your-username/cafe-pos-frontend:v1.0.1 ./frontend
docker push your-username/cafe-pos-frontend:v1.0.1

# Test locally
docker-compose -f docker-compose.hub.yml up -d
curl http://localhost
curl http://localhost:3000/api/health
docker-compose -f docker-compose.hub.yml down

# Deploy to EC2
./deploy-to-ec2.sh
```

---

**Version**: 1.0.0  
**Last Updated**: January 2026  
**Status**: Ready to Build & Deploy
