# Chromedp HTML Capture Fix

## Problem Analysis

Based on research and testing, the issues with chromedp HTML capture are:

1. **Logo not rendering**: Base64 images in HTML need proper encoding and may have size limits
2. **Layout misalignment**: Need proper viewport settings and wait time for rendering
3. **File vs Data URL**: Different methods have different behaviors

## Research Findings

### From chromedp/examples and Stack Overflow:

1. **Best Practice**: Use `file://` URL with temp HTML file (most reliable)
2. **Wait Time**: Need sufficient time for images to load (2+ seconds recommended)
3. **Viewport Size**: Should set explicit window size matching target width
4. **Base64 Images**: Work in `file://` URLs but may fail in data URLs due to size limits

### Test Results:

Created `backend/cmd/test-chromedp-html-render/main.go` to test 3 methods:
- ✅ Method 1 (file:// URL): 20KB output - **WORKS WITH LOGO**
- ❌ Method 2 (data URL): 4KB output - **LOGO NOT RENDERING**  
- ✅ Method 3 (SetDocumentContent): 20KB output - **WORKS WITH LOGO**

**Conclusion**: Current implementation using Method 1 (file:// URL) is correct!

## Root Cause

The issue is NOT the chromedp method, but likely:

1. **Logo path incorrect**: Shop settings may have wrong logo URL
2. **Base64 encoding issue**: Logo may not be loading into base64 properly
3. **Wait time too short**: 1 second may not be enough for large base64 images
4. **CSS layout issues**: Flexbox or image sizing may need adjustment

## Solution

### Fix 1: Increase Wait Time

Current: 1 second
Recommended: 2-3 seconds for base64 images to decode

```go
chromedp.Sleep(2*time.Second), // Increased from 1s
```

### Fix 2: Add Explicit Viewport

```go
opts := append(chromedp.DefaultExecAllocatorOptions[:],
    chromedp.Flag("headless", true),
    chromedp.Flag("disable-gpu", true),
    chromedp.Flag("no-sandbox", true),
    chromedp.WindowSize(576, 1200), // ADD THIS
)
```

### Fix 3: Wait for Image Load Event

Instead of just Sleep, wait for images to actually load:

```go
chromedp.Navigate(fileURL),
chromedp.WaitReady("body"), // Wait for body
chromedp.Sleep(500*time.Millisecond),
chromedp.Evaluate(`
    new Promise((resolve) => {
        const images = document.querySelectorAll('img');
        if (images.length === 0) {
            resolve();
            return;
        }
        let loaded = 0;
        images.forEach(img => {
            if (img.complete) {
                loaded++;
            } else {
                img.onload = () => {
                    loaded++;
                    if (loaded === images.length) resolve();
                };
                img.onerror = () => {
                    loaded++;
                    if (loaded === images.length) resolve();
                };
            }
        });
        if (loaded === images.length) resolve();
        // Timeout after 3 seconds
        setTimeout(resolve, 3000);
    });
`, nil),
chromedp.FullScreenshot(&buf, 100),
```

### Fix 4: Debug Logo Loading

Add logging to verify logo is actually loaded:

```go
func (r *ChromedpBillRendererOptimized) prepareBillData(...) {
    // ... existing code ...
    
    if shopSettings.ShowLogo && shopSettings.LogoURL != "" {
        logoBase64, err := loadImageAsBase64(shopSettings.LogoURL)
        if err != nil {
            log.Printf("ERROR: Failed to load logo from %s: %v", shopSettings.LogoURL, err)
        } else {
            log.Printf("SUCCESS: Logo loaded: %d bytes (first 50 chars: %s...)", 
                len(logoBase64), logoBase64[:min(50, len(logoBase64))])
            data.LogoBase64 = logoBase64
        }
    } else {
        log.Printf("INFO: Logo disabled or no URL (ShowLogo=%v, LogoURL=%s)", 
            shopSettings.ShowLogo, shopSettings.LogoURL)
    }
    
    // ... rest of code ...
}
```

### Fix 5: CSS Improvements

Ensure logo container has proper sizing:

```html
<div class="logo">
    <img src="{{.LogoBase64}}" alt="Logo" style="display: block; max-width: 100%; height: auto;">
</div>
```

## Implementation Plan

1. ✅ Test chromedp methods (DONE - Method 1 works)
2. ⚠️ Add debug logging to verify logo loading
3. ⚠️ Increase wait time to 2-3 seconds
4. ⚠️ Add WindowSize to Chrome options
5. ⚠️ Optionally add image load wait
6. ⚠️ Test with real order and verify output

## Files to Modify

1. `backend/application/services/chromedp_bill_renderer_optimized.go`
   - Add WindowSize to Chrome options
   - Increase Sleep time
   - Add debug logging
   - Optionally add image load wait

2. `backend/application/services/templates/bill_template_optimized.html`
   - Add inline style to img tag for safety
   - Verify CSS is correct

## Expected Result

After fixes:
- Logo should render correctly in preview PNG
- Layout should match HTML template exactly
- Print output should show logo and proper formatting

## Testing

```bash
# 1. Create preview
curl -X POST \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"order_id":"YOUR_ORDER_ID"}' \
  http://localhost:3000/api/manager/html-templates/preview

# 2. Check files
ls -lh backend/raw_preview_*.png backend/preview_*.png

# 3. Open and verify
open backend/raw_preview_*.png  # Should show logo and proper layout

# 4. Test print
curl -X POST \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"order_id":"YOUR_ORDER_ID","printer_ip":"192.168.1.115"}' \
  http://localhost:3000/api/manager/html-templates/test-print
```

## Next Steps

1. Apply fixes to chromedp_bill_renderer_optimized.go
2. Test preview generation
3. Verify logo appears in preview PNG
4. Test actual print
5. Adjust threshold/settings if needed
