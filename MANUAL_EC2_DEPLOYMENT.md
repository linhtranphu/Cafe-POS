# 🚀 Manual EC2 Deployment Guide

## ⚠️ SSH Key Issue

The PEM file `EC2_PEM/OngTaPOS.pem` doesn't match the EC2 instance key pair.

**Solution**: Deploy manually using the files that have been created.

---

## 📋 Files Ready for Deployment

The following files have been created and are ready to copy to EC2:

1. **docker-compose.hub.yml.ec2** - Docker Compose configuration
2. **.env.ec2** - Environment template
3. **ec2-deploy.sh** - Deployment script

---

## 🚀 Manual Deployment Steps

### Step 1: Copy Files to EC2

Use your EC2 key pair to copy files:

```bash
# Replace YOUR_KEY.pem with your actual EC2 key pair
scp -i YOUR_KEY.pem docker-compose.hub.yml.ec2 ec2-user@13.229.74.162:~/cafe-pos/
scp -i YOUR_KEY.pem .env.ec2 ec2-user@13.229.74.162:~/cafe-pos/
scp -i YOUR_KEY.pem ec2-deploy.sh ec2-user@13.229.74.162:~/cafe-pos/
```

### Step 2: SSH to EC2

```bash
ssh -i YOUR_KEY.pem ec2-user@13.229.74.162
```

### Step 3: Setup Files on EC2

```bash
cd ~/cafe-pos

# Rename files
mv docker-compose.hub.yml.ec2 docker-compose.hub.yml
mv .env.ec2 .env

# View current .env
cat .env
```

### Step 4: Generate Secure Passwords

```bash
# Generate MongoDB password (32 characters)
echo "MongoDB Password:"
openssl rand -base64 32

# Generate JWT secret (64 characters)
echo ""
echo "JWT Secret:"
openssl rand -base64 64
```

### Step 5: Edit .env File

```bash
nano .env
```

**Update these lines:**

```bash
# Change from:
MONGO_INITDB_ROOT_PASSWORD=CHANGE_THIS_TO_SECURE_PASSWORD
MONGODB_URI=mongodb://admin:CHANGE_THIS_TO_SECURE_PASSWORD@mongodb:27017
JWT_SECRET=CHANGE_THIS_TO_SECURE_JWT_SECRET

# To:
MONGO_INITDB_ROOT_PASSWORD=<your-generated-password>
MONGODB_URI=mongodb://admin:<your-generated-password>@mongodb:27017
JWT_SECRET=<your-generated-jwt-secret>
```

**Save and exit:** `Ctrl+X`, then `Y`, then `Enter`

### Step 6: Set Permissions

```bash
chmod 600 .env
chmod +x ec2-deploy.sh
```

### Step 7: Deploy

```bash
./ec2-deploy.sh
```

The script will:
- Install Docker
- Install Docker Compose
- Pull images from Docker Hub
- Start services
- Verify MongoDB authentication
- Seed initial data

### Step 8: Verify Deployment

```bash
# Check services
docker-compose -f docker-compose.hub.yml ps

# Test frontend
curl http://localhost

# Test backend
curl http://localhost:3000/api/health

# View logs
docker-compose -f docker-compose.hub.yml logs -f
```

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
Frontend: http://13.229.74.162
Backend:  http://13.229.74.162:3000
```

### Login

```
Username: admin
Password: admin123

⚠️ Change password immediately after first login!
```

---

## 📊 Configuration Summary

| Item | Value |
|------|-------|
| Docker Hub Username | linhtranphu |
| EC2 IP | 13.229.74.162 |
| EC2 User | ec2-user |
| Frontend Port | 80 |
| Backend Port | 3000 |
| MongoDB Port | 27017 (internal only) |

---

## 🔐 Security Notes

1. **MongoDB Password**
   - Generate with: `openssl rand -base64 32`
   - Min 16 characters
   - Use in both MONGO_INITDB_ROOT_PASSWORD and MONGODB_URI

2. **JWT Secret**
   - Generate with: `openssl rand -base64 64`
   - Min 32 characters
   - Keep secure

3. **Admin Password**
   - Default: admin123
   - Change immediately after first login
   - Use strong password (min 12 chars)

---

## 🆘 Troubleshooting

### SSH Connection Failed

```bash
# Check if key pair is correct
ssh -i YOUR_KEY.pem -v ec2-user@13.229.74.162

# If permission denied, the key doesn't match
# You need to use the correct EC2 key pair
```

### Docker Not Found

```bash
# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER
```

### Docker Compose Not Found

```bash
# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

### MongoDB Authentication Failed

```bash
# Check .env
cat .env | grep MONGO_INITDB

# Check MongoDB logs
docker-compose -f docker-compose.hub.yml logs mongodb

# Restart MongoDB
docker-compose -f docker-compose.hub.yml restart mongodb
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

1. **PEM File Issue**
   - The provided PEM file doesn't match the EC2 instance
   - Use your actual EC2 key pair for SSH/SCP commands
   - Replace `YOUR_KEY.pem` with your actual key file

2. **Files Created**
   - docker-compose.hub.yml.ec2 ✅
   - .env.ec2 ✅
   - ec2-deploy.sh ✅

3. **Next Steps**
   - Copy files to EC2 with correct key pair
   - Edit .env with secure passwords
   - Run ec2-deploy.sh
   - Access application

---

## 🎯 Quick Commands

```bash
# Copy files (replace YOUR_KEY.pem)
scp -i YOUR_KEY.pem docker-compose.hub.yml.ec2 ec2-user@13.229.74.162:~/cafe-pos/
scp -i YOUR_KEY.pem .env.ec2 ec2-user@13.229.74.162:~/cafe-pos/
scp -i YOUR_KEY.pem ec2-deploy.sh ec2-user@13.229.74.162:~/cafe-pos/

# SSH to EC2
ssh -i YOUR_KEY.pem ec2-user@13.229.74.162

# On EC2
cd ~/cafe-pos
mv docker-compose.hub.yml.ec2 docker-compose.hub.yml
mv .env.ec2 .env
nano .env  # Edit with secure passwords
chmod 600 .env
chmod +x ec2-deploy.sh
./ec2-deploy.sh
```

---

**Status**: Ready for Manual Deployment  
**Last Updated**: January 2026
