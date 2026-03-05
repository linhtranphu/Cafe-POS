package printing

import (
	"testing"
)

func TestNewFontSizeManager(t *testing.T) {
	manager := NewFontSizeManager()

	if manager.normalSize != 18.0 {
		t.Errorf("Expected normalSize to be 18.0, got %f", manager.normalSize)
	}
	if manager.headerSize != 22.0 {
		t.Errorf("Expected headerSize to be 22.0, got %f", manager.headerSize)
	}
	if manager.totalSize != 20.0 {
		t.Errorf("Expected totalSize to be 20.0, got %f", manager.totalSize)
	}
}

func TestNewFontSizeManagerWithSizes(t *testing.T) {
	tests := []struct {
		name       string
		normal     float64
		header     float64
		total      float64
		wantNormal float64
		wantHeader float64
		wantTotal  float64
	}{
		{
			name:       "Valid custom sizes",
			normal:     16.0,
			header:     20.0,
			total:      18.0,
			wantNormal: 16.0,
			wantHeader: 20.0,
			wantTotal:  18.0,
		},
		{
			name:       "Invalid sizes - use defaults",
			normal:     -1.0,
			header:     0.0,
			total:      -5.0,
			wantNormal: 18.0,
			wantHeader: 22.0,
			wantTotal:  20.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := NewFontSizeManagerWithSizes(tt.normal, tt.header, tt.total)
			if manager.normalSize != tt.wantNormal {
				t.Errorf("normalSize = %f, want %f", manager.normalSize, tt.wantNormal)
			}
			if manager.headerSize != tt.wantHeader {
				t.Errorf("headerSize = %f, want %f", manager.headerSize, tt.wantHeader)
			}
			if manager.totalSize != tt.wantTotal {
				t.Errorf("totalSize = %f, want %f", manager.totalSize, tt.wantTotal)
			}
		})
	}
}

func TestGetFontSizeForLine_Headers(t *testing.T) {
	manager := NewFontSizeManager()

	tests := []struct {
		name     string
		line     string
		wantSize float64
		wantBold bool
	}{
		{
			name:     "HÓA ĐƠN BÁN HÀNG",
			line:     "HÓA ĐƠN BÁN HÀNG",
			wantSize: 22.0,
			wantBold: true,
		},
		{
			name:     "Shop name",
			line:     "Cafe ABC",
			wantSize: 22.0,
			wantBold: true,
		},
		{
			name:     "INVOICE",
			line:     "INVOICE",
			wantSize: 22.0,
			wantBold: true,
		},
		{
			name:     "RECEIPT",
			line:     "RECEIPT",
			wantSize: 22.0,
			wantBold: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := manager.GetFontSizeForLine(tt.line)
			if config.Size != tt.wantSize {
				t.Errorf("Size = %f, want %f", config.Size, tt.wantSize)
			}
			if config.Bold != tt.wantBold {
				t.Errorf("Bold = %v, want %v", config.Bold, tt.wantBold)
			}
		})
	}
}

func TestGetFontSizeForLine_TotalLine(t *testing.T) {
	manager := NewFontSizeManager()

	tests := []struct {
		name     string
		line     string
		wantSize float64
		wantBold bool
	}{
		{
			name:     "TỔNG CỘNG",
			line:     "TỔNG CỘNG: 100,000 VND",
			wantSize: 20.0,
			wantBold: true,
		},
		{
			name:     "TOTAL",
			line:     "TOTAL: 100,000 VND",
			wantSize: 20.0,
			wantBold: true,
		},
		{
			name:     "GRAND TOTAL",
			line:     "GRAND TOTAL: 100,000 VND",
			wantSize: 20.0,
			wantBold: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := manager.GetFontSizeForLine(tt.line)
			if config.Size != tt.wantSize {
				t.Errorf("Size = %f, want %f", config.Size, tt.wantSize)
			}
			if config.Bold != tt.wantBold {
				t.Errorf("Bold = %v, want %v", config.Bold, tt.wantBold)
			}
		})
	}
}

func TestGetFontSizeForLine_TableHeader(t *testing.T) {
	manager := NewFontSizeManager()

	tests := []struct {
		name     string
		line     string
		wantSize float64
		wantBold bool
	}{
		{
			name:     "Vietnamese table header",
			line:     "Tên món              SL  Đơn giá    Thành tiền",
			wantSize: 18.0,
			wantBold: true,
		},
		{
			name:     "English table header",
			line:     "Item Name            Qty  Price      Total",
			wantSize: 18.0,
			wantBold: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := manager.GetFontSizeForLine(tt.line)
			if config.Size != tt.wantSize {
				t.Errorf("Size = %f, want %f", config.Size, tt.wantSize)
			}
			if config.Bold != tt.wantBold {
				t.Errorf("Bold = %v, want %v", config.Bold, tt.wantBold)
			}
		})
	}
}

func TestGetFontSizeForLine_RegularContent(t *testing.T) {
	manager := NewFontSizeManager()

	tests := []struct {
		name     string
		line     string
		wantSize float64
		wantBold bool
	}{
		{
			name:     "Order info",
			line:     "Order: #12345",
			wantSize: 18.0,
			wantBold: false,
		},
		{
			name:     "Date",
			line:     "Ngày: 01/01/2024 10:30",
			wantSize: 18.0,
			wantBold: false,
		},
		{
			name:     "Table number",
			line:     "Bàn: 5",
			wantSize: 18.0,
			wantBold: false,
		},
		{
			name:     "Item line",
			line:     "Cafe Latte            2   45,000       90,000",
			wantSize: 18.0,
			wantBold: false,
		},
		{
			name:     "Variant line",
			line:     "  (Variant: Đặc biệt)",
			wantSize: 18.0,
			wantBold: false,
		},
		{
			name:     "Thank you message",
			line:     "Cảm ơn quý khách!",
			wantSize: 18.0,
			wantBold: false,
		},
		{
			name:     "Empty line",
			line:     "",
			wantSize: 18.0,
			wantBold: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := manager.GetFontSizeForLine(tt.line)
			if config.Size != tt.wantSize {
				t.Errorf("Size = %f, want %f", config.Size, tt.wantSize)
			}
			if config.Bold != tt.wantBold {
				t.Errorf("Bold = %v, want %v", config.Bold, tt.wantBold)
			}
		})
	}
}

func TestGetFontSizeForLine_BoldLines(t *testing.T) {
	manager := NewFontSizeManager()

	tests := []struct {
		name     string
		line     string
		wantSize float64
		wantBold bool
	}{
		{
			name:     "Subtotal",
			line:     "Tổng tiền: 90,000 VND",
			wantSize: 18.0,
			wantBold: true,
		},
		{
			name:     "Discount",
			line:     "Giảm giá: -10,000 VND",
			wantSize: 18.0,
			wantBold: true,
		},
		{
			name:     "Tax",
			line:     "Thuế: 5,000 VND",
			wantSize: 18.0,
			wantBold: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := manager.GetFontSizeForLine(tt.line)
			if config.Size != tt.wantSize {
				t.Errorf("Size = %f, want %f", config.Size, tt.wantSize)
			}
			if config.Bold != tt.wantBold {
				t.Errorf("Bold = %v, want %v", config.Bold, tt.wantBold)
			}
		})
	}
}

func TestGetFontSizeForLine_Separators(t *testing.T) {
	manager := NewFontSizeManager()

	tests := []struct {
		name     string
		line     string
		wantSize float64
		wantBold bool
	}{
		{
			name:     "Equals separator",
			line:     "================================",
			wantSize: 18.0,
			wantBold: false,
		},
		{
			name:     "Dash separator",
			line:     "------------------------------------------------",
			wantSize: 18.0,
			wantBold: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := manager.GetFontSizeForLine(tt.line)
			if config.Size != tt.wantSize {
				t.Errorf("Size = %f, want %f", config.Size, tt.wantSize)
			}
			if config.Bold != tt.wantBold {
				t.Errorf("Bold = %v, want %v", config.Bold, tt.wantBold)
			}
		})
	}
}
