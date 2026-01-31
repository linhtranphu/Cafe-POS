#!/bin/bash

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║     🚀 Deploy Café POS from GitHub to EC2                ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════════╝${NC}"
echo ""

# Configuration
GITHUB_REPO="https://github.com/linhtranphu/Cafe-POS.git"
DEPLOY_DIR="$HOME/cafe-pos-deploy"
DOCKER_USERNAME="linhtranphu"

echo -e "${BLUE}📋 Configuration:${NC}"
echo "   GitHub Repo: $GITHUB_REPO"
echo "   Deploy Dir: $DEPLOY_DIR"
echo "   Docker Username: $DOCKER_USERNAME"
echo ""

# ============================================
# PART 1: Install Prerequisites
# ============================================

echo -e "${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}PART 1: Install Prerequisites${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo ""

# Install Git
echo -e "${BLUE}📦 Step 1: Installing Git...${NC}"
if ! command -v git &> /dev/null; then
    sudo yum install -y git
    echo -e "${GREEN}✅ Git installed${NC}"
else
    echo -e "${GREEN}✅ Git already installed${NC}"
fi
echo ""

# Install Docker
echo -e "${BLUE}📦 Step 2: Installing Docker...${NC}"
if ! command -v docker &> /dev/null; then
    curl -fsSL https://get.docker.com -o get-docker.sh
    sudo sh get-docker.sh
    sudo usermod -aG docker $USER
    echo -e "${GREEN}✅ Docker installed${NC}"
else
    echo -e "${GREEN}✅ Docker already installed${NC}"
fi
echo ""

# Install Docker Compose
echo -e "${BLUE}📦 Step 3: Installing Docker Compose...${NC}"
if ! command -v docker-compose &> /dev/null; then
    sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    sudo chmod +x /usr/local/bin/docker-compose
    echo -e "${GREEN}✅ Docker Compose installed${NC}"
else
    echo -e "${GREEN}✅ Docker Compose already installed${NC}"
fi
echo ""

# ============================================
# PART 2: Clone Repository
# ============================================

echo -e "${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}PART 2: Clone Repository from GitHub${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo ""

# Create deploy directory
echo -e "${BLUE}📁 Creating deploy directory...${NC}"
mkdir -p "$DEPLOY_DIR"
cd "$DEPLOY_DIR"
echo -e "${GREEN}✅ Directory created: $DEPLOY_DIR${NC}"
echo ""

# Clone repository
echo -e "${BLUE}📥 Cloning repository from GitHub...${NC}"
if [ -d ".git" ]; then
    echo -e "${YELLOW}⚠️  Repository already exists, pulling latest changes...${NC}"
    git pull origin main
else
    git clone "$GITHUB_REPO" .
fi
echo -e "${GREEN}✅ Repository cloned/updated${NC}"
echo ""

# ============================================
# PART 3: Setup Environment
# ============================================

echo -e "${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}PART 3: Setup Environment${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo ""

# Check if .env exists
if [ ! -f ".env" ]; then
    echo -e "${YELLOW}⚠️  .env file not found!${NC}"
    echo ""
    echo -e "${BLUE}📝 Creating .env file...${NC}"
    
    # Create .env with empty passwords
    cat > .env << EOF
# ============================================
# MongoDB Configuration
# ============================================
MONGO_INITDB_ROOT_USERNAME=admin
MONGO_INITDB_ROOT_PASSWORD=
MONGO_INITDB_DATABASE=cafe_pos

# ============================================
# Backend Configuration
# ============================================
MONGODB_URI=mongodb://admin:@mongodb:27017
MONGODB_DATABASE=cafe_pos
JWT_SECRET=
PORT=3000
EOF
    
    echo -e "${GREEN}✅ .env file created${NC}"
    echo ""
    echo -e "${YELLOW}⚠️  IMPORTANT: Edit .env file with your passwords!${NC}"
    echo ""
    echo -e "${BLUE}📝 To set passwords:${NC}"
    echo "   1. nano .env"
    echo "   2. Generate MongoDB password: openssl rand -base64 32"
    echo "   3. Generate JWT secret: openssl rand -base64 64"
    echo "   4. Update MONGO_INITDB_ROOT_PASSWORD and MONGODB_URI"
    echo "   5. Update JWT_SECRET"
    echo "   6. Save and exit (Ctrl+X, Y, Enter)"
    echo ""
else
    echo -e "${GREEN}✅ .env file already exists${NC}"
fi
echo ""

# Set permissions
chmod 600 .env
echo -e "${GREEN}✅ .env permissions set to 600${NC}"
echo ""

# ============================================
# PART 4: Deploy with Docker Compose
# ============================================

echo -e "${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}PART 4: Deploy with Docker Compose${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo ""

# Pull images
echo -e "${BLUE}📥 Step 1: Pulling Docker images...${NC}"
docker-compose -f docker-compose.hub.yml pull
echo -e "${GREEN}✅ Images pulled${NC}"
echo ""

# Stop existing services
echo -e "${BLUE}🛑 Step 2: Stopping existing services...${NC}"
docker-compose -f docker-compose.hub.yml down || true
echo -e "${GREEN}✅ Services stopped${NC}"
echo ""

# Start services
echo -e "${BLUE}🚀 Step 3: Starting services...${NC}"
docker-compose -f docker-compose.hub.yml up -d
echo -e "${GREEN}✅ Services started${NC}"
echo ""

# Wait for services
echo -e "${BLUE}⏳ Step 4: Waiting for services to be ready...${NC}"
sleep 30
echo -e "${GREEN}✅ Services ready${NC}"
echo ""

# Verify MongoDB
echo -e "${BLUE}🔐 Step 5: Verifying MongoDB authentication...${NC}"
source .env
if docker exec cafe-pos-mongodb mongosh \
    -u "$MONGO_INITDB_ROOT_USERNAME" \
    -p "$MONGO_INITDB_ROOT_PASSWORD" \
    --authenticationDatabase admin \
    --eval "db.adminCommand('ping')" > /dev/null 2>&1; then
    echo -e "${GREEN}✅ MongoDB authentication successful${NC}"
else
    echo -e "${RED}❌ MongoDB authentication failed${NC}"
    exit 1
fi
echo ""

# Seed data
echo -e "${BLUE}🌱 Step 6: Seeding initial data...${NC}"
docker exec -it cafe-pos-backend ./cafe-pos-server seed
echo -e "${GREEN}✅ Data seeded${NC}"
echo ""

# ============================================
# PART 5: Verification
# ============================================

echo -e "${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}PART 5: Verification${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo ""

# Check services
echo -e "${BLUE}📊 Services Status:${NC}"
docker-compose -f docker-compose.hub.yml ps
echo ""

# Test frontend
echo -e "${BLUE}🧪 Testing Frontend...${NC}"
if curl -s http://localhost > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Frontend is accessible${NC}"
else
    echo -e "${YELLOW}⚠️  Frontend not responding yet${NC}"
fi
echo ""

# Test backend
echo -e "${BLUE}🧪 Testing Backend...${NC}"
if curl -s http://localhost:3000/api/health > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Backend is healthy${NC}"
else
    echo -e "${YELLOW}⚠️  Backend not responding yet${NC}"
fi
echo ""

# ============================================
# SUMMARY
# ============================================

echo -e "${GREEN}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║          ✅ Deployment Completed Successfully!           ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════════════════════════╝${NC}"
echo ""

# Get public IP
PUBLIC_IP=$(curl -s http://169.254.169.254/latest/meta-data/public-ipv4)

echo -e "${BLUE}🌐 Access Application:${NC}"
echo "   Frontend: http://$PUBLIC_IP"
echo "   Backend:  http://$PUBLIC_IP:3000"
echo ""

echo -e "${BLUE}🔑 Credentials:${NC}"
echo "   App Username: admin"
echo "   App Password: admin123"
echo "   ⚠️  Change password after first login!"
echo ""

echo -e "${BLUE}📁 Deployment Directory:${NC}"
echo "   $DEPLOY_DIR"
echo ""

echo -e "${BLUE}📚 Useful Commands:${NC}"
echo "   View logs:    docker-compose -f docker-compose.hub.yml logs -f"
echo "   Restart:      docker-compose -f docker-compose.hub.yml restart"
echo "   Stop:         docker-compose -f docker-compose.hub.yml down"
echo "   Update:       git pull && docker-compose -f docker-compose.hub.yml pull && docker-compose -f docker-compose.hub.yml up -d"
echo ""

echo -e "${BLUE}📝 Environment File:${NC}"
echo "   Location: $DEPLOY_DIR/.env"
echo "   Permissions: 600 (secure)"
echo ""

echo -e "${YELLOW}⚠️  Important:${NC}"
echo "   - Change admin password immediately"
echo "   - Keep .env file secure"
echo "   - Regular backups recommended"
echo "   - Monitor logs for errors"
echo ""

echo -e "${GREEN}✅ Deployment Ready!${NC}"
echo ""
