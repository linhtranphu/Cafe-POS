# Task 2.4 Implementation Summary: CalculateShiftOrderCosts Method

## Overview
Implemented the `CalculateShiftOrderCosts` method in the Cost Calculator Service to calculate and store accounting costs for all orders in a shift when the shift is closed.

## Requirements Addressed
- **Requirement 5.2**: Calculate cost for ALL orders in a shift using cost_per_unit values at the time of shift closure
- **Requirement 5.3**: Use the same calculation method as current menu item cost (sum of ingredient.quantity * ingredient.cost_per_unit)
- **Requirement 5.4**: Store the calculated cost in the order_items table with a timestamp

## Implementation Details

### 1. Service Method Signature
```go
func (s *CostCalculatorService) CalculateShiftOrderCosts(ctx context.Context, shiftID primitive.ObjectID) (*ShiftCostCalculationResult, error)
```

### 2. Key Features

#### Batch Processing
- Fetches all orders in the shift in a single query
- Fetches all ingredients once for efficiency
- Fetches all menu items once for efficiency
- Builds lookup maps to avoid repeated database queries

#### Cost Calculation
- For each order item:
  1. Looks up the menu item by ID
  2. Calculates cost using the same formula as `CalculateMenuItemCost`:
     - `cost = sum(quantity * cost_per_unit * conversion_rate * (1 + wastage_percentage/100))`
  3. Multiplies by order item quantity to get total accounting cost
  4. Rounds to 2 decimal places

#### Cost Status Handling
- **FINAL**: All ingredients have valid cost_per_unit values
- **INCOMPLETE**: One or more ingredients have missing or zero cost_per_unit
- Converts between `menu.CostStatus` and `order.CostStatus` types

#### Data Persistence
- Creates `OrderItemWithCost` records for each order item
- Saves all records to the `order_items` collection using batch insert
- Records include:
  - Order ID and Menu Item ID references
  - Item details (name, price, quantity, subtotal)
  - Accounting cost and calculation timestamp
  - Cost status (FINAL or INCOMPLETE)

### 3. Result Summary
Returns `ShiftCostCalculationResult` with:
- `TotalOrders`: Number of orders processed
- `TotalItems`: Number of order items processed
- `ItemsWithFinalCost`: Count of items with complete cost data
- `ItemsWithIncompleteCost`: Count of items with missing ingredient costs
- `TotalAccountingCost`: Sum of all accounting costs (only FINAL items)

### 4. Error Handling
- Gracefully handles missing menu items (marks as INCOMPLETE)
- Gracefully handles missing ingredient costs (marks as INCOMPLETE)
- Returns descriptive errors for database failures
- Continues processing even if some items have incomplete data

## Code Changes

### Modified Files
1. **backend/application/services/cost_calculator_service.go**
   - Added `OrderItemRepository` interface
   - Updated `CostCalculatorService` struct to include order and order item repositories
   - Updated `NewCostCalculatorService` constructor
   - Added `ShiftCostCalculationResult` struct
   - Implemented `CalculateShiftOrderCosts` method

2. **backend/application/services/cost_calculator_service_test.go**
   - Added `mockOrderRepository` implementation
   - Added `mockOrderItemRepository` implementation
   - Added `TestCalculateShiftOrderCosts` - main test case
   - Added `TestCalculateShiftOrderCosts_EmptyShift` - edge case test
   - Added `TestCalculateShiftOrderCosts_IncompleteCost` - incomplete data test
   - Updated all existing test cases to pass new repository parameters

## Test Coverage

### Test Cases
1. **TestCalculateShiftOrderCosts**: 
   - Tests calculation with multiple orders and items
   - Verifies correct cost calculation with conversion rate and wastage
   - Validates result statistics
   - Confirms order items are saved to database

2. **TestCalculateShiftOrderCosts_EmptyShift**:
   - Tests behavior with no orders in shift
   - Verifies graceful handling of empty result

3. **TestCalculateShiftOrderCosts_IncompleteCost**:
   - Tests handling of missing ingredient costs
   - Verifies INCOMPLETE status is set correctly
   - Confirms items are still saved with zero cost

### Test Data
- **Cappuccino**: 
  - Espresso (30ml @ 200/ml, 5% wastage) = 6,300
  - Milk (150ml @ 50/ml, 10% wastage) = 8,250
  - Total per item: 14,550

- **Latte**:
  - Espresso (30ml @ 200/ml, 5% wastage) = 6,300
  - Milk (200ml @ 50/ml, 10% wastage) = 11,000
  - Total per item: 17,300

## Integration Points

### Dependencies
- `MenuRepository`: To fetch menu items and their recipes
- `IngredientRepository`: To fetch current ingredient costs
- `OrderRepository`: To fetch orders by shift ID
- `OrderItemRepository`: To save order items with accounting costs

### Usage
This method will be called by the Shift Service when closing a shift:
```go
result, err := costCalculatorService.CalculateShiftOrderCosts(ctx, shiftID)
if err != nil {
    // Handle error
}
// Use result for reporting
```

## Design Decisions

### 1. Separate Collection for Order Items
- Uses the `order_items` collection (separate from embedded order items)
- Enables efficient querying and aggregation for profit reports
- Follows the design decision documented in the design.md

### 2. Immutable Accounting Cost
- Once saved, accounting costs are not recalculated
- Reflects the cost at the time of shift closure
- Aligns with accounting best practices

### 3. Graceful Degradation
- Continues processing even if some items have incomplete data
- Marks incomplete items clearly with INCOMPLETE status
- Allows partial cost tracking rather than failing completely

### 4. Type Conversion
- Handles conversion between `menu.CostStatus` and `order.CostStatus`
- Both packages define their own CostStatus type
- Explicit conversion ensures type safety

## Next Steps

### Task 2.5: Property-Based Test
- Implement property test for shift closure cost calculation
- Validate that accounting costs are immutable after shift closure
- Test with various order and ingredient combinations

### Integration with Shift Service
- Modify the shift closure endpoint to call `CalculateShiftOrderCosts`
- Return cost calculation summary in the shift closure response
- Handle cost calculation errors gracefully

## Notes

- The implementation follows the same cost calculation formula as `CalculateMenuItemCost`
- All costs are rounded to 2 decimal places for consistency
- The method is designed for batch processing efficiency
- Pre-existing test failures in other services do not affect this implementation
