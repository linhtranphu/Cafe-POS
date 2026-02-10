package services

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// MetricType represents the type of metric being collected
type MetricType string

const (
	MetricTypeCostCalculation     MetricType = "cost_calculation"
	MetricTypeRecalculationJob    MetricType = "recalculation_job"
	MetricTypeShiftClosure        MetricType = "shift_closure"
	MetricTypeProfitAnalysis      MetricType = "profit_analysis"
	MetricTypeOperatingExpense    MetricType = "operating_expense"
)

// MetricStatus represents the status of a metric event
type MetricStatus string

const (
	MetricStatusSuccess MetricStatus = "success"
	MetricStatusFailure MetricStatus = "failure"
	MetricStatusWarning MetricStatus = "warning"
)

// Metric represents a single metric event
type Metric struct {
	Type      MetricType
	Status    MetricStatus
	Timestamp time.Time
	Duration  time.Duration
	Message   string
	Metadata  map[string]interface{}
}

// AlertLevel represents the severity of an alert
type AlertLevel string

const (
	AlertLevelInfo     AlertLevel = "info"
	AlertLevelWarning  AlertLevel = "warning"
	AlertLevelCritical AlertLevel = "critical"
)

// Alert represents a monitoring alert
type Alert struct {
	Level     AlertLevel
	Type      string
	Message   string
	Timestamp time.Time
	Metadata  map[string]interface{}
}

// AlertRule defines conditions for triggering alerts
type AlertRule struct {
	Name              string
	Type              MetricType
	FailureThreshold  int           // Number of failures before alerting
	TimeWindow        time.Duration // Time window to check failures
	ErrorRateThreshold float64      // Error rate percentage (0-100)
	Enabled           bool
}

// MonitoringMetrics holds aggregated metrics
type MonitoringMetrics struct {
	// Cost calculation metrics
	TotalCostCalculations      int64
	SuccessfulCostCalculations int64
	FailedCostCalculations     int64
	AverageCostCalcDuration    time.Duration
	
	// Recalculation job metrics
	TotalRecalcJobs            int64
	SuccessfulRecalcJobs       int64
	FailedRecalcJobs           int64
	AverageRecalcDuration      time.Duration
	
	// Shift closure metrics
	TotalShiftClosures         int64
	SuccessfulShiftClosures    int64
	FailedShiftClosures        int64
	AverageShiftClosureDuration time.Duration
	
	// Profit analysis metrics
	TotalProfitAnalyses        int64
	SuccessfulProfitAnalyses   int64
	FailedProfitAnalyses       int64
	AverageProfitAnalysisDuration time.Duration
	
	// Operating expense metrics
	TotalExpenseOperations     int64
	SuccessfulExpenseOperations int64
	FailedExpenseOperations    int64
	
	// Alert metrics
	TotalAlerts                int64
	CriticalAlerts             int64
	WarningAlerts              int64
	
	LastUpdated                time.Time
}

// MonitoringService handles metrics collection and alerting
type MonitoringService struct {
	mu              sync.RWMutex
	metrics         []Metric
	alerts          []Alert
	aggregatedMetrics *MonitoringMetrics
	alertRules      []AlertRule
	maxMetricsSize  int
	maxAlertsSize   int
}

// NewMonitoringService creates a new monitoring service
func NewMonitoringService() *MonitoringService {
	return &MonitoringService{
		metrics:           make([]Metric, 0),
		alerts:            make([]Alert, 0),
		aggregatedMetrics: &MonitoringMetrics{
			LastUpdated: time.Now(),
		},
		alertRules:     getDefaultAlertRules(),
		maxMetricsSize: 10000, // Keep last 10k metrics
		maxAlertsSize:  1000,  // Keep last 1k alerts
	}
}

// getDefaultAlertRules returns the default alert rules
func getDefaultAlertRules() []AlertRule {
	return []AlertRule{
		{
			Name:              "High Cost Calculation Failure Rate",
			Type:              MetricTypeCostCalculation,
			FailureThreshold:  10,
			TimeWindow:        5 * time.Minute,
			ErrorRateThreshold: 20.0, // Alert if >20% error rate
			Enabled:           true,
		},
		{
			Name:              "High Recalculation Job Failure Rate",
			Type:              MetricTypeRecalculationJob,
			FailureThreshold:  5,
			TimeWindow:        5 * time.Minute,
			ErrorRateThreshold: 15.0, // Alert if >15% error rate
			Enabled:           true,
		},
		{
			Name:              "Shift Closure Failures",
			Type:              MetricTypeShiftClosure,
			FailureThreshold:  3,
			TimeWindow:        10 * time.Minute,
			ErrorRateThreshold: 10.0, // Alert if >10% error rate
			Enabled:           true,
		},
		{
			Name:              "Profit Analysis Failures",
			Type:              MetricTypeProfitAnalysis,
			FailureThreshold:  10,
			TimeWindow:        5 * time.Minute,
			ErrorRateThreshold: 25.0, // Alert if >25% error rate
			Enabled:           true,
		},
	}
}

// RecordMetric records a new metric event
func (s *MonitoringService) RecordMetric(metric Metric) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// Set timestamp if not provided
	if metric.Timestamp.IsZero() {
		metric.Timestamp = time.Now()
	}
	
	// Add metric to the list
	s.metrics = append(s.metrics, metric)
	
	// Trim metrics if exceeding max size
	if len(s.metrics) > s.maxMetricsSize {
		s.metrics = s.metrics[len(s.metrics)-s.maxMetricsSize:]
	}
	
	// Update aggregated metrics
	s.updateAggregatedMetrics(metric)
	
	// Check alert rules
	s.checkAlertRules(metric)
}

// updateAggregatedMetrics updates the aggregated metrics based on a new metric
func (s *MonitoringService) updateAggregatedMetrics(metric Metric) {
	switch metric.Type {
	case MetricTypeCostCalculation:
		s.aggregatedMetrics.TotalCostCalculations++
		if metric.Status == MetricStatusSuccess {
			s.aggregatedMetrics.SuccessfulCostCalculations++
		} else if metric.Status == MetricStatusFailure {
			s.aggregatedMetrics.FailedCostCalculations++
		}
		// Update average duration
		if metric.Duration > 0 {
			total := s.aggregatedMetrics.TotalCostCalculations
			avg := s.aggregatedMetrics.AverageCostCalcDuration
			s.aggregatedMetrics.AverageCostCalcDuration = time.Duration(
				(int64(avg)*int64(total-1) + int64(metric.Duration)) / int64(total),
			)
		}
		
	case MetricTypeRecalculationJob:
		s.aggregatedMetrics.TotalRecalcJobs++
		if metric.Status == MetricStatusSuccess {
			s.aggregatedMetrics.SuccessfulRecalcJobs++
		} else if metric.Status == MetricStatusFailure {
			s.aggregatedMetrics.FailedRecalcJobs++
		}
		// Update average duration
		if metric.Duration > 0 {
			total := s.aggregatedMetrics.TotalRecalcJobs
			avg := s.aggregatedMetrics.AverageRecalcDuration
			s.aggregatedMetrics.AverageRecalcDuration = time.Duration(
				(int64(avg)*int64(total-1) + int64(metric.Duration)) / int64(total),
			)
		}
		
	case MetricTypeShiftClosure:
		s.aggregatedMetrics.TotalShiftClosures++
		if metric.Status == MetricStatusSuccess {
			s.aggregatedMetrics.SuccessfulShiftClosures++
		} else if metric.Status == MetricStatusFailure {
			s.aggregatedMetrics.FailedShiftClosures++
		}
		// Update average duration
		if metric.Duration > 0 {
			total := s.aggregatedMetrics.TotalShiftClosures
			avg := s.aggregatedMetrics.AverageShiftClosureDuration
			s.aggregatedMetrics.AverageShiftClosureDuration = time.Duration(
				(int64(avg)*int64(total-1) + int64(metric.Duration)) / int64(total),
			)
		}
		
	case MetricTypeProfitAnalysis:
		s.aggregatedMetrics.TotalProfitAnalyses++
		if metric.Status == MetricStatusSuccess {
			s.aggregatedMetrics.SuccessfulProfitAnalyses++
		} else if metric.Status == MetricStatusFailure {
			s.aggregatedMetrics.FailedProfitAnalyses++
		}
		// Update average duration
		if metric.Duration > 0 {
			total := s.aggregatedMetrics.TotalProfitAnalyses
			avg := s.aggregatedMetrics.AverageProfitAnalysisDuration
			s.aggregatedMetrics.AverageProfitAnalysisDuration = time.Duration(
				(int64(avg)*int64(total-1) + int64(metric.Duration)) / int64(total),
			)
		}
		
	case MetricTypeOperatingExpense:
		s.aggregatedMetrics.TotalExpenseOperations++
		if metric.Status == MetricStatusSuccess {
			s.aggregatedMetrics.SuccessfulExpenseOperations++
		} else if metric.Status == MetricStatusFailure {
			s.aggregatedMetrics.FailedExpenseOperations++
		}
	}
	
	s.aggregatedMetrics.LastUpdated = time.Now()
}

// checkAlertRules checks if any alert rules are triggered by the new metric
func (s *MonitoringService) checkAlertRules(metric Metric) {
	// Only check for failures
	if metric.Status != MetricStatusFailure {
		return
	}
	
	for _, rule := range s.alertRules {
		if !rule.Enabled || rule.Type != metric.Type {
			continue
		}
		
		// Count failures in the time window
		failures := s.countMetricsInWindow(rule.Type, MetricStatusFailure, rule.TimeWindow)
		total := s.countMetricsInWindow(rule.Type, "", rule.TimeWindow)
		
		// Check failure threshold
		if failures >= rule.FailureThreshold {
			alert := Alert{
				Level:     AlertLevelCritical,
				Type:      fmt.Sprintf("%s_failure_threshold", rule.Type),
				Message:   fmt.Sprintf("%s: %d failures in %v (threshold: %d)", rule.Name, failures, rule.TimeWindow, rule.FailureThreshold),
				Timestamp: time.Now(),
				Metadata: map[string]interface{}{
					"rule":     rule.Name,
					"failures": failures,
					"total":    total,
					"window":   rule.TimeWindow.String(),
				},
			}
			s.addAlert(alert)
		}
		
		// Check error rate threshold
		if total > 0 {
			errorRate := float64(failures) / float64(total) * 100
			if errorRate >= rule.ErrorRateThreshold {
				alert := Alert{
					Level:     AlertLevelWarning,
					Type:      fmt.Sprintf("%s_error_rate", rule.Type),
					Message:   fmt.Sprintf("%s: %.2f%% error rate in %v (threshold: %.2f%%)", rule.Name, errorRate, rule.TimeWindow, rule.ErrorRateThreshold),
					Timestamp: time.Now(),
					Metadata: map[string]interface{}{
						"rule":       rule.Name,
						"error_rate": errorRate,
						"failures":   failures,
						"total":      total,
						"window":     rule.TimeWindow.String(),
					},
				}
				s.addAlert(alert)
			}
		}
	}
}

// countMetricsInWindow counts metrics of a specific type and status within a time window
func (s *MonitoringService) countMetricsInWindow(metricType MetricType, status MetricStatus, window time.Duration) int {
	cutoff := time.Now().Add(-window)
	count := 0
	
	for i := len(s.metrics) - 1; i >= 0; i-- {
		metric := s.metrics[i]
		
		// Stop if we've gone past the time window
		if metric.Timestamp.Before(cutoff) {
			break
		}
		
		// Count matching metrics
		if metric.Type == metricType {
			if status == "" || metric.Status == status {
				count++
			}
		}
	}
	
	return count
}

// addAlert adds a new alert to the list
func (s *MonitoringService) addAlert(alert Alert) {
	s.alerts = append(s.alerts, alert)
	
	// Update alert counters
	s.aggregatedMetrics.TotalAlerts++
	if alert.Level == AlertLevelCritical {
		s.aggregatedMetrics.CriticalAlerts++
	} else if alert.Level == AlertLevelWarning {
		s.aggregatedMetrics.WarningAlerts++
	}
	
	// Trim alerts if exceeding max size
	if len(s.alerts) > s.maxAlertsSize {
		s.alerts = s.alerts[len(s.alerts)-s.maxAlertsSize:]
	}
}

// GetMetrics returns all metrics (or filtered by type)
func (s *MonitoringService) GetMetrics(ctx context.Context, metricType MetricType, limit int) ([]Metric, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	var filtered []Metric
	
	// Filter by type if specified
	if metricType != "" {
		for i := len(s.metrics) - 1; i >= 0 && len(filtered) < limit; i-- {
			if s.metrics[i].Type == metricType {
				filtered = append(filtered, s.metrics[i])
			}
		}
	} else {
		// Return all metrics up to limit
		start := len(s.metrics) - limit
		if start < 0 {
			start = 0
		}
		filtered = make([]Metric, len(s.metrics)-start)
		copy(filtered, s.metrics[start:])
	}
	
	return filtered, nil
}

// GetAlerts returns recent alerts
func (s *MonitoringService) GetAlerts(ctx context.Context, level AlertLevel, limit int) ([]Alert, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	var filtered []Alert
	
	// Filter by level if specified
	if level != "" {
		for i := len(s.alerts) - 1; i >= 0 && len(filtered) < limit; i-- {
			if s.alerts[i].Level == level {
				filtered = append(filtered, s.alerts[i])
			}
		}
	} else {
		// Return all alerts up to limit
		start := len(s.alerts) - limit
		if start < 0 {
			start = 0
		}
		filtered = make([]Alert, len(s.alerts)-start)
		copy(filtered, s.alerts[start:])
	}
	
	return filtered, nil
}

// GetAggregatedMetrics returns the aggregated metrics
func (s *MonitoringService) GetAggregatedMetrics(ctx context.Context) (*MonitoringMetrics, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	// Create a copy to avoid race conditions
	metricsCopy := *s.aggregatedMetrics
	return &metricsCopy, nil
}

// GetHealthStatus returns the overall health status of the system
func (s *MonitoringService) GetHealthStatus(ctx context.Context) (map[string]interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	// Calculate error rates
	costCalcErrorRate := 0.0
	if s.aggregatedMetrics.TotalCostCalculations > 0 {
		costCalcErrorRate = float64(s.aggregatedMetrics.FailedCostCalculations) / float64(s.aggregatedMetrics.TotalCostCalculations) * 100
	}
	
	recalcErrorRate := 0.0
	if s.aggregatedMetrics.TotalRecalcJobs > 0 {
		recalcErrorRate = float64(s.aggregatedMetrics.FailedRecalcJobs) / float64(s.aggregatedMetrics.TotalRecalcJobs) * 100
	}
	
	shiftClosureErrorRate := 0.0
	if s.aggregatedMetrics.TotalShiftClosures > 0 {
		shiftClosureErrorRate = float64(s.aggregatedMetrics.FailedShiftClosures) / float64(s.aggregatedMetrics.TotalShiftClosures) * 100
	}
	
	// Determine overall health
	health := "healthy"
	if s.aggregatedMetrics.CriticalAlerts > 0 {
		health = "critical"
	} else if s.aggregatedMetrics.WarningAlerts > 0 || costCalcErrorRate > 10 || recalcErrorRate > 10 || shiftClosureErrorRate > 5 {
		health = "degraded"
	}
	
	return map[string]interface{}{
		"status": health,
		"error_rates": map[string]float64{
			"cost_calculation":  costCalcErrorRate,
			"recalculation_job": recalcErrorRate,
			"shift_closure":     shiftClosureErrorRate,
		},
		"recent_alerts": map[string]int64{
			"critical": s.aggregatedMetrics.CriticalAlerts,
			"warning":  s.aggregatedMetrics.WarningAlerts,
			"total":    s.aggregatedMetrics.TotalAlerts,
		},
		"last_updated": s.aggregatedMetrics.LastUpdated,
	}, nil
}

// ResetMetrics resets all metrics and alerts (useful for testing)
func (s *MonitoringService) ResetMetrics() {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.metrics = make([]Metric, 0)
	s.alerts = make([]Alert, 0)
	s.aggregatedMetrics = &MonitoringMetrics{
		LastUpdated: time.Now(),
	}
}
