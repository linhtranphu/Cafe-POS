# 🚀 Café POS - Docker Deployment Guide

## 📋 Prerequisites

- Docker Engine 20.10+
- Docker Compose 2.0+
- Make (optional, for convenience)

## 🏗️ Architecture

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  Frontend   │────▶│   Backend   │────▶│   MongoDB   │
│  (Vue.js)   │     │    (Go)     │     │             │
│   Port 80   │     │  Port 8080  │     │ Port 27017  │
└─────────────┘     └─────────────┘     └─────────────┘
```

## 🚀 Quick Start

### Option 1: Using Make (Recommended)

```bash
# Build and deploy all services
make deploy

# View logs
make logs

# Stop all services
make down
```

### Option 2: Using Docker Compose

```bash
# Build images
docker-compose build

# Start services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

## 📦 Services

### Frontend (Vue.js + Nginx)
- **URL**: http://localhost
- **Port**: 80
- **Image**: cafe-pos-frontend

### Backend (Go + Gin)
- **URL**: http://localhost:8080
- **Port**: 8080
- **Image**: cafe-pos-backend

### MongoDB
- **URL**: mongodb://localhost:27017
- **Port**: 27017
- **Username**: admin
- **Password**: admin123
- **Database**: cafe_pos

## 🔧 Configuration

### Environment Variables

Edit `docker-compose.yml` to change:

**Backend:**
- `MONGODB_URI`: MongoDB connection string
- `JWT_SECRET`: JWT secret key (⚠️ CHANGE IN PRODUCTION)
- `PORT`: Backend port

**MongoDB:**
- `MONGO_INITDB_ROOT_USERNAME`: MongoDB admin username
- `MONGO_INITDB_ROOT_PASSWORD`: MongoDB admin password

## 📝 Common Commands

```bash
# View all running containers
make ps
# or
docker-compose ps

# View backend logs
make logs-backend

# View frontend logs
make logs-frontend

# View MongoDB logs
make logs-mongodb

# Execute shell in backend container
make exec-backend

# Execute MongoDB shell
make exec-mongodb

# Restart all services
make restart

# Clean up everything (⚠️ removes volumes)
make clean
```

## 🔍 Health Checks

All services have health checks configured:

- **Backend**: `http://localhost:8080/api/health`
- **Frontend**: `http://localhost`
- **MongoDB**: Internal ping command

## 📊 Monitoring

View real-time logs:
```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f backend
docker-compose logs -f frontend
docker-compose logs -f mongodb
```

## 🛠️ Troubleshooting

### Port already in use
```bash
# Check what's using the port
lsof -ti:80 | xargs kill -9   # Frontend
lsof -ti:8080 | xargs kill -9 # Backend
lsof -ti:27017 | xargs kill -9 # MongoDB
```

### Rebuild images
```bash
# Rebuild all
docker-compose build --no-cache

# Rebuild specific service
docker-compose build --no-cache backend
```

### Reset everything
```bash
# Stop and remove everything
make clean

# Start fresh
make deploy
```

## 🔐 Security Notes

⚠️ **IMPORTANT for Production:**

1. Change `JWT_SECRET` in docker-compose.yml
2. Change MongoDB credentials
3. Use environment variables file (.env)
4. Enable HTTPS/SSL
5. Configure firewall rules
6. Use Docker secrets for sensitive data

## 📈 Scaling

Scale specific services:
```bash
# Scale backend to 3 instances
docker-compose up -d --scale backend=3
```

## 🎯 Default Users

After deployment, login with:
- **Manager**: `admin/admin123`
- **Waiter**: `waiter1/waiter123`
- **Cashier**: `cashier1/cashier123`

## 📞 Support

For issues, check:
1. Service logs: `make logs`
2. Container status: `make ps`
3. Health checks: `curl http://localhost:8080/api/health`
