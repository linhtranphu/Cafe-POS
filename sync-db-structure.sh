#!/bin/bash

# Script to sync database structure from Production to Local
# This creates missing collections and indexes without touching data

set -e

echo "=========================================="
echo "🔄 Database Structure Sync"
echo "=========================================="
echo ""

# Check if running as root or with sudo
if [ "$EUID" -ne 0 ]; then 
    echo "⚠️  Please run with sudo"
    exit 1
fi

# MongoDB credentials
MONGO_USER="admin"
MONGO_PASS="password123"
MONGO_DB="cafe_pos"

echo "📍 Target: LOCAL database"
echo "🔗 MongoDB: mongodb://$MONGO_USER:***@localhost:27017/$MONGO_DB"
echo ""

# Confirm before proceeding
read -p "⚠️  This will sync structure to LOCAL database. Continue? (yes/NO) " -r
echo ""
if [[ ! $REPLY =~ ^[Yy][Ee][Ss]$ ]]; then
    echo "❌ Sync cancelled"
    exit 1
fi

echo ""
echo "🔄 Starting sync..."
echo ""

# Check if MongoDB container is running
if ! docker ps | grep -q mongodb; then
    echo "❌ MongoDB container is not running"
    echo "💡 Start it with: docker-compose up -d mongodb"
    exit 1
fi

MONGO_CONTAINER=$(docker ps --filter "name=mongodb" --format "{{.Names}}" | head -1)
echo "🐳 MongoDB container: $MONGO_CONTAINER"
echo ""

# Step 1: Ensure shop_settings exists with all required fields
echo "1️⃣  Syncing shop_settings..."
docker exec $MONGO_CONTAINER mongosh $MONGO_DB \
    -u "$MONGO_USER" -p "$MONGO_PASS" --authenticationDatabase admin \
    --quiet --eval '
var existing = db.shop_settings.findOne();
if (existing) {
    print("   ✅ Shop settings exists");
    
    // Add missing fields
    var updates = {};
    var hasUpdates = false;
    
    if (existing.print_bridge_url === undefined) {
        updates.print_bridge_url = "http://localhost:3001";
        hasUpdates = true;
    }
    if (existing.auto_print_enabled === undefined) {
        updates.auto_print_enabled = true;
        hasUpdates = true;
    }
    if (existing.show_logo === undefined) {
        updates.show_logo = false;
        hasUpdates = true;
    }
    if (existing.show_address === undefined) {
        updates.show_address = true;
        hasUpdates = true;
    }
    if (existing.show_phone === undefined) {
        updates.show_phone = true;
        hasUpdates = true;
    }
    if (existing.show_custom_message === undefined) {
        updates.show_custom_message = true;
        hasUpdates = true;
    }
    
    if (hasUpdates) {
        updates.updated_at = new Date();
        db.shop_settings.updateOne({ _id: existing._id }, { $set: updates });
        print("   ✅ Added missing fields: " + Object.keys(updates).join(", "));
    } else {
        print("   ✅ All fields present");
    }
} else {
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
        auto_print_enabled: true,
        created_at: now,
        updated_at: now
    });
    print("   ✅ Created shop_settings (ID: " + result.insertedId + ")");
}
'

echo ""

# Step 2: Ensure all required collections exist
echo "2️⃣  Ensuring required collections exist..."
docker exec $MONGO_CONTAINER mongosh $MONGO_DB \
    -u "$MONGO_USER" -p "$MONGO_PASS" --authenticationDatabase admin \
    --quiet --eval '
var requiredCollections = [
    "users",
    "menu_items", 
    "categories",
    "orders",
    "ingredients",
    "batches",
    "expenses",
    "printer_configs",
    "print_templates",
    "print_jobs",
    "shop_settings"
];

var existingCollections = db.getCollectionNames();

requiredCollections.forEach(function(collName) {
    if (existingCollections.indexOf(collName) === -1) {
        db.createCollection(collName);
        print("   ✅ Created collection: " + collName);
    } else {
        print("   ✓ " + collName + " exists");
    }
});
'

echo ""

# Step 3: Create indexes for print_jobs
echo "3️⃣  Creating indexes for print_jobs..."
docker exec $MONGO_CONTAINER mongosh $MONGO_DB \
    -u "$MONGO_USER" -p "$MONGO_PASS" --authenticationDatabase admin \
    --quiet --eval '
try {
    db.print_jobs.createIndex({ "status": 1, "created_at": -1 });
    db.print_jobs.createIndex({ "order_id": 1 });
    db.print_jobs.createIndex({ "printer_id": 1 });
    print("   ✅ Created print_jobs indexes");
} catch(e) {
    if (e.code === 85 || e.code === 86) {
        print("   ✓ Indexes already exist");
    } else {
        print("   ⚠️  " + e.message);
    }
}
'

echo ""

# Step 4: Create indexes for printer_configs
echo "4️⃣  Creating indexes for printer_configs..."
docker exec $MONGO_CONTAINER mongosh $MONGO_DB \
    -u "$MONGO_USER" -p "$MONGO_PASS" --authenticationDatabase admin \
    --quiet --eval '
try {
    db.printer_configs.createIndex({ "type": 1, "is_default": 1 });
    print("   ✅ Created printer_configs indexes");
} catch(e) {
    if (e.code === 85 || e.code === 86) {
        print("   ✓ Indexes already exist");
    } else {
        print("   ⚠️  " + e.message);
    }
}
'

echo ""

# Step 5: Create indexes for print_templates
echo "5️⃣  Creating indexes for print_templates..."
docker exec $MONGO_CONTAINER mongosh $MONGO_DB \
    -u "$MONGO_USER" -p "$MONGO_PASS" --authenticationDatabase admin \
    --quiet --eval '
try {
    db.print_templates.createIndex({ "type": 1, "is_default": 1 });
    print("   ✅ Created print_templates indexes");
} catch(e) {
    if (e.code === 85 || e.code === 86) {
        print("   ✓ Indexes already exist");
    } else {
        print("   ⚠️  " + e.message);
    }
}
'

echo ""

# Step 6: Create indexes for orders
echo "6️⃣  Creating indexes for orders..."
docker exec $MONGO_CONTAINER mongosh $MONGO_DB \
    -u "$MONGO_USER" -p "$MONGO_PASS" --authenticationDatabase admin \
    --quiet --eval '
try {
    db.orders.createIndex({ "created_at": -1 });
    db.orders.createIndex({ "status": 1 });
    print("   ✅ Created orders indexes");
} catch(e) {
    if (e.code === 85 || e.code === 86) {
        print("   ✓ Indexes already exist");
    } else {
        print("   ⚠️  " + e.message);
    }
}
'

echo ""

# Step 7: Summary
echo "7️⃣  Verification..."
docker exec $MONGO_CONTAINER mongosh $MONGO_DB \
    -u "$MONGO_USER" -p "$MONGO_PASS" --authenticationDatabase admin \
    --quiet --eval '
var collections = db.getCollectionNames();
print("   ✅ Total collections: " + collections.length);

var settings = db.shop_settings.findOne();
if (settings) {
    print("   ✅ shop_settings fields: " + Object.keys(settings).length);
    print("   ✅ Has print_bridge_url: " + (settings.print_bridge_url !== undefined));
}
'

echo ""
echo "=========================================="
echo "✅ Structure sync completed!"
echo "=========================================="
echo ""
echo "📝 Next steps:"
echo "   1. Compare with production: sudo bash compare-db-structure.sh"
echo "   2. Restart backend: docker-compose restart backend"
echo "   3. Test application"
echo ""
