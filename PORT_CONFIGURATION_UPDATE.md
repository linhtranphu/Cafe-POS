# 🔧 Port Configuration Update

## 📋 Thay đổi thực hiện

### 1. Frontend Port Configuration

**File**: `frontend/src/services/api.js`
- **Trước**: `baseURL: 'http://localhost:8080/api'`
- **Sau**: `baseURL: 'http://localhost:3000/api'`

### 2. Nginx Proxy Configuration

**File**: `frontend/nginx.conf`
- **Trước**: `proxy_pass http://backend:8080;`
- **Sau**: `proxy_pass http://backend:3000;`

### 3. Docker Compose Configuration

**File**: `docker-compose.hub.yml`
- **Backend PORT**: `8080` → `3000`
- **Backend ports**: `"8080:8080"` → `"3000:3000"`
- **Backend healthcheck**: `http://localhost:8080/api/health` → `http://localhost:3000/api/health`

### 4. Environment Variable

**File**: `.env`
- **Trước**: `PORT=8080`
- **Sau**: `PORT=3000`

### 5. Deploy Script

**File**: `deploy-to-ec2.sh`
- Updated docker-compose.hub.yml.ec2 template
- Updated .env.ec2 template
- PORT=3000 in environment

---

## 🔌 Port Mapping

### Local Development
```
Frontend:  http://localhost:80 (via nginx)
Backend:   http://localhost:3000
MongoDB:   localhost:27017
```

### Docker Compose
```
Frontend:  port 80 → container port 80
Backend:   port 3000 → container port 3000
MongoDB:   port 27017 → container port 27017
```

### EC2 Production
```
Frontend:  http://your-ec2-ip:80
Backend:   http://your-ec2-ip:3000
MongoDB:   your-ec2-ip:27017 (internal only)
```

---

## ✅ Build Status

### Frontend
- ✅ `npm run build` - Success
- ✅ Updated API endpoint to port 3000
- ✅ Updated nginx proxy to port 3000

### Backend
- ✅ `go build ./...` - Success
- ✅ Ready to run on port 3000

---

## 🚀 Next Steps

### 1. Rebuild Docker Images

```bash
# Backend
cd backend
docker build -t your-username/cafe-pos-backend:v1.0.1 .
docker tag your-username/cafe-pos-backend:v1.0.1 your-username/cafe-pos-backend:latest
docker push your-username/cafe-pos-backend:v1.0.1
docker push your-username/cafe-pos-backend:latest
cd ..

# Frontend
cd frontend
docker build -t your-username/cafe-pos-frontend:v1.0.1 .
docker tag your-username/cafe-pos-frontend:v1.0.1 your-username/cafe-pos-frontend:latest
docker push your-username/cafe-pos-frontend:v1.0.1
docker push your-username/cafe-pos-frontend:latest
cd ..
```

### 2. Test Locally

```bash
# Start services
docker-compose -f docker-compose.hub.yml up -d

# Wait for services
sleep 30

# Test frontend
curl http://localhost

# Test backend
curl http://localhost:3000/api/health

# Check logs
docker-compose -f docker-compose.hub.yml logs -f
```

### 3. Deploy to EC2

```bash
# Run deployment script
./deploy-to-ec2.sh

# Or manually:
# 1. Copy files to EC2
# 2. Configure .env
# 3. Run ./ec2-deploy.sh
```

---

## 📊 Configuration Summary

| Component | Port | Protocol | Access |
|-----------|------|----------|--------|
| Frontend | 80 | HTTP | http://localhost |
| Backend | 3000 | HTTP | http://localhost:3000 |
| MongoDB | 27017 | TCP | localhost:27017 |

---

## 🔐 Security Notes

1. **MongoDB Port (27017)**
   - ✅ NOT exposed to external network
   - ✅ Only accessible from Docker network
   - ✅ Requires authentication

2. **Backend Port (3000)**
   - ✅ Exposed for API access
   - ✅ Can be restricted via firewall
   - ✅ Should use HTTPS in production

3. **Frontend Port (80)**
   - ✅ Standard HTTP port
   - ✅ Should use HTTPS (port 443) in production
   - ✅ Nginx handles static files and API proxy

---

## 📝 Files Modified

1. `frontend/src/services/api.js` - API endpoint port
2. `frontend/nginx.conf` - Nginx proxy port
3. `docker-compose.hub.yml` - Backend port mapping
4. `.env` - PORT environment variable
5. `deploy-to-ec2.sh` - Updated templates

---

## ✅ Verification Checklist

- [x] Frontend API endpoint updated to port 3000
- [x] Nginx proxy updated to port 3000
- [x] Docker compose updated to port 3000
- [x] Environment variable updated to PORT=3000
- [x] Frontend build successful
- [x] Backend build successful
- [ ] Docker images rebuilt and pushed
- [ ] Local testing completed
- [ ] EC2 deployment completed

---

## 🎯 Testing Commands

### Local Testing

```bash
# Start services
docker-compose -f docker-compose.hub.yml up -d

# Test frontend
curl http://localhost

# Test backend API
curl http://localhost:3000/api/health

# Test MongoDB
docker exec -it cafe-pos-mongodb mongosh \
  -u admin \
  -p $(grep MONGO_INITDB_ROOT_PASSWORD .env | cut -d= -f2) \
  --authenticationDatabase admin

# View logs
docker-compose -f docker-compose.hub.yml logs -f

# Stop services
docker-compose -f docker-compose.hub.yml down
```

### EC2 Testing

```bash
# SSH to EC2
ssh -i your-key.pem ec2-user@your-ec2-ip

# Test frontend
curl http://localhost

# Test backend
curl http://localhost:3000/api/health

# View logs
docker-compose -f docker-compose.hub.yml logs -f

# Check ports
sudo netstat -tulpn | grep -E ':(80|3000|27017)'
```

---

## 📞 Important Notes

1. **Port 80 for Frontend**
   - Standard HTTP port
   - Nginx listens on port 80
   - Serves static files and proxies API requests

2. **Port 3000 for Backend**
   - Go backend application port
   - Handles API requests
   - Proxied through nginx on frontend

3. **Port 27017 for MongoDB**
   - Internal Docker network only
   - NOT exposed to external network
   - Requires authentication

---

**Version**: 1.0.0  
**Last Updated**: January 2026  
**Status**: ✅ Configuration Updated
