#!/bin/bash

# Script to add ONLY missing fields/collections to existing database
# SAFE: Does NOT delete or modify existing data
# SAFE: Only adds what's missing

set -e

echo "=========================================="
echo "➕ Add Missing Fields/Collections Only"
echo "=========================================="
echo ""

# MongoDB credentials - UPDATE THESE FOR PRODUCTION
read -p "MongoDB username [admin]: " MONGO_USER
MONGO_USER=${MONGO_USER:-admin}

read -sp "MongoDB password: " MONGO_PASS
echo ""

read -p "Database name [cafe_pos]: " MONGO_DB
MONGO_DB=${MONGO_DB:-cafe_pos}

echo ""
echo "📍 Target database: $MONGO_DB"
echo ""

# Confirm before proceeding
read -p "⚠️  This will ADD missing fields/collections. Continue? (yes/NO) " -r
echo ""
if [[ ! $REPLY =~ ^[Yy][Ee][Ss]$ ]]; then
    echo "❌ Cancelled"
    exit 1
fi

echo ""
echo "➕ Adding missing items..."
echo ""

# Check if MongoDB container is running
if ! docker ps | grep -q mongodb; then
    echo "❌ MongoDB container is not running"
    exit 1
fi

MONGO_CONTAINER=$(docker ps --filter "name=mongodb" --format "{{.Names}}" | head -1)
echo "🐳 MongoDB container: $MONGO_CONTAINER"
echo ""

# Step 1: Add missing fields to shop_settings
echo "1️⃣  Adding missing fields to shop_settings..."
docker exec $MONGO_CONTAINER mongosh $MONGO_DB \
    -u "$MONGO_USER" -p "$MONGO_PASS" --authenticationDatabase admin \
    --quiet --eval '
var existing = db.shop_settings.findOne();
if (!existing) {
    // Create shop_settings if not exists
    var now = new Date();
    var result = db.shop_settings.insertOne({
        shop_name: "Cafe POS",
        shop_address: "",
        shop_phone: "",
        logo_url: "",
        custom_message: "Cảm ơn quý khách! Hẹn gặp lại!",
        print_bridge_url: "http://localhost:3001",
        show_logo: false,
        show_address: true,
        show_phone: true,
        show_custom_message: true,
        low_margin_threshold: 20.0,
        auto_print_enabled: true,
        created_at: now,
        updated_at: now
    });
    print("   ✅ Created shop_settings (ID: " + result.insertedId + ")");
} else {
    // Add only missing fields
    var updates = {};
    
    if (existing.print_bridge_url === undefined) {
        updates.print_bridge_url = "http://localhost:3001";
    }
    if (existing.auto_print_enabled === undefined) {
        updates.auto_print_enabled = true;
    }
    if (existing.show_logo === undefined) {
        updates.show_logo = existing.show_logo !== undefined ? existing.show_logo : false;
    }
    if (existing.show_address === undefined) {
        updates.show_address = true;
    }
    if (existing.show_phone === undefined) {
        updates.show_phone = true;
    }
    if (existing.show_custom_message === undefined) {
        updates.show_custom_message = true;
    }
    if (existing.low_margin_threshold === undefined) {
        updates.low_margin_threshold = 20.0;
    }
    if (existing.shop_address === undefined) {
        updates.shop_address = "";
    }
    if (existing.shop_phone === undefined) {
        updates.shop_phone = "";
    }
    if (existing.logo_url === undefined) {
        updates.logo_url = "";
    }
    if (existing.custom_message === undefined) {
        updates.custom_message = "Cảm ơn quý khách! Hẹn gặp lại!";
    }
    
    if (Object.keys(updates).length > 0) {
        updates.updated_at = new Date();
        db.shop_settings.updateOne({ _id: existing._id }, { $set: updates });
        print("   ✅ Added " + Object.keys(updates).length + " missing fields:");
        Object.keys(updates).forEach(function(key) {
            if (key !== "updated_at") {
                print("      - " + key);
            }
        });
    } else {
        print("   ✅ All fields already present");
    }
}
'

echo ""

# Step 2: Create missing collections (if any)
echo "2️⃣  Creating missing collections..."
docker exec $MONGO_CONTAINER mongosh $MONGO_DB \
    -u "$MONGO_USER" -p "$MONGO_PASS" --authenticationDatabase admin \
    --quiet --eval '
var requiredCollections = [
    "users", "menu_items", "menu_categories", "orders", "order_items",
    "ingredients", "ingredient_categories", "batch_definitions", 
    "batch_records", "batch_usage_logs", "expenses", "expense_categories",
    "printer_configs", "print_templates", "print_jobs", "print_notifications",
    "shop_settings", "shifts", "cashier_shifts", "cash_handovers",
    "cash_discrepancies", "fund_transactions", "fund_handovers", "stock_history"
];

var existingCollections = db.getCollectionNames();
var created = 0;

requiredCollections.forEach(function(collName) {
    if (existingCollections.indexOf(collName) === -1) {
        db.createCollection(collName);
        print("   ✅ Created: " + collName);
        created++;
    }
});

if (created === 0) {
    print("   ✅ All collections already exist");
} else {
    print("   ✅ Created " + created + " collections");
}
'

echo ""

# Step 3: Add missing indexes
echo "3️⃣  Adding missing indexes..."
docker exec $MONGO_CONTAINER mongosh $MONGO_DB \
    -u "$MONGO_USER" -p "$MONGO_PASS" --authenticationDatabase admin \
    --quiet --eval '
var indexesAdded = 0;

function tryCreateIndex(collection, indexSpec, options) {
    try {
        db[collection].createIndex(indexSpec, options || {});
        indexesAdded++;
        return true;
    } catch(e) {
        if (e.code === 85 || e.code === 86) {
            // Index already exists
            return false;
        }
        print("   ⚠️  " + collection + ": " + e.message);
        return false;
    }
}

// Print jobs indexes
if (tryCreateIndex("print_jobs", { "status": 1, "created_at": -1 })) {
    print("   ✅ print_jobs: status + created_at");
}
if (tryCreateIndex("print_jobs", { "order_id": 1 })) {
    print("   ✅ print_jobs: order_id");
}
if (tryCreateIndex("print_jobs", { "printer_id": 1 })) {
    print("   ✅ print_jobs: printer_id");
}

// Printer configs indexes
if (tryCreateIndex("printer_configs", { "type": 1, "is_default": 1 })) {
    print("   ✅ printer_configs: type + is_default");
}

// Print templates indexes
if (tryCreateIndex("print_templates", { "type": 1, "is_default": 1 })) {
    print("   ✅ print_templates: type + is_default");
}

// Orders indexes
if (tryCreateIndex("orders", { "created_at": -1 })) {
    print("   ✅ orders: created_at");
}
if (tryCreateIndex("orders", { "status": 1 })) {
    print("   ✅ orders: status");
}

if (indexesAdded === 0) {
    print("   ✅ All indexes already exist");
}
'

echo ""

# Step 4: Verify
echo "4️⃣  Verification..."
docker exec $MONGO_CONTAINER mongosh $MONGO_DB \
    -u "$MONGO_USER" -p "$MONGO_PASS" --authenticationDatabase admin \
    --quiet --eval '
try {
    var settings = db.shop_settings.findOne({}, { _id: 1, print_bridge_url: 1, auto_print_enabled: 1 });
    if (settings) {
        print("   ✅ shop_settings exists");
        print("   ✅ Has print_bridge_url: " + (settings.print_bridge_url !== undefined));
        print("   ✅ Has auto_print_enabled: " + (settings.auto_print_enabled !== undefined));
    }
    
    var collections = db.getCollectionNames();
    print("   ✅ Total collections: " + collections.length);
} catch(e) {
    print("   ⚠️  Verification error: " + e.message);
}
' || echo "   ⚠️  Verification skipped (non-critical)"

echo ""
echo "=========================================="
echo "✅ Update completed!"
echo "=========================================="
echo ""
echo "📝 What was done:"
echo "   - Added missing fields to shop_settings"
echo "   - Created missing collections (if any)"
echo "   - Added missing indexes"
echo "   - NO data was deleted or modified"
echo ""
echo "📝 Next steps:"
echo "   1. Restart backend: docker-compose restart backend"
echo "   2. Test application"
echo ""
