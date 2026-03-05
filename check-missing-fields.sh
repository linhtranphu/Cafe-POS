#!/bin/bash

# Script to check what's missing in production database
# Safe to run - only reads data, doesn't modify anything

set -e

echo "=========================================="
echo "🔍 Checking Missing Fields/Collections"
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
echo "🔗 Connecting to MongoDB..."
echo ""

# Check if MongoDB container is running
if ! docker ps | grep -q mongodb; then
    echo "❌ MongoDB container is not running"
    exit 1
fi

MONGO_CONTAINER=$(docker ps --filter "name=mongodb" --format "{{.Names}}" | head -1)
echo "🐳 MongoDB container: $MONGO_CONTAINER"
echo ""

# Check shop_settings fields
echo "⚙️  Checking shop_settings fields..."
echo "-------------------------------------------"
docker exec $MONGO_CONTAINER mongosh $MONGO_DB \
    -u "$MONGO_USER" -p "$MONGO_PASS" --authenticationDatabase admin \
    --quiet --eval '
var settings = db.shop_settings.findOne();
if (!settings) {
    print("❌ shop_settings collection does not exist or is empty");
} else {
    print("✅ shop_settings exists");
    print("\nCurrent fields:");
    Object.keys(settings).forEach(function(key) {
        print("  ✓ " + key);
    });
    
    print("\nMissing fields:");
    var requiredFields = [
        "shop_name",
        "shop_address", 
        "shop_phone",
        "logo_url",
        "custom_message",
        "print_bridge_url",
        "show_logo",
        "show_address",
        "show_phone",
        "show_custom_message",
        "low_margin_threshold",
        "auto_print_enabled",
        "created_at",
        "updated_at"
    ];
    
    var missing = [];
    requiredFields.forEach(function(field) {
        if (settings[field] === undefined) {
            missing.push(field);
            print("  ❌ " + field);
        }
    });
    
    if (missing.length === 0) {
        print("  ✅ All required fields present");
    }
}
'

echo ""
echo "-------------------------------------------"
echo ""

# Check required collections
echo "📋 Checking required collections..."
echo "-------------------------------------------"
docker exec $MONGO_CONTAINER mongosh $MONGO_DB \
    -u "$MONGO_USER" -p "$MONGO_PASS" --authenticationDatabase admin \
    --quiet --eval '
var requiredCollections = [
    "users",
    "menu_items",
    "menu_categories", 
    "orders",
    "order_items",
    "ingredients",
    "ingredient_categories",
    "batch_definitions",
    "batch_records",
    "batch_usage_logs",
    "expenses",
    "expense_categories",
    "printer_configs",
    "print_templates",
    "print_jobs",
    "print_notifications",
    "shop_settings",
    "shifts",
    "cashier_shifts",
    "cash_handovers",
    "cash_discrepancies",
    "fund_transactions",
    "fund_handovers",
    "stock_history"
];

var existingCollections = db.getCollectionNames();

print("Existing: " + existingCollections.length + " collections");
print("Required: " + requiredCollections.length + " collections\n");

var missing = [];
requiredCollections.forEach(function(collName) {
    if (existingCollections.indexOf(collName) === -1) {
        missing.push(collName);
        print("  ❌ Missing: " + collName);
    } else {
        var count = db[collName].countDocuments();
        print("  ✓ " + collName + " (" + count + " docs)");
    }
});

if (missing.length === 0) {
    print("\n✅ All required collections exist");
} else {
    print("\n⚠️  Missing " + missing.length + " collections");
}
'

echo ""
echo "-------------------------------------------"
echo ""

# Check indexes
echo "🔑 Checking important indexes..."
echo "-------------------------------------------"
docker exec $MONGO_CONTAINER mongosh $MONGO_DB \
    -u "$MONGO_USER" -p "$MONGO_PASS" --authenticationDatabase admin \
    --quiet --eval '
function checkIndexes(collection, requiredIndexes) {
    var existing = db[collection].getIndexes();
    var existingNames = existing.map(function(idx) { return idx.name; });
    
    print("\n" + collection + ":");
    requiredIndexes.forEach(function(idxName) {
        if (existingNames.indexOf(idxName) !== -1) {
            print("  ✓ " + idxName);
        } else {
            print("  ❌ Missing: " + idxName);
        }
    });
}

checkIndexes("print_jobs", ["status_1_created_at_-1", "order_id_1", "printer_id_1"]);
checkIndexes("printer_configs", ["type_1_is_default_1"]);
checkIndexes("print_templates", ["type_1_is_default_1"]);
checkIndexes("orders", ["created_at_-1", "status_1"]);
'

echo ""
echo "-------------------------------------------"
echo ""

echo "✅ Check complete!"
echo ""
echo "📝 Next steps:"
echo "   - If missing fields/collections: sudo bash add-missing-only.sh"
echo "   - If all good: No action needed"
echo ""
