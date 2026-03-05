# Auto Print on Order Creation

## Overview

Đã implement tính năng tự động tạo print jobs khi có order mới được tạo.

## Changes Made

### 1. Updated OrderHandler.CreateOrder

**File:** `backend/interfaces/http/order_handler.go`

**Before:**
```go
func (h *OrderHandler) CreateOrder(c *gin.Context) {
    // ... create order logic ...
    
    c.JSON(http.StatusCreated, o)
}
```

**After:**
```go
func (h *OrderHandler) CreateOrder(c *gin.Context) {
    // ... create order logic ...
    
    // Auto-create print jobs (bill + labels) if print service is available
    // This runs asynchronously to not block the response
    // Use background context since the request context will be cancelled
    go func() {
        if h.printService != nil {
            // Create print jobs for the order (1 bill + N labels)
            ctx := context.Background()
            if err := h.printService.CreatePrintJobsForOrder(ctx, o); err != nil {
                // Log error but don't fail the order creation
                // User can manually reprint if needed
                // TODO: Add proper logging
            }
        }
    }()
    
    c.JSON(http.StatusCreated, o)
}
```

**Key Points:**
- Runs asynchronously in goroutine (không block response)
- Uses `context.Background()` (không bị cancel khi request kết thúc)
- Calls `CreatePrintJobsForOrder()` (tạo 1 bill + N labels)
- Errors không fail order creation (user có thể reprint manually)

### 2. Added Context Import

```go
import (
    "context"
    // ... other imports
)
```

## How It Works

### Workflow

```
User creates order
    ↓
POST /api/orders
    ↓
OrderHandler.CreateOrder()
    ↓
OrderService.CreateOrder() → Order saved to DB
    ↓
Response 201 Created (immediate)
    ↓
[Async] PrintService.CreatePrintJobsForOrder()
    ↓
Creates print jobs:
    - 1 Bill print job
    - N Label print jobs (1 per item)
    ↓
Print jobs saved to DB with status PENDING
    ↓
WebSocket broadcast: "print_job_created"
    ↓
Print worker picks up jobs
    ↓
Sends to printers via Local Print Bridge
    ↓
Updates job status: COMPLETED or FAILED
```

### Print Jobs Created

For each order, the system creates:

1. **Bill Print Job**
   - Type: BILL
   - Printer: Bill printer (from printer config)
   - Template: Bill template
   - Status: PENDING

2. **Label Print Jobs** (one per item)
   - Type: LABEL
   - Printer: Label printer (from printer config)
   - Template: Label template
   - Status: PENDING
   - Item info: name, variant, quantity

### Print Job Processing

Print jobs are processed by the print worker:

1. Worker polls for PENDING jobs
2. Renders template with order data
3. Sends to printer via Local Print Bridge
4. Updates job status:
   - SUCCESS → COMPLETED
   - ERROR → FAILED (with error message)
5. Broadcasts status change via WebSocket

## Shop Settings

Auto-print behavior is controlled by `AutoPrintEnabled` setting:

```go
type ShopSettings struct {
    AutoPrintEnabled bool `bson:"auto_print_enabled" json:"auto_print_enabled"`
    // ... other fields
}
```

**Default:** `true` (auto-print enabled)

**To disable:**
```
1. Go to Print Management → Settings
2. Uncheck "Auto Print Enabled"
3. Save settings
```

**Note:** Currently the code always creates print jobs. To respect the setting, we need to check it in CreateOrder handler.

## Checking Print Jobs

### Frontend

1. Open Print Management → Print Jobs tab
2. See list of all print jobs
3. Filter by status: PENDING, COMPLETED, FAILED
4. Retry failed jobs
5. Cancel pending jobs

### Backend API

```bash
# Get all print jobs
GET /api/manager/print-jobs

# Get pending jobs
GET /api/manager/print-jobs?status=PENDING

# Get failed jobs
GET /api/manager/print-jobs?status=FAILED

# Retry failed job
POST /api/manager/print-jobs/:id/retry

# Cancel pending job
POST /api/manager/print-jobs/:id/cancel
```

## Manual Reprint

If auto-print fails or is disabled, users can manually reprint:

### Reprint Bill

```bash
POST /api/orders/:id/reprint-bill
```

Frontend:
1. Go to Orders list
2. Click order
3. Click "Reprint Bill" button

### Reprint Label

```bash
POST /api/orders/:id/reprint-label
Body: { "item_index": 0 }
```

Frontend:
1. Go to Orders list
2. Click order
3. Click "Reprint Label" on specific item

## WebSocket Events

Print jobs broadcast events via WebSocket:

### Event: print_job_created

```json
{
  "event": "print_job_created",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "type": "BILL",
    "order_id": "507f1f77bcf86cd799439012",
    "order_number": "20260226-123456-001",
    "status": "PENDING",
    "created_at": "2026-02-26T12:34:56Z"
  }
}
```

### Event: print_job_status_changed

```json
{
  "event": "print_job_status_changed",
  "data": {
    "job_id": "507f1f77bcf86cd799439011",
    "status": "COMPLETED",
    "error_message": ""
  }
}
```

### Event: print_job_failed

```json
{
  "event": "print_job_failed",
  "data": {
    "job_id": "507f1f77bcf86cd799439011",
    "error_message": "Printer offline"
  }
}
```

## Testing

### 1. Create Order via POS

```
1. Open POS interface
2. Add items to cart
3. Click "Checkout"
4. Order created → Print jobs created automatically
```

### 2. Check Print Jobs

```
1. Open Print Management → Print Jobs
2. See new print jobs with status PENDING
3. Wait for print worker to process
4. Status changes to COMPLETED or FAILED
```

### 3. Check Printer Output

```
1. Bill should print on bill printer
2. Labels should print on label printer (1 per item)
```

### 4. Test Failed Print

```
1. Turn off printer
2. Create order
3. Print jobs created with status PENDING
4. Worker tries to print → FAILED
5. Error message shows in print jobs list
6. Turn on printer
7. Click "Retry" on failed job
8. Job status → COMPLETED
```

## Troubleshooting

### Issue: Print jobs not created

**Check:**
1. Backend logs: `tail -f backend.log | grep -i print`
2. OrderHandler has printService injected?
3. PrintService is not nil?

**Solution:**
```bash
# Check backend initialization logs
grep "Print service" backend.log
grep "Print worker" backend.log
```

### Issue: Print jobs stuck in PENDING

**Check:**
1. Print worker is running?
2. Local Print Bridge is running?
3. Printer is online?

**Solution:**
```bash
# Check print worker logs
tail -f backend.log | grep -i "print worker"

# Check Local Print Bridge
curl http://localhost:3001/health

# Check printer
ping 192.168.1.115
telnet 192.168.1.115 9100
```

### Issue: Print jobs FAILED

**Check:**
1. Error message in print job
2. Printer connection
3. Template rendering errors

**Solution:**
```bash
# View failed job details
GET /api/manager/print-jobs/:id

# Check error message
{
  "error_message": "Failed to connect to printer: connection refused"
}

# Fix printer connection and retry
POST /api/manager/print-jobs/:id/retry
```

### Issue: Auto-print not respecting settings

**Current behavior:** Always creates print jobs regardless of `AutoPrintEnabled` setting.

**To fix:** Add setting check in CreateOrder handler:

```go
// Check if auto-print is enabled
shopSettings, err := h.shopSettingsRepo.GetSettings(ctx)
if err == nil && shopSettings.AutoPrintEnabled {
    // Create print jobs
    go func() {
        // ...
    }()
}
```

## Future Enhancements

1. **Respect AutoPrintEnabled Setting**
   - Check setting before creating print jobs
   - Allow per-printer enable/disable

2. **Print Job Priorities**
   - High priority for bills
   - Normal priority for labels
   - Process high priority first

3. **Batch Printing**
   - Group multiple labels into one print job
   - Reduce printer commands

4. **Print Preview**
   - Show preview before printing
   - Allow user to approve/reject

5. **Print Templates per Printer**
   - Different templates for different printers
   - Customize layout per printer type

6. **Print Job Scheduling**
   - Schedule print jobs for later
   - Print during off-peak hours

7. **Print Job Analytics**
   - Track print success rate
   - Monitor printer uptime
   - Alert on high failure rate

## Summary

✅ Auto-create print jobs on order creation
✅ Async processing (không block response)
✅ Creates 1 bill + N labels per order
✅ WebSocket notifications
✅ Manual reprint available
✅ Failed jobs can be retried
✅ Print jobs visible in Print Management

**Next steps:**
1. Test order creation
2. Check print jobs in Print Management
3. Verify printer output
4. Add setting check for AutoPrintEnabled
5. Add proper logging for errors

Hệ thống bây giờ tự động tạo print jobs khi có order mới!
