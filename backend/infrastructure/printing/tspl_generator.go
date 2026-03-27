package printing

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type TSPLGenerator struct {
	Width  int
	Height int
	DPI    int
}

type LabelData struct {
	OrderNumber  string
	ItemName     string
	Note         string
	Time         string
	CustomerName string
	Quantity     int
}

func NewTSPLGenerator(width, height, dpi int) *TSPLGenerator {
	return &TSPLGenerator{
		Width:  width,
		Height: height,
		DPI:    dpi,
	}
}

func (g *TSPLGenerator) GenerateLabelCommands(data LabelData, templateContent string) (string, error) {
	funcMap := template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
		"mul": func(a, b int) int { return a * b },
		"div": func(a, b int) int { return a / b },
	}
	tmpl, err := template.New("label").Funcs(funcMap).Parse(templateContent)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	customerName := removeDiacritics(truncateText(data.CustomerName, 30))
	customerNameLines := splitLines(customerName, 15)

	templateData := map[string]interface{}{
		"Width":             g.Width,
		"Height":            g.Height,
		"OrderNumber":       shortOrderNumber(data.OrderNumber),
		"ItemName":          removeDiacritics(truncateText(data.ItemName, 20)),
		"Note":              removeDiacritics(truncateText(data.Note, 25)),
		"Time":              data.Time,
		"CustomerName":      customerName,
		"CustomerNameLines": customerNameLines,
		"Quantity":          data.Quantity,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, templateData); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// removeDiacritics converts Vietnamese/accented text to ASCII.
func removeDiacritics(s string) string {
	s = strings.NewReplacer("đ", "d", "Đ", "D").Replace(s)
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.M)))
	result, _, err := transform.String(t, s)
	if err != nil {
		return s
	}
	return result
}

func truncateText(text string, maxLen int) string {
	if text == "" {
		return ""
	}
	r := []rune(text)
	if len(r) <= maxLen {
		return text
	}
	return string(r[:maxLen-3]) + "..."
}

// splitLines splits text into chunks of maxLen runes each.
func splitLines(text string, maxLen int) []string {
	if text == "" {
		return []string{""}
	}
	r := []rune(text)
	var lines []string
	for len(r) > 0 {
		end := maxLen
		if end > len(r) {
			end = len(r)
		}
		lines = append(lines, string(r[:end]))
		r = r[end:]
	}
	return lines
}

func ValidateTSPLCommands(commands string) error {
	if commands == "" {
		return fmt.Errorf("TSPL commands cannot be empty")
	}
	required := []string{"SIZE", "CLS", "PRINT"}
	for _, cmd := range required {
		if !bytes.Contains([]byte(commands), []byte(cmd)) {
			return fmt.Errorf("TSPL commands missing required command: %s", cmd)
		}
	}
	if !utf8.ValidString(commands) {
		return fmt.Errorf("TSPL commands contain invalid UTF-8")
	}
	return nil
}

// shortOrderNumber returns the part after the last '-' in the order number.
func shortOrderNumber(orderNumber string) string {
	if idx := strings.LastIndex(orderNumber, "-"); idx >= 0 {
		return orderNumber[idx+1:]
	}
	return orderNumber
}
