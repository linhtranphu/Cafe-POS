# ProfitAnalysisView Component

## Overview
The ProfitAnalysisView component provides comprehensive profit analysis functionality for cafe managers, including category-level profit analysis and operating profit analysis with date range filtering.

## Location
`frontend/src/views/ProfitAnalysisView.vue`

## Features

### 1. View Mode Toggle
- **Category View**: Shows profit analysis by menu category
- **Operating View**: Shows operating profit after deducting expenses

### 2. Date Range Picker
- **Presets**:
  - Today: Current date
  - This Week: Monday to Sunday of current week
  - This Month: First to last day of current month
- **Custom Range**: Manual date selection with start and end dates

### 3. Category Profit View
Displays profit metrics for each menu category:
- Total revenue
- Total cost (COGS)
- Total profit
- Average profit margin
- Order count
- Item count

### 4. Operating Profit View
Displays comprehensive operating profit analysis:
- **Gross Profit Section**:
  - Total revenue
  - Cost of goods sold (COGS)
  - Gross profit
  - Gross profit margin

- **Operating Expenses**:
  - Staff salary
  - Rent
  - Utilities
  - Marketing costs
  - Other expenses
  - Total expenses

- **Operating Profit**:
  - Operating profit (gross profit - expenses)
  - Operating profit margin

## Child Components

### CategoryProfitView
**Location**: `frontend/src/components/CategoryProfitView.vue`

**Props**:
- `dateRange`: Object with start and end dates
- `categoryProfits`: Array of category profit data

**Features**:
- Date range display
- Category profit cards
- Color-coded profit indicators
- Empty state handling

### OperatingProfitView
**Location**: `frontend/src/components/OperatingProfitView.vue`

**Props**:
- `operatingProfit`: Object with operating profit data

**Features**:
- Gross profit section
- Expenses breakdown
- Operating profit calculation
- Expense allocation indicator
- Warning for missing expenses

## Usage

### Navigation
```javascript
// Route
/manager/profit-analysis

// Router Link
<router-link to="/manager/profit-analysis">Profit Analysis</router-link>
```

### API Integration
```javascript
import { profitAnalysisService } from '@/services/profitAnalysis'

// Get category profit
const categoryData = await profitAnalysisService.getCategoryProfit({
  start: '2024-01-01',
  end: '2024-01-31'
})

// Get operating profit
const operatingData = await profitAnalysisService.getOperatingProfit({
  start: '2024-01-01',
  end: '2024-01-31'
})
```

## State Management

### Component State
```javascript
{
  loading: boolean,           // Loading state
  error: string | null,        // Error message
  viewMode: 'category' | 'operating',  // Current view mode
  selectedPreset: string | null,       // Selected date preset
  dateRange: {
    start: string,             // Start date (YYYY-MM-DD)
    end: string                // End date (YYYY-MM-DD)
  },
  categoryProfits: Array,      // Category profit data
  operatingProfit: Object      // Operating profit data
}
```

## Date Range Calculation

### Today
```javascript
{
  start: '2024-02-08',
  end: '2024-02-08'
}
```

### This Week
```javascript
{
  start: '2024-02-05',  // Monday
  end: '2024-02-11'     // Sunday
}
```

### This Month
```javascript
{
  start: '2024-02-01',  // First day
  end: '2024-02-29'     // Last day
}
```

## Color Coding

### Profit Margin
- **Green**: Margin >= 20% (good profit)
- **Yellow**: Margin < 20% (low margin)
- **Red**: Margin < 0% (loss)

### Profit Amount
- **Green**: Positive profit
- **Gray**: Break-even (0)
- **Red**: Negative profit (loss)

## Error Handling

### Loading State
```
⏳
Đang tải dữ liệu...
```

### Error State
```
❌
[Error message]
[Thử lại button]
```

### Empty State
```
📭
Không có dữ liệu trong khoảng thời gian này
```

## Responsive Design

### Mobile
- Full-screen layout
- Touch-friendly controls
- Pull-to-refresh support
- Scrollable content area
- Fixed header with controls

### Desktop
- Same layout optimized for larger screens
- Better spacing and typography

## Requirements Mapping

| Requirement | Implementation |
|-------------|----------------|
| 6.1 | Category profit table with all metrics |
| 6.4 | Date range filtering with presets |
| 6.5.1 | Operating profit analysis |
| 6.5.3 | Expense breakdown display |
| 6.5.4 | Operating profit calculation |
| 6.5.9 | Allocation note display |
| 7.1 | Manager view display |

## Testing

### Unit Tests
Location: `frontend/src/views/__tests__/ProfitAnalysisView.test.js`

Run tests:
```bash
npm test
```

### Manual Testing Checklist
- [ ] Navigate to /manager/profit-analysis
- [ ] Toggle between category and operating views
- [ ] Select date presets (today, this week, this month)
- [ ] Select custom date range
- [ ] Verify category profit data displays correctly
- [ ] Verify operating profit data displays correctly
- [ ] Test with missing expense data
- [ ] Test with allocated expenses
- [ ] Test error states
- [ ] Test pull-to-refresh
- [ ] Test on mobile devices

## Dependencies

- Vue 3
- Vue Router
- Pinia (auth store)
- profitAnalysisService (API service)
- PullToRefresh composable
- formatPrice utility

## Related Files

- `frontend/src/services/profitAnalysis.js` - API service
- `frontend/src/components/CategoryProfitView.vue` - Category view
- `frontend/src/components/OperatingProfitView.vue` - Operating view
- `frontend/src/utils/formatters.js` - Formatting utilities
- `frontend/src/router/index.js` - Router configuration

## Notes

### Expense Allocation
When viewing daily reports but expenses are entered monthly, the system allocates expenses proportionally. An indicator shows "⚠️ Phân bổ" and a note explains the allocation.

### Missing Expenses
When no operating expenses are entered, a warning is displayed:
```
⚠️ Chưa nhập chi phí vận hành
Hiện tại chỉ hiển thị lợi nhuận gộp. Vui lòng nhập chi phí vận hành 
để xem lợi nhuận vận hành chính xác.
```

### Date Format
All dates use ISO 8601 format (YYYY-MM-DD) for API communication and internal storage.

## Future Enhancements

- Export to CSV functionality
- Profit trend charts
- Comparison with previous periods
- Drill-down to individual items
- Print/PDF export
- Email reports
- Scheduled reports

---

**Last Updated**: February 8, 2026
