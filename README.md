# ☕ Café POS System

Modern Point of Sale system for café management with role-based access control.

## 🚀 Quick Start

### Prerequisites
- Docker & Docker Compose
- 2GB RAM minimum
- 10GB disk space

### Production Deployment

1. **Clone repository**
```bash
git clone <repository-url>
cd cafe-pos
```

2. **Configure environment**
```bash
cp .env.example .env
nano .env  # Edit with secure credentials
```

3. **Deploy**
```bash
# Option 1: Build from source
docker-compose up -d --build

# Option 2: Use Docker Hub images
docker-compose -f docker-compose.hub.yml up -d
```

4. **Seed initial data**
```bash
docker exec -it cafe-pos-backend ./cafe-pos-server seed
```

5. **Access application**
- URL: `http://localhost` or `http://your-server-ip`
- Default login: `admin` / `admin123`
- **⚠️ Change password immediately after first login!**

## 📚 Documentation

- **[Production Deployment Guide](PRODUCTION_DEPLOYMENT.md)** - Complete deployment instructions
- **[Security Checklist](PRODUCTION_DEPLOYMENT.md#-security-checklist)** - Security best practices
- **[Backup & Restore](PRODUCTION_DEPLOYMENT.md#-backup--restore)** - Data backup procedures

## 👥 User Roles

- **Manager**: Full system access, user management, reports
- **Waiter**: Order management, table service
- **Barista**: Drink preparation, order fulfillment
- **Cashier**: Payment processing, shift management

## 🔒 Security

- ✅ JWT-based authentication
- ✅ Role-based access control
- ✅ Password hashing (bcrypt)
- ✅ Environment-based configuration
- ✅ MongoDB authentication

**Important**: 
- Never commit `.env` file
- Change all default passwords
- Use strong JWT secret (min 32 chars)
- Enable HTTPS in production

## 🛠️ Tech Stack

**Backend**:
- Go 1.21+
- Gin Web Framework
- MongoDB
- JWT Authentication

**Frontend**:
- Vue.js 3
- Vite
- Tailwind CSS
- Pinia (State Management)

**Infrastructure**:
- Docker & Docker Compose
- Nginx
- MongoDB 7.0

## 📊 Features

### Manager Features
- User management (create, edit, delete users)
- Ingredient inventory tracking
- Facility & equipment management
- Expense tracking (manual & auto-generated)
- Reports & analytics
- Menu management

### Waiter Features
- Order creation & management
- Table management
- Order status tracking
- Shift management

### Barista Features
- Order queue management
- Drink preparation workflow
- Shift tracking
- Order completion

### Cashier Features
- Payment processing
- Cash reconciliation
- Shift closure
- Payment discrepancy handling
- Daily reports

## 🔄 Updates

```bash
# Pull latest changes
git pull

# Rebuild and restart
docker-compose down
docker-compose up -d --build
```

## 🆘 Troubleshooting

### Services won't start
```bash
# Check logs
docker-compose logs -f

# Restart services
docker-compose restart
```

### Cannot connect to backend
```bash
# Check backend health
curl http://localhost:8080/api/health

# Check if backend is running
docker-compose ps
```

### Database issues
```bash
# Check MongoDB logs
docker-compose logs mongodb

# Verify MongoDB is running
docker exec -it cafe-pos-mongodb mongosh --eval "db.adminCommand('ping')"
```

## 📝 License

Proprietary - All rights reserved

## 📞 Support

For deployment issues, refer to [PRODUCTION_DEPLOYMENT.md](PRODUCTION_DEPLOYMENT.md)

---

**Version**: 1.0.0  
**Last Updated**: January 2026
