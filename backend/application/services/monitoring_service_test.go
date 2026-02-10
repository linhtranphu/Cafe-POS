package services

import (
	"context"
	"testing"
	"time"
)

func TestMonitoringService_RecordMetric(t *testing.T) {
	service := NewMonitoringService()
	
	// Record a success metric
	metric := Metric{
		Type:      MetricTypeCostCalculation,
		Status:    MetricStatusSuccess,
		Timestamp: time.Now(),
		Duration:  100 * time.Millisecond,
		Message:   "Test metric",
		Metadata: map[string]interface{}{
			"test": "value",
		},
	}
	
	service.RecordMetric(metric)
	
	// Verify aggregated metrics
	aggregated, err := service.GetAggregatedMetrics(context.Background())
	if err != nil {
		t.Fatalf("Failed to get aggregated metrics: %v", err)
	}
	
	if aggregated.TotalCostCalculations != 1 {
		t.Errorf("Expected 1 total cost calculation, got %d", aggregated.TotalCostCalculations)
	}
	
	if aggregated.SuccessfulCostCalculations != 1 {
		t.Errorf("Expected 1 successful cost calculation, got %d", aggregated.SuccessfulCostCalculations)
	}
	
	if aggregated.FailedCostCalculations != 0 {
		t.Errorf("Expected 0 failed cost calculations, got %d", aggregated.FailedCostCalculations)
	}
}

func TestMonitoringService_AlertOnFailureThreshold(t *testing.T) {
	service := NewMonitoringService()
	
	// Record multiple failures to trigger alert
	for i := 0; i < 11; i++ {
		metric := Metric{
			Type:      MetricTypeCostCalculation,
			Status:    MetricStatusFailure,
			Timestamp: time.Now(),
			Message:   "Test failure",
		}
		service.RecordMetric(metric)
	}
	
	// Check if alert was triggered
	alerts, err := service.GetAlerts(context.Background(), "", 100)
	if err != nil {
		t.Fatalf("Failed to get alerts: %v", err)
	}
	
	if len(alerts) == 0 {
		t.Error("Expected at least one alert to be triggered")
	}
	
	// Verify alert is critical
	foundCritical := false
	for _, alert := range alerts {
		if alert.Level == AlertLevelCritical {
			foundCritical = true
			break
		}
	}
	
	if !foundCritical {
		t.Error("Expected a critical alert to be triggered")
	}
	
	// Verify aggregated metrics
	aggregated, err := service.GetAggregatedMetrics(context.Background())
	if err != nil {
		t.Fatalf("Failed to get aggregated metrics: %v", err)
	}
	
	if aggregated.CriticalAlerts == 0 {
		t.Error("Expected critical alerts counter to be incremented")
	}
}

func TestMonitoringService_AlertOnErrorRate(t *testing.T) {
	service := NewMonitoringService()
	
	// Record mix of success and failures to trigger error rate alert
	// 3 successes, 7 failures = 70% error rate (threshold is 20%)
	for i := 0; i < 3; i++ {
		metric := Metric{
			Type:      MetricTypeRecalculationJob,
			Status:    MetricStatusSuccess,
			Timestamp: time.Now(),
		}
		service.RecordMetric(metric)
	}
	
	for i := 0; i < 7; i++ {
		metric := Metric{
			Type:      MetricTypeRecalculationJob,
			Status:    MetricStatusFailure,
			Timestamp: time.Now(),
		}
		service.RecordMetric(metric)
	}
	
	// Check if alert was triggered
	alerts, err := service.GetAlerts(context.Background(), AlertLevelWarning, 100)
	if err != nil {
		t.Fatalf("Failed to get alerts: %v", err)
	}
	
	if len(alerts) == 0 {
		t.Error("Expected at least one warning alert to be triggered")
	}
}

func TestMonitoringService_GetHealthStatus(t *testing.T) {
	service := NewMonitoringService()
	
	// Initially should be healthy
	health, err := service.GetHealthStatus(context.Background())
	if err != nil {
		t.Fatalf("Failed to get health status: %v", err)
	}
	
	status, ok := health["status"].(string)
	if !ok {
		t.Fatal("Health status should have a status field")
	}
	
	if status != "healthy" {
		t.Errorf("Expected healthy status, got %s", status)
	}
	
	// Record many failures to degrade health
	for i := 0; i < 15; i++ {
		metric := Metric{
			Type:      MetricTypeCostCalculation,
			Status:    MetricStatusFailure,
			Timestamp: time.Now(),
		}
		service.RecordMetric(metric)
	}
	
	// Health should now be degraded or critical
	health, err = service.GetHealthStatus(context.Background())
	if err != nil {
		t.Fatalf("Failed to get health status: %v", err)
	}
	
	status, ok = health["status"].(string)
	if !ok {
		t.Fatal("Health status should have a status field")
	}
	
	if status == "healthy" {
		t.Error("Expected degraded or critical status after failures")
	}
}

func TestMonitoringService_MetricFiltering(t *testing.T) {
	service := NewMonitoringService()
	
	// Record different types of metrics
	service.RecordMetric(Metric{
		Type:      MetricTypeCostCalculation,
		Status:    MetricStatusSuccess,
		Timestamp: time.Now(),
	})
	
	service.RecordMetric(Metric{
		Type:      MetricTypeShiftClosure,
		Status:    MetricStatusSuccess,
		Timestamp: time.Now(),
	})
	
	service.RecordMetric(Metric{
		Type:      MetricTypeCostCalculation,
		Status:    MetricStatusSuccess,
		Timestamp: time.Now(),
	})
	
	// Get only cost calculation metrics
	metrics, err := service.GetMetrics(context.Background(), MetricTypeCostCalculation, 100)
	if err != nil {
		t.Fatalf("Failed to get metrics: %v", err)
	}
	
	if len(metrics) != 2 {
		t.Errorf("Expected 2 cost calculation metrics, got %d", len(metrics))
	}
	
	// Verify all returned metrics are cost calculation type
	for _, metric := range metrics {
		if metric.Type != MetricTypeCostCalculation {
			t.Errorf("Expected only cost calculation metrics, got %s", metric.Type)
		}
	}
}

func TestMonitoringService_AverageDuration(t *testing.T) {
	service := NewMonitoringService()
	
	// Record metrics with different durations
	service.RecordMetric(Metric{
		Type:      MetricTypeCostCalculation,
		Status:    MetricStatusSuccess,
		Timestamp: time.Now(),
		Duration:  100 * time.Millisecond,
	})
	
	service.RecordMetric(Metric{
		Type:      MetricTypeCostCalculation,
		Status:    MetricStatusSuccess,
		Timestamp: time.Now(),
		Duration:  200 * time.Millisecond,
	})
	
	// Get aggregated metrics
	aggregated, err := service.GetAggregatedMetrics(context.Background())
	if err != nil {
		t.Fatalf("Failed to get aggregated metrics: %v", err)
	}
	
	// Average should be 150ms
	expectedAvg := 150 * time.Millisecond
	if aggregated.AverageCostCalcDuration != expectedAvg {
		t.Errorf("Expected average duration %v, got %v", expectedAvg, aggregated.AverageCostCalcDuration)
	}
}

func TestMonitoringService_ResetMetrics(t *testing.T) {
	service := NewMonitoringService()
	
	// Record some metrics
	service.RecordMetric(Metric{
		Type:      MetricTypeCostCalculation,
		Status:    MetricStatusSuccess,
		Timestamp: time.Now(),
	})
	
	// Reset metrics
	service.ResetMetrics()
	
	// Verify metrics are cleared
	aggregated, err := service.GetAggregatedMetrics(context.Background())
	if err != nil {
		t.Fatalf("Failed to get aggregated metrics: %v", err)
	}
	
	if aggregated.TotalCostCalculations != 0 {
		t.Errorf("Expected 0 total cost calculations after reset, got %d", aggregated.TotalCostCalculations)
	}
	
	metrics, err := service.GetMetrics(context.Background(), "", 100)
	if err != nil {
		t.Fatalf("Failed to get metrics: %v", err)
	}
	
	if len(metrics) != 0 {
		t.Errorf("Expected 0 metrics after reset, got %d", len(metrics))
	}
}

func TestMonitoringService_ShiftClosureMetrics(t *testing.T) {
	service := NewMonitoringService()
	
	// Record shift closure metrics
	service.RecordMetric(Metric{
		Type:      MetricTypeShiftClosure,
		Status:    MetricStatusSuccess,
		Timestamp: time.Now(),
		Duration:  500 * time.Millisecond,
		Metadata: map[string]interface{}{
			"shift_id":     "test-shift",
			"total_orders": 10,
		},
	})
	
	// Verify aggregated metrics
	aggregated, err := service.GetAggregatedMetrics(context.Background())
	if err != nil {
		t.Fatalf("Failed to get aggregated metrics: %v", err)
	}
	
	if aggregated.TotalShiftClosures != 1 {
		t.Errorf("Expected 1 total shift closure, got %d", aggregated.TotalShiftClosures)
	}
	
	if aggregated.SuccessfulShiftClosures != 1 {
		t.Errorf("Expected 1 successful shift closure, got %d", aggregated.SuccessfulShiftClosures)
	}
	
	if aggregated.AverageShiftClosureDuration != 500*time.Millisecond {
		t.Errorf("Expected average duration 500ms, got %v", aggregated.AverageShiftClosureDuration)
	}
}

func TestMonitoringService_ProfitAnalysisMetrics(t *testing.T) {
	service := NewMonitoringService()
	
	// Record profit analysis metrics
	service.RecordMetric(Metric{
		Type:      MetricTypeProfitAnalysis,
		Status:    MetricStatusSuccess,
		Timestamp: time.Now(),
		Duration:  200 * time.Millisecond,
	})
	
	service.RecordMetric(Metric{
		Type:      MetricTypeProfitAnalysis,
		Status:    MetricStatusFailure,
		Timestamp: time.Now(),
		Duration:  100 * time.Millisecond,
	})
	
	// Verify aggregated metrics
	aggregated, err := service.GetAggregatedMetrics(context.Background())
	if err != nil {
		t.Fatalf("Failed to get aggregated metrics: %v", err)
	}
	
	if aggregated.TotalProfitAnalyses != 2 {
		t.Errorf("Expected 2 total profit analyses, got %d", aggregated.TotalProfitAnalyses)
	}
	
	if aggregated.SuccessfulProfitAnalyses != 1 {
		t.Errorf("Expected 1 successful profit analysis, got %d", aggregated.SuccessfulProfitAnalyses)
	}
	
	if aggregated.FailedProfitAnalyses != 1 {
		t.Errorf("Expected 1 failed profit analysis, got %d", aggregated.FailedProfitAnalyses)
	}
}
