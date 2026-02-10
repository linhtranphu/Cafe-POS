# Menu Cost & Profit Analysis API Services

Quick reference guide for using the Menu Cost & Profit Analysis API services in Vue components.

## Import Services

```javascript
import { menuCostService } from '@/services/menuCost'
import { profitAnalysisService } from '@/services/profitAnalysis'
import { operatingExpenseService } from '@/services/operatingExpense'
```

## Menu Cost Service

### Get Menu Costs

Retrieve all menu items with cost and profit information.

```javascript
// Get all menu costs
const response = await menuCostService.getMenuCosts()

// Filter by category
const coffeeItems = await menuCostService.getMenuCosts({
  category: 'Coffee'
})

// Sort by profit margin (descending)
const sortedByMargin = await menuCostService.getMenuCosts({
  sort_by: 'profit_margin',
  sort_order: 'desc'
})

// Sort by absolute profit (ascending)
const sortedByProfit = await menuCostService.getMenuCosts({
  sort_by: 'absolute_profit',
  sort_order: 'asc'
})

// Response structure:
// {
//   items: [
//     {
//       menu_item_id: "...",
//       name: "Cappuccino",
//       category: "Coffee",
//       price: 45000,
//       current_cost: 15000,
//       profit_margin: 66.67,
//       absolute_profit: 30000,
//       cost_status: "FINAL",
//       cost_last_calculated_at: "2024-01-15T10:30:00Z",
//       warning_status: "none"
//     }
//   ],
//   summary: {
//     total_items: 50,
//     loss_count: 2,
//     low_margin_count: 5,
//     average_profit_margin: 55.5
//   },
//   recalculation_status: {
//     in_progress: false,
//     queued_items: 0,
//     processed_items: 50,
//     last_updated: "2024-01-15T10:30:00Z"
//   }
// }
```

### Get Menu Cost Detail

Retrieve detailed cost breakdown for a specific menu item.

```javascript
const menuItemId = "65a1b2c3d4e5f6g7h8i9j0k1"
const breakdown = await menuCostService.getMenuCostDetail(menuItemId)

// Response structure:
// {
//   menu_item: {
//     id: "...",
//     name: "Cappuccino",
//     price: 45000,
//     current_cost: 15000,
//     cost_status: "FINAL"
//   },
//   ingredients: [
//     {
//       name: "Espresso",
//       quantity: 30,
//       unit: "ml",
//       cost_per_unit: 200,
//       conversion_rate: 1.0,
//       wastage_percentage: 5.0,
//       total_cost: 6300
//     }
//   ],
//   total_cost: 15000
// }
```

### Get Menu Warnings

Retrieve menu items with warnings (loss or low margin).

```javascript
// Use default threshold from shop settings
const warnings = await menuCostService.getMenuWarnings()

// Use custom threshold (e.g., 25%)
const customWarnings = await menuCostService.getMenuWarnings(25)

// Response structure:
// {
//   loss_items: [
//     {
//       menu_item_id: "...",
//       name: "Special Promo Coffee",
//       price: 20000,
//       current_cost: 25000,
//       profit_margin: -25.0,
//       absolute_profit: -5000,
//       warning_status: "loss"
//     }
//   ],
//   low_margin_items: [
//     {
//       menu_item_id: "...",
//       name: "Basic Coffee",
//       price: 30000,
//       current_cost: 25000,
//       profit_margin: 16.67,
//       absolute_profit: 5000,
//       warning_status: "low_margin"
//     }
//   ],
//   loss_count: 1,
//   low_margin_count: 1,
//   threshold: 20.0
// }
```

## Profit Analysis Service

### Get Category Profit

Retrieve profit analysis aggregated by category.

```javascript
const dateRange = {
  start: '2024-01-01',
  end: '2024-01-31'
}

const categoryProfit = await profitAnalysisService.getCategoryProfit(dateRange)

// Response structure:
// {
//   date_range: {
//     start: "2024-01-01",
//     end: "2024-01-31"
//   },
//   categories: [
//     {
//       category: "Coffee",
//       total_revenue: 5000000,
//       total_cost: 1500000,
//       total_profit: 3500000,
//       average_profit_margin: 70.0,
//       order_count: 150,
//       item_count: 200
//     }
//   ]
// }
```

### Get Operating Profit

Retrieve operating profit analysis including expenses.

```javascript
const dateRange = {
  start: '2024-01-01',
  end: '2024-01-31'
}

const operatingProfit = await profitAnalysisService.getOperatingProfit(dateRange)

// Response structure:
// {
//   date_range: {
//     start: "2024-01-01",
//     end: "2024-01-31"
//   },
//   total_revenue: 10000000,
//   total_cogs: 3000000,
//   gross_profit: 7000000,
//   gross_profit_margin: 70.0,
//   staff_salary: 2000000,
//   rent: 1000000,
//   utilities: 500000,
//   marketing_costs: 300000,
//   other_expenses: 200000,
//   total_expenses: 4000000,
//   operating_profit: 3000000,
//   operating_profit_margin: 30.0,
//   expense_allocated: false
// }
```

## Operating Expense Service

### Create Operating Expense

Create or update an operating expense record.

```javascript
const expenseData = {
  period_start: '2024-01-01',
  period_end: '2024-01-31',
  staff_salary: 2000000,
  rent: 1000000,
  utilities: 500000,
  marketing_costs: 300000,
  other_expenses: 200000
}

const createdExpense = await operatingExpenseService.createOperatingExpense(expenseData)

// Response structure:
// {
//   id: "...",
//   period_start: "2024-01-01",
//   period_end: "2024-01-31",
//   staff_salary: 2000000,
//   rent: 1000000,
//   utilities: 500000,
//   marketing_costs: 300000,
//   other_expenses: 200000,
//   total_expenses: 4000000  // Auto-calculated
// }
```

### Get Operating Expenses

Retrieve operating expenses with optional date range filter.

```javascript
// Get all expenses
const allExpenses = await operatingExpenseService.getOperatingExpenses()

// Get expenses for a specific date range
const dateRange = {
  start: '2024-01-01',
  end: '2024-03-31'
}
const filteredExpenses = await operatingExpenseService.getOperatingExpenses(dateRange)

// Response structure:
// {
//   expenses: [
//     {
//       id: "...",
//       period_start: "2024-01-01",
//       period_end: "2024-01-31",
//       staff_salary: 2000000,
//       rent: 1000000,
//       utilities: 500000,
//       marketing_costs: 300000,
//       other_expenses: 200000,
//       total_expenses: 4000000
//     }
//   ]
// }
```

## Error Handling

All services use the axios instance with interceptors. Errors should be handled in components:

```javascript
try {
  const response = await menuCostService.getMenuCosts()
  // Handle success
} catch (error) {
  if (error.response) {
    // Server responded with error status
    console.error('Server error:', error.response.status, error.response.data)
  } else if (error.request) {
    // Request made but no response
    console.error('Network error:', error.message)
  } else {
    // Other errors
    console.error('Error:', error.message)
  }
}
```

## Usage in Vue Components

### Composition API Example

```javascript
import { ref, onMounted } from 'vue'
import { menuCostService } from '@/services/menuCost'

export default {
  setup() {
    const menuCosts = ref([])
    const loading = ref(false)
    const error = ref(null)

    const fetchMenuCosts = async () => {
      loading.value = true
      error.value = null
      
      try {
        const response = await menuCostService.getMenuCosts({
          sort_by: 'profit_margin',
          sort_order: 'desc'
        })
        menuCosts.value = response.items
      } catch (err) {
        error.value = err.message
      } finally {
        loading.value = false
      }
    }

    onMounted(() => {
      fetchMenuCosts()
    })

    return {
      menuCosts,
      loading,
      error,
      fetchMenuCosts
    }
  }
}
```

### Options API Example

```javascript
import { menuCostService } from '@/services/menuCost'

export default {
  data() {
    return {
      menuCosts: [],
      loading: false,
      error: null
    }
  },
  
  mounted() {
    this.fetchMenuCosts()
  },
  
  methods: {
    async fetchMenuCosts() {
      this.loading = true
      this.error = null
      
      try {
        const response = await menuCostService.getMenuCosts({
          sort_by: 'profit_margin',
          sort_order: 'desc'
        })
        this.menuCosts = response.items
      } catch (err) {
        this.error = err.message
      } finally {
        this.loading = false
      }
    }
  }
}
```

## Type Safety with JSDoc

All services include JSDoc type definitions. Enable type checking in your IDE:

```javascript
/**
 * @type {import('@/services/types/menuCost').MenuItemCost[]}
 */
const menuCosts = []

/**
 * @type {import('@/services/types/menuCost').DateRange}
 */
const dateRange = {
  start: '2024-01-01',
  end: '2024-01-31'
}
```

## Notes

- All date parameters use ISO 8601 format (YYYY-MM-DD)
- All monetary values are in VND (Vietnamese Dong)
- Percentages are returned as numbers (e.g., 66.67 for 66.67%)
- The API automatically handles authentication via axios interceptors
- 401 errors trigger automatic logout and redirect to login page
