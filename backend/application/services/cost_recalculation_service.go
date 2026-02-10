package services

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RecalculationStatus represents the status of background cost recalculation
type RecalculationStatus struct {
	InProgress     bool      `json:"in_progress"`
	QueuedItems    int       `json:"queued_items"`
	ProcessedItems int       `json:"processed_items"`
	FailedItems    int       `json:"failed_items"`
	LastUpdated    time.Time `json:"last_updated"`
}

// CostRecalculationService handles background cost recalculation jobs
type CostRecalculationService struct {
	costCalculator *CostCalculatorService
	menuRepo       MenuRepository
	
	// Worker pool configuration
	numWorkers     int
	workerTimeout  time.Duration
	maxRetries     int
	
	// Queue and status tracking
	recalcQueue    chan primitive.ObjectID
	queueSize      int
	
	// Status tracking with mutex for thread safety
	statusMu       sync.RWMutex
	status         *RecalculationStatus
	
	// Context for graceful shutdown
	ctx            context.Context
	cancel         context.CancelFunc
	wg             sync.WaitGroup
	
	// Monitoring service for metrics and alerts
	monitoringService *MonitoringService
}

// NewCostRecalculationService creates a new cost recalculation service
func NewCostRecalculationService(
	costCalculator *CostCalculatorService,
	menuRepo MenuRepository,
	numWorkers int,
	queueSize int,
) *CostRecalculationService {
	ctx, cancel := context.WithCancel(context.Background())
	
	return &CostRecalculationService{
		costCalculator: costCalculator,
		menuRepo:       menuRepo,
		numWorkers:     numWorkers,
		workerTimeout:  5 * time.Second, // 5 second timeout per item
		maxRetries:     3,
		recalcQueue:    make(chan primitive.ObjectID, queueSize),
		queueSize:      queueSize,
		status: &RecalculationStatus{
			InProgress:     false,
			QueuedItems:    0,
			ProcessedItems: 0,
			FailedItems:    0,
			LastUpdated:    time.Now(),
		},
		ctx:               ctx,
		cancel:            cancel,
		monitoringService: NewMonitoringService(),
	}
}

// SetMonitoringService sets the monitoring service for metrics collection
func (s *CostRecalculationService) SetMonitoringService(monitoringService *MonitoringService) {
	s.monitoringService = monitoringService
}

// Start starts the worker pool for processing recalculation jobs
// This should be called once when the service starts up
func (s *CostRecalculationService) Start() {
	s.updateStatus(func(status *RecalculationStatus) {
		status.InProgress = true
		status.LastUpdated = time.Now()
	})
	
	// Start worker goroutines
	for i := 0; i < s.numWorkers; i++ {
		s.wg.Add(1)
		go s.worker(i)
	}
}

// Stop gracefully stops the worker pool
func (s *CostRecalculationService) Stop() {
	// Cancel context to signal workers to stop
	s.cancel()
	
	// Wait for all workers to finish
	s.wg.Wait()
	
	// Close the queue channel
	close(s.recalcQueue)
	
	s.updateStatus(func(status *RecalculationStatus) {
		status.InProgress = false
		status.LastUpdated = time.Now()
	})
}

// worker is a background worker that processes recalculation jobs from the queue
func (s *CostRecalculationService) worker(workerID int) {
	defer s.wg.Done()
	
	for {
		select {
		case <-s.ctx.Done():
			// Context cancelled, stop worker
			return
			
		case menuItemID, ok := <-s.recalcQueue:
			if !ok {
				// Channel closed, stop worker
				return
			}
			
			// Process the recalculation job with retries
			err := s.processRecalculationWithRetry(menuItemID)
			
			// Update status
			s.updateStatus(func(status *RecalculationStatus) {
				status.QueuedItems = len(s.recalcQueue)
				if err != nil {
					status.FailedItems++
				} else {
					status.ProcessedItems++
				}
				status.LastUpdated = time.Now()
			})
		}
	}
}

// processRecalculationWithRetry processes a single recalculation job with retry logic
func (s *CostRecalculationService) processRecalculationWithRetry(menuItemID primitive.ObjectID) error {
	var lastErr error
	startTime := time.Now()
	
	for attempt := 0; attempt < s.maxRetries; attempt++ {
		// Create context with timeout
		ctx, cancel := context.WithTimeout(s.ctx, s.workerTimeout)
		
		// Attempt to recalculate cost
		err := s.processRecalculation(ctx, menuItemID)
		cancel()
		
		if err == nil {
			// Success - record metric
			duration := time.Since(startTime)
			if s.monitoringService != nil {
				s.monitoringService.RecordMetric(Metric{
					Type:      MetricTypeRecalculationJob,
					Status:    MetricStatusSuccess,
					Timestamp: time.Now(),
					Duration:  duration,
					Message:   "Cost recalculation completed successfully",
					Metadata: map[string]interface{}{
						"menu_item_id": menuItemID.Hex(),
						"attempts":     attempt + 1,
					},
				})
			}
			return nil
		}
		
		lastErr = err
		
		// Exponential backoff before retry
		if attempt < s.maxRetries-1 {
			backoffDuration := time.Duration(1<<uint(attempt)) * time.Second
			select {
			case <-time.After(backoffDuration):
				// Continue to next retry
			case <-s.ctx.Done():
				// Context cancelled, stop retrying
				return s.ctx.Err()
			}
		}
	}
	
	// All retries failed - record metric
	duration := time.Since(startTime)
	if s.monitoringService != nil {
		s.monitoringService.RecordMetric(Metric{
			Type:      MetricTypeRecalculationJob,
			Status:    MetricStatusFailure,
			Timestamp: time.Now(),
			Duration:  duration,
			Message:   fmt.Sprintf("Cost recalculation failed after %d retries: %v", s.maxRetries, lastErr),
			Metadata: map[string]interface{}{
				"menu_item_id": menuItemID.Hex(),
				"attempts":     s.maxRetries,
				"error":        lastErr.Error(),
			},
		})
	}
	
	return fmt.Errorf("failed after %d retries: %w", s.maxRetries, lastErr)
}

// processRecalculation processes a single recalculation job
func (s *CostRecalculationService) processRecalculation(ctx context.Context, menuItemID primitive.ObjectID) error {
	// Calculate new cost
	costResult, err := s.costCalculator.CalculateMenuItemCost(ctx, menuItemID)
	if err != nil {
		return fmt.Errorf("failed to calculate cost: %w", err)
	}
	
	// Fetch menu item
	menuItem, err := s.menuRepo.FindByID(ctx, menuItemID)
	if err != nil {
		return fmt.Errorf("failed to fetch menu item: %w", err)
	}
	
	// Update menu item with new cost
	menuItem.CurrentCost = costResult.CurrentCost
	menuItem.CostStatus = costResult.CostStatus
	menuItem.CostLastCalculatedAt = costResult.CostLastCalculatedAt
	
	// Save to database
	err = s.menuRepo.Update(ctx, menuItemID, menuItem)
	if err != nil {
		return fmt.Errorf("failed to update menu item: %w", err)
	}
	
	return nil
}

// ProcessRecalculationQueue processes all pending recalculation jobs in the queue
// This method blocks until all jobs are processed or context is cancelled
// Requirements: 9.2, 9.3
func (s *CostRecalculationService) ProcessRecalculationQueue(ctx context.Context) error {
	// Get current queue size
	queueSize := len(s.recalcQueue)
	
	if queueSize == 0 {
		return nil
	}
	
	// Wait for all queued items to be processed
	// We'll poll the queue size until it's empty or context is cancelled
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if len(s.recalcQueue) == 0 {
				// Queue is empty, all jobs processed
				return nil
			}
		}
	}
}

// QueueRecalculation queues a menu item for cost recalculation
// Returns error if queue is full
func (s *CostRecalculationService) QueueRecalculation(menuItemID primitive.ObjectID) error {
	select {
	case s.recalcQueue <- menuItemID:
		// Successfully queued
		s.updateStatus(func(status *RecalculationStatus) {
			status.QueuedItems = len(s.recalcQueue)
			status.LastUpdated = time.Now()
		})
		return nil
	default:
		// Queue is full
		return fmt.Errorf("recalculation queue is full (size: %d)", s.queueSize)
	}
}

// GetRecalculationStatus returns the current status of the recalculation service
// Requirements: 9.4
func (s *CostRecalculationService) GetRecalculationStatus(ctx context.Context) (*RecalculationStatus, error) {
	s.statusMu.RLock()
	defer s.statusMu.RUnlock()
	
	// Create a copy to avoid race conditions
	statusCopy := &RecalculationStatus{
		InProgress:     s.status.InProgress,
		QueuedItems:    len(s.recalcQueue), // Get real-time queue size
		ProcessedItems: s.status.ProcessedItems,
		FailedItems:    s.status.FailedItems,
		LastUpdated:    s.status.LastUpdated,
	}
	
	return statusCopy, nil
}

// updateStatus is a helper method to safely update the status with a mutex
func (s *CostRecalculationService) updateStatus(updateFn func(*RecalculationStatus)) {
	s.statusMu.Lock()
	defer s.statusMu.Unlock()
	updateFn(s.status)
}

// ResetStatistics resets the processed and failed item counters
func (s *CostRecalculationService) ResetStatistics() {
	s.updateStatus(func(status *RecalculationStatus) {
		status.ProcessedItems = 0
		status.FailedItems = 0
		status.LastUpdated = time.Now()
	})
}
