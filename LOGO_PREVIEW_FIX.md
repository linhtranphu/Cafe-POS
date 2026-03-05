# Logo Preview Fix in Print Management

## Problem
Logo không hiển thị trong preview ở màn hình Print Management (`http://localhost:5173/#/print-management`) mặc dù backend đã load được logo base64.

## Root Causes

### 1. Frontend không load shop settings
- `HTMLTemplateEditor.vue` chỉ load template HTML
- Không load shop settings để lấy logo
- `sampleData` có `ShowLogo: false` và `LogoBase64: ''`

### 2. Template processing không đúng
- Go template syntax: `{{range .Items}}` nhưng code xử lý `{{range $index, $item := .Items}}`
- Conditional `{{if .ShowCustomMessage}}` nhưng code check `{{if .ShowCustomMsg}}`
- Missing `{{.LogoBase64}}` replacement

### 3. Vite proxy không forward `/uploads`
- Vite chỉ proxy `/api` → backend
- `/uploads` không được proxy → 404 khi fetch logo

## Solutions

### 1. Added `loadShopSettings()` function
```javascript
const loadShopSettings = async () => {
  try {
    const response = await api.get('/manager/shop-settings')
    const settings = response.data
    
    // Update sample data with real shop settings
    sampleData.value.ShopName = settings.shop_name
    sampleData.value.ShopAddress = settings.shop_address
    sampleData.value.ShopPhone = settings.shop_phone
    sampleData.value.ShowLogo = settings.show_logo
    sampleData.value.ShowAddress = settings.show_address
    sampleData.value.ShowPhone = settings.show_phone
    sampleData.value.CustomMessage = settings.custom_message
    
    // Load logo as base64
    if (settings.show_logo && settings.logo_url) {
      const logoResponse = await fetch(settings.logo_url)
      const logoBlob = await logoResponse.blob()
      const reader = new FileReader()
      
      reader.onloadend = () => {
        sampleData.value.LogoBase64 = reader.result
        updatePreview()
      }
      
      reader.readAsDataURL(logoBlob)
    }
  } catch (error) {
    console.error('Failed to load shop settings:', error)
  }
}
```

### 2. Fixed `processTemplate()` function
**Before:**
```javascript
// Wrong: Always remove logo
processed = processed.replace(/\{\{if \.ShowLogo\}\}[\s\S]*?\{\{end\}\}/g, '')

// Wrong: Old template syntax
const rangeMatch = processed.match(/\{\{range \$index, \$item := \.Items\}\}/)
```

**After:**
```javascript
// Correct: Show logo if available
if (sampleData.value.ShowLogo && sampleData.value.LogoBase64) {
  processed = processed.replace(/\{\{if \.ShowLogo\}\}([\s\S]*?)\{\{end\}\}/g, '$1')
} else {
  processed = processed.replace(/\{\{if \.ShowLogo\}\}[\s\S]*?\{\{end\}\}/g, '')
}

// Correct: Match actual template syntax
const rangeMatch = processed.match(/\{\{range \.Items\}\}/)

// Add logo base64 replacement
processed = processed.replace(/\{\{\.LogoBase64\}\}/g, sampleData.value.LogoBase64)
```

### 3. Added Vite proxy for `/uploads`
**File**: `frontend/vite.config.js`

```javascript
server: {
  port: 5173,
  proxy: {
    '/api': {
      target: 'http://localhost:3000',
      changeOrigin: true,
      rewrite: (path) => path
    },
    '/uploads': {  // NEW
      target: 'http://localhost:3000',
      changeOrigin: true,
      rewrite: (path) => path
    }
  }
}
```

### 4. Updated `sampleData` structure
```javascript
const sampleData = ref({
  // ... other fields
  ShowLogo: false,        // Will be updated from shop settings
  LogoBase64: '',         // Will be loaded from backend
  // ... other fields
})
```

## How It Works Now

1. **On component mount**:
   - Load HTML template from backend
   - Load shop settings (including logo URL)
   - Fetch logo file from `/uploads/logos/...` (proxied to backend)
   - Convert logo to base64 using FileReader
   - Update preview

2. **Template processing**:
   - Replace `{{.LogoBase64}}` with actual base64 data
   - Show/hide logo section based on `ShowLogo` flag
   - Match actual Go template syntax from `bill_template_optimized.html`

3. **Preview rendering**:
   - Process template with real data
   - Inject into iframe
   - Logo displays correctly

## Files Modified
- `frontend/src/components/printing/HTMLTemplateEditor.vue`
  - Added `loadShopSettings()` function
  - Fixed `processTemplate()` to handle logo correctly
  - Updated `sampleData` structure
  - Call `loadShopSettings()` in `loadTemplate()`
  
- `frontend/vite.config.js`
  - Added `/uploads` proxy configuration

## Testing
1. Open `http://localhost:5173/#/print-management`
2. Go to "Templates" tab
3. Logo should appear in preview (if `show_logo: true` in settings)
4. Check browser console for "Logo loaded successfully" message

## Next Steps
- Verify logo displays in actual print output
- Test with different logo sizes
- Test with PNG vs JPEG logos
