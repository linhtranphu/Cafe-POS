# 🚀 START HERE - Café POS Deployment

**Welcome! Your Café POS application is ready to deploy.**

---

## ⚡ Quick Start (5 minutes)

### Step 1: SSH to EC2
```bash
ssh -i YOUR_KEY.pem ec2-user@13.229.74.162
```

### Step 2: Run Deployment
```bash
curl -fsSL https://raw.githubusercontent.com/linhtranphu/Cafe-POS/main/ec2-deploy-from-github.sh | bash
```

### Step 3: Wait for Completion
The script will:
- Install Docker & Docker Compose
- Clone repository from GitHub
- Generate secure passwords
- Deploy all services
- Seed initial data

**Time**: ~5-10 minutes

### Step 4: Access Application
```
Frontend: http://13.229.74.162
Backend:  http://13.229.74.162:3000
Login: admin / admin123
```

### Step 5: Change Admin Password
1. Go to Profile
2. Change password
3. Use strong password (min 12 chars)

---

## 📚 Documentation

### For Quick Deployment
- **Read**: `EC2_QUICK_START.md` (5-minute guide)
- **Reference**: `DEPLOYMENT_QUICK_REFERENCE.md` (commands)

### For Detailed Deployment
- **Read**: `DEPLOYMENT_COMPLETE_SUMMARY.md` (overview)
- **Follow**: `DEPLOY_TO_EC2_WITH_DOCKER_HUB.md` (step-by-step)
- **Reference**: `DEPLOYMENT_QUICK_REFERENCE.md` (commands)

### For Manual Deployment
- **Follow**: `MANUAL_EC2_DEPLOYMENT.md` (manual steps)
- **Reference**: `DEPLOYMENT_QUICK_REFERENCE.md` (commands)

### For Understanding Architecture
- **Read**: `DEPLOYMENT_INFRASTRUCTURE_SUMMARY.md` (architecture)
- **Read**: `DEPLOYMENT_COMPLETE_SUMMARY.md` (overview)

### For Verification
- **Check**: `DEPLOYMENT_READY_CHECKLIST.md` (pre-deployment)
- **Check**: `DEPLOYMENT_VERIFICATION.md` (verification report)

### For Navigation
- **Use**: `DEPLOYMENT_DOCUMENTATION_INDEX.md` (find any document)

---

## 🎯 Choose Your Path

### Path 1: Quick Deployment (Recommended)
```
1. SSH to EC2
2. Run one command
3. Wait for completion
4. Access application
```
**Time**: ~5-10 minutes  
**Complexity**: Simple  
**Guide**: `EC2_QUICK_START.md`

---

### Path 2: Detailed Deployment
```
1. Read overview
2. Follow step-by-step guide
3. Deploy to EC2
4. Verify deployment
5. Access application
```
**Time**: ~15-20 minutes  
**Complexity**: Medium  
**Guide**: `DEPLOY_TO_EC2_WITH_DOCKER_HUB.md`

---

### Path 3: Manual Deployment
```
1. SSH to EC2
2. Install Docker manually
3. Clone repository
4. Create .env file
5. Run docker-compose
6. Seed data
7. Access application
```
**Time**: ~20-30 minutes  
**Complexity**: High  
**Guide**: `MANUAL_EC2_DEPLOYMENT.md`

---

## ✅ Before You Deploy

Make sure you have:
- [ ] EC2 instance running at 13.229.74.162
- [ ] EC2 SSH key pair (YOUR_KEY.pem)
- [ ] Access to GitHub (public repository)
- [ ] Docker Hub account (images already pushed)

---

## 🔐 Security

### Credentials
```
MongoDB Username: admin
MongoDB Password: (empty - you must set)
JWT Secret: (empty - you must set)
Admin Username: admin
Admin Password: admin123 (change after login)
```

### Important
- ⚠️ You must set MongoDB password and JWT secret in .env
- ⚠️ Generate passwords with: openssl rand -base64 32 (MongoDB) and openssl rand -base64 64 (JWT)
- ✅ .env file not committed to git
- ✅ No hardcoded credentials in code
- ⚠️ Change admin password immediately after login

---

## 🌐 Access URLs

### After Deployment
```
Frontend: http://13.229.74.162
Backend:  http://13.229.74.162:3000
```

### Login
```
Username: admin
Password: admin123
```

---

## 📞 Need Help?

### Quick Issues
- **Port already in use**: See `DEPLOYMENT_QUICK_REFERENCE.md` troubleshooting
- **MongoDB connection failed**: See `DEPLOYMENT_QUICK_REFERENCE.md` troubleshooting
- **Frontend shows blank page**: See `DEPLOYMENT_QUICK_REFERENCE.md` troubleshooting

### Detailed Help
- **Manual deployment**: See `MANUAL_EC2_DEPLOYMENT.md`
- **Architecture questions**: See `DEPLOYMENT_INFRASTRUCTURE_SUMMARY.md`
- **Security questions**: See `PRODUCTION_DEPLOYMENT.md`

### Find Any Document
- **Navigation guide**: See `DEPLOYMENT_DOCUMENTATION_INDEX.md`

---

## 🎉 You're Ready!

Everything is prepared and tested. Choose your deployment method and get started!

### Recommended Command
```bash
# SSH to EC2
ssh -i YOUR_KEY.pem ec2-user@13.229.74.162

# Deploy from GitHub
curl -fsSL https://raw.githubusercontent.com/linhtranphu/Cafe-POS/main/ec2-deploy-from-github.sh | bash

# Wait for completion (~5-10 minutes)
# Access at http://13.229.74.162
```

---

## 📋 Deployment Checklist

- [ ] SSH to EC2 successful
- [ ] Deployment script running
- [ ] Services starting
- [ ] MongoDB authentication verified
- [ ] Initial data seeded
- [ ] Frontend accessible
- [ ] Backend accessible
- [ ] Login successful
- [ ] Admin password changed

---

## 🚀 Next Steps

1. **Deploy**: Choose your deployment method above
2. **Verify**: Check that all services are running
3. **Access**: Open http://13.229.74.162 in browser
4. **Login**: Use admin / admin123
5. **Change Password**: Update admin password immediately
6. **Test**: Try all features
7. **Monitor**: Check logs for errors

---

## 📚 Documentation Files

| File | Purpose | When to Read |
|------|---------|------------|
| DEPLOYMENT_START_HERE.md | This file - quick start | Now |
| EC2_QUICK_START.md | 5-minute deployment | For quick deployment |
| DEPLOYMENT_COMPLETE_SUMMARY.md | Complete overview | For understanding |
| DEPLOY_TO_EC2_WITH_DOCKER_HUB.md | Detailed guide | For step-by-step |
| MANUAL_EC2_DEPLOYMENT.md | Manual deployment | For manual setup |
| DEPLOYMENT_QUICK_REFERENCE.md | Quick commands | During deployment |
| DEPLOYMENT_INFRASTRUCTURE_SUMMARY.md | Architecture | For understanding |
| DEPLOYMENT_READY_CHECKLIST.md | Pre-deployment | Before deploying |
| DEPLOYMENT_VERIFICATION.md | Verification | To verify readiness |
| DEPLOYMENT_DOCUMENTATION_INDEX.md | Navigation | To find documents |

---

## ⚡ One-Liner Deployment

```bash
ssh -i YOUR_KEY.pem ec2-user@13.229.74.162 && curl -fsSL https://raw.githubusercontent.com/linhtranphu/Cafe-POS/main/ec2-deploy-from-github.sh | bash
```

---

**Status**: ✅ Ready for Deployment  
**Last Updated**: January 31, 2026  
**Next Action**: Deploy to EC2

🚀 **Let's go!**

