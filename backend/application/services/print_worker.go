package services

import (
	"context"
	"fmt"
	"log"
	"math"
	"time"

	"cafe-pos/backend/domain/printing"
)

// PrintWorker defines the interface for background print job processing
type PrintWorker interface {
	// Start begins the worker goroutine
	Start(ctx context.Context)

	// Stop gracefully stops the worker
	Stop()

	// ProcessJob processes a single print job
	ProcessJob(ctx context.Context, job *printing.PrintJob) error
}

// printWorker implements the PrintWorker interface
type printWorker struct {
	printJobRepo         printing.PrintJobRepository
	printerConfigRepo    printing.PrinterConfigRepository
	printerManager       PrinterManager
	notificationService  PrintNotificationService
	stopChan             chan struct{}
	pollInterval         time.Duration
}

// PrintWorkerConfig contains configuration for the print worker
type PrintWorkerConfig struct {
	PrintJobRepo         printing.PrintJobRepository
	PrinterConfigRepo    printing.PrinterConfigRepository
	PrinterManager       PrinterManager
	NotificationService  PrintNotificationService
	PollInterval         time.Duration // Default: 10 seconds
}

// NewPrintWorker creates a new print worker instance
func NewPrintWorker(config PrintWorkerConfig) PrintWorker {
	pollInterval := config.PollInterval
	if pollInterval == 0 {
		pollInterval = 10 * time.Second
	}

	return &printWorker{
		printJobRepo:        config.PrintJobRepo,
		printerConfigRepo:   config.PrinterConfigRepo,
		printerManager:      config.PrinterManager,
		notificationService: config.NotificationService,
		stopChan:            make(chan struct{}),
		pollInterval:        pollInterval,
	}
}

// Start begins the worker goroutine that continuously processes the print queue
func (w *printWorker) Start(ctx context.Context) {
	log.Println("Print worker started")

	ticker := time.NewTicker(w.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Print worker stopped: context cancelled")
			return
		case <-w.stopChan:
			log.Println("Print worker stopped: stop signal received")
			return
		case <-ticker.C:
			// Process pending jobs
			if err := w.processPendingJobs(ctx); err != nil {
				log.Printf("Error processing pending jobs: %v", err)
			}
		}
	}
}

// Stop gracefully stops the worker
func (w *printWorker) Stop() {
	close(w.stopChan)
}

// processPendingJobs fetches and processes all pending print jobs
func (w *printWorker) processPendingJobs(ctx context.Context) error {
	// Fetch pending jobs (limit to 50 per batch)
	jobs, err := w.printJobRepo.FindPending(ctx, 50)
	if err != nil {
		return fmt.Errorf("failed to fetch pending jobs: %w", err)
	}

	if len(jobs) == 0 {
		return nil
	}

	log.Printf("Processing %d pending print jobs", len(jobs))

	// Process each job
	for _, job := range jobs {
		if err := w.ProcessJob(ctx, job); err != nil {
			log.Printf("Error processing job %s: %v", job.ID.Hex(), err)
			// Continue processing other jobs even if one fails
		}
	}

	return nil
}

// ProcessJob processes a single print job
func (w *printWorker) ProcessJob(ctx context.Context, job *printing.PrintJob) error {
	if job == nil {
		return fmt.Errorf("job cannot be nil")
	}

	// Only process pending jobs
	if job.Status != printing.PrintJobStatusPending {
		return fmt.Errorf("job %s is not pending (status: %s)", job.ID.Hex(), job.Status)
	}

	// Check if we should retry based on retry count
	if job.RetryCount > 0 {
		// Calculate backoff delay
		delay := w.calculateBackoffDelay(job.RetryCount)
		timeSinceUpdate := time.Since(job.UpdatedAt)

		if timeSinceUpdate < delay {
			// Not ready to retry yet
			log.Printf("Job %s not ready for retry (waiting %v more)", job.ID.Hex(), delay-timeSinceUpdate)
			return nil
		}
	}

	log.Printf("[PRINT] Starting print attempt - job_id=%s, order_id=%s, order_number=%s, job_type=%s, retry_count=%d, timestamp=%s",
		job.ID.Hex(), job.OrderID.Hex(), job.OrderNumber, job.Type, job.RetryCount, time.Now().Format(time.RFC3339))

	// Update status to PRINTING
	if err := w.printJobRepo.UpdateStatus(ctx, job.ID, printing.PrintJobStatusPrinting, ""); err != nil {
		return fmt.Errorf("failed to update job status to PRINTING: %w", err)
	}

	// Get printer configuration
	printerConfig, err := w.printerConfigRepo.FindByID(ctx, job.PrinterID)
	if err != nil {
		return w.handlePrintError(ctx, job, fmt.Errorf("failed to get printer config: %w", err))
	}
	if printerConfig == nil {
		return w.handlePrintError(ctx, job, fmt.Errorf("printer not found: %s", job.PrinterID.Hex()))
	}

	// Check if printer is enabled
	if !printerConfig.IsEnabled {
		return w.handlePrintError(ctx, job, fmt.Errorf("printer is disabled: %s", printerConfig.Name))
	}

	// Get printer instance
	printer, err := w.printerManager.GetPrinter(printerConfig)
	if err != nil {
		return w.handlePrintError(ctx, job, fmt.Errorf("failed to get printer instance: %w", err))
	}

	// Connect to printer
	if err := printer.Connect(); err != nil {
		log.Printf("[PRINT ERROR] Printer connection failed - job_id=%s, order_id=%s, printer=%s, connection_type=%s, address=%s:%d, error=%v, timestamp=%s",
			job.ID.Hex(), job.OrderID.Hex(), printerConfig.Name, printerConfig.ConnectionType, printerConfig.IPAddress, printerConfig.Port, err, time.Now().Format(time.RFC3339))
		return w.handlePrintError(ctx, job, fmt.Errorf("failed to connect to printer: %w", err))
	}
	defer func() {
		if err := printer.Disconnect(); err != nil {
			log.Printf("Warning: failed to disconnect printer for job %s: %v", job.ID.Hex(), err)
		}
	}()

	// Send print command
	if err := printer.Print(job.Content); err != nil {
		log.Printf("[PRINT ERROR] Print command failed - job_id=%s, order_id=%s, printer=%s, error=%v, timestamp=%s",
			job.ID.Hex(), job.OrderID.Hex(), printerConfig.Name, err, time.Now().Format(time.RFC3339))
		return w.handlePrintError(ctx, job, fmt.Errorf("failed to print: %w", err))
	}

	// Print successful - update status to COMPLETED
	now := time.Now()
	job.Status = printing.PrintJobStatusCompleted
	job.PrintedAt = &now
	job.UpdatedAt = now

	if err := w.printJobRepo.UpdateStatus(ctx, job.ID, printing.PrintJobStatusCompleted, ""); err != nil {
		log.Printf("Warning: failed to update job status to COMPLETED: %v", err)
		// Don't return error since print was successful
	}

	log.Printf("[PRINT SUCCESS] Print completed - job_id=%s, order_id=%s, order_number=%s, job_type=%s, printer=%s, timestamp=%s",
		job.ID.Hex(), job.OrderID.Hex(), job.OrderNumber, job.Type, printerConfig.Name, time.Now().Format(time.RFC3339))
	return nil
}

// handlePrintError handles print errors with retry logic
func (w *printWorker) handlePrintError(ctx context.Context, job *printing.PrintJob, err error) error {
	log.Printf("[PRINT ERROR] Print failed - job_id=%s, order_id=%s, order_number=%s, job_type=%s, retry_count=%d, error=%v, timestamp=%s",
		job.ID.Hex(), job.OrderID.Hex(), job.OrderNumber, job.Type, job.RetryCount, err, time.Now().Format(time.RFC3339))

	// Get printer config for notification
	printerConfig, _ := w.printerConfigRepo.FindByID(ctx, job.PrinterID)
	printerName := "Unknown"
	if printerConfig != nil {
		printerName = printerConfig.Name
	}

	// Increment retry count
	if err := w.printJobRepo.IncrementRetry(ctx, job.ID); err != nil {
		log.Printf("Warning: failed to increment retry count: %v", err)
	}

	// Check if we've exceeded max retries
	if job.RetryCount >= job.MaxRetries {
		// Mark as FAILED
		if err := w.printJobRepo.UpdateStatus(ctx, job.ID, printing.PrintJobStatusFailed, err.Error()); err != nil {
			log.Printf("Warning: failed to update job status to FAILED: %v", err)
		}
		
		// Send notification for permanent failure
		if w.notificationService != nil {
			w.notificationService.NotifyPrintFailure(job, printerName, err)
		}
		
		log.Printf("[PRINT FAILED] Job permanently failed - job_id=%s, order_id=%s, order_number=%s, job_type=%s, total_retries=%d, final_error=%v, timestamp=%s",
			job.ID.Hex(), job.OrderID.Hex(), job.OrderNumber, job.Type, job.RetryCount, err, time.Now().Format(time.RFC3339))
		return fmt.Errorf("job failed after %d retries: %w", job.RetryCount, err)
	}

	// Set back to PENDING for retry
	if err := w.printJobRepo.UpdateStatus(ctx, job.ID, printing.PrintJobStatusPending, err.Error()); err != nil {
		log.Printf("Warning: failed to update job status to PENDING: %v", err)
	}

	nextRetry := job.RetryCount + 1
	delay := w.calculateBackoffDelay(nextRetry)
	log.Printf("[PRINT RETRY] Job scheduled for retry - job_id=%s, order_id=%s, retry_attempt=%d/%d, next_retry_in=%v, timestamp=%s",
		job.ID.Hex(), job.OrderID.Hex(), nextRetry, job.MaxRetries, delay, time.Now().Format(time.RFC3339))

	return err
}

// calculateBackoffDelay calculates the exponential backoff delay for retries
// Retry 1: 30 seconds
// Retry 2: 1 minute
// Retry 3: 2 minutes
func (w *printWorker) calculateBackoffDelay(retryCount int) time.Duration {
	if retryCount <= 0 {
		return 0
	}

	// Base delay: 30 seconds
	baseDelay := 30 * time.Second

	// Exponential backoff: 30s * 2^(retryCount-1)
	// Retry 1: 30s * 2^0 = 30s
	// Retry 2: 30s * 2^1 = 60s (1m)
	// Retry 3: 30s * 2^2 = 120s (2m)
	multiplier := math.Pow(2, float64(retryCount-1))
	delay := time.Duration(float64(baseDelay) * multiplier)

	// Cap at 5 minutes
	maxDelay := 5 * time.Minute
	if delay > maxDelay {
		delay = maxDelay
	}

	return delay
}
