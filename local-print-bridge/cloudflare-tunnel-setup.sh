#!/bin/bash

# Cloudflare Tunnel Setup for Local Print Bridge (Go version)
# Expose local Print Bridge to internet via Cloudflare Tunnel

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${BLUE}╔════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║   ☁️  Cloudflare Tunnel Setup for Print Bridge       ║${NC}"
echo -e "${BLUE}║              (Go + chromedp version)                   ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════╝${NC}"
echo ""

# Check if cloudflared is installed
if ! command -v cloudflared &> /dev/null; then
    echo -e "${YELLOW}⚠️  cloudflared not found. Installing...${NC}"
    echo ""
    
    # Detect OS
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        if command -v brew &> /dev/null; then
            echo "Installing via Homebrew..."
            brew install cloudflared
        else
            echo -e "${RED}Please install Homebrew first: https://brew.sh${NC}"
            exit 1
        fi
    elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
        # Linux
        echo "Downloading cloudflared for Linux..."
        wget -q https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-amd64.deb
        sudo dpkg -i cloudflared-linux-amd64.deb
        rm cloudflared-linux-amd64.deb
    elif [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "win32" ]]; then
        # Windows
        echo -e "${YELLOW}For Windows, please download from:${NC}"
        echo "https://github.com/cloudflare/cloudflared/releases/latest"
        echo "Then run: cloudflared.exe service install"
        exit 1
    else
        echo -e "${RED}Unsupported OS. Please install cloudflared manually:${NC}"
        echo "https://developers.cloudflare.com/cloudflare-one/connections/connect-apps/install-and-setup/installation/"
        exit 1
    fi
fi

echo -e "${GREEN}✅ cloudflared installed${NC}"
cloudflared --version
echo ""

# Login to Cloudflare
echo -e "${YELLOW}📝 Step 1: Login to Cloudflare${NC}"
echo "This will open a browser window for authentication..."
echo "Please login with your Cloudflare account that has access to tacafe.store domain"
echo ""
read -p "Press Enter to continue..."

cloudflared tunnel login

if [ $? -ne 0 ]; then
    echo -e "${RED}❌ Login failed${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Logged in to Cloudflare${NC}"
echo ""

# Create tunnel
TUNNEL_NAME="print-bridge-$(date +%s)"
echo -e "${YELLOW}📝 Step 2: Creating tunnel: ${TUNNEL_NAME}${NC}"
cloudflared tunnel create ${TUNNEL_NAME}

if [ $? -ne 0 ]; then
    echo -e "${RED}❌ Tunnel creation failed${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Tunnel created${NC}"
echo ""

# Get tunnel ID
TUNNEL_ID=$(cloudflared tunnel list | grep ${TUNNEL_NAME} | awk '{print $1}')
echo "Tunnel ID: ${TUNNEL_ID}"
echo "Tunnel Name: ${TUNNEL_NAME}"
echo ""

# Determine config path based on OS
if [[ "$OSTYPE" == "darwin"* ]]; then
    CONFIG_DIR="$HOME/.cloudflared"
    CREDENTIALS_FILE="$HOME/.cloudflared/${TUNNEL_ID}.json"
elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
    CONFIG_DIR="$HOME/.cloudflared"
    CREDENTIALS_FILE="$HOME/.cloudflared/${TUNNEL_ID}.json"
else
    CONFIG_DIR="$HOME/.cloudflared"
    CREDENTIALS_FILE="$HOME/.cloudflared/${TUNNEL_ID}.json"
fi

# Create config file
echo -e "${YELLOW}📝 Step 3: Creating tunnel configuration${NC}"
mkdir -p ${CONFIG_DIR}

cat > ${CONFIG_DIR}/config.yml << EOF
tunnel: ${TUNNEL_ID}
credentials-file: ${CREDENTIALS_FILE}

ingress:
  # Route print.tacafe.store to local print bridge
  - hostname: print.tacafe.store
    service: http://localhost:3001
  
  # Catch-all rule (required)
  - service: http_status:404

# Optional: Logging
loglevel: info
EOF

echo -e "${GREEN}✅ Configuration created at: ${CONFIG_DIR}/config.yml${NC}"
echo ""

# Save tunnel info to .env
echo -e "${YELLOW}📝 Step 4: Saving tunnel info to .env${NC}"
cat >> .env << EOF

# Cloudflare Tunnel Configuration
CLOUDFLARE_TUNNEL_ID=${TUNNEL_ID}
CLOUDFLARE_TUNNEL_NAME=${TUNNEL_NAME}
CLOUDFLARE_TUNNEL_URL=https://print.tacafe.store
EOF

echo -e "${GREEN}✅ Tunnel info saved to .env${NC}"
echo ""

# Instructions
echo ""
echo -e "${BLUE}╔════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║                  🎉 Setup Complete!                    ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════════╝${NC}"
echo ""
echo -e "${YELLOW}Next Steps:${NC}"
echo ""
echo -e "${GREEN}1. Add DNS record in Cloudflare Dashboard:${NC}"
echo "   - Go to: https://dash.cloudflare.com"
echo "   - Select domain: tacafe.store"
echo "   - Go to: DNS → Records"
echo "   - Click: Add record"
echo "   - Fill in:"
echo "     Type: CNAME"
echo "     Name: print"
echo "     Target: ${TUNNEL_ID}.cfargotunnel.com"
echo "     Proxy status: Proxied (orange cloud ☁️)"
echo "     TTL: Auto"
echo "   - Click: Save"
echo ""
echo -e "${GREEN}2. Start the Print Bridge:${NC}"
echo "   # Option A: Run directly"
echo "   ./print-bridge"
echo ""
echo "   # Option B: Run with Go"
echo "   go run main.go"
echo ""
echo "   # Option C: Run with Docker"
echo "   docker-compose up -d"
echo ""
echo -e "${GREEN}3. Start the Cloudflare Tunnel:${NC}"
echo "   # Option A: Run in foreground (for testing)"
echo "   cloudflared tunnel run ${TUNNEL_NAME}"
echo ""
echo "   # Option B: Run as background service (recommended)"
echo "   cloudflared service install"
echo ""
echo "   # On Linux:"
echo "   sudo systemctl start cloudflared"
echo "   sudo systemctl enable cloudflared"
echo "   sudo systemctl status cloudflared"
echo ""
echo "   # On macOS:"
echo "   sudo launchctl load /Library/LaunchDaemons/com.cloudflare.cloudflared.plist"
echo ""
echo -e "${GREEN}4. Test the tunnel:${NC}"
echo "   # Wait 30 seconds for DNS propagation, then:"
echo "   curl https://print.tacafe.store/health"
echo ""
echo "   # Expected response:"
echo '   {"status":"ok","service":"Local Print Bridge (Go + chromedp)",...}'
echo ""
echo -e "${GREEN}5. Update Backend Settings:${NC}"
echo "   - Open: https://tacafe.store/#/print-management"
echo "   - Set Print Bridge URL to: https://print.tacafe.store"
echo "   - Click: Kiểm tra kết nối"
echo "   - Click: Lưu cài đặt"
echo ""
echo -e "${YELLOW}Tunnel Information:${NC}"
echo "   Tunnel ID: ${TUNNEL_ID}"
echo "   Tunnel Name: ${TUNNEL_NAME}"
echo "   Public URL: https://print.tacafe.store"
echo "   Local Service: http://localhost:3001"
echo "   Config File: ${CONFIG_DIR}/config.yml"
echo "   Credentials: ${CREDENTIALS_FILE}"
echo ""
echo -e "${YELLOW}Useful Commands:${NC}"
echo "   # List all tunnels"
echo "   cloudflared tunnel list"
echo ""
echo "   # Check tunnel status"
echo "   cloudflared tunnel info ${TUNNEL_NAME}"
echo ""
echo "   # View tunnel logs"
echo "   cloudflared tunnel run ${TUNNEL_NAME} --loglevel debug"
echo ""
echo "   # Stop tunnel service"
echo "   sudo systemctl stop cloudflared  # Linux"
echo "   sudo launchctl unload /Library/LaunchDaemons/com.cloudflare.cloudflared.plist  # macOS"
echo ""
echo "   # Delete tunnel (if needed)"
echo "   cloudflared tunnel delete ${TUNNEL_NAME}"
echo ""
echo -e "${GREEN}✅ Setup complete! Follow the steps above to activate the tunnel.${NC}"
echo ""
