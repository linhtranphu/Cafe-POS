package printing

import (
	"strings"
)

// Alignment represents text alignment options
type Alignment int

const (
	AlignLeft Alignment = iota
	AlignCenter
	AlignRight
)

// LineFormat represents the formatting information for a single line of text
type LineFormat struct {
	Text        string
	Bold        bool
	Alignment   Alignment
	IsSeparator bool
	FontSize    float64 // Font size in points (18pt, 20pt, 22pt)
	IsTableRow  bool    // True if line is part of a table
}

// FormatParser parses text content and identifies formatting requirements for each line
type FormatParser struct {
	paperWidth int
}

// NewFormatParser creates a new FormatParser instance
func NewFormatParser(paperWidth int) *FormatParser {
	return &FormatParser{
		paperWidth: paperWidth,
	}
}

// Parse parses text content and returns formatted lines
func (p *FormatParser) Parse(content string) []LineFormat {
	lines := strings.Split(content, "\n")
	result := make([]LineFormat, 0, len(lines))

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		
		lineFormat := LineFormat{
			Text:        line,
			Bold:        p.detectBold(trimmed),
			Alignment:   p.detectAlignment(trimmed),
			IsSeparator: p.isSeparator(trimmed),
			FontSize:    p.detectFontSize(trimmed),
			IsTableRow:  p.isTableRow(trimmed),
		}
		
		result = append(result, lineFormat)
	}

	return result
}

// detectAlignment determines alignment based on line content
func (p *FormatParser) detectAlignment(line string) Alignment {
	if line == "" {
		return AlignLeft
	}

	// Separator lines are centered
	if p.isSeparator(line) {
		return AlignCenter
	}

	upper := strings.ToUpper(line)

	// Header lines (short lines without colons) are centered
	if !strings.Contains(line, ":") && len(line) < 30 {
		return AlignCenter
	}

	// Total lines are centered
	if strings.HasPrefix(upper, "TOTAL") || strings.HasPrefix(upper, "TỔNG") {
		return AlignCenter
	}

	// Thank you messages are centered
	if strings.Contains(upper, "THANK") || strings.Contains(upper, "CẢM ƠN") {
		return AlignCenter
	}

	// Footer messages are centered
	if strings.Contains(upper, "VISIT") || strings.Contains(upper, "WELCOME") || 
	   strings.Contains(upper, "HẸN GẶP") || strings.Contains(upper, "CHÀO") {
		return AlignCenter
	}

	// Default to left alignment for item lines and other content
	return AlignLeft
}

// detectBold determines if line should be bold
func (p *FormatParser) detectBold(line string) bool {
	if line == "" {
		return false
	}

	upper := strings.ToUpper(line)

	// Header lines (short lines without colons) are bold
	if !strings.Contains(line, ":") && len(line) < 30 && len(line) > 0 {
		return true
	}

	// Lines with TOTAL, SUBTOTAL, DISCOUNT are bold
	if strings.Contains(upper, "TOTAL") || strings.Contains(upper, "TỔNG") {
		return true
	}

	if strings.Contains(upper, "SUBTOTAL") || strings.Contains(upper, "TỔNG PHỤ") {
		return true
	}

	if strings.Contains(upper, "DISCOUNT") || strings.Contains(upper, "GIẢM GIÁ") {
		return true
	}

	return false
}

// isSeparator checks if a line is a separator (===, ---)
func (p *FormatParser) isSeparator(line string) bool {
	if len(line) == 0 {
		return false
	}

	// Check if line consists only of = or - characters (but not mixed)
	firstChar := rune(line[0])
	if firstChar != '=' && firstChar != '-' {
		return false
	}

	// All characters must be the same as the first character
	for _, ch := range line {
		if ch != firstChar {
			return false
		}
	}

	return true
}

// detectFontSize determines the font size for a line based on its content
// Returns: 22pt for headers, 20pt for totals, 18pt for regular content
func (p *FormatParser) detectFontSize(line string) float64 {
	if line == "" {
		return 18.0 // Default size for empty lines
	}

	upper := strings.ToUpper(line)

	// Total line: "TỔNG CỘNG", "TOTAL", "GRAND TOTAL"
	// Check this first before header detection
	if strings.Contains(upper, "TỔNG CỘNG") || 
	   strings.Contains(upper, "GRAND TOTAL") ||
	   (strings.HasPrefix(upper, "TOTAL") && strings.Contains(line, ":")) {
		return 20.0
	}

	// Footer messages should be 18pt (not headers)
	// Check for common footer keywords
	if strings.Contains(upper, "THANK") || strings.Contains(upper, "CẢM ƠN") ||
	   strings.Contains(upper, "VISIT") || strings.Contains(upper, "WELCOME") ||
	   strings.Contains(upper, "HẸN GẶP") || strings.Contains(upper, "CHÀO") {
		return 18.0
	}

	// Header lines: Shop name, "HÓA ĐƠN BÁN HÀNG", etc.
	// These are typically short lines without colons
	if !strings.Contains(line, ":") && len(line) < 30 && len(line) > 0 {
		// Check if it's a header-like line (not a separator, not empty)
		if !p.isSeparator(line) {
			// Check for specific header keywords
			if strings.Contains(upper, "HÓA ĐƠN") || 
			   strings.Contains(upper, "INVOICE") ||
			   strings.Contains(upper, "RECEIPT") ||
			   strings.Contains(upper, "BÁN HÀNG") {
				return 22.0
			}
			// Shop name is typically the first non-empty line
			// For now, we'll use 22pt for short centered lines that look like headers
			return 22.0
		}
	}

	// All other content: table rows, order info, footer, etc.
	return 18.0
}

// isTableRow determines if a line is part of a table
// Table rows are identified by being between [TABLE_START] and [TABLE_END] markers
// or by having a specific structure (multiple columns with spacing)
func (p *FormatParser) isTableRow(line string) bool {
	if line == "" {
		return false
	}

	// Check for table markers
	if strings.Contains(line, "[TABLE_START]") || strings.Contains(line, "[TABLE_END]") {
		return false // Markers themselves are not table rows
	}

	// Table header row typically contains column names
	upper := strings.ToUpper(line)
	if strings.Contains(upper, "TÊN MÓN") || 
	   strings.Contains(upper, "SỐ LƯỢNG") || 
	   strings.Contains(upper, "ĐƠN GIÁ") ||
	   strings.Contains(upper, "THÀNH TIỀN") ||
	   strings.Contains(upper, "ITEM NAME") ||
	   strings.Contains(upper, "QUANTITY") ||
	   strings.Contains(upper, "PRICE") {
		return true
	}

	// Table separator lines (dashes between header and content)
	if p.isSeparator(line) && strings.Contains(line, "-") {
		// This could be a table separator
		// We'll mark it as a table row if it's a dash separator
		return true
	}

	// Item rows typically have multiple spaces (column separation)
	// and contain numbers (quantity, prices)
	// This is a heuristic and may need refinement
	spaceCount := strings.Count(line, "  ") // Count double spaces
	hasNumbers := false
	for _, ch := range line {
		if ch >= '0' && ch <= '9' {
			hasNumbers = true
			break
		}
	}

	// If line has multiple column-like spacing and numbers, it's likely a table row
	if spaceCount >= 2 && hasNumbers {
		return true
	}

	return false
}
