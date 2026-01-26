#!/bin/bash

# Café POS - Auto Deploy Script for EC2
# Usage: ./deploy-ec2.sh

set -e

echo "🚀 Café POS - EC2 Deployment Script"
echo "===================================="

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if running on EC2
if [ ! -f /sys/hypervisor/uuid ] || ! grep -q ec2 /sys/hypervisor/uuid 2>/dev/null; then
    echo -e "${YELLOW}⚠️  Warning: This doesn't appear to be an EC2 instance${NC}"
    read -p "Continue anyway? (y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Step 1: Update system
echo -e "\n${GREEN}Step 1: Updating system...${NC}"
sudo apt update && sudo apt upgrade -y

# Step 2: Install Docker
if ! command_exists docker; then
    echo -e "\n${GREEN}Step 2: Installing Docker...${NC}"
    curl -fsSL https://get.docker.com -o get-docker.sh
    sudo sh get-docker.sh
    sudo usermod -aG docker $USER
    rm get-docker.sh
else
    echo -e "\n${GREEN}Step 2: Docker already installed${NC}"
fi

# Step 3: Install Docker Compose
if ! command_exists docker-compose; then
    echo -e "\n${GREEN}Step 3: Installing Docker Compose...${NC}"
    sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    sudo chmod +x /usr/local/bin/docker-compose
else
    echo -e "\n${GREEN}Step 3: Docker Compose already installed${NC}"
fi

# Step 4: Install Git
if ! command_exists git; then
    echo -e "\n${GREEN}Step 4: Installing Git...${NC}"
    sudo apt install git -y
else
    echo -e "\n${GREEN}Step 4: Git already installed${NC}"
fi

# Step 5: Setup environment
echo -e "\n${GREEN}Step 5: Setting up environment...${NC}"

if [ ! -f .env ]; then
    echo "Creating .env file..."
    cat > .env << 'EOF'
# MongoDB
MONGO_INITDB_ROOT_USERNAME=admin
MONGO_INITDB_ROOT_PASSWORD=CHANGE_THIS_PASSWORD_123

# Backend
JWT_SECRET=CHANGE_THIS_SECRET_KEY_IN_PRODUCTION
MONGODB_URI=mongodb://admin:CHANGE_THIS_PASSWORD_123@mongodb:27017
MONGODB_DATABASE=cafe_pos
PORT=8080

# Frontend
VITE_API_URL=http://localhost:8080
EOF
    chmod 600 .env
    echo -e "${YELLOW}⚠️  Please edit .env file and change default passwords!${NC}"
    echo "Run: nano .env"
    read -p "Press enter to continue after editing .env..."
fi

# Step 6: Build and deploy
echo -e "\n${GREEN}Step 6: Building and deploying...${NC}"
docker-compose build
docker-compose up -d

# Step 7: Wait for services to be ready
echo -e "\n${GREEN}Step 7: Waiting for services to start...${NC}"
sleep 10

# Step 8: Check status
echo -e "\n${GREEN}Step 8: Checking deployment status...${NC}"
docker-compose ps

# Step 9: Setup firewall
echo -e "\n${GREEN}Step 9: Setting up firewall...${NC}"
if command_exists ufw; then
    sudo ufw --force enable
    sudo ufw allow ssh
    sudo ufw allow 80/tcp
    sudo ufw allow 443/tcp
    sudo ufw allow 8080/tcp
    sudo ufw status
fi

# Step 10: Create backup script
echo -e "\n${GREEN}Step 10: Creating backup script...${NC}"
cat > backup-mongodb.sh << 'EOF'
#!/bin/bash
BACKUP_DIR="$HOME/backups"
DATE=$(date +%Y%m%d_%H%M%S)
mkdir -p $BACKUP_DIR

docker exec cafe-pos-mongodb mongodump \
  --username admin \
  --password CHANGE_THIS_PASSWORD_123 \
  --authenticationDatabase admin \
  --out /tmp/backup

docker cp cafe-pos-mongodb:/tmp/backup $BACKUP_DIR/mongodb_$DATE
tar -czf $BACKUP_DIR/mongodb_$DATE.tar.gz -C $BACKUP_DIR mongodb_$DATE
rm -rf $BACKUP_DIR/mongodb_$DATE

# Keep only last 7 backups
ls -t $BACKUP_DIR/mongodb_*.tar.gz | tail -n +8 | xargs -r rm

echo "Backup completed: $BACKUP_DIR/mongodb_$DATE.tar.gz"
EOF
chmod +x backup-mongodb.sh

# Get EC2 public IP
EC2_IP=$(curl -s http://169.254.169.254/latest/meta-data/public-ipv4 2>/dev/null || echo "localhost")

# Final message
echo -e "\n${GREEN}✅ Deployment completed!${NC}"
echo "===================================="
echo -e "Frontend: ${GREEN}http://$EC2_IP${NC}"
echo -e "Backend:  ${GREEN}http://$EC2_IP:8080${NC}"
echo ""
echo "Default users:"
echo "  - Manager: admin/admin123"
echo "  - Waiter:  waiter1/waiter123"
echo "  - Cashier: cashier1/cashier123"
echo ""
echo -e "${YELLOW}⚠️  Important:${NC}"
echo "1. Change default passwords in .env file"
echo "2. Update passwords in backup-mongodb.sh"
echo "3. Setup SSL/HTTPS for production"
echo ""
echo "Useful commands:"
echo "  - View logs:    docker-compose logs -f"
echo "  - Restart:      docker-compose restart"
echo "  - Stop:         docker-compose down"
echo "  - Backup:       ./backup-mongodb.sh"
echo ""
echo -e "${GREEN}Happy deploying! ☕${NC}"
