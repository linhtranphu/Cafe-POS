package services

import (
	"testing"

	"cafe-pos/backend/domain/printing"
	infraPrinting "cafe-pos/backend/infrastructure/printing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPrinterManager(t *testing.T) {
	pm := NewPrinterManager()
	assert.NotNil(t, pm)
}

func TestPrinterManager_GetPrinter(t *testing.T) {
	pm := NewPrinterManager()

	tests := []struct {
		name        string
		config      *printing.PrinterConfig
		expectError bool
		errorMsg    string
	}{
		{
			name:        "nil config should return error",
			config:      nil,
			expectError: true,
			errorMsg:    "printer config cannot be nil",
		},
		{
			name: "bill printer should return ESC/POS printer",
			config: &printing.PrinterConfig{
				Type:           printing.PrinterTypeBill,
				ConnectionType: printing.ConnectionTypeNetwork,
				IPAddress:      "192.168.1.100",
				Port:           9100,
				PaperWidth:     80,
			},
			expectError: false,
		},
		{
			name: "label printer should return label printer",
			config: &printing.PrinterConfig{
				Type:           printing.PrinterTypeLabel,
				ConnectionType: printing.ConnectionTypeNetwork,
				IPAddress:      "192.168.1.101",
				Port:           9100,
			},
			expectError: false,
		},
		{
			name: "unsupported printer type should return error",
			config: &printing.PrinterConfig{
				Type: "UNKNOWN",
			},
			expectError: true,
			errorMsg:    "unsupported printer type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			printer, err := pm.GetPrinter(tt.config)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
				assert.Nil(t, printer)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, printer)
			}
		})
	}
}

func TestPrinterManager_TestConnection(t *testing.T) {
	pm := NewPrinterManager()

	tests := []struct {
		name        string
		config      *printing.PrinterConfig
		expectError bool
		errorMsg    string
	}{
		{
			name:        "nil config should return error",
			config:      nil,
			expectError: true,
			errorMsg:    "printer config cannot be nil",
		},
		{
			name: "valid network printer config should fail without actual printer",
			config: &printing.PrinterConfig{
				Type:           printing.PrinterTypeBill,
				ConnectionType: printing.ConnectionTypeNetwork,
				IPAddress:      "192.168.1.100",
				Port:           9100,
				PaperWidth:     80,
			},
			expectError: true, // Will fail because no actual printer
			errorMsg:    "connection test failed",
		},
		{
			name: "USB printer should fail (not supported for ESC/POS)",
			config: &printing.PrinterConfig{
				Type:           printing.PrinterTypeBill,
				ConnectionType: printing.ConnectionTypeUSB,
				USBPath:        "/dev/usb/lp0",
				PaperWidth:     58,
			},
			expectError: true,
			errorMsg:    "only supports network connection",
		},
		{
			name: "network printer without IP should fail",
			config: &printing.PrinterConfig{
				Type:           printing.PrinterTypeBill,
				ConnectionType: printing.ConnectionTypeNetwork,
				Port:           9100,
				PaperWidth:     80,
			},
			expectError: true,
			errorMsg:    "IP address is required",
		},
		{
			name: "label printer with valid network config should fail without actual printer",
			config: &printing.PrinterConfig{
				Type:           printing.PrinterTypeLabel,
				ConnectionType: printing.ConnectionTypeNetwork,
				IPAddress:      "192.168.1.101",
				Port:           9100,
			},
			expectError: true, // Will fail because no actual printer
			errorMsg:    "connection test failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pm.TestConnection(tt.config)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestESCPOSPrinter_Connect(t *testing.T) {
	tests := []struct {
		name        string
		config      *printing.PrinterConfig
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid network config should fail without actual printer",
			config: &printing.PrinterConfig{
				Type:           printing.PrinterTypeBill,
				ConnectionType: printing.ConnectionTypeNetwork,
				IPAddress:      "192.168.1.100",
				Port:           9100,
			},
			expectError: true, // Will fail because no actual printer at this address
			errorMsg:    "failed to connect",
		},
		{
			name: "network config without IP should fail",
			config: &printing.PrinterConfig{
				Type:           printing.PrinterTypeBill,
				ConnectionType: printing.ConnectionTypeNetwork,
				Port:           9100,
			},
			expectError: true,
			errorMsg:    "IP address is required",
		},
		{
			name: "USB connection should fail (not supported)",
			config: &printing.PrinterConfig{
				Type:           printing.PrinterTypeBill,
				ConnectionType: printing.ConnectionTypeUSB,
				USBPath:        "/dev/usb/lp0",
			},
			expectError: true,
			errorMsg:    "only supports network connection",
		},
		{
			name: "network config without port should fail",
			config: &printing.PrinterConfig{
				Type:           printing.PrinterTypeBill,
				ConnectionType: printing.ConnectionTypeNetwork,
				IPAddress:      "192.168.1.100",
			},
			expectError: true,
			errorMsg:    "port is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			printer := infraPrinting.NewESCPOSPrinter(tt.config)
			err := printer.Connect()

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
				// Clean up connection
				printer.Disconnect()
			}
		})
	}
}

func TestESCPOSPrinter_Print(t *testing.T) {
	config := &printing.PrinterConfig{
		Type:           printing.PrinterTypeBill,
		ConnectionType: printing.ConnectionTypeNetwork,
		IPAddress:      "192.168.1.100",
		Port:           9100,
		PaperWidth:     80,
	}
	printer := infraPrinting.NewESCPOSPrinter(config)

	tests := []struct {
		name        string
		content     string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "print without connection should fail",
			content:     "Test print content",
			expectError: true,
			errorMsg:    "printer not connected",
		},
		{
			name:        "empty content should fail",
			content:     "",
			expectError: true,
			errorMsg:    "print content cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := printer.Print(tt.content)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestESCPOSPrinter_GetStatus(t *testing.T) {
	config := &printing.PrinterConfig{
		Type:           printing.PrinterTypeBill,
		ConnectionType: printing.ConnectionTypeNetwork,
		IPAddress:      "192.168.1.100",
		Port:           9100,
		PaperWidth:     80,
	}
	printer := infraPrinting.NewESCPOSPrinter(config)

	// Without connection, status should show offline
	status, err := printer.GetStatus()
	require.NoError(t, err)
	assert.False(t, status.IsOnline)
	assert.Equal(t, "UNKNOWN", status.PaperStatus)
	assert.Equal(t, "Not connected", status.ErrorMsg)
}

func TestESCPOSPrinter_Disconnect(t *testing.T) {
	config := &printing.PrinterConfig{
		Type:           printing.PrinterTypeBill,
		ConnectionType: printing.ConnectionTypeNetwork,
		IPAddress:      "192.168.1.100",
		Port:           9100,
		PaperWidth:     80,
	}
	printer := infraPrinting.NewESCPOSPrinter(config)

	// Disconnect without connection should succeed (no-op)
	err := printer.Disconnect()
	assert.NoError(t, err)
}

func TestLabelPrinter_Connect(t *testing.T) {
	tests := []struct {
		name        string
		config      *printing.PrinterConfig
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid network config should fail without actual printer",
			config: &printing.PrinterConfig{
				Type:           printing.PrinterTypeLabel,
				ConnectionType: printing.ConnectionTypeNetwork,
				IPAddress:      "192.168.1.101",
				Port:           9100,
			},
			expectError: true, // Will fail because no actual printer at this address
			errorMsg:    "failed to connect",
		},
		{
			name: "network config without IP should fail",
			config: &printing.PrinterConfig{
				Type:           printing.PrinterTypeLabel,
				ConnectionType: printing.ConnectionTypeNetwork,
				Port:           9100,
			},
			expectError: true,
			errorMsg:    "IP address is required",
		},
		{
			name: "USB connection should fail (not supported)",
			config: &printing.PrinterConfig{
				Type:           printing.PrinterTypeLabel,
				ConnectionType: printing.ConnectionTypeUSB,
				USBPath:        "/dev/usb/lp0",
			},
			expectError: true,
			errorMsg:    "only supports network connection",
		},
		{
			name: "network config without port should fail",
			config: &printing.PrinterConfig{
				Type:           printing.PrinterTypeLabel,
				ConnectionType: printing.ConnectionTypeNetwork,
				IPAddress:      "192.168.1.101",
			},
			expectError: true,
			errorMsg:    "port is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			printer := infraPrinting.NewLabelPrinter(tt.config)
			err := printer.Connect()

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
				// Clean up connection
				printer.Disconnect()
			}
		})
	}
}

func TestLabelPrinter_Print(t *testing.T) {
	config := &printing.PrinterConfig{
		Type:           printing.PrinterTypeLabel,
		ConnectionType: printing.ConnectionTypeNetwork,
		IPAddress:      "192.168.1.101",
		Port:           9100,
		PaperWidth:     50, // 50x30mm label
	}
	printer := infraPrinting.NewLabelPrinter(config)

	tests := []struct {
		name        string
		content     string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "print without connection should fail",
			content:     "Label content",
			expectError: true,
			errorMsg:    "printer not connected",
		},
		{
			name:        "empty content should fail",
			content:     "",
			expectError: true,
			errorMsg:    "print content cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := printer.Print(tt.content)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestLabelPrinter_GetStatus(t *testing.T) {
	config := &printing.PrinterConfig{
		Type:           printing.PrinterTypeLabel,
		ConnectionType: printing.ConnectionTypeNetwork,
		IPAddress:      "192.168.1.101",
		Port:           9100,
		PaperWidth:     50,
	}
	printer := infraPrinting.NewLabelPrinter(config)

	// Without connection, status should show offline
	status, err := printer.GetStatus()
	require.NoError(t, err)
	assert.False(t, status.IsOnline)
	assert.Equal(t, "UNKNOWN", status.PaperStatus)
	assert.Equal(t, "Not connected", status.ErrorMsg)
}

func TestLabelPrinter_Disconnect(t *testing.T) {
	config := &printing.PrinterConfig{
		Type:           printing.PrinterTypeLabel,
		ConnectionType: printing.ConnectionTypeNetwork,
		IPAddress:      "192.168.1.101",
		Port:           9100,
		PaperWidth:     50,
	}
	printer := infraPrinting.NewLabelPrinter(config)

	// Disconnect without connection should succeed (no-op)
	err := printer.Disconnect()
	assert.NoError(t, err)
}
