#!/bin/bash

# Fix Bottom Padding - Reduce from pb-24 to pb-20
# pb-24 (96px) is too much for webapp mode without address bar
# pb-20 (80px) is better for BottomNav (~60px height)

echo "🔧 Fixing Bottom Padding in All Views"
echo "======================================"
echo ""
echo "Changing pb-24 (96px) → pb-20 (80px)"
echo ""

# List of files to update
files=(
  "frontend/src/views/BaristaView.vue"
  "frontend/src/views/CashierDashboard.vue"
  "frontend/src/views/CashierHandoverView.vue"
  "frontend/src/views/CashierReports.vue"
  "frontend/src/views/CashierShiftClosure.vue"
  "frontend/src/views/DashboardView.vue"
  "frontend/src/views/ExpenseManagementView.vue"
  "frontend/src/views/FacilityManagementView.vue"
  "frontend/src/views/IngredientManagementView.vue"
  "frontend/src/views/ManagerShiftView.vue"
  "frontend/src/views/MenuView.vue"
  "frontend/src/views/OrderView.vue"
  "frontend/src/views/ProfileView.vue"
  "frontend/src/views/ShiftView.vue"
  "frontend/src/views/UserManagementView.vue"
)

count=0
for file in "${files[@]}"; do
  if [ -f "$file" ]; then
    # Replace pb-24 with pb-20
    if grep -q "pb-24" "$file"; then
      sed -i '' 's/pb-24/pb-20/g' "$file"
      echo "✅ Updated: $file"
      ((count++))
    else
      echo "⏭️  Skipped: $file (no pb-24 found)"
    fi
  else
    echo "❌ Not found: $file"
  fi
done

echo ""
echo "======================================"
echo "✅ Updated $count files"
echo ""
echo "Summary:"
echo "  pb-24 (96px) → pb-20 (80px)"
echo "  Better for webapp mode (no address bar)"
echo "  BottomNav height: ~60px"
echo "  Visible space: ~20px"
echo ""
echo "Next steps:"
echo "  1. Test on iPhone 14 webapp mode"
echo "  2. Verify bottom spacing looks good"
echo "  3. Check all views"
