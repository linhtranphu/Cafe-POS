# Task 2.6 Usage Example: QueueCostRecalculation

## Overview

This document provides examples of how to use the `QueueCostRecalculation` method in the Cost Calculator Service.

## Setup

### 1. Initialize the Service

```go
// In your service initialization (e.g., main.go)
menuRepo := mongodb.NewMenuRepository(db)
ingredientRepo := mongodb.NewIngredientRepository(db)
orderRepo := mongodb.NewOrderRepository(db)
orderItemRepo := mongodb.NewOrderItemRepository(db)

costCalculatorService := services.NewCostCalculatorService(
    menuRepo,
    ingredientRepo,
    orderRepo,
    orderItemRepo,
)
```

### 2. Start Background Workers

```go
// Start workers when application starts
// Recommended: 3-5 workers for typical load
ctx := context.Background()
numWorkers := 5

costCalculatorService.StartRecalculationWorker(ctx, numWorkers)
```

## Usage in Ingredient Update Handler

### Example: Ingredient Handler Integration

```go
// In ingredient_handler.go
func (h *IngredientHandler) UpdateIngredient(c *gin.Context) {
    var req ingredient.UpdateIngredientRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    ingredientID, err := primitive.ObjectIDFromHex(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ingredient ID"})
        return
    }

    // Get the current ingredient to check if cost changed
    currentIngredient, err := h.ingredientService.GetIngredient(c.Request.Context(), ingredientID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Ingredient not found"})
        return
    }

    // Update the ingredient
    updatedIngredient, err := h.ingredientService.UpdateIngredient(c.Request.Context(), ingredientID, &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Check if cost_per_unit changed
    if req.CostPerUnit != nil && *req.CostPerUnit != currentIngredient.CostPerUnit {
        // Cost changed, queue background recalculation for all affected menu items
        err := h.costCalculatorService.QueueCostRecalculation(c.Request.Context(), ingredientID)
        if err != nil {
            // Log error but don't fail the update
            // The ingredient update was successful, recalculation is secondary
            log.Printf("Warning: Failed to queue cost recalculation for ingredient %s: %v", 
                ingredientID.Hex(), err)
        } else {
            log.Printf("Queued cost recalculation for menu items using ingredient: %s", 
                updatedIngredient.Name)
        }
    }

    c.JSON(http.StatusOK, updatedIngredient)
}
```

## Monitoring Queue Status

### Check Queue Size

```go
// Get current queue size for monitoring
queueSize := costCalculatorService.GetRecalculationQueueSize()
log.Printf("Current recalculation queue size: %d", queueSize)
```

### Health Check Endpoint

```go
// In your health check handler
func (h *HealthHandler) GetStatus(c *gin.Context) {
    queueSize := h.costCalculatorService.GetRecalculationQueueSize()
    
    status := gin.H{
        "status": "healthy",
        "recalculation_queue_size": queueSize,
    }
    
    // Warn if queue is getting full
    if queueSize > 800 {
        status["warning"] = "Recalculation queue is nearly full"
    }
    
    c.JSON(http.StatusOK, status)
}
```

## Example Workflow

### Scenario: Update Espresso Cost

```go
// 1. Manager updates Espresso cost from 200 to 250
ingredientID := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")
updateReq := &ingredient.UpdateIngredientRequest{
    CostPerUnit: float64Ptr(250.0),
}

// 2. Ingredient service updates the ingredient
updatedIngredient, err := ingredientService.UpdateIngredient(ctx, ingredientID, updateReq)

// 3. Handler queues cost recalculation
err = costCalculatorService.QueueCostRecalculation(ctx, ingredientID)
// This finds all menu items using "Espresso" and queues them:
// - Cappuccino (ID: xxx)
// - Latte (ID: yyy)
// - Americano (ID: zzz)

// 4. Background workers process the queue
// Worker 1: Recalculates Cappuccino cost
// Worker 2: Recalculates Latte cost
// Worker 3: Recalculates Americano cost

// 5. Menu items are updated with new costs
// Cappuccino: current_cost updated from 14550 to 16050
// Latte: current_cost updated from 16250 to 18250
// Americano: current_cost updated from 12600 to 15600
```

## Error Handling

### Queue Full Scenario

```go
// If queue is full (1000 items), QueueCostRecalculation will:
// 1. Try to queue each menu item
// 2. Skip items that don't fit (non-blocking)
// 3. Log warning: "Warning: recalculation queue is full, skipping menu item xxx"
// 4. Return success (ingredient update succeeded)

// The skipped items will have stale costs until:
// - Next ingredient update triggers recalculation
// - Manual recalculation via CalculateAllMenuItemCosts
```

### Worker Error Handling

```go
// If a worker fails to recalculate a menu item:
// 1. Error is logged: "Worker 1: failed to recalculate cost for menu item xxx: error"
// 2. Worker continues processing other items
// 3. Failed item remains with old cost
// 4. Can be retried manually or on next ingredient update
```

## Performance Tuning

### Adjust Number of Workers

```go
// For light load (< 100 menu items)
numWorkers := 2

// For moderate load (100-500 menu items)
numWorkers := 5

// For heavy load (> 500 menu items)
numWorkers := 10

// Note: More workers = more database connections
// Ensure your connection pool can handle it
```

### Adjust Queue Size

```go
// In NewCostCalculatorService, modify queueSize:
queueSize := 2000 // Increase for larger menus

return &CostCalculatorService{
    // ...
    recalcQueue: make(chan primitive.ObjectID, queueSize),
    queueSize:   queueSize,
}
```

## Testing

### Manual Test

```bash
# 1. Start the server with workers
go run main.go

# 2. Update an ingredient cost
curl -X PATCH http://localhost:8080/api/ingredients/507f1f77bcf86cd799439011 \
  -H "Content-Type: application/json" \
  -d '{"cost_per_unit": 250.0}'

# 3. Check logs for queue activity
# Expected output:
# "Queued cost recalculation for menu items using ingredient: Espresso"
# "Worker 1: Recalculating cost for menu item xxx"
# "Worker 2: Recalculating cost for menu item yyy"

# 4. Verify menu item costs updated
curl http://localhost:8080/api/menu/costs
```

## Best Practices

1. **Always start workers on application startup**
   - Workers should run for the lifetime of the application
   - Use context for graceful shutdown

2. **Don't block on queue operations**
   - QueueCostRecalculation uses non-blocking send
   - Ingredient updates succeed even if queue is full

3. **Monitor queue size**
   - Add metrics/logging for queue size
   - Alert if queue is consistently full

4. **Handle worker errors gracefully**
   - Log errors but don't crash workers
   - Consider adding retry logic for transient errors

5. **Test with realistic data**
   - Test with 100+ menu items
   - Test with multiple concurrent ingredient updates
   - Verify workers don't overwhelm database

## Future Enhancements

1. **Add retry logic**: Retry failed recalculations with exponential backoff
2. **Add metrics**: Track queue size, processing time, error rate
3. **Add priority**: Prioritize popular menu items
4. **Add batching**: Batch multiple ingredient updates
5. **Add persistence**: Use Redis for durable queue

## Helper Function

```go
// Helper to convert float64 to pointer
func float64Ptr(f float64) *float64 {
    return &f
}
```
