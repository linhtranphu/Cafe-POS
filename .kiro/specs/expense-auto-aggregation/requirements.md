# Requirements Document: Expense Auto-Aggregation

## Introduction

This feature enables automatic aggregation of individual expense records into operating expense categories. Users can categorize expenses during creation and later aggregate them by date range and type into the Operating Expenses section of Shop Settings. This provides traceability between individual expenses and operating expense summaries while reducing manual data entry.

## Glossary

- **Expense**: An individual expense record created in the Expense Management module
- **Expense_Type**: A categorization field for expenses with five predefined values
- **Operating_Expense**: A summary expense category stored in Shop Settings for financial reporting
- **Aggregation**: The process of summing individual expenses by type and date range
- **Breakdown**: Detailed list showing which individual expenses contributed to an aggregated amount
- **Date_Range**: A start and end date used to filter expenses for aggregation
- **Expense_Management**: The module where individual expenses are created and tracked
- **Shop_Settings**: The configuration area containing operating expense summaries
- **Traceability**: The ability to identify which individual expenses contributed to an operating expense total

## Requirements

### Requirement 1: Expense Type Selection During Creation

**User Story:** As a user creating an expense, I want to select an expense type from predefined categories, so that expenses can be automatically categorized for later aggregation.

#### Acceptance Criteria

1. WHEN a user creates a new expense, THE Expense_Form SHALL display an expense_type selection field
2. THE Expense_Type_Field SHALL provide exactly five options: staff_salary, rent, utilities, marketing_costs, other_expenses
3. WHEN a user attempts to save an expense without selecting an expense_type, THE System SHALL prevent the save operation and display a validation error
4. WHEN an expense is saved with a valid expense_type, THE System SHALL store the expense_type value in the expense record
5. THE Expense_Type_Field SHALL display Vietnamese labels: "Lương nhân viên", "Thuê mặt bằng", "Điện nước", "Marketing", "Khác"

### Requirement 2: Expense Type Data Storage

**User Story:** As a system, I want to store expense_type with each expense record, so that expenses can be filtered and aggregated by type.

#### Acceptance Criteria

1. THE Expense_Model SHALL include an expense_type field
2. THE expense_type field SHALL accept only the five predefined values: staff_salary, rent, utilities, marketing_costs, other_expenses
3. WHEN an expense is created, THE System SHALL persist the expense_type value to the database
4. WHEN an expense is retrieved, THE System SHALL return the expense_type value
5. THE expense_type field SHALL be indexed for efficient querying

### Requirement 3: Expense Aggregation API

**User Story:** As a backend system, I want to provide an API endpoint that aggregates expenses by date range and type, so that the frontend can display aggregated totals.

#### Acceptance Criteria

1. THE System SHALL provide a GET endpoint at /api/expenses/aggregate
2. WHEN the endpoint receives a start_date and end_date parameter, THE System SHALL return aggregated totals for all expense types within that date range
3. THE Aggregation_Response SHALL include totals for each of the five expense types: staff_salary, rent, utilities, marketing_costs, other_expenses
4. THE Aggregation_Response SHALL include a breakdown array for each expense type showing individual expense IDs, amounts, dates, and descriptions
5. WHEN no expenses exist in the date range, THE System SHALL return zero values for all expense types
6. WHEN the date range is invalid (end_date before start_date), THE System SHALL return an error response
7. THE System SHALL calculate totals by summing the amount field of all expenses matching the expense_type and date range

### Requirement 4: Operating Expense Aggregation UI

**User Story:** As a user managing operating expenses, I want to access an auto-aggregation interface in Settings, so that I can automatically populate operating expenses from individual expense records.

#### Acceptance Criteria

1. WHEN a user navigates to Settings → Operating Expenses, THE System SHALL display an "Auto-aggregate from Expenses" button
2. WHEN the user clicks the auto-aggregation button, THE System SHALL display a modal or panel with date range selection
3. THE Date_Range_Selector SHALL include a start_date field and an end_date field
4. WHEN the user selects a date range and clicks "Aggregate", THE System SHALL call the aggregation API and display results
5. THE Aggregation_Results SHALL display five rows, one for each expense type, showing the Vietnamese label and aggregated total
6. THE Aggregation_Results SHALL allow users to click on each expense type to view the breakdown of individual expenses
7. THE Breakdown_View SHALL display a table with columns: date, description, amount for each contributing expense
8. THE Aggregation_Results SHALL include editable amount fields allowing users to modify totals before saving
9. THE Aggregation_Results SHALL include a "Save to Operating Expenses" button
10. WHEN the user clicks "Save to Operating Expenses", THE System SHALL update the Shop Settings operating expense fields with the aggregated (or edited) values

### Requirement 5: Traceability and Breakdown Display

**User Story:** As a user reviewing aggregated expenses, I want to see which individual expenses contributed to each operating expense category, so that I can verify accuracy and understand the composition.

#### Acceptance Criteria

1. WHEN displaying aggregated results, THE System SHALL provide a clickable link or expand button for each expense type
2. WHEN a user clicks to view breakdown details, THE System SHALL display all individual expenses that contributed to that category's total
3. THE Breakdown_Display SHALL show expense date, description, and amount for each contributing expense
4. THE Breakdown_Display SHALL include a link to view the full expense record in Expense Management
5. THE Breakdown_Display SHALL calculate and display the sum of all breakdown items matching the aggregated total

### Requirement 6: Operating Expense Update Integration

**User Story:** As a user, I want to save aggregated expense totals directly to Operating Expenses, so that I can update financial summaries without manual data entry.

#### Acceptance Criteria

1. WHEN a user clicks "Save to Operating Expenses", THE System SHALL update the Shop Settings document with the aggregated values
2. THE System SHALL map aggregated values to the corresponding operating expense fields: staff_salary, rent, utilities, marketing_costs, other_expenses
3. WHEN the save operation succeeds, THE System SHALL display a success message
4. WHEN the save operation fails, THE System SHALL display an error message and preserve the aggregated data for retry
5. THE System SHALL allow users to edit aggregated amounts before saving
6. WHEN a user edits an aggregated amount, THE System SHALL preserve the edited value for saving (not recalculate from breakdown)

### Requirement 7: Date Range Validation and Filtering

**User Story:** As a system, I want to validate date ranges and filter expenses accurately, so that aggregation results are correct and reliable.

#### Acceptance Criteria

1. WHEN a user selects a date range, THE System SHALL validate that start_date is not after end_date
2. IF the date range is invalid, THEN THE System SHALL display a validation error and prevent aggregation
3. THE System SHALL filter expenses inclusively (including expenses on start_date and end_date)
4. WHEN aggregating expenses, THE System SHALL use the expense creation date or transaction date for filtering
5. THE System SHALL handle timezone considerations consistently between frontend and backend

### Requirement 8: Error Handling and Edge Cases

**User Story:** As a system, I want to handle errors and edge cases gracefully, so that users have a reliable experience.

#### Acceptance Criteria

1. WHEN the aggregation API fails, THE System SHALL display a user-friendly error message
2. WHEN no expenses exist in the selected date range, THE System SHALL display a message indicating zero expenses found
3. WHEN network errors occur during aggregation, THE System SHALL allow users to retry the operation
4. WHEN saving to Operating Expenses fails, THE System SHALL preserve the aggregated data and allow retry
5. THE System SHALL handle concurrent updates to Operating Expenses by using appropriate locking or versioning mechanisms

### Requirement 9: User Experience and Feedback

**User Story:** As a user, I want clear feedback during the aggregation process, so that I understand what the system is doing and can verify results.

#### Acceptance Criteria

1. WHEN aggregation is in progress, THE System SHALL display a loading indicator
2. WHEN aggregation completes, THE System SHALL display the results immediately
3. THE System SHALL display the date range used for aggregation in the results view
4. THE System SHALL display the total number of expenses included in the aggregation
5. WHEN saving to Operating Expenses, THE System SHALL display a confirmation message with the updated values

### Requirement 10: Existing Expense Records Migration

**User Story:** As a system administrator, I want existing expense records to be compatible with the new expense_type field, so that the system continues to function with historical data.

#### Acceptance Criteria

1. WHEN an existing expense record without expense_type is retrieved, THE System SHALL handle it gracefully without errors
2. THE System SHALL allow existing expenses to be edited and assigned an expense_type
3. WHEN aggregating expenses, THE System SHALL exclude expenses without an expense_type from the results
4. THE System SHALL provide a way to identify expenses missing the expense_type field for data cleanup
