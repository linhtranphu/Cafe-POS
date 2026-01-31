# 📋 EC2 Deployment Summary

## ✅ Những gì đã được chuẩn bị

### 1. Documentation Files ✅
- `DEPLOY_TO_EC2_WITH_DOCKER_HUB.md` - Hướng dẫn chi tiết đầy đủ
- `EC2_QUICK_START.md` - Quick start guide (5 phút)
- `EC2_DEPLOYMENT_SUMMARY.md` - File này

### 2. Scripts ✅
- `deploy-to-ec2.sh` - Tự động build, push, và tạo files cho EC2
- `ec2-deploy.sh` - Tự động deploy trên EC2 (được tạo bởi deploy-to-ec2.sh)

### 3. Configuration Files ✅
- `docker-compose.hub.yml.ec2` - Docker compose cho EC2 (được tạo bởi deploy-to-ec2.sh)
- `.env.ec2` - Environment template (được tạo bởi deploy-to-ec2.sh)

---

## 🚀 Quy trình Deploy (3 bước chính)

### Bước 1: Build & Push Images (Local Machine)

```bash
./deploy-to-ec2.sh
```

**Script sẽ:**
- ✅ Yêu cầu Docker Hub username
- ✅ Build backend image
- ✅ Build frontend image
- ✅ Push to Docker Hub
- ✅ Tạo files cho EC2:
  - docker-compose.hub.yml.ec2
  - .env.ec2
  - ec2-deploy.sh

**Thời gian:** ~5-10 phút

### Bước 2: Copy Files to EC2

```bash
# Thay your-key.pem và your-ec2-ip
scp -i your-key.pem docker-compose.hub.yml.ec2 ec2-user@your-ec2-ip:~/cafe-pos/
scp -i your-key.pem .env.ec2 ec2-user@your-ec2-ip:~/cafe-pos/
scp -i your-key.pem ec2-deploy.sh ec2-user@your-ec2-ip:~/cafe-pos/
```

**Thời gian:** ~1 phút

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

# Deploy
chmod 600 .env
chmod +x ec2-deploy.sh
./ec2-deploy.sh
```

**Thời gian:** ~5 phút

---

## 📊 MongoDB Password Setup

### Tạo Secure Password

```bash
# On EC2
openssl rand -base64 32

# Example:
# v0ys4k/cduilPyonooIa23XgBWtNv+f6qEmYyNEIzfM=
```

### Update .env

```bash
# Edit .env
nano .env

# Change:
MONGO_INITDB_ROOT_PASSWORD=v0ys4k/cduilPyonooIa23XgBWtNv+f6qEmYyNEIzfM=
MONGODB_URI=mongodb://admin:v0ys4k/cduilPyonooIa23XgBWtNv+f6qEmYyNEIzfM=@mongodb:27017
JWT_SECRET=<generated-secret>

# Save (Ctrl+X, Y, Enter)
```

### Verify

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

## 🎯 Credentials

### MongoDB
- **Username**: admin
- **Password**: (from .env - MONGO_INITDB_ROOT_PASSWORD)
- **Connection**: mongodb://admin:password@mongodb:27017

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

## 📋 Checklist

### Local Machine
- [ ] Docker installed and running
- [ ] Run: `./deploy-to-ec2.sh`
- [ ] Verify images on Docker Hub

### EC2 Server
- [ ] SSH to EC2
- [ ] Create ~/cafe-pos directory
- [ ] Copy files from local machine
- [ ] Rename files (remove .ec2 extension)
- [ ] Edit .env with secure MongoDB password
- [ ] Set permissions: `chmod 600 .env`
- [ ] Run: `./ec2-deploy.sh`
- [ ] Verify services: `docker-compose ps`
- [ ] Test backend: `curl http://localhost:8080/api/health`

### Post-Deployment
- [ ] Access application: http://your-ec2-ip
- [ ] Login with admin/admin123
- [ ] Change admin password
- [ ] Create additional users
- [ ] Configure Security Group (firewall)
- [ ] Setup SSL/TLS (recommended)
- [ ] Setup backups
- [ ] Monitor logs

---

## 🛠️ Useful Commands

### View Logs

```bash
docker-compose -f docker-compose.hub.yml logs -f
docker-compose -f docker-compose.hub.yml logs -f backend
docker-compose -f docker-compose.hub.yml logs -f mongodb
```

### Check Status

```bash
docker-compose -f docker-compose.hub.yml ps
docker stats
df -h
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
   - Always use strong passwords (min 16 chars)
   - Keep .env file secure (chmod 600)
   - Enable SSL/TLS in production
   - Regular backups (weekly)
   - Monitor logs daily

3. **Maintenance:**
   - Update images regularly
   - Backup MongoDB weekly
   - Check disk space monthly
   - Review logs for errors

---

## 🆘 Troubleshooting

### Docker not found

```bash
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER
```

### Images won't pull

```bash
docker login
docker pull your-username/cafe-pos-backend:latest
```

### MongoDB authentication failed

```bash
cat .env | grep MONGO_INITDB
docker-compose -f docker-compose.hub.yml logs mongodb
docker-compose -f docker-compose.hub.yml restart mongodb
```

### Port already in use

```bash
sudo lsof -i :80
sudo lsof -i :8080
sudo lsof -i :27017
```

---

## 📞 Files Reference

| File | Purpose |
|------|---------|
| `deploy-to-ec2.sh` | Build & push images, create EC2 files |
| `DEPLOY_TO_EC2_WITH_DOCKER_HUB.md` | Detailed guide |
| `EC2_QUICK_START.md` | Quick start (5 min) |
| `EC2_DEPLOYMENT_SUMMARY.md` | This file |

---

## 🎯 Next Steps

1. **Build & Push**: `./deploy-to-ec2.sh`
2. **Copy to EC2**: `scp ...`
3. **Configure**: `nano .env`
4. **Deploy**: `./ec2-deploy.sh`
5. **Access**: `http://your-ec2-ip`
6. **Secure**: Change passwords, setup SSL
7. **Monitor**: Check logs regularly

---

## 📊 Timeline

| Step | Time | Action |
|------|------|--------|
| 1 | 5-10 min | Build & push images |
| 2 | 1 min | Copy files to EC2 |
| 3 | 5 min | Deploy on EC2 |
| 4 | 2 min | Verify & test |
| **Total** | **~15 min** | **Complete deployment** |

---

**Version**: 1.0.0  
**Last Updated**: January 2026  
**Status**: ✅ Ready for Production Deployment

**Start with**: `./deploy-to-ec2.sh` 🚀
