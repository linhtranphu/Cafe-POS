#!/bin/bash

echo "=== Test Print Simple Text ==="
echo ""

# Check if print bridge is running
if ! curl -s http://localhost:3001/health > /dev/null 2>&1; then
    echo "❌ Print bridge is not running at http://localhost:3001"
    echo "Please start it with: docker start local-print-bridge"
    exit 1
fi

echo "✅ Print bridge is running"
echo ""

# Get printer IP from user
read -p "Enter printer IP address (e.g., 192.168.1.100): " PRINTER_IP

if [ -z "$PRINTER_IP" ]; then
    echo "❌ Printer IP is required"
    exit 1
fi

echo ""
echo "Testing with printer: $PRINTER_IP"
echo ""

# Create simple ESC/POS text command
# ESC @ (initialize) + text + line feed + cut
TEXT="Hello from Cafe POS!\nTest print successful.\n\n"

# Convert to hex bytes
ESCPOS_HEX="1b40"  # ESC @ (initialize)
ESCPOS_HEX="${ESCPOS_HEX}$(echo -n "$TEXT" | xxd -p | tr -d '\n')"  # Text
ESCPOS_HEX="${ESCPOS_HEX}1b6403"  # ESC d 3 (feed 3 lines)
ESCPOS_HEX="${ESCPOS_HEX}1d5600"  # GS V 0 (cut)

# Convert hex to binary then to base64
BASE64_DATA=$(echo "$ESCPOS_HEX" | xxd -r -p | base64 | tr -d '\n')

echo "Sending simple text to printer..."
RESPONSE=$(curl -s -X POST http://localhost:3001/print \
    -H "Content-Type: application/json" \
    -d "{
        \"printer_ip\": \"$PRINTER_IP\",
        \"data\": \"$BASE64_DATA\",
        \"encoding\": \"base64\"
    }")

echo "Response: $RESPONSE"

if echo "$RESPONSE" | grep -q "success"; then
    echo "✅ Text sent successfully"
    echo ""
    echo "Check your printer - it should print:"
    echo "  Hello from Cafe POS!"
    echo "  Test print successful."
else
    echo "❌ Failed to send"
fi

echo ""
echo "=== Test completed ==="
