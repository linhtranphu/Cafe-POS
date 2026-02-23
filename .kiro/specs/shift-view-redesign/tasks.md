# Shift View Redesign - Implementation Tasks

## 1. Refactor ShiftView Component Structure
- [x] 1.1 Add role detection logic to conditionally render Cashier vs Waiter/Barista UI
- [x] 1.2 Wrap existing UI in conditional block for Waiter/Barista roles
- [x] 1.3 Create new section for Cashier UI
- [x] 1.4 Import required components and stores (cashierStore, modals)

## 2. Implement Cashier Tab Navigation
- [x] 2.1 Add tab state management (activeTab: 'waiter' | 'barista')
- [x] 2.2 Create tab buttons UI with icons and labels
- [x] 2.3 Add tab switching logic
- [x] 2.4 Apply active tab styling (blue background for active, gray for inactive)
- [x] 2.5 Clear selectedShiftId when tab changes

## 3. Implement Date Filter and Shift Selector
- [x] 3.1 Add selectedDate state (default to today)
- [x] 3.2 Add selectedShiftId state
- [x] 3.3 Create date picker input
- [x] 3.4 Create filteredShifts computed property (filter by date and role_type)
- [x] 3.5 Create shift selector dropdown
- [x] 3.6 Populate dropdown with filtered shifts
- [x] 3.7 Display shift info in dropdown options (date, user name, role type)
- [x] 3.8 Clear selectedShiftId when date changes

## 4. Implement Shift Summary Card for Waiter Shifts
- [x] 4.1 Create gradient card container (yellow-orange gradient)
- [x] 4.2 Display shift basic info (user name, shift type, time, status)
- [x] 4.3 Display financial stats grid (4 cards: total orders, revenue, cash, transfer)
- [x] 4.4 Use backdrop-blur effect for stat cards
- [x] 4.5 Add "Chốt ca" button for open shifts
- [x] 4.6 Fetch shift status when shift is selected

## 5. Implement Shift Summary Card for Barista Shifts
- [x] 5.1 Create gradient card container (purple-indigo gradient)
- [x] 5.2 Display shift basic info (user name, shift type, time, status)
- [x] 5.3 Display order stats grid (2 cards: total orders, in progress)
- [x] 5.4 Use backdrop-blur effect for stat cards
- [x] 5.5 Add "Chốt ca" button for open shifts

## 6. Implement Payment List for Waiter Shifts
- [x] 6.1 Add section header with payment count
- [x] 6.2 Fetch payments when waiter shift is selected
- [x] 6.3 Create payment card component
- [x] 6.4 Display payment info (customer name, time, amount, payment method)
- [x] 6.5 Display payment status badge
- [x] 6.6 Add payment action buttons (override, discrepancy, lock)
- [x] 6.7 Handle empty state (no payments)

## 7. Implement Payment Actions
- [x] 7.1 Add modal state management (showOverride, showDiscrepancy, selectedPayment)
- [x] 7.2 Implement override payment handler
- [x] 7.3 Implement report discrepancy handler
- [x] 7.4 Implement lock order handler with confirmation
- [x] 7.5 Refresh payments after each action
- [x] 7.6 Show success/error messages

## 8. Implement Close Shift Functionality
- [x] 8.1 Add showCloseShift modal state
- [x] 8.2 Create close shift modal with end_cash input
- [x] 8.3 Implement closeShift handler
- [x] 8.4 Call shiftStore.closeShift API
- [x] 8.5 Refresh shifts list after closing
- [x] 8.6 Clear selectedShiftId after closing

## 9. Add Loading and Error States
- [x] 9.1 Add loading spinner while fetching shifts
- [x] 9.2 Add loading spinner while fetching shift details
- [x] 9.3 Add loading spinner while fetching payments
- [x] 9.4 Handle API errors gracefully
- [x] 9.5 Show error messages to user

## 10. Style and Polish
- [x] 10.1 Apply consistent spacing and padding
- [x] 10.2 Ensure mobile responsiveness
- [x] 10.3 Add touch-friendly button sizes
- [x] 10.4 Add transition animations
- [x] 10.5 Test on different screen sizes
- [x] 10.6 Ensure color contrast meets accessibility standards

## 11. Testing and Validation
- [-] 11.1 Test tab switching functionality
- [-] 11.2 Test date filter functionality
- [-] 11.3 Test shift selection and details loading
- [-] 11.4 Test payment actions (override, discrepancy, lock)
- [-] 11.5 Test close shift functionality
- [-] 11.6 Test with different user roles (cashier, waiter, barista)
- [-] 11.7 Test error scenarios
- [x] 11.8 Test on mobile devices

## 12. Documentation and Cleanup
- [x] 12.1 Add code comments for complex logic
- [x] 12.2 Remove unused code and imports
- [x] 12.3 Update component documentation
- [x] 12.4 Verify no console errors or warnings
