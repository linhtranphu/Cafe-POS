/**
 * Debug utility for menu cost issues
 * Usage: Open browser console and run: debugMenuCost('áddd')
 */

import api from '../services/api'

export async function debugMenuCost(menuName) {
  console.log(`🔍 Debugging menu item: "${menuName}"`)
  console.log('=====================================\n')
  
  try {
    // 1. Get all menu items
    console.log('📋 Step 1: Fetching menu items...')
    const menuResponse = await api.get('/manager/menu')
    const allMenuItems = menuResponse.data
    
    // Find the menu item
    const menuItem = allMenuItems.find(item => 
      item.name.toLowerCase().includes(menuName.toLowerCase())
    )
    
    if (!menuItem) {
      console.error(`❌ Menu item "${menuName}" not found`)
      console.log('Available menu items:', allMenuItems.map(m => m.name))
      return
    }
    
    console.log('✅ Found menu item:')
    console.log('  ID:', menuItem._id)
    console.log('  Name:', menuItem.name)
    console.log('  Price:', menuItem.price)
    console.log('  Category:', menuItem.category)
    console.log('  Has Variants:', menuItem.has_variants)
    console.log('  Current Cost:', menuItem.current_cost)
    console.log('  Cost Status:', menuItem.cost_status)
    console.log('  Cost Last Calculated:', menuItem.cost_last_calculated_at)
    console.log('')
    
    // 2. Check ingredients
    console.log('🥬 Step 2: Checking ingredients...')
    if (!menuItem.ingredients || menuItem.ingredients.length === 0) {
      console.error('❌ PROBLEM: Menu item has NO INGREDIENTS')
      console.log('   → Cannot calculate cost without ingredients')
      console.log('   → Solution: Add ingredients to this menu item')
      return
    }
    
    console.log(`✅ Menu item has ${menuItem.ingredients.length} ingredients:`)
    for (const ing of menuItem.ingredients) {
      console.log(`  - ${ing.name}: ${ing.quantity} ${ing.unit}`)
      console.log(`    Ingredient ID: ${ing.ingredient_id}`)
    }
    console.log('')
    
    // 3. Get all ingredients to check if they exist
    console.log('📦 Step 3: Verifying ingredients in database...')
    const ingredientsResponse = await api.get('/manager/ingredients')
    const allIngredients = ingredientsResponse.data
    
    let missingCount = 0
    let noCostCount = 0
    
    for (const ing of menuItem.ingredients) {
      const ingredientDoc = allIngredients.find(i => i._id === ing.ingredient_id)
      
      if (!ingredientDoc) {
        console.error(`  ❌ Ingredient "${ing.name}" (ID: ${ing.ingredient_id}) NOT FOUND in database`)
        missingCount++
      } else {
        if (!ingredientDoc.cost_per_unit || ingredientDoc.cost_per_unit === 0) {
          console.warn(`  ⚠️  Ingredient "${ing.name}" has NO COST (cost_per_unit = ${ingredientDoc.cost_per_unit})`)
          noCostCount++
        } else {
          console.log(`  ✅ Ingredient "${ing.name}": ${ingredientDoc.cost_per_unit} ${ingredientDoc.unit}`)
        }
      }
    }
    console.log('')
    
    // 4. Analyze the problem
    console.log('⚠️  Step 4: Problem Analysis')
    console.log('=====================================')
    
    if (missingCount > 0) {
      console.error(`🔴 CRITICAL: ${missingCount}/${menuItem.ingredients.length} ingredients DO NOT EXIST`)
      console.log('   → Cost Status will be: INCOMPLETE')
      console.log('   → Solution: Remove invalid ingredients and add correct ones')
    } else if (noCostCount > 0) {
      console.warn(`🟡 WARNING: ${noCostCount}/${menuItem.ingredients.length} ingredients have NO COST`)
      console.log('   → Cost Status will be: ESTIMATED or INCOMPLETE')
      console.log('   → Solution: Update cost_per_unit for these ingredients')
    } else {
      console.log('✅ All ingredients exist and have costs')
      
      if (!menuItem.current_cost || menuItem.current_cost === 0) {
        console.warn('⚠️  But current_cost is 0 or null')
        console.log('   → Cost may not have been calculated yet')
        console.log('   → Solution: Trigger recalculation by updating the menu item or ingredient')
      } else {
        console.log('✅ Everything looks good!')
      }
    }
    console.log('')
    
    // 5. Get cost breakdown
    console.log('💰 Step 5: Getting cost breakdown...')
    try {
      const costResponse = await api.get(`/manager/menu/costs/${menuItem._id}`)
      console.log('Cost breakdown:', costResponse.data)
    } catch (error) {
      console.error('❌ Failed to get cost breakdown:', error.response?.data || error.message)
    }
    
    console.log('\n=====================================')
    console.log('🏁 Debug complete!')
    
  } catch (error) {
    console.error('❌ Error during debug:', error)
    console.error('Details:', error.response?.data || error.message)
  }
}

// Make it available globally in browser console
if (typeof window !== 'undefined') {
  window.debugMenuCost = debugMenuCost
}
