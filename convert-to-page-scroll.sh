#!/bin/bash

# Convert views from container scroll to page scroll pattern
# Based on ManagerShiftView which has perfect bottom padding

echo "🔄 Converting Views to Page Scroll Pattern"
echo "=========================================="
echo ""
echo "Pattern: h-screen (container) → min-h-screen (page)"
echo ""

# Files to convert
files=(
  "frontend/src/views/DashboardView.vue"
  "frontend/src/views/FacilityManagementView.vue"
  "frontend/src/views/IngredientManagementView.vue"
  "frontend/src/views/MenuView.vue"
  "frontend/src/views/ProfileView.vue"
  "frontend/src/views/UserManagementView.vue"
)

# Note: FacilityAddEditView is a form, keep as is

count=0
for file in "${files[@]}"; do
  if [ -f "$file" ]; then
    echo "Processing: $file"
    
    # 1. Change outer container
    sed -i '' 's/<div class="h-screen w-screen overflow-hidden flex flex-col bg-gray-50">/<div class="min-h-screen bg-gray-50">/g' "$file"
    
    # 2. Remove flex-shrink-0 from header
    sed -i '' 's/class="sticky top-0 z-40 bg-white shadow-sm flex-shrink-0"/class="sticky top-0 z-40 bg-white shadow-sm"/g' "$file"
    
    # 3. Change content container
    sed -i '' 's/<div class="flex-1 overflow-y-auto px-4 py-4 pb-24">/<div class="px-4 py-4 pb-24">/g' "$file"
    
    echo "  ✅ Converted"
    ((count++))
  else
    echo "  ❌ Not found: $file"
  fi
done

echo ""
echo "=========================================="
echo "✅ Converted $count files"
echo ""
echo "Changes made:"
echo "  1. h-screen w-screen overflow-hidden flex flex-col → min-h-screen"
echo "  2. Removed flex-shrink-0 from headers"
echo "  3. flex-1 overflow-y-auto → regular div (page scroll)"
echo ""
echo "Result:"
echo "  - Page scrolls naturally (like ManagerShiftView)"
echo "  - Bottom padding works correctly"
echo "  - BottomNav sits at bottom properly"
echo ""
echo "Next steps:"
echo "  1. Test each view on iPhone 14"
echo "  2. Verify bottom spacing is consistent"
echo "  3. Check pull-to-refresh still works"
