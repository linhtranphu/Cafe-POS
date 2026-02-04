# Cash Handover API Documentation

## Overview

The Cash Handover API provides endpoints for managing cash handovers between waiters and cashiers in the cafe POS system. This includes creating handover requests, confirming/rejecting handovers, handling discrepancies, and manager approval workflows.

## Authentication

All endpoints require authentication via JWT token in the Authorization header:
```
Authorization: Bearer <jwt_token>
```

## Role-Based Access Control

- **Waiter**: Can create handovers for their own shifts
- **Cashier**: Can confirm/reject handovers, view pending handovers
- **Manager**: Full access including discrepancy approval

## API Endpoints

### Waiter Endpoints

#### Create Handover
Creates a new cash handover request.

**Endpoint:** `POST /api/shifts/:id/handover`  
**Role:** Waiter, Manager  
**Parameters:**
- `id` (path): Shift ID

**Request Body:**
```json
{
  "type": "PARTIAL|FULL|END_SHIFT",
  "requested_amount": 500000,
  "waiter_notes": "Bàn giao tiền ca sáng"
}
```

**Response:**
```json
{
  "id": "64f...",
  "type": "PARTIAL",
  "status": "PENDING",
  "waiter_shift_id": "64f...",
  "waiter_id": "64f...",
  "waiter_name": "Nguyễn Văn A",
  "requested_amount": 500000,
  "waiter_notes": "Bàn giao tiền ca sáng",
  "requested_at": "2024-02-04T10:30:00Z",
  "created_at": "2024-02-04T10:30:00Z",
  "updated_at": "2024-02-04T10:30:00Z"
}
```

#### Create Handover and End Shift
Creates a handover request and ends the shift when confirmed.

**Endpoint:** `POST /api/shifts/:id/handover-and-end`  
**Role:** Waiter, Manager  
**Parameters:**
- `id` (path): Shift ID

**Request Body:** Same as Create Handover

#### Get Pending Handover
Gets the pending handover for a shift.

**Endpoint:** `GET /api/shifts/:id/pending-handover`  
**Role:** Waiter, Manager  
**Parameters:**
- `id` (path): Shift ID

**Response:** Handover object or 404 if no pending handover

#### Get Handover History
Gets all handovers for a shift.

**Endpoint:** `GET /api/shifts/:id/handovers`  
**Role:** Waiter, Manager  
**Parameters:**
- `id` (path): Shift ID

**Response:**
```json
{
  "handovers": [
    {
      "id": "64f...",
      "type": "PARTIAL",
      "status": "CONFIRMED",
      // ... other handover fields
    }
  ],
  "count": 1
}
```

#### Cancel Handover
Cancels a pending handover (waiter can cancel their own).

**Endpoint:** `DELETE /api/cash-handovers/:id`  
**Role:** Waiter (own handovers), Manager  
**Parameters:**
- `id` (path): Handover ID

**Response:**
```json
{
  "message": "Handover cancelled successfully"
}
```

### Cashier Endpoints

#### Get Pending Handovers
Gets all pending handovers for the cashier.

**Endpoint:** `GET /api/cash-handovers/pending`  
**Role:** Cashier, Manager

**Response:**
```json
{
  "handovers": [
    {
      "id": "64f...",
      "type": "PARTIAL",
      "status": "PENDING",
      "waiter_name": "Nguyễn Văn A",
      "requested_amount": 500000,
      "waiter_notes": "Bàn giao tiền ca sáng",
      "requested_at": "2024-02-04T10:30:00Z"
    }
  ],
  "count": 1
}
```

#### Get Today's Handovers
Gets all handovers for today.

**Endpoint:** `GET /api/cash-handovers/today`  
**Role:** Cashier, Manager

**Response:** Same format as pending handovers

#### Quick Confirm Handover
Quickly confirms a handover with exact amount (no discrepancy).

**Endpoint:** `POST /api/cash-handovers/:id/quick-confirm`  
**Role:** Cashier, Manager  
**Parameters:**
- `id` (path): Handover ID

**Request Body:**
```json
{
  "cashier_notes": "Xác nhận số tiền chính xác"
}
```

**Response:**
```json
{
  "message": "Handover confirmed successfully"
}
```

#### Reconcile Handover
Reconciles a handover with actual amount and discrepancy details.

**Endpoint:** `POST /api/cash-handovers/:id/reconcile`  
**Role:** Cashier, Manager  
**Parameters:**
- `id` (path): Handover ID

**Request Body:**
```json
{
  "actual_amount": 480000,
  "discrepancy_reason": "Thiếu tiền lẻ",
  "responsibility": "WAITER|CASHIER|SYSTEM|UNKNOWN",
  "cashier_notes": "Thiếu 20k so với yêu cầu"
}
```

**Response:**
```json
{
  "message": "Handover reconciled successfully"
}
```

#### Reject Handover
Rejects a handover request.

**Endpoint:** `POST /api/cash-handovers/:id/reject`  
**Role:** Cashier, Manager  
**Parameters:**
- `id` (path): Handover ID

**Request Body:**
```json
{
  "reason": "Số tiền không khớp với báo cáo"
}
```

**Response:**
```json
{
  "message": "Handover rejected successfully"
}
```

#### Get Discrepancy Statistics
Gets discrepancy statistics for a date range.

**Endpoint:** `GET /api/cash-handovers/discrepancy-stats`  
**Role:** Cashier, Manager  
**Query Parameters:**
- `start_date` (optional): Start date (YYYY-MM-DD), defaults to start of current month
- `end_date` (optional): End date (YYYY-MM-DD), defaults to end of today

**Response:**
```json
{
  "stats": {
    "total_discrepancies": 5,
    "total_shortages": 3,
    "total_overages": 2,
    "total_shortage_amount": 150000,
    "total_overage_amount": 75000,
    "net_discrepancy": -75000,
    "pending_count": 1,
    "resolved_count": 4,
    "escalated_count": 0
  },
  "start_date": "2024-02-01",
  "end_date": "2024-02-04"
}
```

### Manager Endpoints

#### Get Pending Approvals
Gets handovers requiring manager approval.

**Endpoint:** `GET /api/manager/cash-handovers/pending-approval`  
**Role:** Manager

**Response:**
```json
{
  "handovers": [],
  "count": 0,
  "message": "Manager approval functionality - to be implemented with repository access"
}
```

#### Approve/Reject Discrepancy
Approves or rejects a discrepancy.

**Endpoint:** `POST /api/manager/cash-handovers/:id/approve`  
**Role:** Manager  
**Parameters:**
- `id` (path): Handover ID

**Request Body:**
```json
{
  "approved": true,
  "manager_notes": "Phê duyệt chênh lệch do lỗi hệ thống"
}
```

**Response:**
```json
{
  "message": "Discrepancy approved successfully"
}
```

#### Get Manager Discrepancy Stats
Gets discrepancy statistics for managers.

**Endpoint:** `GET /api/manager/discrepancies/stats`  
**Role:** Manager  
**Query Parameters:** Same as cashier discrepancy stats

**Response:** Same format as cashier discrepancy stats

## Error Responses

All endpoints return errors in the following format:

```json
{
  "error": "Error message description"
}
```

Common HTTP status codes:
- `400 Bad Request`: Invalid request data or business rule violation
- `401 Unauthorized`: Missing or invalid authentication token
- `403 Forbidden`: Insufficient permissions for the requested action
- `404 Not Found`: Requested resource not found
- `500 Internal Server Error`: Server-side error

## Business Rules

### Handover Creation
- Only waiters can create handovers for their own shifts
- Shift must be open (not closed)
- Requested amount must be > 0 and <= remaining cash in shift
- Only one pending handover allowed per shift
- Maximum handover amount: 10,000,000 VND

### Handover Confirmation
- Only cashiers can confirm/reject handovers
- Cashier must have an open shift
- Discrepancies >= 50,000 VND require manager approval
- END_SHIFT handovers automatically close the waiter shift when confirmed

### Manager Approval
- Required for discrepancies >= 50,000 VND (configurable)
- Managers can approve/reject any discrepancy
- Approval/rejection includes mandatory notes

## Data Types

### HandoverType
- `PARTIAL`: Partial cash handover during shift
- `FULL`: Full cash handover (all remaining cash)
- `END_SHIFT`: Handover at end of shift (closes shift when confirmed)

### HandoverStatus
- `PENDING`: Waiting for cashier confirmation
- `CONFIRMED`: Confirmed by cashier (no discrepancy or approved discrepancy)
- `REJECTED`: Rejected by cashier
- `DISCREPANCY`: Has discrepancy, may require manager approval

### ResponsibilityType
- `WAITER`: Discrepancy is waiter's responsibility
- `CASHIER`: Discrepancy is cashier's responsibility
- `SYSTEM`: Discrepancy due to system error
- `UNKNOWN`: Responsibility not determined

## Configuration

### Environment Variables
- `DISCREPANCY_THRESHOLD`: Amount threshold for manager approval (default: 50000)
- `MAX_HANDOVER_AMOUNT`: Maximum amount per handover (default: 10000000)

### Service Configuration
The discrepancy threshold can be configured via the service:
```go
cashHandoverService.SetDiscrepancyThreshold(100000) // 100k VND
```

## Testing

Use the provided test script to verify API functionality:
```bash
./scripts/test-cash-handover-api.sh
```

This script tests basic endpoint connectivity and authentication. For full functionality testing, you'll need active waiter and cashier shifts.

## Integration Notes

### Database Collections
- `cash_handovers`: Main handover records
- `cash_discrepancies`: Detailed discrepancy tracking
- `shifts`: Updated with cash tracking fields
- `cashier_shifts`: Updated with received cash tracking

### State Management
- Handovers follow a state machine pattern
- State transitions are validated and logged
- Audit trail maintained for all operations

### Performance Considerations
- Database indexes on common query fields
- Pagination for large result sets
- Efficient aggregation queries for statistics