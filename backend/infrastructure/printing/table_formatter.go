package printing

import (
	"fmt"
	"strings"
)

// TableFormatter formats order items into a structured table
type TableFormatter struct {
	paperWidth int
	margin     int
	columnGap  int // Space between columns in characters
}

// TableColumn defines a column in the table
type TableColumn struct {
	Header    string
	Width     int       // Width in characters
	Alignment Alignment // Left, Center, Right
}

// TableRow defines a row in the table
type TableRow struct {
	Cells []string
}

// OrderItem represents an item in an order (for table formatting)
type OrderItem struct {
	Name        string
	VariantName string
	Quantity    int
	UnitPrice   float64
	TotalPrice  float64
}

// NewTableFormatter creates a new TableFormatter instance
func NewTableFormatter(paperWidth, margin, columnGap int) *TableFormatter {
	if paperWidth <= 0 {
		paperWidth = 576 // Default to 80mm paper width
	}
	if margin < 0 {
		margin = 0
	}
	if columnGap < 0 {
		columnGap = 1 // Default to 1 character gap
	}

	return &TableFormatter{
		paperWidth: paperWidth,
		margin:     margin,
		columnGap:  columnGap,
	}
}

// FormatItemsTable formats order items into table lines
// Returns a slice of strings representing the formatted table
func (f *TableFormatter) FormatItemsTable(items []OrderItem, paperWidth int) []string {
	if len(items) == 0 {
		return []string{}
	}

	// Update paper width if provided
	if paperWidth > 0 {
		f.paperWidth = paperWidth
	}

	// Calculate column widths
	columns := f.calculateColumnWidths(f.paperWidth)

	// Build table lines
	var lines []string

	// Add header row
	headerRow := TableRow{
		Cells: []string{
			columns[0].Header,
			columns[1].Header,
			columns[2].Header,
			columns[3].Header,
		},
	}
	lines = append(lines, f.formatRow(headerRow, columns))

	// Add separator line
	separatorWidth := 0
	for i, col := range columns {
		separatorWidth += col.Width
		if i < len(columns)-1 {
			separatorWidth += f.columnGap
		}
	}
	lines = append(lines, strings.Repeat("-", separatorWidth))

	// Add item rows
	for _, item := range items {
		// Format prices
		unitPriceStr := formatPrice(item.UnitPrice)
		totalPriceStr := formatPrice(item.TotalPrice)
		quantityStr := fmt.Sprintf("%d", item.Quantity)

		// Create main row
		itemRow := TableRow{
			Cells: []string{
				item.Name,
				quantityStr,
				unitPriceStr,
				totalPriceStr,
			},
		}

		// Format the row (this handles text wrapping for long names)
		formattedRows := f.formatRowWithWrapping(itemRow, columns)
		lines = append(lines, formattedRows...)

		// Add variant row if present
		if item.VariantName != "" {
			variantRow := TableRow{
				Cells: []string{
					"  (" + item.VariantName + ")",
					"",
					"",
					"",
				},
			}
			lines = append(lines, f.formatRow(variantRow, columns))
		}
	}

	return lines
}

// calculateColumnWidths calculates optimal column widths based on paper width
// Column distribution: Name (50%), Quantity (15%), Unit Price (17.5%), Total (17.5%)
func (f *TableFormatter) calculateColumnWidths(paperWidth int) []TableColumn {
	// Calculate available width (excluding margins)
	availableWidth := paperWidth - 2*f.margin

	// Calculate character width (approximate: 1 character ≈ 8 pixels for 18pt font)
	charWidth := 8
	availableChars := availableWidth / charWidth

	// Subtract gaps between columns (3 gaps for 4 columns)
	gapChars := f.columnGap * 3
	availableChars -= gapChars

	// Ensure minimum available characters
	if availableChars < 20 {
		availableChars = 20
	}

	// Calculate column widths based on percentages
	nameWidth := int(float64(availableChars) * 0.50)
	qtyWidth := int(float64(availableChars) * 0.15)
	unitPriceWidth := int(float64(availableChars) * 0.175)
	totalWidth := int(float64(availableChars) * 0.175)

	// Ensure minimum widths
	if nameWidth < 10 {
		nameWidth = 10
	}
	if qtyWidth < 3 {
		qtyWidth = 3
	}
	if unitPriceWidth < 6 {
		unitPriceWidth = 6
	}
	if totalWidth < 6 {
		totalWidth = 6
	}

	return []TableColumn{
		{Header: "Tên món", Width: nameWidth, Alignment: AlignLeft},
		{Header: "SL", Width: qtyWidth, Alignment: AlignRight},
		{Header: "Đơn giá", Width: unitPriceWidth, Alignment: AlignRight},
		{Header: "Thành tiền", Width: totalWidth, Alignment: AlignRight},
	}
}

// formatRow formats a single row with column alignment
func (f *TableFormatter) formatRow(row TableRow, columns []TableColumn) string {
	if len(row.Cells) != len(columns) {
		// Handle mismatch by padding with empty cells
		for len(row.Cells) < len(columns) {
			row.Cells = append(row.Cells, "")
		}
	}

	var parts []string
	for i, cell := range row.Cells {
		if i >= len(columns) {
			break
		}

		col := columns[i]
		formatted := f.alignText(cell, col.Width, col.Alignment)
		parts = append(parts, formatted)
	}

	return strings.Join(parts, strings.Repeat(" ", f.columnGap))
}

// formatRowWithWrapping formats a row and handles text wrapping for cells that are too long
// Returns multiple lines if wrapping is needed
func (f *TableFormatter) formatRowWithWrapping(row TableRow, columns []TableColumn) []string {
	if len(row.Cells) != len(columns) {
		// Handle mismatch by padding with empty cells
		for len(row.Cells) < len(columns) {
			row.Cells = append(row.Cells, "")
		}
	}

	// Wrap cells that are too long
	wrappedCells := make([][]string, len(row.Cells))
	maxLines := 1

	for i, cell := range row.Cells {
		if i >= len(columns) {
			break
		}

		col := columns[i]
		wrapped := f.wrapCellText(cell, col.Width)
		wrappedCells[i] = wrapped

		if len(wrapped) > maxLines {
			maxLines = len(wrapped)
		}
	}

	// Build output lines
	var lines []string
	for lineIdx := 0; lineIdx < maxLines; lineIdx++ {
		var parts []string
		for colIdx, col := range columns {
			var cellText string
			if lineIdx < len(wrappedCells[colIdx]) {
				cellText = wrappedCells[colIdx][lineIdx]
			} else {
				cellText = ""
			}

			formatted := f.alignText(cellText, col.Width, col.Alignment)
			parts = append(parts, formatted)
		}
		lines = append(lines, strings.Join(parts, strings.Repeat(" ", f.columnGap)))
	}

	return lines
}

// wrapCellText wraps text in a cell if it exceeds maxWidth
// Returns a slice of strings, one for each line
func (f *TableFormatter) wrapCellText(text string, maxWidth int) []string {
	if text == "" {
		return []string{""}
	}

	// Calculate text width in characters (approximate for Vietnamese text)
	textWidth := f.calculateTextWidth(text)

	// If text fits, return as-is
	if textWidth <= maxWidth {
		return []string{text}
	}

	// Split text into words
	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{text}
	}

	var lines []string
	var currentLine strings.Builder
	currentWidth := 0

	for _, word := range words {
		wordWidth := f.calculateTextWidth(word)
		spaceWidth := 1 // Space is approximately 1 character width

		// Check if adding this word would exceed maxWidth
		widthWithWord := currentWidth
		if currentLine.Len() > 0 {
			widthWithWord += spaceWidth
		}
		widthWithWord += wordWidth

		if widthWithWord > maxWidth && currentLine.Len() > 0 {
			// Current line is full, start new line
			lines = append(lines, currentLine.String())
			currentLine.Reset()
			currentLine.WriteString(word)
			currentWidth = wordWidth
		} else {
			// Add word to current line
			if currentLine.Len() > 0 {
				currentLine.WriteString(" ")
				currentWidth += spaceWidth
			}
			currentLine.WriteString(word)
			currentWidth += wordWidth
		}
	}

	// Add the last line
	if currentLine.Len() > 0 {
		lines = append(lines, currentLine.String())
	}

	// If no lines were created, return original text
	if len(lines) == 0 {
		return []string{text}
	}

	return lines
}

// alignText aligns text within a given width
func (f *TableFormatter) alignText(text string, width int, alignment Alignment) string {
	textWidth := f.calculateTextWidth(text)

	// If text is longer than width, truncate it
	if textWidth > width {
		// Truncate text to fit (approximate)
		runes := []rune(text)
		if len(runes) > width {
			text = string(runes[:width])
			textWidth = width
		}
	}

	padding := width - textWidth
	if padding < 0 {
		padding = 0
	}

	switch alignment {
	case AlignLeft:
		return text + strings.Repeat(" ", padding)
	case AlignRight:
		return strings.Repeat(" ", padding) + text
	case AlignCenter:
		leftPad := padding / 2
		rightPad := padding - leftPad
		return strings.Repeat(" ", leftPad) + text + strings.Repeat(" ", rightPad)
	default:
		return text + strings.Repeat(" ", padding)
	}
}

// calculateTextWidth calculates the approximate width of text in characters
// This is a simplified calculation that works reasonably well for Vietnamese text
func (f *TableFormatter) calculateTextWidth(text string) int {
	// For Vietnamese text with mixed ASCII and Unicode characters,
	// we approximate the width by counting runes
	return len([]rune(text))
}

// formatPrice formats a price value as a string with thousand separators
func formatPrice(price float64) string {
	// Convert to integer (assuming prices are in VND without decimals)
	priceInt := int(price)

	// Convert to string
	priceStr := fmt.Sprintf("%d", priceInt)

	// Add thousand separators
	if len(priceStr) <= 3 {
		return priceStr
	}

	// Build result with commas
	var result strings.Builder
	for i, digit := range priceStr {
		if i > 0 && (len(priceStr)-i)%3 == 0 {
			result.WriteRune(',')
		}
		result.WriteRune(digit)
	}

	return result.String()
}
