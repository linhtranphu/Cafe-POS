package services

import (
	"context"
	"log"
	"time"

	"cafe-pos/backend/domain/batch"
)

// BatchBackgroundJobs handles scheduled background jobs for batch management
type BatchBackgroundJobs struct {
	batchRecordService *BatchRecordService
	batchAlertService  *BatchAlertService
	stopChan           chan struct{}
	running            bool
}

// NewBatchBackgroundJobs creates a new batch background jobs manager
func NewBatchBackgroundJobs(
	batchRecordService *BatchRecordService,
	batchAlertService *BatchAlertService,
) *BatchBackgroundJobs {
	return &BatchBackgroundJobs{
		batchRecordService: batchRecordService,
		batchAlertService:  batchAlertService,
		stopChan:           make(chan struct{}),
		running:            false,
	}
}

// Start starts all background jobs
func (j *BatchBackgroundJobs) Start(ctx context.Context) {
	if j.running {
		return
	}

	j.running = true
	log.Println("Starting batch background jobs...")

	// Start expired batch marker job (runs every hour)
	go j.runExpiredBatchMarker(ctx)

	// Start alert checker job (runs every 5 minutes)
	go j.runAlertChecker(ctx)

	log.Println("Batch background jobs started")
}

// Stop stops all background jobs
func (j *BatchBackgroundJobs) Stop() {
	if !j.running {
		return
	}

	log.Println("Stopping batch background jobs...")
	close(j.stopChan)
	j.running = false
	log.Println("Batch background jobs stopped")
}

// runExpiredBatchMarker runs the expired batch marker job every hour
func (j *BatchBackgroundJobs) runExpiredBatchMarker(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	// Run immediately on start
	j.markExpiredBatches(ctx)

	for {
		select {
		case <-ticker.C:
			j.markExpiredBatches(ctx)
		case <-j.stopChan:
			return
		case <-ctx.Done():
			return
		}
	}
}

// runAlertChecker runs the alert checker job every 5 minutes
func (j *BatchBackgroundJobs) runAlertChecker(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	// Run immediately on start
	j.checkAlerts(ctx)

	for {
		select {
		case <-ticker.C:
			j.checkAlerts(ctx)
		case <-j.stopChan:
			return
		case <-ctx.Done():
			return
		}
	}
}

// markExpiredBatches marks all expired batches
func (j *BatchBackgroundJobs) markExpiredBatches(ctx context.Context) {
	log.Println("Running expired batch marker job...")
	
	err := j.batchRecordService.MarkExpiredBatches(ctx)
	if err != nil {
		log.Printf("Error marking expired batches: %v", err)
		return
	}

	log.Println("Expired batch marker job completed")
}

// checkAlerts checks and caches all alerts
func (j *BatchBackgroundJobs) checkAlerts(ctx context.Context) {
	log.Println("Running alert checker job...")
	
	// This will check all alerts and update the cache
	alerts, err := j.batchAlertService.GetAlerts(ctx)
	if err != nil {
		log.Printf("Error checking alerts: %v", err)
		return
	}

	// Log alert counts
	log.Printf("Alert check completed: %d low stock, %d expiring, %d expired",
		len(alerts.LowStock), len(alerts.Expiring), len(alerts.Expired))
}

// RunExpiredBatchMarkerNow runs the expired batch marker job immediately (for testing/manual trigger)
func (j *BatchBackgroundJobs) RunExpiredBatchMarkerNow(ctx context.Context) error {
	return j.batchRecordService.MarkExpiredBatches(ctx)
}

// RunAlertCheckerNow runs the alert checker job immediately (for testing/manual trigger)
func (j *BatchBackgroundJobs) RunAlertCheckerNow(ctx context.Context) (*batch.BatchAlerts, error) {
	return j.batchAlertService.GetAlerts(ctx)
}

// BatchAlerts is an alias for batch.BatchAlerts for convenience
type BatchAlerts struct {
	LowStock    int
	Expiring    int
	Expired     int
	LastChecked time.Time
}
