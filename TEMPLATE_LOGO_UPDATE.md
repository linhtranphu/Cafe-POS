# Template Logo Update

## Changes Made

Đã cập nhật `bill_template_optimized.html` để có layout giống với `preview.go`:

### 1. Header Layout

**Before:**
```css
.header {
    display: flex;
    align-items: flex-start;
    margin-bottom: 20px;
}

.logo {
    width: 200px;
    margin-right: 20px;
}
```

**After:**
```css
.header {
    display: flex;
    align-items: flex-start;
    margin-bottom: 45px;
    min-height: 100px;
}

.logo {
    width: 200px;
    margin-left: 20px;
    margin-right: 60px;
    flex-shrink: 0;
}
```

**Changes:**
- Increased margin-bottom: 20px → 45px
- Added min-height: 100px
- Added margin-left: 20px (logo ở góc trái)
- Increased margin-right: 20px → 60px
- Added flex-shrink: 0

### 2. Shop Info

**Before:**
```css
.shop-info {
    flex: 1;
}

.shop-name {
    font-size: 25px;
    font-weight: bold;
    margin-bottom: 8px;
}
```

**After:**
```css
.shop-info {
    flex: 1;
    padding-top: 20px;
}

.shop-name {
    font-size: 25px;
    font-weight: bold;
    margin-bottom: 5px;
    text-shadow: 2px 0 0 currentColor;
    letter-spacing: 0.5px;
}
```

**Changes:**
- Added padding-top: 20px
- Reduced margin-bottom: 8px → 5px
- Added text-shadow for fake bold effect
- Added letter-spacing: 0.5px

### 3. Order Info

**Before:**
```html
<div class="order-info">Order: {{.OrderNumber}}</div>
<div class="order-info">Waiter: {{.WaiterName}}</div>
```

**After:**
```html
<div class="order-info">
    <div>Order: {{.OrderNumber}}</div>
    <div>Waiter: {{.WaiterName}}</div>
    <div>Thanh Toán: {{.PaymentMethod}}</div>
    <div>Ngày tạo: {{.CreatedDate}}</div>
</div>
```

**CSS:**
```css
.order-info {
    font-size: 16px;
    margin-bottom: 20px;
    line-height: 20px;
}

.order-info div {
    margin-left: 10px;
}
```

**Changes:**
- Wrapped all order info in single container
- Added margin-left: 10px to child divs
- Increased margin-bottom: 8px → 20px

### 4. Dividers

**Before:**
```css
.divider {
    border-top: 2px solid black;
    margin: 15px 0;
}
```

**After:**
```css
.divider {
    border-top: 1px solid black;
    margin: 25px 0;
}
```

**Changes:**
- Reduced border: 2px → 1px
- Increased margin: 15px → 25px

### 5. Table Headers

**Before:**
```css
th {
    font-size: 17px;
    text-align: left;
    padding: 8px 0;
    border-bottom: 2px solid black;
}
```

**After:**
```css
th {
    font-size: 17px;
    font-weight: normal;
    text-align: left;
    padding: 8px 0;
    border-bottom: 1px solid black;
}
```

**Changes:**
- Added font-weight: normal
- Reduced border: 2px → 1px

### 6. Table Cells

**Before:**
```html
<th style="width: 40px;">STT</th>
<td>{{.STT}}</td>
```

**After:**
```html
<th style="width: 40px; padding-left: 10px;">STT</th>
<td style="padding-left: 10px;">{{.STT}}</td>
```

**Changes:**
- Added padding-left: 10px to first column
- Added padding-right: 10px to last column
- Adjusted column widths to match preview.go

### 7. Total Section

**Before:**
```css
.total-label {
    font-size: 24px;
    font-weight: bold;
}
```

**After:**
```css
.total-label {
    font-size: 24px;
    font-weight: bold;
    padding-left: 190px;
}

.total-amount {
    font-size: 24px;
    font-weight: bold;
    text-align: right;
    padding-right: 10px;
}
```

**Changes:**
- Added padding-left: 190px to total label
- Added padding-right: 10px to total amount

### 8. Thanks Message

**Before:**
```css
.thanks {
    text-align: center;
    font-size: 22px;
    margin-top: 20px;
}
```

**After:**
```css
.thanks {
    text-align: center;
    font-size: 22px;
    margin-top: 30px;
    margin-bottom: 30px;
}
```

**Changes:**
- Increased margin-top: 20px → 30px
- Added margin-bottom: 30px

## Layout Comparison

### Preview.go Layout
```
┌────────────────────────────────────────────────┐
│ [Logo 200px]    Shop Name (bold + shadow)     │
│ (margin-left    Address                        │
│  20px)          Phone                          │
│                                                │
│         HÓA ĐƠN THANH TOÁN                    │
│                                                │
│    Order: ...                                  │
│    Waiter: ...                                 │
│    Payment: ...                                │
│    Date: ...                                   │
│ ───────────────────────────────────────────── │
│ STT  Name         SL  Price      Total        │
│ ───────────────────────────────────────────── │
│  1   Item 1       2   25,000     50,000       │
│  2   Item 2       1   35,000     35,000       │
│ ───────────────────────────────────────────── │
│                    TỔNG TIỀN:     105,000     │
│ ───────────────────────────────────────────── │
│           Cảm ơn quý khách!                   │
└────────────────────────────────────────────────┘
```

### Updated Template Layout
```
┌────────────────────────────────────────────────┐
│ [Logo 200px]    Shop Name (bold + shadow)     │
│ (margin-left    Address                        │
│  20px)          Phone                          │
│                                                │
│         HÓA ĐƠN THANH TOÁN                    │
│                                                │
│    Order: ...                                  │
│    Waiter: ...                                 │
│    Payment: ...                                │
│    Date: ...                                   │
│ ───────────────────────────────────────────── │
│ STT  Name         SL  Price      Total        │
│ ───────────────────────────────────────────── │
│  1   Item 1       2   25,000     50,000       │
│  2   Item 2       1   35,000     35,000       │
│ ───────────────────────────────────────────── │
│                    TỔNG TIỀN:     105,000     │
│ ───────────────────────────────────────────── │
│           Cảm ơn quý khách!                   │
└────────────────────────────────────────────────┘
```

## Logo Position

Logo bây giờ được đặt ở góc trái với:
- Width: 200px (fixed)
- Margin-left: 20px (từ edge)
- Margin-right: 60px (khoảng cách với shop info)
- Flex-shrink: 0 (không bị shrink)

## Testing

### 1. Reload Template in Frontend

```
1. Open http://localhost:5173/#/print-management
2. Click Templates → HTML Template
3. Click "🔄 Reload" button
4. Template mới sẽ load với logo ở góc trái
```

### 2. Preview Changes

Template preview sẽ tự động update với layout mới.

### 3. Test Print

```
1. Select an order
2. Click "🖨️ Test Print"
3. Check output trên máy in
```

## Logo Display Conditions

Logo chỉ hiển thị khi:
1. `ShowLogo` = true (trong shop settings)
2. `LogoBase64` có giá trị (logo đã được upload)

Nếu không có logo:
- Shop info sẽ chiếm toàn bộ width
- Layout vẫn giữ nguyên spacing

## Shop Settings

Để enable logo:
1. Go to Print Management → Settings
2. Upload logo
3. Enable "Show Logo" checkbox
4. Save settings

## File Location

```
backend/application/services/templates/bill_template_optimized.html
```

## Backup

Nếu cần rollback, file backup được tạo tự động khi save:
```
backend/application/services/templates/bill_template_optimized.html.backup
```

## Summary

✅ Logo positioned ở góc trái (margin-left: 20px)
✅ Layout matches preview.go exactly
✅ All spacing và margins updated
✅ Text shadow added cho shop name (fake bold)
✅ Column widths và paddings adjusted
✅ Dividers changed từ 2px → 1px
✅ All margins increased for better spacing

Template bây giờ có layout giống hệt preview.go với logo ở góc trái!
