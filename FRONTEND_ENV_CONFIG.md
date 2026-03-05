# Frontend Environment Configuration

## Files

- `.env` - Development (local testing)
- `.env.production` - Production (deployed on EC2)
- `.env.example` - Template/documentation

## Current Configuration

### Development (.env)
```env
VITE_API_URL=http://localhost:3000
```

### Production (.env.production)
```env
VITE_API_URL=https://tacafe.store
VITE_PRINT_BRIDGE_URL=http://localhost:3001
```

## Print Bridge Configuration

### Scenario 1: Print Bridge on Same Machine as Browser

**Current setup** - Browser và print bridge cùng máy:

```env
VITE_PRINT_BRIDGE_URL=http://localhost:3001
```

✅ Works when: Staff accesses from the machine running print bridge

❌ Doesn't work when: Staff accesses from different device

---

### Scenario 2: Print Bridge on Different Machine (Same Network)

**Recommended for cafe** - Print bridge chạy trên máy riêng:

1. **Find print bridge machine IP:**
```bash
# On print bridge machine
ipconfig  # Windows
ifconfig  # macOS/Linux

# Example output: 192.168.1.100
```

2. **Update `.env.production`:**
```env
VITE_API_URL=https://tacafe.store
VITE_PRINT_BRIDGE_URL=http://192.168.1.100:3001
```

3. **Rebuild frontend:**
```bash
cd frontend
npm run build
```

4. **Rebuild Docker image:**
```bash
docker build -t linhtranphu/cafe-pos-frontend:2.0.1 .
docker push linhtranphu/cafe-pos-frontend:2.0.1
```

5. **Deploy on EC2:**
```bash
ssh ubuntu@your-ec2
sudo docker-compose pull frontend
sudo docker-compose up -d frontend
```

✅ Works when: Staff accesses from any device on same network

❌ Doesn't work when: Accessed from outside network

⚠️ **Note:** May show "Mixed Content" warning (HTTPS → HTTP)

---

### Scenario 3: Print Bridge via Cloudflare Tunnel (Most Flexible)

**Best for production** - Accessible from anywhere:

1. **Setup Cloudflare Tunnel** (see PRINT_BRIDGE_CONFIG.md)

2. **Update `.env.production`:**
```env
VITE_API_URL=https://tacafe.store
VITE_PRINT_BRIDGE_URL=https://print.tacafe.store
```

3. **Rebuild and deploy** (same as Scenario 2)

✅ Works from anywhere
✅ No mixed content warning (HTTPS)
✅ Secure

---

## Quick Setup Guide

### For Local Network Access (Recommended)

1. **Find print bridge IP:**
```bash
ipconfig  # On print bridge machine
```

2. **Edit `frontend/.env.production`:**
```env
VITE_PRINT_BRIDGE_URL=http://YOUR_IP:3001
```

3. **Rebuild:**
```bash
./build_docker_hub.sh
```

4. **Deploy:**
```bash
# On EC2
sudo docker-compose pull
sudo docker-compose up -d
```

---

## Verification

### Check Environment Variables in Browser

1. Open https://tacafe.store
2. Open browser console (F12)
3. Type:
```javascript
import.meta.env.VITE_API_URL
import.meta.env.VITE_PRINT_BRIDGE_URL
```

Should show:
```
"https://tacafe.store"
"http://192.168.1.100:3001"  // or your configured URL
```

### Check Print Bridge Connection

1. Go to Print Management page
2. Check console for:
```
[LocalPrint] Bridge available: true
```

If false, check:
- Print bridge is running
- IP/URL is correct
- Firewall allows port 3001
- CORS is configured

---

## Troubleshooting

### "Bridge not available"

**Check 1:** Can access from browser?
```
http://192.168.1.100:3001/health
```

**Check 2:** Print bridge running?
```bash
# On print bridge machine
curl http://localhost:3001/health
```

**Check 3:** Firewall blocking?
- Windows: Allow port 3001 in Windows Firewall
- macOS: System Preferences → Security → Firewall

### Mixed Content Warning

Browser console shows:
```
Mixed Content: The page at 'https://tacafe.store' was loaded over HTTPS,
but requested an insecure resource 'http://192.168.1.100:3001'
```

**Solutions:**
1. Use Cloudflare Tunnel (HTTPS)
2. Allow insecure content in browser (not recommended)
3. Access via HTTP instead of HTTPS (not recommended)

### Wrong API URL

If API calls fail, check `.env.production`:
```env
VITE_API_URL=https://tacafe.store  # Must match your domain
```

---

## Build Process

### Development Build
```bash
npm run dev
# Uses .env
```

### Production Build
```bash
npm run build
# Uses .env.production
```

### Docker Build
```bash
docker build -t linhtranphu/cafe-pos-frontend:2.0.1 .
# Copies .env.production into image
```

---

## Important Notes

1. **Environment variables are embedded at build time**
   - Changing `.env.production` requires rebuild
   - Cannot change after deployment without rebuild

2. **Security**
   - Don't put secrets in .env files (they're in frontend code)
   - Use for configuration only

3. **Git**
   - `.env` is gitignored
   - `.env.production` should be committed (no secrets)
   - `.env.example` is for documentation

---

## Summary

**Current setup:**
- Development: `localhost:3000` (backend), `localhost:3001` (print bridge)
- Production: `tacafe.store` (backend), `localhost:3001` (print bridge)

**To change print bridge URL:**
1. Edit `frontend/.env.production`
2. Run `./build_docker_hub.sh`
3. Deploy on EC2

**Recommended for cafe:**
```env
VITE_PRINT_BRIDGE_URL=http://192.168.1.100:3001
```
(Replace with actual IP of print bridge machine)
