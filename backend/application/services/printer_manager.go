package services

import (
	"fmt"

	"cafe-pos/backend/domain/printing"
	infraPrinting "cafe-pos/backend/infrastructure/printing"
	"cafe-pos/backend/infrastructure/printbridge"
)

// PrinterStatus represents the status of a printer
type PrinterStatus = infraPrinting.PrinterStatus

// Printer defines the interface for printer operations
type Printer = infraPrinting.Printer

// PrinterManager defines the interface for managing printers
type PrinterManager interface {
	GetPrinter(config *printing.PrinterConfig) (Printer, error)
	TestConnection(config *printing.PrinterConfig) error
	SetPrintBridgeClient(client *printbridge.Client)
}

// printerManager implements the PrinterManager interface
type printerManager struct{
	printBridgeClient *printbridge.Client
}

// NewPrinterManager creates a new printer manager instance
func NewPrinterManager() PrinterManager {
	return &printerManager{
		printBridgeClient: nil,
	}
}

// SetPrintBridgeClient sets the print bridge client for using bridge mode
func (pm *printerManager) SetPrintBridgeClient(client *printbridge.Client) {
	pm.printBridgeClient = client
}

// GetPrinter returns a printer instance based on the configuration (factory pattern)
func (pm *printerManager) GetPrinter(config *printing.PrinterConfig) (Printer, error) {
	if config == nil {
		return nil, fmt.Errorf("printer config cannot be nil")
	}

	// If print bridge client is configured and available, use bridge printer directly
	// No need to create innerPrinter (which requires fonts, etc.)
	if pm.printBridgeClient != nil && pm.printBridgeClient.IsAvailable() {
		return infraPrinting.NewBridgePrinter(config, pm.printBridgeClient, nil), nil
	}

	// No print bridge - create direct printer
	var innerPrinter Printer
	var err error

	// Factory pattern: create appropriate printer based on type
	switch config.Type {
	case printing.PrinterTypeBill:
		innerPrinter, err = infraPrinting.NewESCPOSPrinter(config)
		if err != nil {
			return nil, fmt.Errorf("failed to create ESC/POS printer: %w", err)
		}
	case printing.PrinterTypeLabel:
		innerPrinter = infraPrinting.NewLabelPrinter(config)
	default:
		return nil, fmt.Errorf("unsupported printer type: %s", config.Type)
	}

	return innerPrinter, nil
}

// TestConnection tests if a printer is reachable with the given configuration
func (pm *printerManager) TestConnection(config *printing.PrinterConfig) error {
	if config == nil {
		return fmt.Errorf("printer config cannot be nil")
	}

	// Get printer instance
	printer, err := pm.GetPrinter(config)
	if err != nil {
		return fmt.Errorf("failed to get printer: %w", err)
	}

	// Try to connect
	if err := printer.Connect(); err != nil {
		return fmt.Errorf("connection test failed: %w", err)
	}

	// Disconnect after test
	if err := printer.Disconnect(); err != nil {
		return fmt.Errorf("failed to disconnect after test: %w", err)
	}

	return nil
}
