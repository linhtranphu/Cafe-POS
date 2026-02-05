#!/bin/bash

# Seed Menu Items Script
# This script seeds the database with sample menu items

echo "🌱 Seeding menu items..."
echo ""

# Run MongoDB script
mongosh "mongodb://admin:password123@localhost:27017/cafe_pos" \
  --authenticationDatabase admin \
  --file scripts/seed-menu-items.js

echo ""
echo "✅ Done! You can now create orders with these menu items."
echo ""
echo "📱 To test:"
echo "1. Login as waiter"
echo "2. Start a shift"
echo "3. Create an order"
echo "4. Select items from menu"
echo "5. Complete payment"
echo ""
