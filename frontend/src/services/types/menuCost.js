/**
 * @fileoverview Type definitions for Menu Cost & Profit Analysis
 * These JSDoc types provide type safety for the cost and profit analysis features
 */

/**
 * Cost status enum
 * @typedef {'FINAL' | 'ESTIMATED' | 'INCOMPLETE'} CostStatus
 */

/**
 * Warning status enum
 * @typedef {'none' | 'low_margin' | 'loss'} WarningStatus
 */

/**
 * Menu item with cost information
 * @typedef {Object} MenuItemCost
 * @property {string} menu_item_id - Menu item ID
 * @property {string} name - Menu item name
 * @property {string} category - Menu item category
 * @property {number} price - Menu item price
 * @property {number} current_cost - Current cost of goods sold
 * @property {number} profit_margin - Profit margin percentage
 * @property {number} absolute_profit - Absolute profit amount
 * @property {CostStatus} cost_status - Status of cost calculation
 * @property {string} cost_last_calculated_at - ISO timestamp of last calculation
 * @property {WarningStatus} warning_status - Warning status (none, low_margin, loss)
 */

/**
 * Ingredient cost detail for breakdown
 * @typedef {Object} IngredientCostDetail
 * @property {string} name - Ingredient name
 * @property {number} quantity - Quantity used
 * @property {string} unit - Unit of measurement
 * @property {number} cost_per_unit - Cost per unit
 * @property {number} conversion_rate - Unit conversion rate
 * @property {number} wastage_percentage - Wastage percentage
 * @property {number} total_cost - Total cost for this ingredient
 */

/**
 * Menu item cost breakdown response
 * @typedef {Object} MenuItemCostBreakdown
 * @property {Object} menu_item - Menu item basic info
 * @property {string} menu_item.id - Menu item ID
 * @property {string} menu_item.name - Menu item name
 * @property {number} menu_item.price - Menu item price
 * @property {number} menu_item.current_cost - Current cost
 * @property {CostStatus} menu_item.cost_status - Cost status
 * @property {IngredientCostDetail[]} ingredients - Array of ingredient cost details
 * @property {number} total_cost - Total cost
 */

/**
 * Summary statistics for menu costs
 * @typedef {Object} MenuCostSummary
 * @property {number} total_items - Total number of menu items
 * @property {number} loss_count - Number of items sold at loss
 * @property {number} low_margin_count - Number of items with low margin
 * @property {number} average_profit_margin - Average profit margin percentage
 */

/**
 * Recalculation status
 * @typedef {Object} RecalculationStatus
 * @property {boolean} in_progress - Whether recalculation is in progress
 * @property {number} queued_items - Number of items queued for recalculation
 * @property {number} processed_items - Number of items processed
 * @property {string} last_updated - ISO timestamp of last update
 */

/**
 * Menu costs response
 * @typedef {Object} MenuCostsResponse
 * @property {MenuItemCost[]} items - Array of menu items with cost info
 * @property {MenuCostSummary} summary - Summary statistics
 * @property {RecalculationStatus} recalculation_status - Recalculation status
 */

/**
 * Profit warnings response
 * @typedef {Object} ProfitWarnings
 * @property {MenuItemCost[]} loss_items - Items sold at loss
 * @property {MenuItemCost[]} low_margin_items - Items with low margin
 * @property {number} loss_count - Count of loss items
 * @property {number} low_margin_count - Count of low margin items
 * @property {number} threshold - Low margin threshold used
 */

/**
 * Category profit analysis
 * @typedef {Object} CategoryProfit
 * @property {string} category - Category name
 * @property {number} total_revenue - Total revenue for category
 * @property {number} total_cost - Total cost for category
 * @property {number} total_profit - Total profit for category
 * @property {number} average_profit_margin - Average profit margin percentage
 * @property {number} order_count - Number of orders
 * @property {number} item_count - Number of items sold
 */

/**
 * Date range filter
 * @typedef {Object} DateRange
 * @property {string} start - Start date (ISO 8601 format)
 * @property {string} end - End date (ISO 8601 format)
 */

/**
 * Category profit response
 * @typedef {Object} CategoryProfitResponse
 * @property {DateRange} date_range - Date range for the report
 * @property {CategoryProfit[]} categories - Array of category profits
 */

/**
 * Operating profit report
 * @typedef {Object} OperatingProfitReport
 * @property {DateRange} date_range - Date range for the report
 * @property {number} total_revenue - Total revenue
 * @property {number} total_cogs - Total cost of goods sold
 * @property {number} gross_profit - Gross profit
 * @property {number} gross_profit_margin - Gross profit margin percentage
 * @property {number} staff_salary - Staff salary expenses
 * @property {number} rent - Rent expenses
 * @property {number} utilities - Utilities expenses
 * @property {number} marketing_costs - Marketing expenses
 * @property {number} other_expenses - Other expenses
 * @property {number} total_expenses - Total operating expenses
 * @property {number} operating_profit - Operating profit
 * @property {number} operating_profit_margin - Operating profit margin percentage
 * @property {boolean} expense_allocated - Whether expenses were allocated
 * @property {string} [allocation_note] - Note about expense allocation
 */

/**
 * Operating expense data
 * @typedef {Object} OperatingExpense
 * @property {string} [id] - Expense ID (optional for create)
 * @property {string} period_start - Period start date (ISO 8601)
 * @property {string} period_end - Period end date (ISO 8601)
 * @property {number} staff_salary - Staff salary amount
 * @property {number} rent - Rent amount
 * @property {number} utilities - Utilities amount
 * @property {number} marketing_costs - Marketing costs amount
 * @property {number} other_expenses - Other expenses amount
 * @property {number} total_expenses - Total expenses (auto-calculated)
 */

/**
 * Operating expenses list response
 * @typedef {Object} OperatingExpensesResponse
 * @property {OperatingExpense[]} expenses - Array of operating expenses
 */

/**
 * Profit filter options
 * @typedef {Object} ProfitFilter
 * @property {string} [category] - Filter by category
 * @property {'profit_margin' | 'absolute_profit' | 'name'} [sort_by] - Sort field
 * @property {'asc' | 'desc'} [sort_order] - Sort order
 */

export {}
