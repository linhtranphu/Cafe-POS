# Ingredient Management - Comprehensive Analysis & User Stories

## Current State Analysis

### ✅ What's Working Well

1. **Create Ingredient**
   - Modal with total/unit price toggle
   - Auto-expense tracking
   - Initial stock history creation
   - Category management
   - Mobile-optimized UI

2. **Stock Adjustment**
   - Add/Remove/Set stock
   - Price input for purchases
   - Total/unit price toggle
   - Auto-expense for purchases
   - Weighted average cost calculation

3. **View History**
   - Compact mobile design
   - Price information display
   - Price comparison with current
   - User and timestamp tracking

4. **List View**
   - Search functionality
   - Stock status indicators
   - Stats cards (total, in stock, low stock, out of stock)
   - Quick actions per ingredient
   - Pull to refresh

5. **Technical Features**
   - Double-submit prevention
   - Loading states
   - Input validation
   - Safe area handling (iPhone)
   - Responsive design

### 🔴 Missing Features & Pain Points

1. **No Quick Stock IN/OUT Actions**
   - Current: Must open modal, select type, enter details
   - Need: Quick buttons for common operations

2. **No Bulk Operations**
   - Can't select multiple ingredients
   - Can't bulk adjust prices
   - Can't bulk import/export

3. **No Stock Alerts**
   - Low stock items not prominently shown
   - No notifications
   - No reorder suggestions

4. **Limited Filtering**
   - Only search by name
   - Can't filter by category
   - Can't filter by stock status
   - Can't sort by different criteria

5. **No Analytics**
   - No cost trends
   - No usage patterns
   - No supplier comparison
   - No waste tracking

6. **No Barcode/QR Support**
   - Manual entry only
   - No quick scan

7. **Limited History**
   - Only shows 50 records
   - No date range filter
   - No export capability

8. **No Recipe Integration**
   - Can't see which menu items use this ingredient
   - Can't calculate recipe costs

## User Stories

### Epic 1: Quick Stock Operations

#### US-1.1: Quick Stock IN Button
**As a** manager  
**I want to** quickly add stock with a single tap  
**So that** I can record purchases faster during busy times

**Acceptance Criteria:**
- [ ] "+" button visible on each ingredient card
- [ ] Tapping opens compact modal with quantity and price only
- [ ] Default price is current cost_per_unit
- [ ] Can toggle to enter total price
- [ ] Auto-expense created
- [ ] Stock history recorded
- [ ] Modal closes after success
- [ ] Success feedback shown

**Priority:** HIGH  
**Effort:** 2 points

#### US-1.2: Quick Stock OUT Button
**As a** manager  
**I want to** quickly record stock usage/waste  
**So that** inventory stays accurate

**Acceptance Criteria:**
- [ ] "-" button visible on each ingredient card
- [ ] Tapping opens compact modal with quantity and reason
- [ ] Predefined reasons: "Sử dụng", "Hỏng", "Hết hạn", "Khác"
- [ ] Stock history recorded
- [ ] No expense created (already purchased)
- [ ] Modal closes after success

**Priority:** HIGH  
**Effort:** 2 points

#### US-1.3: Quick View Stock Level
**As a** manager  
**I want to** see stock levels at a glance  
**So that** I know what needs ordering

**Acceptance Criteria:**
- [ ] Progress bar showing current vs min stock
- [ ] Color coded: green (good), yellow (low), red (out)
- [ ] Percentage or days remaining
- [ ] Visual indicator on card

**Priority:** MEDIUM  
**Effort:** 1 point

### Epic 2: Advanced Filtering & Sorting

#### US-2.1: Filter by Category
**As a** manager  
**I want to** filter ingredients by category  
**So that** I can focus on specific types

**Acceptance Criteria:**
- [ ] Category chips/pills at top of list
- [ ] "Tất cả" shows all
- [ ] Tapping category filters list
- [ ] Active category highlighted
- [ ] Count shown per category

**Priority:** HIGH  
**Effort:** 2 points

#### US-2.2: Filter by Stock Status
**As a** manager  
**I want to** filter by stock status  
**So that** I can see what needs attention

**Acceptance Criteria:**
- [ ] Filter options: "Tất cả", "Đủ hàng", "Sắp hết", "Hết hàng"
- [ ] Quick toggle buttons
- [ ] Count shown per status
- [ ] Persists during session

**Priority:** HIGH  
**Effort:** 1 point

#### US-2.3: Sort Options
**As a** manager  
**I want to** sort ingredients by different criteria  
**So that** I can find what I need quickly

**Acceptance Criteria:**
- [ ] Sort by: Name, Category, Quantity, Cost, Last Updated
- [ ] Ascending/Descending toggle
- [ ] Sort icon in header
- [ ] Persists during session

**Priority:** MEDIUM  
**Effort:** 2 points

### Epic 3: Stock Alerts & Notifications

#### US-3.1: Low Stock Banner
**As a** manager  
**I want to** see low stock items prominently  
**So that** I don't run out

**Acceptance Criteria:**
- [ ] Banner at top showing low stock count
- [ ] Tapping shows filtered list
- [ ] Dismissible but reappears on refresh
- [ ] Shows top 3 critical items

**Priority:** HIGH  
**Effort:** 2 points

#### US-3.2: Reorder Suggestions
**As a** manager  
**I want to** get reorder quantity suggestions  
**So that** I order the right amount

**Acceptance Criteria:**
- [ ] Based on usage history
- [ ] Shows average daily usage
- [ ] Suggests quantity to last X days
- [ ] Calculates estimated cost
- [ ] One-tap to create adjustment

**Priority:** MEDIUM  
**Effort:** 5 points

#### US-3.3: Expiry Tracking
**As a** manager  
**I want to** track ingredient expiry dates  
**So that** I can use items before they expire

**Acceptance Criteria:**
- [ ] Optional expiry date field
- [ ] Warning when approaching expiry
- [ ] Filter by expiring soon
- [ ] Batch tracking (FIFO)

**Priority:** LOW  
**Effort:** 8 points

### Epic 4: Bulk Operations

#### US-4.1: Bulk Price Update
**As a** manager  
**I want to** update prices for multiple ingredients  
**So that** I can adjust for supplier price changes

**Acceptance Criteria:**
- [ ] Select multiple ingredients
- [ ] Apply percentage increase/decrease
- [ ] Or set specific prices
- [ ] Preview changes before applying
- [ ] Confirmation dialog
- [ ] History recorded for each

**Priority:** MEDIUM  
**Effort:** 5 points

#### US-4.2: Bulk Stock Adjustment
**As a** manager  
**I want to** adjust stock for multiple ingredients  
**So that** I can do inventory counts efficiently

**Acceptance Criteria:**
- [ ] Select multiple ingredients
- [ ] Enter new quantities
- [ ] Reason applies to all
- [ ] Preview changes
- [ ] Confirmation dialog
- [ ] History recorded for each

**Priority:** MEDIUM  
**Effort:** 5 points

#### US-4.3: Import/Export
**As a** manager  
**I want to** import/export ingredient data  
**So that** I can backup or bulk update

**Acceptance Criteria:**
- [ ] Export to CSV/Excel
- [ ] Import from CSV/Excel
- [ ] Template download
- [ ] Validation on import
- [ ] Preview before import
- [ ] Error reporting

**Priority:** LOW  
**Effort:** 8 points

### Epic 5: Analytics & Insights

#### US-5.1: Cost Trends
**As a** manager  
**I want to** see cost trends over time  
**So that** I can budget better

**Acceptance Criteria:**
- [ ] Line chart showing price changes
- [ ] Selectable time range
- [ ] Compare multiple ingredients
- [ ] Show average, min, max
- [ ] Export chart data

**Priority:** MEDIUM  
**Effort:** 5 points

#### US-5.2: Usage Analytics
**As a** manager  
**I want to** see usage patterns  
**So that** I can optimize ordering

**Acceptance Criteria:**
- [ ] Daily/weekly/monthly usage
- [ ] Peak usage times
- [ ] Waste percentage
- [ ] Cost per period
- [ ] Trend indicators

**Priority:** MEDIUM  
**Effort:** 5 points

#### US-5.3: Supplier Comparison
**As a** manager  
**I want to** compare suppliers  
**So that** I can get best prices

**Acceptance Criteria:**
- [ ] Group by supplier
- [ ] Average prices per supplier
- [ ] Quality ratings
- [ ] Delivery times
- [ ] Total spend per supplier

**Priority:** LOW  
**Effort:** 5 points

### Epic 6: Mobile Optimization

#### US-6.1: Offline Support
**As a** manager  
**I want to** record stock changes offline  
**So that** I can work without internet

**Acceptance Criteria:**
- [ ] Queue operations when offline
- [ ] Sync when back online
- [ ] Show offline indicator
- [ ] Conflict resolution
- [ ] Local storage

**Priority:** LOW  
**Effort:** 13 points

#### US-6.2: Barcode Scanning
**As a** manager  
**I want to** scan barcodes  
**So that** I can quickly find/add ingredients

**Acceptance Criteria:**
- [ ] Camera access
- [ ] Barcode recognition
- [ ] Link barcode to ingredient
- [ ] Quick lookup by scan
- [ ] Add new with scan

**Priority:** MEDIUM  
**Effort:** 8 points

#### US-6.3: Voice Input
**As a** manager  
**I want to** use voice to record stock  
**So that** I can work hands-free

**Acceptance Criteria:**
- [ ] Voice button on forms
- [ ] Speech recognition
- [ ] Parse quantity and item
- [ ] Confirmation before save
- [ ] Works offline

**Priority:** LOW  
**Effort:** 8 points

### Epic 7: Integration Features

#### US-7.1: Recipe Integration
**As a** manager  
**I want to** see which menu items use each ingredient  
**So that** I know impact of stock outs

**Acceptance Criteria:**
- [ ] "Used in" section on ingredient detail
- [ ] List of menu items
- [ ] Quantity needed per item
- [ ] Can navigate to menu item
- [ ] Shows if menu item is available

**Priority:** HIGH  
**Effort:** 5 points

#### US-7.2: Auto-Deduct on Order
**As a** manager  
**I want** ingredients auto-deducted when orders complete  
**So that** inventory stays accurate

**Acceptance Criteria:**
- [ ] Recipe defines ingredient quantities
- [ ] Order completion triggers deduction
- [ ] Stock history shows order reference
- [ ] Can disable per ingredient
- [ ] Batch processing for performance

**Priority:** MEDIUM  
**Effort:** 8 points

#### US-7.3: Purchase Order Generation
**As a** manager  
**I want to** generate purchase orders  
**So that** I can order from suppliers easily

**Acceptance Criteria:**
- [ ] Select ingredients to order
- [ ] Suggested quantities
- [ ] Group by supplier
- [ ] Generate PDF
- [ ] Email to supplier
- [ ] Track order status

**Priority:** LOW  
**Effort:** 13 points

## Implementation Priority

### Phase 1: Quick Wins (Sprint 1-2)
1. US-1.1: Quick Stock IN Button
2. US-1.2: Quick Stock OUT Button
3. US-2.1: Filter by Category
4. US-2.2: Filter by Stock Status
5. US-3.1: Low Stock Banner

**Total Effort:** 9 points  
**Business Value:** HIGH

### Phase 2: Core Features (Sprint 3-4)
1. US-1.3: Quick View Stock Level
2. US-2.3: Sort Options
3. US-7.1: Recipe Integration
4. US-3.2: Reorder Suggestions

**Total Effort:** 10 points  
**Business Value:** HIGH

### Phase 3: Advanced Features (Sprint 5-6)
1. US-4.1: Bulk Price Update
2. US-4.2: Bulk Stock Adjustment
3. US-5.1: Cost Trends
4. US-5.2: Usage Analytics
5. US-6.2: Barcode Scanning

**Total Effort:** 28 points  
**Business Value:** MEDIUM

### Phase 4: Nice to Have (Future)
1. US-3.3: Expiry Tracking
2. US-4.3: Import/Export
3. US-5.3: Supplier Comparison
4. US-6.1: Offline Support
5. US-6.3: Voice Input
6. US-7.2: Auto-Deduct on Order
7. US-7.3: Purchase Order Generation

**Total Effort:** 63 points  
**Business Value:** LOW-MEDIUM

## Technical Considerations

### Frontend
- Vue 3 Composition API
- Pinia for state management
- Tailwind CSS for styling
- Mobile-first responsive design
- Progressive Web App features

### Backend
- Go with Gin framework
- MongoDB for data storage
- RESTful API design
- JWT authentication
- Background jobs for heavy operations

### Performance
- Pagination for large lists
- Lazy loading for images
- Debounced search
- Optimistic UI updates
- Caching strategies

### Security
- Role-based access control
- Input validation
- SQL injection prevention
- XSS protection
- CSRF tokens

## Success Metrics

### User Satisfaction
- Time to record stock change < 10 seconds
- Error rate < 1%
- User adoption > 80%
- Daily active users

### Business Impact
- Reduce stock-outs by 50%
- Reduce waste by 30%
- Improve inventory accuracy to 95%+
- Reduce time spent on inventory by 40%

### Technical Metrics
- Page load time < 2 seconds
- API response time < 500ms
- Uptime > 99.9%
- Mobile performance score > 90

## Next Steps

1. **Review & Prioritize**: Stakeholder meeting to confirm priorities
2. **Design Phase**: Create detailed mockups for Phase 1
3. **Technical Spike**: Investigate barcode scanning options
4. **Sprint Planning**: Break down Phase 1 into tasks
5. **Development**: Start with US-1.1 (Quick Stock IN)
