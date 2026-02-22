# Pixel Width Calculation Implementation

## Overview

Implemented the `CalculatePixelWidth` helper function that calculates pixel width from paper width in millimeters for ESC/POS thermal printers.

## Implementation

**Location:** `backend/infrastructure/printing/escpos_printer.go`

**Function:**
```go
func CalculatePixelWidth(paperWidthMM int) int {
	return int((float64(paperWidthMM) / 25.4) * float64(DPI))
}
```

**Formula:** `(paper_width_mm / 25.4) * 203 DPI`

This converts:
1. Millimeters to inches (divide by 25.4)
2. Inches to pixels (multiply by DPI)

## Test Coverage

**Location:** `backend/infrastructure/printing/escpos_printer_test.go`

Tests include:
- `TestCalculatePixelWidth`: Validates the formula implementation
- `TestCalculatePixelWidth_Formula`: Verifies the mathematical calculation
- `TestPixelWidthConstants`: Validates the constant definitions
- `TestCalculatePixelWidth_MatchesConstants`: Documents the relationship between formula and constants

All tests pass successfully.

## Important Note: Formula vs Constants Discrepancy

There is a discrepancy between the formula results and the design constants:

| Paper Width | Formula Result | Design Constant | Difference |
|-------------|---------------|-----------------|------------|
| 58mm        | 463 pixels    | 384 pixels      | +79 pixels |
| 80mm        | 639 pixels    | 576 pixels      | +63 pixels |

### Possible Explanations

1. **Printable Area vs Total Width**: The constants (384, 576) may represent the actual printable area, while the formula calculates the total paper width in pixels.

2. **Different DPI**: The constants may be based on a different DPI:
   - 384 pixels at 58mm = ~168 DPI
   - 576 pixels at 80mm = ~183 DPI

3. **Printer Hardware Specifications**: The constants may be based on actual printer hardware specifications rather than theoretical calculations.

### Recommendation

The `CalculatePixelWidth` function implements the exact formula specified in requirements 1.3, 7.1, and 7.2. When this function is integrated into the image rendering pipeline (tasks 7.1-7.3), the actual pixel width values should be tested with the physical printer to determine which values produce the best results.

If the constants (384, 576) are found to be more accurate for the actual printer hardware, the function can be updated to use a lookup table or adjusted formula.

## Requirements Validated

- ✅ Requirement 1.3: Calculate pixel width based on paper width
- ✅ Requirement 7.1: Support 58mm paper width
- ✅ Requirement 7.2: Support 80mm paper width

## Next Steps

This helper function will be used by:
- Task 7.1: ESCPOSPrinter initialization to calculate pixel width from config
- Task 3.1: TextRenderer to determine image width
- Task 5.1: ImageConverter to validate image dimensions
