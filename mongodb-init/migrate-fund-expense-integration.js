// Migration Script: Fund-Expense Integration
// Purpose: Add fund integration fields to expenses and fund_transactions collections
// Date: 2026-03-05
// Requirements: 8.1, 8.2, 8.3, 8.4, 8.5, 8.6

print('🚀 Starting Fund-Expense Integration Migration...');

// Switch to cafe_pos database
db = db.getSiblingDB('cafe_pos');

print('📋 Step 1: Adding fund integration fields to expenses collection...');

// Add new fields to existing expenses with default values
const expensesUpdateResult = db.expenses.updateMany(
  {},
  {
    $set: {
      paid_from_fund: false,           // Default: not paid from fund
      fund_transaction_id: null        // Default: no fund transaction link
    }
  }
);

print(`✅ Updated ${expensesUpdateResult.modifiedCount} expense records`);

print('📋 Step 2: Creating indexes for expenses collection...');

// Index for filtering expenses paid from fund
db.expenses.createIndex(
  { "paid_from_fund": 1, "created_at": -1 },
  { name: "idx_expenses_paid_from_fund" }
);

// Index for looking up expense by fund transaction
db.expenses.createIndex(
  { "fund_transaction_id": 1 },
  { name: "idx_expenses_fund_transaction_id", sparse: true }
);

print('✅ Indexes created for expenses collection');

print('📋 Step 3: Adding source tracking fields to fund_transactions collection...');

// Add new fields to existing fund transactions with default values
const fundTxUpdateResult = db.fund_transactions.updateMany(
  {},
  {
    $set: {
      source_type: null,               // Default: no source (manual transaction)
      source_id: null                  // Default: no source link
    }
  }
);

print(`✅ Updated ${fundTxUpdateResult.modifiedCount} fund transaction records`);

print('📋 Step 4: Creating indexes for fund_transactions collection...');

// Unique compound index for source_type + source_id (sparse to allow nulls)
// This prevents duplicate fund transactions for the same source
// Note: Drop existing index if it exists to avoid duplicate key error
try {
  db.fund_transactions.dropIndex("idx_fund_transactions_source_unique");
  print('   Dropped existing idx_fund_transactions_source_unique');
} catch (e) {
  // Index doesn't exist, continue
}

db.fund_transactions.createIndex(
  { "source_type": 1, "source_id": 1 },
  { 
    name: "idx_fund_transactions_source_unique",
    unique: true,
    sparse: true,  // Allow null values, only enforce uniqueness when both fields are present
    partialFilterExpression: {
      source_type: { $ne: null },
      source_id: { $ne: null }
    }
  }
);

// Index for querying fund transactions by source type
db.fund_transactions.createIndex(
  { "source_type": 1, "timestamp": -1 },
  { name: "idx_fund_transactions_source_type" }
);

print('✅ Indexes created for fund_transactions collection');

print('📋 Step 5: Creating ingredient_restock_history collection...');

// Create the collection if it doesn't exist
if (!db.getCollectionNames().includes('ingredient_restock_history')) {
  db.createCollection('ingredient_restock_history');
  print('   Created ingredient_restock_history collection');
} else {
  print('   Collection ingredient_restock_history already exists');
}

print('📋 Step 6: Creating indexes for ingredient_restock_history collection...');

// Index for querying restock history by ingredient
db.ingredient_restock_history.createIndex(
  { "ingredient_id": 1, "created_at": -1 },
  { name: "idx_restock_history_ingredient" }
);

// Index for looking up restock by fund transaction
db.ingredient_restock_history.createIndex(
  { "fund_transaction_id": 1 },
  { name: "idx_restock_history_fund_transaction", sparse: true }
);

// Index for looking up restock by expense
db.ingredient_restock_history.createIndex(
  { "expense_id": 1 },
  { name: "idx_restock_history_expense", sparse: true }
);

print('✅ Indexes created for ingredient_restock_history collection');

print('📋 Step 7: Verification...');

// Verify expenses collection
const expensesCount = db.expenses.countDocuments({});
const expensesWithFundFields = db.expenses.countDocuments({
  paid_from_fund: { $exists: true },
  fund_transaction_id: { $exists: true }
});

print(`   Expenses total: ${expensesCount}`);
print(`   Expenses with fund fields: ${expensesWithFundFields}`);

// Verify fund_transactions collection
const fundTxCount = db.fund_transactions.countDocuments({});
const fundTxWithSourceFields = db.fund_transactions.countDocuments({
  source_type: { $exists: true },
  source_id: { $exists: true }
});

print(`   Fund transactions total: ${fundTxCount}`);
print(`   Fund transactions with source fields: ${fundTxWithSourceFields}`);

// Verify indexes
const expensesIndexes = db.expenses.getIndexes();
const fundTxIndexes = db.fund_transactions.getIndexes();

print(`   Expenses indexes: ${expensesIndexes.length}`);
print(`   Fund transactions indexes: ${fundTxIndexes.length}`);

// Verify ingredient_restock_history collection
const restockHistoryExists = db.getCollectionNames().includes('ingredient_restock_history');
print(`   Ingredient restock history collection exists: ${restockHistoryExists}`);

if (restockHistoryExists) {
  const restockHistoryIndexes = db.ingredient_restock_history.getIndexes();
  print(`   Ingredient restock history indexes: ${restockHistoryIndexes.length}`);
}

if (expensesWithFundFields === expensesCount && fundTxWithSourceFields === fundTxCount && restockHistoryExists) {
  print('✅ Migration completed successfully!');
  print('📊 Summary:');
  print(`   - ${expensesCount} expenses updated with fund integration fields`);
  print(`   - ${fundTxCount} fund transactions updated with source tracking fields`);
  print(`   - ingredient_restock_history collection created`);
  print(`   - 2 new indexes created on expenses collection`);
  print(`   - 2 new indexes created on fund_transactions collection`);
  print(`   - 3 new indexes created on ingredient_restock_history collection`);
} else {
  print('⚠️  Warning: Some records may not have been updated correctly');
  print('   Please verify the migration manually');
}

print('🎉 Fund-Expense Integration Migration Complete!');
