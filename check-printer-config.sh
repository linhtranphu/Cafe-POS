#!/bin/bash

# Script to check printer configuration in MongoDB

echo "=== Checking Printer Configuration ==="
echo ""

# Check if mongo or mongosh is available
if command -v mongosh &> /dev/null; then
    MONGO_CMD="mongosh"
elif command -v mongo &> /dev/null; then
    MONGO_CMD="mongo"
else
    echo "❌ MongoDB client not found. Please install mongo or mongosh."
    exit 1
fi

MONGO_URI="mongodb://admin:password123@localhost:27017/cafe_pos?replicaSet=rs0&authSource=admin"

echo "1. Checking all printers..."
$MONGO_CMD "$MONGO_URI" --quiet --eval "
    db.printer_configs.find({}, {
        name: 1, 
        type: 1, 
        is_default: 1, 
        connection_type: 1,
        ip_address: 1,
        port: 1
    }).forEach(function(p) {
        print('---');
        print('Name: ' + p.name);
        print('Type: ' + p.type);
        print('Is Default: ' + p.is_default);
        print('Connection: ' + p.connection_type);
        print('IP: ' + p.ip_address + ':' + p.port);
    });
"

echo ""
echo "2. Checking default BILL printer..."
$MONGO_CMD "$MONGO_URI" --quiet --eval "
    var billPrinter = db.printer_configs.findOne({type: 'BILL', is_default: true});
    if (billPrinter) {
        print('✅ Found: ' + billPrinter.name);
    } else {
        print('❌ No default BILL printer found');
    }
"

echo ""
echo "3. Checking default LABEL printer..."
$MONGO_CMD "$MONGO_URI" --quiet --eval "
    var labelPrinter = db.printer_configs.findOne({type: 'LABEL', is_default: true});
    if (labelPrinter) {
        print('✅ Found: ' + labelPrinter.name);
    } else {
        print('❌ No default LABEL printer found');
    }
"

echo ""
echo "4. Checking templates..."
$MONGO_CMD "$MONGO_URI" --quiet --eval "
    var billTemplate = db.print_templates.findOne({type: 'BILL', is_default: true});
    var labelTemplate = db.print_templates.findOne({type: 'LABEL', is_default: true});
    
    if (billTemplate) {
        print('✅ Bill template: ' + billTemplate.name);
    } else {
        print('❌ No default BILL template');
    }
    
    if (labelTemplate) {
        print('✅ Label template: ' + labelTemplate.name);
    } else {
        print('❌ No default LABEL template');
    }
"

echo ""
echo "5. Checking auto-print setting..."
$MONGO_CMD "$MONGO_URI" --quiet --eval "
    var settings = db.shop_settings.findOne({});
    if (settings) {
        print('Auto Print Enabled: ' + settings.auto_print_enabled);
    } else {
        print('⚠️  No shop settings found');
    }
"

echo ""
echo "=== Check Complete ==="
