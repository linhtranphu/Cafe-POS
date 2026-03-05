#!/bin/bash

# Quick script to add ONLY missing items - no verification step
# SAFE: Does NOT delete or modify existing data

set -e

echo "=========================================="
echo "➕ Quick Add Missing Items"
echo "=========================================="
echo ""

# MongoDB credentials
read -p "MongoDB username [admin]: " MONGO_USER
MONGO_USER=${MONGO_USER:-admin}

read -sp "MongoDB password: " MONGO_PASS
echo ""

read -p "Database name [cafe_pos]: " MONGO_DB
MONGO_DB=${MONGO_DB:-cafe_pos}

echo ""

# Confirm
read -p "⚠️  Add missing fields/collections? (yes/NO) " -r
echo ""
if [[ ! $REPLY =~ ^[Yy][Ee][Ss]$ ]]; then
    echo "❌ Cancelled"
    exit 1
fi

# Check container
if ! docker ps | grep -q mongodb; then
    echo "❌ MongoDB container is not running"
    exit 1
fi

MONGO_CONTAINER=$(docker ps --filter "name=mongodb" --format "{{.Names}}" | head -1)
echo "🐳 Container: $MONGO_CONTAINER"
echo ""

# Single command to do everything
echo "➕ Adding missing items..."
docker exec $MONGO_CONTAINER mongosh $MONGO_DB \
    -u "$MONGO_USER" -p "$MONGO_PASS" --authenticationDatabase admin \
    --quiet --eval '
// 1. Fix shop_settings
var existing = db.shop_settings.findOne();
if (!existing) {
    db.shop_settings.insertOne({
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
        created_at: new Date(),
        updated_at: new Date()
    });
    print("✅ Created shop_settings");
} else {
    var updates = {};
    if (!existing.print_bridge_url) updates.print_bridge_url = "http://localhost:3001";
    if (existing.auto_print_enabled === undefined) updates.auto_print_enabled = true;
    if (existing.show_logo === undefined) updates.show_logo = false;
    if (existing.show_address === undefined) updates.show_address = true;
    if (existing.show_phone === undefined) updates.show_phone = true;
    if (existing.show_custom_message === undefined) updates.show_custom_message = true;
    if (existing.low_margin_threshold === undefined) updates.low_margin_threshold = 20.0;
    
    if (Object.keys(updates).length > 0) {
        updates.updated_at = new Date();
        db.shop_settings.updateOne({ _id: existing._id }, { $set: updates });
        print("✅ Added " + Object.keys(updates).length + " fields to shop_settings");
    } else {
        print("✅ shop_settings OK");
    }
}

// 2. Create missing collections
var required = ["users","menu_items","menu_categories","orders","order_items","ingredients","ingredient_categories","batch_definitions","batch_records","batch_usage_logs","expenses","expense_categories","printer_configs","print_templates","print_jobs","print_notifications","shop_settings","shifts","cashier_shifts","cash_handovers","cash_discrepancies","fund_transactions","fund_handovers","stock_history"];
var existing = db.getCollectionNames();
var created = 0;
required.forEach(function(c) {
    if (existing.indexOf(c) === -1) {
        db.createCollection(c);
        created++;
    }
});
print("✅ Collections: " + (created > 0 ? "created " + created : "all exist"));

// 3. Add indexes (ignore errors if exist)
try { db.print_jobs.createIndex({"status":1,"created_at":-1}); } catch(e) {}
try { db.print_jobs.createIndex({"order_id":1}); } catch(e) {}
try { db.print_jobs.createIndex({"printer_id":1}); } catch(e) {}
try { db.printer_configs.createIndex({"type":1,"is_default":1}); } catch(e) {}
try { db.print_templates.createIndex({"type":1,"is_default":1}); } catch(e) {}
try { db.orders.createIndex({"created_at":-1}); } catch(e) {}
try { db.orders.createIndex({"status":1}); } catch(e) {}
print("✅ Indexes added");

print("\n✅ Done!");
'

echo ""
echo "=========================================="
echo "✅ Completed!"
echo "=========================================="
echo ""
echo "Next: docker-compose restart backend"
echo ""
