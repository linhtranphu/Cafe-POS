#!/bin/bash

echo "🚀 Creating new bill template with logo and table..."
cd backend/cmd/create-logo-bill-template
go run main.go
cd ../../..
echo ""
echo "✅ Done! Template created."
echo "📋 Next steps:"
echo "   1. Go to http://localhost:5173/#/print-management"
echo "   2. Click 'Templates' tab"
echo "   3. Find 'Bill với Logo và Bảng'"
echo "   4. Set it as default template"
