#!/bin/bash

echo "=== Test Logo Rendering Integration ==="
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if backend is running
echo "1. Checking if backend is running..."
if curl -s http://localhost:3000/api/health > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Backend is running${NC}"
else
    echo -e "${RED}✗ Backend is not running${NC}"
    echo "Please start the backend first: cd backend && go run main.go"
    exit 1
fi

# Check if logo exists
echo ""
echo "2. Checking if logo file exists..."
if [ -f "./uploads/logos/logo.png" ] || [ -f "./uploads/logos/logo.jpg" ]; then
    echo -e "${GREEN}✓ Logo file found${NC}"
    ls -lh ./uploads/logos/
else
    echo -e "${YELLOW}⚠ No logo file found in ./uploads/logos/${NC}"
    echo "Please upload a logo via the UI at http://localhost:5173/#/print-management"
fi

# Check shop settings
echo ""
echo "3. Checking shop settings..."
SETTINGS=$(curl -s http://localhost:3000/api/settings)
SHOW_LOGO=$(echo $SETTINGS | grep -o '"show_logo":[^,}]*' | cut -d':' -f2)
LOGO_URL=$(echo $SETTINGS | grep -o '"logo_url":"[^"]*"' | cut -d'"' -f4)

if [ "$SHOW_LOGO" = "true" ]; then
    echo -e "${GREEN}✓ show_logo is enabled${NC}"
else
    echo -e "${YELLOW}⚠ show_logo is disabled${NC}"
fi

if [ ! -z "$LOGO_URL" ]; then
    echo -e "${GREEN}✓ logo_url is set: $LOGO_URL${NC}"
else
    echo -e "${YELLOW}⚠ logo_url is not set${NC}"
fi

# Check if template with [LOGO] marker exists
echo ""
echo "4. Checking for template with [LOGO] marker..."
TEMPLATES=$(curl -s http://localhost:3000/api/manager/print-templates?type=BILL)
HAS_LOGO_MARKER=$(echo $TEMPLATES | grep -o '\[LOGO\]')

if [ ! -z "$HAS_LOGO_MARKER" ]; then
    echo -e "${GREEN}✓ Template with [LOGO] marker found${NC}"
else
    echo -e "${YELLOW}⚠ No template with [LOGO] marker found${NC}"
    echo "You can create one using the template from BILL_TEMPLATE_WITH_LOGO.txt"
fi

# Test preview functionality
echo ""
echo "5. Testing template preview..."
TEMPLATE_ID=$(echo $TEMPLATES | grep -o '"_id":"[^"]*"' | head -1 | cut -d'"' -f4)

if [ ! -z "$TEMPLATE_ID" ]; then
    echo "Testing preview for template: $TEMPLATE_ID"
    
    # Get template details
    TEMPLATE=$(curl -s http://localhost:3000/api/manager/print-templates/$TEMPLATE_ID)
    CONTENT=$(echo $TEMPLATE | grep -o '"content":"[^"]*"' | cut -d'"' -f4)
    TYPE=$(echo $TEMPLATE | grep -o '"type":"[^"]*"' | cut -d'"' -f4)
    
    # Test preview
    PREVIEW_RESULT=$(curl -s -X POST http://localhost:3000/api/manager/print-templates/$TEMPLATE_ID/preview \
        -H "Content-Type: application/json" \
        -d "{\"content\":\"$CONTENT\",\"type\":\"$TYPE\"}")
    
    if echo $PREVIEW_RESULT | grep -q '"success":true'; then
        echo -e "${GREEN}✓ Preview API works${NC}"
    else
        echo -e "${RED}✗ Preview API failed${NC}"
        echo "Response: $PREVIEW_RESULT"
    fi
else
    echo -e "${YELLOW}⚠ No templates found to test preview${NC}"
fi

echo ""
echo "=== Summary ==="
echo "To test logo rendering:"
echo "1. Upload a logo at http://localhost:5173/#/print-management"
echo "2. Create a template with [LOGO] marker (copy from BILL_TEMPLATE_WITH_LOGO.txt)"
echo "3. Set the template as default"
echo "4. Create a new order and check if logo appears in the printed bill"
echo ""
