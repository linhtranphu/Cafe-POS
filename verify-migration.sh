#!/bin/bash

# Script to verify migration was successful
# Can run on both local and production

echo "=========================================="
echo "🔍 Verifying Migration Status"
echo "=========================================="
echo ""

# Check if MongoDB is accessible
if docker ps | grep -q mongodb; then
    MONGO_CONTAINER=$(docker ps --filter "name=mongodb" --format "{{.Names}}" | head -1)
    echo "✅ MongoDB container found: $MONGO_CONTAINER"
else
    echo "❌ MongoDB container not found"
    exit 1
fi

echo ""
echo "📊 Checking database structure..."
echo ""

# Check shop_settings collection
echo "1️⃣  Shop Settings:"
SHOP_SETTINGS=$(docker exec $MONGO_CONTAINER mongosh cafe_pos \
    --quiet --eval "db.shop_settings.countDocuments()")

if [ "$SHOP_SETTINGS" -gt 0 ]; then
    echo "   ✅ shop_settings exists ($SHOP_SETTINGS documents)"
    docker exec $MONGO_CONTAINER mongosh cafe_pos \
        --quiet --eval "db.shop_settings.findOne({}, {shop_name: 1, auto_print_enabled: 1})" | head -5
else
    echo "   ❌ shop_settings collection is empty"
    echo "   → Run migration: sudo bash migrate-v2.0-mongodb.sh"
fi

echo ""

# Check print collections
echo "2️⃣  Print Collections:"
for collection in print_jobs printer_configs print_templates print_notifications; do
    EXISTS=$(docker exec $MONGO_CONTAINER mongosh cafe_pos \
        --quiet --eval "db.getCollectionNames().includes('$collection')")
    
    if [ "$EXISTS" = "true" ]; then
        COUNT=$(docker exec $MONGO_CONTAINER mongosh cafe_pos \
            --quiet --eval "db.$collection.countDocuments()")
        echo "   ✅ $collection exists ($COUNT documents)"
    else
        echo "   ❌ $collection not found"
    fi
done

echo ""

# Check indexes
echo "3️⃣  Indexes:"
PRINT_JOBS_INDEXES=$(docker exec $MONGO_CONTAINER mongosh cafe_pos \
    --quiet --eval "db.print_jobs.getIndexes().length")
echo "   print_jobs: $PRINT_JOBS_INDEXES indexes"

PRINTER_CONFIGS_INDEXES=$(docker exec $MONGO_CONTAINER mongosh cafe_pos \
    --quiet --eval "db.printer_configs.getIndexes().length")
echo "   printer_configs: $PRINTER_CONFIGS_INDEXES indexes"

PRINT_TEMPLATES_INDEXES=$(docker exec $MONGO_CONTAINER mongosh cafe_pos \
    --quiet --eval "db.print_templates.getIndexes().length")
echo "   print_templates: $PRINT_TEMPLATES_INDEXES indexes"

echo ""

# Check existing data (should not be affected)
echo "4️⃣  Existing Data (should be preserved):"
ORDERS=$(docker exec $MONGO_CONTAINER mongosh cafe_pos \
    --quiet --eval "db.orders.countDocuments()")
echo "   orders: $ORDERS documents"

USERS=$(docker exec $MONGO_CONTAINER mongosh cafe_pos \
    --quiet --eval "db.users.countDocuments()")
echo "   users: $USERS documents"

MENU=$(docker exec $MONGO_CONTAINER mongosh cafe_pos \
    --quiet --eval "db.menu_items.countDocuments()")
echo "   menu_items: $MENU documents"

echo ""
echo "=========================================="

# Summary
if [ "$SHOP_SETTINGS" -gt 0 ]; then
    echo "✅ Migration Status: SUCCESS"
    echo ""
    echo "📝 Next steps:"
    echo "   1. Configure shop settings in UI"
    echo "   2. Add printers in Print Management"
    echo "   3. Test bill printing"
else
    echo "⚠️  Migration Status: INCOMPLETE"
    echo ""
    echo "📝 Action required:"
    echo "   Run: sudo bash migrate-v2.0-production.sh"
fi

echo "=========================================="
