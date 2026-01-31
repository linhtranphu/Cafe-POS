// MongoDB script to fix admin role
// Run with: mongo cafe_pos fix-admin-role-mongo.js

print("🔧 Fixing admin user role...");
print("");

// Update admin user role to manager
var result = db.users.updateOne(
  { username: "admin" },
  { $set: { role: "manager" } }
);

print("Matched:", result.matchedCount);
print("Modified:", result.modifiedCount);
print("");

// Verify the change
var admin = db.users.findOne({ username: "admin" });
print("Admin user:");
print("  Username:", admin.username);
print("  Role:", admin.role);
print("  Name:", admin.name);
print("");

if (admin.role === "manager") {
  print("✅ Admin role successfully updated to manager");
} else {
  print("❌ Failed to update admin role");
}

print("");
print("Next steps:");
print("1. Logout from the app");
print("2. Clear browser cache (Ctrl+Shift+Delete)");
print("3. Run: localStorage.clear() in console");
print("4. Login again with admin/admin123");
print("5. You should see all manager menus");
print("");
