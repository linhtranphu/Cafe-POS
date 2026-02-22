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
