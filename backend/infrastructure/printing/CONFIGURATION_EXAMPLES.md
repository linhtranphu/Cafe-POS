# Printer Configuration Examples

## Overview

This document provides configuration examples for the Vietnamese image-based printing system.

## Basic Printer Configuration

### Network Printer (58mm paper)

```go
config := &printing.PrinterConfig{
    Name:           "Kitchen Printer",
    Type:           printing.PrinterTypeBill,
    ConnectionType: printing.ConnectionTypeNetwork,
    IPAddress:      "192.168.1.100",
    Port:           9100,
    PaperWidth:     58, // 58mm thermal paper
    IsDefault:      false,
    IsEnabled:      true,
}
```

### Network Printer (80mm paper)

```go
config := &printing.PrinterConfig{
    Name:           "Counter Printer",
    Type:           printing.PrinterTypeBill,
    ConnectionType: printing.ConnectionTypeNetwork,
    IPAddress:      "192.168.1.101",
    Port:           9100,
    PaperWidth:     80, // 80mm thermal paper
    IsDefault:      true,
    IsEnabled:      true,
}
```

## Font Configuration

### Using System Fonts (Recommended)

The system automatically discovers and uses system fonts. No additional configuration is required.

```go
// Create printer - fonts are automatically discovered
printer, err := infraPrinting.NewESCPOSPrinter(config)
if err != nil {
    log.Fatalf("Failed to create printer: %v", err)
}
```

**Supported System Fonts:**
- Windows: Arial Unicode MS, Arial, DejaVu Sans
- macOS: Arial Unicode MS, Arial
- Linux: DejaVu Sans, Liberation Sans, Noto Sans, Roboto

### Using Custom Font Path

If you need to use a specific font file, you can modify the `NewESCPOSPrinter` function to accept a font path parameter, or use environment variables.

**Option 1: Environment Variable**

```bash
# Set font path via environment variable
export CAFE_POS_FONT_PATH="/opt/cafe-pos/fonts/Roboto-Regular.ttf"
```

```go
// Read font path from environment
fontPath := os.Getenv("CAFE_POS_FONT_PATH")

// Create renderer config with custom font
rendererConfig := &RendererConfig{
    PixelWidth:  pixelWidth,
    FontPath:    fontPath, // Use custom font if set, empty for system fonts
    FontSize:    14.0,
    LineSpacing: 4,
    Margin:      8,
}
```

**Option 2: Configuration File**

Create a configuration file (e.g., `printer_config.yaml`):

```yaml
printers:
  - name: "Kitchen Printer"
    type: "BILL"
    connection_type: "NETWORK"
    ip_address: "192.168.1.100"
    port: 9100
    paper_width: 58
    font_path: "/opt/cafe-pos/fonts/Roboto-Regular.ttf"  # Optional
    font_size: 14.0
    line_spacing: 4
    margin: 8
```

## Rendering Configuration

### Default Settings

```go
rendererConfig := &RendererConfig{
    PixelWidth:  463,  // For 58mm paper (calculated automatically)
    FontPath:    "",   // Empty = use system fonts
    FontSize:    14.0, // 14pt font
    LineSpacing: 4,    // 4 pixels between lines
    Margin:      8,    // 8 pixels left/right margin
}
```

### Compact Settings (Save Paper)

```go
rendererConfig := &RendererConfig{
    PixelWidth:  463,
    FontPath:    "",
    FontSize:    12.0, // Smaller font
    LineSpacing: 2,    // Tighter line spacing
    Margin:      4,    // Smaller margins
}
```

### Large Text Settings (Better Readability)

```go
rendererConfig := &RendererConfig{
    PixelWidth:  639,  // For 80mm paper
    FontPath:    "",
    FontSize:    16.0, // Larger font
    Line_spacing: 6,   // More spacing
    Margin:      12,   // Larger margins
}
```

## Complete Example

### Production Setup with Error Handling

```go
package main

import (
    "log"
    "os"
    
    "cafe-pos/backend/domain/printing"
    infraPrinting "cafe-pos/backend/infrastructure/printing"
)

func main() {
    // Create printer configuration
    config := &printing.PrinterConfig{
        Name:           "Main Printer",
        Type:           printing.PrinterTypeBill,
        ConnectionType: printing.ConnectionTypeNetwork,
        IPAddress:      os.Getenv("PRINTER_IP"),
        Port:           9100,
        PaperWidth:     80,
        IsDefault:      true,
        IsEnabled:      true,
    }
    
    // Validate configuration
    if config.IPAddress == "" {
        log.Fatal("PRINTER_IP environment variable not set")
    }
    
    // Create printer instance
    printer, err := infraPrinting.NewESCPOSPrinter(config)
    if err != nil {
        log.Fatalf("Failed to create printer: %v", err)
    }
    
    // Connect to printer
    if err := printer.Connect(); err != nil {
        log.Fatalf("Failed to connect to printer: %v", err)
    }
    defer printer.Disconnect()
    
    // Check printer status
    status, err := printer.GetStatus()
    if err != nil {
        log.Fatalf("Failed to get printer status: %v", err)
    }
    
    if !status.IsOnline {
        log.Fatalf("Printer is offline: %s", status.ErrorMsg)
    }
    
    // Print receipt
    receipt := `Cafe ABC
123 Nguyen Hue, Q1, TP.HCM
================================
HOA DON BAN HANG
Order: ORD-001
================================
Cafe Latte
  2 x 45,000 = 90,000
================================
TONG CONG: 90,000 VND
================================
Cam on quy khach!`
    
    if err := printer.Print(receipt); err != nil {
        log.Fatalf("Failed to print: %v", err)
    }
    
    log.Println("Receipt printed successfully")
}
```

### Docker Deployment

**Dockerfile:**

```dockerfile
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Copy source code
WORKDIR /app
COPY . .

# Build application
RUN go build -o cafe-pos ./cmd/server

# Runtime image
FROM alpine:latest

# Install fonts
RUN apk add --no-cache fontconfig ttf-dejavu

# Copy application
COPY --from=builder /app/cafe-pos /usr/local/bin/

# Set environment variables
ENV PRINTER_IP=192.168.1.100
ENV PRINTER_PORT=9100
ENV PAPER_WIDTH=80

# Run application
CMD ["cafe-pos"]
```

**docker-compose.yml:**

```yaml
version: '3.8'

services:
  cafe-pos:
    build: .
    environment:
      - PRINTER_IP=192.168.1.100
      - PRINTER_PORT=9100
      - PAPER_WIDTH=80
      - CAFE_POS_FONT_PATH=/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf
    volumes:
      - ./fonts:/usr/share/fonts/truetype/cafe-pos  # Optional: custom fonts
    networks:
      - cafe-network

networks:
  cafe-network:
    driver: bridge
```

### Kubernetes Deployment

**ConfigMap:**

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: printer-config
data:
  PRINTER_IP: "192.168.1.100"
  PRINTER_PORT: "9100"
  PAPER_WIDTH: "80"
  CAFE_POS_FONT_PATH: "/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf"
```

**Deployment:**

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: cafe-pos
spec:
  replicas: 1
  selector:
    matchLabels:
      app: cafe-pos
  template:
    metadata:
      labels:
        app: cafe-pos
    spec:
      containers:
      - name: cafe-pos
        image: cafe-pos:latest
        envFrom:
        - configMapRef:
            name: printer-config
        volumeMounts:
        - name: fonts
          mountPath: /usr/share/fonts/truetype/cafe-pos
      volumes:
      - name: fonts
        configMap:
          name: custom-fonts  # Optional: custom fonts
```

## Environment Variables

### Supported Variables

| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `PRINTER_IP` | Printer IP address | - | `192.168.1.100` |
| `PRINTER_PORT` | Printer port | `9100` | `9100` |
| `PAPER_WIDTH` | Paper width in mm | `80` | `58` or `80` |
| `CAFE_POS_FONT_PATH` | Custom font file path | System fonts | `/opt/fonts/Roboto.ttf` |
| `FONT_SIZE` | Font size in points | `14.0` | `12.0` to `18.0` |
| `LINE_SPACING` | Line spacing in pixels | `4` | `2` to `8` |
| `MARGIN` | Left/right margin in pixels | `8` | `4` to `16` |

### Example .env File

```bash
# Printer Configuration
PRINTER_IP=192.168.1.100
PRINTER_PORT=9100
PAPER_WIDTH=80

# Font Configuration (optional)
CAFE_POS_FONT_PATH=/opt/cafe-pos/fonts/Roboto-Regular.ttf
FONT_SIZE=14.0
LINE_SPACING=4
MARGIN=8

# Application Configuration
LOG_LEVEL=info
DEBUG=false
```

## Testing Configuration

### Unit Test Configuration

```go
func TestPrinterConfiguration(t *testing.T) {
    config := &printing.PrinterConfig{
        Type:           printing.PrinterTypeBill,
        ConnectionType: printing.ConnectionTypeNetwork,
        IPAddress:      "192.168.1.100",
        Port:           9100,
        PaperWidth:     80,
    }
    
    printer, err := infraPrinting.NewESCPOSPrinter(config)
    require.NoError(t, err)
    assert.NotNil(t, printer)
}
```

### Integration Test Configuration

```go
func TestPrinterIntegration(t *testing.T) {
    // Skip if printer not available
    if os.Getenv("PRINTER_IP") == "" {
        t.Skip("PRINTER_IP not set, skipping integration test")
    }
    
    config := &printing.PrinterConfig{
        Type:           printing.PrinterTypeBill,
        ConnectionType: printing.ConnectionTypeNetwork,
        IPAddress:      os.Getenv("PRINTER_IP"),
        Port:           9100,
        PaperWidth:     80,
    }
    
    printer, err := infraPrinting.NewESCPOSPrinter(config)
    require.NoError(t, err)
    
    err = printer.Connect()
    require.NoError(t, err)
    defer printer.Disconnect()
    
    // Test print
    err = printer.Print("Test Receipt\n================================\nTest Item: 100,000 VND")
    require.NoError(t, err)
}
```

## Troubleshooting

### Common Configuration Issues

1. **Printer Not Found**
   - Verify IP address and port
   - Check network connectivity
   - Ensure printer is powered on

2. **Font Loading Failed**
   - Install recommended system fonts
   - Check font file path and permissions
   - Verify font file is valid TrueType format

3. **Vietnamese Characters Not Rendering**
   - Ensure font supports Vietnamese Unicode
   - Try different font (Roboto, DejaVu Sans)
   - Check font file is not corrupted

4. **Paper Width Mismatch**
   - Verify paper width setting (58mm or 80mm)
   - Check printer physical paper size
   - Adjust configuration to match printer

## Best Practices

1. **Use Environment Variables** for deployment-specific settings
2. **Use System Fonts** when possible for simplicity
3. **Test Configuration** before production deployment
4. **Monitor Printer Status** regularly
5. **Handle Errors Gracefully** with proper logging
6. **Document Custom Settings** for team reference
7. **Version Control** configuration files
8. **Backup Configurations** before changes

## Additional Resources

- [Font Deployment Guide](./FONT_DEPLOYMENT.md)
- [Font Manager Implementation](./FONT_MANAGER_IMPLEMENTATION.md)
- [Vietnamese Image-Based Printing Design](../../.kiro/specs/vietnamese-image-printing/design.md)
