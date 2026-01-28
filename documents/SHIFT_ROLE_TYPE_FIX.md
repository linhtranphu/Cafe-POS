# Shift Role Type Conversion Fix

## Problem

When starting a shift, the application crashed with error:
```
interface conversion: interface {} is user.Role, not string
/usr/local/Cellar/go/1.25.5/libexec/src/runtime/iface.go:275
```

## Root Cause

In `shift_handler.go`, the code was trying to convert `role` from context as `string`:

```go
role, _ := c.Get("role")
switch role.(string) {  // ❌ WRONG: role is user.Role, not string
    case "waiter":
        roleType = order.RoleWaiter
    // ...
}
```

However, the middleware sets `role` as `user.Role` type:

```go
// middleware.go
c.Set("role", claims.Role)  // claims.Role is user.Role type
```

## Solution

### 1. Created Helper Function in Domain

Added `ParseRoleType()` function in `backend/domain/order/shift.go`:

```go
// ParseRoleType converts a string to RoleType
func ParseRoleType(role string) RoleType {
	switch role {
	case "waiter":
		return RoleWaiter
	case "cashier":
		return RoleCashier
	case "barista":
		return RoleBarista
	default:
		return RoleWaiter // default fallback
	}
}

// IsValid checks if the RoleType is valid
func (r RoleType) IsValid() bool {
	switch r {
	case RoleWaiter, RoleCashier, RoleBarista:
		return true
	default:
		return false
	}
}

// String returns the string representation
func (r RoleType) String() string {
	return string(r)
}
```

### 2. Fixed Handler to Convert Properly

Updated `shift_handler.go` to:
1. Import `user` package
2. Type assert as `user.Role`
3. Convert to string
4. Use `ParseRoleType()` helper

```go
import (
	"cafe-pos/backend/domain/user"
	// ...
)

func (h *ShiftHandler) StartShift(c *gin.Context) {
	// ...
	role, _ := c.Get("role")
	
	// Convert user.Role to order.RoleType
	roleType := order.ParseRoleType(string(role.(user.Role)))
	
	shift, err := h.shiftService.StartShift(c.Request.Context(), &req, userID.(string), username.(string), roleType)
	// ...
}
```

### 3. Applied Same Fix to All Handler Methods

- `StartShift()` ✅
- `GetCurrentShift()` ✅
- `GetMyShifts()` ✅

## Type Hierarchy

```
user.Role (domain/user)
    ↓ string conversion
order.RoleType (domain/order)
```

Both are `string` types but in different packages:
- `user.Role`: Used for authentication and authorization
- `order.RoleType`: Used for shift management

## Why Two Different Types?

### Separation of Concerns
- **user.Role**: Authentication domain
  - Includes "manager" role
  - Used for access control
  - Stored in user collection

- **order.RoleType**: Shift management domain
  - Only operational roles (waiter, barista, cashier)
  - Used for shift tracking
  - Stored in shift collection

### Benefits
1. **Domain Isolation**: Each domain has its own types
2. **Type Safety**: Compiler catches misuse
3. **Flexibility**: Can add roles to one domain without affecting the other
4. **Clear Intent**: Code explicitly shows conversion between domains

## Testing

After fix, all operations work correctly:
- ✅ Start shift as waiter
- ✅ Start shift as barista
- ✅ Start shift as cashier
- ✅ Get current shift by role
- ✅ Get shift history by role

## Lessons Learned

1. **Type Assertions Need Care**: Always check the actual type in context
2. **Helper Functions**: Create conversion helpers to avoid duplication
3. **Domain Types**: Different domains can have similar but separate types
4. **Error Messages**: Go's error messages clearly indicate type mismatches

## Related Files

- `backend/domain/order/shift.go` - Added ParseRoleType()
- `backend/interfaces/http/shift_handler.go` - Fixed type conversion
- `backend/interfaces/http/middleware.go` - Sets user.Role in context
- `backend/domain/user/user.go` - Defines user.Role type
