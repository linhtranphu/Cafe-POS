// Clean up incorrect maintenance records
// Run with: mongosh cafe_pos clean-incorrect-maintenance.js

print("=== Cleaning Incorrect Maintenance Records ===\n");

// Find the facility
const facility = db.facilities.findOne({ _id: ObjectId("697629c9d0ac2facdbb23baa") });
print("Facility:", facility.name, "-", facility.type);

// Find maintenance records for this facility
const maintenanceRecords = db.maintenance_records.find({ 
  facility_id: ObjectId("697629c9d0ac2facdbb23baa") 
}).toArray();

print("\nFound", maintenanceRecords.length, "maintenance records");
print("\nMaintenance records:");
maintenanceRecords.forEach(record => {
  print("  -", record.description);
});

// These records are for "Máy pha cà phê" but attached to "Bàn khách 2 chỗ"
// This is incorrect test data
print("\n⚠️  These maintenance records are incorrect!");
print("They describe 'Máy pha cà phê' but are attached to 'Bàn khách 2 chỗ'");

print("\nDeleting incorrect maintenance records...");
const result = db.maintenance_records.deleteMany({ 
  facility_id: ObjectId("697629c9d0ac2facdbb23baa") 
});

print("✅ Deleted", result.deletedCount, "incorrect maintenance records");

// Verify
const remaining = db.maintenance_records.countDocuments({ 
  facility_id: ObjectId("697629c9d0ac2facdbb23baa") 
});
print("Remaining maintenance records for this facility:", remaining);

print("\n=== Cleanup Complete ===");
