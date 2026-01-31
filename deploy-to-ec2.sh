#!/bin/bash

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║     🚀 Deploy Café POS to EC2 via Docker Hub             ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════════╝${NC}"
echo ""

# Configuration
DOCKER_USERNAME="linhtranphu"
EC2_IP="13.229.74.162"
EC2_USER="ec2-user"
PEM_FILE="EC2_PEM/OngTaPOS.pem"
EC2_PATH="~/cafe-pos"

echo -e "${BLUE}📋 Configuration:${NC}"
echo "   Docker Hub Username: $DOCKER_USERNAME"
echo "   EC2 IP: $EC2_IP"
echo "   EC2 User: $EC2_USER"
echo "   PEM File: $PEM_FILE"
echo "   EC2 Path: $EC2_PATH"
echo ""

# Check if PEM file exists
if [ ! -f "$PEM_FILE" ]; then
    echo -e "${RED}❌ Error: PEM file not found at $PEM_FILE${NC}"
    echo "Please ensure your EC2 key pair is in the EC2_PEM directory"
    exit 1
fi

# Set PEM file permissions
chmod 600 "$PEM_FILE"

# Check if running on local machine or EC2
if [ -f /.dockerenv ]; then
    echo -e "${RED}❌ This script should run on local machine, not in Docker!${NC}"
    exit 1
fi

# ============================================
# PART 1: Build and Push Images
# ============================================

echo -e "${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}PART 1: Build and Push Docker Images${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo ""

# Login Docker Hub
echo -e "${BLUE}📝 Logging in to Docker Hub...${NC}"
docker login

# Build Backend
echo ""
echo -e "${BLUE}🔨 Building backend image...${NC}"
cd backend
docker build -t $DOCKER_USERNAME/cafe-pos-backend:v1.0.1 .
docker tag $DOCKER_USERNAME/cafe-pos-backend:v1.0.1 $DOCKER_USERNAME/cafe-pos-backend:latest
echo -e "${GREEN}✅ Backend image built${NC}"

# Push Backend
echo -e "${BLUE}📤 Pushing backend image...${NC}"
docker push $DOCKER_USERNAME/cafe-pos-backend:v1.0.1
docker push $DOCKER_USERNAME/cafe-pos-backend:latest
echo -e "${GREEN}✅ Backend image pushed${NC}"
cd ..

# Build Frontend
echo ""
echo -e "${BLUE}🔨 Building frontend image...${NC}"
cd frontend
docker build -t $DOCKER_USERNAME/cafe-pos-frontend:v1.0.1 .
docker tag $DOCKER_USERNAME/cafe-pos-frontend:v1.0.1 $DOCKER_USERNAME/cafe-pos-frontend:latest
echo -e "${GREEN}✅ Frontend image built${NC}"

# Push Frontend
echo -e "${BLUE}� Pushing frontend image...${NC}"
docker push $DOCKER_USERNAME/cafe-pos-frontend:v1.0.1
docker push $DOCKER_USERNAME/cafe-pos-frontend:latest
echo -e "${GREEN}✅ Frontend image pushed${NC}"
cd ..

echo ""
echo -e "${GREEN}✅ All images pushed to Docker Hub!${NC}"
echo ""

# ============================================
# PART 2: Prepare EC2 Deployment Files
# ============================================

echo -e "${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}PART 2: Prepare EC2 Deployment Files${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo ""

# Create docker-compose.hub.yml for EC2
echo -e "${BLUE}📝 Creating docker-compose.hub.yml...${NC}"
cat > docker-compose.hub.yml.ec2 << EOF
version: '3.8'

services:
  mongodb:
    image: mongo:7.0
    container_name: cafe-pos-mongodb
    restart: always
    environment:
      MONGO_INITDB_ROOT_USERNAME: \${MONGO_INITDB_ROOT_USERNAME}
      MONGO_INITDB_ROOT_PASSWORD: \${MONGO_INITDB_ROOT_PASSWORD}
      MONGO_INITDB_DATABASE: \${MONGO_INITDB_DATABASE}
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
    image: $DOCKER_USERNAME/cafe-pos-backend:latest
    container_name: cafe-pos-backend
    restart: always
    environment:
      - MONGODB_URI=\${MONGODB_URI}
      - MONGODB_DATABASE=\${MONGODB_DATABASE}
      - JWT_SECRET=\${JWT_SECRET}
      - PORT=3000
    ports:
      - "3000:3000"
    depends_on:
      mongodb:
        condition: service_healthy
    networks:
      - cafe-pos-network
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:3000/api/health"]
      interval: 30s
      timeout: 10s
      retries: 3

  frontend:
    image: $DOCKER_USERNAME/cafe-pos-frontend:latest
    container_name: cafe-pos-frontend
    restart: always
    ports:
      - "80:80"
    depends_on:
      - backend
    networks:
      - cafe-pos-network
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost"]
      interval: 30s
      timeout: 10s
      retries: 3

volumes:
  mongodb_data:
    driver: local
  mongodb_config:
    driver: local

networks:
  cafe-pos-network:
    driver: bridge
EOF

echo -e "${GREEN}✅ docker-compose.hub.yml created${NC}"

# Create .env.ec2 template
echo -e "${BLUE}📝 Creating .env template...${NC}"
cat > .env.ec2 << 'EOF'
# ============================================
# MongoDB Configuration
# ============================================
MONGO_INITDB_ROOT_USERNAME=admin
MONGO_INITDB_ROOT_PASSWORD=CHANGE_THIS_TO_SECURE_PASSWORD
MONGO_INITDB_DATABASE=cafe_pos

# ============================================
# Backend Configuration
# ============================================
MONGODB_URI=mongodb://admin:CHANGE_THIS_TO_SECURE_PASSWORD@mongodb:27017
MONGODB_DATABASE=cafe_pos
JWT_SECRET=CHANGE_THIS_TO_SECURE_JWT_SECRET
PORT=3000
EOF

echo -e "${GREEN}✅ .env template created${NC}"

# Create EC2 deployment script
echo -e "${BLUE}📝 Creating EC2 deployment script...${NC}"
cat > ec2-deploy.sh << 'EOF'
#!/bin/bash

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║     🚀 EC2 Deployment - Setup & Deploy                   ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════════╝${NC}"
echo ""

# Check if .env exists
if [ ! -f .env ]; then
    echo -e "${RED}❌ Error: .env file not found!${NC}"
    echo "Please create .env file with MongoDB password"
    exit 1
fi

source .env

# Step 1: Install Docker
echo -e "${BLUE}📦 Step 1: Installing Docker...${NC}"
if ! command -v docker &> /dev/null; then
    curl -fsSL https://get.docker.com -o get-docker.sh
    sudo sh get-docker.sh
    sudo usermod -aG docker $USER
    echo -e "${GREEN}✅ Docker installed${NC}"
else
    echo -e "${GREEN}✅ Docker already installed${NC}"
fi
echo ""

# Step 2: Install Docker Compose
echo -e "${BLUE}� Step 2: Installing Docker Compose...${NC}"
if ! command -v docker-compose &> /dev/null; then
    sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    sudo chmod +x /usr/local/bin/docker-compose
    echo -e "${GREEN}✅ Docker Compose installed${NC}"
else
    echo -e "${GREEN}✅ Docker Compose already installed${NC}"
fi
echo ""

# Step 3: Pull images
echo -e "${BLUE}� Step 3: Pulling Docker images...${NC}"
docker-compose -f docker-compose.hub.yml pull
echo -e "${GREEN}✅ Images pulled${NC}"
echo ""

# Step 4: Start services
echo -e "${BLUE}🚀 Step 4: Starting services...${NC}"
docker-compose -f docker-compose.hub.yml up -d
echo -e "${GREEN}✅ Services started${NC}"
echo ""

# Step 5: Wait for services
echo -e "${BLUE}⏳ Step 5: Waiting for services to be ready...${NC}"
sleep 30
echo -e "${GREEN}✅ Services ready${NC}"
echo ""

# Step 6: Verify MongoDB
echo -e "${BLUE}🔐 Step 6: Verifying MongoDB authentication...${NC}"
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

# Step 7: Seed data
echo -e "${BLUE}🌱 Step 7: Seeding initial data...${NC}"
docker exec -it cafe-pos-backend ./cafe-pos-server seed
echo -e "${GREEN}✅ Data seeded${NC}"
echo ""

# Summary
echo -e "${GREEN}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║          ✅ EC2 Deployment Completed Successfully!        ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════════════════════════╝${NC}"
echo ""

echo -e "${BLUE}🌐 Access Application:${NC}"
echo "   Frontend: http://$(curl -s http://169.254.169.254/latest/meta-data/public-ipv4)"
echo "   Backend:  http://$(curl -s http://169.254.169.254/latest/meta-data/public-ipv4):3000"
echo ""

echo -e "${BLUE}🔑 Credentials:${NC}"
echo "   App Username: admin"
echo "   App Password: admin123"
echo "   ⚠️  Change password after first login!"
echo ""

echo -e "${BLUE}� Useful Commands:${NC}"
echo "   View logs:    docker-compose -f docker-compose.hub.yml logs -f"
echo "   Restart:      docker-compose -f docker-compose.hub.yml restart"
echo "   Stop:         docker-compose -f docker-compose.hub.yml down"
echo ""
EOF

chmod +x ec2-deploy.sh
echo -e "${GREEN}✅ EC2 deployment script created${NC}"

# ============================================
# PART 3: Copy Files to EC2
# ============================================

echo ""
echo -e "${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}PART 3: Copy Files to EC2${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo ""

echo -e "${BLUE}📋 Files to copy:${NC}"
echo "   1. docker-compose.hub.yml.ec2"
echo "   2. .env.ec2"
echo "   3. ec2-deploy.sh"
echo ""

echo -e "${BLUE}🔐 Using PEM file: $PEM_FILE${NC}"
echo ""

# Create cafe-pos directory on EC2
echo -e "${BLUE}📁 Creating ~/cafe-pos directory on EC2...${NC}"
ssh -i "$PEM_FILE" $EC2_USER@$EC2_IP "mkdir -p ~/cafe-pos" || true
echo -e "${GREEN}✅ Directory created${NC}"
echo ""

# Copy files to EC2
echo -e "${BLUE}📤 Copying files to EC2...${NC}"
scp -i "$PEM_FILE" docker-compose.hub.yml.ec2 $EC2_USER@$EC2_IP:~/cafe-pos/
echo -e "${GREEN}✅ docker-compose.hub.yml.ec2 copied${NC}"

scp -i "$PEM_FILE" .env.ec2 $EC2_USER@$EC2_IP:~/cafe-pos/
echo -e "${GREEN}✅ .env.ec2 copied${NC}"

scp -i "$PEM_FILE" ec2-deploy.sh $EC2_USER@$EC2_IP:~/cafe-pos/
echo -e "${GREEN}✅ ec2-deploy.sh copied${NC}"

echo ""
echo -e "${GREEN}╔════════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║     ✅ All Files Copied to EC2 Successfully!             ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════════════════════════╝${NC}"
echo ""

echo -e "${BLUE}🚀 Next steps on EC2:${NC}"
echo ""
echo "1. SSH to EC2:"
echo "   ssh -i $PEM_FILE $EC2_USER@$EC2_IP"
echo ""
echo "2. Setup files:"
echo "   cd ~/cafe-pos"
echo "   mv docker-compose.hub.yml.ec2 docker-compose.hub.yml"
echo "   mv .env.ec2 .env"
echo ""
echo "3. Edit .env with secure MongoDB password:"
echo "   nano .env"
echo ""
echo "   Generate secure password:"
echo "   openssl rand -base64 32"
echo ""
echo "   Update these lines:"
echo "   MONGO_INITDB_ROOT_PASSWORD=<generated-password>"
echo "   MONGODB_URI=mongodb://admin:<generated-password>@mongodb:27017"
echo "   JWT_SECRET=<generated-jwt-secret>"
echo ""
echo "4. Set permissions and deploy:"
echo "   chmod 600 .env"
echo "   chmod +x ec2-deploy.sh"
echo "   ./ec2-deploy.sh"
echo ""
echo "5. Access application:"
echo "   Frontend: http://$EC2_IP"
echo "   Backend:  http://$EC2_IP:3000"
echo ""

echo -e "${YELLOW}⚠️  Important:${NC}"
echo "   - Edit .env with secure MongoDB password"
echo "   - Never commit .env to git"
echo "   - Change admin password after first login"
echo "   - Keep PEM file secure"
echo ""
