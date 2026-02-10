# Task 17: Responsive Design Visual Guide

## Overview
This guide demonstrates the responsive design improvements implemented for MenuCostView and ProfitAnalysisView.

## MenuCostView Responsive Layouts

### Mobile View (< 768px)
```
┌─────────────────────────────────┐
│ 💰 Chi phí món                  │
│                                 │
│ [Search: Tìm kiếm món...]       │
│                                 │
│ [📁 Tất cả] [Coffee] [Tea] →    │
│                                 │
│ Sắp xếp: [Lợi nhuận %] [↓ Giảm]│
├─────────────────────────────────┤
│ ┌─────────────────────────────┐ │
│ │ Tổng quan chi phí món       │ │
│ │ [50] [2] [5] [66.67%]       │ │
│ └─────────────────────────────┘ │
│                                 │
│ ┌─────────────────────────────┐ │
│ │ Cappuccino        [Coffee]  │ │
│ │ Giá: 45.000₫               │ │
│ │ Chi phí: 15.000₫           │ │
│ │ LN %: 66.67% | LN: 30.000₫ │ │
│ └─────────────────────────────┘ │
│                                 │
│ ┌─────────────────────────────┐ │
│ │ Latte             [Coffee]  │ │
│ │ ...                         │ │
│ └─────────────────────────────┘ │
└─────────────────────────────────┘
```

### Desktop View (≥ 768px)
```
┌───────────────────────────────────────────────────────────────┐
│ 💰 Chi phí món                    [📱 Thẻ] [📊 Bảng]         │
│                                                               │
│ [Search: Tìm kiếm món...]                                     │
│                                                               │
│ [📁 Tất cả] [Coffee] [Tea] [Smoothie] [Snacks]               │
│                                                               │
│ Sắp xếp: [Lợi nhuận %] [↓ Giảm]                              │
├───────────────────────────────────────────────────────────────┤
│ ┌───────────────────────────────────────────────────────────┐ │
│ │ Tổng quan chi phí món                                     │ │
│ │ [50 Tổng món] [2 Bán lỗ] [5 LN thấp] [66.67% LN TB]      │ │
│ └───────────────────────────────────────────────────────────┘ │
│                                                               │
│ ┌───────────────────────────────────────────────────────────┐ │
│ │ Tên món    │ Danh mục │ Giá bán │ Chi phí │ LN % │ LN tiền│ │
│ ├───────────────────────────────────────────────────────────┤ │
│ │ Cappuccino │ Coffee   │ 45.000₫ │ 15.000₫ │ 66.67% │ 30k  │ │
│ │ Latte      │ Coffee   │ 50.000₫ │ 18.000₫ │ 64.00% │ 32k  │ │
│ │ Espresso   │ Coffee   │ 35.000₫ │ 12.000₫ │ 65.71% │ 23k  │ │
│ └───────────────────────────────────────────────────────────┘ │
└───────────────────────────────────────────────────────────────┘
```

## ProfitAnalysisView Responsive Layouts

### Mobile View (< 768px)
```
┌─────────────────────────────────┐
│ 📊 Phân tích lợi nhuận          │
│                                 │
│ [📁 Theo danh mục] [💼 Vận hành]│
│                                 │
│ [Hôm nay] [Tuần này] [Tháng]   │
│                                 │
│ [2024-01-01]                    │
│ [2024-01-31]                    │
├─────────────────────────────────┤
│ ┌─────────────────────────────┐ │
│ │ Khoảng thời gian            │ │
│ │ 01/01/2024 → 31/01/2024     │ │
│ └─────────────────────────────┘ │
│                                 │
│ ┌─────────────────────────────┐ │
│ │ Coffee                      │ │
│ │ 150 đơn • 200 món           │ │
│ │ LN %: 70.00%                │ │
│ │                             │ │
│ │ Doanh thu: 5.000.000₫       │ │
│ │ Chi phí: 1.500.000₫         │ │
│ │ Lợi nhuận: 3.500.000₫       │ │
│ └─────────────────────────────┘ │
└─────────────────────────────────┘
```

### Desktop View (≥ 768px)
```
┌───────────────────────────────────────────────────────────────┐
│ 📊 Phân tích lợi nhuận                                        │
│                                                               │
│ [📁 Theo danh mục] [💼 Lợi nhuận vận hành]                    │
│                                                               │
│ [Hôm nay] [Tuần này] [Tháng này]                              │
│                                                               │
│ [2024-01-01] → [2024-01-31]                                   │
├───────────────────────────────────────────────────────────────┤
│ ┌───────────────────────────────────────────────────────────┐ │
│ │ Khoảng thời gian: 01/01/2024 → 31/01/2024                │ │
│ └───────────────────────────────────────────────────────────┘ │
│                                                               │
│ ┌───────────────────────────────────────────────────────────┐ │
│ │ Danh mục │ Đơn hàng │ Doanh thu │ Chi phí │ LN │ LN %     │ │
│ ├───────────────────────────────────────────────────────────┤ │
│ │ Coffee   │ 150      │ 5.000.000₫│ 1.5M₫   │ 3.5M│ 70.00% │ │
│ │ Tea      │ 80       │ 2.000.000₫│ 800k₫   │ 1.2M│ 60.00% │ │
│ │ Smoothie │ 60       │ 3.000.000₫│ 1.2M₫   │ 1.8M│ 60.00% │ │
│ └───────────────────────────────────────────────────────────┘ │
└───────────────────────────────────────────────────────────────┘
```

## Loading States

### Skeleton Loader - Card
```
┌─────────────────────────────────┐
│ ████████████        ████        │
│ ████████                        │
│                                 │
│ ████ ████████                   │
│                                 │
│ ████████  ████████              │
└─────────────────────────────────┘
```

### Skeleton Loader - Table
```
┌───────────────────────────────────────────────────────────────┐
│ Header │ Header │ Header │ Header │ Header │ Header          │
├───────────────────────────────────────────────────────────────┤
│ ████   │ ████   │ ████   │ ████   │ ████   │ ████            │
│ ████   │ ████   │ ████   │ ████   │ ████   │ ████            │
│ ████   │ ████   │ ████   │ ████   │ ████   │ ████            │
└───────────────────────────────────────────────────────────────┘
```

## Empty States

### No Data
```
┌─────────────────────────────────┐
│                                 │
│           📭                    │
│                                 │
│    Không tìm thấy món nào       │
│                                 │
│  Thử thay đổi bộ lọc hoặc      │
│       tìm kiếm                  │
│                                 │
└─────────────────────────────────┘
```

### Error State
```
┌─────────────────────────────────┐
│                                 │
│           ❌                    │
│                                 │
│      Lỗi tải dữ liệu            │
│                                 │
│  Không thể kết nối đến server   │
│                                 │
│      [🔄 Thử lại]               │
│                                 │
└─────────────────────────────────┘
```

## Responsive Breakpoints

### Tailwind CSS Breakpoints Used
- **sm**: 640px (small tablets)
- **md**: 768px (tablets)
- **lg**: 1024px (laptops)
- **xl**: 1280px (desktops)

### Component Behavior by Screen Size

| Component | Mobile (<768px) | Tablet (768-1024px) | Desktop (>1024px) |
|-----------|----------------|---------------------|-------------------|
| MenuCostView | Card layout | Card + Table toggle | Card + Table toggle |
| CategoryProfitView | Card layout | Card layout | Table layout |
| OperatingProfitView | Stacked sections | Stacked sections | Stacked sections |
| Date Picker | Vertical stack | Horizontal | Horizontal |
| Filters | Horizontal scroll | Wrap | Wrap |
| Summary Stats | 4 columns | 4 columns | 4 columns (larger) |

## Number Formatting Examples

### Vietnamese Locale Formatting
```javascript
// Prices
formatPrice(45000)        // "45.000 ₫"
formatPrice(1500000)      // "1.500.000 ₫"

// Percentages
formatPercentage(66.67)   // "66,67%"
formatPercentage(100)     // "100,00%"

// Numbers
formatNumber(1000)        // "1.000"
formatNumber(1500000)     // "1.500.000"
```

## Accessibility Features

### Keyboard Navigation
- Tab through interactive elements
- Enter to activate buttons
- Arrow keys for dropdowns

### Screen Reader Support
- Semantic HTML elements
- ARIA labels where needed
- Descriptive button text

### Color Contrast
- Text meets WCAG AA standards
- Icons have sufficient contrast
- Focus indicators visible

## Performance Optimizations

### CSS
- Tailwind CSS utilities (minimal CSS)
- No custom CSS animations (use Tailwind)
- Efficient class names

### JavaScript
- Lazy loading for heavy components
- Debounced search input
- Memoized computed properties

### Images
- No images used (CSS-based design)
- Icon fonts or SVG icons
- Optimized bundle size

## Browser Support
- Chrome 90+
- Firefox 88+
- Safari 14+
- Edge 90+
- Mobile browsers (iOS Safari, Chrome Mobile)

## Testing Checklist

### Visual Testing
- [ ] Test on iPhone SE (375px)
- [ ] Test on iPhone 12 (390px)
- [ ] Test on iPad (768px)
- [ ] Test on iPad Pro (1024px)
- [ ] Test on MacBook (1440px)
- [ ] Test on Desktop (1920px)

### Functional Testing
- [ ] View toggle works
- [ ] Filters work on all sizes
- [ ] Sort works on all sizes
- [ ] Date picker works on all sizes
- [ ] Loading states display correctly
- [ ] Empty states display correctly
- [ ] Error states display correctly
- [ ] Number formatting correct

### Performance Testing
- [ ] Page loads in < 3 seconds
- [ ] Smooth scrolling
- [ ] No layout shifts
- [ ] Responsive to user input

## Conclusion
The responsive design implementation provides:
- Optimal viewing experience on all devices
- Consistent UX across screen sizes
- Better perceived performance with skeletons
- Clear feedback with empty/error states
- Proper Vietnamese locale formatting

All components are mobile-first and progressively enhanced for larger screens.
