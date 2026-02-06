#!/bin/bash

# Revert back to container scroll pattern

echo "↩️  Reverting to Container Scroll Pattern"
echo "=========================================="
echo ""

files=(
  "frontend/src/views/DashboardView.vue"
  "frontend/src/views/FacilityManagementView.vue"
  "frontend/src/views/IngredientManagementView.vue"
  "frontend/src/views/MenuView.vue"
  "frontend/src/views/ProfileView.vue"
  "frontend/src/views/UserManagementView.vue"
)

count=0
for file in "${files[@]}"; do
  if [ -f "$file" ]; then
    echo "Processing: $file"
    
    # 1. Revert outer container
    sed -i '' 's/<div class="min-h-screen bg-gray-50">/<div class="h-screen w-screen overflow-hidden flex flex-col bg-gray-50">/g' "$file"
    
    # 2. Add back flex-shrink-0 to header
    sed -i '' 's/class="sticky top-0 z-40 bg-white shadow-sm">/class="sticky top-0 z-40 bg-white shadow-sm flex-shrink-0">/g' "$file"
    
    # 3. Revert content container
    sed -i '' 's/<div class="px-4 py-4 pb-24">/<div class="flex-1 overflow-y-auto px-4 py-4 pb-24">/g' "$file"
    
    echo "  ✅ Reverted"
    ((count++))
  fi
done

echo ""
echo "=========================================="
echo "✅ Reverted $count files back to container scroll"
