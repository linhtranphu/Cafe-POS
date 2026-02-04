# Phân Tích Tính Năng Cash Handover - Bàn Giao Tiền Waiter → Cashier

## 🎯 Tổng Quan

Tính năng Cash Handover cho phép waiter bàn giao tiền thu được từ khách hàng cho cashier với **tương tác trực tiếp** giữa hai bên trong quá trình ca làm việc hoặc khi kết thúc ca.

## 📋 Yêu Cầu Chức Năng

### 1. Quy Trình Handover với Đối Soát Tiền
- **Waiter** khởi tạo yêu cầu bàn giao (một phần hoặc toàn bộ)
- **Cashier** nhận thông báo và kiểm tra thông tin
- **Đối soát vật lý**: Cashier đếm tiền thực tế và so sánh
- **Xác nhận/Từ chối**: Cashier xác nhận với số tiền thực nhận hoặc từ chối với lý do
- **Ghi nhận chênh lệch**: Nếu có sai khác giữa khai báo và thực tế
- **Cập nhật hệ thống**: Tự động cập nhật số dư cho cả hai bên
- **Audit trail**: Ghi lại toàn bộ quá trình để kiểm toán

### 2. Đối Soát Chi Tiết
- **Số tiền khai báo**: Waiter khai báo số tiền bàn giao
- **Số tiền thực nhận**: Cashier đếm và xác nhận số tiền thực tế
- **Phát hiện chênh lệch**: Hệ thống tự động tính toán sai khác
- **Xử lý chênh lệch**: Ghi nhận lý do và trách nhiệm
- **Báo cáo sai lệch**: Tạo báo cáo cho quản lý nếu cần

### 3. Tích Hợp UI với Đối Soát
- **Waiter**: Form khai báo số tiền với breakdown chi tiết
- **Cashier**: Interface đối soát với calculator và form xác nhận
- **Discrepancy Handling**: Modal xử lý chênh lệch với các tùy chọn
- **Real-time**: Updates và notifications cho cả hai bên
- **History**: Lịch sử đầy đủ với thông tin đối soát

### 4. Validation Rules với Đối Soát
- Waiter chỉ có thể handover tiền <= số tiền hiện có
- Cashier phải xác nhận số tiền thực nhận (có thể khác khai báo)
- Chênh lệch > threshold phải có lý do và approval
- Tất cả giao dịch phải có audit trail đầy đủ
- Không thể handover khi không có cashier shift mở

## 🏗️ Thiết Kế Database với Đối Soát

### Cash Handover Collection (Mở Rộng)
```javascript
{
  _id: ObjectId,
  waiter_shift_id: ObjectId,        // ID ca waiter
  cashier_shift_id: ObjectId,       // ID ca cashier
  waiter_id: ObjectId,              // ID waiter
  waiter_name: String,              // Tên waiter
  cashier_id: ObjectId,             // ID cashier
  cashier_name: String,             // Tên cashier
  
  // Thông tin bàn giao
  declared_amount: Number,          // Số tiền waiter khai báo
  actual_amount: Number,            // Số tiền cashier thực nhận
  discrepancy: Number,              // Chênh lệch (actual - declared)
  
  handover_type: String,            // "PARTIAL" | "FULL" | "END_SHIFT"
  status: String,                   // "PENDING" | "CONFIRMED" | "REJECTED" | "DISCREPANCY"
  
  // Ghi chú và lý do
  waiter_note: String,              // Ghi chú từ waiter
  cashier_note: String,             // Ghi chú từ cashier
  discrepancy_reason: String,       // Lý do chênh lệch
  discrepancy_responsibility: String, // "WAITER" | "CASHIER" | "SYSTEM" | "UNKNOWN"
  
  // Thời gian
  handover_at: Date,                // Thời gian bàn giao
  confirmed_at: Date,               // Thời gian xác nhận
  reconciled_at: Date,              // Thời gian đối soát
  
  // Metadata
  end_cash: Number,                 // Tiền cuối ca (cho END_SHIFT)
  requires_approval: Boolean,       // Cần approval từ manager
  approved_by: ObjectId,            // ID người approve
  approved_at: Date,                // Thời gian approve
  
  created_at: Date,
  updated_at: Date
}
```

### Cash Discrepancy Collection (Mới)
```javascript
{
  _id: ObjectId,
  handover_id: ObjectId,            // Liên kết với handover
  waiter_shift_id: ObjectId,
  cashier_shift_id: ObjectId,
  
  // Thông tin chênh lệch
  declared_amount: Number,
  actual_amount: Number,
  discrepancy_amount: Number,       // Số tiền chênh lệch
  discrepancy_type: String,         // "SHORTAGE" | "OVERAGE"
  
  // Phân tích nguyên nhân
  reason_category: String,          // "COUNTING_ERROR" | "TRANSACTION_ERROR" | "THEFT" | "OTHER"
  detailed_reason: String,          // Mô tả chi tiết
  responsibility: String,           // "WAITER" | "CASHIER" | "SYSTEM" | "CUSTOMER" | "UNKNOWN"
  
  // Xử lý
  resolution_status: String,        // "PENDING" | "RESOLVED" | "ESCALATED"
  resolution_action: String,        // Hành động xử lý
  resolved_by: ObjectId,            // ID người xử lý
  resolved_at: Date,                // Thời gian xử lý
  
  // Approval (nếu cần)
  requires_manager_approval: Boolean,
  manager_approved: Boolean,
  approved_by: ObjectId,
  approved_at: Date,
  manager_note: String,
  
  created_at: Date,
  updated_at: Date
}
```

### Cập Nhật Shift Models
```go
// Thêm vào Shift struct (Waiter)
type Shift struct {
    // ... existing fields
    CurrentCash         float64 `bson:"current_cash" json:"current_cash"`           // Tiền hiện có
    HandedOverCash      float64 `bson:"handed_over_cash" json:"handed_over_cash"`   // Tổng tiền đã bàn giao
    RemainingCash       float64 `bson:"remaining_cash" json:"remaining_cash"`       // Tiền còn lại
    TotalDiscrepancy    float64 `bson:"total_discrepancy" json:"total_discrepancy"` // Tổng chênh lệch
    HandoverCount       int     `bson:"handover_count" json:"handover_count"`       // Số lần bàn giao
}

// Thêm vào CashierShift struct
type CashierShift struct {
    // ... existing fields
    ReceivedCash        float64 `bson:"received_cash" json:"received_cash"`         // Tiền nhận từ waiter
    TotalDiscrepancy    float64 `bson:"total_discrepancy" json:"total_discrepancy"` // Tổng chênh lệch
    HandoverCount       int     `bson:"handover_count" json:"handover_count"`       // Số lần nhận bàn giao
    DiscrepancyCount    int     `bson:"discrepancy_count" json:"discrepancy_count"` // Số lần có chênh lệch
}
```

## 🔧 Backend Implementation

### 1. Domain Models

#### Cash Handover Domain (Mở Rộng)
```go
// backend/domain/handover/cash_handover.go
package handover

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type HandoverStatus string
type HandoverType string
type DiscrepancyType string
type ResponsibilityType string

const (
    StatusPending     HandoverStatus = "PENDING"
    StatusConfirmed   HandoverStatus = "CONFIRMED"
    StatusRejected    HandoverStatus = "REJECTED"
    StatusDiscrepancy HandoverStatus = "DISCREPANCY"  // Có chênh lệch cần xử lý
    
    TypePartial   HandoverType = "PARTIAL"
    TypeFull      HandoverType = "FULL"
    TypeEndShift  HandoverType = "END_SHIFT"
    
    DiscrepancyShortage DiscrepancyType = "SHORTAGE"  // Thiếu tiền
    DiscrepancyOverage  DiscrepancyType = "OVERAGE"   // Thừa tiền
    
    ResponsibilityWaiter   ResponsibilityType = "WAITER"
    ResponsibilityCashier  ResponsibilityType = "CASHIER"
    ResponsibilitySystem   ResponsibilityType = "SYSTEM"
    ResponsibilityCustomer ResponsibilityType = "CUSTOMER"
    ResponsibilityUnknown  ResponsibilityType = "UNKNOWN"
)

// Cash breakdown structure - REMOVED
// Calculate total from breakdown - REMOVED

type CashHandover struct {
    ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    WaiterShiftID   primitive.ObjectID `bson:"waiter_shift_id" json:"waiter_shift_id"`
    CashierShiftID  primitive.ObjectID `bson:"cashier_shift_id" json:"cashier_shift_id"`
    WaiterID        primitive.ObjectID `bson:"waiter_id" json:"waiter_id"`
    WaiterName      string             `bson:"waiter_name" json:"waiter_name"`
    CashierID       primitive.ObjectID `bson:"cashier_id" json:"cashier_id"`
    CashierName     string             `bson:"cashier_name" json:"cashier_name"`
    
    // Amounts
    DeclaredAmount  float64            `bson:"declared_amount" json:"declared_amount"`   // Waiter khai báo
    ActualAmount    float64            `bson:"actual_amount" json:"actual_amount"`       // Cashier thực nhận
    Discrepancy     float64            `bson:"discrepancy" json:"discrepancy"`           // Chênh lệch
    
    HandoverType    HandoverType       `bson:"handover_type" json:"handover_type"`
    Status          HandoverStatus     `bson:"status" json:"status"`
    
    // Notes and reasons
    WaiterNote              string             `bson:"waiter_note,omitempty" json:"waiter_note,omitempty"`
    CashierNote             string             `bson:"cashier_note,omitempty" json:"cashier_note,omitempty"`
    DiscrepancyReason       string             `bson:"discrepancy_reason,omitempty" json:"discrepancy_reason,omitempty"`
    DiscrepancyResponsibility ResponsibilityType `bson:"discrepancy_responsibility,omitempty" json:"discrepancy_responsibility,omitempty"`
    
    // Timestamps
    HandoverAt      time.Time          `bson:"handover_at" json:"handover_at"`
    ConfirmedAt     *time.Time         `bson:"confirmed_at,omitempty" json:"confirmed_at,omitempty"`
    ReconciledAt    *time.Time         `bson:"reconciled_at,omitempty" json:"reconciled_at,omitempty"`
    
    // Metadata
    EndCash         float64            `bson:"end_cash,omitempty" json:"end_cash,omitempty"`
    RequiresApproval bool              `bson:"requires_approval" json:"requires_approval"`
    ApprovedBy      primitive.ObjectID `bson:"approved_by,omitempty" json:"approved_by,omitempty"`
    ApprovedAt      *time.Time         `bson:"approved_at,omitempty" json:"approved_at,omitempty"`
    
    CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
    UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
}

// Check if handover has discrepancy
func (h *CashHandover) HasDiscrepancy() bool {
    return h.Discrepancy != 0
}

// Get discrepancy type
func (h *CashHandover) GetDiscrepancyType() DiscrepancyType {
    if h.Discrepancy < 0 {
        return DiscrepancyShortage
    } else if h.Discrepancy > 0 {
        return DiscrepancyOverage
    }
    return ""
}

// Check if requires manager approval (large discrepancy)
func (h *CashHandover) RequiresManagerApproval(threshold float64) bool {
    return h.HasDiscrepancy() && (h.Discrepancy > threshold || h.Discrepancy < -threshold)
}

// Request structures
type CreateHandoverRequest struct {
    DeclaredAmount float64      `json:"declared_amount" binding:"required,gt=0"`
    HandoverType   HandoverType `json:"handover_type" binding:"required"`
    WaiterNote     string       `json:"waiter_note"`
}

type CreateHandoverAndEndShiftRequest struct {
    DeclaredAmount float64 `json:"declared_amount" binding:"required,gt=0"`
    WaiterNote     string  `json:"waiter_note"`
    EndCash        float64 `json:"end_cash" binding:"min=0"`
}

type ConfirmHandoverRequest struct {
    ActualAmount            float64            `json:"actual_amount" binding:"required,gte=0"`
    Status                  HandoverStatus     `json:"status" binding:"required"`
    CashierNote             string             `json:"cashier_note"`
    DiscrepancyReason       string             `json:"discrepancy_reason"`
    DiscrepancyResponsibility ResponsibilityType `json:"discrepancy_responsibility"`
}

// Cash Discrepancy model
type CashDiscrepancy struct {
    ID                      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    HandoverID              primitive.ObjectID `bson:"handover_id" json:"handover_id"`
    WaiterShiftID           primitive.ObjectID `bson:"waiter_shift_id" json:"waiter_shift_id"`
    CashierShiftID          primitive.ObjectID `bson:"cashier_shift_id" json:"cashier_shift_id"`
    
    // Discrepancy details
    DeclaredAmount          float64            `bson:"declared_amount" json:"declared_amount"`
    ActualAmount            float64            `bson:"actual_amount" json:"actual_amount"`
    DiscrepancyAmount       float64            `bson:"discrepancy_amount" json:"discrepancy_amount"`
    DiscrepancyType         DiscrepancyType    `bson:"discrepancy_type" json:"discrepancy_type"`
    
    // Analysis
    ReasonCategory          string             `bson:"reason_category" json:"reason_category"`
    DetailedReason          string             `bson:"detailed_reason" json:"detailed_reason"`
    Responsibility          ResponsibilityType `bson:"responsibility" json:"responsibility"`
    
    // Resolution
    ResolutionStatus        string             `bson:"resolution_status" json:"resolution_status"`
    ResolutionAction        string             `bson:"resolution_action" json:"resolution_action"`
    ResolvedBy              primitive.ObjectID `bson:"resolved_by,omitempty" json:"resolved_by,omitempty"`
    ResolvedAt              *time.Time         `bson:"resolved_at,omitempty" json:"resolved_at,omitempty"`
    
    // Manager approval
    RequiresManagerApproval bool               `bson:"requires_manager_approval" json:"requires_manager_approval"`
    ManagerApproved         bool               `bson:"manager_approved" json:"manager_approved"`
    ApprovedBy              primitive.ObjectID `bson:"approved_by,omitempty" json:"approved_by,omitempty"`
    ApprovedAt              *time.Time         `bson:"approved_at,omitempty" json:"approved_at,omitempty"`
    ManagerNote             string             `bson:"manager_note,omitempty" json:"manager_note,omitempty"`
    
    CreatedAt               time.Time          `bson:"created_at" json:"created_at"`
    UpdatedAt               time.Time          `bson:"updated_at" json:"updated_at"`
}
```

### 3. Service Layer với Đối Soát
```go
// backend/application/services/cash_handover_service.go
type CashHandoverService struct {
    handoverRepo        CashHandoverRepository
    discrepancyRepo     CashDiscrepancyRepository
    shiftRepo           ShiftRepository
    cashierShiftRepo    CashierShiftRepository
    stateMachineManager *domain.StateMachineManager
    discrepancyThreshold float64  // Ngưỡng chênh lệch cần approval
}

func (s *CashHandoverService) CreateHandover(ctx context.Context, waiterShiftID primitive.ObjectID, req *handover.CreateHandoverRequest, waiterID, waiterName string) (*handover.CashHandover, error) {
    // 1. Validate waiter shift exists and is open
    waiterShift, err := s.shiftRepo.FindByID(ctx, waiterShiftID)
    if err != nil || waiterShift.Status != order.ShiftOpen {
        return nil, errors.New("waiter shift not found or not open")
    }
    
    // 2. Check if waiter owns the shift
    waiterOID, _ := primitive.ObjectIDFromHex(waiterID)
    if waiterShift.UserID != waiterOID {
        return nil, errors.New("unauthorized: not your shift")
    }
    
    // 3. Validate declared amount
    if req.DeclaredAmount > waiterShift.RemainingCash {
        return nil, errors.New("declared amount exceeds remaining cash")
    }
    
    // 4. Validate cash breakdown if provided
    if req.CashBreakdown != nil {
        breakdownTotal := req.CashBreakdown.Total()
        if breakdownTotal != req.DeclaredAmount {
            return nil, errors.New("cash breakdown total does not match declared amount")
        }
    }
    
    // 5. Find active cashier shift
    cashierShift, err := s.cashierShiftRepo.FindOpenShift(ctx)
    if err != nil {
        return nil, errors.New("no active cashier shift found")
    }
    
    // 6. Create handover record
    handover := &handover.CashHandover{
        WaiterShiftID:   waiterShiftID,
        CashierShiftID:  cashierShift.ID,
        WaiterID:        waiterOID,
        WaiterName:      waiterName,
        CashierID:       cashierShift.CashierID,
        CashierName:     cashierShift.CashierName,
        DeclaredAmount:  req.DeclaredAmount,
        ActualAmount:    0, // Will be set by cashier
        Discrepancy:     0, // Will be calculated
        HandoverType:    req.HandoverType,
        Status:          handover.StatusPending,
        WaiterNote:      req.WaiterNote,
        HandoverAt:      time.Now(),
        CreatedAt:       time.Now(),
        UpdatedAt:       time.Now(),
    }
    
    if err := s.handoverRepo.Create(ctx, handover); err != nil {
        return nil, err
    }
    
    return handover, nil
}

func (s *CashHandoverService) ConfirmHandoverWithReconciliation(ctx context.Context, handoverID primitive.ObjectID, req *handover.ConfirmHandoverRequest, cashierID string) error {
    // 1. Get handover record
    handover, err := s.handoverRepo.FindByID(ctx, handoverID)
    if err != nil {
        return err
    }
    
    // 2. Validate cashier authorization
    cashierOID, _ := primitive.ObjectIDFromHex(cashierID)
    if handover.CashierID != cashierOID {
        return errors.New("unauthorized: not assigned to you")
    }
    
    // 3. Calculate discrepancy
    discrepancy := req.ActualAmount - handover.DeclaredAmount
    
    // 4. Update handover with reconciliation data
    now := time.Now()
    handover.ActualAmount = req.ActualAmount
    handover.Discrepancy = discrepancy
    handover.Status = req.Status
    handover.CashierNote = req.CashierNote
    handover.ConfirmedAt = &now
    handover.ReconciledAt = &now
    handover.UpdatedAt = now
    
    // 5. Handle discrepancy if exists
    if handover.HasDiscrepancy() {
        handover.DiscrepancyReason = req.DiscrepancyReason
        handover.DiscrepancyResponsibility = req.DiscrepancyResponsibility
        
        // Check if requires manager approval
        if handover.RequiresManagerApproval(s.discrepancyThreshold) {
            handover.RequiresApproval = true
            handover.Status = handover.StatusDiscrepancy
        }
        
        // Create discrepancy record
        if err := s.createDiscrepancyRecord(ctx, handover); err != nil {
            return err
        }
    }
    
    // 6. Update handover record
    if err := s.handoverRepo.Update(ctx, handoverID, handover); err != nil {
        return err
    }
    
    // 7. If confirmed (and not requiring approval), update cash amounts
    if req.Status == handover.StatusConfirmed && !handover.RequiresApproval {
        if err := s.updateCashAmounts(ctx, handover); err != nil {
            return err
        }
    }
    
    return nil
}

func (s *CashHandoverService) createDiscrepancyRecord(ctx context.Context, handover *handover.CashHandover) error {
    discrepancy := &handover.CashDiscrepancy{
        HandoverID:              handover.ID,
        WaiterShiftID:           handover.WaiterShiftID,
        CashierShiftID:          handover.CashierShiftID,
        DeclaredAmount:          handover.DeclaredAmount,
        ActualAmount:            handover.ActualAmount,
        DiscrepancyAmount:       handover.Discrepancy,
        DiscrepancyType:         handover.GetDiscrepancyType(),
        DetailedReason:          handover.DiscrepancyReason,
        Responsibility:          handover.DiscrepancyResponsibility,
        ResolutionStatus:        "PENDING",
        RequiresManagerApproval: handover.RequiresApproval,
        CreatedAt:               time.Now(),
        UpdatedAt:               time.Now(),
    }
    
    return s.discrepancyRepo.Create(ctx, discrepancy)
}

func (s *CashHandoverService) updateCashAmounts(ctx context.Context, handover *handover.CashHandover) error {
    now := time.Now()
    
    // Update waiter shift - use actual amount received
    waiterShift, _ := s.shiftRepo.FindByID(ctx, handover.WaiterShiftID)
    waiterShift.HandedOverCash += handover.ActualAmount
    waiterShift.RemainingCash -= handover.DeclaredAmount  // Reduce by declared amount
    waiterShift.TotalDiscrepancy += handover.Discrepancy
    waiterShift.HandoverCount++
    waiterShift.UpdatedAt = now
    
    // Handle END_SHIFT type
    if handover.HandoverType == handover.TypeEndShift {
        // Calculate total revenue and orders
        orders, _ := s.orderRepo.FindByShiftID(ctx, handover.WaiterShiftID)
        totalRevenue := 0.0
        for _, o := range orders {
            if o.Status == order.StatusPaid || o.Status == order.StatusInProgress || o.Status == order.StatusServed {
                totalRevenue += o.Total
            }
        }
        
        // End the shift
        waiterShift.Status = order.ShiftClosed
        waiterShift.EndCash = handover.EndCash
        waiterShift.TotalRevenue = totalRevenue
        waiterShift.TotalOrders = len(orders)
        waiterShift.EndedAt = &now
        
        // Lock completed orders
        for _, o := range orders {
            if o.Status == order.StatusServed || o.Status == order.StatusCancelled {
                o.Status = order.StatusLocked
                o.LockedAt = &now
                s.orderRepo.Update(ctx, o.ID, o)
            }
        }
    }
    
    s.shiftRepo.Update(ctx, handover.WaiterShiftID, waiterShift)
    
    // Update cashier shift
    cashierShift, _ := s.cashierShiftRepo.FindByID(ctx, handover.CashierShiftID)
    cashierShift.ReceivedCash += handover.ActualAmount
    cashierShift.TotalDiscrepancy += handover.Discrepancy
    cashierShift.HandoverCount++
    if handover.HasDiscrepancy() {
        cashierShift.DiscrepancyCount++
    }
    cashierShift.UpdatedAt = now
    
    s.cashierShiftRepo.Update(ctx, handover.CashierShiftID, cashierShift)
    
    return nil
}

// Manager approval for large discrepancies
func (s *CashHandoverService) ApproveDiscrepancy(ctx context.Context, handoverID primitive.ObjectID, managerID string, approved bool, note string) error {
    handover, err := s.handoverRepo.FindByID(ctx, handoverID)
    if err != nil {
        return err
    }
    
    if !handover.RequiresApproval {
        return errors.New("handover does not require approval")
    }
    
    now := time.Now()
    managerOID, _ := primitive.ObjectIDFromHex(managerID)
    
    handover.ApprovedBy = managerOID
    handover.ApprovedAt = &now
    handover.UpdatedAt = now
    
    if approved {
        handover.Status = handover.StatusConfirmed
        // Update cash amounts after approval
        if err := s.updateCashAmounts(ctx, handover); err != nil {
            return err
        }
    } else {
        handover.Status = handover.StatusRejected
        handover.CashierNote += " | Manager rejected: " + note
    }
    
    // Update discrepancy record
    discrepancy, _ := s.discrepancyRepo.FindByHandoverID(ctx, handoverID)
    if discrepancy != nil {
        discrepancy.ManagerApproved = approved
        discrepancy.ApprovedBy = managerOID
        discrepancy.ApprovedAt = &now
        discrepancy.ManagerNote = note
        discrepancy.ResolutionStatus = "RESOLVED"
        discrepancy.UpdatedAt = now
        s.discrepancyRepo.Update(ctx, discrepancy.ID, discrepancy)
    }
    
    return s.handoverRepo.Update(ctx, handoverID, handover)
}

// Get discrepancy statistics
func (s *CashHandoverService) GetDiscrepancyStats(ctx context.Context, startDate, endDate time.Time) (*DiscrepancyStats, error) {
    handovers, err := s.handoverRepo.FindByDateRange(ctx, startDate, endDate)
    if err != nil {
        return nil, err
    }
    
    stats := &DiscrepancyStats{
        TotalHandovers:    len(handovers),
        TotalDiscrepancy:  0,
        ShortageCount:     0,
        OverageCount:      0,
        ShortageAmount:    0,
        OverageAmount:     0,
        RequiredApproval:  0,
    }
    
    for _, h := range handovers {
        if h.HasDiscrepancy() {
            stats.TotalDiscrepancy += h.Discrepancy
            if h.Discrepancy < 0 {
                stats.ShortageCount++
                stats.ShortageAmount += -h.Discrepancy
            } else {
                stats.OverageCount++
                stats.OverageAmount += h.Discrepancy
            }
            if h.RequiresApproval {
                stats.RequiredApproval++
            }
        }
    }
    
    return stats, nil
}

type DiscrepancyStats struct {
    TotalHandovers   int     `json:"total_handovers"`
    TotalDiscrepancy float64 `json:"total_discrepancy"`
    ShortageCount    int     `json:"shortage_count"`
    OverageCount     int     `json:"overage_count"`
    ShortageAmount   float64 `json:"shortage_amount"`
    OverageAmount    float64 `json:"overage_amount"`
    RequiredApproval int     `json:"required_approval"`
}
```
### 2. Repository Layer
```go
// backend/infrastructure/mongodb/cash_handover_repository.go
type CashHandoverRepository struct {
    collection *mongo.Collection
}

func (r *CashHandoverRepository) Create(ctx context.Context, handover *handover.CashHandover) error
func (r *CashHandoverRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*handover.CashHandover, error)
func (r *CashHandoverRepository) Update(ctx context.Context, id primitive.ObjectID, handover *handover.CashHandover) error
func (r *CashHandoverRepository) FindByWaiterShift(ctx context.Context, shiftID primitive.ObjectID) ([]*handover.CashHandover, error)
func (r *CashHandoverRepository) FindByCashierShift(ctx context.Context, shiftID primitive.ObjectID) ([]*handover.CashHandover, error)
func (r *CashHandoverRepository) FindPendingByCashier(ctx context.Context, cashierID primitive.ObjectID) ([]*handover.CashHandover, error)
func (r *CashHandoverRepository) FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*handover.CashHandover, error)
func (r *CashHandoverRepository) FindWithDiscrepancies(ctx context.Context) ([]*handover.CashHandover, error)
func (r *CashHandoverRepository) FindRequiringApproval(ctx context.Context) ([]*handover.CashHandover, error)

// backend/infrastructure/mongodb/cash_discrepancy_repository.go
type CashDiscrepancyRepository struct {
    collection *mongo.Collection
}

func (r *CashDiscrepancyRepository) Create(ctx context.Context, discrepancy *handover.CashDiscrepancy) error
func (r *CashDiscrepancyRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*handover.CashDiscrepancy, error)
func (r *CashDiscrepancyRepository) Update(ctx context.Context, id primitive.ObjectID, discrepancy *handover.CashDiscrepancy) error
func (r *CashDiscrepancyRepository) FindByHandoverID(ctx context.Context, handoverID primitive.ObjectID) (*handover.CashDiscrepancy, error)
func (r *CashDiscrepancyRepository) FindPendingResolution(ctx context.Context) ([]*handover.CashDiscrepancy, error)
func (r *CashDiscrepancyRepository) FindRequiringApproval(ctx context.Context) ([]*handover.CashDiscrepancy, error)
func (r *CashDiscrepancyRepository) FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*handover.CashDiscrepancy, error)
```

### 3. Service Layer
```go
// backend/application/services/cash_handover_service.go
type CashHandoverService struct {
    handoverRepo        CashHandoverRepository
    shiftRepo           ShiftRepository
    cashierShiftRepo    CashierShiftRepository
    stateMachineManager *domain.StateMachineManager
}

func (s *CashHandoverService) CreateHandover(ctx context.Context, waiterShiftID primitive.ObjectID, req *handover.CreateHandoverRequest, waiterID, waiterName string) (*handover.CashHandover, error) {
    // 1. Validate waiter shift exists and is open
    waiterShift, err := s.shiftRepo.FindByID(ctx, waiterShiftID)
    if err != nil || waiterShift.Status != order.ShiftOpen {
        return nil, errors.New("waiter shift not found or not open")
    }
    
    // 2. Check if waiter owns the shift
    waiterOID, _ := primitive.ObjectIDFromHex(waiterID)
    if waiterShift.UserID != waiterOID {
        return nil, errors.New("unauthorized: not your shift")
    }
    
    // 3. Find active cashier shift
    cashierShift, err := s.cashierShiftRepo.FindOpenShift(ctx)
    if err != nil {
        return nil, errors.New("no active cashier shift found")
    }
    
    // 4. Validate amount <= remaining cash
    if req.Amount > waiterShift.RemainingCash {
        return nil, errors.New("amount exceeds remaining cash")
    }
    
    // 5. Create handover record
    handover := &handover.CashHandover{
        WaiterShiftID:  waiterShiftID,
        CashierShiftID: cashierShift.ID,
        WaiterID:       waiterOID,
        WaiterName:     waiterName,
        CashierID:      cashierShift.CashierID,
        CashierName:    cashierShift.CashierName,
        Amount:         req.Amount,
        HandoverType:   req.HandoverType,
        Status:         handover.StatusPending,
        WaiterNote:     req.WaiterNote,
        HandoverAt:     time.Now(),
        CreatedAt:      time.Now(),
        UpdatedAt:      time.Now(),
    }
    
    if err := s.handoverRepo.Create(ctx, handover); err != nil {
        return nil, err
    }
    
    return handover, nil
}

func (s *CashHandoverService) ConfirmHandover(ctx context.Context, handoverID primitive.ObjectID, req *handover.ConfirmHandoverRequest, cashierID string) error {
    // 1. Get handover record
    handover, err := s.handoverRepo.FindByID(ctx, handoverID)
    if err != nil {
        return err
    }
    
    // 2. Validate cashier authorization
    cashierOID, _ := primitive.ObjectIDFromHex(cashierID)
    if handover.CashierID != cashierOID {
        return errors.New("unauthorized: not assigned to you")
    }
    
    // 3. Update handover status
    now := time.Now()
    handover.Status = req.Status
    handover.CashierNote = req.CashierNote
    handover.ConfirmedAt = &now
    handover.UpdatedAt = now
    
    if err := s.handoverRepo.Update(ctx, handoverID, handover); err != nil {
        return err
    }
    
    // 4. If confirmed, update waiter shift cash amounts
    if req.Status == handover.StatusConfirmed {
        waiterShift, _ := s.shiftRepo.FindByID(ctx, handover.WaiterShiftID)
        waiterShift.HandedOverCash += handover.Amount
        waiterShift.RemainingCash -= handover.Amount
        waiterShift.UpdatedAt = now
        
        s.shiftRepo.Update(ctx, handover.WaiterShiftID, waiterShift)
        
        // Update cashier shift received cash
        cashierShift, _ := s.cashierShiftRepo.FindByID(ctx, handover.CashierShiftID)
        cashierShift.ReceivedCash += handover.Amount
        cashierShift.UpdatedAt = now
        
        s.cashierShiftRepo.Update(ctx, handover.CashierShiftID, cashierShift)
    }
    
    return nil
}

func (s *CashHandoverService) CreateHandoverAndEndShift(ctx context.Context, waiterShiftID primitive.ObjectID, req *handover.CreateHandoverAndEndShiftRequest, waiterID, waiterName string) (*handover.CashHandover, error) {
    // 1. Validate waiter shift exists and is open
    waiterShift, err := s.shiftRepo.FindByID(ctx, waiterShiftID)
    if err != nil || waiterShift.Status != order.ShiftOpen {
        return nil, errors.New("waiter shift not found or not open")
    }
    
    // 2. Check if waiter owns the shift
    waiterOID, _ := primitive.ObjectIDFromHex(waiterID)
    if waiterShift.UserID != waiterOID {
        return nil, errors.New("unauthorized: not your shift")
    }
    
    // 3. Find active cashier shift
    cashierShift, err := s.cashierShiftRepo.FindOpenShift(ctx)
    if err != nil {
        return nil, errors.New("no active cashier shift found")
    }
    
    // 4. Amount must equal remaining cash for END_SHIFT
    handoverAmount := waiterShift.RemainingCash
    if handoverAmount <= 0 {
        return nil, errors.New("no remaining cash to handover")
    }
    
    // 5. Create handover record
    handover := &handover.CashHandover{
        WaiterShiftID:  waiterShiftID,
        CashierShiftID: cashierShift.ID,
        WaiterID:       waiterOID,
        WaiterName:     waiterName,
        CashierID:      cashierShift.CashierID,
        CashierName:    cashierShift.CashierName,
        Amount:         handoverAmount,
        HandoverType:   handover.TypeEndShift,
        Status:         handover.StatusPending,
        WaiterNote:     req.WaiterNote,
        EndCash:        req.EndCash,  // Store end cash for later use
        HandoverAt:     time.Now(),
        CreatedAt:      time.Now(),
        UpdatedAt:      time.Now(),
    }
    
    if err := s.handoverRepo.Create(ctx, handover); err != nil {
        return nil, err
    }
    
    return handover, nil
}

func (s *CashHandoverService) ConfirmHandoverAndEndShift(ctx context.Context, handoverID primitive.ObjectID, req *handover.ConfirmHandoverRequest, cashierID string) error {
    // 1. Get handover record
    handover, err := s.handoverRepo.FindByID(ctx, handoverID)
    if err != nil {
        return err
    }
    
    // 2. Validate cashier authorization
    cashierOID, _ := primitive.ObjectIDFromHex(cashierID)
    if handover.CashierID != cashierOID {
        return errors.New("unauthorized: not assigned to you")
    }
    
    // 3. Update handover status
    now := time.Now()
    handover.Status = req.Status
    handover.CashierNote = req.CashierNote
    handover.ConfirmedAt = &now
    handover.UpdatedAt = now
    
    if err := s.handoverRepo.Update(ctx, handoverID, handover); err != nil {
        return err
    }
    
    // 4. If confirmed and END_SHIFT type, update waiter shift and end it
    if req.Status == handover.StatusConfirmed && handover.HandoverType == handover.TypeEndShift {
        // Update waiter shift cash amounts
        waiterShift, _ := s.shiftRepo.FindByID(ctx, handover.WaiterShiftID)
        waiterShift.HandedOverCash += handover.Amount
        waiterShift.RemainingCash = 0  // All cash handed over
        waiterShift.EndCash = handover.EndCash
        waiterShift.UpdatedAt = now
        
        // Calculate total revenue and orders
        orders, _ := s.orderRepo.FindByShiftID(ctx, handover.WaiterShiftID)
        totalRevenue := 0.0
        for _, o := range orders {
            if o.Status == order.StatusPaid || o.Status == order.StatusInProgress || o.Status == order.StatusServed {
                totalRevenue += o.Total
            }
        }
        
        // End the shift
        waiterShift.Status = order.ShiftClosed
        waiterShift.TotalRevenue = totalRevenue
        waiterShift.TotalOrders = len(orders)
        waiterShift.EndedAt = &now
        
        s.shiftRepo.Update(ctx, handover.WaiterShiftID, waiterShift)
        
        // Lock completed orders
        for _, o := range orders {
            if o.Status == order.StatusServed || o.Status == order.StatusCancelled {
                o.Status = order.StatusLocked
                o.LockedAt = &now
                s.orderRepo.Update(ctx, o.ID, o)
            }
        }
        
        // Update cashier shift received cash
        cashierShift, _ := s.cashierShiftRepo.FindByID(ctx, handover.CashierShiftID)
        cashierShift.ReceivedCash += handover.Amount
        cashierShift.UpdatedAt = now
        
        s.cashierShiftRepo.Update(ctx, handover.CashierShiftID, cashierShift)
    }
    
    return nil
}
```

### 4. HTTP Handlers
```go
// backend/interfaces/http/cash_handover_handler.go
type CashHandoverHandler struct {
    handoverService *services.CashHandoverService
}

func (h *CashHandoverHandler) CreateHandover(c *gin.Context) {
    shiftID := c.Param("shift_id")
    shiftOID, _ := primitive.ObjectIDFromHex(shiftID)
    
    var req handover.CreateHandoverRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    userID, _ := c.Get("user_id")
    username, _ := c.Get("username")
    
    handover, err := h.handoverService.CreateHandover(
        c.Request.Context(), 
        shiftOID, 
        &req, 
        userID.(string), 
        username.(string),
    )
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusCreated, handover)
}

func (h *CashHandoverHandler) ConfirmHandover(c *gin.Context) {
    handoverID := c.Param("id")
    handoverOID, _ := primitive.ObjectIDFromHex(handoverID)
    
    var req handover.ConfirmHandoverRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    userID, _ := c.Get("user_id")
    
    err := h.handoverService.ConfirmHandover(
        c.Request.Context(),
        handoverOID,
        &req,
        userID.(string),
    )
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"message": "handover confirmed"})
}

func (h *CashHandoverHandler) GetPendingHandovers(c *gin.Context) {
    userID, _ := c.Get("user_id")
    userOID, _ := primitive.ObjectIDFromHex(userID.(string))
    
    handovers, err := h.handoverService.GetPendingByCashier(c.Request.Context(), userOID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, handovers)
}
```

## 🎨 Frontend Implementation

### 1. Cập Nhật ShiftView.vue (Waiter Interface)

#### Template Updates
```vue
<!-- Thêm vào phần "Ca đang mở" -->
<div v-if="currentShift" class="bg-gradient-to-r from-blue-500 to-purple-500 text-white rounded-2xl p-6 mb-4 shadow-lg">
  <!-- ... existing content ... -->
  
  <!-- Cash Status for Waiter -->
  <div v-if="isWaiter" class="grid grid-cols-3 gap-3 mb-4">
    <div class="bg-white bg-opacity-20 rounded-xl p-3">
      <p class="text-sm text-blue-100">Tiền hiện có</p>
      <p class="font-bold">{{ formatPrice(currentShift.remaining_cash || currentShift.current_cash) }}</p>
    </div>
    <div class="bg-white bg-opacity-20 rounded-xl p-3">
      <p class="text-sm text-blue-100">Đã bàn giao</p>
      <p class="font-bold">{{ formatPrice(currentShift.handed_over_cash || 0) }}</p>
    </div>
    <div class="bg-white bg-opacity-20 rounded-xl p-3">
      <p class="text-sm text-blue-100">Tổng thu</p>
      <p class="font-bold">{{ formatPrice(currentShift.total_collected || 0) }}</p>
    </div>
  </div>

  <!-- Pending Handover Status -->
  <div v-if="isWaiter && pendingHandover" class="bg-yellow-500 bg-opacity-20 rounded-xl p-3 mb-4">
    <div class="flex items-center justify-between">
      <div>
        <p class="text-sm text-yellow-100">🕐 Đang chờ xác nhận bàn giao</p>
        <p class="font-bold">{{ formatPrice(pendingHandover.amount) }}</p>
        <p class="text-xs text-yellow-200">{{ pendingHandover.handover_type === 'END_SHIFT' ? 'Bàn giao và đóng ca' : 'Bàn giao một phần' }}</p>
      </div>
      <button @click="cancelHandover(pendingHandover.id)" 
        class="bg-red-500 hover:bg-red-600 text-white px-3 py-1 rounded-lg text-sm">
        Hủy
      </button>
    </div>
  </div>

  <!-- Action Buttons for Waiter -->
  <div v-if="isWaiter" class="space-y-2">
    <!-- Partial Handover Button -->
    <button v-if="(currentShift.remaining_cash || 0) > 0 && !pendingHandover" 
      @click="showPartialHandoverForm = true"
      class="w-full bg-yellow-500 hover:bg-yellow-600 text-white px-4 py-3 rounded-xl font-bold active:scale-95 transition-transform">
      💰 Bàn giao một phần
    </button>
    
    <!-- Handover and End Shift Button -->
    <button v-if="(currentShift.remaining_cash || 0) > 0 && !pendingHandover"
      @click="showHandoverEndShiftForm = true"
      class="w-full bg-orange-500 hover:bg-orange-600 text-white px-4 py-3 rounded-xl font-bold active:scale-95 transition-transform">
      🏁 Bàn giao và đóng ca
    </button>
    
    <!-- Regular End Shift Button (only when no remaining cash) -->
    <button v-if="(currentShift.remaining_cash || 0) === 0 && !pendingHandover"
      @click="showEndShiftForm = true" 
      class="w-full bg-white text-blue-600 hover:bg-blue-50 px-4 py-3 rounded-xl font-bold active:scale-95 transition-transform">
      Kết thúc ca
    </button>
    
    <!-- Disabled state when pending -->
    <div v-if="pendingHandover" class="w-full bg-gray-400 text-gray-200 px-4 py-3 rounded-xl font-bold text-center">
      Chờ cashier xác nhận...
    </div>
  </div>
  
  <!-- Action Buttons for Non-Waiter -->
  <div v-else>
    <button @click="showEndShiftForm = true" 
      class="w-full bg-white text-blue-600 hover:bg-blue-50 px-4 py-3 rounded-xl font-bold active:scale-95 transition-transform">
      Kết thúc ca
    </button>
  </div>
</div>

<!-- Handover History Section for Waiter -->
<div v-if="isWaiter && handoverHistory.length > 0" class="bg-white rounded-2xl p-6 shadow-sm mb-4">
  <h3 class="text-xl font-bold mb-4">📋 Lịch sử bàn giao</h3>
  <div class="space-y-3">
    <div v-for="handover in handoverHistory" :key="handover.id" 
      class="border rounded-xl p-4">
      <div class="flex justify-between items-start mb-2">
        <div>
          <p class="font-bold">{{ formatPrice(handover.amount) }}</p>
          <p class="text-sm text-gray-500">{{ formatDate(handover.handover_at) }}</p>
          <p class="text-xs text-blue-600">{{ getHandoverTypeText(handover.handover_type) }}</p>
        </div>
        <span :class="getHandoverStatusClass(handover.status)"
          class="px-3 py-1 rounded-full text-xs font-medium">
          {{ getHandoverStatusText(handover.status) }}
        </span>
      </div>
      <div v-if="handover.waiter_note" class="text-sm text-gray-600 mb-2">
        <strong>Ghi chú:</strong> {{ handover.waiter_note }}
      </div>
      <div v-if="handover.cashier_note" class="text-sm text-green-600">
        <strong>Phản hồi cashier:</strong> {{ handover.cashier_note }}
      </div>
    </div>
  </div>
</div>
```

#### Handover Modals
```vue
<!-- Partial Handover Modal -->
<transition name="slide-up">
  <div v-if="showPartialHandoverForm" class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-end">
    <div class="bg-white rounded-t-3xl w-full p-6">
      <h3 class="text-xl font-bold mb-4">💰 Bàn giao một phần tiền</h3>
      
      <!-- Current Cash Info -->
      <div class="bg-blue-50 p-4 rounded-xl mb-4">
        <div class="flex justify-between items-center">
          <span class="text-sm text-gray-600">Tiền hiện có</span>
          <span class="font-bold text-2xl text-blue-600">{{ formatPrice(currentShift?.remaining_cash || 0) }}</span>
        </div>
      </div>
      
      <form @submit.prevent="createPartialHandover" class="space-y-4">
        <!-- Amount Input -->
        <div>
          <label class="block text-sm font-medium mb-2">Số tiền bàn giao (VNĐ) *</label>
          <input v-model.number="partialHandoverForm.amount" 
            type="number" 
            :max="currentShift?.remaining_cash || 0"
            min="1000" 
            step="1000" 
            required 
            class="w-full p-3 border rounded-xl text-lg font-bold focus:ring-2 focus:ring-yellow-500">
        </div>
        
        <!-- Note -->
        <div>
          <label class="block text-sm font-medium mb-2">Ghi chú (tùy chọn)</label>
          <textarea v-model="partialHandoverForm.waiter_note" 
            rows="3" 
            class="w-full p-3 border rounded-xl focus:ring-2 focus:ring-yellow-500"
            placeholder="Ghi chú về việc bàn giao..."></textarea>
        </div>
        
        <!-- Action Buttons -->
        <div class="flex gap-2">
          <button type="button" @click="showPartialHandoverForm = false" 
            class="flex-1 bg-gray-200 text-gray-700 px-4 py-3 rounded-xl font-medium">
            Hủy
          </button>
          <button type="submit" 
            class="flex-1 bg-yellow-500 hover:bg-yellow-600 text-white px-4 py-3 rounded-xl font-medium">
            Bàn giao
          </button>
        </div>
      </form>
    </div>
  </div>
</transition>

<!-- Handover and End Shift Modal -->
<transition name="slide-up">
  <div v-if="showHandoverEndShiftForm" class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-end">
    <div class="bg-white rounded-t-3xl w-full p-6">
      <h3 class="text-xl font-bold mb-4">🏁 Bàn giao toàn bộ và đóng ca</h3>
      
      <!-- Warning Notice -->
      <div class="bg-orange-50 border-l-4 border-orange-400 p-4 mb-4">
        <div class="flex">
          <div class="flex-shrink-0">
            <svg class="h-5 w-5 text-orange-400" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
            </svg>
          </div>
          <div class="ml-3">
            <p class="text-sm text-orange-700">
              <strong>Lưu ý:</strong> Thao tác này sẽ bàn giao toàn bộ tiền còn lại và tự động đóng ca sau khi cashier xác nhận.
            </p>
          </div>
        </div>
      </div>
      
      <!-- Cash Summary -->
      <div class="bg-orange-50 p-4 rounded-xl mb-4">
        <div class="space-y-2">
          <div class="flex justify-between items-center">
            <span class="text-sm text-gray-600">Tiền sẽ bàn giao</span>
            <span class="font-bold text-2xl text-orange-600">{{ formatPrice(currentShift?.remaining_cash || 0) }}</span>
          </div>
          <div class="flex justify-between items-center text-sm">
            <span class="text-gray-500">Tiền cuối ca</span>
            <span class="font-medium">{{ formatPrice(handoverEndShiftForm.end_cash) }}</span>
          </div>
        </div>
      </div>
      
      <form @submit.prevent="createHandoverAndEndShift" class="space-y-4">
        <!-- End Cash Input -->
        <div>
          <label class="block text-sm font-medium mb-2">Tiền cuối ca (VNĐ) *</label>
          <input v-model.number="handoverEndShiftForm.end_cash" 
            type="number" 
            min="0" 
            step="1000" 
            required 
            class="w-full p-3 border rounded-xl text-lg font-bold focus:ring-2 focus:ring-orange-500">
          <p class="text-xs text-gray-500 mt-1">Tiền còn lại sau khi bàn giao (thường là 0)</p>
        </div>
        
        <!-- Note -->
        <div>
          <label class="block text-sm font-medium mb-2">Ghi chú (tùy chọn)</label>
          <textarea v-model="handoverEndShiftForm.waiter_note" 
            rows="3" 
            class="w-full p-3 border rounded-xl focus:ring-2 focus:ring-orange-500"
            placeholder="Ghi chú về việc bàn giao và đóng ca..."></textarea>
        </div>
        
        <!-- Action Buttons -->
        <div class="flex gap-2">
          <button type="button" @click="showHandoverEndShiftForm = false" 
            class="flex-1 bg-gray-200 text-gray-700 px-4 py-3 rounded-xl font-medium">
            Hủy
          </button>
          <button type="submit" 
            class="flex-1 bg-orange-500 hover:bg-orange-600 text-white px-4 py-3 rounded-xl font-medium">
            Bàn giao và đóng ca
          </button>
        </div>
      </form>
    </div>
  </div>
</transition>
```

### 2. Tạo CashierHandoverView.vue (Cashier Interface)

```vue
<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Mobile Header -->
    <div class="sticky top-0 z-40 bg-white shadow-sm">
      <div class="px-4 py-3">
        <h1 class="text-xl font-bold text-gray-800">💰 Quản lý bàn giao</h1>
      </div>
    </div>

    <!-- Content -->
    <div class="px-4 py-4 pb-24">
      <!-- Pending Handovers -->
      <div class="bg-white rounded-2xl p-6 shadow-sm mb-4">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-xl font-bold">🕐 Chờ xác nhận</h3>
          <span class="bg-red-100 text-red-800 px-3 py-1 rounded-full text-sm font-medium">
            {{ pendingHandovers.length }}
          </span>
        </div>
        
        <div v-if="loading" class="text-center py-10">
          <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500"></div>
        </div>
        
        <div v-else-if="pendingHandovers.length === 0" class="text-center py-10">
          <div class="text-4xl mb-2">✅</div>
          <p class="text-gray-500">Không có yêu cầu bàn giao nào</p>
        </div>
        
        <div v-else class="space-y-4">
          <div v-for="handover in pendingHandovers" :key="handover.id" 
            class="border-2 border-yellow-200 rounded-xl p-4 bg-yellow-50">
            
            <!-- Handover Header -->
            <div class="flex justify-between items-start mb-3">
              <div>
                <h4 class="font-bold text-lg">{{ handover.waiter_name }}</h4>
                <p class="text-sm text-gray-500">{{ formatDate(handover.handover_at) }}</p>
                <span :class="getHandoverTypeClass(handover.handover_type)"
                  class="inline-block px-2 py-1 rounded-full text-xs font-medium mt-1">
                  {{ getHandoverTypeText(handover.handover_type) }}
                </span>
              </div>
              <div class="text-right">
                <p class="text-2xl font-bold text-green-600">{{ formatPrice(handover.amount) }}</p>
                <p v-if="handover.handover_type === 'END_SHIFT'" class="text-sm text-gray-500">
                  Tiền cuối ca: {{ formatPrice(handover.end_cash || 0) }}
                </p>
              </div>
            </div>
            
            <!-- Waiter Note -->
            <div v-if="handover.waiter_note" class="bg-blue-50 p-3 rounded-lg mb-3">
              <p class="text-sm text-blue-800">
                <strong>Ghi chú từ waiter:</strong><br>
                {{ handover.waiter_note }}
              </p>
            </div>
            
            <!-- Action Buttons -->
            <div class="flex gap-2">
              <button @click="showConfirmModal(handover, 'CONFIRMED')"
                class="flex-1 bg-green-500 hover:bg-green-600 text-white px-4 py-2 rounded-xl font-medium">
                ✅ Xác nhận
              </button>
              <button @click="showConfirmModal(handover, 'REJECTED')"
                class="flex-1 bg-red-500 hover:bg-red-600 text-white px-4 py-2 rounded-xl font-medium">
                ❌ Từ chối
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Today's Handovers -->
      <div class="bg-white rounded-2xl p-6 shadow-sm">
        <h3 class="text-xl font-bold mb-4">📋 Bàn giao hôm nay</h3>
        
        <div v-if="todayHandovers.length === 0" class="text-center py-10">
          <div class="text-4xl mb-2">📭</div>
          <p class="text-gray-500">Chưa có bàn giao nào hôm nay</p>
        </div>
        
        <div v-else class="space-y-3">
          <div v-for="handover in todayHandovers" :key="handover.id" 
            class="border rounded-xl p-4">
            <div class="flex justify-between items-start mb-2">
              <div>
                <h4 class="font-bold">{{ handover.waiter_name }}</h4>
                <p class="text-sm text-gray-500">{{ formatTime(handover.handover_at) }}</p>
                <span :class="getHandoverTypeClass(handover.handover_type)"
                  class="inline-block px-2 py-1 rounded-full text-xs font-medium mt-1">
                  {{ getHandoverTypeText(handover.handover_type) }}
                </span>
              </div>
              <div class="text-right">
                <p class="font-bold text-lg">{{ formatPrice(handover.amount) }}</p>
                <span :class="getHandoverStatusClass(handover.status)"
                  class="px-2 py-1 rounded-full text-xs font-medium">
                  {{ getHandoverStatusText(handover.status) }}
                </span>
              </div>
            </div>
            
            <div v-if="handover.cashier_note" class="text-sm text-gray-600 mt-2">
              <strong>Ghi chú của bạn:</strong> {{ handover.cashier_note }}
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Confirm Modal -->
    <transition name="slide-up">
      <div v-if="showConfirmForm" class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-end">
        <div class="bg-white rounded-t-3xl w-full p-6">
          <h3 class="text-xl font-bold mb-4">
            {{ confirmAction === 'CONFIRMED' ? '✅ Xác nhận bàn giao' : '❌ Từ chối bàn giao' }}
          </h3>
          
          <!-- Handover Summary -->
          <div class="bg-gray-50 p-4 rounded-xl mb-4">
            <div class="flex justify-between items-center mb-2">
              <span class="text-sm text-gray-600">Waiter</span>
              <span class="font-medium">{{ selectedHandover?.waiter_name }}</span>
            </div>
            <div class="flex justify-between items-center mb-2">
              <span class="text-sm text-gray-600">Số tiền</span>
              <span class="font-bold text-lg">{{ formatPrice(selectedHandover?.amount || 0) }}</span>
            </div>
            <div class="flex justify-between items-center">
              <span class="text-sm text-gray-600">Loại</span>
              <span class="text-sm">{{ getHandoverTypeText(selectedHandover?.handover_type) }}</span>
            </div>
          </div>
          
          <form @submit.prevent="confirmHandover" class="space-y-4">
            <!-- Cashier Note -->
            <div>
              <label class="block text-sm font-medium mb-2">
                {{ confirmAction === 'CONFIRMED' ? 'Ghi chú xác nhận' : 'Lý do từ chối' }}
                {{ confirmAction === 'REJECTED' ? ' *' : '' }}
              </label>
              <textarea v-model="confirmForm.cashier_note" 
                :required="confirmAction === 'REJECTED'"
                rows="3" 
                class="w-full p-3 border rounded-xl focus:ring-2 focus:ring-blue-500"
                :placeholder="confirmAction === 'CONFIRMED' ? 'Ghi chú về việc nhận tiền...' : 'Lý do từ chối bàn giao...'"></textarea>
            </div>
            
            <!-- Action Buttons -->
            <div class="flex gap-2">
              <button type="button" @click="showConfirmForm = false" 
                class="flex-1 bg-gray-200 text-gray-700 px-4 py-3 rounded-xl font-medium">
                Hủy
              </button>
              <button type="submit" 
                :class="[
                  'flex-1 px-4 py-3 rounded-xl font-medium',
                  confirmAction === 'CONFIRMED' 
                    ? 'bg-green-500 hover:bg-green-600 text-white' 
                    : 'bg-red-500 hover:bg-red-600 text-white'
                ]">
                {{ confirmAction === 'CONFIRMED' ? 'Xác nhận' : 'Từ chối' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useCashierStore } from '../stores/cashier'
import { useAuthStore } from '../stores/auth'

const cashierStore = useCashierStore()
const authStore = useAuthStore()

const showConfirmForm = ref(false)
const selectedHandover = ref(null)
const confirmAction = ref('')
const confirmForm = ref({
  cashier_note: ''
})

const loading = computed(() => cashierStore.loading)
const pendingHandovers = computed(() => cashierStore.pendingHandovers)
const todayHandovers = computed(() => cashierStore.todayHandovers)

onMounted(async () => {
  await cashierStore.fetchPendingHandovers()
  await cashierStore.fetchTodayHandovers()
})

const showConfirmModal = (handover, action) => {
  selectedHandover.value = handover
  confirmAction.value = action
  confirmForm.value.cashier_note = ''
  showConfirmForm.value = true
}

const confirmHandover = async () => {
  try {
    await cashierStore.confirmHandover(selectedHandover.value.id, {
      status: confirmAction.value,
      cashier_note: confirmForm.value.cashier_note
    })
    
    showConfirmForm.value = false
    selectedHandover.value = null
    confirmForm.value.cashier_note = ''
    
    // Refresh data
    await cashierStore.fetchPendingHandovers()
    await cashierStore.fetchTodayHandovers()
    
    const message = confirmAction.value === 'CONFIRMED' 
      ? 'Đã xác nhận bàn giao thành công!' 
      : 'Đã từ chối bàn giao!'
    alert(message)
  } catch (error) {
    alert('Lỗi: ' + (error.response?.data?.error || error.message))
  }
}

// Helper functions
const formatPrice = (price) => {
  return new Intl.NumberFormat('vi-VN', { 
    style: 'currency', 
    currency: 'VND',
    maximumFractionDigits: 0
  }).format(price)
}

const formatDate = (date) => {
  return new Date(date).toLocaleString('vi-VN')
}

const formatTime = (date) => {
  return new Date(date).toLocaleTimeString('vi-VN', { 
    hour: '2-digit', 
    minute: '2-digit' 
  })
}

const getHandoverTypeText = (type) => {
  const types = {
    'PARTIAL': 'Một phần',
    'END_SHIFT': 'Toàn bộ + Đóng ca'
  }
  return types[type] || type
}

const getHandoverTypeClass = (type) => {
  const classes = {
    'PARTIAL': 'bg-yellow-100 text-yellow-800',
    'END_SHIFT': 'bg-orange-100 text-orange-800'
  }
  return classes[type] || 'bg-gray-100 text-gray-800'
}

const getHandoverStatusText = (status) => {
  const statuses = {
    'PENDING': 'Chờ xác nhận',
    'CONFIRMED': 'Đã xác nhận',
    'REJECTED': 'Đã từ chối'
  }
  return statuses[status] || status
}

const getHandoverStatusClass = (status) => {
  const classes = {
    'PENDING': 'bg-yellow-100 text-yellow-800',
    'CONFIRMED': 'bg-green-100 text-green-800',
    'REJECTED': 'bg-red-100 text-red-800'
  }
  return classes[status] || 'bg-gray-100 text-gray-800'
}
</script>
```
### 3. Cập Nhật CashierDashboard.vue

```vue
<!-- Thêm vào CashierDashboard.vue -->
<template>
  <!-- ... existing content ... -->
  
  <!-- Handover Notifications Section -->
  <div v-if="pendingHandovers.length > 0" class="bg-yellow-50 border-l-4 border-yellow-400 p-4 mb-4">
    <div class="flex items-center justify-between">
      <div class="flex items-center">
        <div class="flex-shrink-0">
          <svg class="h-5 w-5 text-yellow-400" viewBox="0 0 20 20" fill="currentColor">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-12a1 1 0 10-2 0v4a1 1 0 00.293.707l2.828 2.829a1 1 0 101.415-1.415L11 9.586V6z" clip-rule="evenodd" />
          </svg>
        </div>
        <div class="ml-3">
          <p class="text-sm text-yellow-700">
            <strong>{{ pendingHandovers.length }} yêu cầu bàn giao</strong> đang chờ xác nhận
          </p>
        </div>
      </div>
      <button @click="$router.push('/cashier/handovers')" 
        class="bg-yellow-500 hover:bg-yellow-600 text-white px-4 py-2 rounded-lg text-sm font-medium">
        Xem ngay
      </button>
    </div>
  </div>

  <!-- Quick Handover Actions -->
  <div v-if="pendingHandovers.length > 0" class="bg-white rounded-2xl p-6 shadow-sm mb-4">
    <h3 class="text-lg font-bold mb-4">⚡ Bàn giao nhanh</h3>
    <div class="space-y-3">
      <div v-for="handover in pendingHandovers.slice(0, 3)" :key="handover.id" 
        class="flex items-center justify-between p-3 bg-gray-50 rounded-xl">
        <div>
          <p class="font-medium">{{ handover.waiter_name }}</p>
          <p class="text-sm text-gray-500">{{ formatPrice(handover.amount) }}</p>
        </div>
        <div class="flex gap-2">
          <button @click="quickConfirm(handover.id, 'CONFIRMED')"
            class="bg-green-500 hover:bg-green-600 text-white px-3 py-1 rounded-lg text-sm">
            ✅
          </button>
          <button @click="quickConfirm(handover.id, 'REJECTED')"
            class="bg-red-500 hover:bg-red-600 text-white px-3 py-1 rounded-lg text-sm">
            ❌
          </button>
        </div>
      </div>
    </div>
    <button v-if="pendingHandovers.length > 3" @click="$router.push('/cashier/handovers')"
      class="w-full mt-3 text-blue-600 hover:text-blue-800 text-sm font-medium">
      Xem tất cả {{ pendingHandovers.length }} yêu cầu →
    </button>
  </div>
</template>
```

### 4. Script Updates cho ShiftView.vue (Waiter)
```javascript
// Thêm vào script setup của ShiftView.vue
const showPartialHandoverForm = ref(false)
const showHandoverEndShiftForm = ref(false)
const pendingHandover = ref(null)
const handoverHistory = ref([])

const partialHandoverForm = ref({
  amount: 0,
  waiter_note: ''
})

const handoverEndShiftForm = ref({
  end_cash: 0,
  waiter_note: ''
})

const isWaiter = computed(() => authStore.user?.role === 'waiter')

// Fetch pending handover and history
onMounted(async () => {
  await shiftStore.fetchCurrentShift()
  if (isWaiter.value) {
    await fetchHandoverData()
  }
  // ... existing onMounted code
})

const fetchHandoverData = async () => {
  try {
    pendingHandover.value = await shiftStore.getPendingHandover(currentShift.value?.id)
    handoverHistory.value = await shiftStore.getHandoverHistory(currentShift.value?.id)
  } catch (error) {
    console.error('Error fetching handover data:', error)
  }
}

// Partial Handover Function
const createPartialHandover = async () => {
  try {
    const handoverData = {
      amount: partialHandoverForm.value.amount,
      handover_type: 'PARTIAL',
      waiter_note: partialHandoverForm.value.waiter_note
    }
    
    await shiftStore.createCashHandover(currentShift.value.id, handoverData)
    showPartialHandoverForm.value = false
    partialHandoverForm.value = { amount: 0, waiter_note: '' }
    
    // Refresh data
    await shiftStore.fetchCurrentShift()
    await fetchHandoverData()
    
    alert('Đã gửi yêu cầu bàn giao một phần tiền. Chờ thu ngân xác nhận.')
  } catch (error) {
    alert('Lỗi: ' + (error.response?.data?.error || error.message))
  }
}

// Handover and End Shift Function
const createHandoverAndEndShift = async () => {
  try {
    const handoverData = {
      amount: currentShift.value?.remaining_cash || 0,
      handover_type: 'END_SHIFT',
      waiter_note: handoverEndShiftForm.value.waiter_note,
      end_cash: handoverEndShiftForm.value.end_cash
    }
    
    await shiftStore.createHandoverAndEndShift(currentShift.value.id, handoverData)
    showHandoverEndShiftForm.value = false
    handoverEndShiftForm.value = { end_cash: 0, waiter_note: '' }
    
    // Refresh data
    await shiftStore.fetchCurrentShift()
    await fetchHandoverData()
    
    alert('Đã gửi yêu cầu bàn giao toàn bộ và đóng ca. Chờ thu ngân xác nhận.')
  } catch (error) {
    alert('Lỗi: ' + (error.response?.data?.error || error.message))
  }
}

// Cancel handover
const cancelHandover = async (handoverId) => {
  if (confirm('Bạn có chắc muốn hủy yêu cầu bàn giao này?')) {
    try {
      await shiftStore.cancelHandover(handoverId)
      await fetchHandoverData()
      alert('Đã hủy yêu cầu bàn giao!')
    } catch (error) {
      alert('Lỗi: ' + (error.response?.data?.error || error.message))
    }
  }
}

// Helper functions
const getHandoverTypeText = (type) => {
  const types = {
    'PARTIAL': 'Một phần',
    'END_SHIFT': 'Toàn bộ + Đóng ca'
  }
  return types[type] || type
}

const getHandoverStatusText = (status) => {
  const statuses = {
    'PENDING': 'Chờ xác nhận',
    'CONFIRMED': 'Đã xác nhận',
    'REJECTED': 'Đã từ chối'
  }
  return statuses[status] || status
}

const getHandoverStatusClass = (status) => {
  const classes = {
    'PENDING': 'bg-yellow-100 text-yellow-800',
    'CONFIRMED': 'bg-green-100 text-green-800',
    'REJECTED': 'bg-red-100 text-red-800'
  }
  return classes[status] || 'bg-gray-100 text-gray-800'
}
```

### 5. Store Updates
```javascript
// frontend/src/stores/shift.js (Waiter Store)
export const useShiftStore = defineStore('shift', () => {
  // ... existing code ...
  
  const createCashHandover = async (shiftId, handoverData) => {
    try {
      const response = await api.post(`/api/shifts/${shiftId}/handover`, handoverData)
      return response.data
    } catch (error) {
      console.error('Error creating cash handover:', error)
      throw error
    }
  }
  
  const createHandoverAndEndShift = async (shiftId, handoverData) => {
    try {
      const response = await api.post(`/api/shifts/${shiftId}/handover-and-end`, handoverData)
      return response.data
    } catch (error) {
      console.error('Error creating handover and end shift:', error)
      throw error
    }
  }
  
  const getPendingHandover = async (shiftId) => {
    try {
      const response = await api.get(`/api/shifts/${shiftId}/pending-handover`)
      return response.data
    } catch (error) {
      console.error('Error fetching pending handover:', error)
      return null
    }
  }
  
  const getHandoverHistory = async (shiftId) => {
    try {
      const response = await api.get(`/api/shifts/${shiftId}/handovers`)
      return response.data
    } catch (error) {
      console.error('Error fetching handover history:', error)
      return []
    }
  }
  
  const cancelHandover = async (handoverId) => {
    try {
      const response = await api.delete(`/api/cash-handovers/${handoverId}`)
      return response.data
    } catch (error) {
      console.error('Error canceling handover:', error)
      throw error
    }
  }
  
  return {
    // ... existing returns ...
    createCashHandover,
    createHandoverAndEndShift,
    getPendingHandover,
    getHandoverHistory,
    cancelHandover
  }
})

// frontend/src/stores/cashier.js (Cashier Store)
export const useCashierStore = defineStore('cashier', () => {
  const pendingHandovers = ref([])
  const todayHandovers = ref([])
  const loading = ref(false)
  
  const fetchPendingHandovers = async () => {
    loading.value = true
    try {
      const response = await api.get('/api/cash-handovers/pending')
      pendingHandovers.value = response.data
    } catch (error) {
      console.error('Error fetching pending handovers:', error)
      throw error
    } finally {
      loading.value = false
    }
  }
  
  const fetchTodayHandovers = async () => {
    try {
      const response = await api.get('/api/cash-handovers/today')
      todayHandovers.value = response.data
    } catch (error) {
      console.error('Error fetching today handovers:', error)
      throw error
    }
  }
  
  const confirmHandover = async (handoverId, confirmData) => {
    try {
      const response = await api.post(`/api/cash-handovers/${handoverId}/confirm`, confirmData)
      return response.data
    } catch (error) {
      console.error('Error confirming handover:', error)
      throw error
    }
  }
  
  const quickConfirm = async (handoverId, status) => {
    try {
      const response = await api.post(`/api/cash-handovers/${handoverId}/quick-confirm`, { status })
      // Refresh pending handovers
      await fetchPendingHandovers()
      return response.data
    } catch (error) {
      console.error('Error quick confirming handover:', error)
      throw error
    }
  }
  
  return {
    pendingHandovers,
    todayHandovers,
    loading,
    fetchPendingHandovers,
    fetchTodayHandovers,
    confirmHandover,
    quickConfirm
  }
})
```

## 🔗 API Endpoints
```javascript
// frontend/src/router/index.js
const routes = [
  // ... existing routes ...
  {
    path: '/cashier/handovers',
    name: 'CashierHandovers',
    component: () => import('../views/CashierHandoverView.vue'),
    meta: { requiresAuth: true, roles: ['cashier', 'manager'] }
  }
]
```

### 7. Navigation Updates
```vue
<!-- Thêm vào Navigation.vue cho cashier -->
<router-link v-if="isCashier" to="/cashier/handovers" 
  class="flex items-center px-4 py-2 text-gray-700 hover:bg-gray-100 rounded-lg">
  <span class="mr-3">💰</span>
  <span>Bàn giao tiền</span>
  <span v-if="pendingHandoversCount > 0" 
    class="ml-auto bg-red-500 text-white text-xs px-2 py-1 rounded-full">
    {{ pendingHandoversCount }}
  </span>
</router-link>
```

```
# Waiter Endpoints
POST   /api/shifts/:id/handover               # Tạo yêu cầu bàn giao một phần
POST   /api/shifts/:id/handover-and-end       # Tạo yêu cầu bàn giao toàn bộ và đóng ca
GET    /api/shifts/:id/pending-handover       # Lấy handover đang pending của ca
GET    /api/shifts/:id/handovers              # Lịch sử bàn giao của ca
DELETE /api/cash-handovers/:id                # Hủy yêu cầu bàn giao (chỉ khi PENDING)

# Cashier Endpoints
GET    /api/cash-handovers/pending            # Lấy danh sách chờ xác nhận
GET    /api/cash-handovers/today              # Lấy bàn giao hôm nay
POST   /api/cash-handovers/:id/confirm        # Xác nhận/từ chối bàn giao với ghi chú
POST   /api/cash-handovers/:id/quick-confirm  # Xác nhận nhanh không cần ghi chú

# Shared Endpoints
GET    /api/cash-handovers/my-requests        # Yêu cầu bàn giao của tôi (waiter)
GET    /api/cash-handovers/history            # Lịch sử bàn giao (cả hai role)
```

## 🎯 User Experience Flow với Tương Tác

### Waiter Flow:

#### 1. Bàn Giao Một Phần:
1. **Trong ca**: Nhìn thấy "Tiền hiện có" và nút "💰 Bàn giao một phần"
2. **Click nút**: Mở modal với form nhập số tiền và ghi chú
3. **Nhập thông tin**: Số tiền (≤ tiền hiện có) và ghi chú tùy chọn
4. **Gửi yêu cầu**: Hệ thống tạo handover record với status PENDING
5. **Trạng thái chờ**: 
   - Hiển thị banner "🕐 Đang chờ xác nhận bàn giao"
   - Disable các nút handover khác
   - Có nút "Hủy" để hủy yêu cầu
6. **Nhận phản hồi**: 
   - Nếu CONFIRMED: Cập nhật tiền hiện có, có thể tiếp tục bàn giao
   - Nếu REJECTED: Hiển thị lý do từ chối, có thể tạo yêu cầu mới
7. **Lịch sử**: Xem tất cả handover đã thực hiện với status và ghi chú

#### 2. Bàn Giao và Đóng Ca:
1. **Trong ca**: Nhìn thấy nút "🏁 Bàn giao và đóng ca"
2. **Click nút**: Mở modal với cảnh báo và form
3. **Nhập thông tin**: 
   - Tiền cuối ca (thường là 0)
   - Ghi chú tùy chọn
   - Hiển thị số tiền sẽ bàn giao (= toàn bộ remaining_cash)
4. **Gửi yêu cầu**: Hệ thống tạo handover record với type END_SHIFT
5. **Trạng thái chờ**: 
   - Hiển thị "Chờ cashier xác nhận để đóng ca"
   - Không thể thực hiện thao tác nào khác
6. **Tự động đóng ca**: Sau khi cashier xác nhận, ca tự động đóng
7. **Hoàn thành**: Waiter không thể thao tác gì thêm với ca này

#### 3. Tương Tác Real-time:
- **Notifications**: Nhận thông báo khi cashier phản hồi
- **Status Updates**: Cập nhật trạng thái handover real-time
- **Cancel Option**: Có thể hủy yêu cầu khi đang PENDING

### Cashier Flow:

#### 1. Nhận Yêu Cầu:
1. **Dashboard Alert**: Hiển thị số lượng yêu cầu pending
2. **Quick Actions**: Có thể xác nhận/từ chối nhanh ngay từ dashboard
3. **Detailed View**: Click "Xem ngay" để vào trang quản lý handover

#### 2. Trang Quản Lý Handover:
1. **Pending Section**: 
   - Danh sách yêu cầu chờ xác nhận
   - Hiển thị: Waiter, số tiền, loại handover, thời gian, ghi chú
   - Nút "Xác nhận" và "Từ chối" cho từng yêu cầu
2. **Today's Handovers**: Lịch sử bàn giao hôm nay với status

#### 3. Xác Nhận/Từ Chối:
1. **Click nút**: Mở modal xác nhận với thông tin chi tiết
2. **Nhập ghi chú**: 
   - Xác nhận: Ghi chú tùy chọn
   - Từ chối: Bắt buộc nhập lý do
3. **Submit**: Cập nhật status và gửi phản hồi cho waiter
4. **Auto Update**: 
   - Nếu CONFIRMED: Cập nhật cash amounts cho cả hai ca
   - Nếu END_SHIFT + CONFIRMED: Tự động đóng ca waiter

#### 4. Quick Actions:
1. **Dashboard**: Có thể xác nhận/từ chối nhanh với ✅/❌
2. **No Note Required**: Quick confirm không cần ghi chú
3. **Bulk Actions**: Xử lý nhiều yêu cầu nhanh chóng

### Real-time Interactions:

#### 1. Notifications:
- **Waiter → Cashier**: "Yêu cầu bàn giao mới từ [Waiter Name]"
- **Cashier → Waiter**: "Bàn giao đã được xác nhận/từ chối"

#### 2. Status Synchronization:
- **Real-time Updates**: Cả hai bên thấy status changes ngay lập tức
- **Auto Refresh**: Tự động refresh data khi có thay đổi

#### 3. Error Handling:
- **Connection Issues**: Retry mechanism và offline indicators
- **Conflict Resolution**: Xử lý khi có thay đổi đồng thời

### Communication Flow:
```
Waiter                    System                    Cashier
  |                         |                         |
  |-- Create Handover ----->|                         |
  |                         |-- Notification ------->|
  |<-- Pending Status ------|                         |
  |                         |<-- View Request -------|
  |                         |<-- Confirm/Reject -----|
  |<-- Status Update -------|-- Update Status ------>|
  |<-- Cash Updated --------|-- Cash Updated ------->|
```

## 🔒 Security & Validation

### Backend Validation:
- Waiter chỉ có thể bàn giao tiền từ ca của mình
- Số tiền không được vượt quá `remaining_cash`
- Chỉ cashier được phân công mới có thể xác nhận
- Không thể bàn giao khi không có ca cashier mở

### Frontend Validation:
- Disable nút "Kết thúc ca" nếu còn tiền chưa bàn giao
- Giới hạn input số tiền tối đa = `remaining_cash`
- Hiển thị cảnh báo khi cố gắng đóng ca mà chưa bàn giao hết

## 📊 Reporting & Audit

### Báo Cáo Bàn Giao:
- Tổng tiền bàn giao theo ca
- Lịch sử bàn giao theo waiter
- Thống kê thời gian xác nhận
- Báo cáo sai lệch (nếu có)

### Audit Trail:
- Ghi lại tất cả thao tác bàn giao
- Timestamp cho mỗi bước
- User tracking đầy đủ
- Immutable records sau khi confirmed

---

Tính năng này đảm bảo tính minh bạch và kiểm soát chặt chẽ trong việc quản lý tiền mặt giữa waiter và cashier, đồng thời tích hợp mượt mà vào quy trình làm việc hiện tại.

---

## 🎨 Frontend Implementation với Đối Soát Chi Tiết

### 1. Waiter Interface - Enhanced ShiftView.vue

#### Partial Handover Modal (Simplified)
```vue
<!-- Partial Handover Modal -->
<transition name="slide-up">
  <div v-if="showPartialHandoverForm" class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-end">
    <div class="bg-white rounded-t-3xl w-full p-6">
      <h3 class="text-xl font-bold mb-4">💰 Bàn giao một phần tiền</h3>
      
      <!-- Current Cash Info -->
      <div class="bg-blue-50 p-4 rounded-xl mb-4">
        <div class="flex justify-between items-center">
          <span class="text-sm text-gray-600">Tiền hiện có</span>
          <span class="font-bold text-2xl text-blue-600">{{ formatPrice(currentShift?.remaining_cash || 0) }}</span>
        </div>
      </div>
      
      <form @submit.prevent="createPartialHandover" class="space-y-4">
        <!-- Amount Input -->
        <div>
          <label class="block text-sm font-medium mb-2">Số tiền bàn giao (VNĐ) *</label>
          <input v-model.number="partialHandoverForm.declared_amount" 
            type="number" 
            :max="currentShift?.remaining_cash || 0"
            min="1000" 
            step="1000" 
            required 
            class="w-full p-3 border rounded-xl text-lg font-bold focus:ring-2 focus:ring-yellow-500">
        </div>
        
        <!-- Note -->
        <div>
          <label class="block text-sm font-medium mb-2">Ghi chú (tùy chọn)</label>
          <textarea v-model="partialHandoverForm.waiter_note" 
            rows="3" 
            class="w-full p-3 border rounded-xl focus:ring-2 focus:ring-yellow-500"
            placeholder="Ghi chú về việc bàn giao..."></textarea>
        </div>
        
        <!-- Action Buttons -->
        <div class="flex gap-2">
          <button type="button" @click="showPartialHandoverForm = false" 
            class="flex-1 bg-gray-200 text-gray-700 px-4 py-3 rounded-xl font-medium">
            Hủy
          </button>
          <button type="submit" 
            class="flex-1 bg-yellow-500 hover:bg-yellow-600 text-white px-4 py-3 rounded-xl font-medium">
            Bàn giao
          </button>
        </div>
      </form>
    </div>
  </div>
</transition>
```

### 2. Cashier Interface - Enhanced CashierHandoverView.vue

#### Reconciliation Modal với Discrepancy Handling
```vue
<!-- Reconciliation Modal -->
<transition name="slide-up">
  <div v-if="showReconcileForm" class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-end">
    <div class="bg-white rounded-t-3xl w-full p-6 max-h-[90vh] overflow-y-auto">
      <h3 class="text-xl font-bold mb-4">🔍 Đối soát bàn giao</h3>
      
      <!-- Handover Summary -->
      <div class="bg-gray-50 p-4 rounded-xl mb-4">
        <div class="flex justify-between items-center mb-2">
          <span class="text-sm text-gray-600">Waiter</span>
          <span class="font-medium">{{ selectedHandover?.waiter_name }}</span>
        </div>
        <div class="flex justify-between items-center mb-2">
          <span class="text-sm text-gray-600">Số tiền khai báo</span>
          <span class="font-bold text-lg">{{ formatPrice(selectedHandover?.declared_amount || 0) }}</span>
        </div>
        <div class="flex justify-between items-center">
          <span class="text-sm text-gray-600">Loại bàn giao</span>
          <span class="text-sm">{{ getHandoverTypeText(selectedHandover?.handover_type) }}</span>
        </div>
      </div>
      
      <!-- Cash Breakdown Display (if provided) - REMOVED -->
      
      <form @submit.prevent="reconcileHandover" class="space-y-4">
        <!-- Actual Amount Input -->
        <div>
          <label class="block text-sm font-medium mb-2">Số tiền thực nhận (VNĐ) *</label>
          <input v-model.number="reconcileForm.actual_amount" 
            type="number" 
            min="0" 
            step="1000" 
            required 
            @input="calculateDiscrepancy"
            class="w-full p-3 border rounded-xl text-lg font-bold focus:ring-2 focus:ring-blue-500">
        </div>
        
        <!-- Discrepancy Display -->
        <div v-if="discrepancy !== 0" class="p-4 rounded-xl" :class="discrepancy > 0 ? 'bg-green-50 border border-green-200' : 'bg-red-50 border border-red-200'">
          <div class="flex justify-between items-center mb-2">
            <span class="text-sm font-medium">Chênh lệch:</span>
            <span class="font-bold text-lg" :class="discrepancy > 0 ? 'text-green-600' : 'text-red-600'">
              {{ discrepancy > 0 ? '+' : '' }}{{ formatPrice(discrepancy) }}
            </span>
          </div>
          <p class="text-xs" :class="discrepancy > 0 ? 'text-green-700' : 'text-red-700'">
            {{ discrepancy > 0 ? '✅ Thừa tiền' : '⚠️ Thiếu tiền' }}
          </p>
        </div>
        
        <!-- Discrepancy Reason (if discrepancy exists) -->
        <div v-if="discrepancy !== 0">
          <label class="block text-sm font-medium mb-2">Lý do chênh lệch *</label>
          <select v-model="reconcileForm.discrepancy_reason" required 
            class="w-full p-3 border rounded-xl focus:ring-2 focus:ring-blue-500 mb-2">
            <option value="">-- Chọn lý do --</option>
            <option value="COUNTING_ERROR">Lỗi đếm tiền</option>
            <option value="TRANSACTION_ERROR">Lỗi giao dịch</option>
            <option value="CUSTOMER_ISSUE">Vấn đề khách hàng</option>
            <option value="SYSTEM_ERROR">Lỗi hệ thống</option>
            <option value="OTHER">Khác</option>
          </select>
          
          <label class="block text-sm font-medium mb-2">Trách nhiệm</label>
          <select v-model="reconcileForm.discrepancy_responsibility" required 
            class="w-full p-3 border rounded-xl focus:ring-2 focus:ring-blue-500">
            <option value="">-- Chọn trách nhiệm --</option>
            <option value="WAITER">Waiter</option>
            <option value="CASHIER">Cashier</option>
            <option value="CUSTOMER">Khách hàng</option>
            <option value="SYSTEM">Hệ thống</option>
            <option value="UNKNOWN">Không rõ</option>
          </select>
        </div>
        
        <!-- Large Discrepancy Warning -->
        <div v-if="Math.abs(discrepancy) > discrepancyThreshold" class="bg-orange-50 border border-orange-200 p-4 rounded-xl">
          <div class="flex items-center">
            <svg class="h-5 w-5 text-orange-400 mr-2" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
            </svg>
            <div>
              <p class="text-sm font-medium text-orange-800">Chênh lệch lớn</p>
              <p class="text-xs text-orange-700">Cần sự phê duyệt từ quản lý</p>
            </div>
          </div>
        </div>
        
        <!-- Cashier Note -->
        <div>
          <label class="block text-sm font-medium mb-2">Ghi chú đối soát</label>
          <textarea v-model="reconcileForm.cashier_note" 
            rows="3" 
            class="w-full p-3 border rounded-xl focus:ring-2 focus:ring-blue-500"
            placeholder="Ghi chú về quá trình đối soát..."></textarea>
        </div>
        
        <!-- Action Buttons -->
        <div class="flex gap-2">
          <button type="button" @click="showReconcileForm = false" 
            class="flex-1 bg-gray-200 text-gray-700 px-4 py-3 rounded-xl font-medium">
            Hủy
          </button>
          <button type="button" @click="rejectHandover"
            class="flex-1 bg-red-500 hover:bg-red-600 text-white px-4 py-3 rounded-xl font-medium">
            Từ chối
          </button>
          <button type="submit" 
            class="flex-1 bg-green-500 hover:bg-green-600 text-white px-4 py-3 rounded-xl font-medium">
            Xác nhận
          </button>
        </div>
      </form>
    </div>
  </div>
</transition>
```

---

## 🔗 API Endpoints Mở Rộng với Đối Soát

```
# Waiter Endpoints
POST   /api/shifts/:id/handover               # Tạo yêu cầu bàn giao với breakdown
POST   /api/shifts/:id/handover-and-end       # Tạo yêu cầu bàn giao toàn bộ và đóng ca
GET    /api/shifts/:id/pending-handover       # Lấy handover đang pending của ca
GET    /api/shifts/:id/handovers              # Lịch sử bàn giao của ca
DELETE /api/cash-handovers/:id                # Hủy yêu cầu bàn giao (chỉ khi PENDING)

# Cashier Endpoints
GET    /api/cash-handovers/pending            # Lấy danh sách chờ xác nhận
GET    /api/cash-handovers/today              # Lấy bàn giao hôm nay
POST   /api/cash-handovers/:id/reconcile      # Đối soát với actual amount
POST   /api/cash-handovers/:id/quick-confirm  # Xác nhận nhanh không cần ghi chú
GET    /api/cash-handovers/discrepancy-stats  # Thống kê chênh lệch

# Manager Endpoints
GET    /api/cash-handovers/pending-approval   # Chênh lệch cần phê duyệt
POST   /api/cash-handovers/:id/approve        # Phê duyệt/từ chối chênh lệch
GET    /api/discrepancies/stats               # Thống kê chênh lệch chi tiết
GET    /api/discrepancies/history             # Lịch sử chênh lệch

# Shared Endpoints
GET    /api/cash-handovers/my-requests        # Yêu cầu bàn giao của tôi (waiter)
GET    /api/cash-handovers/history            # Lịch sử bàn giao (cả hai role)
```

---

## 🎯 Quy Trình Đối Soát Chi Tiết

### 1. **Waiter Handover Process:**
```
1. Waiter khai báo số tiền bàn giao
2. Tạo handover record với status PENDING
3. Gửi notification cho cashier
```

### 2. **Cashier Reconciliation Process:**
```
1. Cashier nhận notification
2. Cashier xem thông tin handover
3. Cashier đếm tiền thực tế
4. Cashier nhập actual amount
5. Hệ thống tự động tính discrepancy
6. Nếu có chênh lệch:
   - Cashier chọn lý do và trách nhiệm
   - Nếu chênh lệch > threshold → cần manager approval
7. Cashier xác nhận hoặc từ chối
```

### 3. **Manager Approval Process (nếu cần):**
```
1. Manager nhận notification về chênh lệch lớn
2. Manager xem chi tiết handover và discrepancy
3. Manager phê duyệt hoặc từ chối với ghi chú
4. Nếu phê duyệt → cập nhật cash amounts
5. Nếu từ chối → handover status = REJECTED
```

### 4. **System Updates:**
```
1. Waiter shift: 
   - handed_over_cash += actual_amount
   - remaining_cash -= declared_amount
   - total_discrepancy += discrepancy
2. Cashier shift:
   - received_cash += actual_amount
   - total_discrepancy += discrepancy
3. Audit trail: Ghi lại tất cả thay đổi
```

---

## 📊 Báo Cáo Đối Soát & Chênh Lệch

### 1. **Discrepancy Dashboard:**
- Tổng số lần handover
- Số lần có chênh lệch
- Tổng số tiền chênh lệch
- Phân tích theo nguyên nhân
- Top waiter/cashier có chênh lệch nhiều

### 2. **Audit Reports:**
- Chi tiết từng giao dịch handover
- Timeline đầy đủ với timestamps
- User actions và approvals
- Discrepancy resolution tracking

### 3. **Performance Metrics:**
- Accuracy rate per user
- Average discrepancy amount
- Resolution time
- Manager approval frequency

---

## 🔒 Security & Compliance

### 1. **Data Integrity:**
- Immutable audit trail
- Cryptographic signatures cho critical data
- Backup và recovery procedures
- Data retention policies

### 2. **Access Control:**
- Role-based permissions
- Manager approval workflows
- Audit log access restrictions
- Sensitive data encryption

### 3. **Compliance Features:**
- SOX compliance reporting
- Financial audit trails
- Regulatory reporting
- Data privacy protection

---

Thiết kế đối soát chi tiết này đảm bảo:
- ✅ **Accuracy**: Đối soát chính xác giữa khai báo và thực tế
- ✅ **Transparency**: Theo dõi đầy đủ mọi chênh lệch
- ✅ **Accountability**: Xác định trách nhiệm rõ ràng
- ✅ **Control**: Manager approval cho chênh lệch lớn
- ✅ **Audit**: Audit trail hoàn chỉnh cho compliance
- ✅ **Reporting**: Báo cáo chi tiết và thống kê