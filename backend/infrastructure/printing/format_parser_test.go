package printing

import (
	"strings"
	"testing"
)

// Unit Tests

func TestNewFormatParser(t *testing.T) {
	parser := NewFormatParser(384)
	if parser == nil {
		t.Fatal("NewFormatParser returned nil")
	}
	if parser.paperWidth != 384 {
		t.Errorf("Expected paperWidth 384, got %d", parser.paperWidth)
	}
}

func TestParse_EmptyContent(t *testing.T) {
	parser := NewFormatParser(384)
	result := parser.Parse("")
	
	if len(result) != 1 {
		t.Errorf("Expected 1 line for empty content, got %d", len(result))
	}
	
	if result[0].Text != "" {
		t.Errorf("Expected empty text, got %q", result[0].Text)
	}
}

func TestParse_SingleLine(t *testing.T) {
	parser := NewFormatParser(384)
	result := parser.Parse("Hello World")
	
	if len(result) != 1 {
		t.Errorf("Expected 1 line, got %d", len(result))
	}
	
	if result[0].Text != "Hello World" {
		t.Errorf("Expected 'Hello World', got %q", result[0].Text)
	}
}

func TestParse_MultipleLines(t *testing.T) {
	parser := NewFormatParser(384)
	content := "Line 1\nLine 2\nLine 3"
	result := parser.Parse(content)
	
	if len(result) != 3 {
		t.Errorf("Expected 3 lines, got %d", len(result))
	}
	
	expected := []string{"Line 1", "Line 2", "Line 3"}
	for i, line := range result {
		if line.Text != expected[i] {
			t.Errorf("Line %d: expected %q, got %q", i, expected[i], line.Text)
		}
	}
}

func TestDetectAlignment_Separator(t *testing.T) {
	parser := NewFormatParser(384)
	
	tests := []struct {
		line     string
		expected Alignment
	}{
		{"===", AlignCenter},
		{"---", AlignCenter},
		{"==========", AlignCenter},
		{"----------", AlignCenter},
	}
	
	for _, tt := range tests {
		result := parser.detectAlignment(tt.line)
		if result != tt.expected {
			t.Errorf("detectAlignment(%q) = %v, want %v", tt.line, result, tt.expected)
		}
	}
}

func TestDetectAlignment_Header(t *testing.T) {
	parser := NewFormatParser(384)
	
	tests := []struct {
		line     string
		expected Alignment
	}{
		{"MY SHOP", AlignCenter},
		{"Cửa hàng của tôi", AlignCenter},
		{"123 Main St", AlignCenter},
	}
	
	for _, tt := range tests {
		result := parser.detectAlignment(tt.line)
		if result != tt.expected {
			t.Errorf("detectAlignment(%q) = %v, want %v", tt.line, result, tt.expected)
		}
	}
}

func TestDetectAlignment_Total(t *testing.T) {
	parser := NewFormatParser(384)
	
	tests := []struct {
		line     string
		expected Alignment
	}{
		{"TOTAL: 100,000đ", AlignCenter},
		{"Total: $50.00", AlignCenter},
		{"TỔNG: 100,000đ", AlignCenter},
	}
	
	for _, tt := range tests {
		result := parser.detectAlignment(tt.line)
		if result != tt.expected {
			t.Errorf("detectAlignment(%q) = %v, want %v", tt.line, result, tt.expected)
		}
	}
}

func TestDetectAlignment_ThankYou(t *testing.T) {
	parser := NewFormatParser(384)
	
	tests := []struct {
		line     string
		expected Alignment
	}{
		{"Thank you!", AlignCenter},
		{"CẢM ƠN QUÝ KHÁCH", AlignCenter},
		{"Thank you for your visit", AlignCenter},
	}
	
	for _, tt := range tests {
		result := parser.detectAlignment(tt.line)
		if result != tt.expected {
			t.Errorf("detectAlignment(%q) = %v, want %v", tt.line, result, tt.expected)
		}
	}
}

func TestDetectAlignment_ItemLine(t *testing.T) {
	parser := NewFormatParser(384)
	
	tests := []struct {
		line     string
		expected Alignment
	}{
		{"Coffee: 50,000đ", AlignLeft},
		{"Cà phê sữa đá: 25,000đ", AlignLeft},
		{"Item with a longer description: 100,000đ", AlignLeft},
	}
	
	for _, tt := range tests {
		result := parser.detectAlignment(tt.line)
		if result != tt.expected {
			t.Errorf("detectAlignment(%q) = %v, want %v", tt.line, result, tt.expected)
		}
	}
}

func TestDetectBold_Header(t *testing.T) {
	parser := NewFormatParser(384)
	
	tests := []struct {
		line     string
		expected bool
	}{
		{"MY SHOP", true},
		{"Cửa hàng", true},
		{"123 Main St", true},
	}
	
	for _, tt := range tests {
		result := parser.detectBold(tt.line)
		if result != tt.expected {
			t.Errorf("detectBold(%q) = %v, want %v", tt.line, result, tt.expected)
		}
	}
}

func TestDetectBold_Total(t *testing.T) {
	parser := NewFormatParser(384)
	
	tests := []struct {
		line     string
		expected bool
	}{
		{"TOTAL: 100,000đ", true},
		{"Total: $50.00", true},
		{"TỔNG: 100,000đ", true},
		{"SUBTOTAL: 90,000đ", true},
		{"TỔNG PHỤ: 90,000đ", true},
	}
	
	for _, tt := range tests {
		result := parser.detectBold(tt.line)
		if result != tt.expected {
			t.Errorf("detectBold(%q) = %v, want %v", tt.line, result, tt.expected)
		}
	}
}

func TestDetectBold_Discount(t *testing.T) {
	parser := NewFormatParser(384)
	
	tests := []struct {
		line     string
		expected bool
	}{
		{"DISCOUNT: 10,000đ", true},
		{"Discount: $5.00", true},
		{"GIẢM GIÁ: 10,000đ", true},
	}
	
	for _, tt := range tests {
		result := parser.detectBold(tt.line)
		if result != tt.expected {
			t.Errorf("detectBold(%q) = %v, want %v", tt.line, result, tt.expected)
		}
	}
}

func TestDetectBold_ItemLine(t *testing.T) {
	parser := NewFormatParser(384)
	
	tests := []struct {
		line     string
		expected bool
	}{
		{"Coffee: 50,000đ", false},
		{"Cà phê sữa đá: 25,000đ", false},
		{"Item with a longer description: 100,000đ", false},
	}
	
	for _, tt := range tests {
		result := parser.detectBold(tt.line)
		if result != tt.expected {
			t.Errorf("detectBold(%q) = %v, want %v", tt.line, result, tt.expected)
		}
	}
}

func TestIsSeparator(t *testing.T) {
	parser := NewFormatParser(384)
	
	tests := []struct {
		line     string
		expected bool
	}{
		{"===", true},
		{"---", true},
		{"==========", true},
		{"----------", true},
		{"=-=", false},
		{"= = =", false},
		{"", false},
		{"Hello", false},
		{"123", false},
	}
	
	for _, tt := range tests {
		result := parser.isSeparator(tt.line)
		if result != tt.expected {
			t.Errorf("isSeparator(%q) = %v, want %v", tt.line, result, tt.expected)
		}
	}
}

func TestParse_VietnameseReceipt(t *testing.T) {
	parser := NewFormatParser(384)
	content := `CỬA HÀNG CÀ PHÊ
123 Đường Nguyễn Huệ
===
Cà phê sữa đá: 25,000đ
Bánh mì: 15,000đ
---
TỔNG: 40,000đ
===
CẢM ƠN QUÝ KHÁCH`

	result := parser.Parse(content)
	
	// Verify line count
	expectedLines := 9
	if len(result) != expectedLines {
		t.Errorf("Expected %d lines, got %d", expectedLines, len(result))
	}
	
	// Verify header formatting (line 0)
	if !result[0].Bold {
		t.Error("Header line should be bold")
	}
	if result[0].Alignment != AlignCenter {
		t.Error("Header line should be centered")
	}
	
	// Verify separator (line 2)
	if !result[2].IsSeparator {
		t.Error("Line 2 should be a separator")
	}
	if result[2].Alignment != AlignCenter {
		t.Error("Separator should be centered")
	}
	
	// Verify item lines (lines 3-4)
	if result[3].Alignment != AlignLeft {
		t.Error("Item line should be left-aligned")
	}
	if result[3].Bold {
		t.Error("Item line should not be bold")
	}
	
	// Verify total line (line 6)
	if !result[6].Bold {
		t.Error("Total line should be bold")
	}
	if result[6].Alignment != AlignCenter {
		t.Error("Total line should be centered")
	}
	
	// Verify footer (line 8)
	if result[8].Alignment != AlignCenter {
		t.Error("Footer line should be centered")
	}
}

func TestParse_PreservesOriginalText(t *testing.T) {
	parser := NewFormatParser(384)
	content := "  Indented text  \nNormal text"
	result := parser.Parse(content)
	
	// Original text should be preserved, including whitespace
	if result[0].Text != "  Indented text  " {
		t.Errorf("Expected '  Indented text  ', got %q", result[0].Text)
	}
	if result[1].Text != "Normal text" {
		t.Errorf("Expected 'Normal text', got %q", result[1].Text)
	}
}

func TestParse_EmptyLines(t *testing.T) {
	parser := NewFormatParser(384)
	content := "Line 1\n\nLine 3"
	result := parser.Parse(content)
	
	if len(result) != 3 {
		t.Errorf("Expected 3 lines, got %d", len(result))
	}
	
	if result[1].Text != "" {
		t.Errorf("Expected empty line, got %q", result[1].Text)
	}
	
	if result[1].Bold {
		t.Error("Empty line should not be bold")
	}
	
	if result[1].Alignment != AlignLeft {
		t.Error("Empty line should be left-aligned")
	}
}

// Property-Based Tests
// These tests validate the correctness properties defined in the design document

// Property 13: Receipt Header Formatting
// For any receipt template with header sections, the FormatParser should identify them
// and apply centered, bold formatting.
// **Validates: Requirements 5.1**
func TestProperty_ReceiptHeaderFormatting(t *testing.T) {
	parser := NewFormatParser(384)
	
	// Test various header patterns
	headers := []string{
		"MY SHOP",
		"CỬA HÀNG",
		"Café ABC",
		"Restaurant XYZ",
		"123 Main Street",
		"Hà Nội, Việt Nam",
	}
	
	for _, header := range headers {
		result := parser.Parse(header)
		
		if len(result) != 1 {
			t.Errorf("Expected 1 line for header %q, got %d", header, len(result))
			continue
		}
		
		line := result[0]
		
		// Headers should be centered
		if line.Alignment != AlignCenter {
			t.Errorf("Header %q should be centered, got alignment %v", header, line.Alignment)
		}
		
		// Headers should be bold
		if !line.Bold {
			t.Errorf("Header %q should be bold", header)
		}
		
		// Headers should not be separators
		if line.IsSeparator {
			t.Errorf("Header %q should not be a separator", header)
		}
	}
}

// Property 14: Receipt Item Formatting
// For any receipt template with item lines, the FormatParser should identify them
// and apply left-aligned formatting.
// **Validates: Requirements 5.2**
func TestProperty_ReceiptItemFormatting(t *testing.T) {
	parser := NewFormatParser(384)
	
	// Test various item line patterns
	items := []string{
		"Coffee: 50,000đ",
		"Cà phê sữa đá: 25,000đ",
		"Bánh mì thịt: 15,000đ",
		"Trà sữa trân châu: 30,000đ",
		"Item with long description: 100,000đ",
	}
	
	for _, item := range items {
		result := parser.Parse(item)
		
		if len(result) != 1 {
			t.Errorf("Expected 1 line for item %q, got %d", item, len(result))
			continue
		}
		
		line := result[0]
		
		// Item lines should be left-aligned
		if line.Alignment != AlignLeft {
			t.Errorf("Item %q should be left-aligned, got alignment %v", item, line.Alignment)
		}
		
		// Item lines should not be bold
		if line.Bold {
			t.Errorf("Item %q should not be bold", item)
		}
		
		// Item lines should not be separators
		if line.IsSeparator {
			t.Errorf("Item %q should not be a separator", item)
		}
	}
}

// Property 15: Receipt Total Formatting
// For any receipt template with total lines, the FormatParser should identify them
// and apply bold formatting.
// **Validates: Requirements 5.3**
func TestProperty_ReceiptTotalFormatting(t *testing.T) {
	parser := NewFormatParser(384)
	
	// Test various total line patterns
	totals := []string{
		"TOTAL: 100,000đ",
		"Total: $50.00",
		"TỔNG: 100,000đ",
		"SUBTOTAL: 90,000đ",
		"TỔNG PHỤ: 90,000đ",
		"DISCOUNT: 10,000đ",
		"GIẢM GIÁ: 10,000đ",
	}
	
	for _, total := range totals {
		result := parser.Parse(total)
		
		if len(result) != 1 {
			t.Errorf("Expected 1 line for total %q, got %d", total, len(result))
			continue
		}
		
		line := result[0]
		
		// Total lines should be bold
		if !line.Bold {
			t.Errorf("Total %q should be bold", total)
		}
		
		// Total lines should not be separators
		if line.IsSeparator {
			t.Errorf("Total %q should not be a separator", total)
		}
	}
}

// Property 16: Receipt Footer Formatting
// For any receipt template with footer sections, the FormatParser should identify them
// and apply centered formatting.
// **Validates: Requirements 5.4**
func TestProperty_ReceiptFooterFormatting(t *testing.T) {
	parser := NewFormatParser(384)
	
	// Test various footer patterns
	footers := []string{
		"Thank you!",
		"CẢM ƠN QUÝ KHÁCH",
		"Thank you for your visit",
		"HẸN GẶP LẠI",
		"WELCOME BACK",
		"CHÀO MỪNG",
	}
	
	for _, footer := range footers {
		result := parser.Parse(footer)
		
		if len(result) != 1 {
			t.Errorf("Expected 1 line for footer %q, got %d", footer, len(result))
			continue
		}
		
		line := result[0]
		
		// Footer lines should be centered
		if line.Alignment != AlignCenter {
			t.Errorf("Footer %q should be centered, got alignment %v", footer, line.Alignment)
		}
		
		// Footer lines should not be separators
		if line.IsSeparator {
			t.Errorf("Footer %q should not be a separator", footer)
		}
	}
}

// Test that separator lines are properly identified and formatted
func TestProperty_SeparatorFormatting(t *testing.T) {
	parser := NewFormatParser(384)
	
	separators := []string{
		"===",
		"---",
		"==========",
		"----------",
		strings.Repeat("=", 50),
		strings.Repeat("-", 50),
	}
	
	for _, sep := range separators {
		result := parser.Parse(sep)
		
		if len(result) != 1 {
			t.Errorf("Expected 1 line for separator %q, got %d", sep, len(result))
			continue
		}
		
		line := result[0]
		
		// Separators should be identified
		if !line.IsSeparator {
			t.Errorf("Line %q should be identified as separator", sep)
		}
		
		// Separators should be centered
		if line.Alignment != AlignCenter {
			t.Errorf("Separator %q should be centered, got alignment %v", sep, line.Alignment)
		}
	}
}

// Test that the parser handles various Vietnamese characters correctly
func TestProperty_VietnameseCharacterSupport(t *testing.T) {
	parser := NewFormatParser(384)
	
	// Test various Vietnamese text with different tones and diacritics
	vietnameseTexts := []string{
		"Cà phê sữa đá",
		"Bánh mì thịt",
		"Trà sữa trân châu",
		"Phở bò",
		"Bún chả",
		"Gỏi cuốn",
		"Cơm tấm",
		"Hủ tiếu",
		"TỔNG CỘNG",
		"GIẢM GIÁ",
		"CẢM ƠN QUÝ KHÁCH",
	}
	
	for _, text := range vietnameseTexts {
		result := parser.Parse(text)
		
		if len(result) != 1 {
			t.Errorf("Expected 1 line for text %q, got %d", text, len(result))
			continue
		}
		
		// Text should be preserved exactly
		if result[0].Text != text {
			t.Errorf("Text not preserved: expected %q, got %q", text, result[0].Text)
		}
	}
}

// Test that multi-line content is parsed correctly
func TestProperty_MultiLinePreservation(t *testing.T) {
	parser := NewFormatParser(384)
	
	content := `Line 1
Line 2
Line 3
Line 4
Line 5`
	
	result := parser.Parse(content)
	
	expectedLines := 5
	if len(result) != expectedLines {
		t.Errorf("Expected %d lines, got %d", expectedLines, len(result))
	}
	
	lines := strings.Split(content, "\n")
	for i, expected := range lines {
		if i >= len(result) {
			break
		}
		if result[i].Text != expected {
			t.Errorf("Line %d: expected %q, got %q", i, expected, result[i].Text)
		}
	}
}

// Tests for new FontSize and IsTableRow functionality

func TestDetectFontSize_Headers(t *testing.T) {
	parser := NewFormatParser(384)
	
	tests := []struct {
		line     string
		expected float64
	}{
		{"HÓA ĐƠN BÁN HÀNG", 22.0},
		{"INVOICE", 22.0},
		{"RECEIPT", 22.0},
		{"MY SHOP", 22.0},
		{"Cửa hàng ABC", 22.0},
	}
	
	for _, tt := range tests {
		result := parser.detectFontSize(tt.line)
		if result != tt.expected {
			t.Errorf("detectFontSize(%q) = %v, want %v", tt.line, result, tt.expected)
		}
	}
}

func TestDetectFontSize_Totals(t *testing.T) {
	parser := NewFormatParser(384)
	
	tests := []struct {
		line     string
		expected float64
	}{
		{"TỔNG CỘNG: 100,000đ", 20.0},
		{"GRAND TOTAL: $50.00", 20.0},
		{"TOTAL: 100,000đ", 20.0},
	}
	
	for _, tt := range tests {
		result := parser.detectFontSize(tt.line)
		if result != tt.expected {
			t.Errorf("detectFontSize(%q) = %v, want %v", tt.line, result, tt.expected)
		}
	}
}

func TestDetectFontSize_RegularContent(t *testing.T) {
	parser := NewFormatParser(384)
	
	tests := []struct {
		line     string
		expected float64
	}{
		{"Coffee: 50,000đ", 18.0},
		{"Order: #12345", 18.0},
		{"Ngày: 01/01/2024", 18.0},
		{"Bàn: 5", 18.0},
		{"Cảm ơn quý khách!", 18.0},
		{"", 18.0},
	}
	
	for _, tt := range tests {
		result := parser.detectFontSize(tt.line)
		if result != tt.expected {
			t.Errorf("detectFontSize(%q) = %v, want %v", tt.line, result, tt.expected)
		}
	}
}

func TestDetectFontSize_Separators(t *testing.T) {
	parser := NewFormatParser(384)
	
	tests := []struct {
		line     string
		expected float64
	}{
		{"===", 18.0},
		{"---", 18.0},
		{"==========", 18.0},
	}
	
	for _, tt := range tests {
		result := parser.detectFontSize(tt.line)
		if result != tt.expected {
			t.Errorf("detectFontSize(%q) = %v, want %v", tt.line, result, tt.expected)
		}
	}
}

func TestIsTableRow_TableHeaders(t *testing.T) {
	parser := NewFormatParser(384)
	
	tests := []struct {
		line     string
		expected bool
	}{
		{"Tên món              SL  Đơn giá    Thành tiền", true},
		{"TÊN MÓN              SL  ĐƠN GIÁ    THÀNH TIỀN", true},
		{"ITEM NAME            QTY PRICE      TOTAL", true},
	}
	
	for _, tt := range tests {
		result := parser.isTableRow(tt.line)
		if result != tt.expected {
			t.Errorf("isTableRow(%q) = %v, want %v", tt.line, result, tt.expected)
		}
	}
}

func TestIsTableRow_TableSeparators(t *testing.T) {
	parser := NewFormatParser(384)
	
	tests := []struct {
		line     string
		expected bool
	}{
		{"------------------------------------------------", true},
		{"---", true},
		{"----------", true},
	}
	
	for _, tt := range tests {
		result := parser.isTableRow(tt.line)
		if result != tt.expected {
			t.Errorf("isTableRow(%q) = %v, want %v", tt.line, result, tt.expected)
		}
	}
}

func TestIsTableRow_ItemRows(t *testing.T) {
	parser := NewFormatParser(384)
	
	tests := []struct {
		line     string
		expected bool
	}{
		{"Cafe Latte            2   45,000       90,000", true},
		{"Banh Mi Thit          1   35,000       35,000", true},
		{"Coffee  2  50000  100000", true},
	}
	
	for _, tt := range tests {
		result := parser.isTableRow(tt.line)
		if result != tt.expected {
			t.Errorf("isTableRow(%q) = %v, want %v", tt.line, result, tt.expected)
		}
	}
}

func TestIsTableRow_NonTableLines(t *testing.T) {
	parser := NewFormatParser(384)
	
	tests := []struct {
		line     string
		expected bool
	}{
		{"HÓA ĐƠN BÁN HÀNG", false},
		{"Order: #12345", false},
		{"Ngày: 01/01/2024", false},
		{"Cảm ơn quý khách!", false},
		{"", false},
		{"[TABLE_START]", false},
		{"[TABLE_END]", false},
		{"===", false},
	}
	
	for _, tt := range tests {
		result := parser.isTableRow(tt.line)
		if result != tt.expected {
			t.Errorf("isTableRow(%q) = %v, want %v", tt.line, result, tt.expected)
		}
	}
}

func TestParse_WithFontSizeAndTableRow(t *testing.T) {
	parser := NewFormatParser(384)
	content := `CỬA HÀNG CÀ PHÊ
HÓA ĐƠN BÁN HÀNG
===
Order: #12345
Ngày: 01/01/2024
===
Tên món              SL  Đơn giá    Thành tiền
------------------------------------------------
Cafe Latte            2   45,000       90,000
Banh Mi Thit          1   35,000       35,000
===
TỔNG CỘNG: 125,000đ
===
Cảm ơn quý khách!`

	result := parser.Parse(content)
	
	// Verify line count
	expectedLines := 14
	if len(result) != expectedLines {
		t.Errorf("Expected %d lines, got %d", expectedLines, len(result))
	}
	
	// Verify header font size (line 0 and 1)
	if result[0].FontSize != 22.0 {
		t.Errorf("Header line 0 should have font size 22.0, got %v", result[0].FontSize)
	}
	if result[1].FontSize != 22.0 {
		t.Errorf("Header line 1 should have font size 22.0, got %v", result[1].FontSize)
	}
	
	// Verify regular content font size (line 3, 4)
	if result[3].FontSize != 18.0 {
		t.Errorf("Regular line 3 should have font size 18.0, got %v", result[3].FontSize)
	}
	if result[4].FontSize != 18.0 {
		t.Errorf("Regular line 4 should have font size 18.0, got %v", result[4].FontSize)
	}
	
	// Verify table header (line 6)
	if !result[6].IsTableRow {
		t.Error("Line 6 should be identified as table row")
	}
	if result[6].FontSize != 18.0 {
		t.Errorf("Table header should have font size 18.0, got %v", result[6].FontSize)
	}
	
	// Verify table separator (line 7)
	if !result[7].IsTableRow {
		t.Error("Line 7 should be identified as table row")
	}
	
	// Verify table item rows (line 8, 9)
	if !result[8].IsTableRow {
		t.Error("Line 8 should be identified as table row")
	}
	if !result[9].IsTableRow {
		t.Error("Line 9 should be identified as table row")
	}
	
	// Verify total line (line 11)
	if result[11].FontSize != 20.0 {
		t.Errorf("Total line should have font size 20.0, got %v", result[11].FontSize)
	}
	
	// Verify footer (line 13)
	if result[13].FontSize != 18.0 {
		t.Errorf("Footer should have font size 18.0, got %v", result[13].FontSize)
	}
}

// Property test: Font size consistency
// For any line, font size should be one of the three defined sizes: 18pt, 20pt, or 22pt
func TestProperty_FontSizeConsistency(t *testing.T) {
	parser := NewFormatParser(384)
	
	testLines := []string{
		"HÓA ĐƠN BÁN HÀNG",
		"MY SHOP",
		"TỔNG CỘNG: 100,000đ",
		"GRAND TOTAL: $50.00",
		"Coffee: 50,000đ",
		"Order: #12345",
		"Cảm ơn quý khách!",
		"===",
		"---",
		"",
		"Cafe Latte            2   45,000       90,000",
		"Tên món              SL  Đơn giá    Thành tiền",
	}
	
	validSizes := map[float64]bool{
		18.0: true,
		20.0: true,
		22.0: true,
	}
	
	for _, line := range testLines {
		result := parser.Parse(line)
		if len(result) != 1 {
			t.Errorf("Expected 1 line for %q, got %d", line, len(result))
			continue
		}
		
		fontSize := result[0].FontSize
		if !validSizes[fontSize] {
			t.Errorf("Line %q has invalid font size %v, expected one of 18.0, 20.0, or 22.0", line, fontSize)
		}
	}
}

// Property test: Table row detection consistency
// For any line identified as a table row, it should have specific characteristics
func TestProperty_TableRowDetectionConsistency(t *testing.T) {
	parser := NewFormatParser(384)
	
	// Lines that should be table rows
	tableRows := []string{
		"Tên món              SL  Đơn giá    Thành tiền",
		"Cafe Latte            2   45,000       90,000",
		"------------------------------------------------",
		"Coffee  2  50000  100000",
	}
	
	for _, line := range tableRows {
		result := parser.Parse(line)
		if len(result) != 1 {
			t.Errorf("Expected 1 line for %q, got %d", line, len(result))
			continue
		}
		
		if !result[0].IsTableRow {
			t.Errorf("Line %q should be identified as table row", line)
		}
	}
	
	// Lines that should NOT be table rows
	nonTableRows := []string{
		"HÓA ĐƠN BÁN HÀNG",
		"Order: #12345",
		"Cảm ơn quý khách!",
		"[TABLE_START]",
		"[TABLE_END]",
		"===",
	}
	
	for _, line := range nonTableRows {
		result := parser.Parse(line)
		if len(result) != 1 {
			t.Errorf("Expected 1 line for %q, got %d", line, len(result))
			continue
		}
		
		if result[0].IsTableRow {
			t.Errorf("Line %q should NOT be identified as table row", line)
		}
	}
}
