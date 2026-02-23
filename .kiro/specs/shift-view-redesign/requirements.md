# Shift View Redesign - Requirements

## Overview
Redesign the Shift View (`/shifts`) for Cashier role to provide a better UX for viewing shifts, payments, and managing payment adjustments. The redesign will separate waiter shifts and barista shifts since they have different management needs.

**IMPORTANT**: This redesign ONLY affects the Cashier role view. Waiter and Barista views remain unchanged.

## User Stories

### 1. As a Cashier, I want to view and manage waiter shifts efficiently
**Acceptance Criteria:**
- 1.1 I can see a list of waiter shifts (open and closed) in a clean, organized layout
- 1.2 I can filter shifts by date to find specific shifts quickly
- 1.3 I can select a shift to view detailed information including:
  - Shift basic info (type, start time, end time, user name)
  - Financial summary (start cash, end cash, total revenue, cash collected, transfer collected)
  - List of all payments in that shift
- 1.4 I can see payment details for each order in the selected shift
- 1.5 I can perform payment adjustments (override, report discrepancy, lock order) directly from the shift view
- 1.6 The UI clearly distinguishes between open and closed shifts
- 1.7 I can close open waiter shifts from this view

### 2. As a Cashier, I want to view barista shifts separately
**Acceptance Criteria:**
- 2.1 I can switch between "Waiter Shifts" and "Barista Shifts" views using tabs
- 2.2 Barista shift view shows simplified information (no payment management)
- 2.3 Barista shift view shows:
  - Shift basic info (type, start time, end time, user name)
  - Order statistics (total orders processed in the shift)
- 2.4 The UI is simpler for barista shifts since there are fewer management options
- 2.5 I can close open barista shifts from this view

### 3. As a Waiter, my shift view remains unchanged
**Acceptance Criteria:**
- 3.1 I see the same UI as before (no changes)
- 3.2 I can manage my own shift (start, handover, end)
- 3.3 I can view my handover history
- 3.4 I can see my shift history

### 4. As a Barista, my shift view remains unchanged
**Acceptance Criteria:**
- 4.1 I see the same UI as before (no changes)
- 4.2 I can manage my own shift (start, end)
- 4.3 I can see my shift history

## Design Goals (Cashier View Only)

### Separation of Concerns
- **Waiter Shifts Tab**: Focus on payment management, cash handling, financial overview
- **Barista Shifts Tab**: Focus on order statistics, simpler view

### Information Hierarchy
1. **Primary**: Tab selector (Waiter Shifts / Barista Shifts)
2. **Secondary**: Date filter + Shift selector dropdown
3. **Tertiary**: Shift summary card (financial/stats overview)
4. **Quaternary**: Detailed information (payments list with actions)

### UI/UX Improvements (Cashier Only)
- Use tabs to separate waiter and barista shifts
- Use dropdown selector to choose specific shift (similar to CashierDashboard)
- Show shift summary in gradient card (matching CashierDashboard style)
- Show payment list inline with shift details (no need to navigate away)
- Group payment actions together (override, discrepancy, lock)
- Use color coding for shift status (open = green, closed = gray)
- Use consistent styling with CashierDashboard for familiarity

### No Changes for Waiter/Barista
- Waiter view: Keep existing UI with handover functionality
- Barista view: Keep existing UI with simple shift management

## Technical Requirements

### Conditional Rendering
- Detect user role (Cashier vs Waiter vs Barista)
- Show redesigned UI only for Cashier
- Show existing UI for Waiter and Barista (no changes)

### Data Fetching (Cashier View)
- Fetch all shifts (waiter and barista) for cashier
- Fetch payments for selected waiter shift
- Fetch shift statistics for selected shift

### State Management (Cashier View)
- Track active tab (waiter/barista)
- Track selected date
- Track selected shift ID
- Track payment actions (override, discrepancy, lock)

### API Endpoints Used
- `GET /api/shifts` - Fetch all shifts (cashier can see all)
- `GET /api/cashier/shifts/:id/status` - Get shift status with stats
- `GET /api/cashier/shifts/:id/payments` - Get payments for a shift
- `PATCH /api/cashier/payments/:orderId/override` - Override payment
- `POST /api/cashier/discrepancies` - Report discrepancy
- `PATCH /api/cashier/orders/:orderId/lock` - Lock order
- `POST /api/manager/shifts/:id/close` - Close shift (cashier/manager only)

## Non-Functional Requirements

### Performance
- Load shifts efficiently with pagination if needed
- Cache shift data to avoid repeated API calls
- Lazy load payment details when shift is selected

### Accessibility
- Use semantic HTML
- Ensure proper color contrast
- Support keyboard navigation
- Provide clear labels and descriptions

### Responsive Design
- Mobile-first design
- Touch-friendly buttons and controls
- Proper spacing for mobile devices
- Scrollable content areas

## Out of Scope
- Changes to Waiter shift view (keep existing UI)
- Changes to Barista shift view (keep existing UI)
- Creating new shifts (handled elsewhere)
- Closing cashier's own shifts (handled in CashierDashboard)
- Order creation/modification (handled in Orders view)
- Handover management (handled in separate Handover view)
- Reprint functionality (can be added later if needed)

## Success Metrics (Cashier View)
- Reduced time to find and view shift details
- Reduced clicks to perform payment adjustments
- Clear separation between waiter and barista shift management
- Improved visual clarity and information hierarchy
- Consistent UX with CashierDashboard
