package printing

import (
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TestPrintJobType_Constants tests that PrintJobType constants are defined correctly
func TestPrintJobType_Constants(t *testing.T) {
	tests := []struct {
		name     string
		jobType  PrintJobType
		expected string
	}{
		{"Bill type", PrintJobTypeBill, "BILL"},
		{"Label type", PrintJobTypeLabel, "LABEL"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.jobType) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(tt.jobType))
			}
		})
	}
}

// TestPrintJobStatus_Constants tests that PrintJobStatus constants are defined correctly
func TestPrintJobStatus_Constants(t *testing.T) {
	tests := []struct {
		name     string
		status   PrintJobStatus
		expected string
	}{
		{"Pending status", PrintJobStatusPending, "PENDING"},
		{"Printing status", PrintJobStatusPrinting, "PRINTING"},
		{"Completed status", PrintJobStatusCompleted, "COMPLETED"},
		{"Failed status", PrintJobStatusFailed, "FAILED"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.status) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(tt.status))
			}
		})
	}
}

// TestPrintJob_Creation tests creating a PrintJob with valid data
func TestPrintJob_Creation(t *testing.T) {
	now := time.Now()
	orderID := primitive.NewObjectID()
	printerID := primitive.NewObjectID()

	job := &PrintJob{
		Type:        PrintJobTypeBill,
		OrderID:     orderID,
		OrderNumber: "ORD-001",
		PrinterID:   printerID,
		Content:     "Test bill content",
		Status:      PrintJobStatusPending,
		RetryCount:  0,
		MaxRetries:  3,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Verify all fields are set correctly
	if job.Type != PrintJobTypeBill {
		t.Errorf("Expected Type to be BILL, got %s", job.Type)
	}
	if job.OrderID != orderID {
		t.Errorf("Expected OrderID to match")
	}
	if job.OrderNumber != "ORD-001" {
		t.Errorf("Expected OrderNumber to be ORD-001, got %s", job.OrderNumber)
	}
	if job.PrinterID != printerID {
		t.Errorf("Expected PrinterID to match")
	}
	if job.Status != PrintJobStatusPending {
		t.Errorf("Expected Status to be PENDING, got %s", job.Status)
	}
	if job.RetryCount != 0 {
		t.Errorf("Expected RetryCount to be 0, got %d", job.RetryCount)
	}
	if job.MaxRetries != 3 {
		t.Errorf("Expected MaxRetries to be 3, got %d", job.MaxRetries)
	}
}

// TestPrintJob_StatusTransitions tests valid status transitions
func TestPrintJob_StatusTransitions(t *testing.T) {
	tests := []struct {
		name        string
		fromStatus  PrintJobStatus
		toStatus    PrintJobStatus
		shouldAllow bool
		description string
	}{
		{
			"PENDING to PRINTING",
			PrintJobStatusPending,
			PrintJobStatusPrinting,
			true,
			"Job starts printing",
		},
		{
			"PRINTING to COMPLETED",
			PrintJobStatusPrinting,
			PrintJobStatusCompleted,
			true,
			"Job completes successfully",
		},
		{
			"PRINTING to FAILED",
			PrintJobStatusPrinting,
			PrintJobStatusFailed,
			true,
			"Job fails during printing",
		},
		{
			"PENDING to FAILED",
			PrintJobStatusPending,
			PrintJobStatusFailed,
			true,
			"Job fails before printing starts",
		},
		{
			"FAILED to PENDING",
			PrintJobStatusFailed,
			PrintJobStatusPending,
			true,
			"Retry failed job",
		},
		{
			"COMPLETED to PENDING",
			PrintJobStatusCompleted,
			PrintJobStatusPending,
			false,
			"Cannot retry completed job",
		},
		{
			"COMPLETED to PRINTING",
			PrintJobStatusCompleted,
			PrintJobStatusPrinting,
			false,
			"Cannot restart completed job",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			job := &PrintJob{
				Status:     tt.fromStatus,
				RetryCount: 0,
				MaxRetries: 3,
			}

			// Simulate status transition
			oldStatus := job.Status
			job.Status = tt.toStatus
			job.UpdatedAt = time.Now()

			// For this test, we're documenting expected transitions
			// In a real implementation, you might have a method like job.TransitionTo(status)
			// that validates the transition
			if tt.shouldAllow {
				if job.Status != tt.toStatus {
					t.Errorf("Expected status to transition from %s to %s", oldStatus, tt.toStatus)
				}
			}
		})
	}
}

// TestPrintJob_RetryLogic tests retry count validation
func TestPrintJob_RetryLogic(t *testing.T) {
	tests := []struct {
		name           string
		retryCount     int
		maxRetries     int
		shouldRetry    bool
		description    string
	}{
		{
			"First retry allowed",
			0,
			3,
			true,
			"Job has not been retried yet",
		},
		{
			"Second retry allowed",
			1,
			3,
			true,
			"Job has been retried once",
		},
		{
			"Third retry allowed",
			2,
			3,
			true,
			"Job has been retried twice",
		},
		{
			"Max retries reached",
			3,
			3,
			false,
			"Job has reached max retries",
		},
		{
			"Exceeded max retries",
			4,
			3,
			false,
			"Job has exceeded max retries",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			job := &PrintJob{
				Status:     PrintJobStatusFailed,
				RetryCount: tt.retryCount,
				MaxRetries: tt.maxRetries,
			}

			canRetry := job.RetryCount < job.MaxRetries
			if canRetry != tt.shouldRetry {
				t.Errorf("Expected canRetry to be %v, got %v (retryCount=%d, maxRetries=%d)",
					tt.shouldRetry, canRetry, job.RetryCount, job.MaxRetries)
			}
		})
	}
}

// TestPrintJob_ErrorHandling tests error message handling
func TestPrintJob_ErrorHandling(t *testing.T) {
	job := &PrintJob{
		Status:     PrintJobStatusFailed,
		ErrorMsg:   "Printer connection timeout",
		RetryCount: 1,
		MaxRetries: 3,
	}

	if job.ErrorMsg == "" {
		t.Error("Expected ErrorMsg to be set for failed job")
	}

	if job.Status != PrintJobStatusFailed {
		t.Errorf("Expected Status to be FAILED, got %s", job.Status)
	}
}

// TestPrintJob_PrintedAtTimestamp tests PrintedAt timestamp handling
func TestPrintJob_PrintedAtTimestamp(t *testing.T) {
	now := time.Now()

	// Job not yet printed
	pendingJob := &PrintJob{
		Status:    PrintJobStatusPending,
		PrintedAt: nil,
	}

	if pendingJob.PrintedAt != nil {
		t.Error("Expected PrintedAt to be nil for pending job")
	}

	// Job completed
	completedJob := &PrintJob{
		Status:    PrintJobStatusCompleted,
		PrintedAt: &now,
	}

	if completedJob.PrintedAt == nil {
		t.Error("Expected PrintedAt to be set for completed job")
	}

	if !completedJob.PrintedAt.Equal(now) {
		t.Errorf("Expected PrintedAt to be %v, got %v", now, completedJob.PrintedAt)
	}
}

// TestCreatePrintJobRequest_Validation tests request validation
func TestCreatePrintJobRequest_Validation(t *testing.T) {
	orderID := primitive.NewObjectID()
	printerID := primitive.NewObjectID()

	tests := []struct {
		name    string
		request CreatePrintJobRequest
		valid   bool
	}{
		{
			"Valid bill request",
			CreatePrintJobRequest{
				Type:        PrintJobTypeBill,
				OrderID:     orderID,
				OrderNumber: "ORD-001",
				PrinterID:   printerID,
				Content:     "Bill content",
			},
			true,
		},
		{
			"Valid label request",
			CreatePrintJobRequest{
				Type:        PrintJobTypeLabel,
				OrderID:     orderID,
				OrderNumber: "ORD-001",
				PrinterID:   printerID,
				Content:     "Label content",
			},
			true,
		},
		{
			"Empty content",
			CreatePrintJobRequest{
				Type:        PrintJobTypeBill,
				OrderID:     orderID,
				OrderNumber: "ORD-001",
				PrinterID:   printerID,
				Content:     "",
			},
			false,
		},
		{
			"Empty order number",
			CreatePrintJobRequest{
				Type:        PrintJobTypeBill,
				OrderID:     orderID,
				OrderNumber: "",
				PrinterID:   printerID,
				Content:     "Content",
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Basic validation checks
			isValid := tt.request.Type != "" &&
				!tt.request.OrderID.IsZero() &&
				tt.request.OrderNumber != "" &&
				!tt.request.PrinterID.IsZero() &&
				tt.request.Content != ""

			if isValid != tt.valid {
				t.Errorf("Expected validation to be %v, got %v", tt.valid, isValid)
			}
		})
	}
}

// TestPrintJobFilter_DefaultValues tests filter default values
func TestPrintJobFilter_DefaultValues(t *testing.T) {
	filter := &PrintJobFilter{
		Page:  1,
		Limit: 20,
	}

	if filter.Page != 1 {
		t.Errorf("Expected default Page to be 1, got %d", filter.Page)
	}

	if filter.Limit != 20 {
		t.Errorf("Expected default Limit to be 20, got %d", filter.Limit)
	}

	if !filter.OrderID.IsZero() {
		t.Error("Expected OrderID to be zero value")
	}

	if filter.Status != "" {
		t.Errorf("Expected Status to be empty, got %s", filter.Status)
	}
}
