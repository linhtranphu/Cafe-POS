# Monitoring and Alerts System

## Overview

The Menu Cost & Profit Analysis feature includes a comprehensive monitoring and alerting system that tracks metrics, detects anomalies, and alerts managers to potential issues.

## Features

### 1. Metrics Collection

The system automatically collects metrics for:

- **Cost Calculations**: Track success/failure rates and duration of menu item cost calculations
- **Recalculation Jobs**: Monitor background job processing for cost updates
- **Shift Closures**: Track shift closure cost calculation performance
- **Profit Analysis**: Monitor profit analysis query performance
- **Operating Expenses**: Track expense management operations

### 2. Alert Rules

The system includes predefined alert rules that trigger based on:

#### High Cost Calculation Failure Rate
- **Threshold**: 10 failures in 5 minutes OR >20% error rate
- **Level**: Critical
- **Action**: Investigate ingredient data quality or database connectivity

#### High Recalculation Job Failure Rate
- **Threshold**: 5 failures in 5 minutes OR >15% error rate
- **Level**: Critical
- **Action**: Check worker pool health and queue capacity

#### Shift Closure Failures
- **Threshold**: 3 failures in 10 minutes OR >10% error rate
- **Level**: Critical
- **Action**: Investigate order data integrity or ingredient availability

#### Profit Analysis Failures
- **Threshold**: 10 failures in 5 minutes OR >25% error rate
- **Level**: Warning
- **Action**: Check database query performance

### 3. Health Status

The system provides an overall health status:

- **Healthy**: All systems operating normally, error rates below thresholds
- **Degraded**: Warning alerts present or error rates elevated (10-20%)
- **Critical**: Critical alerts present or error rates very high (>20%)

## API Endpoints

### Get Metrics

```
GET /api/manager/monitoring/metrics?type=cost_calculation&limit=100
```

**Query Parameters**:
- `type` (optional): Filter by metric type (`cost_calculation`, `recalculation_job`, `shift_closure`, `profit_analysis`, `operating_expense`)
- `limit` (optional): Number of metrics to return (default: 100, max: 1000)

**Response**:
```json
{
  "metrics": [
    {
      "type": "cost_calculation",
      "status": "success",
      "timestamp": "2024-01-15T10:30:00Z",
      "duration": 150000000,
      "message": "Cost calculation completed",
      "metadata": {
        "menu_item_id": "...",
        "cost": 15000
      }
    }
  ],
  "count": 1
}
```

### Get Alerts

```
GET /api/manager/monitoring/alerts?level=critical&limit=50
```

**Query Parameters**:
- `level` (optional): Filter by alert level (`info`, `warning`, `critical`)
- `limit` (optional): Number of alerts to return (default: 50, max: 500)

**Response**:
```json
{
  "alerts": [
    {
      "level": "critical",
      "type": "cost_calculation_failure_threshold",
      "message": "High Cost Calculation Failure Rate: 12 failures in 5m0s (threshold: 10)",
      "timestamp": "2024-01-15T10:35:00Z",
      "metadata": {
        "rule": "High Cost Calculation Failure Rate",
        "failures": 12,
        "total": 50,
        "window": "5m0s"
      }
    }
  ],
  "count": 1
}
```

### Get Aggregated Metrics

```
GET /api/manager/monitoring/metrics/aggregated
```

**Response**:
```json
{
  "total_cost_calculations": 1000,
  "successful_cost_calculations": 980,
  "failed_cost_calculations": 20,
  "average_cost_calc_duration": 150000000,
  "total_recalc_jobs": 500,
  "successful_recalc_jobs": 495,
  "failed_recalc_jobs": 5,
  "average_recalc_duration": 200000000,
  "total_shift_closures": 50,
  "successful_shift_closures": 49,
  "failed_shift_closures": 1,
  "average_shift_closure_duration": 500000000,
  "total_profit_analyses": 200,
  "successful_profit_analyses": 198,
  "failed_profit_analyses": 2,
  "average_profit_analysis_duration": 300000000,
  "total_expense_operations": 100,
  "successful_expense_operations": 100,
  "failed_expense_operations": 0,
  "total_alerts": 5,
  "critical_alerts": 2,
  "warning_alerts": 3,
  "last_updated": "2024-01-15T10:40:00Z"
}
```

### Get Health Status

```
GET /api/manager/monitoring/health
```

**Response**:
```json
{
  "status": "healthy",
  "error_rates": {
    "cost_calculation": 2.0,
    "recalculation_job": 1.0,
    "shift_closure": 2.0
  },
  "recent_alerts": {
    "critical": 0,
    "warning": 1,
    "total": 1
  },
  "last_updated": "2024-01-15T10:40:00Z"
}
```

**HTTP Status Codes**:
- `200 OK`: System is healthy or degraded
- `503 Service Unavailable`: System is in critical state

## Integration

### Cost Calculator Service

The cost calculator service automatically records metrics for:

- Menu item cost calculations (success/failure/warning)
- Shift closure cost calculations
- Ingredient cost updates

Example metric recorded:
```go
monitoringService.RecordMetric(Metric{
    Type:      MetricTypeCostCalculation,
    Status:    MetricStatusSuccess,
    Timestamp: time.Now(),
    Duration:  150 * time.Millisecond,
    Message:   "Cost calculation completed",
    Metadata: map[string]interface{}{
        "menu_item_id": menuItemID.Hex(),
        "cost":         15000.0,
        "cost_status":  "FINAL",
    },
})
```

### Cost Recalculation Service

The recalculation service records metrics for background jobs:

- Job success/failure
- Retry attempts
- Processing duration

Example metric recorded:
```go
monitoringService.RecordMetric(Metric{
    Type:      MetricTypeRecalculationJob,
    Status:    MetricStatusSuccess,
    Timestamp: time.Now(),
    Duration:  200 * time.Millisecond,
    Message:   "Cost recalculation completed successfully",
    Metadata: map[string]interface{}{
        "menu_item_id": menuItemID.Hex(),
        "attempts":     1,
    },
})
```

## Troubleshooting

### High Error Rates

If you see high error rates in cost calculations:

1. **Check ingredient data quality**:
   - Verify all ingredients have valid `cost_per_unit` values
   - Check for missing or deleted ingredients referenced in menu items

2. **Check database connectivity**:
   - Verify MongoDB connection is stable
   - Check for slow queries or timeouts

3. **Check worker pool health**:
   - Verify recalculation workers are running
   - Check queue capacity (default: 1000 items)

### Critical Alerts

When critical alerts are triggered:

1. **Review recent changes**:
   - Check if ingredient costs were recently updated
   - Verify menu item recipes are valid

2. **Check system resources**:
   - Monitor CPU and memory usage
   - Check database performance

3. **Review error logs**:
   - Check application logs for detailed error messages
   - Look for patterns in failures

### Performance Issues

If average durations are high:

1. **Optimize database queries**:
   - Ensure proper indexes are in place
   - Consider caching frequently accessed data

2. **Adjust worker pool size**:
   - Increase number of workers for recalculation jobs
   - Adjust queue size if needed

3. **Review data volume**:
   - Check if menu items have excessive ingredients
   - Consider pagination for large result sets

## Best Practices

1. **Regular Monitoring**: Check health status and metrics daily
2. **Alert Response**: Investigate critical alerts within 1 hour
3. **Metric Review**: Review aggregated metrics weekly to identify trends
4. **Capacity Planning**: Monitor queue sizes and adjust worker pool as needed
5. **Data Quality**: Maintain complete ingredient cost data to minimize warnings

## Configuration

Alert rules can be customized by modifying the `getDefaultAlertRules()` function in `monitoring_service.go`:

```go
{
    Name:              "High Cost Calculation Failure Rate",
    Type:              MetricTypeCostCalculation,
    FailureThreshold:  10,  // Adjust threshold
    TimeWindow:        5 * time.Minute,  // Adjust time window
    ErrorRateThreshold: 20.0,  // Adjust error rate %
    Enabled:           true,
}
```

## Future Enhancements

Potential improvements for the monitoring system:

1. **External Alerting**: Integration with email, Slack, or PagerDuty
2. **Metric Persistence**: Store metrics in time-series database (e.g., InfluxDB)
3. **Dashboards**: Grafana integration for visualization
4. **Custom Alert Rules**: Allow managers to configure alert rules via UI
5. **Metric Retention**: Configurable retention policies for historical data
6. **Performance Profiling**: Detailed performance traces for slow operations
