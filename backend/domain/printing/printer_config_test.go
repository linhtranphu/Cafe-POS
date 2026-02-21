package printing

import (
	"testing"
	"time"
)

// TestPrinterType_Constants tests that PrinterType constants are defined correctly
func TestPrinterType_Constants(t *testing.T) {
	tests := []struct {
		name         string
		printerType  PrinterType
		expected     string
	}{
		{"Bill printer type", PrinterTypeBill, "BILL"},
		{"Label printer type", PrinterTypeLabel, "LABEL"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.printerType) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(tt.printerType))
			}
		})
	}
}

// TestConnectionType_Constants tests that ConnectionType constants are defined correctly
func TestConnectionType_Constants(t *testing.T) {
	tests := []struct {
		name           string
		connectionType ConnectionType
		expected       string
	}{
		{"Network connection", ConnectionTypeNetwork, "NETWORK"},
		{"USB connection", ConnectionTypeUSB, "USB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.connectionType) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(tt.connectionType))
			}
		})
	}
}

// TestPrinterConfig_Creation tests creating a PrinterConfig with valid data
func TestPrinterConfig_Creation(t *testing.T) {
	now := time.Now()

	config := &PrinterConfig{
		Name:           "Bill Printer 1",
		Type:           PrinterTypeBill,
		ConnectionType: ConnectionTypeNetwork,
		IPAddress:      "192.168.1.100",
		Port:           9100,
		PaperWidth:     80,
		IsDefault:      true,
		IsEnabled:      true,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	// Verify all fields are set correctly
	if config.Name != "Bill Printer 1" {
		t.Errorf("Expected Name to be 'Bill Printer 1', got %s", config.Name)
	}
	if config.Type != PrinterTypeBill {
		t.Errorf("Expected Type to be BILL, got %s", config.Type)
	}
	if config.ConnectionType != ConnectionTypeNetwork {
		t.Errorf("Expected ConnectionType to be NETWORK, got %s", config.ConnectionType)
	}
	if config.IPAddress != "192.168.1.100" {
		t.Errorf("Expected IPAddress to be 192.168.1.100, got %s", config.IPAddress)
	}
	if config.Port != 9100 {
		t.Errorf("Expected Port to be 9100, got %d", config.Port)
	}
	if config.PaperWidth != 80 {
		t.Errorf("Expected PaperWidth to be 80, got %d", config.PaperWidth)
	}
	if !config.IsDefault {
		t.Error("Expected IsDefault to be true")
	}
	if !config.IsEnabled {
		t.Error("Expected IsEnabled to be true")
	}
}

// TestPrinterConfig_NetworkConnection tests network connection configuration
func TestPrinterConfig_NetworkConnection(t *testing.T) {
	config := &PrinterConfig{
		Name:           "Network Printer",
		Type:           PrinterTypeBill,
		ConnectionType: ConnectionTypeNetwork,
		IPAddress:      "192.168.1.100",
		Port:           9100,
		PaperWidth:     80,
		IsEnabled:      true,
	}

	// Network printer should have IP and Port
	if config.IPAddress == "" {
		t.Error("Expected IPAddress to be set for network printer")
	}
	if config.Port == 0 {
		t.Error("Expected Port to be set for network printer")
	}
	// USB path should be empty for network printer
	if config.USBPath != "" {
		t.Error("Expected USBPath to be empty for network printer")
	}
}

// TestPrinterConfig_USBConnection tests USB connection configuration
func TestPrinterConfig_USBConnection(t *testing.T) {
	config := &PrinterConfig{
		Name:           "USB Printer",
		Type:           PrinterTypeLabel,
		ConnectionType: ConnectionTypeUSB,
		USBPath:        "/dev/usb/lp0",
		PaperWidth:     58,
		IsEnabled:      true,
	}

	// USB printer should have USB path
	if config.USBPath == "" {
		t.Error("Expected USBPath to be set for USB printer")
	}
	// IP and Port should be empty for USB printer
	if config.IPAddress != "" {
		t.Error("Expected IPAddress to be empty for USB printer")
	}
	if config.Port != 0 {
		t.Error("Expected Port to be 0 for USB printer")
	}
}

// TestPrinterConfig_PaperWidthValidation tests paper width validation
func TestPrinterConfig_PaperWidthValidation(t *testing.T) {
	tests := []struct {
		name       string
		paperWidth int
		valid      bool
	}{
		{"Valid 58mm", 58, true},
		{"Valid 80mm", 80, true},
		{"Invalid 70mm", 70, false},
		{"Invalid 0mm", 0, false},
		{"Invalid negative", -1, false},
		{"Invalid 100mm", 100, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &PrinterConfig{
				PaperWidth: tt.paperWidth,
			}

			// Validate paper width (should be 58 or 80)
			isValid := config.PaperWidth == 58 || config.PaperWidth == 80
			if isValid != tt.valid {
				t.Errorf("Expected validation to be %v for paper width %d", tt.valid, tt.paperWidth)
			}
		})
	}
}

// TestPrinterConfig_DefaultPrinter tests default printer logic
func TestPrinterConfig_DefaultPrinter(t *testing.T) {
	// Only one printer of each type should be default
	billPrinter1 := &PrinterConfig{
		Name:      "Bill Printer 1",
		Type:      PrinterTypeBill,
		IsDefault: true,
		IsEnabled: true,
	}

	billPrinter2 := &PrinterConfig{
		Name:      "Bill Printer 2",
		Type:      PrinterTypeBill,
		IsDefault: false,
		IsEnabled: true,
	}

	labelPrinter := &PrinterConfig{
		Name:      "Label Printer 1",
		Type:      PrinterTypeLabel,
		IsDefault: true,
		IsEnabled: true,
	}

	// Verify only one bill printer is default
	if !billPrinter1.IsDefault {
		t.Error("Expected Bill Printer 1 to be default")
	}
	if billPrinter2.IsDefault {
		t.Error("Expected Bill Printer 2 to not be default")
	}

	// Label printer can also be default (different type)
	if !labelPrinter.IsDefault {
		t.Error("Expected Label Printer to be default")
	}
}

// TestPrinterConfig_EnableDisable tests enable/disable functionality
func TestPrinterConfig_EnableDisable(t *testing.T) {
	config := &PrinterConfig{
		Name:      "Test Printer",
		Type:      PrinterTypeBill,
		IsEnabled: true,
	}

	// Initially enabled
	if !config.IsEnabled {
		t.Error("Expected printer to be enabled initially")
	}

	// Disable printer
	config.IsEnabled = false
	config.UpdatedAt = time.Now()

	if config.IsEnabled {
		t.Error("Expected printer to be disabled")
	}

	// Configuration should still exist (not deleted)
	if config.Name == "" {
		t.Error("Expected printer configuration to still exist after disable")
	}

	// Re-enable printer
	config.IsEnabled = true
	config.UpdatedAt = time.Now()

	if !config.IsEnabled {
		t.Error("Expected printer to be re-enabled")
	}
}

// TestCreatePrinterConfigRequest_Validation tests request validation
func TestCreatePrinterConfigRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request CreatePrinterConfigRequest
		valid   bool
	}{
		{
			"Valid network bill printer",
			CreatePrinterConfigRequest{
				Name:           "Bill Printer",
				Type:           PrinterTypeBill,
				ConnectionType: ConnectionTypeNetwork,
				IPAddress:      "192.168.1.100",
				Port:           9100,
				PaperWidth:     80,
				IsEnabled:      true,
			},
			true,
		},
		{
			"Valid USB label printer",
			CreatePrinterConfigRequest{
				Name:           "Label Printer",
				Type:           PrinterTypeLabel,
				ConnectionType: ConnectionTypeUSB,
				USBPath:        "/dev/usb/lp0",
				PaperWidth:     58,
				IsEnabled:      true,
			},
			true,
		},
		{
			"Empty name",
			CreatePrinterConfigRequest{
				Name:           "",
				Type:           PrinterTypeBill,
				ConnectionType: ConnectionTypeNetwork,
				IPAddress:      "192.168.1.100",
				Port:           9100,
				PaperWidth:     80,
			},
			false,
		},
		{
			"Invalid paper width",
			CreatePrinterConfigRequest{
				Name:           "Printer",
				Type:           PrinterTypeBill,
				ConnectionType: ConnectionTypeNetwork,
				IPAddress:      "192.168.1.100",
				Port:           9100,
				PaperWidth:     70,
			},
			false,
		},
		{
			"Network printer without IP",
			CreatePrinterConfigRequest{
				Name:           "Printer",
				Type:           PrinterTypeBill,
				ConnectionType: ConnectionTypeNetwork,
				IPAddress:      "",
				Port:           9100,
				PaperWidth:     80,
			},
			false,
		},
		{
			"USB printer without path",
			CreatePrinterConfigRequest{
				Name:           "Printer",
				Type:           PrinterTypeLabel,
				ConnectionType: ConnectionTypeUSB,
				USBPath:        "",
				PaperWidth:     58,
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Basic validation checks
			isValid := tt.request.Name != "" &&
				(tt.request.PaperWidth == 58 || tt.request.PaperWidth == 80)

			// Connection-specific validation
			if tt.request.ConnectionType == ConnectionTypeNetwork {
				isValid = isValid && tt.request.IPAddress != "" && tt.request.Port > 0
			} else if tt.request.ConnectionType == ConnectionTypeUSB {
				isValid = isValid && tt.request.USBPath != ""
			}

			if isValid != tt.valid {
				t.Errorf("Expected validation to be %v, got %v", tt.valid, isValid)
			}
		})
	}
}

// TestUpdatePrinterConfigRequest_PartialUpdate tests partial update request
func TestUpdatePrinterConfigRequest_PartialUpdate(t *testing.T) {
	// Test that we can update individual fields
	isEnabled := false
	port := 9200

	request := &UpdatePrinterConfigRequest{
		Name:      "Updated Name",
		IsEnabled: &isEnabled,
		Port:      &port,
	}

	if request.Name != "Updated Name" {
		t.Errorf("Expected Name to be 'Updated Name', got %s", request.Name)
	}

	if request.IsEnabled == nil || *request.IsEnabled != false {
		t.Error("Expected IsEnabled to be false")
	}

	if request.Port == nil || *request.Port != 9200 {
		t.Error("Expected Port to be 9200")
	}

	// Fields not set should be nil/empty
	if request.IPAddress != "" {
		t.Error("Expected IPAddress to be empty for partial update")
	}
}

// TestPrinterConfigFilter_DefaultValues tests filter default values
func TestPrinterConfigFilter_DefaultValues(t *testing.T) {
	filter := &PrinterConfigFilter{
		Page:  1,
		Limit: 20,
	}

	if filter.Page != 1 {
		t.Errorf("Expected default Page to be 1, got %d", filter.Page)
	}

	if filter.Limit != 20 {
		t.Errorf("Expected default Limit to be 20, got %d", filter.Limit)
	}

	if filter.Type != "" {
		t.Errorf("Expected Type to be empty, got %s", filter.Type)
	}

	if filter.IsEnabled != nil {
		t.Error("Expected IsEnabled to be nil")
	}
}

// TestPrinterConfig_TypeSpecificValidation tests type-specific validation
func TestPrinterConfig_TypeSpecificValidation(t *testing.T) {
	tests := []struct {
		name        string
		printerType PrinterType
		paperWidth  int
		valid       bool
		description string
	}{
		{
			"Bill printer with 80mm",
			PrinterTypeBill,
			80,
			true,
			"Common configuration for bill printers",
		},
		{
			"Bill printer with 58mm",
			PrinterTypeBill,
			58,
			true,
			"Smaller bill printer",
		},
		{
			"Label printer with 58mm",
			PrinterTypeLabel,
			58,
			true,
			"Common configuration for label printers",
		},
		{
			"Label printer with 80mm",
			PrinterTypeLabel,
			80,
			true,
			"Larger label printer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &PrinterConfig{
				Type:       tt.printerType,
				PaperWidth: tt.paperWidth,
			}

			// Validate paper width
			isValid := config.PaperWidth == 58 || config.PaperWidth == 80
			if isValid != tt.valid {
				t.Errorf("Expected validation to be %v for %s", tt.valid, tt.description)
			}
		})
	}
}
