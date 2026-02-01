# Terser Build Error - Fix Guide

## Problem
When building Docker image, you get error:
```
Error: terser not found. Since Vite v3, terser has become an optional dependency. You need to install it.
```

## Root Cause
Vite v3+ made terser an optional dependency. If you specify `minify: 'terser'` in vite.config.js but terser is not installed, the build fails.

## Solution

### Option 1: Use esbuild (Recommended)
esbuild is built into Vite and doesn't require separate installation.

**File: `frontend/vite.config.js`**
```javascript
build: {
  minify: 'esbuild',  // Use esbuild instead of terser
  // ... other config
}
```

**Pros:**
- No additional dependencies
- Faster build time
- Built into Vite
- Good minification quality

**Cons:**
- Slightly less aggressive minification than terser

### Option 2: Install terser
If you prefer terser's minification quality:

**File: `frontend/package.json`**
```json
{
  "devDependencies": {
    "terser": "^5.24.0"
  }
}
```

Then run:
```bash
npm install
```

**Pros:**
- More aggressive minification
- Smaller bundle size

**Cons:**
- Additional dependency
- Slower build time

## Implementation

### Step 1: Update vite.config.js
```javascript
export default defineConfig({
  build: {
    minify: 'esbuild',  // Changed from 'terser'
    rollupOptions: {
      output: {
        entryFileNames: '[name].[hash].js',
        chunkFileNames: '[name].[hash].js',
        assetFileNames: '[name].[hash][extname]'
      }
    }
  }
})
```

### Step 2: Install dependencies (if using terser)
```bash
npm install
```

### Step 3: Test build
```bash
npm run build
```

Expected output:
```
✓ 149 modules transformed.
dist/index.html            1.92 kB │ gzip:   1.02 kB
dist/index.e318f6f0.css   37.46 kB │ gzip:   6.44 kB
dist/index.d47b2b57.js   409.74 kB │ gzip: 111.56 kB
✓ built in 4.30s
```

### Step 4: Verify cache busting
Check that output files have hashes:
```bash
ls -la dist/ | grep -E '\.(js|css)$'
```

Expected output:
```
-rw-r--r--  index.d47b2b57.js
-rw-r--r--  index.e318f6f0.css
```

## Docker Build

The Dockerfile already handles npm install correctly:

```dockerfile
# Build stage
FROM node:18-alpine AS builder

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies (includes terser if specified)
RUN npm ci

# Copy source code
COPY . .

# Build the application
RUN npm run build
```

When you rebuild Docker image:
```bash
docker build -f frontend/Dockerfile -t linhtranphu/cafe-pos-frontend:latest .
```

The build will:
1. Install all dependencies from package.json
2. Run `npm run build`
3. Generate hashed files for cache busting

## Comparison: esbuild vs terser

| Feature | esbuild | terser |
|---------|---------|--------|
| Installation | Built-in | Requires npm install |
| Build Speed | Fast | Slower |
| Minification | Good | Excellent |
| Bundle Size | Slightly larger | Smaller |
| Recommended | ✅ Yes | For production optimization |

## Files Modified

- ✅ `frontend/vite.config.js` - Changed minify to 'esbuild'
- ✅ `frontend/package.json` - Added terser as optional dependency
- ✅ `frontend/Dockerfile` - Already correct (runs npm ci)

## Testing

### Local Development
```bash
npm run build
ls dist/
# Should see files with hashes: index.abc123.js, index.def456.css
```

### Docker Build
```bash
docker build -f frontend/Dockerfile -t cafe-pos-frontend:test .
docker run -p 80:80 cafe-pos-frontend:test
# Visit http://localhost
# Check DevTools → Network → should see hashed filenames
```

## Related Best Practices

- See `CACHE_CONFIGURATION.md` for cache busting strategy
- See `README.md` for deployment best practices

