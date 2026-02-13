#!/bin/bash

# Auto seed menu variants without confirmation

echo "🌱 Seeding menu variants data..."
echo ""

cd backend

# Run seed with auto-confirm
echo "y" | go run cmd/seed-menu-variants/main.go

cd ..

echo ""
echo "✅ Seed complete!"
echo ""
echo "📊 Next steps:"
echo "  1. Calculate costs: ./calculate-all-menu-costs.sh"
echo "  2. View analysis: ./view-menu-cost-analysis.sh"
