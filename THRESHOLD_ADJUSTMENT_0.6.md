# Threshold Adjustment to 0.6

## Change Summary
Điều chỉnh threshold từ 0.5 lên 0.6 để in nhạt hơn (ít pixel đen hơn).

## Reason
Print output quá đậm/đen so với preview. Cần giảm số lượng pixel đen để match với preview tốt hơn.

## Changes Made

### chromedp_bill_renderer_optimized.go
```go
converter := &raster.Converter{
    MaxWidth:  576,
    Threshold: 0.6,  // Changed from 0.5 to 0.6
}
```

### visual_bill_renderer.go
```go
converter := &raster.Converter{
    MaxWidth:  576,
    Threshold: 0.6,  // Changed from 0.5 to 0.6
}
```

## How Threshold Works

**Threshold Range**: 0.0 to 1.0

- **0.0**: All pixels print black (darkest)
- **0.5**: 50% threshold (balanced)
- **0.6**: 60% threshold (lighter) ← Current setting
- **1.0**: No pixels print black (lightest)

**Logic:**
```go
if lightness(pixel) >= Threshold {
    // Print black dot
}
```

So với threshold 0.6:
- Pixels với lightness < 0.6 → White (không in)
- Pixels với lightness ≥ 0.6 → Black (in đen)

## Effect
- **Fewer black pixels**: Chỉ những pixel rất sáng mới được in đen
- **Lighter print**: Text và graphics nhạt hơn
- **Better match**: Gần với preview PNG hơn

## Testing
1. Create test order
2. Print to Zywell ZY303 (192.168.1.115:9100)
3. Compare with preview PNG
4. If still too dark → increase to 0.65 or 0.7
5. If too light → decrease to 0.55 or 0.5

## Threshold Tuning Guide

| Threshold | Effect | Use When |
|-----------|--------|----------|
| 0.3-0.4 | Very dark | Print is too light/faded |
| 0.5 | Balanced | Default setting |
| 0.6 | Lighter | Print is too dark (current) |
| 0.7-0.8 | Very light | Print is extremely dark |

## Files Modified
- `backend/application/services/chromedp_bill_renderer_optimized.go`
- `backend/application/services/visual_bill_renderer.go`

## Backend Status
✅ Backend restarted with threshold 0.6
- Server running on port 3000
- Ready to test print

## Next Steps
1. Test print với order thật
2. So sánh với preview PNG
3. Điều chỉnh threshold nếu cần:
   - Nếu vẫn đậm: tăng lên 0.65 hoặc 0.7
   - Nếu quá nhạt: giảm xuống 0.55 hoặc 0.5
