// Check recent shifts
db = db.getSiblingDB('cafe_pos');

print("=== ALL WAITER SHIFTS (RECENT 5) ===\n");

db.shifts.find({ role: "WAITER" })
  .sort({ created_at: -1 })
  .limit(5)
  .forEach(shift => {
    print("\n=== Shift " + shift._id + " ===");
    print("User: " + shift.user_name);
    print("Status: " + shift.status);
    print("Created: " + shift.created_at);
    print("Start Cash: " + (shift.start_cash || 0));
    print("Current Cash: " + (shift.current_cash || 0));
    print("Remaining Cash: " + (shift.remaining_cash || 0));
    print("Handed Over Cash: " + (shift.handed_over_cash || 0));
    print("Transfer Revenue: " + (shift.transfer_revenue || 0));
    print("Remaining Transfer: " + (shift.remaining_transfer || 0));
    print("Handed Over Transfer: " + (shift.handed_over_transfer || 0));
    print("Total Revenue: " + (shift.total_revenue || 0));
    print("Handover Count: " + (shift.handover_count || 0));
    
    // Check handovers for this shift
    const handoverCount = db.cash_handovers.countDocuments({
      waiter_shift_id: shift._id
    });
    print("Actual Handovers in DB: " + handoverCount);
  });
