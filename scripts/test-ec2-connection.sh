#!/bin/bash

# Test EC2 SSH Connection
# This script helps diagnose SSH connection issues

echo "🔍 Testing EC2 SSH Connection"
echo "=============================="
echo ""

EC2_HOST="13.212.27.222"
EC2_KEY="/Volumes/MacOS/users/tranphulinh/EC2PEM/OngTaPOS.pem"

# Check if key file exists
echo "1️⃣ Checking key file..."
if [ ! -f "$EC2_KEY" ]; then
  echo "❌ Key file not found: $EC2_KEY"
  exit 1
fi
echo "✅ Key file exists"
echo ""

# Check key permissions
echo "2️⃣ Checking key permissions..."
PERMS=$(stat -f "%Lp" "$EC2_KEY" 2>/dev/null || stat -c "%a" "$EC2_KEY" 2>/dev/null)
echo "   Permissions: $PERMS"
if [ "$PERMS" != "400" ] && [ "$PERMS" != "600" ]; then
  echo "⚠️  Warning: Permissions should be 400 or 600"
  echo "   Fix with: chmod 400 $EC2_KEY"
else
  echo "✅ Permissions OK"
fi
echo ""

# Check key format
echo "3️⃣ Checking key format..."
FIRST_LINE=$(head -1 "$EC2_KEY")
if [[ "$FIRST_LINE" == "-----BEGIN"* ]]; then
  echo "✅ Key format looks valid"
  echo "   Type: $FIRST_LINE"
else
  echo "❌ Invalid key format"
  exit 1
fi
echo ""

# Test connection with different usernames
echo "4️⃣ Testing SSH connections..."
echo ""

USERS=("ubuntu" "ec2-user" "admin" "root")

for USER in "${USERS[@]}"; do
  echo "   Testing: $USER@$EC2_HOST"
  
  # Try to connect with timeout
  timeout 5 ssh -o ConnectTimeout=5 \
    -o StrictHostKeyChecking=no \
    -o BatchMode=yes \
    -i "$EC2_KEY" \
    "$USER@$EC2_HOST" \
    "echo 'Success'" 2>&1 | grep -q "Success"
  
  if [ $? -eq 0 ]; then
    echo "   ✅ SUCCESS with user: $USER"
    echo ""
    echo "🎉 Working configuration:"
    echo "   User: $USER"
    echo "   Host: $EC2_HOST"
    echo "   Key: $EC2_KEY"
    echo ""
    echo "Update your script with:"
    echo "   EC2_USER=$USER"
    exit 0
  else
    echo "   ❌ Failed"
  fi
done

echo ""
echo "❌ All connection attempts failed"
echo ""
echo "Possible issues:"
echo "  1. Wrong PEM key for this EC2 instance"
echo "  2. EC2 instance not running"
echo "  3. Security group blocking SSH (port 22)"
echo "  4. Wrong EC2 IP address"
echo ""
echo "Troubleshooting steps:"
echo "  1. Check EC2 instance is running in AWS Console"
echo "  2. Verify the correct PEM key was used when launching instance"
echo "  3. Check Security Group allows SSH from your IP"
echo "  4. Verify EC2 public IP: $EC2_HOST"
