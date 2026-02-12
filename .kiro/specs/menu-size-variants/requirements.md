# Menu Size Variants - Requirements

## 1. Overview

### 1.1 Feature Description
Cho phép một món trong menu có nhiều kích cỡ (size variants) khác nhau, mỗi size có giá và công thức nguyên liệu riêng.

### 1.2 Business Value
- Giảm 66% số lượng menu items cần quản lý
- UX tốt hơn: khách chọn món → chọn size (natural flow)
- Dễ maintain: update description/category 1 lần cho tất cả sizes
- Scalable: foundation cho future features (toppings, customizations)

### 1.3 Context
- **Current State**: Chưa có menu items trong database
- **Advantage**: Không cần migration, implement clean từ đầu
- **Risk Level**: Medium (chỉ implementation complexity)

## 2. User Stories

### 2.1 As a Manager - Menu Management

**US-1: Create Single-Size Menu Item**
```
As a manager
I want to create a menu item without size variants
So that I can add simple items like "Bánh mì" with one price
```

**Acceptance Criteria**:
- [ ] AC-1.1: Form có checkbox "Món có nhiều size"
- [ ] AC-1.2: Khi unchecked, hiển thị single price input
- [ ] AC-1.3: Khi unchecked, hiển thị single ingredients selector
- [ ] AC-1.4: Save thành công với has_variants = false
- [ ] AC-1.5: Item hiển thị với 1 giá duy nhất

**US-2: Create Multi-Size Menu Item**
```
As a manager
I want to create a menu item with multiple size variants
So that I can offer "Cà phê sữa đá" in sizes M, L, XL
```

**Acceptance Criteria**:
- [ ] AC-2.1: Form có checkbox "Món có nhiều size"
- [ ] AC-2.2: Khi checked, hiển thị variants form
- [ ] AC-2.3: Có thể thêm nhiều variants (button "Thêm size")
- [ ] AC-2.4: Mỗi variant có: ID, Name, Price, Ingredients
- [ ] AC-2.5: Mỗi variant có checkbox "Mặc định"
- [ ] AC-2.6: Có thể xóa variant (button "Xóa")
- [ ] AC-2.7: Save thành công với has_variants = true
- [ ] AC-2.8: Item hiển thị với nhiều size options

**US-3: Edit Menu Item**
```
As a manager
I want to edit an existing menu item
So that I can update prices, ingredients, or add/remove sizes
```

**Acceptance Criteria**:
- [ ] AC-3.1: Load existing item data vào form
- [ ] AC-3.2: Có thể toggle has_variants
- [ ] AC-3.3: Có thể edit variant details
- [ ] AC-3.4: Có thể add/remove variants
- [ ] AC-3.5: Save updates thành công
- [ ] AC-3.6: Changes reflect immediately trong menu

**US-4: View Menu Items**
```
As a manager
I want to view all menu items with their variants
So that I can see the complete menu structure
```

**Acceptance Criteria**:
- [ ] AC-4.1: Single-size items hiển thị với 1 giá
- [ ] AC-4.2: Multi-size items hiển thị tất cả variants
- [ ] AC-4.3: Mỗi variant hiển thị name và price
- [ ] AC-4.4: Có thể edit từng item
- [ ] AC-4.5: Có thể delete item (cả variants)

### 2.2 As a Waiter - Order Taking

**US-5: Order Single-Size Item**
```
As a waiter
I want to add a single-size item to an order
So that I can quickly add items without size selection
```

**Acceptance Criteria**:
- [ ] AC-5.1: Tap item → immediately added to order
- [ ] AC-5.2: No size selection required
- [ ] AC-5.3: Order shows item name and price
- [ ] AC-5.4: Can adjust quantity

**US-6: Order Multi-Size Item**
```
As a waiter
I want to select a size when adding a multi-size item
So that I can order the correct variant
```

**Acceptance Criteria**:
- [ ] AC-6.1: Tap item → size options appear
- [ ] AC-6.2: Each size shows name and price
- [ ] AC-6.3: Tap size → added to order
- [ ] AC-6.4: Order shows item name + variant name
- [ ] AC-6.5: Order shows correct price for selected size
- [ ] AC-6.6: Can adjust quantity

**US-7: View Order with Variants**
```
As a waiter
I want to see variant names in the order
So that I know exactly what was ordered
```

**Acceptance Criteria**:
- [ ] AC-7.1: Single-size items show: "Bánh mì"
- [ ] AC-7.2: Multi-size items show: "Cà phê sữa đá (Size M)"
- [ ] AC-7.3: Receipt displays variant names
- [ ] AC-7.4: Kitchen display shows variant names

### 2.3 As a Cashier - Payment Processing

**US-8: Process Payment with Variants**
```
As a cashier
I want to see variant details in orders
So that I can verify the correct items and prices
```

**Acceptance Criteria**:
- [ ] AC-8.1: Order list shows variant names
- [ ] AC-8.2: Prices match selected variants
- [ ] AC-8.3: Receipt prints variant names
- [ ] AC-8.4: Total calculation is correct

### 2.4 As a Barista - Order Fulfillment

**US-9: View Order Items with Variants**
```
As a barista
I want to see variant details in orders
So that I can prepare the correct size
```

**Acceptance Criteria**:
- [ ] AC-9.1: Order queue shows variant names
- [ ] AC-9.2: Can see size clearly (M/L/XL)
- [ ] AC-9.3: Ingredients list matches variant

## 3. Functional Requirements

### 3.1 Data Model

**FR-1: MenuItem Structure**
- [ ] FR-1.1: MenuItem has `has_variants` boolean field
- [ ] FR-1.2: MenuItem has `variants` array field (optional)
- [ ] FR-1.3: MenuItem has `price` field for single-size (optional)
- [ ] FR-1.4: MenuItem has `ingredients` array for single-size (optional)
- [ ] FR-1.5: Backward compatible with existing structure

**FR-2: MenuItemVariant Structure**
- [ ] FR-2.1: Variant has `id` string field (e.g., "M", "L", "XL")
- [ ] FR-2.2: Variant has `name` string field (e.g., "Size M")
- [ ] FR-2.3: Variant has `price` float64 field
- [ ] FR-2.4: Variant has `ingredients` array field
- [ ] FR-2.5: Variant has `available` boolean field
- [ ] FR-2.6: Variant has `is_default` boolean field
- [ ] FR-2.7: Variant has cost tracking fields (current_cost, cost_status, cost_last_calculated_at)

**FR-3: OrderItem Structure**
- [ ] FR-3.1: OrderItem has `variant_id` string field (optional)
- [ ] FR-3.2: OrderItem has `variant_name` string field (optional)
- [ ] FR-3.3: OrderItem stores price from selected variant
- [ ] FR-3.4: OrderItem references menu_item_id

### 3.2 Business Logic

**FR-4: Validation Rules**
- [ ] FR-4.1: If has_variants = true, must have at least 1 variant
- [ ] FR-4.2: If has_variants = true, must have exactly 1 default variant
- [ ] FR-4.3: Variant IDs must be unique within a menu item
- [ ] FR-4.4: Variant prices must be > 0
- [ ] FR-4.5: If has_variants = false, must have price > 0
- [ ] FR-4.6: Cannot delete last variant
- [ ] FR-4.7: Cannot have 0 or multiple default variants

**FR-5: Order Creation**
- [ ] FR-5.1: For single-size items, use item.price
- [ ] FR-5.2: For multi-size items, variant_id is required
- [ ] FR-5.3: For multi-size items, validate variant_id exists
- [ ] FR-5.4: For multi-size items, use variant.price
- [ ] FR-5.5: Store variant_name in order for display

**FR-6: Cost Calculation**
- [ ] FR-6.1: For single-size items, calculate cost from item.ingredients
- [ ] FR-6.2: For multi-size items, calculate cost per variant
- [ ] FR-6.3: Store cost per variant separately
- [ ] FR-6.4: Update costs when ingredient prices change

### 3.3 API Requirements

**FR-7: Menu API**
- [ ] FR-7.1: POST /api/menu - Create item with/without variants
- [ ] FR-7.2: GET /api/menu - List all items with variants
- [ ] FR-7.3: GET /api/menu/:id - Get item with variants
- [ ] FR-7.4: PUT /api/menu/:id - Update item and variants
- [ ] FR-7.5: DELETE /api/menu/:id - Delete item and variants

**FR-8: Order API**
- [ ] FR-8.1: POST /api/orders - Create order with variant_id
- [ ] FR-8.2: Validate variant_id if provided
- [ ] FR-8.3: Return error if variant_id invalid
- [ ] FR-8.4: Return error if variant_id missing for multi-size item

## 4. Non-Functional Requirements

### 4.1 Performance

**NFR-1: Response Time**
- [ ] NFR-1.1: Menu list API < 500ms
- [ ] NFR-1.2: Create menu item < 1s
- [ ] NFR-1.3: Order creation < 1s
- [ ] NFR-1.4: Cost calculation per variant < 2s

**NFR-2: Scalability**
- [ ] NFR-2.1: Support up to 10 variants per item
- [ ] NFR-2.2: Support up to 1000 menu items
- [ ] NFR-2.3: Efficient queries with indexes

### 4.2 Usability

**NFR-3: User Interface**
- [ ] NFR-3.1: Form is intuitive and easy to use
- [ ] NFR-3.2: Variants display clearly on mobile
- [ ] NFR-3.3: Size selection is quick (< 3 taps)
- [ ] NFR-3.4: Error messages are clear
- [ ] NFR-3.5: Responsive design for all screen sizes

**NFR-4: Accessibility**
- [ ] NFR-4.1: Keyboard navigation works
- [ ] NFR-4.2: Touch targets are at least 44x44px
- [ ] NFR-4.3: Color contrast meets WCAG AA
- [ ] NFR-4.4: Labels are descriptive

### 4.3 Reliability

**NFR-5: Data Integrity**
- [ ] NFR-5.1: Validation prevents invalid states
- [ ] NFR-5.2: Transactions are atomic
- [ ] NFR-5.3: No orphaned variants
- [ ] NFR-5.4: Cost calculations are accurate

**NFR-6: Error Handling**
- [ ] NFR-6.1: Graceful error messages
- [ ] NFR-6.2: Rollback on failure
- [ ] NFR-6.3: Log errors for debugging

### 4.4 Maintainability

**NFR-7: Code Quality**
- [ ] NFR-7.1: Unit test coverage > 80%
- [ ] NFR-7.2: Integration tests for critical flows
- [ ] NFR-7.3: Clear code documentation
- [ ] NFR-7.4: Follow existing code patterns

## 5. Constraints

### 5.1 Technical Constraints
- Must use existing tech stack (Go, Vue.js, MongoDB)
- Must maintain backward compatibility
- Must not break existing order system
- Must work on mobile devices

### 5.2 Business Constraints
- No migration needed (no existing data)
- Implementation timeline: 2 weeks
- Must be production-ready
- Must be tested thoroughly

## 6. Assumptions

- Users understand the concept of size variants
- Most items will have 2-4 variants (M, L, XL)
- Variant IDs will be simple strings (M, L, XL)
- Default variant is the most commonly ordered size
- Ingredients can differ significantly between sizes

## 7. Dependencies

### 7.1 Internal Dependencies
- Existing menu management system
- Existing order system
- Existing cost calculation system
- Existing ingredient management

### 7.2 External Dependencies
- None

## 8. Success Metrics

### 8.1 Functional Metrics
- [ ] Can create single-size items
- [ ] Can create multi-size items with variants
- [ ] Can order items with variant selection
- [ ] Variants display correctly in all views
- [ ] Cost calculation works per variant
- [ ] All acceptance criteria met

### 8.2 Quality Metrics
- [ ] 0 critical bugs in production
- [ ] Unit test coverage > 80%
- [ ] All integration tests passing
- [ ] Performance requirements met

### 8.3 User Satisfaction
- [ ] Manager can create menu items in < 3 minutes
- [ ] Waiter can order items in < 5 seconds per item
- [ ] No user confusion reported
- [ ] Positive feedback from users

## 9. Out of Scope

The following are explicitly out of scope for this feature:

- Multi-dimensional variants (size + temperature)
- Toppings/add-ons system
- Customizations beyond size
- Promotional pricing per variant
- Variant-specific images
- Nutritional information per variant
- Migration of existing data (no data exists)
- Bulk import/export of variants
- Variant templates/presets (future enhancement)

## 10. Future Enhancements

Potential future additions (not in current scope):

- Quick size templates (M/L/XL preset)
- Bulk edit variants across multiple items
- Variant analytics and reporting
- A/B testing different variant prices
- Seasonal variants
- Time-based variant availability
- Variant-specific promotions
- Multi-dimensional variants (size + temperature + sweetness)
