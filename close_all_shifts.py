#!/usr/bin/env python3
"""
Script to close all open waiter shifts directly in MongoDB
"""

from pymongo import MongoClient
from datetime import datetime
from bson import ObjectId

# MongoDB connection
MONGO_URI = "mongodb://admin:password123@localhost:27017/cafe_pos?replicaSet=rs0&authSource=admin"
DB_NAME = "cafe_pos"

def main():
    print("🔄 Closing All Open Waiter Shifts")
    print("==================================\n")
    
    try:
        # Connect to MongoDB
        client = MongoClient(MONGO_URI)
        db = client[DB_NAME]
        
        # Find all open shifts
        open_shifts = list(db.shifts.find({"status": "OPEN"}))
        
        if not open_shifts:
            print("✅ No open shifts found. All shifts are already closed.\n")
            return
        
        print(f"Found {len(open_shifts)} open shift(s)\n")
        
        closed_count = 0
        now = datetime.now()
        
        # Close each shift
        for i, shift in enumerate(open_shifts):
            print(f"[{i+1}/{len(open_shifts)}] Closing shift: {shift['_id']}")
            print(f"   User: {shift.get('user_name', 'Unknown')} ({shift.get('role_type', 'Unknown')})")
            print(f"   Type: {shift.get('type', 'Unknown')}")
            print(f"   Started: {shift.get('started_at', 'Unknown')}")
            
            # Get orders for this shift
            orders = list(db.orders.find({"shift_id": shift["_id"]}))
            
            # Calculate revenue
            total_revenue = 0
            cash_revenue = 0
            transfer_revenue = 0
            total_orders = 0
            
            for order in orders:
                if order.get("status") in ["PAID", "IN_PROGRESS", "SERVED"]:
                    total_revenue += order.get("total", 0)
                    total_orders += 1
                    
                    payment_method = order.get("payment_method")
                    if payment_method == "CASH":
                        cash_revenue += order.get("total", 0)
                    elif payment_method in ["TRANSFER", "QR"]:
                        transfer_revenue += order.get("total", 0)
            
            # Update shift
            result = db.shifts.update_one(
                {"_id": shift["_id"]},
                {
                    "$set": {
                        "status": "CLOSED",
                        "ended_at": now,
                        "total_revenue": total_revenue,
                        "total_orders": total_orders,
                        "end_cash": shift.get("current_cash", 0),
                        "updated_at": now
                    }
                }
            )
            
            if result.modified_count > 0:
                # Lock completed orders
                lock_result = db.orders.update_many(
                    {
                        "shift_id": shift["_id"],
                        "status": {"$in": ["SERVED", "CANCELLED"]}
                    },
                    {
                        "$set": {
                            "status": "LOCKED",
                            "locked_at": now
                        }
                    }
                )
                
                print(f"   ✅ Shift closed successfully")
                print(f"   📊 Summary:")
                print(f"      - Total Orders: {total_orders}")
                print(f"      - Total Revenue: {total_revenue:.0f} VND")
                print(f"      - Cash Revenue: {cash_revenue:.0f} VND")
                print(f"      - Transfer Revenue: {transfer_revenue:.0f} VND")
                print(f"      - Current Cash: {shift.get('current_cash', 0):.0f} VND")
                print(f"      - Remaining Cash: {shift.get('remaining_cash', 0):.0f} VND")
                print(f"      - Remaining Transfer: {shift.get('remaining_transfer', 0):.0f} VND")
                print(f"      - Locked Orders: {lock_result.modified_count}")
                print()
                
                closed_count += 1
            else:
                print(f"   ❌ Failed to close shift")
                print()
        
        print(f"✅ Successfully closed {closed_count}/{len(open_shifts)} shift(s)\n")
        
    except Exception as e:
        print(f"❌ Error: {e}")
        return 1
    
    return 0

if __name__ == "__main__":
    exit(main())
