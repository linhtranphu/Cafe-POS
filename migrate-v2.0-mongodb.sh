#!/bin/bash

# Migration script for Cafe POS v2.0 - Direct MongoDB approach
# This script runs migration directly in MongoDB without needing Go code

set -e

echo "=========================================="
echo "🚀 Cafe POS v2.0 - Production Migration"
echo "=========================================="
echo ""

# Check if running as root or with sudo
if [ "$EUID" -ne 0 ]; then 
    echo "⚠️  Please run with sudo"
    exit 1
fi

# Load environment variables from .env if exists
if [ -f ".env" ]; then
    set -a
    source .env
    set +a
    echo "📄 Loaded .env file"
elif [ -f ".env.production" ]; then
    set -a
    source .env.production
    set +a
    echo "📄 Loaded .env.production file"
fi

# Get MongoDB credentials from environment or use defaults
MONGO_USER="${MONGO_INITDB_ROOT_USERNAME:-admin}"
MONGO_PASS="${MONGO_INITDB_ROOT_PASSWORD:-password123}"
MONGO_DB="${MONGODB_DATABASE:-cafe_pos}"

echo "📍 Environment: PRODUCTION"
echo "🔗 MongoDB: mongodb://$MONGO_USER:***@localhost:27017/$MONGO_DB"
echo ""

# Confirm before proceeding
read -p "⚠️  This will modify the PRODUCTION database. Are you sure? (yes/NO) " -r
echo ""
if [[ ! $REPLY =~ ^[Yy][Ee][Ss]$ ]]; then
    echo "❌ Migration cancelled"
    exit 1
fi

echo ""
echo "🔄 Starting migration..."
echo ""

# Check if MongoDB container is running
if ! docker ps | grep -q mongodb; then
    echo "❌ MongoDB container is not running"
    exit 1
fi

MONGO_CONTAINER=$(docker ps --filter "name=mongodb" --format "{{.Names}}" | head -1)
echo "🐳 MongoDB container: $MONGO_CONTAINER"
echo ""

# Migration 1: Create shop_settings if not exists
echo "1️⃣  Creating shop_settings..."
docker exec $MONGO_CONTAINER mongosh $MONGO_DB \
    -u "$MONGO_USER" -p "$MONGO_PASS" --authenticationDatabase admin \
    --quiet --eval '
var existing = db.shop_settings.findOne();
if (existing) {
    print("   ✅ Shop settings already exist (ID: " + existing._id + ")");
    // Update existing settings to add print_bridge_url if missing
    if (!existing.print_bridge_url) {
        db.shop_settings.updateOne(
            { _id: existing._id },
            { $set: { print_bridge_url: "http://localhost:3001", updated_at: new Date() } }
        );
        print("   ✅ Added print_bridge_url to existing settings");
    }
} else {
    var now = new Date();
    var result = db.shop_settings.insertOne({
        shop_name: "Cafe POS",
        shop_address: "123 Main Street",
        shop_phone: "0123-456-789",
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
    print("   ✅ Created shop settings (ID: " + result.insertedId + ")");
}
'

echo ""

# Migration 2: Create indexes for print_jobs
echo "2️⃣  Creating indexes for print_jobs..."
docker exec $MONGO_CONTAINER mongosh $MONGO_DB \
    -u "$MONGO_USER" -p "$MONGO_PASS" --authenticationDatabase admin \
    --quiet --eval '
db.print_jobs.createIndex({ "status": 1, "created_at": -1 });
db.print_jobs.createIndex({ "order_id": 1 });
db.print_jobs.createIndex({ "printer_id": 1 });
print("   ✅ Created print_jobs indexes");
'

echo ""

# Migration 3: Create indexes for printer_configs
echo "3️⃣  Creating indexes for printer_configs..."
docker exec $MONGO_CONTAINER mongosh $MONGO_DB \
    -u "$MONGO_USER" -p "$MONGO_PASS" --authenticationDatabase admin \
    --quiet --eval '
db.printer_configs.createIndex({ "type": 1, "is_default": 1 });
print("   ✅ Created printer_configs indexes");
'

echo ""

# Migration 4: Create indexes for print_templates
echo "4️⃣  Creating indexes for print_templates..."
docker exec $MONGO_CONTAINER mongosh $MONGO_DB \
    -u "$MONGO_USER" -p "$MONGO_PASS" --authenticationDatabase admin \
    --quiet --eval '
db.print_templates.createIndex({ "type": 1, "is_default": 1 });
print("   ✅ Created print_templates indexes");
'

echo ""

# Migration 5: Verify existing data
echo "5️⃣  Verifying existing data..."
docker exec $MONGO_CONTAINER mongosh $MONGO_DB \
    -u "$MONGO_USER" -p "$MONGO_PASS" --authenticationDatabase admin \
    --quiet --eval '
var orders = db.orders.countDocuments();
var users = db.users.countDocuments();
var menu = db.menu_items.countDocuments();
print("   ✅ orders: " + orders + " documents");
print("   ✅ users: " + users + " documents");
print("   ✅ menu_items: " + menu + " documents");
'

echo ""
echo "=========================================="
echo "✅ Migration completed successfully!"
echo "=========================================="
echo ""
echo "📝 Next steps:"
echo "   1. Restart backend: docker-compose restart backend"
echo "   2. Check logs: docker-compose logs -f backend"
echo "   3. Test frontend: https://tacafe.store"
echo ""
