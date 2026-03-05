#!/bin/bash

# Script to compare database structure between Production and Local
# This helps identify differences in collections, indexes, and schema

set -e

echo "=========================================="
echo "🔍 Database Structure Comparison"
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

echo "📊 Analyzing database structure..."
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

# Get all collections
echo "📋 Collections in database:"
echo "-------------------------------------------"
docker exec $MONGO_CONTAINER mongosh $MONGO_DB \
    -u "$MONGO_USER" -p "$MONGO_PASS" --authenticationDatabase admin \
    --quiet --eval '
db.getCollectionNames().forEach(function(collection) {
    var count = db[collection].countDocuments();
    print(collection + ": " + count + " documents");
});
'

echo ""
echo "-------------------------------------------"
echo ""

# Get indexes for each collection
echo "🔑 Indexes for each collection:"
echo "-------------------------------------------"
docker exec $MONGO_CONTAINER mongosh $MONGO_DB \
    -u "$MONGO_USER" -p "$MONGO_PASS" --authenticationDatabase admin \
    --quiet --eval '
db.getCollectionNames().forEach(function(collection) {
    print("\n" + collection + ":");
    var indexes = db[collection].getIndexes();
    indexes.forEach(function(index) {
        print("  - " + index.name + ": " + JSON.stringify(index.key));
    });
});
'

echo ""
echo "-------------------------------------------"
echo ""

# Sample document structure from each collection
echo "📄 Sample document structure (first document from each collection):"
echo "-------------------------------------------"
docker exec $MONGO_CONTAINER mongosh $MONGO_DB \
    -u "$MONGO_USER" -p "$MONGO_PASS" --authenticationDatabase admin \
    --quiet --eval '
db.getCollectionNames().forEach(function(collection) {
    var sample = db[collection].findOne();
    if (sample) {
        print("\n" + collection + ":");
        print("  Fields: " + Object.keys(sample).join(", "));
    }
});
'

echo ""
echo "-------------------------------------------"
echo ""

# Check for shop_settings specifically
echo "⚙️  Shop Settings Details:"
echo "-------------------------------------------"
docker exec $MONGO_CONTAINER mongosh $MONGO_DB \
    -u "$MONGO_USER" -p "$MONGO_PASS" --authenticationDatabase admin \
    --quiet --eval '
var settings = db.shop_settings.findOne();
if (settings) {
    print("Shop settings exists:");
    print("  Fields: " + Object.keys(settings).join(", "));
    print("\n  Has print_bridge_url: " + (settings.print_bridge_url !== undefined));
    print("  Has auto_print_enabled: " + (settings.auto_print_enabled !== undefined));
} else {
    print("❌ Shop settings collection does not exist or is empty");
}
'

echo ""
echo "-------------------------------------------"
echo ""

echo "✅ Analysis complete!"
echo ""
echo "💡 To export this structure to a file:"
echo "   sudo bash compare-db-structure.sh > db-structure.txt"
echo ""
