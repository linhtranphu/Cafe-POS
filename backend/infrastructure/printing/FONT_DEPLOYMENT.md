# Font File Deployment Guide

## Overview

The Vietnamese image-based printing system requires TrueType fonts that support Vietnamese Unicode characters. This guide explains how to deploy and configure fonts for the printing system.

## Font Requirements

The system requires fonts that support:
- Vietnamese Unicode characters (U+0000 to U+024F)
- All Vietnamese diacritical marks (tones and accents)
- Standard Latin characters

### Recommended Fonts

1. **Arial Unicode MS** (Windows, macOS)
   - Comprehensive Unicode support
   - Excellent Vietnamese character coverage
   - Professional appearance

2. **Roboto** (Linux, Android, cross-platform)
   - Modern, clean design
   - Good Vietnamese support
   - Open source (Apache License 2.0)

3. **DejaVu Sans** (Linux, cross-platform)
   - Excellent Unicode coverage
   - Open source (free to use)
   - Good fallback option

4. **Noto Sans** (cross-platform)
   - Google's comprehensive font family
   - Excellent Vietnamese support
   - Open source (SIL Open Font License)

## Deployment Options

### Option 1: System Fonts (Recommended)

The system automatically discovers and uses system fonts. No additional deployment is required if one of the recommended fonts is installed on the system.

**Automatic Font Discovery:**
- Windows: Searches `C:\Windows\Fonts\`
- macOS: Searches `/Library/Fonts/` and `/System/Library/Fonts/`
- Linux: Searches `/usr/share/fonts/truetype/`

**Advantages:**
- No additional deployment steps
- Fonts are already installed and maintained by the OS
- Automatic updates through OS updates

**Disadvantages:**
- Depends on OS font availability
- May vary across different systems

### Option 2: Bundled Font Files

Include font files in the application deployment package for consistent rendering across all systems.

**Steps:**

1. **Create fonts directory:**
   ```bash
   mkdir -p /opt/cafe-pos/fonts
   ```

2. **Copy font files:**
   ```bash
   # Example: Copy Roboto font
   cp Roboto-Regular.ttf /opt/cafe-pos/fonts/
   cp Roboto-Bold.ttf /opt/cafe-pos/fonts/
   ```

3. **Configure font paths in application:**
   ```go
   // In your printer configuration
   config := &printing.PrinterConfig{
       Type:           printing.PrinterTypeBill,
       ConnectionType: printing.ConnectionTypeNetwork,
       IPAddress:      "192.168.1.100",
       Port:           9100,
       PaperWidth:     80,
   }
   
   // Font path will be used by TextRenderer
   rendererConfig := &RendererConfig{
       PixelWidth:  pixelWidth,
       FontPath:    "/opt/cafe-pos/fonts/Roboto-Regular.ttf",
       FontSize:    14.0,
       LineSpacing: 4,
       Margin:      8,
   }
   ```

**Advantages:**
- Consistent rendering across all systems
- No dependency on OS fonts
- Full control over font versions

**Disadvantages:**
- Requires font file distribution
- Need to manage font licenses
- Larger deployment package

### Option 3: Environment Variable Configuration

Configure font paths through environment variables for flexible deployment.

**Steps:**

1. **Set environment variable:**
   ```bash
   export CAFE_POS_FONT_PATH="/path/to/fonts/Roboto-Regular.ttf"
   ```

2. **Read in application:**
   ```go
   fontPath := os.Getenv("CAFE_POS_FONT_PATH")
   if fontPath == "" {
       fontPath = "" // Empty triggers system font discovery
   }
   ```

**Advantages:**
- Flexible configuration
- Easy to change without code modification
- Supports different environments (dev, staging, prod)

**Disadvantages:**
- Requires environment configuration
- May be forgotten during deployment

## Font Licensing

### Open Source Fonts (Recommended)

Use open source fonts to avoid licensing issues:

1. **Roboto** - Apache License 2.0
   - Free for commercial use
   - Can be bundled with application
   - Download: https://fonts.google.com/specimen/Roboto

2. **DejaVu Sans** - Free license
   - Free for any use
   - Can be bundled with application
   - Download: https://dejavu-fonts.github.io/

3. **Noto Sans** - SIL Open Font License
   - Free for commercial use
   - Can be bundled with application
   - Download: https://fonts.google.com/noto/specimen/Noto+Sans

### Commercial Fonts

If using commercial fonts (e.g., Arial Unicode MS):
- Ensure proper licensing for deployment
- Check if font can be bundled with application
- Consider per-installation licensing requirements

## Installation Instructions

### Windows

1. **System Font Installation:**
   - Download font file (e.g., `Roboto-Regular.ttf`)
   - Right-click font file → "Install"
   - Font will be available in `C:\Windows\Fonts\`

2. **Application Bundle:**
   - Create `fonts` folder in application directory
   - Copy font files to `fonts` folder
   - Configure application to use bundled fonts

### macOS

1. **System Font Installation:**
   - Download font file
   - Double-click font file
   - Click "Install Font" in Font Book
   - Font will be available in `/Library/Fonts/`

2. **Application Bundle:**
   - Create `fonts` folder in application bundle
   - Copy font files to `fonts` folder
   - Configure application to use bundled fonts

### Linux

1. **System Font Installation:**
   ```bash
   # Ubuntu/Debian
   sudo apt-get install fonts-roboto
   
   # Or manually install
   sudo mkdir -p /usr/share/fonts/truetype/roboto
   sudo cp Roboto-*.ttf /usr/share/fonts/truetype/roboto/
   sudo fc-cache -f -v
   ```

2. **Application Bundle:**
   ```bash
   mkdir -p /opt/cafe-pos/fonts
   cp Roboto-*.ttf /opt/cafe-pos/fonts/
   ```

## Verification

### Test Font Loading

Run the following test to verify font loading:

```bash
cd backend
go test ./infrastructure/printing/... -v -run TestLoadFonts
```

### Test Vietnamese Character Rendering

Run integration tests to verify Vietnamese character support:

```bash
cd backend
go test ./infrastructure/printing/... -v -run TestIntegration_VietnameseCharacters
```

## Troubleshooting

### Font Not Found Error

**Error:** `font loading error: no configured fonts found and system font discovery failed`

**Solutions:**
1. Install a recommended font on the system
2. Provide explicit font path in configuration
3. Check font file permissions (must be readable)

### Vietnamese Characters Not Rendering

**Error:** Characters appear as boxes or question marks

**Solutions:**
1. Verify font supports Vietnamese Unicode range
2. Check font file is not corrupted
3. Try a different font (e.g., Roboto or DejaVu Sans)

### Font File Permissions

**Error:** `font file access error: cannot access /path/to/font.ttf`

**Solutions:**
```bash
# Make font file readable
chmod 644 /path/to/font.ttf

# Ensure directory is accessible
chmod 755 /path/to/fonts/
```

## Configuration Examples

### Development Environment

```go
// Use system fonts for development
rendererConfig := &RendererConfig{
    PixelWidth:  pixelWidth,
    FontPath:    "", // Empty = system font discovery
    FontSize:    14.0,
    LineSpacing: 4,
    Margin:      8,
}
```

### Production Environment (Docker)

```dockerfile
# Dockerfile
FROM golang:1.21-alpine

# Install fonts
RUN apk add --no-cache fontconfig ttf-dejavu

# Or copy bundled fonts
COPY fonts/ /usr/share/fonts/truetype/cafe-pos/
RUN fc-cache -f -v

# Copy application
COPY . /app
WORKDIR /app
```

### Production Environment (Bare Metal)

```bash
# Install fonts during deployment
sudo apt-get install fonts-roboto fonts-dejavu

# Or copy bundled fonts
sudo mkdir -p /opt/cafe-pos/fonts
sudo cp fonts/*.ttf /opt/cafe-pos/fonts/
sudo chmod 644 /opt/cafe-pos/fonts/*.ttf
```

## Best Practices

1. **Use System Fonts When Possible**
   - Reduces deployment complexity
   - Automatic updates
   - No licensing concerns

2. **Bundle Fonts for Consistency**
   - Ensures consistent rendering
   - No dependency on OS
   - Better for production

3. **Test on Target Platform**
   - Verify font availability
   - Test Vietnamese character rendering
   - Check performance

4. **Document Font Requirements**
   - List required fonts in README
   - Provide installation instructions
   - Include troubleshooting guide

5. **Use Open Source Fonts**
   - Avoid licensing issues
   - Free to bundle and distribute
   - Good Vietnamese support

## Font Fallback Chain

The system uses the following fallback chain:

1. **Configured font path** (if provided)
2. **System fonts** (automatic discovery)
   - Windows: Arial Unicode MS → Arial → DejaVu Sans
   - macOS: Arial Unicode MS → Arial
   - Linux: DejaVu Sans → Liberation Sans → Noto Sans → Roboto

3. **Error** if no suitable font found

## Performance Considerations

- **Font Loading:** Fonts are loaded once during printer initialization
- **Memory Usage:** Each font face uses ~1-2 MB of memory
- **Rendering Speed:** Font rendering is fast (~10-50ms per receipt)
- **Caching:** Font faces are cached for the lifetime of the printer instance

## Support

For font-related issues:
1. Check this documentation
2. Run font loading tests
3. Verify font file integrity
4. Check system font availability
5. Review application logs for detailed error messages
