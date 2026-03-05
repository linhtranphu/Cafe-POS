// MongoDB Initialization Script
// This script runs automatically when MongoDB container starts for the first time
// It creates the database structure, collections, indexes, and default data

print('🚀 Initializing Cafe POS Database...');

// Switch to cafe_pos database
db = db.getSiblingDB('cafe_pos');

print('📋 Creating collections...');

// Create collections
db.createCollection('users');
db.createCollection('menu_items');
db.createCollection('menu_categories');
db.createCollection('orders');
db.createCollection('order_items');
db.createCollection('ingredients');
db.createCollection('ingredient_categories');
db.createCollection('batch_definitions');
db.createCollection('batch_records');
db.createCollection('batch_usage_logs');
db.createCollection('expenses');
db.createCollection('expense_categories');
db.createCollection('printer_configs');
db.createCollection('print_templates');
db.createCollection('print_jobs');
db.createCollection('print_notifications');
db.createCollection('shop_settings');
db.createCollection('shifts');
db.createCollection('cashier_shifts');
db.createCollection('cash_handovers');
db.createCollection('cash_discrepancies');
db.createCollection('fund_transactions');
db.createCollection('fund_handovers');
db.createCollection('stock_history');

print('✅ Collections created');

print('🔑 Creating indexes...');

// Indexes for orders
db.orders.createIndex({ "created_at": -1 });
db.orders.createIndex({ "status": 1 });
db.orders.createIndex({ "order_number": 1 }, { unique: true });

// Indexes for print_jobs
db.print_jobs.createIndex({ "status": 1, "created_at": -1 });
db.print_jobs.createIndex({ "order_id": 1 });
db.print_jobs.createIndex({ "printer_id": 1 });

// Indexes for printer_configs
db.printer_configs.createIndex({ "type": 1, "is_default": 1 });

// Indexes for print_templates
db.print_templates.createIndex({ "type": 1, "is_default": 1 });

// Indexes for users
db.users.createIndex({ "username": 1 }, { unique: true });
db.users.createIndex({ "role": 1 });

// Indexes for menu_items
db.menu_items.createIndex({ "name": 1 });
db.menu_items.createIndex({ "category_id": 1 });

// Indexes for ingredients
db.ingredients.createIndex({ "name": 1 });

// Indexes for batches
db.batch_records.createIndex({ "ingredient_id": 1 });
db.batch_records.createIndex({ "created_at": -1 });

// Indexes for expenses
db.expenses.createIndex({ "date": -1 });
db.expenses.createIndex({ "category_id": 1 });

print('✅ Indexes created');

print('⚙️  Creating shop settings...');

// Create default shop settings
var now = new Date();
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
    created_at: now,
    updated_at: now
});

print('✅ Shop settings created');

print('✅ Database initialization completed!');
print('📊 Collections: ' + db.getCollectionNames().length);
