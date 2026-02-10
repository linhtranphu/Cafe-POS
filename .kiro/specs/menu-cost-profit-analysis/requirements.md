# Requirements Document

## Introduction

Hệ thống quản lý cafe POS cần tính năng phân tích chi phí và lợi nhuận cho từng menu item. Tính năng này cho phép manager theo dõi giá vốn (cost) của mỗi món dựa trên ingredients và cost_per_unit hiện tại, từ đó đánh giá profit margin và phát hiện các món bán lỗ. Điều này giúp manager đưa ra quyết định kinh doanh chính xác về giá bán, menu optimization và quản lý chi phí.

## Glossary

- **System**: Hệ thống POS cafe bao gồm backend API và frontend manager interface
- **MenuItem**: Món trong menu với thông tin name, price, category, và danh sách ingredients
- **Ingredient**: Nguyên liệu với cost_per_unit được cập nhật theo weighted average cost, có thể có conversion_rate để quy đổi đơn vị
- **Current_Cost**: Giá vốn hiện tại của menu item, tính từ cost_per_unit hiện tại của ingredients (dùng cho pricing decisions)
- **Accounting_Cost**: Giá vốn chính thức được lưu khi kết ca, dùng cho báo cáo profit/loss
- **Cost_Status**: Trạng thái của cost data - FINAL (đã chốt ca), ESTIMATED (ca chưa đóng), INCOMPLETE (thiếu giá nguyên liệu)
- **Gross_Profit**: Lợi nhuận gộp = Revenue - Cost of Goods Sold (COGS). Đây là lợi nhuận từ bán hàng trước khi trừ chi phí vận hành
- **Operating_Profit**: Lợi nhuận vận hành = Gross Profit - Operating Expenses (lương nhân viên + mặt bằng + điện nước + marketing)
- **Profit_Margin**: Tỷ lệ lợi nhuận = (price - cost) / price * 100%
- **Absolute_Profit**: Lợi nhuận tuyệt đối = price - cost (tiền mặt thực tế)
- **Manager**: Người dùng có quyền xem báo cáo chi phí và lợi nhuận
- **Cost_Calculator**: Module tính toán cost cho menu items
- **Profit_Analyzer**: Module phân tích profit/loss theo món và category
- **Shift**: Ca làm việc với thời gian bắt đầu và kết thúc, chứa các orders

## Requirements

### Requirement 1: Current Menu Item Cost Calculation

**User Story:** As a manager, I want the system to automatically calculate the current cost of each menu item based on its ingredients, so that I can know the current cost of goods sold for pricing decisions.

#### Acceptance Criteria

1. WHEN calculating current_cost for a menu item, THE Cost_Calculator SHALL use the current cost_per_unit values of all ingredients
2. WHEN a menu item has ingredients with valid cost_per_unit values, THE Cost_Calculator SHALL compute the total current_cost by summing (ingredient.quantity * ingredient.cost_per_unit) for all ingredients
3. WHEN an ingredient's cost_per_unit is updated, THE System SHALL queue a background job to recalculate current_cost for all affected menu items asynchronously
4. WHEN a menu item has no ingredients, THE Cost_Calculator SHALL return zero current_cost with cost_status = "FINAL"
5. WHEN an ingredient in a menu item has null or undefined cost_per_unit, THE Cost_Calculator SHALL mark the menu item with cost_status = "INCOMPLETE" and SHALL display a warning "⚠ Thiếu giá nguyên liệu" in the UI
6. WHEN a menu item has cost_status = "INCOMPLETE", THE System SHALL NOT include that item in profit calculations or reports
7. THE Cost_Calculator SHALL round the final current_cost to two decimal places
8. WHEN cost recalculation is in progress, THE System SHALL display "Costs are being updated..." indicator in the manager view without blocking other operations

### Requirement 2: Profit Margin Calculation

**User Story:** As a manager, I want to see the profit margin and absolute profit for each menu item, so that I can identify which items are most profitable.

#### Acceptance Criteria

1. WHEN a menu item has both cost and price values, THE Profit_Analyzer SHALL calculate profit_margin as ((price - cost) / price) * 100
2. WHEN a menu item's price equals its cost, THE Profit_Analyzer SHALL return zero profit_margin and zero absolute_profit
3. WHEN a menu item's cost exceeds its price, THE Profit_Analyzer SHALL return a negative profit_margin and negative absolute_profit
4. WHEN a menu item's price is zero or negative (promotional/gifted items), THE Profit_Analyzer SHALL mark profit_margin as "N/A" and SHALL NOT include it in average calculations
5. THE Profit_Analyzer SHALL round profit_margin to two decimal places
6. THE Profit_Analyzer SHALL calculate absolute_profit as (price - cost) to show actual cash profit/loss
7. WHEN absolute_profit is negative, THE System SHALL display it with a red indicator showing actual cash loss
8. WHEN absolute_profit is approximately zero (within 0.01), THE System SHALL display "Hòa vốn" (break-even) indicator
9. WHEN a menu item has cost_status = "INCOMPLETE", THE Profit_Analyzer SHALL NOT calculate profit_margin and SHALL display "Data Incomplete" indicator

### Requirement 3: Loss Detection and Warning

**User Story:** As a manager, I want to be alerted when menu items are sold at a loss or have low margins, so that I can adjust pricing or recipe to maintain profitability.

#### Acceptance Criteria

1. WHEN a menu item's cost exceeds its price, THE System SHALL mark the item with "loss" status (red warning)
2. WHEN a menu item's profit_margin is below the configured low_margin_threshold (default 20%), THE System SHALL mark the item with "low_margin" status (yellow warning)
3. THE System SHALL allow managers to configure the low_margin_threshold value (e.g., 15%, 20%, 25%) per shop
4. WHEN displaying menu items with loss or low_margin warnings, THE System SHALL highlight them with distinct visual indicators (red for loss, yellow for low margin)
5. WHEN a manager views the menu cost report, THE System SHALL display the count of items with loss and low_margin warnings at the top
6. WHEN a menu item transitions between warning states, THE System SHALL update the warning status immediately

### Requirement 4: Menu Item Cost Report API

**User Story:** As a manager, I want to retrieve cost and profit information for all menu items via API, so that I can analyze profitability in the management interface.

#### Acceptance Criteria

1. THE System SHALL provide an API endpoint that returns cost, price, profit_margin, and absolute_profit for each menu item
2. WHEN the API is called, THE System SHALL include the timestamp of the last cost calculation for each item
3. THE System SHALL support filtering menu items by category in the cost report API
4. THE System SHALL support sorting menu items by profit_margin (ascending or descending) in the API response
5. WHEN an error occurs during cost calculation, THE System SHALL return an error response with a descriptive message

### Requirement 5: Shift Closure Cost Calculation

**User Story:** As a manager, I want the system to calculate and store menu item costs when I close a shift (chốt sổ), so that I can accurately track cost of goods sold based on the accounting period.

#### Acceptance Criteria

1. WHEN an order is created during a shift, THE System SHALL NOT calculate or store cost for the order items
2. WHEN a shift is closed (kết ca), THE System SHALL calculate cost for ALL orders in that shift using the cost_per_unit values at the time of shift closure
3. WHEN calculating cost during shift closure, THE Cost_Calculator SHALL use the same calculation method as current menu item cost (sum of ingredient.quantity * ingredient.cost_per_unit)
4. THE System SHALL store the calculated cost in the order_items table with a timestamp indicating when the cost was calculated
5. WHEN an order is retrieved after shift closure, THE System SHALL return the stored cost calculated at shift closure time
6. WHEN viewing reports for orders in an open (unclosed) shift, THE System SHALL display cost as "pending" or use current cost as an estimate with a clear "Estimated" indicator
7. WHEN historical cost data is not available for an order item (orders before this feature), THE System SHALL use the current menu item cost as a fallback with an "Estimated" indicator
8. WHEN an ingredient's cost_per_unit changes during an open shift, THE System SHALL NOT recalculate costs for orders already in that shift until the shift is closed

### Requirement 6: Category-Level Profit Analysis

**User Story:** As a manager, I want to see profit analysis aggregated by menu category, so that I can identify which categories are most profitable.

#### Acceptance Criteria

1. THE Profit_Analyzer SHALL calculate total_revenue, total_cost, and total_profit for each menu category based on historical order data
2. WHEN calculating category profit, THE System SHALL use the stored accounting_cost from shift closure (not current_cost)
3. THE System SHALL calculate average_profit_margin for each category as (total_profit / total_revenue) * 100
4. THE System SHALL support filtering category profit analysis by date range
5. WHEN a category has no orders in the selected date range, THE Profit_Analyzer SHALL return zero values for that category
6. WHEN viewing profit for orders in an unclosed shift, THE System SHALL clearly indicate that costs are estimates (cost_status = "ESTIMATED") and not final

### Requirement 6.5: Operating Profit Analysis

**User Story:** As a manager, I want to see operating profit after deducting operating expenses (lương, mặt bằng, điện nước, marketing), so that I can understand the true profitability of my business.

#### Acceptance Criteria

1. THE System SHALL calculate gross_profit as (total_revenue - total_cost_of_goods_sold) for a given period
2. THE System SHALL allow managers to input operating expenses for a period, including: staff_salary, rent, utilities, marketing_costs, other_expenses
3. THE System SHALL calculate operating_profit as: gross_profit - (staff_salary + rent + utilities + marketing_costs + other_expenses)
4. THE System SHALL calculate operating_profit_margin as (operating_profit / total_revenue) * 100
5. WHEN displaying operating profit analysis, THE System SHALL show breakdown of all expense categories with individual amounts
6. THE System SHALL support filtering operating profit analysis by date range (daily, weekly, monthly)
7. WHEN operating expenses are not entered for a period, THE System SHALL display gross_profit only with a note "Chưa nhập chi phí vận hành"
8. WHEN viewing daily reports but expenses are entered monthly, THE System SHALL allocate expenses proportionally (monthly_expense / days_in_month) with an indicator "Chi phí được phân bổ từ tháng"

### Requirement 7: Manager View Display

**User Story:** As a manager, I want to view cost and profit information in an intuitive interface, so that I can quickly understand the financial performance of menu items.

#### Acceptance Criteria

1. THE System SHALL display a table showing menu item name, current_cost, price, profit_margin, absolute_profit, and cost_status
2. WHEN displaying the menu cost view, THE System SHALL sort items by profit_margin in descending order by default
3. THE System SHALL provide filter controls for category selection in the manager view
4. THE System SHALL display summary statistics including total items, items at loss, items with low margin, average profit_margin
5. WHEN a manager clicks on a menu item row, THE System SHALL show detailed ingredient cost breakdown
6. THE System SHALL use color coding: green for profitable items, yellow for low margin, red for loss, gray for incomplete data

### Requirement 8: Menu Item Cost History Tracking

**User Story:** As a manager, I want to track how menu item costs change over time with visual trends, so that I can understand cost trends and make informed pricing decisions.

#### Acceptance Criteria

1. WHEN an ingredient's cost_per_unit changes, THE System SHALL record the timestamp of the change
2. THE System SHALL provide an API to retrieve historical cost data for a specific menu item over a date range
3. WHEN displaying historical cost data, THE System SHALL show cost value and calculation timestamp for each data point
4. THE System SHALL provide a line chart visualization showing cost trends over time for each menu item
5. THE System SHALL support exporting historical cost data in CSV format
6. WHEN a menu item's recipe is modified, THE System SHALL record the change with a timestamp and SHALL NOT recalculate historical costs for closed shifts

### Requirement 9: Real-Time Cost Updates

**User Story:** As a manager, I want current_cost to update when ingredient costs change, so that I can see up-to-date cost information for pricing decisions.

#### Acceptance Criteria

1. WHEN an ingredient cost_per_unit is updated via the ingredient management interface, THE System SHALL queue a background job to recalculate current_cost for all menu items using that ingredient
2. WHEN multiple ingredients are updated in a batch operation, THE System SHALL recalculate affected menu items once after all updates complete
3. THE System SHALL complete cost recalculation within 5 seconds for up to 1000 menu items using asynchronous processing with eventual consistency
4. WHEN cost recalculation is in progress, THE System SHALL display a "Costs are being updated..." indicator in the manager view without blocking user interactions
5. WHEN cost recalculation completes, THE System SHALL refresh the manager view automatically or provide a "Refresh" button
6. THE System SHALL NOT update accounting_cost for orders in closed shifts when ingredient costs change (accounting_cost is immutable after shift closure)

### Requirement 10: Ingredient Unit Conversion and Wastage

**User Story:** As a manager, I want the system to handle unit conversions and wastage factors, so that cost calculations reflect real-world usage accurately.

#### Acceptance Criteria

1. WHEN an ingredient has a conversion_rate defined, THE Cost_Calculator SHALL apply the conversion when calculating cost (e.g., stock in kg but recipe uses grams)
2. WHEN an ingredient has a wastage_percentage defined, THE Cost_Calculator SHALL increase the cost by that percentage to account for waste (e.g., 10% wastage means multiply cost by 1.10)
3. THE System SHALL allow managers to configure conversion_rate and wastage_percentage for each ingredient
4. WHEN calculating menu item cost with conversion and wastage, THE formula SHALL be: cost = (quantity * conversion_rate * cost_per_unit) * (1 + wastage_percentage/100)
5. WHEN conversion_rate or wastage_percentage is not defined, THE System SHALL treat them as 1.0 and 0% respectively

## Notes

### Key Concepts

- **Current_Cost vs Accounting_Cost**:
  - **current_cost**: Giá vốn hiện tại, tính từ cost_per_unit hiện tại của ingredients. Dùng cho pricing decisions và real-time analysis. Được update asynchronously khi ingredient costs thay đổi.
  - **accounting_cost**: Giá vốn chính thức được lưu khi kết ca (shift closure). Dùng cho báo cáo profit/loss và accounting. Immutable sau khi shift đóng.

- **Gross Profit vs Operating Profit**:
  - **Gross Profit (Lợi nhuận gộp)**: Revenue - Cost of Goods Sold (COGS). Đây là lợi nhuận từ bán hàng trước khi trừ chi phí vận hành. Được tính tự động từ orders và accounting_cost.
  - **Operating Profit (Lợi nhuận vận hành)**: Gross Profit - Operating Expenses (lương nhân viên + mặt bằng + điện nước + marketing + chi phí khác). Đây là lợi nhuận thực tế sau khi trừ tất cả chi phí vận hành. Manager cần nhập operating expenses thủ công.
  - **Expense Allocation**: Khi manager nhập chi phí theo tháng nhưng xem báo cáo theo ngày, hệ thống sẽ phân bổ tự động: `daily_expense = monthly_expense / days_in_month`. Hiển thị indicator "Chi phí được phân bổ từ tháng" để manager biết đây là ước tính.

- **Cost Status**:
  - **FINAL**: Cost đã được tính và lưu chính thức (shift đã đóng hoặc menu item không có ingredients)
  - **ESTIMATED**: Cost tạm tính từ current_cost (shift chưa đóng)
  - **INCOMPLETE**: Thiếu giá nguyên liệu, không thể tính cost chính xác

- **Shift Closure Cost**: Cost được tính và lưu khi kết ca (chốt sổ), sử dụng cost_per_unit tại thời điểm đó. Đây là cost chính thức cho accounting và profit analysis

- **Cost Timing**: Orders trong ca chưa đóng sẽ không có accounting_cost, chỉ có current_cost với status = "ESTIMATED"

- **Weighted Average Cost**: Ingredient cost_per_unit được cập nhật theo weighted average cost method khi nhập hàng

- **Fallback Strategy**: Nếu accounting_cost chưa được tính (shift chưa đóng hoặc orders cũ), sử dụng current_cost làm estimate với indicator rõ ràng

- **Accounting Accuracy**: Logic này phù hợp với quy trình kế toán cafe, cost phản ánh giá vốn tại thời điểm chốt sổ

- **Incomplete Data Handling**: Menu items thiếu cost data sẽ được đánh dấu cost_status = "INCOMPLETE" với warning rõ ràng, không tính vào reports (tránh hiểu lầm)

- **Asynchronous Processing**: Cost recalculation chạy background với eventual consistency để tránh treo hệ thống khi có nhiều menu items

- **Unit Conversion**: Hỗ trợ quy đổi đơn vị (kg → gram) và wastage factor để tính cost chính xác

- **Warning Thresholds**: 
  - Loss (cost > price) = red warning
  - Low margin (profit_margin < low_margin_threshold) = yellow warning
  - low_margin_threshold is configurable per shop (default: 20%)
  - Manager có thể điều chỉnh threshold tùy theo business model (e.g., cafe cao cấp có thể set 30%, cafe bình dân có thể set 15%)

- **Recipe Changes**: Khi sửa recipe, không được tính lại accounting_cost của các shifts đã đóng (immutable historical data)

- **Promotional Items**: Items có price = 0 hoặc âm sẽ được đánh dấu "N/A" cho profit margin

- **Absolute Profit**: Hiển thị lợi nhuận tiền mặt thực tế (price - cost), manager thường nhìn tiền trước %

- Manager view cần responsive design để hỗ trợ cả desktop và mobile
- API cần pagination cho danh sách menu items khi số lượng lớn

### Data Model Overview

**MenuItem**:
- `current_cost`: decimal - Giá vốn hiện tại (real-time)
- `cost_last_calculated_at`: timestamp - Thời điểm tính current_cost
- `cost_status`: enum - FINAL | ESTIMATED | INCOMPLETE

**OrderItem**:
- `accounting_cost`: decimal - Giá vốn chính thức (lưu khi kết ca)
- `cost_calculated_at`: timestamp - Thời điểm tính accounting_cost
- `cost_status`: enum - FINAL | ESTIMATED | INCOMPLETE

**Ingredient**:
- `cost_per_unit`: decimal - Giá vốn trên đơn vị (weighted average)
- `cost_updated_at`: timestamp - Thời điểm cập nhật cost
- `conversion_rate`: decimal - Tỷ lệ quy đổi đơn vị (optional)
- `wastage_percentage`: decimal - Tỷ lệ hao hụt % (optional)

**OperatingExpense** (NEW):
- `period_start`: date - Ngày bắt đầu kỳ (e.g., 2024-01-01)
- `period_end`: date - Ngày kết thúc kỳ (e.g., 2024-01-31)
- `staff_salary`: decimal - Tổng lương nhân viên
- `rent`: decimal - Tiền thuê mặt bằng
- `utilities`: decimal - Điện nước
- `marketing_costs`: decimal - Chi phí marketing
- `other_expenses`: decimal - Chi phí khác
- `total_expenses`: decimal - Tổng chi phí vận hành (tự động tính)
- `created_at`: timestamp
- `updated_at`: timestamp

**ShopSettings** (NEW):
- `low_margin_threshold`: decimal - Ngưỡng cảnh báo lợi nhuận thấp (default: 20.0)
- `updated_at`: timestamp

## Out of Scope (Phase 1)

Các tính năng sau sẽ được xem xét trong các phase tiếp theo:

- **Combo Items**: Menu items lồng ghép (combo gồm nhiều items con) - cần logic tính cost phức tạp hơn
- **Tax and Fees**: Tính toán VAT, phí dịch vụ, phí app (ShopeeFood, Grab) - cần tích hợp với accounting system
- **Net Profit Analysis**: Profit sau khi trừ thuế và phí - phụ thuộc vào tax/fees feature
- **Multi-location Cost**: Cost khác nhau giữa các chi nhánh - cần architecture phức tạp hơn
