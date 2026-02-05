// MongoDB Migration Script for Cash Handover Feature
// Adds new fields to existing shifts and cashier_shifts collections

print("=== Cash Handover Migration Script ===\n");

// Migration for shifts collection (Waiter/Barista shifts)
print("Migrating shifts collection...");

const shiftsResult = db.shifts.updateMany(
  { 
    type: { $in: ["WAITER", "BARISTA"] },
    // Only update if fields don't exist
    current_cash: { $exists: false }
  },
  {
    $set: {
      current_cash: 0,
      handed_over_cash: 0,
      remaining_cash: 0,
      total_discrepancy: 0,
      handover_count: 0
    }
  }
);

print(`✓ Shifts updated: ${shiftsResult.modifiedCount} documents`);
print(`  Matched: ${shiftsResult.matchedCount}`);

// Migration for cashier_shifts collection
print("\nMigrating cashier_shifts collection...");

const cashierShiftsResult = db.cashier_shifts.updateMany(
  {
    // Only update if fields don't exist
    received_cash: { $exists: false }
  },
  {
    $set: {
      received_cash: 0,
      total_discrepancy: 0,
      handover_count: 0,
      discrepancy_count: 0
    }
  }
);

print(`✓ Cashier shifts updated: ${cashierShiftsResult.modifiedCount} documents`);
print(`  Matched: ${cashierShiftsResult.matchedCount}`);

// Verification
print("\n=== Verification ===");

print("\nSample shift document:");
const sampleShift = db.shifts.findOne({ type: "WAITER" });
if (sampleShift) {
  print("Fields present:");
  print(`  - current_cash: ${sampleShift.current_cash !== undefined ? '✓' : '✗'}`);
  print(`  - handed_over_cash: ${sampleShift.handed_over_cash !== undefined ? '✓' : '✗'}`);
  print(`  - remaining_cash: ${sampleShift.remaining_cash !== undefined ? '✓' : '✗'}`);
  print(`  - total_discrepancy: ${sampleShift.total_discrepancy !== undefined ? '✓' : '✗'}`);
  print(`  - handover_count: ${sampleShift.handover_count !== undefined ? '✓' : '✗'}`);
} else {
  print("  No waiter shifts found");
}

print("\nSample cashier shift document:");
const sampleCashierShift = db.cashier_shifts.findOne();
if (sampleCashierShift) {
  print("Fields present:");
  print(`  - received_cash: ${sampleCashierShift.received_cash !== undefined ? '✓' : '✗'}`);
  print(`  - total_discrepancy: ${sampleCashierShift.total_discrepancy !== undefined ? '✓' : '✗'}`);
  print(`  - handover_count: ${sampleCashierShift.handover_count !== undefined ? '✓' : '✗'}`);
  print(`  - discrepancy_count: ${sampleCashierShift.discrepancy_count !== undefined ? '✓' : '✗'}`);
} else {
  print("  No cashier shifts found");
}

// Statistics
print("\n=== Statistics ===");
print(`Total shifts in database: ${db.shifts.countDocuments()}`);
print(`Total cashier shifts in database: ${db.cashier_shifts.countDocuments()}`);
print(`Waiter shifts: ${db.shifts.countDocuments({ type: "WAITER" })}`);
print(`Barista shifts: ${db.shifts.countDocuments({ type: "BARISTA" })}`);

print("\n✅ Migration completed successfully!");
print("\nUsage: mongosh <database_name> < mongodb-handover-migration.js");
print("\nNote: This script is idempotent - safe to run multiple times");
