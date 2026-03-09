# Fix 404 Frontend Error

## Vấn đề
```
GET https://tacafe.store/assets/FundManagementView-1a6e0574.js net::ERR_ABORTED 404
```

## Nguyên nhân
- Browser cache file JS cũ với hash cũ
- File JS cũ reference đến chunk files với hash cũ
- Sau khi build mới, hash thay đổi → 404

## Giải pháp

### Option 1: Quick Fix (Recommended)
Chỉ cần clear browser cache trên client:

**Desktop:**
- Chrome/Edge: `Ctrl + Shift + R` (Windows) hoặc `Cmd + Shift + R` (Mac)
- Firefox: `Ctrl + Shift + Delete` → Clear cache

**Mobile:**
- Settings → Clear browsing data → Cached images and files

### Option 2: Rebuild & Redeploy (Nếu Option 1 không work)

**Bước 1: Rebuild và push image mới**
```bash
./fix-404-frontend.sh
```

**Bước 2: Update trên server**
```bash
ssh user@tacafe.store
cd /path/to/project
docker-compose pull frontend
docker-compose up -d frontend
```

**Bước 3: Clear browser cache**
- Hard refresh: `Ctrl + Shift + R` hoặc `Cmd + Shift + R`

### Option 3: Manual Rebuild

```bash
# 1. Build frontend
cd frontend
rm -rf dist
npm run build
cd ..

# 2. Build Docker image
docker build -t linhtranphu/cafe-pos-frontend:latest -f frontend/Dockerfile frontend/

# 3. Push to Docker Hub
docker push linhtranphu/cafe-pos-frontend:latest

# 4. Update on server
ssh user@tacafe.store
docker-compose pull frontend
docker-compose up -d frontend
```

## Phòng tránh trong tương lai

### 1. Nginx đã config cache busting
- HTML files: NEVER cache
- JS/CSS với hash: Cache 1 year
- Assets: Cache 30 days

### 2. Sau mỗi lần deploy
- Luôn clear browser cache
- Hoặc dùng incognito mode để test

### 3. Version tagging
Khi build, có thể tag với timestamp:
```bash
TIMESTAMP=$(date +%Y%m%d-%H%M%S)
docker build -t linhtranphu/cafe-pos-frontend:$TIMESTAMP frontend/
docker tag linhtranphu/cafe-pos-frontend:$TIMESTAMP linhtranphu/cafe-pos-frontend:latest
```

## Kiểm tra

Sau khi fix, verify:
1. Open DevTools (F12)
2. Network tab → Disable cache
3. Hard refresh (Ctrl+Shift+R)
4. Check tất cả JS files load thành công (200 OK)
