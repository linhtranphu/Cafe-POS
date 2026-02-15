package services

import (
	"context"
	"time"

	"cafe-pos/backend/domain/batch"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BatchReportService provides reporting functionality for batch operations
type BatchReportService struct {
	batchRecordRepo   batch.BatchRecordRepository
	batchUsageLogRepo batch.BatchUsageLogRepository
	batchDefRepo      batch.BatchDefinitionRepository
}

// NewBatchReportService creates a new batch report service
func NewBatchReportService(
	batchRecordRepo batch.BatchRecordRepository,
	batchUsageLogRepo batch.BatchUsageLogRepository,
	batchDefRepo batch.BatchDefinitionRepository,
) *BatchReportService {
	return &BatchReportService{
		batchRecordRepo:   batchRecordRepo,
		batchUsageLogRepo: batchUsageLogRepo,
		batchDefRepo:      batchDefRepo,
	}
}

// ProductionReportParams contains parameters for production report
type ProductionReportParams struct {
	FromDate          time.Time
	ToDate            time.Time
	BatchDefinitionID *primitive.ObjectID
	PreparedBy        string
}

// ProductionReport contains production statistics
type ProductionReport struct {
	TotalBatchesProduced int                         `json:"total_batches_produced"`
	TotalQuantityProduced float64                    `json:"total_quantity_produced"`
	TotalCost            float64                     `json:"total_cost"`
	ByBatchType          []ProductionByBatchType     `json:"by_batch_type"`
	ByPreparer           []ProductionByPreparer      `json:"by_preparer"`
}

// ProductionByBatchType contains production stats by batch type
type ProductionByBatchType struct {
	BatchName     string  `json:"batch_name"`
	Count         int     `json:"count"`
	TotalQuantity float64 `json:"total_quantity"`
	TotalCost     float64 `json:"total_cost"`
}

// ProductionByPreparer contains production stats by preparer
type ProductionByPreparer struct {
	Preparer      string  `json:"preparer"`
	Count         int     `json:"count"`
	TotalQuantity float64 `json:"total_quantity"`
}

// GetProductionReport generates a production report for the specified time period
func (s *BatchReportService) GetProductionReport(ctx context.Context, params ProductionReportParams) (*ProductionReport, error) {
	// Build filter for batch records
	filter := batch.BatchRecordFilter{
		FromDate:          &params.FromDate,
		ToDate:            &params.ToDate,
		BatchDefinitionID: params.BatchDefinitionID,
		PreparedBy:        params.PreparedBy,
	}

	// Get all batch records in the time period
	records, _, err := s.batchRecordRepo.FindAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Initialize report
	report := &ProductionReport{
		ByBatchType: []ProductionByBatchType{},
		ByPreparer:  []ProductionByPreparer{},
	}

	// Maps for aggregation
	byBatchType := make(map[string]*ProductionByBatchType)
	byPreparer := make(map[string]*ProductionByPreparer)

	// Aggregate data
	for _, record := range records {
		// Overall totals
		report.TotalBatchesProduced++
		report.TotalQuantityProduced += record.QuantityProduced
		report.TotalCost += record.TotalCost

		// By batch type
		if _, exists := byBatchType[record.BatchName]; !exists {
			byBatchType[record.BatchName] = &ProductionByBatchType{
				BatchName: record.BatchName,
			}
		}
		byBatchType[record.BatchName].Count++
		byBatchType[record.BatchName].TotalQuantity += record.QuantityProduced
		byBatchType[record.BatchName].TotalCost += record.TotalCost

		// By preparer
		if _, exists := byPreparer[record.PreparedBy]; !exists {
			byPreparer[record.PreparedBy] = &ProductionByPreparer{
				Preparer: record.PreparedBy,
			}
		}
		byPreparer[record.PreparedBy].Count++
		byPreparer[record.PreparedBy].TotalQuantity += record.QuantityProduced
	}

	// Convert maps to slices
	for _, v := range byBatchType {
		report.ByBatchType = append(report.ByBatchType, *v)
	}
	for _, v := range byPreparer {
		report.ByPreparer = append(report.ByPreparer, *v)
	}

	return report, nil
}

// WastageReportParams contains parameters for wastage report
type WastageReportParams struct {
	FromDate          time.Time
	ToDate            time.Time
	BatchDefinitionID *primitive.ObjectID
}

// WastageReport contains wastage statistics
type WastageReport struct {
	TotalExpiredBatches int                   `json:"total_expired_batches"`
	TotalQuantityWasted float64               `json:"total_quantity_wasted"`
	TotalCostWasted     float64               `json:"total_cost_wasted"`
	WastageByType       []WastageByBatchType  `json:"wastage_by_type"`
}

// WastageByBatchType contains wastage stats by batch type
type WastageByBatchType struct {
	BatchName       string  `json:"batch_name"`
	ExpiredCount    int     `json:"expired_count"`
	QuantityWasted  float64 `json:"quantity_wasted"`
	CostWasted      float64 `json:"cost_wasted"`
}

// GetWastageReport generates a wastage report for the specified time period
func (s *BatchReportService) GetWastageReport(ctx context.Context, params WastageReportParams) (*WastageReport, error) {
	// Build filter for expired batch records
	filter := batch.BatchRecordFilter{
		FromDate:          &params.FromDate,
		ToDate:            &params.ToDate,
		BatchDefinitionID: params.BatchDefinitionID,
		Status:            batch.BatchStatusExpired,
	}

	// Get all expired batch records in the time period
	records, _, err := s.batchRecordRepo.FindAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Initialize report
	report := &WastageReport{
		WastageByType: []WastageByBatchType{},
	}

	// Map for aggregation
	byBatchType := make(map[string]*WastageByBatchType)

	// Aggregate data
	for _, record := range records {
		// Calculate wasted quantity (quantity remaining when expired)
		wastedQuantity := record.QuantityRemaining
		wastedCost := wastedQuantity * record.CostPerUnit

		// Overall totals
		report.TotalExpiredBatches++
		report.TotalQuantityWasted += wastedQuantity
		report.TotalCostWasted += wastedCost

		// By batch type
		if _, exists := byBatchType[record.BatchName]; !exists {
			byBatchType[record.BatchName] = &WastageByBatchType{
				BatchName: record.BatchName,
			}
		}
		byBatchType[record.BatchName].ExpiredCount++
		byBatchType[record.BatchName].QuantityWasted += wastedQuantity
		byBatchType[record.BatchName].CostWasted += wastedCost
	}

	// Convert map to slice
	for _, v := range byBatchType {
		report.WastageByType = append(report.WastageByType, *v)
	}

	return report, nil
}

// UsageReportParams contains parameters for usage report
type UsageReportParams struct {
	FromDate          time.Time
	ToDate            time.Time
	BatchDefinitionID *primitive.ObjectID
	MenuItemID        *primitive.ObjectID
}

// UsageReport contains usage statistics
type UsageReport struct {
	TotalUsageCount  int                `json:"total_usage_count"`
	TotalQuantityUsed float64           `json:"total_quantity_used"`
	TotalCost        float64            `json:"total_cost"`
	ByMenuItem       []UsageByMenuItem  `json:"by_menu_item"`
	UsageTrend       []UsageTrendPoint  `json:"usage_trend"`
}

// UsageByMenuItem contains usage stats by menu item
type UsageByMenuItem struct {
	MenuItemName  string  `json:"menu_item_name"`
	UsageCount    int     `json:"usage_count"`
	QuantityUsed  float64 `json:"quantity_used"`
	Cost          float64 `json:"cost"`
}

// UsageTrendPoint represents usage on a specific date
type UsageTrendPoint struct {
	Date         string  `json:"date"`
	QuantityUsed float64 `json:"quantity_used"`
	Cost         float64 `json:"cost"`
}

// GetUsageReport generates a usage report for the specified time period
func (s *BatchReportService) GetUsageReport(ctx context.Context, params UsageReportParams) (*UsageReport, error) {
	// Build filter for usage logs
	filter := batch.BatchUsageLogFilter{
		FromDate:   &params.FromDate,
		ToDate:     &params.ToDate,
		MenuItemID: params.MenuItemID,
	}

	// If batch definition ID is provided, we need to filter by batch records
	// For now, we'll get all usage logs and filter later if needed
	// TODO: Optimize this with aggregation pipeline if performance becomes an issue

	// Get all usage logs in the time period
	logs, _, err := s.batchUsageLogRepo.FindAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	// If batch definition ID is provided, filter logs by batch records
	if params.BatchDefinitionID != nil {
		// Get all batch records for this definition
		recordFilter := batch.BatchRecordFilter{
			BatchDefinitionID: params.BatchDefinitionID,
		}
		records, _, err := s.batchRecordRepo.FindAll(ctx, recordFilter)
		if err != nil {
			return nil, err
		}

		// Create a set of batch record IDs
		recordIDs := make(map[primitive.ObjectID]bool)
		for _, record := range records {
			recordIDs[record.ID] = true
		}

		// Filter logs to only include those from these batch records
		filteredLogs := []*batch.BatchUsageLog{}
		for _, log := range logs {
			if recordIDs[log.BatchRecordID] {
				filteredLogs = append(filteredLogs, log)
			}
		}
		logs = filteredLogs
	}

	// Initialize report
	report := &UsageReport{
		ByMenuItem: []UsageByMenuItem{},
		UsageTrend: []UsageTrendPoint{},
	}

	// Maps for aggregation
	byMenuItem := make(map[string]*UsageByMenuItem)
	byDate := make(map[string]*UsageTrendPoint)

	// Aggregate data
	for _, log := range logs {
		// Overall totals
		report.TotalUsageCount++
		report.TotalQuantityUsed += log.QuantityUsed
		report.TotalCost += log.TotalCost

		// By menu item
		if _, exists := byMenuItem[log.MenuItemName]; !exists {
			byMenuItem[log.MenuItemName] = &UsageByMenuItem{
				MenuItemName: log.MenuItemName,
			}
		}
		byMenuItem[log.MenuItemName].UsageCount++
		byMenuItem[log.MenuItemName].QuantityUsed += log.QuantityUsed
		byMenuItem[log.MenuItemName].Cost += log.TotalCost

		// By date (for trend)
		dateStr := log.UsedAt.Format("2006-01-02")
		if _, exists := byDate[dateStr]; !exists {
			byDate[dateStr] = &UsageTrendPoint{
				Date: dateStr,
			}
		}
		byDate[dateStr].QuantityUsed += log.QuantityUsed
		byDate[dateStr].Cost += log.TotalCost
	}

	// Convert maps to slices
	for _, v := range byMenuItem {
		report.ByMenuItem = append(report.ByMenuItem, *v)
	}
	for _, v := range byDate {
		report.UsageTrend = append(report.UsageTrend, *v)
	}

	return report, nil
}
