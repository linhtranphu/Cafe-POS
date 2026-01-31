// MongoDB shell script to clean all ingredient data
// Run with: mongosh cafe_pos clean-ingredients-mongo.js

print('🧹 Cleaning all ingredient data...\n');

// Delete all ingredients
const ingredientsResult = db.ingredients.deleteMany({});
print(`🗑️  Deleted ${ingredientsResult.deletedCount} ingredients`);

// Delete all stock history
const historyResult = db.stock_history.deleteMany({});
print(`🗑️  Deleted ${historyResult.deletedCount} stock history records`);

print('\n✅ All ingredient data cleaned successfully!');
