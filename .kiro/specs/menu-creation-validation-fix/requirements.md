# Requirements: Menu Creation Validation Fix

## 1. Problem Statement

When users attempt to create a new menu item through the menu management interface, the API request is sent correctly but the backend returns an error: **"Lỗi tạo món"** (Error creating menu item).

### Root Cause

The backend validation in `MenuItem.Validate()` requires:
- For single-size items: `price must be > 0`
- For multi-size items: `variants` array must not be empty and have exactly one default variant

However, the frontend form allows users to submit without entering a price, or with price = 0, which fails backend validation.

**Backend Validation Code** (`backend/domain/menu/menu.go`):
```go
if m.HasVariants {
    // Multi-size validation
    if len(m.Variants) == 0 {
        return fmt.Errorf("variants required when has_variants=true")
    }
    // Must have exactly 1 default
    // ...
} else {
    // Single-size validation
    if m.Price <= 0 {
        return fmt.Errorf("price must be > 0 for single-size item")
    }
}
```

**Frontend Form Initialization** (`frontend/src/views/Men