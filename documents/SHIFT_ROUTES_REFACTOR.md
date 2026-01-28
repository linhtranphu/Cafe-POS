# Shift Routes Refactoring

## Problem

Barista users were getting "insufficient permissions" error when trying to open a shift because shift routes were only available in the `/waiter` group.

## Root Cause

**Before**:
```go
// Shift routes only in waiter group
waiter := protected.Group("/waiter")
waiter.Use(http.RequireRole(user.RoleWaiter, user.RoleCashier, user.RoleManager))
{
    waiter.POST("/shifts/start", shiftHandler.StartShift)
    waiter.POST("/shifts/:id/end", shiftHandler.EndShift)
    waiter.GET("/shifts/current", shiftHandler.GetCurrentShift)
    waiter.GET("/shifts", shiftHandler.GetMyShifts)
}

// Barista group had NO shift routes
barista := protected.Group("/barista")
barista.Use(http.RequireRole(user.RoleBarista, user.RoleManager))
{
    // Only order routes, no shift routes
}
```

**Issue**: Barista couldn't access `/api/shifts/start` because:
1. Frontend calls `/api/shifts/start` (not `/api/waiter/shifts/start`)
2. Route only existed in waiter group with waiter role requirement
3. Barista role was rejected by `RequireRole(RoleWaiter, ...)`

## Solution

Created a common `/shifts` group accessible to all authenticated users:

```go
// Shift management - available for all roles (waiter, barista, cashier)
shifts := protected.Group("/shifts")
{
    shifts.POST("/start", shiftHandler.StartShift)
    shifts.POST("/:id/end", shiftHandler.EndShift)
    shifts.POST("/:id/close", shiftHandler.CloseShift)
    shifts.GET("/current", shiftHandler.GetCurrentShift)
    shifts.GET("/my", shiftHandler.GetMyShifts)
    shifts.GET("/:id", shiftHandler.GetShift)
}
```

## Route Structure After Refactoring

### Common Routes (All Authenticated Users)

```
POST   /api/shifts/start          - Start a shift (any role)
POST   /api/shifts/:id/end        - End a shift (any role)
POST   /api/shifts/:id/close      - Close shift and lock orders (cashier)
GET    /api/shifts/current        - Get current open shift (any role)
GET    /api/shifts/my             - Get my shift history (any role)
GET    /api/shifts/:id            - Get shift details (any role)
```

### Waiter Routes

```
POST   /api/waiter/orders         - Create order
POST   /api/waiter/orders/:id/payment - Collect payment
PUT    /api/waiter/orders/:id/edit - Edit order
POST   /api/waiter/orders/:id/send - Send to bar
POST   /api/waiter/orders/:id/serve - Serve order
GET    /api/waiter/orders         - Get my orders
GET    /api/waiter/orders/:id     - Get order details
GET    /api/waiter/menu           - View menu
GET    /api/waiter/ingredients    - View ingredients
GET    /api/waiter/facilities     - View facilities
```

### Barista Routes

```
GET    /api/barista/orders/queue  - View queued orders
GET    /api/barista/orders/my     - View my orders (in progress + ready)
POST   /api/barista/orders/:id/accept - Accept order from queue
POST   /api/barista/orders/:id/ready - Mark order as ready
GET    /api/barista/orders/:id    - Get order details
```

### Cashier Routes

```
POST   /api/cashier/shifts/:id/close - Close shift (duplicate, can use common)
GET    /api/cashier/shifts        - Get all shifts
GET    /api/cashier/shifts/:id    - Get shift details
GET    /api/cashier/shifts/:id/status - Get shift status
GET    /api/cashier/shifts/:id/payments - Get payments by shift
POST   /api/cashier/discrepancies - Report discrepancy
GET    /api/cashier/discrepancies/pending - Get pending discrepancies
```

### Manager Routes

```
GET    /api/manager/shifts        - Get all shifts
GET    /api/manager/shifts/:id    - Get shift details
(Plus all menu, ingredient, facility, user management routes)
```

## Benefits

### 1. Role Independence
- Each role can manage their own shifts
- No cross-role dependencies
- Barista doesn't need waiter permissions

### 2. Consistent API
- Frontend uses same endpoint for all roles: `/api/shifts/start`
- No need for role-specific shift endpoints
- Simpler frontend code

### 3. Security
- Role validation happens in service layer (checks `role_type`)
- Each user can only access their own shifts
- Middleware still validates authentication

### 4. Scalability
- Easy to add new roles (e.g., kitchen staff)
- New roles automatically get shift management
- No route duplication needed

## Security Considerations

### Authentication
- All shift routes require valid JWT token
- Middleware: `AuthMiddleware(jwtService)`

### Authorization
- Service layer validates shift belongs to user
- `FindOpenShiftByUser(userID, roleType)` ensures role match
- Cannot access other users' shifts

### Role-Specific Logic
- `StartShift()` uses role from JWT to set `role_type`
- `GetCurrentShift()` filters by user's role
- `GetMyShifts()` returns only shifts for user's role

## Migration Notes

### Breaking Changes
None - frontend already uses `/api/shifts/*` endpoints

### Backward Compatibility
- Old waiter routes removed (were redundant)
- All existing frontend code works without changes
- Database schema unchanged

### Testing Required
- ✅ Waiter can open/close shift
- ✅ Barista can open/close shift
- ✅ Cashier can open/close shift
- ✅ Each role sees only their own shifts
- ✅ Cannot access other users' shifts

## Code Changes

### File: `backend/main.go`

**Removed**:
```go
// From waiter group
waiter.POST("/shifts/start", shiftHandler.StartShift)
waiter.POST("/shifts/:id/end", shiftHandler.EndShift)
waiter.GET("/shifts/current", shiftHandler.GetCurrentShift)
waiter.GET("/shifts", shiftHandler.GetMyShifts)
```

**Added**:
```go
// Common shifts group
shifts := protected.Group("/shifts")
{
    shifts.POST("/start", shiftHandler.StartShift)
    shifts.POST("/:id/end", shiftHandler.EndShift)
    shifts.POST("/:id/close", shiftHandler.CloseShift)
    shifts.GET("/current", shiftHandler.GetCurrentShift)
    shifts.GET("/my", shiftHandler.GetMyShifts)
    shifts.GET("/:id", shiftHandler.GetShift)
}
```

### No Changes Required In:
- `shift_handler.go` - Already uses role from JWT
- `shift_service.go` - Already validates role
- Frontend code - Already uses `/api/shifts/*`

## Testing

### Manual Test

```bash
# Test barista can open shift
curl -X POST http://localhost:8080/api/shifts/start \
  -H "Authorization: Bearer <barista_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "MORNING",
    "start_cash": 0
  }'

# Expected: 201 Created with shift data

# Test waiter can open shift
curl -X POST http://localhost:8080/api/shifts/start \
  -H "Authorization: Bearer <waiter_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "MORNING",
    "start_cash": 1000000
  }'

# Expected: 201 Created with shift data
```

### Automated Test

```go
func TestShiftRoutes_AllRoles(t *testing.T) {
    roles := []string{"waiter", "barista", "cashier"}
    
    for _, role := range roles {
        t.Run(role, func(t *testing.T) {
            // Login as role
            token := loginAs(role)
            
            // Start shift
            resp := startShift(token)
            assert.Equal(t, 201, resp.StatusCode)
            
            // Get current shift
            resp = getCurrentShift(token)
            assert.Equal(t, 200, resp.StatusCode)
        })
    }
}
```

## Deployment

### Steps
1. Build new backend: `go build -o cafe-pos-server`
2. Stop old backend: `pkill -f cafe-pos-server`
3. Start new backend: `./cafe-pos-server`
4. Test with each role
5. Monitor logs for errors

### Rollback Plan
If issues occur:
1. Stop new backend
2. Start old backend
3. Investigate logs
4. Fix and redeploy

### Monitoring
Watch for:
- 403 Forbidden errors (permission issues)
- 404 Not Found errors (route issues)
- Shift creation failures
- Role-specific problems

## Related Documents

- `documents/SHIFT_BY_ROLE.md` - Shift management by role
- `documents/BR13_BARISTA_SHIFT_VALIDATION.md` - Barista shift requirement
- `documents/BARISTA_SHIFT_UI_NOTIFICATION.md` - UI notifications
- `RESTART_BACKEND.md` - How to restart backend

## Summary

Refactored shift routes from role-specific groups to a common group accessible by all authenticated users. This fixes the "insufficient permissions" error for barista users while maintaining security through service-layer validation. No frontend changes required.
