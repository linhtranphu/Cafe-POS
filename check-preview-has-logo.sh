#!/bin/bash

echo "Checking if preview PNG files have logo..."
echo ""

# Find latest preview files
RAW_FILE=$(ls -t backend/raw_preview_html_template_*.png 2>/dev/null | head -1)
PREVIEW_FILE=$(ls -t backend/preview_html_template_*.png 2>/dev/null | head -1)

if [ -z "$RAW_FILE" ]; then
    echo "❌ No raw preview file found"
    echo "Run: ./test-chromedp-preview.sh first"
    exit 1
fi

echo "Latest files:"
echo "  Raw: $RAW_FILE"
echo "  Preview: $PREVIEW_FILE"
echo ""

# Check file sizes
echo "File sizes:"
ls -lh "$RAW_FILE" "$PREVIEW_FILE" 2>/dev/null
echo ""

# Check image dimensions
echo "Image dimensions:"
file "$RAW_FILE" "$PREVIEW_FILE"
echo ""

echo "To visually check if logo appears:"
echo "  open $RAW_FILE"
echo "  open $PREVIEW_FILE"
echo ""

# Try to detect if image has significant content (not all white/black)
echo "Analyzing image content..."
if command -v identify &> /dev/null; then
    echo "Raw preview stats:"
    identify -verbose "$RAW_FILE" 2>/dev/null | grep -E "mean|standard deviation" | head -4
    echo ""
    echo "Processed preview stats:"
    identify -verbose "$PREVIEW_FILE" 2>/dev/null | grep -E "mean|standard deviation" | head -4
else
    echo "Install ImageMagick to analyze image content: brew install imagemagick"
fi
