# Cache Busting Strategy - Giải Thích Chi Tiết

## Vấn Đề: Webapp Không Hiển Thị Feature Mới

Khi deploy code mới nhưng user vẫn thấy code cũ → Do browser cache!

## Cách Vite Build Process Hoạt Động

### 1. Source Code (Development)
```
frontend/src/
├── App.vue
├── main.js
├── components/
│   └── BatchManagementView.vue  ← Feature mới
└── router/index.js
```

### 2. Build Output (Production)
```bash
npm run build
```

Vite tự động tạo:
```
frontend/dist/
├── index.html                      ← KHÔNG có hash
├── assets/
│   ├── index-a3f2b8c9.js          ← Hash tự động
│   ├── index-d4e1f7a2.css         ← Hash tự động
│   └── BatchManagement-f5e8a1b3.js
```

### 3. index.html Reference
```html
<!DOCTYPE html>
<html>
<head>
  <link rel="stylesheet" href="/assets/index-d4e1f7a2.css">
</head>
<body>
  <div id="app"></div>
  <script src="/assets/index-a3f2b8c9.js"></script>
</body>
</html>
```

## Cache Busting Strategy

### Nguyên Lý
1. **Files CÓ hash** (index-abc123.js) → Cache lâu dài (1 năm)
   - Khi code thay đổi → Hash thay đổi → Browser tự động download file mới
   
2. **Files KHÔNG hash** (index.html) → KHÔNG cache
   - Luôn download mới → Luôn có reference đến file JS/CSS mới nhất

### Flow Hoạt Động

#### Deploy Lần 1 (Code Cũ)
```
Browser request: http://localhost/
↓
Nginx: index.html (no cache) → Download
↓
index.html chứa: <script src="/assets/index-a3f2b8c9.js">
↓
Browser: Download index-a3f2b8c9.js → Cache 1 năm ✅
```

#### Deploy Lần 2 (Code Mới - Có Batch Feature)
```
Browser request: http://localhost/
↓
Nginx: index.html (no cache) → Download MỚI ✅
↓
index.html chứa: <script src="/assets/index-b7d9e4f1.js">  ← Hash MỚI!
                                        ^^^^^^^^
↓
Browser: Thấy tên file khác → Download index-b7d9e4f1.js ✅
```

## Nginx Configuration Chi Tiết

```nginx
# 1. HTML - NEVER CACHE
location ~* \.html$ {
    expires -1;
    add_header Cache-Control "no-store, no-cache, must-revalidate";
}
# → index.html luôn được download mới
# → Luôn có reference đến file JS/CSS mới nhất

# 2. JS/CSS with hash - CACHE FOREVER
location ~* \.[0-9a-f]{8,}\.(js|css)$ {
    expires 1y;
    add_header Cache-Control "public, immutable";
}
# → Regex match: index.a3f2b8c9.js ✅
# → Regex match: style.d4e1f7a2.css ✅
# → Cache 1 năm vì hash đảm bảo uniqueness

# 3. Images/Fonts - CACHE 30 DAYS
location ~* \.(png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|eot)$ {
    expires 30d;
    add_header Cache-Control "public";
}

# 4. JS/CSS without hash - SHORT CACHE (fallback)
location ~* \.(js|css)$ {
    expires 1h;
    add_header Cache-Control "public, max-age=3600";
}
# → Safety net cho files không có hash
```

## Tại Sao Config Cũ Không Hoạt Động?

### Config Cũ (SAI)
```nginx
location ~* \.(js|css|png|jpg|...)$ {
    expires 1y;
    add_header Cache-Control "public, immutable";
}
```

**Vấn đề**: 
- Cache TẤT CẢ file JS/CSS trong 1 năm
- Không phân biệt file có hash hay không
- `index.html` cũng bị cache (vì match pattern `\.html$` không được exclude)

**Kết quả**:
```
Deploy code mới → index.html bị cache → Vẫn reference file JS cũ → Feature mới không hiện!
```

## Cách Deploy Với Cache Busting

### Bước 1: Update nginx.conf
```bash
# File frontend/nginx.conf đã được update với cache busting strategy
```

### Bước 2: Rebuild Frontend Image
```bash
cd frontend
npm run build  # Build code mới với hash mới
cd ..
docker build -t linhtranphu/cafe-pos-frontend:latest ./frontend
```

### Bước 3: Deploy
```bash
docker-compose -f docker-compose.prod.yml up -d frontend
```

### Bước 4: Verify
```bash
# Check cache headers
curl -I http://localhost/index.html
# Expect: Cache-Control: no-store, no-cache

curl -I http://localhost/assets/index.abc123.js
# Expect: Cache-Control: public, immutable
```

## Test Cache Busting

### Test 1: HTML Không Cache
```bash
# Request 2 lần
curl -I http://localhost/index.html
curl -I http://localhost/index.html

# Expect: Không có "304 Not Modified"
# → Luôn download mới
```

### Test 2: JS/CSS Có Hash Cache Lâu
```bash
curl -I http://localhost/assets/index.a3f2b8c9.js

# Expect:
# Cache-Control: public, immutable
# Expires: [1 năm sau]
```

### Test 3: Deploy Code Mới
```bash
# 1. Build code mới
npm run build

# 2. Check hash mới
ls frontend/dist/assets/
# → index-NEW_HASH.js (hash khác với cũ)

# 3. Deploy
docker build -t linhtranphu/cafe-pos-frontend:latest ./frontend
docker-compose -f docker-compose.prod.yml up -d frontend

# 4. Browser tự động load file mới vì hash khác!
```

## Lợi Ích Cache Busting

### ✅ Performance
- Files có hash cache lâu → Giảm bandwidth
- Chỉ download file thay đổi → Nhanh hơn

### ✅ Reliability
- Luôn load code mới nhất
- Không cần clear cache thủ công
- Không cần version query string (?v=1.2.3)

### ✅ User Experience
- User tự động nhận update
- Không cần hard refresh (Ctrl+F5)
- Không cần hướng dẫn clear cache

## So Sánh Các Phương Pháp

### 1. No Cache (Option 1) ❌
```nginx
location ~* \.(js|css)$ {
    expires -1;
    add_header Cache-Control "no-cache";
}
```
- ❌ Luôn download lại → Chậm
- ❌ Tốn bandwidth
- ✅ Luôn có code mới

### 2. Cache Busting (Option 2) ✅ RECOMMENDED
```nginx
# HTML: no cache
# JS/CSS with hash: cache forever
```
- ✅ Nhanh (cache files không đổi)
- ✅ Tiết kiệm bandwidth
- ✅ Luôn có code mới (khi hash thay đổi)

### 3. Query String Versioning ⚠️
```html
<script src="/app.js?v=1.2.3"></script>
```
- ⚠️ Phải manual update version
- ⚠️ Proxy có thể ignore query string
- ⚠️ Không tự động

## Troubleshooting

### Vấn đề: Vẫn thấy code cũ sau khi deploy

**Nguyên nhân có thể**:
1. Browser cache index.html
2. Service Worker cache
3. CDN/Proxy cache

**Giải pháp**:
```bash
# 1. Hard refresh browser
Ctrl + Shift + R (Windows/Linux)
Cmd + Shift + R (Mac)

# 2. Check nginx config
docker exec cafe-pos-frontend cat /etc/nginx/conf.d/default.conf

# 3. Check build output
docker exec cafe-pos-frontend ls -la /usr/share/nginx/html/assets/

# 4. Verify cache headers
curl -I http://localhost/index.html
curl -I http://localhost/assets/index.*.js
```

### Vấn đề: Hash không thay đổi

**Nguyên nhân**: Code không thực sự thay đổi

**Giải pháp**:
```bash
# Force rebuild
rm -rf frontend/dist
npm run build
```

## Kết Luận

Cache Busting với Vite hash là best practice vì:
1. Tự động (không cần manual versioning)
2. Hiệu quả (cache files không đổi, load files mới khi cần)
3. Đơn giản (chỉ cần config nginx đúng)

**Không cần clear cache thủ công nữa!** 🎉
