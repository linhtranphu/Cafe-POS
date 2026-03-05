# Requirements Document

## Introduction

This feature implements automatic and manual menu cost recalculation to ensure menu item costs stay synchronized with batch costs. Currently, menu item costs are read from the database without recalculation, causing stale cost data when new batches are created with different costs. This feature adds two mechanisms: automatic recalculation when batches are created, and a manual "Recalculate All" button for on-demand updates.

## Glossary

- **Menu_Item**: A sellable item on the menu with a price and cost
- **Batch_Definition**: A recipe/formula for producing a batch (e.g., "Cà phê cốt")
- **Batch_Record**: An actual production instance of a batch with specific cost and expiry
- **Cost_Recalculation_Service**: Background service that processes cost recalculation jobs
- **Menu_Repository**: Data access layer for menu items
- **Batch_Record_Service**: Service handling batch creation and management
- **FIFO**: First-In-First-Out inventory tracking method
- **Cost_Per_Unit**: The cost of one unit of a batch or ingredient

## Requirements

### Requirement 1: Automatic Cost Recalculation on Batch Creation

**User Story:** As a barista, I want menu costs to automatically update when I create a new batch, so that cost data stays accurate without manual intervention.

#### Acceptance Criteria

1. WHEN a barista creates a new batch record, THE Batch_Record_Service SHALL identify all menu items using that batch definition
2. WHEN menu items are identified, THE Batch_Record_Service SHALL queue cost recalculation for each affected menu item
3. WHEN queueing recalculation, THE Batch_Record_Service SHALL use a background goroutine to prevent blocking batch creation
4. WHEN the recalculation queue is processed, THE Cost_Recalculation_Service SHALL update the current_cost field for each menu item
5. IF the recalculation queueing fails, THEN THE Batch_Record_Service SHALL log a warning and continue batch creation successfully

### Requirement 2: Menu Repository Batch Lookup

**User Story:** As the system, I want to efficiently find all menu items using a specific batch definition, so that I can recalculate only affected items.

#### Acceptance Criteria

1. THE Menu_Repository SHALL provide a FindByBatchDefinitionID method that accepts a batch definition ID
2. WHEN searching for menu items, THE Menu_Repository SHALL check both single-size ingredients and multi-size variant ingredients
3. WHEN a batch definition ID is provided, THE Menu_Repository SHALL return all menu items containing that batch in any ingredient list
4. THE Menu_Repository SHALL return an empty list if no menu items use the specified batch definition

### Requirement 3: Manual Recalculation API

**User Story:** As a manager, I want to manually trigger cost recalculation for all menu items via an API, so that I can refresh all costs on demand.

#### Acceptance Criteria

1. THE System SHALL provide a POST endpoint at /api/manager/menu/costs/recalculate-all
2. WHEN the recalculate-all endpoint is called, THE System SHALL fetch all menu items from the database
3. WHEN menu items are fetched, THE System SHALL queue cost recalculation for each menu item
4. WHEN queueing is complete, THE System SHALL return a response containing total items count, queued count, and failed count
5. THE System SHALL process recalculation jobs in the background without blocking the API response

### Requirement 4: Manual Recalculation UI

**User Story:** As a manager, I want a "Recalculate All" button in the Menu Costs view, so that I can manually refresh all menu costs with one click.

#### Acceptance Criteria

1. WHEN a manager views the Menu Costs page, THE System SHALL display a "Tính lại tất cả" button
2. WHEN the button is clicked, THE System SHALL disable the button and show a loading state
3. WHEN the recalculation request completes successfully, THE System SHALL display a success message with the count of items processed
4. WHEN the recalculation request completes, THE System SHALL refresh the menu costs data after a short delay
5. IF the recalculation request fails, THEN THE System SHALL display an error message with details
6. WHERE the device is a desktop, THE System SHALL display the button with full text and icon
7. WHERE the device is mobile, THE System SHALL display the button with icon only

### Requirement 5: Background Processing

**User Story:** As the system, I want to process cost recalculations in the background, so that batch creation and UI interactions remain responsive.

#### Acceptance Criteria

1. WHEN cost recalculation is triggered, THE System SHALL use the existing Cost_Recalculation_Service worker pool
2. THE Cost_Recalculation_Service SHALL process recalculation jobs using FIFO queue ordering
3. WHEN processing a recalculation job, THE Cost_Recalculation_Service SHALL calculate cost using the oldest available batch (FIFO inventory)
4. WHEN a recalculation job completes, THE Cost_Recalculation_Service SHALL update the menu_items.current_cost field in the database
5. IF a recalculation job fails, THEN THE Cost_Recalculation_Service SHALL log the error and continue processing other jobs

### Requirement 6: Dependency Wiring

**User Story:** As a developer, I want all service dependencies properly wired, so that the auto-recalculation feature works correctly.

#### Acceptance Criteria

1. THE Batch_Record_Service SHALL have a reference to the Cost_Recalculation_Service
2. THE Batch_Record_Service SHALL have a reference to the Menu_Repository
3. THE Menu_Cost_Handler SHALL have a reference to the Menu_Repository
4. WHEN the application starts, THE System SHALL wire all dependencies in main.go before starting the server

### Requirement 7: Error Handling and Logging

**User Story:** As a developer, I want comprehensive error handling and logging, so that I can troubleshoot issues with cost recalculation.

#### Acceptance Criteria

1. WHEN a batch is created and recalculation is queued, THE System SHALL log the count of menu items queued
2. IF finding menu items by batch definition fails, THEN THE System SHALL log a warning with the batch definition ID and error details
3. IF queueing a recalculation job fails, THEN THE System SHALL log a warning with the menu item ID and error details
4. WHEN the recalculate-all endpoint is called, THE System SHALL log the total count of items queued
5. THE System SHALL ensure batch creation never fails due to recalculation errors
