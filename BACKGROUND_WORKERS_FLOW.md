# Background Workers Flow Diagram

## Tổng Quan Flow

```
┌─────────────────────────────────────────────────────────────────────┐
│                         SERVER STARTUP                               │
│                                                                      │
│  1. MongoDB Connection                                               │
│  2. Create Services                                                  │
│  3. Create Monitoring Service                                        │
│  4. Create Cost Recalculation Service (4 workers, queue 1000)       │
│  5. Wire up services                                                 │
│  6. ✅ START WORKERS → costRecalculationService.Start()             │
│  7. Setup graceful shutdown                                          │
│  8. Start HTTP server                                                │
└─────────────────────────────────────────────────────────────────────┘
                                ↓
┌─────────────────────────────────────────────────────────────────────┐
│                    WORKERS RUNNING (4 goroutines)                    │
│                                                                      │
│  Worker 1: [IDLE] ──┐                                               │
│  Worker 2: [IDLE] ──┼─→ Waiting for jobs from queue                 │
│  Worker 3: [IDLE] ──┤                                               │
│  Worker 4: [IDLE] ──┘                                               │
│                                                                      │
│  Queue: [ empty ] (capacity: 1000)                                  │
└─────────────────────────────────────────────────────────────────────┘
```

## Flow Khi Cập Nhật Ingredient

```
┌──────────────────────────────────────────────────────────────────────┐
│                    MANAGER CẬP NHẬT INGREDIENT                        │
└──────────────────────────────────────────────────────────────────────┘
                                ↓
                    PATCH /api/manager/ingredients/:id
                    { "cost_per_unit": 250000 }
                                ↓
┌──────────────────────────────────────────────────────────────────────┐
│                      INGREDIENT SERVICE                               │
│                                                                       │
│  1. Update ingredient in database                                    │
│  2. Call costCalculatorService.QueueCostRecalculation(ingredientID)  │
└──────────────────────────────────────────────────────────────────────┘
                                ↓
┌──────────────────────────────────────────────────────────────────────┐
│                    COST CALCULATOR SERVICE                            │
│                                                                       │
│  1. Find all menu items using this ingredient                        │
│     Example: Cappuccino, Latte, Espresso (3 items)                  │
│                                                                       │
│  2. Queue each menu item for recalculation:                          │
│     - recalcService.QueueRecalculation(cappuccino_id)                │
│     - recalcService.QueueRecalculation(latte_id)                     │
│     - recalcService.QueueRecalculation(espresso_id)                  │
└──────────────────────────────────────────────────────────────────────┘
                                ↓
┌──────────────────────────────────────────────────────────────────────┐
│                    RECALCULATION QUEUE                                │
│                                                                       │
│  Queue: [cappuccino_id, latte_id, espresso_id]                      │
│                                                                       │
│  Status:                                                             │
│  - QueuedItems: 3                                                    │
│  - ProcessedItems: 0                                                 │
└──────────────────────────────────────────────────────────────────────┘
                                ↓
┌──────────────────────────────────────────────────────────────────────┐
│                    WORKERS PICK UP JOBS                               │
│                                                                       │
│  Worker 1: [PROCESSING cappuccino_id] ──┐                           │
│  Worker 2: [PROCESSING latte_id]      ──┼─→ Processing in parallel  │
│  Worker 3: [PROCESSING espresso_id]   ──┘                           │
│  Worker 4: [IDLE]                                                    │
│                                                                       │
│  Queue: [ empty ]                                                    │
└──────────────────────────────────────────────────────────────────────┘
                                ↓
┌──────────────────────────────────────────────────────────────────────┐
│                    EACH WORKER PROCESSES JOB                          │
│                                                                       │
│  For each menu item:                                                 │
│  1. Fetch menu item from database                                    │
│  2. Fetch all ingredients                                            │
│  3. Calculate new cost:                                              │
│     cost = Σ(quantity × cost_per_unit × conversion × (1+wastage))   │
│  4. Update menu item in database:                                    │
│     - current_cost = calculated_cost                                 │
│     - cost_status = FINAL or INCOMPLETE                              │
│     - cost_last_calculated_at = now()                                │
│  5. Record metric to monitoring service                              │
└──────────────────────────────────────────────────────────────────────┘
                                ↓
┌──────────────────────────────────────────────────────────────────────┐
│                    MONITORING SERVICE                                 │
│                                                                       │
│  Record metrics for each job:                                        │
│  - Type: recalculation_job                                           │
│  - Status: success/failure                                           │
│  - Duration: 150ms                                                   │
│  - Metadata: menu_item_id, attempts                                  │
│                                                                       │
│  Update aggregated stats:                                            │
│  - total_recalc_jobs: 3                                              │
│  - successful_recalc_jobs: 3                                         │
│  - average_recalc_duration: 150ms                                    │
│                                                                       │
│  Check alert rules:                                                  │
│  - Error rate < 15% ✅                                               │
│  - Failures < 5 in 5min ✅                                           │
└──────────────────────────────────────────────────────────────────────┘
                                ↓
┌──────────────────────────────────────────────────────────────────────┐
│                    RESULT                                             │
│                                                                       │
│  ✅ Cappuccino: current_cost updated to 15,000 VND                   │
│  ✅ Latte: current_cost updated to 18,000 VND                        │
│  ✅ Espresso: current_cost updated to 12,000 VND                     │
│                                                                       │
│  Manager can now see updated costs:                                  │
│  GET /api/manager/menu/costs                                         │
└──────────────────────────────────────────────────────────────────────┘
```

## Flow Khi Server Shutdown

```
┌──────────────────────────────────────────────────────────────────────┐
│                    MANAGER NHẤN Ctrl+C                                │
└──────────────────────────────────────────────────────────────────────┘
                                ↓
┌──────────────────────────────────────────────────────────────────────┐
│                    GRACEFUL SHUTDOWN                                  │
│                                                                       │
│  1. Signal workers to stop (cancel context)                          │
│  2. Wait for current jobs to complete                                │
│  3. Close queue channel                                              │
│  4. Wait for all workers to finish (WaitGroup)                       │
│  5. ✅ All workers stopped safely                                    │
│  6. Close MongoDB connection                                         │
│  7. Exit                                                             │
└──────────────────────────────────────────────────────────────────────┘
```

## Worker Pool Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                    COST RECALCULATION SERVICE                        │
│                                                                      │
│  ┌────────────────────────────────────────────────────────────┐    │
│  │                    QUEUE (Channel)                          │    │
│  │  Capacity: 1000 jobs                                        │    │
│  │  [job1, job2, job3, ...]                                    │    │
│  └────────────────────────────────────────────────────────────┘    │
│                           ↓  ↓  ↓  ↓                                │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐          │
│  │ Worker 1 │  │ Worker 2 │  │ Worker 3 │  │ Worker 4 │          │
│  │          │  │          │  │          │  │          │          │
│  │ Process  │  │ Process  │  │ Process  │  │ Process  │          │
│  │ job from │  │ job from │  │ job from │  │ job from │          │
│  │ queue    │  │ queue    │  │ queue    │  │ queue    │          │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘          │
│       ↓              ↓              ↓              ↓               │
│  ┌────────────────────────────────────────────────────────────┐   │
│  │              MONITORING SERVICE                             │   │
│  │  - Record metrics                                           │   │
│  │  - Check alert rules                                        │   │
│  │  - Update aggregated stats                                  │   │
│  └────────────────────────────────────────────────────────────┘   │
│       ↓              ↓              ↓              ↓               │
│  ┌────────────────────────────────────────────────────────────┐   │
│  │              MONGODB                                        │   │
│  │  - Update menu_items.current_cost                           │   │
│  │  - Update menu_items.cost_status                            │   │
│  │  - Update menu_items.cost_last_calculated_at                │   │
│  └────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────┘
```

## Retry Logic

```
Job Processing Flow:
┌─────────────────────────────────────────────────────────────────────┐
│  Attempt 1                                                           │
│  ├─ Try to process job                                               │
│  ├─ Success? → ✅ Done, record success metric                       │
│  └─ Failure? → Wait 1 second, retry                                 │
│                                                                      │
│  Attempt 2                                                           │
│  ├─ Try to process job                                               │
│  ├─ Success? → ✅ Done, record success metric                       │
│  └─ Failure? → Wait 2 seconds, retry                                │
│                                                                      │
│  Attempt 3 (final)                                                   │
│  ├─ Try to process job                                               │
│  ├─ Success? → ✅ Done, record success metric                       │
│  └─ Failure? → ❌ Give up, record failure metric                    │
│                                                                      │
│  Exponential backoff: 1s → 2s → 4s                                  │
└─────────────────────────────────────────────────────────────────────┘
```

## Monitoring Flow

```
┌─────────────────────────────────────────────────────────────────────┐
│                    CONTINUOUS MONITORING                             │
│                                                                      │
│  Every job completion:                                               │
│  1. Record metric (success/failure, duration)                        │
│  2. Update aggregated stats                                          │
│  3. Check alert rules:                                               │
│     - Count failures in last 5 minutes                               │
│     - Calculate error rate                                           │
│     - Trigger alert if threshold exceeded                            │
│                                                                      │
│  Manager can check anytime:                                          │
│  - GET /api/manager/monitoring/health                                │
│  - GET /api/manager/monitoring/metrics/aggregated                    │
│  - GET /api/manager/monitoring/alerts                                │
└─────────────────────────────────────────────────────────────────────┘
```

## Key Points

1. **Tự động**: Workers start khi server khởi động, không cần config
2. **Song song**: 4 workers xử lý jobs đồng thời
3. **Retry**: Tự động retry 3 lần với exponential backoff
4. **Monitoring**: Mọi operation đều được track và alert
5. **Graceful**: Shutdown an toàn, không mất jobs
6. **Scalable**: Có thể tăng số workers và queue size dễ dàng

## Performance

- **Throughput**: ~20 jobs/second (4 workers × 5 jobs/second/worker)
- **Latency**: ~150-200ms per job
- **Queue capacity**: 1000 jobs
- **Max processing time**: 5 seconds per job (timeout)
- **Retry attempts**: 3 times with exponential backoff
