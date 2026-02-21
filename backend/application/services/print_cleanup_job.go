package services

import (
	"context"
	"log"
	"time"

	"cafe-pos/backend/domain/printing"
)

// PrintCleanupJob handles periodic cleanup of old print jobs and notifications
type PrintCleanupJob interface {
	// Start begins the cleanup job scheduler
	Start(ctx context.Context)

	// Stop gracefully stops the cleanup job
	Stop()

	// RunCleanup executes a single cleanup cycle
	RunCleanup(ctx context.Context) error
}

type printCleanupJob struct {
	printJobRepo     printing.PrintJobRepository
	notificationRepo printing.PrintNotificationRepository
	stopChan         chan struct{}
	interval         time.Duration
	retentionDays    int
}

// PrintCleanupJobConfig contains configuration for the cleanup job
type PrintCleanupJobConfig struct {
	PrintJobRepo     printing.PrintJobRepository
	NotificationRepo printing.PrintNotificationRepository
	Interval         time.Duration // Default: 24 hours
	RetentionDays    int           // Default: 7 days
}

// NewPrintCleanupJob creates a new cleanup job instance
func NewPrintCleanupJob(config PrintCleanupJobConfig) PrintCleanupJob {
	interval := config.Interval
	if interval == 0 {
		interval = 24 * time.Hour // Run daily by default
	}

	retentionDays := config.RetentionDays
	if retentionDays == 0 {
		retentionDays = 7 // Keep for 7 days by default
	}

	return &printCleanupJob{
		printJobRepo:     config.PrintJobRepo,
		notificationRepo: config.NotificationRepo,
		stopChan:         make(chan struct{}),
		interval:         interval,
		retentionDays:    retentionDays,
	}
}

// Start begins the cleanup job scheduler
func (j *printCleanupJob) Start(ctx context.Context) {
	log.Printf("[CLEANUP] Print cleanup job started - interval=%v, retention=%d days", j.interval, j.retentionDays)

	// Run cleanup immediately on start
	if err := j.RunCleanup(ctx); err != nil {
		log.Printf("[CLEANUP ERROR] Initial cleanup failed: %v", err)
	}

	ticker := time.NewTicker(j.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("[CLEANUP] Print cleanup job stopped: context cancelled")
			return
		case <-j.stopChan:
			log.Println("[CLEANUP] Print cleanup job stopped: stop signal received")
			return
		case <-ticker.C:
			if err := j.RunCleanup(ctx); err != nil {
				log.Printf("[CLEANUP ERROR] Cleanup failed: %v", err)
			}
		}
	}
}

// Stop gracefully stops the cleanup job
func (j *printCleanupJob) Stop() {
	close(j.stopChan)
}

// RunCleanup executes a single cleanup cycle
func (j *printCleanupJob) RunCleanup(ctx context.Context) error {
	startTime := time.Now()
	log.Printf("[CLEANUP] Starting cleanup cycle - timestamp=%s", startTime.Format(time.RFC3339))

	// Calculate cutoff time
	cutoffTime := time.Now().AddDate(0, 0, -j.retentionDays)
	log.Printf("[CLEANUP] Deleting completed print jobs older than %s", cutoffTime.Format(time.RFC3339))

	// Delete old completed print jobs
	if err := j.printJobRepo.DeleteOldCompleted(ctx, cutoffTime); err != nil {
		log.Printf("[CLEANUP ERROR] Failed to delete old print jobs: %v", err)
		return err
	}

	// Delete old read notifications (keep for 30 days)
	notificationCutoff := time.Now().AddDate(0, 0, -30)
	if j.notificationRepo != nil {
		if err := j.notificationRepo.DeleteOld(notificationCutoff); err != nil {
			log.Printf("[CLEANUP ERROR] Failed to delete old notifications: %v", err)
			// Don't return error, continue with cleanup
		}
	}

	duration := time.Since(startTime)
	log.Printf("[CLEANUP SUCCESS] Cleanup completed - duration=%v, timestamp=%s", duration, time.Now().Format(time.RFC3339))

	return nil
}
