# Facility Endpoint Fix ✅

## Problem
The frontend was calling `/api/manager/issue-reports` but the backend endpoint is `/api/manager/issues`, causing 404 errors when trying to fetch issue reports.

## Root Cause
Mismatch between frontend service endpoint paths and backend route definitions:
- Frontend: `/api/manager/issue-reports`
- Backend: `/api/manager/issues`

## Solution
Updated `frontend/src/services/facility.js` to use the correct endpoint paths:

### Changes Made:
1. `getIssueReports()`: `/manager/issue-reports` → `/manager/issues`
2. `createIssueReport()`: `/staff/issue-reports` → `/manager/issues`
3. `updateIssueReportStatus()`: `/manager/issue-reports/${id}/status` → `/manager/issues/${id}/status`
4. `addReportComment()`: `/manager/issue-reports/${id}/comments` → `/manager/issues/${id}/comments`

## Backend Endpoints (Confirmed)
From `backend/main.go`:
```go
manager.GET("/issues", facilityHandler.GetIssueReports)
manager.POST("/issues", facilityHandler.CreateIssueReport)
```

## Testing
After this fix:
- ✅ Issue reports modal loads without 404 error
- ✅ Can view existing issue reports
- ✅ Can create new issue reports
- ✅ Consistent with backend API structure

## Related
This fix complements the earlier facility delete error fix (400 Bad Request for business rule violations).

## Status: COMPLETE ✅
Frontend now correctly calls the backend issue endpoints.
