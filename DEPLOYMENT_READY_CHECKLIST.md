# ✅ Café POS Deployment Ready Checklist

**Status**: Ready for EC2 Deployment  
**Date**: January 31, 2026  
**Version**: 1.0.0

---

## 📋 Current State

### ✅ Completed Tasks

1. **Security Hardening**
   - ✅ Removed hardcoded credentials from frontend
   - ✅ Removed test user accounts from backend seed
   - ✅ Created `.env.example` template
   - ✅ Added `.gitignore` to prevent committing `.env`
   - ✅ Generated secure passwords (MongoDB & JWT)

2. **Port Configuration**
   - ✅ Frontend: Port 80 (via Nginx)
   - ✅ Backend: Port 3000
   - ✅ MongoDB: Port 27017 (internal only)
   - ✅ Updated `frontend/src/services/api.js` to use port 3000
   - ✅ Updated `frontend/nginx.conf` to proxy to port 3000

3. **Docker Images**
   - ✅ Backend image: `linhtranphu/cafe-pos-backend:v1.0.1` (latest)
   - ✅ Frontend image: `linhtranphu/cafe-pos-frontend:v1.0.1` (latest)
   - ✅ Both images pushed to Docker Hub
   - ✅ Images use port 3000 for backend

4. **Deployment Scripts**
   - ✅ `ec2-deploy-from-github.sh` - Deploy from GitHub on EC2
   - ✅ `deploy-to-ec2.sh` - Local script to prepare and copy files
   - ✅ `docker-compose.hub.yml` - Production Docker Compose config
   - ✅ `.env` - Secure environment variables with generated passwords

5. **Documentation**
   - ✅ `DEPLOY_TO_EC2_WITH_DOCKER_HUB.md` - Comprehensive guide
   - ✅ `EC2_QUICK_START.md` - 5-minute quick start
   - ✅ `MANUAL_EC2_DEPLOYMENT.md` - Manual deployment steps
   - ✅ `EC2_DEPLOYMENT_SUMMARY.md` - Summary with checklist

---

## 🚀 Two Deployment Options

### Option 1: Deploy from GitHub (Recommended)

**Best for**: Fresh EC2 instance, no local files needed

**Steps**:
1. SSH to EC2
2. Run `ec2-deploy-from-github.sh` (clones from GitHub)
3. Script handles everything: Docker, Docker Compose, deployment

**Advantages**:
- No need to copy files locally
- Always gets latest code from GitHub
- Automatic setup and deployment

**Command**:
```bash
# On EC2
curl -fsSL https://raw.githubusercontent.com/linhtranphu/Cafe-POS/main/ec2-deploy-from-github.sh | bash
```

---

### Option 2: Deploy from Local Files

**Best for**: Testing locally first, then deploying

**Steps**:
1. Run `deploy-to-ec2.sh` locally (builds, pushes images, copies files)
2. SSH to EC2 and run `ec2-deploy.sh`
3. Edit `.env` with secure passwords

**Advantages**:
- Full control over build process
- Can test locally first
- Verify images before deployment

**Commands**:
```bash
# Locally
./deploy-to-ec2.sh

# On EC2
cd ~/cafe-pos
nano .env  # Edit with secure passwords
./ec2-deploy.sh
```

---

## 🔐 Security Configuration

### MongoDB Password
- **Location**: `.env` file
- **Variable**: `MONGO_INITDB_ROOT_PASSWORD`
- **Generation**: `openssl rand -base64 32`
- **Min Length**: 16 characters
- **Current**: Empty (user must set) ⚠️

### JWT Secret
- **Location**: `.env` file
- **Variable**: `JWT_SECRET`
- **Generation**: `openssl rand -base64 64`
- **Min Length**: 32 characters
- **Current**: Empty (user must set) ⚠️

### Admin Account
- **Username**: `admin`
- **Password**: `admin123` (default)
- **Action**: Change immediately after first login

---

## 📊 Configuration Summary

| Component | Port | Status |
|-----------|------|--------|
| Frontend (Nginx) | 80 | ✅ Configured |
| Backend (Go) | 3000 | ✅ Configured |
| MongoDB | 27017 | ✅ Internal only |
| Docker Hub | - | ✅ Images pushed |
| GitHub | - | ✅ Repository ready |

---

## 🌐 Access URLs

### After Deployment

```
Frontend: http://13.229.74.162
Backend:  http://13.229.74.162:3000
```

### Login Credentials

```
Username: admin
Password: admin123
```

⚠️ **Change password immediately after first login!**

---

## 📁 Key Files

### Deployment Scripts
- `ec2-deploy-from-github.sh` - Deploy from GitHub (recommended)
- `deploy-to-ec2.sh` - Local deployment script
- `ec2-deploy.sh` - EC2 deployment script

### Configuration
- `docker-compose.hub.yml` - Production Docker Compose
- `.env` - Environment variables (secure)
- `.env.example` - Template for reference

### Documentation
- `DEPLOY_TO_EC2_WITH_DOCKER_HUB.md` - Full guide
- `EC2_QUICK_START.md` - Quick start
- `MANUAL_EC2_DEPLOYMENT.md` - Manual steps
- `EC2_DEPLOYMENT_SUMMARY.md` - Summary

---

## ✅ Pre-Deployment Checklist

- [ ] Docker images built and pushed to Docker Hub
- [ ] `.env` file has secure MongoDB password
- [ ] `.env` file has secure JWT secret
- [ ] `docker-compose.hub.yml` configured correctly
- [ ] Frontend API endpoint set to port 3000
- [ ] Backend port set to 3000
- [ ] GitHub repository is public and accessible
- [ ] EC2 instance is running and accessible
- [ ] EC2 security group allows ports 80 and 3000

---

## 🚀 Quick Start (GitHub Deployment)

### On Local Machine

```bash
# Verify everything is ready
cat .env | grep MONGO_INITDB_ROOT_PASSWORD
cat .env | grep JWT_SECRET

# Verify images are on Docker Hub
docker images | grep cafe-pos
```

### On EC2

```bash
# SSH to EC2
ssh -i YOUR_KEY.pem ec2-user@13.229.74.162

# Deploy from GitHub
curl -fsSL https://raw.githubusercontent.com/linhtranphu/Cafe-POS/main/ec2-deploy-from-github.sh | bash

# Wait for deployment to complete
# Access application at http://13.229.74.162
```

---

## 🆘 Troubleshooting

### MongoDB Authentication Failed
```bash
# Check .env
cat .env | grep MONGO

# Check MongoDB logs
docker-compose -f docker-compose.hub.yml logs mongodb

# Restart MongoDB
docker-compose -f docker-compose.hub.yml restart mongodb
```

### Backend Not Responding
```bash
# Check backend logs
docker-compose -f docker-compose.hub.yml logs backend

# Check if backend is running
docker ps | grep backend

# Restart backend
docker-compose -f docker-compose.hub.yml restart backend
```

### Frontend Not Loading
```bash
# Check frontend logs
docker-compose -f docker-compose.hub.yml logs frontend

# Check if frontend is running
docker ps | grep frontend

# Restart frontend
docker-compose -f docker-compose.hub.yml restart frontend
```

### Port Already in Use
```bash
# Check what's using the port
sudo lsof -i :80
sudo lsof -i :3000

# Kill the process if needed
sudo kill -9 <PID>
```

---

## 📞 Important Notes

1. **Never commit `.env` to git** - It contains sensitive credentials
2. **Change admin password immediately** after first login
3. **Keep MongoDB password secure** - Use strong, randomly generated passwords
4. **Regular backups recommended** - MongoDB data is in Docker volumes
5. **Monitor logs regularly** - Check for errors and issues

---

## 🎯 Next Steps

### Immediate (Today)
1. Verify `.env` has secure passwords ✅
2. Verify Docker images are on Docker Hub ✅
3. Verify EC2 instance is running ✅

### Deployment (When Ready)
1. SSH to EC2
2. Run `ec2-deploy-from-github.sh`
3. Wait for deployment to complete
4. Access application at `http://13.229.74.162`
5. Login with `admin` / `admin123`
6. Change admin password

### Post-Deployment
1. Test all features
2. Monitor logs for errors
3. Set up regular backups
4. Configure monitoring/alerts

---

## 📝 Version History

| Version | Date | Changes |
|---------|------|---------|
| 1.0.0 | Jan 31, 2026 | Initial deployment ready |

---

**Status**: ✅ Ready for Deployment  
**Last Updated**: January 31, 2026  
**Next Action**: Deploy to EC2

