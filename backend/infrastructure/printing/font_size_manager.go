package printing

import (
	"strings"
)

// FontSizeManager manages font sizes for different types of content in bills
type FontSizeManager struct {
	normalSize float64 // 18pt - for regular content
	headerSize float64 // 22pt - for headers (shop name, "HÓA ĐƠN BÁN HÀNG")
	totalSize  float64 // 20pt - for total line ("TỔNG CỘNG")
}

// FontSizeConfig defines font size and weight for a line
type FontSizeConfig struct {
	Size float64
	Bold bool
}

// NewFontSizeManager creates a new FontSizeManager with default sizes
func NewFontSizeManager() *FontSizeManager {
	return &FontSizeManager{
		normalSize: 18.0,
		headerSize: 22.0,
		totalSize:  20.0,
	}
}

// NewFontSizeManagerWithSizes creates a FontSizeManager with custom sizes
func NewFontSizeManagerWithSizes(normalSize, headerSize, totalSize float64) *FontSizeManager {
	// Validate sizes - use defaults if invalid
	if normalSize <= 0 {
		normalSize = 18.0
	}
	if headerSize <= 0 {
		headerSize = 22.0
	}
	if totalSize <= 0 {
		totalSize = 20.0
	}

	return &FontSizeManager{
		normalSize: normalSize,
		headerSize: headerSize,
		totalSize:  totalSize,
	}
}

// GetFontSizeForLine determines the appropriate font size and weight for a line
// based on its content
func (m *FontSizeManager) GetFontSizeForLine(line string) FontSizeConfig {
	trimmed := strings.TrimSpace(line)
	
	// Empty lines use normal size
	if trimmed == "" {
		return FontSizeConfig{
			Size: m.normalSize,
			Bold: false,
		}
	}

	upper := strings.ToUpper(trimmed)

	// Check for header lines (22pt, bold)
	if m.isHeader(upper, trimmed) {
		return FontSizeConfig{
			Size: m.headerSize,
			Bold: true,
		}
	}

	// Check for total line (20pt, bold)
	if m.isTotalLine(upper) {
		return FontSizeConfig{
			Size: m.totalSize,
			Bold: true,
		}
	}

	// Check for table header (18pt, bold)
	if m.isTableHeader(trimmed) {
		return FontSizeConfig{
			Size: m.normalSize,
			Bold: true,
		}
	}

	// Check for other bold lines (18pt, bold)
	if m.shouldBeBold(upper) {
		return FontSizeConfig{
			Size: m.normalSize,
			Bold: true,
		}
	}

	// Default: regular content (18pt, normal)
	return FontSizeConfig{
		Size: m.normalSize,
		Bold: false,
	}
}

// isHeader checks if a line is a header (shop name or "HÓA ĐƠN BÁN HÀNG")
func (m *FontSizeManager) isHeader(upper, original string) bool {
	// "HÓA ĐƠN BÁN HÀNG" or variations
	if strings.Contains(upper, "HÓA ĐƠN") {
		return true
	}
	if strings.Contains(upper, "INVOICE") {
		return true
	}
	if strings.Contains(upper, "RECEIPT") {
		return true
	}

	// Exclude thank you messages and footer content
	if strings.Contains(upper, "CẢM ƠN") || strings.Contains(upper, "THANK") {
		return false
	}
	if strings.Contains(upper, "HẸN GẶP") || strings.Contains(upper, "WELCOME") {
		return false
	}
	if strings.Contains(upper, "VISIT") || strings.Contains(upper, "CHÀO") {
		return false
	}

	// Shop name detection: short lines without colons at the top
	// (typically the first non-empty, non-separator line)
	if !strings.Contains(original, ":") && 
	   len(original) < 40 && 
	   len(original) > 0 &&
	   !m.isSeparator(original) {
		// Additional check: not a table header
		if !strings.Contains(upper, "TÊN MÓN") && 
		   !strings.Contains(upper, "SL") &&
		   !strings.Contains(upper, "ĐƠN GIÁ") &&
		   !strings.Contains(upper, "THÀNH TIỀN") {
			return true
		}
	}

	return false
}

// isTotalLine checks if a line is the total line ("TỔNG CỘNG")
func (m *FontSizeManager) isTotalLine(upper string) bool {
	// "TỔNG CỘNG" or "TOTAL" (final total)
	if strings.Contains(upper, "TỔNG CỘNG") {
		return true
	}
	if strings.HasPrefix(upper, "TOTAL:") || strings.HasPrefix(upper, "TOTAL ") {
		return true
	}
	if strings.HasPrefix(upper, "GRAND TOTAL") {
		return true
	}

	return false
}

// isTableHeader checks if a line is a table header
func (m *FontSizeManager) isTableHeader(line string) bool {
	upper := strings.ToUpper(line)
	
	// Check for common table header patterns
	hasItemColumn := strings.Contains(upper, "TÊN MÓN") || 
	                 strings.Contains(upper, "ITEM") ||
	                 strings.Contains(upper, "NAME")
	
	hasQtyColumn := strings.Contains(upper, "SL") || 
	                strings.Contains(upper, "QTY") ||
	                strings.Contains(upper, "QUANTITY")
	
	hasPriceColumn := strings.Contains(upper, "ĐƠN GIÁ") || 
	                  strings.Contains(upper, "PRICE") ||
	                  strings.Contains(upper, "THÀNH TIỀN") ||
	                  strings.Contains(upper, "TOTAL")

	// Table header should have at least 2 of these columns
	count := 0
	if hasItemColumn {
		count++
	}
	if hasQtyColumn {
		count++
	}
	if hasPriceColumn {
		count++
	}

	return count >= 2
}

// shouldBeBold checks if a line should be bold (but not header or total)
func (m *FontSizeManager) shouldBeBold(upper string) bool {
	// Subtotal, discount, tax lines
	if strings.Contains(upper, "SUBTOTAL") || strings.Contains(upper, "TỔNG PHỤ") {
		return true
	}
	if strings.Contains(upper, "DISCOUNT") || strings.Contains(upper, "GIẢM GIÁ") {
		return true
	}
	if strings.Contains(upper, "TAX") || strings.Contains(upper, "THUẾ") {
		return true
	}
	if strings.Contains(upper, "TỔNG TIỀN") && !strings.Contains(upper, "TỔNG CỘNG") {
		return true
	}

	return false
}

// isSeparator checks if a line is a separator (===, ---)
func (m *FontSizeManager) isSeparator(line string) bool {
	if len(line) == 0 {
		return false
	}

	trimmed := strings.TrimSpace(line)
	if len(trimmed) == 0 {
		return false
	}

	// Check if line consists only of = or - characters
	firstChar := rune(trimmed[0])
	if firstChar != '=' && firstChar != '-' {
		return false
	}

	// All characters must be the same as the first character
	for _, ch := range trimmed {
		if ch != firstChar {
			return false
		}
	}

	return true
}
