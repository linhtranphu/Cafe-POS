#!/bin/bash

# Deploy Café POS from Docker Hub
# Usage: curl -fsSL https://raw.githubusercontent.com/linhtranphu/Cafe-POS/main/deploy-from-hub.sh | bash
# Or: ./deploy-from-hub.sh

set -e

echo "🚀 Deploying Café POS from Docker Hub"
echo "======================================"

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Detect OS first
if [ -f /etc/os-release ]; then
    . /etc/os-release
    OS=$ID
fi

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "Installing Docker..."
    
    # Install Docker based on OS
    if [ "$OS" = "amzn" ]; then
        echo "Detected Amazon Linux"
        sudo yum update -y
        sudo yum install -y docker
        sudo systemctl start docker
        sudo systemctl enable docker
        sudo usermod -aG docker $USER
    else
        echo "Detected $OS"
        curl -fsSL https://get.docker.com | sh
        sudo usermod -aG docker $USER
    fi
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo "Installing Docker Compose..."
    sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    sudo chmod +x /usr/local/bin/docker-compose
fi

# Create project directory
PROJECT_DIR="$HOME/cafe-pos"
mkdir -p $PROJECT_DIR
cd $PROJECT_DIR

# Download docker-compose file
echo -e "\n${GREEN}Downloading docker-compose configuration...${NC}"
cat > docker-compose.yml << 'EOF'
version: '3.8'

services:
  mongodb:
    image: mongo:7.0
    container_name: cafe-pos-mongodb
    restart: always
    environment:
      MONGO_INITDB_ROOT_USERNAME: ${MONGO_INITDB_ROOT_USERNAME:-admin}
      MONGO_INITDB_ROOT_PASSWORD: ${MONGO_INITDB_ROOT_PASSWORD:-admin123}
      MONGO_INITDB_DATABASE: ${MONGODB_DATABASE:-cafe_pos}
    ports:
      - "27017:27017"
    volumes:
      - mongodb_data:/data/db
      - mongodb_config:/data/configdb
    networks:
      - cafe-pos-network
    healthcheck:
      test: echo 'db.runCommand("ping").ok' | mongosh localhost:27017/test --quiet
      interval: 10s
      timeout: 5s
      retries: 5

  backend:
    image: linhtranphu/cafe-pos-backend:latest
    container_name: cafe-pos-backend
    restart: always
    environment:
      - MONGODB_URI=${MONGODB_URI:-mongodb://admin:admin123@mongodb:27017}
      - MONGODB_DATABASE=${MONGODB_DATABASE:-cafe_pos}
      - JWT_SECRET=${JWT_SECRET:-your-secret-key-change-in-production}
      - PORT=8080
    ports:
      - "8080:8080"
    depends_on:
      mongodb:
        condition: service_healthy
    networks:
      - cafe-pos-network

  frontend:
    image: linhtranphu/cafe-pos-frontend:latest
    container_name: cafe-pos-frontend
    restart: always
    ports:
      - "80:80"
    depends_on:
      - backend
    networks:
      - cafe-pos-network

volumes:
  mongodb_data:
  mongodb_config:

networks:
  cafe-pos-network:
    driver: bridge
EOF

# Create .env file if not exists
if [ ! -f .env ]; then
    echo -e "\n${GREEN}Creating .env file...${NC}"
    cat > .env << 'EOF'
# MongoDB
MONGO_INITDB_ROOT_USERNAME=admin
MONGO_INITDB_ROOT_PASSWORD=CHANGE_THIS_PASSWORD_123

# Backend
JWT_SECRET=CHANGE_THIS_SECRET_KEY_IN_PRODUCTION
MONGODB_URI=mongodb://admin:CHANGE_THIS_PASSWORD_123@mongodb:27017
MONGODB_DATABASE=cafe_pos

# Frontend
VITE_API_URL=http://localhost:8080
EOF
    chmod 600 .env
    echo -e "${YELLOW}⚠️  Please edit .env file and change default passwords!${NC}"
fi

# Pull latest images
echo -e "\n${GREEN}Pulling latest images from Docker Hub...${NC}"
docker-compose pull

# Start services
echo -e "\n${GREEN}Starting services...${NC}"
docker-compose up -d

# Wait for services
echo -e "\n${GREEN}Waiting for services to start...${NC}"
sleep 10

# Check status
echo -e "\n${GREEN}Checking deployment status...${NC}"
docker-compose ps

# Get IP address
if command -v curl &> /dev/null; then
    EC2_IP=$(curl -s http://169.254.169.254/latest/meta-data/public-ipv4 2>/dev/null || hostname -I | awk '{print $1}')
else
    EC2_IP=$(hostname -I | awk '{print $1}')
fi

# Final message
echo ""
echo -e "${GREEN}✅ Deployment completed!${NC}"
echo "======================================"
echo -e "Frontend: ${GREEN}http://$EC2_IP${NC}"
echo -e "Backend:  ${GREEN}http://$EC2_IP:8080${NC}"
echo ""
echo "Default users:"
echo "  - Manager: admin/admin123"
echo "  - Waiter:  waiter1/waiter123"
echo "  - Cashier: cashier1/cashier123"
echo ""
echo -e "${YELLOW}⚠️  Important:${NC}"
echo "1. Edit .env file: nano $PROJECT_DIR/.env"
echo "2. Change default passwords"
echo "3. Restart: cd $PROJECT_DIR && docker-compose restart"
echo ""
echo "Useful commands:"
echo "  cd $PROJECT_DIR"
echo "  docker-compose logs -f    # View logs"
echo "  docker-compose restart    # Restart services"
echo "  docker-compose down       # Stop services"
echo "  docker-compose pull && docker-compose up -d  # Update to latest"
