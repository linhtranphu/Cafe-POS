package mongodb

import (
	"context"
	"testing"
	"time"

	"cafe-pos/backend/domain/printing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// setupTestDB creates a test database connection
func setupTestDB(t *testing.T) (*mongo.Database, func()) {
	ctx := context.Background()
	
	// Use localhost MongoDB for integration tests with auth
	// Password matches MONGODB_URI in .env
	mongoURI := "mongodb://admin:password@localhost:27017/?authSource=admin"
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	require.NoError(t, err, "Failed to connect to MongoDB")
	
	// Use a unique test database for each test run
	testDB := client.Database("cafe_pos_print_test_" + primitive.NewObjectID().Hex())
	
	// Cleanup function
	cleanup := func() {
		// Drop test database
		err := testDB.Drop(ctx)
		if err != nil {
			t.Logf("Warning: Failed to drop test database: %v", err)
		}
		
		// Disconnect client
		err = client.Disconnect(ctx)
		if err != nil {
			t.Logf("Warning: Failed to disconnect from MongoDB: %v", err)
		}
	}
	
	return testDB, cleanup
}

// TestPrintJobRepository_Integration tests PrintJobRepository with real MongoDB
func TestPrintJobRepository_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	
	db, cleanup := setupTestDB(t)
	defer cleanup()
	
	ctx := context.Background()
	repo := NewPrintJobRepository(db)
	
	// Create indexes
	err := repo.CreateIndexes(ctx)
	require.NoError(t, err, "Failed to create indexes")
	
	t.Run("Create and FindByID", func(t *testing.T) {
		job := &printing.PrintJob{
			Type:        printing.PrintJobTypeBill,
			OrderID:     primitive.NewObjectID(),
			OrderNumber: "ORD-001",
			PrinterID:   primitive.NewObjectID(),
			Content:     "Test bill content",
			Status:      printing.PrintJobStatusPending,
			RetryCount:  0,
			MaxRetries:  3,
		}
		
		// Create
		err := repo.Create(ctx, job)
		require.NoError(t, err, "Failed to create print job")
		assert.False(t, job.ID.IsZero(), "ID should be set after creation")
		assert.False(t, job.CreatedAt.IsZero(), "CreatedAt should be set")
		assert.False(t, job.UpdatedAt.IsZero(), "UpdatedAt should be set")
		
		// FindByID
		found, err := repo.FindByID(ctx, job.ID)
		require.NoError(t, err, "Failed to find print job")
		assert.Equal(t, job.ID, found.ID)
		assert.Equal(t, job.Type, found.Type)
		assert.Equal(t, job.OrderNumber, found.OrderNumber)
		assert.Equal(t, job.Status, found.Status)
	})
	
	t.Run("FindByOrderID", func(t *testing.T) {
		orderID := primitive.NewObjectID()
		
		// Create multiple jobs for same order
		billJob := &printing.PrintJob{
			Type:        printing.PrintJobTypeBill,
			OrderID:     orderID,
			OrderNumber: "ORD-002",
			PrinterID:   primitive.NewObjectID(),
			Content:     "Bill content",
			Status:      printing.PrintJobStatusPending,
			MaxRetries:  3,
		}
		err := repo.Create(ctx, billJob)
		require.NoError(t, err)
		
		labelJob1 := &printing.PrintJob{
			Type:        printing.PrintJobTypeLabel,
			OrderID:     orderID,
			OrderNumber: "ORD-002",
			PrinterID:   primitive.NewObjectID(),
			Content:     "Label 1 content",
			Status:      printing.PrintJobStatusPending,
			MaxRetries:  3,
		}
		err = repo.Create(ctx, labelJob1)
		require.NoError(t, err)
		
		labelJob2 := &printing.PrintJob{
			Type:        printing.PrintJobTypeLabel,
			OrderID:     orderID,
			OrderNumber: "ORD-002",
			PrinterID:   primitive.NewObjectID(),
			Content:     "Label 2 content",
			Status:      printing.PrintJobStatusPending,
			MaxRetries:  3,
		}
		err = repo.Create(ctx, labelJob2)
		require.NoError(t, err)
		
		// Find by order ID
		jobs, err := repo.FindByOrderID(ctx, orderID)
		require.NoError(t, err, "Failed to find jobs by order ID")
		assert.Len(t, jobs, 3, "Should find 3 jobs for the order")
		
		// Verify jobs are sorted by created_at
		for i := 1; i < len(jobs); i++ {
			assert.True(t, jobs[i].CreatedAt.After(jobs[i-1].CreatedAt) || 
				jobs[i].CreatedAt.Equal(jobs[i-1].CreatedAt),
				"Jobs should be sorted by created_at ascending")
		}
	})
	
	t.Run("FindPending", func(t *testing.T) {
		// Create pending jobs
		for i := 0; i < 5; i++ {
			job := &printing.PrintJob{
				Type:        printing.PrintJobTypeBill,
				OrderID:     primitive.NewObjectID(),
				OrderNumber: "ORD-PENDING",
				PrinterID:   primitive.NewObjectID(),
				Content:     "Pending content",
				Status:      printing.PrintJobStatusPending,
				MaxRetries:  3,
			}
			err := repo.Create(ctx, job)
			require.NoError(t, err)
			time.Sleep(1 * time.Millisecond) // Ensure different timestamps
		}
		
		// Create completed job (should not be returned)
		completedJob := &printing.PrintJob{
			Type:        printing.PrintJobTypeBill,
			OrderID:     primitive.NewObjectID(),
			OrderNumber: "ORD-COMPLETED",
			PrinterID:   primitive.NewObjectID(),
			Content:     "Completed content",
			Status:      printing.PrintJobStatusCompleted,
			MaxRetries:  3,
		}
		err := repo.Create(ctx, completedJob)
		require.NoError(t, err)
		
		// Find pending with limit
		pending, err := repo.FindPending(ctx, 3)
		require.NoError(t, err, "Failed to find pending jobs")
		assert.LessOrEqual(t, len(pending), 3, "Should respect limit")
		
		// Verify all are pending
		for _, job := range pending {
			assert.Equal(t, printing.PrintJobStatusPending, job.Status)
		}
		
		// Find all pending (no limit)
		allPending, err := repo.FindPending(ctx, 0)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(allPending), 5, "Should find at least 5 pending jobs")
	})
	
	t.Run("FindFailed", func(t *testing.T) {
		// Create failed jobs
		for i := 0; i < 3; i++ {
			job := &printing.PrintJob{
				Type:        printing.PrintJobTypeBill,
				OrderID:     primitive.NewObjectID(),
				OrderNumber: "ORD-FAILED",
				PrinterID:   primitive.NewObjectID(),
				Content:     "Failed content",
				Status:      printing.PrintJobStatusFailed,
				RetryCount:  3,
				MaxRetries:  3,
				ErrorMsg:    "Printer offline",
			}
			err := repo.Create(ctx, job)
			require.NoError(t, err)
		}
		
		// Find failed
		failed, err := repo.FindFailed(ctx)
		require.NoError(t, err, "Failed to find failed jobs")
		assert.GreaterOrEqual(t, len(failed), 3, "Should find at least 3 failed jobs")
		
		// Verify all are failed
		for _, job := range failed {
			assert.Equal(t, printing.PrintJobStatusFailed, job.Status)
		}
	})
	
	t.Run("UpdateStatus", func(t *testing.T) {
		job := &printing.PrintJob{
			Type:        printing.PrintJobTypeBill,
			OrderID:     primitive.NewObjectID(),
			OrderNumber: "ORD-UPDATE",
			PrinterID:   primitive.NewObjectID(),
			Content:     "Update test",
			Status:      printing.PrintJobStatusPending,
			MaxRetries:  3,
		}
		err := repo.Create(ctx, job)
		require.NoError(t, err)
		
		// Update to PRINTING
		err = repo.UpdateStatus(ctx, job.ID, printing.PrintJobStatusPrinting, "")
		require.NoError(t, err, "Failed to update status to PRINTING")
		
		found, err := repo.FindByID(ctx, job.ID)
		require.NoError(t, err)
		assert.Equal(t, printing.PrintJobStatusPrinting, found.Status)
		assert.Nil(t, found.PrintedAt, "PrintedAt should be nil for PRINTING status")
		
		// Update to COMPLETED
		err = repo.UpdateStatus(ctx, job.ID, printing.PrintJobStatusCompleted, "")
		require.NoError(t, err, "Failed to update status to COMPLETED")
		
		found, err = repo.FindByID(ctx, job.ID)
		require.NoError(t, err)
		assert.Equal(t, printing.PrintJobStatusCompleted, found.Status)
		assert.NotNil(t, found.PrintedAt, "PrintedAt should be set for COMPLETED status")
		
		// Update to FAILED with error message
		job2 := &printing.PrintJob{
			Type:        printing.PrintJobTypeBill,
			OrderID:     primitive.NewObjectID(),
			OrderNumber: "ORD-FAIL",
			PrinterID:   primitive.NewObjectID(),
			Content:     "Fail test",
			Status:      printing.PrintJobStatusPending,
			MaxRetries:  3,
		}
		err = repo.Create(ctx, job2)
		require.NoError(t, err)
		
		errorMsg := "Connection timeout"
		err = repo.UpdateStatus(ctx, job2.ID, printing.PrintJobStatusFailed, errorMsg)
		require.NoError(t, err)
		
		found, err = repo.FindByID(ctx, job2.ID)
		require.NoError(t, err)
		assert.Equal(t, printing.PrintJobStatusFailed, found.Status)
		assert.Equal(t, errorMsg, found.ErrorMsg)
	})
	
	t.Run("IncrementRetry", func(t *testing.T) {
		job := &printing.PrintJob{
			Type:        printing.PrintJobTypeBill,
			OrderID:     primitive.NewObjectID(),
			OrderNumber: "ORD-RETRY",
			PrinterID:   primitive.NewObjectID(),
			Content:     "Retry test",
			Status:      printing.PrintJobStatusPending,
			RetryCount:  0,
			MaxRetries:  3,
		}
		err := repo.Create(ctx, job)
		require.NoError(t, err)
		
		// Increment retry count
		err = repo.IncrementRetry(ctx, job.ID)
		require.NoError(t, err)
		
		found, err := repo.FindByID(ctx, job.ID)
		require.NoError(t, err)
		assert.Equal(t, 1, found.RetryCount)
		
		// Increment again
		err = repo.IncrementRetry(ctx, job.ID)
		require.NoError(t, err)
		
		found, err = repo.FindByID(ctx, job.ID)
		require.NoError(t, err)
		assert.Equal(t, 2, found.RetryCount)
	})
	
	t.Run("Delete", func(t *testing.T) {
		job := &printing.PrintJob{
			Type:        printing.PrintJobTypeBill,
			OrderID:     primitive.NewObjectID(),
			OrderNumber: "ORD-DELETE",
			PrinterID:   primitive.NewObjectID(),
			Content:     "Delete test",
			Status:      printing.PrintJobStatusPending,
			MaxRetries:  3,
		}
		err := repo.Create(ctx, job)
		require.NoError(t, err)
		
		// Delete
		err = repo.Delete(ctx, job.ID)
		require.NoError(t, err)
		
		// Verify deleted
		_, err = repo.FindByID(ctx, job.ID)
		assert.Error(t, err, "Should return error for deleted job")
		assert.Equal(t, mongo.ErrNoDocuments, err)
	})
	
	t.Run("DeleteOldCompleted", func(t *testing.T) {
		// Create old completed job
		oldJob := &printing.PrintJob{
			Type:        printing.PrintJobTypeBill,
			OrderID:     primitive.NewObjectID(),
			OrderNumber: "ORD-OLD",
			PrinterID:   primitive.NewObjectID(),
			Content:     "Old completed",
			Status:      printing.PrintJobStatusCompleted,
			MaxRetries:  3,
		}
		err := repo.Create(ctx, oldJob)
		require.NoError(t, err)
		
		// Wait to ensure time difference
		time.Sleep(2 * time.Second)
		
		// Create recent completed job
		recentJob := &printing.PrintJob{
			Type:        printing.PrintJobTypeBill,
			OrderID:     primitive.NewObjectID(),
			OrderNumber: "ORD-RECENT",
			PrinterID:   primitive.NewObjectID(),
			Content:     "Recent completed",
			Status:      printing.PrintJobStatusCompleted,
			MaxRetries:  3,
		}
		err = repo.Create(ctx, recentJob)
		require.NoError(t, err)
		
		// Delete jobs older than 1 second from now (should only delete oldJob)
		cutoffTime := time.Now().Add(-1 * time.Second)
		err = repo.DeleteOldCompleted(ctx, cutoffTime)
		require.NoError(t, err)
		
		// Verify oldJob is deleted
		_, err = repo.FindByID(ctx, oldJob.ID)
		assert.Error(t, err, "Old job should be deleted")
		
		// Verify recentJob still exists
		found, err := repo.FindByID(ctx, recentJob.ID)
		require.NoError(t, err, "Recent job should still exist")
		assert.Equal(t, recentJob.ID, found.ID)
	})
}

// TestPrinterConfigRepository_Integration tests PrinterConfigRepository with real MongoDB
func TestPrinterConfigRepository_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	
	db, cleanup := setupTestDB(t)
	defer cleanup()
	
	ctx := context.Background()
	repo := NewPrinterConfigRepository(db)
	
	// Create indexes
	err := repo.CreateIndexes(ctx)
	require.NoError(t, err, "Failed to create indexes")
	
	t.Run("Create and FindByID", func(t *testing.T) {
		config := &printing.PrinterConfig{
			Name:           "Bill Printer 1",
			Type:           printing.PrinterTypeBill,
			ConnectionType: printing.ConnectionTypeNetwork,
			IPAddress:      "192.168.1.100",
			Port:           9100,
			PaperWidth:     80,
			IsDefault:      true,
			IsEnabled:      true,
		}
		
		// Create
		err := repo.Create(ctx, config)
		require.NoError(t, err, "Failed to create printer config")
		assert.False(t, config.ID.IsZero(), "ID should be set after creation")
		assert.False(t, config.CreatedAt.IsZero(), "CreatedAt should be set")
		
		// FindByID
		found, err := repo.FindByID(ctx, config.ID)
		require.NoError(t, err, "Failed to find printer config")
		assert.Equal(t, config.Name, found.Name)
		assert.Equal(t, config.Type, found.Type)
		assert.Equal(t, config.IPAddress, found.IPAddress)
		assert.Equal(t, config.Port, found.Port)
	})
	
	t.Run("FindAll", func(t *testing.T) {
		// Create multiple configs
		configs := []*printing.PrinterConfig{
			{
				Name:           "Bill Printer 2",
				Type:           printing.PrinterTypeBill,
				ConnectionType: printing.ConnectionTypeNetwork,
				IPAddress:      "192.168.1.101",
				Port:           9100,
				PaperWidth:     58,
				IsEnabled:      true,
			},
			{
				Name:           "Label Printer 1",
				Type:           printing.PrinterTypeLabel,
				ConnectionType: printing.ConnectionTypeUSB,
				USBPath:        "/dev/usb/lp0",
				PaperWidth:     40,
				IsEnabled:      true,
			},
		}
		
		for _, config := range configs {
			err := repo.Create(ctx, config)
			require.NoError(t, err)
		}
		
		// FindAll
		all, err := repo.FindAll(ctx)
		require.NoError(t, err, "Failed to find all printer configs")
		assert.GreaterOrEqual(t, len(all), 2, "Should find at least 2 configs")
	})
	
	t.Run("FindByType", func(t *testing.T) {
		// Create bill printer
		billConfig := &printing.PrinterConfig{
			Name:           "Bill Printer 3",
			Type:           printing.PrinterTypeBill,
			ConnectionType: printing.ConnectionTypeNetwork,
			IPAddress:      "192.168.1.102",
			Port:           9100,
			PaperWidth:     80,
			IsEnabled:      true,
		}
		err := repo.Create(ctx, billConfig)
		require.NoError(t, err)
		
		// Create label printer
		labelConfig := &printing.PrinterConfig{
			Name:           "Label Printer 2",
			Type:           printing.PrinterTypeLabel,
			ConnectionType: printing.ConnectionTypeUSB,
			USBPath:        "/dev/usb/lp1",
			PaperWidth:     50,
			IsEnabled:      true,
		}
		err = repo.Create(ctx, labelConfig)
		require.NoError(t, err)
		
		// Find by type BILL
		billPrinters, err := repo.FindByType(ctx, printing.PrinterTypeBill)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(billPrinters), 1, "Should find at least 1 bill printer")
		for _, p := range billPrinters {
			assert.Equal(t, printing.PrinterTypeBill, p.Type)
		}
		
		// Find by type LABEL
		labelPrinters, err := repo.FindByType(ctx, printing.PrinterTypeLabel)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(labelPrinters), 1, "Should find at least 1 label printer")
		for _, p := range labelPrinters {
			assert.Equal(t, printing.PrinterTypeLabel, p.Type)
		}
	})
	
	t.Run("FindDefault", func(t *testing.T) {
		// Clear any existing default bill printers first
		// (in case previous tests created them)
		
		// Create default bill printer
		defaultBill := &printing.PrinterConfig{
			Name:           "Default Bill Printer",
			Type:           printing.PrinterTypeBill,
			ConnectionType: printing.ConnectionTypeNetwork,
			IPAddress:      "192.168.1.200",
			Port:           9100,
			PaperWidth:     80,
			IsDefault:      true,
			IsEnabled:      true,
		}
		err := repo.Create(ctx, defaultBill)
		require.NoError(t, err)
		
		// Find default bill printer
		found, err := repo.FindDefault(ctx, printing.PrinterTypeBill)
		require.NoError(t, err, "Failed to find default bill printer")
		// Note: May find a different default if multiple exist from previous subtests
		// Just verify it's a default bill printer
		assert.True(t, found.IsDefault, "Should be marked as default")
		assert.Equal(t, printing.PrinterTypeBill, found.Type)
	})
	
	t.Run("Update", func(t *testing.T) {
		config := &printing.PrinterConfig{
			Name:           "Updatable Printer",
			Type:           printing.PrinterTypeBill,
			ConnectionType: printing.ConnectionTypeNetwork,
			IPAddress:      "192.168.1.150",
			Port:           9100,
			PaperWidth:     80,
			IsEnabled:      true,
		}
		err := repo.Create(ctx, config)
		require.NoError(t, err)
		
		// Update fields
		config.Name = "Updated Printer Name"
		config.IPAddress = "192.168.1.151"
		config.IsDefault = true
		
		err = repo.Update(ctx, config)
		require.NoError(t, err, "Failed to update printer config")
		
		// Verify update
		found, err := repo.FindByID(ctx, config.ID)
		require.NoError(t, err)
		assert.Equal(t, "Updated Printer Name", found.Name)
		assert.Equal(t, "192.168.1.151", found.IPAddress)
		assert.True(t, found.IsDefault)
	})
	
	t.Run("Delete", func(t *testing.T) {
		config := &printing.PrinterConfig{
			Name:           "Deletable Printer",
			Type:           printing.PrinterTypeBill,
			ConnectionType: printing.ConnectionTypeNetwork,
			IPAddress:      "192.168.1.160",
			Port:           9100,
			PaperWidth:     80,
			IsEnabled:      true,
		}
		err := repo.Create(ctx, config)
		require.NoError(t, err)
		
		// Delete
		err = repo.Delete(ctx, config.ID)
		require.NoError(t, err)
		
		// Verify deleted
		_, err = repo.FindByID(ctx, config.ID)
		assert.Error(t, err, "Should return error for deleted config")
		assert.Equal(t, mongo.ErrNoDocuments, err)
	})
}

// TestPrintTemplateRepository_Integration tests PrintTemplateRepository with real MongoDB
func TestPrintTemplateRepository_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	
	db, cleanup := setupTestDB(t)
	defer cleanup()
	
	ctx := context.Background()
	repo := NewPrintTemplateRepository(db)
	
	// Create indexes
	err := repo.CreateIndexes(ctx)
	require.NoError(t, err, "Failed to create indexes")
	
	t.Run("Create and FindByID", func(t *testing.T) {
		template := &printing.PrintTemplate{
			Type:      printing.TemplateTypeBill,
			Name:      "Default Bill Template",
			Content:   "{{.ShopName}}\n{{.OrderNumber}}",
			IsDefault: true,
		}
		
		// Create
		err := repo.Create(ctx, template)
		require.NoError(t, err, "Failed to create template")
		assert.False(t, template.ID.IsZero(), "ID should be set after creation")
		
		// FindByID
		found, err := repo.FindByID(ctx, template.ID)
		require.NoError(t, err, "Failed to find template")
		assert.Equal(t, template.Name, found.Name)
		assert.Equal(t, template.Content, found.Content)
		assert.Equal(t, template.Type, found.Type)
	})
	
	t.Run("FindByType", func(t *testing.T) {
		// Create bill template
		billTemplate := &printing.PrintTemplate{
			Type:    printing.TemplateTypeBill,
			Name:    "Custom Bill Template",
			Content: "Bill content",
		}
		err := repo.Create(ctx, billTemplate)
		require.NoError(t, err)
		
		// Create label template
		labelTemplate := &printing.PrintTemplate{
			Type:    printing.TemplateTypeLabel,
			Name:    "Custom Label Template",
			Content: "Label content",
		}
		err = repo.Create(ctx, labelTemplate)
		require.NoError(t, err)
		
		// Find by type BILL
		billTemplates, err := repo.FindByType(ctx, printing.TemplateTypeBill)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(billTemplates), 1, "Should find at least 1 bill template")
		for _, tmpl := range billTemplates {
			assert.Equal(t, printing.TemplateTypeBill, tmpl.Type)
		}
		
		// Find by type LABEL
		labelTemplates, err := repo.FindByType(ctx, printing.TemplateTypeLabel)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(labelTemplates), 1, "Should find at least 1 label template")
		for _, tmpl := range labelTemplates {
			assert.Equal(t, printing.TemplateTypeLabel, tmpl.Type)
		}
	})
	
	t.Run("FindDefault", func(t *testing.T) {
		// Create default label template
		defaultLabel := &printing.PrintTemplate{
			Type:      printing.TemplateTypeLabel,
			Name:      "Default Label Template",
			Content:   "{{.OrderNumber}}\n{{.ItemName}}",
			IsDefault: true,
		}
		err := repo.Create(ctx, defaultLabel)
		require.NoError(t, err)
		
		// Find default label template
		found, err := repo.FindDefault(ctx, printing.TemplateTypeLabel)
		require.NoError(t, err, "Failed to find default label template")
		assert.Equal(t, defaultLabel.Name, found.Name)
		assert.True(t, found.IsDefault)
		assert.Equal(t, printing.TemplateTypeLabel, found.Type)
	})
	
	t.Run("Update", func(t *testing.T) {
		template := &printing.PrintTemplate{
			Type:    printing.TemplateTypeBill,
			Name:    "Updatable Template",
			Content: "Original content",
		}
		err := repo.Create(ctx, template)
		require.NoError(t, err)
		
		// Update fields
		template.Name = "Updated Template Name"
		template.Content = "Updated content"
		template.IsDefault = true
		
		err = repo.Update(ctx, template)
		require.NoError(t, err, "Failed to update template")
		
		// Verify update
		found, err := repo.FindByID(ctx, template.ID)
		require.NoError(t, err)
		assert.Equal(t, "Updated Template Name", found.Name)
		assert.Equal(t, "Updated content", found.Content)
		assert.True(t, found.IsDefault)
	})
	
	t.Run("Delete", func(t *testing.T) {
		template := &printing.PrintTemplate{
			Type:    printing.TemplateTypeBill,
			Name:    "Deletable Template",
			Content: "Delete me",
		}
		err := repo.Create(ctx, template)
		require.NoError(t, err)
		
		// Delete
		err = repo.Delete(ctx, template.ID)
		require.NoError(t, err)
		
		// Verify deleted
		_, err = repo.FindByID(ctx, template.ID)
		assert.Error(t, err, "Should return error for deleted template")
		assert.Equal(t, mongo.ErrNoDocuments, err)
	})
}

// TestTTLIndex_Integration tests TTL index behavior for print_jobs
func TestTTLIndex_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	
	t.Skip("TTL index test requires waiting for MongoDB TTL monitor (runs every 60 seconds)")
	
	// This test is skipped because TTL index cleanup happens asynchronously
	// and requires waiting for MongoDB's TTL monitor thread to run
	// In production, jobs older than 7 days will be automatically deleted
	
	db, cleanup := setupTestDB(t)
	defer cleanup()
	
	ctx := context.Background()
	repo := NewPrintJobRepository(db)
	
	// Create indexes with TTL
	err := repo.CreateIndexes(ctx)
	require.NoError(t, err)
	
	// Create a job (in real scenario, would need to wait 7 days + TTL monitor interval)
	job := &printing.PrintJob{
		Type:        printing.PrintJobTypeBill,
		OrderID:     primitive.NewObjectID(),
		OrderNumber: "ORD-TTL",
		PrinterID:   primitive.NewObjectID(),
		Content:     "TTL test",
		Status:      printing.PrintJobStatusCompleted,
		MaxRetries:  3,
	}
	err = repo.Create(ctx, job)
	require.NoError(t, err)
	
	// Note: In production, this job would be automatically deleted after 7 days
	// MongoDB's TTL monitor runs every 60 seconds and deletes expired documents
}
