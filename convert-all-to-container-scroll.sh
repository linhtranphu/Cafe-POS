#!/bin/bash

# Convert all page scroll views to container scroll pattern
# Container scroll provides consistent bottom spacing

echo "🔄 Converting All Views to Container Scroll"
echo "============================================"
echo ""
echo "Pattern: min-h-screen (page) → h-screen (container)"
echo ""

# Files to convert (exclude LoginView)
files=(
  "frontend/src/views/BaristaView.vue"
  "frontend/src/views/CashierDashboard.vue"
  "frontend/src/views/CashierHandoverView.vue"
  "frontend/src/views/CashierReports.vue"
  "frontend/src/views/CashierShiftClosure.vue"
  "frontend/src/views/ExpenseManagementView.vue"
  "frontend/src/views/ManagerShiftView.vue"
  "frontend/src/views/OrderView.vue"
  "frontend/src/views/ShiftView.vue"
)

count=0
for file in "${files[@]}"; do
  if [ -f "$file" ]; then
    echo "Processing: $file"
    
    # 1. Change outer container from min-h-screen to h-screen
    sed -i '' 's/<div class="min-h-screen bg-gray-50">/<div class="h-screen w-screen overflow-hidden flex flex-col bg-gray-50">/g' "$file"
    
    # 2. Add flex-shrink-0 to sticky headers
    sed -i '' 's/class="sticky top-0 z-40 bg-white shadow-sm">/class="sticky top-0 z-40 bg-white shadow-sm flex-shrink-0">/g' "$file"
    
    # 3. Change content div to flex-1 overflow-y-auto
    # Match various patterns of content divs
    sed -i '' 's/<div class="px-4 py-4 pb-24">/<div class="flex-1 overflow-y-auto px-4 py-4 pb-24">/g' "$file"
    
    echo "  ✅ Converted"
    ((count++))
  else
    echo "  ❌ Not found: $file"
  fi
done

echo ""
echo "============================================"
echo "✅ Converted $count files to container scroll"
echo ""
echo "Changes made:"
echo "  1. min-h-screen → h-screen w-screen overflow-hidden flex flex-col"
echo "  2. Added flex-shrink-0 to sticky headers"
echo "  3. Content div → flex-1 overflow-y-auto"
echo ""
echo "Benefits:"
echo "  - Consistent bottom spacing (pb-24 works perfectly)"
echo "  - No safe area needed on BottomNav"
echo "  - Better scroll performance"
echo "  - Unified pattern across all views"
echo ""
echo "Next steps:"
echo "  1. Test each view on iPhone 14"
echo "  2. Verify bottom spacing is consistent"
echo "  3. Check pull-to-refresh still works"
echo "  4. Verify scroll behavior is smooth"
