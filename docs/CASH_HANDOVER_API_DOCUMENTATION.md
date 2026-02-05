# 📚 Cash Handover API Documentation

## 🎯 Overview

This document provides complete API documentation for the cash handover feature, including all endpoints, request/response formats, and error codes.

**Base URL:** `/api`  
**Authentication:** Bearer token required for all endpoints  
**Content-Type:** `application/json`

---

## 🔐 Authentication

All endpoints require a valid JWT token in the Authorization header:

```
Authorization: Bearer <token>
```

---

## 👨‍💼 Waiter Endpoints

### 1. Create Handover

Create a partial or full cash handover request.

**Endpoint:** `POST /shifts/:shift_id/handover`  
**Role:** Waiter  
**Rate Limit:** 10 requests/minute

**Path Parameters:**
- `shift_id` (string, required) - Waiter shift ID

**Request Body:**
```json
{
  "declared_amount": 200000,
  "handover_type": "PARTIAL",
  "waiter_note": "Bàn giao tiền ca sáng"
}
```

**Request Fields:**
- `declared_amount` (number, required) - Amount to handover (VND), must be > 0 and ≤ remaining_cash
- `handover_type` (string, required) - Type: "PARTIAL" or "FULL"
- `waiter_note` (string, optional) - Note from waiter (max 500 chars)

**Response:** `201 Created`
```json
{
  "id": "507f1f77bcf86cd799439011",
  "waiter_shift_id": "507f1f77bcf86cd799439012",
  "cashier_shift_id": "507f1f77bcf86cd799439013",
  "waiter_id": "507f1f77bcf86cd799439014",
  "waiter_name": "Nguyễn Văn A",
  "cashier_id": "507f1f77bcf86cd799439015",
  "cashier_name": "Trần Thị B",
  "declared_amount": 200000,
  "actual_amount": 0,
  "discrepancy": 0,
  "handover_type": "PARTIAL",
  "status": "PENDING",
  "waiter_note": "Bàn giao tiền ca sáng",
  "handover_at": "2026-02-04T10:30:00Z",
  "created_at": "2026-02-04T10:30:00Z",
  "updated_at": "2026-02-04T10:30:00Z"
}
```

**Error Responses:**
- `400 Bad Request` - Invalid input or validation error
- `401 Unauthorized` - Missing or invalid token
- `403 Forbidden` - Not your shift
- `404 Not Found` - Shift not found
- `409 Conflict` - Pending handover already exists

**Example:**
```bash
curl -X POST http://localhost:8080/api/shifts/507f1f77bcf86cd799439012/handover \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "declared_amount": 200000,
    "handover_type": "PARTIAL",
    "waiter_note": "Bàn giao tiền ca sáng"
  }'
```

---

### 2. Create Handover and End Shift

Create handover for all remaining cash and automatically end shift after cashier confirms.

**Endpoint:** `POST /shifts/:shift_id/handover-and-end`  
**Role:** Waiter

**Request Body:**
```json
{
  "declared_amount": 500000,
  "waiter_note": "Kết thúc ca",
  "end_cash": 0
}
```

**Request Fields:**
- `declared_amount` (number, required) - Should equal remaining_cash
- `waiter_note` (string, optional) - Note from waiter
- `end_cash` (number, required) - Final cash amount (usually 0)

**Response:** `201 Created`
```json
{
  "id": "507f1f77bcf86cd799439011",
  "handover_type": "END_SHIFT",
  "status": "PENDING",
  "end_cash": 0,
  ...
}
```

---

### 3. Get Pending Handover

Get the pending handover for current shift.

**Endpoint:** `GET /shifts/:shift_id/pending-handover`  
**Role:** Waiter

**Response:** `200 OK`
```json
{
  "id": "507f1f77bcf86cd799439011",
  "declared_amount": 200000,
  "status": "PENDING",
  "handover_at": "2026-02-04T10:30:00Z",
  ...
}
```

**Error Responses:**
- `404 Not Found` - No pending handover

---

### 4. Get Handover History

Get all handovers for a shift.

**Endpoint:** `GET /shifts/:shift_id/handovers`  
**Role:** Waiter

**Response:** `200 OK`
```json
[
  {
    "id": "507f1f77bcf86cd799439011",
    "declared_amount": 200000,
    "actual_amount": 200000,
    "discrepancy": 0,
    "status": "CONFIRMED",
    "handover_at": "2026-02-04T10:30:00Z",
    "confirmed_at": "2026-02-04T10:35:00Z",
    ...
  },
  {
    "id": "507f1f77bcf86cd799439012",
    "declared_amount": 150000,
    "actual_amount": 145000,
    "discrepancy": -5000,
    "status": "CONFIRMED",
    ...
  }
]
```

---

### 5. Cancel Handover

Cancel a pending handover (only PENDING status).

**Endpoint:** `DELETE /cash-handovers/:handover_id`  
**Role:** Waiter

**Response:** `200 OK`
```json
{
  "message": "handover cancelled successfully"
}
```

**Error Responses:**
- `400 Bad Request` - Cannot cancel non-pending handover
- `403 Forbidden` - Not your handover

---

## 💰 Cashier Endpoints

### 6. Get Pending Handovers

Get all pending handovers assigned to current cashier.

**Endpoint:** `GET /cash-handovers/pending`  
**Role:** Cashier, Manager

**Response:** `200 OK`
```json
[
  {
    "id": "507f1f77bcf86cd799439011",
    "waiter_name": "Nguyễn Văn A",
    "declared_amount": 200000,
    "handover_type": "PARTIAL",
    "status": "PENDING",
    "waiter_note": "Bàn giao tiền ca sáng",
    "handover_at": "2026-02-04T10:30:00Z",
    ...
  }
]
```

---

### 7. Get Today's Handovers

Get all handovers for today.

**Endpoint:** `GET /cash-handovers/today`  
**Role:** Cashier, Manager

**Response:** `200 OK`
```json
[
  {
    "id": "507f1f77bcf86cd799439011",
    "waiter_name": "Nguyễn Văn A",
    "declared_amount": 200000,
    "actual_amount": 200000,
    "discrepancy": 0,
    "status": "CONFIRMED",
    "handover_at": "2026-02-04T10:30:00Z",
    "confirmed_at": "2026-02-04T10:35:00Z",
    ...
  }
]
```

---

### 8. Confirm Handover

Confirm or reject a handover with reconciliation.

**Endpoint:** `POST /cash-handovers/:handover_id/confirm`  
**Role:** Cashier, Manager

**Request Body (Confirm):**
```json
{
  "actual_amount": 200000,
  "status": "CONFIRMED",
  "cashier_note": "Đã nhận đủ tiền"
}
```

**Request Body (Confirm with Discrepancy):**
```json
{
  "actual_amount": 195000,
  "status": "CONFIRMED",
  "cashier_note": "Thiếu 5k",
  "discrepancy_reason": "COUNTING_ERROR",
  "discrepancy_responsibility": "WAITER"
}
```

**Request Body (Reject):**
```json
{
  "status": "REJECTED",
  "cashier_note": "Số tiền không khớp với hệ thống"
}
```

**Request Fields:**
- `actual_amount` (number, required for CONFIRMED) - Actual amount received
- `status` (string, required) - "CONFIRMED" or "REJECTED"
- `cashier_note` (string, required for REJECTED) - Note from cashier
- `discrepancy_reason` (string, optional) - Reason if discrepancy exists
  - Values: "COUNTING_ERROR", "TRANSACTION_ERROR", "CUSTOMER_ISSUE", "SYSTEM_ERROR", "OTHER"
- `discrepancy_responsibility` (string, optional) - Who is responsible
  - Values: "WAITER", "CASHIER", "CUSTOMER", "SYSTEM", "UNKNOWN"

**Response:** `200 OK`
```json
{
  "message": "handover confirmed successfully"
}
```

**Note:** If discrepancy > 100,000 VND, requires manager approval before cash amounts are updated.

---

### 9. Quick Confirm

Quick confirm/reject without detailed reconciliation.

**Endpoint:** `POST /cash-handovers/:handover_id/quick-confirm`  
**Role:** Cashier, Manager

**Request Body:**
```json
{
  "status": "CONFIRMED"
}
```

**Response:** `200 OK`
```json
{
  "message": "handover quick confirmed"
}
```

**Note:** Assumes actual_amount = declared_amount

---

## 👨‍💼 Manager Endpoints

### 10. Get Pending Approvals

Get handovers with large discrepancies requiring approval.

**Endpoint:** `GET /cash-handovers/pending-approval`  
**Role:** Manager

**Response:** `200 OK`
```json
[
  {
    "id": "507f1f77bcf86cd799439011",
    "waiter_name": "Nguyễn Văn A",
    "cashier_name": "Trần Thị B",
    "declared_amount": 500000,
    "actual_amount": 380000,
    "discrepancy": -120000,
    "status": "DISCREPANCY",
    "requires_approval": true,
    "discrepancy_reason": "COUNTING_ERROR",
    ...
  }
]
```

---

### 11. Approve Discrepancy

Approve or reject a large discrepancy.

**Endpoint:** `POST /cash-handovers/:handover_id/approve`  
**Role:** Manager

**Request Body:**
```json
{
  "approved": true,
  "manager_note": "Đã xác minh, chấp nhận chênh lệch"
}
```

**Request Fields:**
- `approved` (boolean, required) - true to approve, false to reject
- `manager_note` (string, required) - Manager's note

**Response:** `200 OK`
```json
{
  "message": "discrepancy approved successfully"
}
```

**Note:** If approved, cash amounts are updated. If rejected, handover status becomes REJECTED.

---

### 12. Get Discrepancy Statistics

Get statistics about discrepancies for a date range.

**Endpoint:** `GET /cash-handovers/discrepancy-stats`  
**Role:** Manager

**Query Parameters:**
- `start` (string, required) - Start date (YYYY-MM-DD)
- `end` (string, required) - End date (YYYY-MM-DD)

**Response:** `200 OK`
```json
{
  "total_handovers": 150,
  "total_discrepancy": -25000,
  "shortage_count": 12,
  "overage_count": 3,
  "shortage_amount": 35000,
  "overage_amount": 10000,
  "required_approval": 2
}
```

---

## 📊 Data Models

### HandoverStatus
- `PENDING` - Waiting for cashier confirmation
- `CONFIRMED` - Confirmed by cashier
- `REJECTED` - Rejected by cashier
- `DISCREPANCY` - Has large discrepancy, waiting for manager approval

### HandoverType
- `PARTIAL` - Partial handover during shift
- `FULL` - Full handover of all remaining cash
- `END_SHIFT` - Handover and end shift

### DiscrepancyType
- `SHORTAGE` - Actual < Declared (thiếu tiền)
- `OVERAGE` - Actual > Declared (thừa tiền)

### ResponsibilityType
- `WAITER` - Waiter's responsibility
- `CASHIER` - Cashier's responsibility
- `CUSTOMER` - Customer issue
- `SYSTEM` - System error
- `UNKNOWN` - Unknown cause

---

## ⚠️ Error Codes

| Code | Message | Description |
|------|---------|-------------|
| 400 | Bad Request | Invalid input or validation error |
| 401 | Unauthorized | Missing or invalid authentication token |
| 403 | Forbidden | Insufficient permissions |
| 404 | Not Found | Resource not found |
| 409 | Conflict | Resource conflict (e.g., pending handover exists) |
| 422 | Unprocessable Entity | Business logic validation failed |
| 429 | Too Many Requests | Rate limit exceeded |
| 500 | Internal Server Error | Server error |

---

## 🔒 Security

### Authentication
- All endpoints require valid JWT token
- Token must be included in Authorization header
- Token expires after 24 hours

### Authorization
- Role-based access control (RBAC)
- Waiter can only access own shifts
- Cashier can only access assigned handovers
- Manager has full access

### Data Validation
- All amounts must be positive numbers
- Declared amount cannot exceed remaining cash
- Status transitions are validated
- Input sanitization to prevent XSS/SQL injection

---

## 📈 Rate Limiting

| Endpoint | Limit |
|----------|-------|
| Create Handover | 10/minute |
| Confirm Handover | 20/minute |
| Get Pending | 30/minute |
| Get History | 30/minute |
| Other endpoints | 60/minute |

---

## 🧪 Testing

### Example Test Flow

```bash
# 1. Login as waiter
TOKEN=$(curl -s -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"waiter1","password":"password123"}' \
  | jq -r '.token')

# 2. Create handover
HANDOVER=$(curl -s -X POST http://localhost:8080/api/shifts/$SHIFT_ID/handover \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"declared_amount":200000,"handover_type":"PARTIAL"}' \
  | jq -r '.id')

# 3. Login as cashier
CASHIER_TOKEN=$(curl -s -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"cashier1","password":"password123"}' \
  | jq -r '.token')

# 4. Confirm handover
curl -X POST http://localhost:8080/api/cash-handovers/$HANDOVER/confirm \
  -H "Authorization: Bearer $CASHIER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"actual_amount":200000,"status":"CONFIRMED"}'
```

---

## 📝 Notes

- All timestamps are in UTC (ISO 8601 format)
- All amounts are in Vietnamese Dong (VND)
- Discrepancy threshold for manager approval: 100,000 VND
- Handovers are immutable after confirmation (audit trail)
- Cancelled handovers are soft-deleted (status changed, not removed)

---

**Last Updated:** 2026-02-04  
**API Version:** 1.0  
**Contact:** dev@cafepos.com
