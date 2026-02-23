// MongoDB script to close all open waiter shifts
// Run with: mongosh cafe_pos -u admin -p password123 --authenticationDatabase admin < close-all-shifts-mongo.js

print("🔄 Closing All Open Waiter Shifts");
print("==================================\n");

// Switch to cafe_pos database
db = db.getSiblingDB('cafe_pos');

// Find all open shifts
const openShifts = db.shifts.find({ status: "OPEN" }).toArray();

if (openShifts.length === 0) {
    print("✅ No open shifts found. All shifts are already closed.\n");
    quit();
}

print(`Found ${openShifts.length} open shift(s)\n`);

let closedCount = 0;
const now = new Date();

// Close each shift
openShifts.forEach((shift, index) => {
    print(`[${index + 1}/${openShifts.length}] Closing shift: ${shift._id}`);
    print(`   User: ${shift.user_name} (${shift.role_type})`);
    print(`   Type: ${shift.type}`);
    print(`   Started: ${shift.started_at}`);
    
    // Get orders for this shift
    const orders = db.orders.find({ shift_id: shift._id }).toArray();
    
    // Calculate revenue
    let totalRevenue = 0;
    let cashRevenue = 0;
    let transferRevenue = 0;
    let totalOrders = 0;
    
    orders.forEach(order => {
        if (order.status === "PAID" || order.status === "IN_PROGRESS" || order.status === "SERVED") {
            totalRevenue += order.total || 0;
            totalOrders++;
            
            if (order.payment_method === "CASH") {
                cashRevenue += order.total || 0;
            } else if (order.payment_method === "TRANSFER" || order.payment_method === "QR") {
                transferRevenue += order.total || 0;
            }
        }
    });
    
    // Update shift
    const updateResult = db.shifts.updateOne(
        { _id: shift._id },
        {
            $set: {
                status: "CLOSED",
                ended_at: now,
                total_revenue: totalRevenue,
                total_orders: totalOrders,
                end_cash: shift.current_cash || 0,
                updated_at: now
            }
        }
    );
    
    if (updateResult.modifiedCount > 0) {
        // Lock completed orders
        const lockResult = db.orders.updateMany(
            { 
                shift_id: shift._id,
                status: { $in: ["SERVED", "CANCELLED"] }
            },
            {
                $set: {
                    status: "LOCKED",
                    locked_at: now
                }
            }
        );
        
        print(`   ✅ Shift closed successfully`);
        print(`   📊 Summary:`);
        print(`      - Total Orders: ${totalOrders}`);
        print(`      - Total Revenue: ${totalRevenue.toFixed(0)} VND`);
        print(`      - Cash Revenue: ${cashRevenue.toFixed(0)} VND`);
        print(`      - Transfer Revenue: ${transferRevenue.toFixed(0)} VND`);
        print(`      - Current Cash: ${(shift.current_cash || 0).toFixed(0)} VND`);
        print(`      - Remaining Cash: ${(shift.remaining_cash || 0).toFixed(0)} VND`);
        print(`      - Remaining Transfer: ${(shift.remaining_transfer || 0).toFixed(0)} VND`);
        print(`      - Locked Orders: ${lockResult.modifiedCount}`);
        print("");
        
        closedCount++;
    } else {
        print(`   ❌ Failed to close shift`);
        print("");
    }
});

print(`✅ Successfully closed ${closedCount}/${openShifts.length} shift(s)\n`);
