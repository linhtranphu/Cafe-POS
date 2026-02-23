# Transaction Patterns: Long-Running vs Optimistic Locking vs Saga

## Vấn đề: Multi-Step Workflow

Quy trình đóng ca cashier có nhiều bước:
1. Bắt đầu đóng ca
2. Nhập tiền thực tế
3. Giải trình chênh lệch (nếu có)
4. Hoàn tất đóng ca

**Câu hỏi**: Làm sao đảm bảo nếu user bấm "Back" giữa chừng thì không có data nào được lưu?

---

## Pattern 1: Long-Running Transaction ❌ (NOT Recommended)

### Concept
Giữ một transaction mở suốt toàn bộ quy trình, chỉ commit khi hoàn tất.

### Implementation
```go
// Start transaction when user clicks "Bắt đầu đóng ca"
session, _ := mongoClient.StartSession()
tx, _ := session.StartTransaction()

// Step 1: Initiate closure (not committed)
shift.Status = CLOSURE_INITIATED
// Data in transaction buffer, not in DB

// User waits 2 minutes...

// Step 2: Record actual cash (not committed)
shift.ActualCash = 500000
// Still in transaction buffer

// User waits 3 minutes...

// Step 3: Document variance (not committed)
shift.Variance = {...}
// Still in transaction buffer

// Finally: Commit or Rollback
if userClicksComplete {
    tx.Commit()  // All changes saved
} else {
    tx.Rollback()  // All changes discarded
}
```

### Flow Diagram
```
User Action          Transaction State           Database State
-----------          -----------------           --------------
Start closure   →    TX OPEN                     No changes
                     [status=CLOSURE_INITIATED]
                     
Wait 2 mins...  →    TX OPEN (holding locks)     No changes
                     
Record cash     →    TX OPEN                     No changes
                     [actual_cash=500k]
                     
Wait 3 mins...  →    TX OPEN (holding locks)     No changes
                     
Complete        →    TX COMMIT                   All saved ✅
OR
Back            →    TX ROLLBACK                 Nothing saved ✅
```

### Vấn đề

#### 1. Transaction Timeout
```
MongoDB default timeout: 60 seconds
User takes 5 minutes → Transaction timeout → All work lost
```

#### 2. Database Locks
```
Transaction holds locks on shift document
Other operations (view shift, handover) blocked
Performance degradation
```

#### 3. Connection Issues
```
Network hiccup during 5-minute wait
→ Connection lost
→ Transaction lost
→ User has to start over
```

#### 4. Memory Usage
```
Transaction buffer holds all changes in memory
Multiple users = Multiple open transactions
High memory usage on DB server
```

#### 5. Scalability
```
10 cashiers doing closure simultaneously
= 10 long-running transactions
= 10x locks, 10x memory, 10x connections
```

### Khi nào dùng?
- ❌ KHÔNG dùng cho user-facing workflows
- ✅ Chỉ dùng cho automated processes (< 1 second)
- ✅ Batch jobs chạy nhanh

---

## Pattern 2: Optimistic Locking ✅ (Recommended)

### Concept
Không giữ transaction lâu. Mỗi bước save ngay, nhưng track version để detect conflicts.

### Implementation

#### Database Schema
```go
type CashierShift struct {
    ID          primitive.ObjectID
    Status      string
    ActualCash  *float64
    Variance    *Variance
    Version     int  // ✅ Version field for optimistic locking
    UpdatedAt   time.Time
}
```

#### Code
```go
// Step 1: Initiate closure
func InitiateClosure(shiftID) {
    shift := getShift(shiftID)
    currentVersion := shift.Version
    
    shift.Status = CLOSURE_INITIATED
    shift.Version++  // Increment version
    
    // Update with version check
    result := db.UpdateOne(
        filter: {_id: shiftID, version: currentVersion},  // ✅ Check version
        update: {$set: shift}
    )
    
    if result.MatchedCount == 0 {
        return error("Conflict: shift was modified by another process")
    }
}

// Step 2: Record actual cash
func RecordActualCash(shiftID, actualCash) {
    shift := getShift(shiftID)
    currentVersion := shift.Version
    
    shift.ActualCash = actualCash
    shift.Version++
    
    result := db.UpdateOne(
        filter: {_id: shiftID, version: currentVersion},
        update: {$set: shift}
    )
    
    if result.MatchedCount == 0 {
        return error("Conflict: shift was modified")
    }
}

// Cancel: Just call cancel API
func CancelClosure(shiftID) {
    shift := getShift(shiftID)
    currentVersion := shift.Version
    
    shift.Status = OPEN
    shift.ActualCash = nil
    shift.Variance = nil
    shift.Version++
    
    result := db.UpdateOne(
        filter: {_id: shiftID, version: currentVersion},
        update: {$set: shift}
    )
}
```

### Flow Diagram
```
User Action          Database State                      Version
-----------          --------------                      -------
Start closure   →    status=CLOSURE_INITIATED           v1 → v2
                     COMMITTED ✅
                     
Wait 2 mins...  →    (no locks, other ops can proceed)  v2
                     
Record cash     →    actual_cash=500k                   v2 → v3
                     COMMITTED ✅
                     
Wait 3 mins...  →    (no locks)                         v3
                     
Back            →    Call Cancel API                    v3 → v4
                     status=OPEN, actual_cash=null
                     COMMITTED ✅
```

### Ưu điểm

#### 1. No Timeout Issues ✅
```
Mỗi operation < 1 second
User có thể đợi bao lâu cũng được
```

#### 2. No Locks ✅
```
Không hold locks
Other operations không bị block
Better performance
```

#### 3. Conflict Detection ✅
```
Nếu 2 users cùng edit:
User A: version 5 → 6 (success)
User B: version 5 → 6 (fail - conflict detected)
User B: Retry with version 6
```

#### 4. Scalability ✅
```
100 concurrent users = OK
Each operation independent
No resource contention
```

### Nhược điểm

#### 1. Data Visible Immediately
```
Step 2 committed → actual_cash visible in DB
Nếu user bấm Back, cần call Cancel API để cleanup
```

#### 2. Requires Cancel Logic
```
Phải implement cancel/rollback logic
Không tự động như transaction rollback
```

### Khi nào dùng?
- ✅ User-facing workflows
- ✅ Long-running processes
- ✅ High concurrency scenarios
- ✅ Distributed systems

---

## Pattern 3: Saga Pattern ✅ (Advanced)

### Concept
Chia workflow thành các bước độc lập, mỗi bước có compensating action để rollback.

### Implementation

#### Define Saga Steps
```go
type SagaStep struct {
    Name           string
    Execute        func() error
    Compensate     func() error  // Rollback action
}

type CashierClosureSaga struct {
    Steps []SagaStep
    CompletedSteps []string
}
```

#### Saga Definition
```go
saga := CashierClosureSaga{
    Steps: []SagaStep{
        {
            Name: "initiate_closure",
            Execute: func() error {
                shift.Status = CLOSURE_INITIATED
                return save(shift)
            },
            Compensate: func() error {
                shift.Status = OPEN
                return save(shift)
            },
        },
        {
            Name: "record_actual_cash",
            Execute: func() error {
                shift.ActualCash = actualCash
                shift.Variance = calculateVariance()
                return save(shift)
            },
            Compensate: func() error {
                shift.ActualCash = nil
                shift.Variance = nil
                return save(shift)
            },
        },
        {
            Name: "document_variance",
            Execute: func() error {
                shift.Variance.Reason = reason
                shift.Variance.Notes = notes
                return save(shift)
            },
            Compensate: func() error {
                shift.Variance.Reason = nil
                shift.Variance.Notes = ""
                return save(shift)
            },
        },
        {
            Name: "close_shift",
            Execute: func() error {
                shift.Status = CLOSED
                shift.EndTime = now()
                return save(shift)
            },
            Compensate: func() error {
                // Cannot compensate close
                return errors.New("cannot rollback close")
            },
        },
    },
}
```

#### Saga Execution
```go
func ExecuteSaga(saga *CashierClosureSaga) error {
    for _, step := range saga.Steps {
        err := step.Execute()
        if err != nil {
            // Rollback all completed steps
            return saga.Rollback()
        }
        saga.CompletedSteps = append(saga.CompletedSteps, step.Name)
    }
    return nil
}

func (s *CashierClosureSaga) Rollback() error {
    // Execute compensating actions in reverse order
    for i := len(s.CompletedSteps) - 1; i >= 0; i-- {
        stepName := s.CompletedSteps[i]
        step := s.findStep(stepName)
        err := step.Compensate()
        if err != nil {
            // Log compensation failure
            log.Error("Failed to compensate step: %s", stepName)
        }
    }
    return nil
}
```

### Flow Diagram
```
User Action          Saga State                    Database State
-----------          ----------                    --------------
Start closure   →    Step 1 EXECUTED              status=CLOSURE_INITIATED ✅
                     CompletedSteps: [step1]
                     
Record cash     →    Step 2 EXECUTED              actual_cash=500k ✅
                     CompletedSteps: [step1, step2]
                     
Back clicked    →    ROLLBACK INITIATED
                     
                     Compensate step2             actual_cash=null ✅
                     Compensate step1             status=OPEN ✅
                     
                     ROLLBACK COMPLETE            Clean state ✅
```

### Ưu điểm

#### 1. Explicit Rollback Logic ✅
```
Mỗi step có compensating action rõ ràng
Dễ test, dễ maintain
```

#### 2. Partial Failure Handling ✅
```
Nếu step 3 fail:
→ Compensate step 2
→ Compensate step 1
→ System về trạng thái ban đầu
```

#### 3. Audit Trail ✅
```
Track được:
- Steps completed
- Steps compensated
- Failure reasons
```

#### 4. Distributed Transactions ✅
```
Có thể span across multiple services
Step 1: Update shift (Service A)
Step 2: Send notification (Service B)
Step 3: Update analytics (Service C)
```

### Nhược điểm

#### 1. Complexity
```
Phải define compensating actions cho mọi step
More code to write and maintain
```

#### 2. Eventual Consistency
```
Trong quá trình rollback, data có thể inconsistent tạm thời
```

#### 3. Compensation Failures
```
Nếu compensating action fail → Manual intervention needed
```

### Khi nào dùng?
- ✅ Complex workflows với nhiều bước
- ✅ Distributed systems
- ✅ Need explicit rollback logic
- ✅ Need audit trail

---

## So sánh 3 Patterns

| Aspect | Long-Running TX | Optimistic Locking | Saga Pattern |
|--------|----------------|-------------------|--------------|
| **Timeout Risk** | ❌ High | ✅ None | ✅ None |
| **Locks** | ❌ Holds locks | ✅ No locks | ✅ No locks |
| **Scalability** | ❌ Poor | ✅ Excellent | ✅ Excellent |
| **Complexity** | ✅ Simple | ✅ Simple | ⚠️ Complex |
| **Rollback** | ✅ Automatic | ⚠️ Manual API | ✅ Explicit |
| **Audit Trail** | ⚠️ Limited | ⚠️ Limited | ✅ Detailed |
| **Distributed** | ❌ No | ⚠️ Limited | ✅ Yes |
| **Use Case** | Batch jobs | User workflows | Complex workflows |

---

## Recommendation cho Cashier Closure

### Current Implementation: Optimistic Locking ✅

Đây là approach tốt nhất vì:

1. **Simple**: Mỗi bước một transaction ngắn
2. **No timeout**: User có thể đợi bao lâu cũng được
3. **Scalable**: Nhiều cashiers cùng lúc OK
4. **Cancel API**: Đã có sẵn để rollback

### Flow hiện tại
```
Bước 1: Initiate → COMMIT (status=CLOSURE_INITIATED)
Bước 2: Record Cash → COMMIT (actual_cash=500k)
Bấm Back → Call Cancel API → COMMIT (rollback to OPEN)
```

### Nếu muốn "không lưu gì khi bấm Back"

**Option A: Frontend-Only** (Simplest) ✅
```
- Giữ tất cả data ở frontend
- Chỉ call API một lần khi "Hoàn tất"
- Bấm Back = discard frontend data
- Backend nhận full data, validate, save trong 1 transaction
```

**Option B: Draft State** (More complex)
```
- Mỗi bước save vào "draft" collection
- Khi "Hoàn tất" → Move draft → main trong transaction
- Bấm Back → Delete draft
```

**Option C: Saga Pattern** (Most complex)
```
- Define compensating actions
- Track completed steps
- Auto-rollback on cancel
```

---

## Code Example: Frontend-Only Approach

### Frontend
```vue
<script setup>
const closureData = ref({
  status: null,
  actualCash: null,
  variance: null,
  varianceReason: null,
  varianceNotes: null
})

// Step 1: Just update local state
const initiateClosure = () => {
  closureData.value.status = 'CLOSURE_INITIATED'
}

// Step 2: Just update local state
const recordActualCash = (amount) => {
  closureData.value.actualCash = amount
  closureData.value.variance = calculateVariance(amount)
}

// Step 3: Just update local state
const documentVariance = (reason, notes) => {
  closureData.value.varianceReason = reason
  closureData.value.varianceNotes = notes
}

// Final: Call API once with all data
const completeClosureconst completeClosure = async () => {
  await cashierShiftService.completeClosure(shiftId, closureData.value)
  // All data saved in one transaction
}

// Back: Just discard local data
const goBack = () => {
  closureData.value = {}  // Discard everything
  router.push('/cashier')
}
</script>
```

### Backend
```go
func CompleteClosure(shiftID, data ClosureData) error {
    // One transaction for everything
    session, _ := mongoClient.StartSession()
    
    _, err := session.WithTransaction(ctx, func(sessCtx) {
        // 1. Get shift
        shift := getShift(shiftID)
        
        // 2. Validate all data
        if err := validateClosureData(data); err != nil {
            return err
        }
        
        // 3. Apply all changes
        shift.Status = CLOSURE_INITIATED
        shift.ActualCash = data.ActualCash
        shift.Variance = data.Variance
        if data.Variance.Amount != 0 {
            shift.Variance.Reason = data.VarianceReason
            shift.Variance.Notes = data.VarianceNotes
        }
        shift.Status = CLOSED
        shift.EndTime = now()
        
        // 4. Save once
        return save(shift)
    })
    
    return err
}
```

Approach này:
- ✅ Simple
- ✅ No timeout
- ✅ Bấm Back = không lưu gì
- ✅ One short transaction
- ✅ Easy to test

Bạn muốn tôi implement approach này không?
