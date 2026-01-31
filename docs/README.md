# ☕ Café POS System

Hệ thống quản lý quán café với đầy đủ tính năng: Order, Payment, Inventory, Facility, Expense Management.

## 🚀 Quick Deploy (1 Command)

```bash
curl -fsSL https://raw.githubusercontent.com/linhtranphu/Cafe-POS/main/deploy-from-hub.sh | bash
```

## 📦 Docker Hub Images

- **Backend**: [linhtranphu/cafe-pos-backend](https://hub.docker.com/r/linhtranphu/cafe-pos-backend)
- **Frontend**: [linhtranphu/cafe-pos-frontend](https://hub.docker.com/r/linhtranphu/cafe-pos-frontend)

## 🎯 Features

- ✅ Order Management (Waiter)
- ✅ Payment & Reconciliation (Cashier)
- ✅ Shift Management
- ✅ Ingredient Inventory
- ✅ Facility Management
- ✅ Expense Tracking
- ✅ Multi-role Authentication (Manager/Waiter/Cashier)

## 🔧 Tech Stack

- **Frontend**: Vue.js 3 + Pinia + Tailwind CSS
- **Backend**: Go + Gin Framework
- **Database**: MongoDB
- **Deployment**: Docker + Docker Compose

## 📚 Documentation

- [Docker Hub Deployment Guide](./DOCKER-HUB-DEPLOY.md)
- [EC2 Deployment Guide](./EC2-DEPLOYMENT.md)
- [Local Development](./DEPLOYMENT.md)
- [Requirements](./requirements.md)

## 🚀 Deployment Options

### Option 1: Quick Deploy from Docker Hub

```bash
curl -fsSL https://raw.githubusercontent.com/linhtranphu/Cafe-POS/main/deploy-from-hub.sh | bash
```

### Option 2: Manual Deploy

```bash
# Clone repository
git clone https://github.com/linhtranphu/Cafe-POS.git
cd Cafe-POS

# Deploy with Docker Compose
docker-compose up -d
```

### Option 3: Deploy to EC2

See [EC2-DEPLOYMENT.md](./EC2-DEPLOYMENT.md) for detailed instructions.

## 🌐 Access

After deployment:
- **Frontend**: http://YOUR-SERVER-IP
- **Backend API**: http://YOUR-SERVER-IP:8080

## 👥 Default Users

- **Manager**: `admin/admin123`
- **Waiter**: `waiter1/waiter123`
- **Cashier**: `cashier1/cashier123`

⚠️ **Change these passwords in production!**

## 🔄 Update Application

```bash
cd ~/cafe-pos
docker-compose pull
docker-compose up -d
```

## 📊 Management Commands

```bash
# View logs
docker-compose logs -f

# Restart services
docker-compose restart

# Stop services
docker-compose down

# Check status
docker-compose ps
```

## 🛠️ Development

### Prerequisites

- Docker & Docker Compose
- Go 1.21+
- Node.js 18+
- MongoDB 7.0+

### Local Setup

```bash
# Clone repository
git clone https://github.com/linhtranphu/Cafe-POS.git
cd Cafe-POS

# Start with Docker Compose
docker-compose up -d

# Or run locally
cd backend && go run main.go
cd frontend && npm install && npm run dev
```

## 📝 Build & Push to Docker Hub

```bash
# Build and push images
./build-push.sh

# Or with version tag
./build-push.sh v1.0.0
```

## 🔐 Security

- Change default passwords in `.env`
- Setup firewall rules
- Enable HTTPS/SSL
- Regular backups
- Monitor logs

## 📈 Monitoring

```bash
# View logs
docker-compose logs -f

# Check resource usage
docker stats

# Check disk space
docker system df
```

## 💾 Backup

```bash
# Backup MongoDB
docker exec cafe-pos-mongodb mongodump --out /tmp/backup
docker cp cafe-pos-mongodb:/tmp/backup ./backup

# Restore MongoDB
docker cp ./backup cafe-pos-mongodb:/tmp/backup
docker exec cafe-pos-mongodb mongorestore /tmp/backup
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

This project is licensed under the MIT License.

## 📞 Support

For issues or questions:
- GitHub Issues: https://github.com/linhtranphu/Cafe-POS/issues
- Check logs: `docker-compose logs -f`
- Documentation: See `/docs` folder

## 🎯 Roadmap

- [ ] Bill printing
- [ ] Advanced reports & analytics
- [ ] Mobile app
- [ ] Multi-store support
- [ ] Integration with payment gateways

---

Made with ☕ by [linhtranphu](https://github.com/linhtranphu)
