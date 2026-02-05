// MongoDB Indexes for Cash Handover Feature
// Run this script to create all necessary indexes

print("Creating indexes for cash_handovers collection...");

db.cash_handovers.createIndex({ "waiter_shift_id": 1 });
print("✓ Index created: waiter_shift_id");

db.cash_handovers.createIndex({ "cashier_shift_id": 1 });
print("✓ Index created: cashier_shift_id");

db.cash_handovers.createIndex({ "status": 1 });
print("✓ Index created: status");

db.cash_handovers.createIndex({ "handover_at": -1 });
print("✓ Index created: handover_at (descending)");

db.cash_handovers.createIndex({ "cashier_id": 1, "status": 1 });
print("✓ Index created: cashier_id + status (compound)");

db.cash_handovers.createIndex({ "requires_approval": 1, "status": 1 });
print("✓ Index created: requires_approval + status (compound)");

db.cash_handovers.createIndex({ "waiter_id": 1, "handover_at": -1 });
print("✓ Index created: waiter_id + handover_at (compound)");

print("\nCreating indexes for cash_discrepancies collection...");

db.cash_discrepancies.createIndex({ "handover_id": 1 });
print("✓ Index created: handover_id");

db.cash_discrepancies.createIndex({ "resolution_status": 1 });
print("✓ Index created: resolution_status");

db.cash_discrepancies.createIndex({ "requires_manager_approval": 1 });
print("✓ Index created: requires_manager_approval");

db.cash_discrepancies.createIndex({ "created_at": -1 });
print("✓ Index created: created_at (descending)");

db.cash_discrepancies.createIndex({ "waiter_shift_id": 1 });
print("✓ Index created: waiter_shift_id");

db.cash_discrepancies.createIndex({ "cashier_shift_id": 1 });
print("✓ Index created: cashier_shift_id");

print("\nVerifying indexes...");

print("\nCash Handovers Indexes:");
printjson(db.cash_handovers.getIndexes());

print("\nCash Discrepancies Indexes:");
printjson(db.cash_discrepancies.getIndexes());

print("\n✅ All indexes created successfully!");
print("\nUsage: mongosh <database_name> < mongodb-handover-indexes.js");
