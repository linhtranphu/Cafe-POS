#!/bin/bash

# Fix Admin Role Script
# This script updates the admin user role from cashier to manager

echo "🔧 Fixing admin user role..."
echo ""

# MongoDB connection
MONGO_URI="mongodb://localhost:27017"
DB_NAME="cafe_pos"

# Update admin user role to manager
mongosh "$MONGO_URI/$DB_NAME" --eval '
db.users.updateOne(
  { username: "admin" },
  { $set: { role: "manager" } }
)
'

echo ""
echo "✅ Admin role updated to manager"
echo ""
echo "Now:"
echo "1. Logout from the app"
echo "2. Login again with admin/admin123"
echo "3. You should see all manager menus"
echo ""
