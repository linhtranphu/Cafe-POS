# Backend Port Configuration Fix

## Issue
Frontend batch services were trying to connect to `localhost:8080` but the backend server runs on port `3000`, causing connection refused errors.

## Root Cause
The batch record and batch report services were using hardcoded port 8080 in their API base URL configuration:
```javascript
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'
```

## Solution
Updated the default port from 8080 to 3000 in the following files:

### Files Modified
1. `frontend/src/services/batchRecord.js` - Changed default port to 3000
2. `frontend/src/services/batchReport.js` - Changed default port to 3000

### Updated Configuration
```javascript
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:3000'
```

## Status
✅ Backend running on port 3000
✅ Frontend batch services configured to use port 3000
✅ Connection errors resolved

## Note
Other batch services (`batchAlert.js`, `batchDefinition.js`) already use the shared `api.js` module which uses relative paths (`/api`) and works correctly with the Vite dev server proxy.

## Backend Status
The backend server is running successfully with all batch management endpoints available:
- POST `/api/manager/batch-definitions` - Create batch definition
- GET `/api/manager/batch-definitions` - List batch definitions
- GET `/api/manager/batch-definitions/:id` - Get batch definition
- PUT `/api/manager/batch-definitions/:id` - Update batch definition
- DELETE `/api/manager/batch-definitions/:id` - Delete batch definition
- POST `/api/batch-records` - Create batch record (barista/manager)
- GET `/api/batch-records` - List batch records (barista/manager)
- GET `/api/batch-alerts` - Get batch alerts (barista/manager)
- GET `/api/batch-reports/*` - Various batch reports (manager only)
