# Solution: Use Visual Renderer Instead of Chromedp

## Problem Summary

1. Chromedp HTML snapshot có vấn đề với logo rendering
2. Convert từ HTML screenshot sang ESC/POS in ra ký tự vô nghĩa
3. Auto-print gây phiền toái

## Root Cause

Chromedp approach có nhiều vấn đề:
- Base64 images trong HTML có thể không render đúng
- Screenshot quality không ổn định
- Grayscale/binarization conversion phức tạp
- ESC/POS conversion từ screenshot không reliable

## Recommended Solution

**Dùng `visual_bill_renderer.go` thay vì `chromedp_bill_renderer_optimized.go`**

### Why Visual Renderer is Better:

1. **Direct Drawing**: Dùng `gg` library vẽ trực tiếp, không qua browser
2. **Proven to Work**: Đã test và hoạt động tốt với logo
3. **Better Control**: Kiểm soát chính xác từng pixel
4. **Simpler**: Không cần chromedp, không cần HTML parsing
5. **Faster**: Không cần khởi động browser

### Implementation Plan:

#### Option A: Modify HTML Template Handler to Use Visual Renderer

```go
// In html_template_handler.go
func (h *HTMLTemplateHandler) TestPrint(c *gin.Context) {
    // ... existing code ...
    
    // CHANGE: Use visual renderer instead of chromedp
    visualRenderer, err := services.NewVisualBillRenderer()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "error": "Failed to create renderer: " + err.Error(),
        })
        return
    }
    
    escposData, err := visualRenderer.RenderBillToESCPOS(ord, shopSettings)
    // ... rest of code ...
}
```

#### Option B: Create Hybrid Approach

- Keep HTML template for **preview only** (display in browser)
- Use visual renderer for **actual printing**

```go
// Preview: Use HTML template
func (h *HTMLTemplateHandler) Preview(c *gin.Context) {
    // Use chromedp to generate PNG preview
    // This is OK because it's just for display
}

// Print: Use visual renderer
func (h *HTMLTemplateHandler) TestPrint(c *gin.Context) {
    // Use visual renderer for actual printing
    // This ensures reliable output
}
```

## Implementation Steps

### Step 1: Disable Auto-Print (DONE)
✅ Commented out auto-print in `order_handler.go`

### Step 2: Update HTML Template Handler

Modify `backend/interfaces/http/html_template_handler.go`:

```go
func (h *HTMLTemplateHandler) TestPrint(c *gin.Context) {
    var req TestPrintRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Get order
    ord, err := h.orderRepo.GetByID(c.Request.Context(), req.OrderID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
        return
    }

    // Get shop settings
    shopSettings, err := h.settingsRepo.GetSettings(c.Request.Context())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get settings"})
        return
    }

    // USE VISUAL RENDERER INSTEAD OF CHROMEDP
    visualRenderer, err := services.NewVisualBillRenderer()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "error": "Failed to create renderer: " + err.Error(),
        })
        return
    }

    // Render to ESC/POS
    escposData, err := visualRenderer.RenderBillToESCPOS(ord, shopSettings)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "error": "Failed to render bill: " + err.Error(),
        })
        return
    }

    // Send to printer
    printerIP := req.PrinterIP
    if printerIP == "" {
        printerIP = "192.168.1.115:9100"
    }

    if err := sendToPrinter(printerIP, escposData); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "error": "Failed to send to printer: " + err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "message": "Test print successful",
        "order_number": ord.OrderNumber,
    })
}
```

### Step 3: Keep Preview Using Chromedp (Optional)

Preview can still use chromedp for HTML rendering since it's just for display:

```go
func (h *HTMLTemplateHandler) Preview(c *gin.Context) {
    // Keep existing chromedp preview code
    // This is fine for display purposes
}
```

## Benefits

1. ✅ Logo sẽ in ra đúng (visual renderer đã test OK)
2. ✅ Không còn ký tự vô nghĩa (ESC/POS conversion đã proven)
3. ✅ Faster rendering (no browser overhead)
4. ✅ More reliable (direct drawing)
5. ✅ Simpler code (less dependencies)

## Trade-offs

- ❌ Mất tính năng customize HTML template cho printing
- ✅ Nhưng vẫn giữ được HTML template cho preview/display
- ✅ Visual renderer có thể customize qua code

## Alternative: Fix Chromedp Approach

Nếu muốn giữ chromedp approach, cần:

1. Research thêm về chromedp image rendering
2. Test với các browser flags khác nhau
3. Có thể cần dùng Rod thay vì chromedp
4. Debug ESC/POS conversion chi tiết hơn

Nhưng approach này phức tạp hơn và không chắc chắn thành công.

## Recommendation

**Implement Option B (Hybrid Approach)**:
- HTML template + chromedp cho **preview** (display only)
- Visual renderer cho **actual printing** (reliable output)

This gives best of both worlds:
- Users can see HTML preview
- Printing is reliable and works with logo
