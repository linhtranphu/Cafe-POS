package services

import (
	"context"
	"sync"
	"time"

	"cafe-pos/backend/domain/batch"
)

// BatchAlertService handles batch alert checking and caching
type BatchAlertService struct {
	batchDefinitionRepo batch.BatchDefinitionRepository
	batchRecordRepo     batch.BatchRecordRepository
	alertCache          *AlertCache
	mu                  sync.RWMutex
}

// AlertCache stores alerts with TTL
type AlertCache struct {
	alerts      *batch.BatchAlerts
	lastChecked time.Time
	mu          sync.RWMutex
	ttl         time.Duration
}

// NewBatchAlertService creates a new batch alert service
func NewBatchAlertService(
	batchDefinitionRepo batch.BatchDefinitionRepository,
	batchRecordRepo batch.BatchRecordRepository,
) *BatchAlertService {
	return &BatchAlertService{
		batchDefinitionRepo: batchDefinitionRepo,
		batchRecordRepo:     batchRecordRepo,
		alertCache: &AlertCache{
			ttl: 5 * time.Minute,
		},
	}
}

// GetAlerts retrieves all alerts (with caching)
func (s *BatchAlertService) GetAlerts(ctx context.Context) (*batch.BatchAlerts, error) {
	// Try cache first
	if alerts, valid := s.alertCache.Get(); valid {
		return alerts, nil
	}

	// Cache miss or expired, check all alerts
	lowStock, err := s.CheckLowStock(ctx)
	if err != nil {
		return nil, err
	}

	expiring, err := s.CheckExpiring(ctx)
	if err != nil {
		return nil, err
	}

	expired, err := s.CheckExpired(ctx)
	if err != nil {
		return nil, err
	}

	alerts := &batch.BatchAlerts{
		LowStock:    lowStock,
		Expiring:    expiring,
		Expired:     expired,
		LastChecked: time.Now(),
	}

	// Store in cache
	s.alertCache.Set(alerts)

	return alerts, nil
}

// CheckLowStock checks for batches with low stock
func (s *BatchAlertService) CheckLowStock(ctx context.Context) ([]*batch.LowStockAlert, error) {
	// Get all batch definitions
	definitions, _, err := s.batchDefinitionRepo.FindAll(ctx, batch.BatchDefinitionFilter{})
	if err != nil {
		return nil, err
	}

	alerts := make([]*batch.LowStockAlert, 0)

	for _, def := range definitions {
		// Get total available quantity for this definition
		totalQty, err := s.batchRecordRepo.GetTotalAvailableQuantity(ctx, def.ID)
		if err != nil {
			continue // Skip on error
		}

		// Check if below threshold
		if totalQty <= def.LowStockThreshold {
			alerts = append(alerts, &batch.LowStockAlert{
				BatchDefinitionID: def.ID,
				BatchName:         def.Name,
				CurrentStock:      totalQty,
				Threshold:         def.LowStockThreshold,
				Unit:              def.Unit,
			})
		}
	}

	return alerts, nil
}

// CheckExpiring checks for batches expiring soon
func (s *BatchAlertService) CheckExpiring(ctx context.Context) ([]*batch.ExpiringAlert, error) {
	// Get all available batch records
	filter := batch.BatchRecordFilter{
		Status: batch.BatchStatusAvailable,
	}
	
	records, _, err := s.batchRecordRepo.FindAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	alerts := make([]*batch.ExpiringAlert, 0)
	now := time.Now()

	for _, record := range records {
		// Skip if already expired
		if now.After(record.ExpiresAt) {
			continue
		}

		// Calculate hours until expiry
		hoursUntilExpiry := int(record.ExpiresAt.Sub(now).Hours())

		// Get batch definition to check warning threshold
		def, err := s.batchDefinitionRepo.FindByID(ctx, record.BatchDefinitionID)
		if err != nil {
			continue // Skip on error
		}

		// Check if within warning period
		if hoursUntilExpiry <= def.ExpiryWarningHours {
			alerts = append(alerts, &batch.ExpiringAlert{
				BatchRecordID:     record.ID,
				BatchName:         record.BatchName,
				QuantityRemaining: record.QuantityRemaining,
				Unit:              record.Unit,
				ExpiresAt:         record.ExpiresAt,
				HoursUntilExpiry:  hoursUntilExpiry,
			})
		}
	}

	return alerts, nil
}

// CheckExpired checks for expired batches
func (s *BatchAlertService) CheckExpired(ctx context.Context) ([]*batch.ExpiredAlert, error) {
	// Get all expired batch records
	filter := batch.BatchRecordFilter{
		Status: batch.BatchStatusExpired,
	}
	
	records, _, err := s.batchRecordRepo.FindAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	alerts := make([]*batch.ExpiredAlert, 0)

	for _, record := range records {
		// Only alert if there's wasted quantity
		if record.QuantityRemaining > 0 {
			costWasted := record.QuantityRemaining * record.CostPerUnit

			alerts = append(alerts, &batch.ExpiredAlert{
				BatchRecordID:   record.ID,
				BatchName:       record.BatchName,
				QuantityWasted:  record.QuantityRemaining,
				Unit:            record.Unit,
				CostWasted:      costWasted,
				ExpiredAt:       record.ExpiresAt,
			})
		}
	}

	return alerts, nil
}

// InvalidateCache clears the alert cache
func (s *BatchAlertService) InvalidateCache() {
	s.alertCache.Clear()
}

// AlertCache methods

// Get retrieves alerts from cache if not expired
func (ac *AlertCache) Get() (*batch.BatchAlerts, bool) {
	ac.mu.RLock()
	defer ac.mu.RUnlock()

	if ac.alerts == nil {
		return nil, false
	}

	// Check if expired
	if time.Since(ac.lastChecked) > ac.ttl {
		return nil, false
	}

	return ac.alerts, true
}

// Set stores alerts in cache
func (ac *AlertCache) Set(alerts *batch.BatchAlerts) {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	ac.alerts = alerts
	ac.lastChecked = time.Now()
}

// Clear removes alerts from cache
func (ac *AlertCache) Clear() {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	ac.alerts = nil
	ac.lastChecked = time.Time{}
}
