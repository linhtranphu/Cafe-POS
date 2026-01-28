# How to Restart Backend After Code Changes

## Problem
After making code changes and rebuilding the backend, the old process is still running.

## Solution

### Step 1: Find the running process
```bash
ps aux | grep cafe-pos-server | grep -v grep
```

Output example:
```
tranphulinh  41318  0.0  0.1  34875652  16012 s020  S+  9:38PM  0:02.18 ./cafe-pos-server
```

The PID is `41318`

### Step 2: Kill the old process
```bash
kill 41318
```

Or kill all cafe-pos-server processes:
```bash
pkill -f cafe-pos-server
```

### Step 3: Rebuild (if needed)
```bash
cd backend
go build -o cafe-pos-server
```

### Step 4: Start the new backend
```bash
cd backend
./cafe-pos-server
```

Or run in background:
```bash
cd backend
nohup ./cafe-pos-server > server.log 2>&1 &
```

### Step 5: Verify it's running
```bash
ps aux | grep cafe-pos-server | grep -v grep
```

Check the start time matches current time.

## Quick Restart Script

Create `restart-backend.sh`:
```bash
#!/bin/bash

echo "Stopping old backend..."
pkill -f cafe-pos-server

echo "Building new backend..."
cd backend
go build -o cafe-pos-server

echo "Starting new backend..."
nohup ./cafe-pos-server > server.log 2>&1 &

echo "Waiting for backend to start..."
sleep 2

echo "Backend status:"
ps aux | grep cafe-pos-server | grep -v grep

echo ""
echo "Backend restarted successfully!"
echo "Logs: tail -f backend/server.log"
```

Then:
```bash
chmod +x restart-backend.sh
./restart-backend.sh
```

## Testing the New Code

After restart, test the barista shift validation:

```bash
# Login as barista
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "barista1", "password": "barista123"}'

# Save the token from response
TOKEN="<your_token>"

# Try to accept order without shift (should fail with 400)
curl -X POST http://localhost:8080/api/barista/orders/<order_id>/accept \
  -H "Authorization: Bearer $TOKEN"

# Expected response:
# {"error":"barista must open a shift before accepting orders"}
```

## Common Issues

### Issue 1: Port already in use
```
bind: address already in use
```

**Solution**: Kill the old process first
```bash
pkill -f cafe-pos-server
# or
lsof -ti:8080 | xargs kill
```

### Issue 2: Permission denied
```
permission denied: ./cafe-pos-server
```

**Solution**: Make it executable
```bash
chmod +x backend/cafe-pos-server
```

### Issue 3: Changes not reflected
**Solution**: Make sure you:
1. Saved all files
2. Rebuilt the binary
3. Killed the old process
4. Started the new process

Check binary timestamp:
```bash
ls -lh backend/cafe-pos-server
```

Check process start time:
```bash
ps -p <PID> -o lstart,command
```

They should match!
