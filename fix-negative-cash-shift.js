// Fix shift with negative cash due to transfer handover bug
// This script finds shifts with negative remaining_cash and fixes them

db = db.getSiblingDB('cafe_pos');

print("=== FINDING SHIFTS WITH NEGATIVE CASH ===\n");

const shiftsWithNegativeCash = db.shifts.find({
  remaining_cash: { $lt: 0 },
  role: "WAITER"
}).toArray();

print("Found " + shiftsWithNegativeCash.length + " shifts with negative cash\n");

shiftsWithNegativeCash.forEach(shift => {
  print("\n=== Shift " + shift._id + " ===");
  print("User: " + shift.user_name);
  print("Status: " + shift.status);
  print("BEFORE FIX:");
  print("  current_cash: " + (shift.current_cash || 0));
  print("  remaining_cash: " + (shift.remaining_cash || 0));
  print("  handed_over_cash: " + (shift.handed_over_cash || 0));
  print("  transfer_revenue: " + (shift.transfer_revenue || 0));
  print("  remaining_transfer: " + (shift.remaining_transfer || 0));
  print("  handed_over_transfer: " + (shift.handed_over_transfer || 0));
  
  // Find all confirmed handovers for this shift
  const handovers = db.cash_handovers.find({
    waiter_shift_id: shift._id,
    status: "CONFIRMED"
  }).toArray();
  
  print("\nFound " + handovers.length + " confirmed handovers");
  
  // Recalculate correct amounts
  let totalCashHandedOver = 0;
  let totalTransferHandedOver = 0;
  
  handovers.forEach(h => {
    print("\nHandover " + h._id + ":");
    print("  cash_declared: " + (h.cash_declared_amount || 0));
    print("  transfer_declared: " + (h.transfer_declared_amount || 0));
    print("  declared (old): " + (h.declared_amount || 0));
    print("  cash_actual: " + (h.cash_actual_amount || 0));
    print("  transfer_actual: " + (h.transfer_actual_amount || 0));
    print("  actual (old): " + (h.actual_amount || 0));
    
    // Determine if this was cash or transfer handover
    if (h.cash_declared_amount > 0 || h.transfer_declared_amount > 0) {
      // New format
      totalCashHandedOver += (h.cash_actual_amount || 0);
      totalTransferHandedOver += (h.transfer_actual_amount || 0);
      print("  -> Type: NEW FORMAT");
      print("     Cash: " + (h.cash_actual_amount || 0));
      print("     Transfer: " + (h.transfer_actual_amount || 0));
    } else if (h.declared_amount > 0) {
      // Old format - need to guess based on shift data
      // If transfer_revenue exists and declared_amount matches, it's likely transfer
      if (shift.transfer_revenue > 0 && h.declared_amount <= shift.transfer_revenue) {
        totalTransferHandedOver += (h.actual_amount || 0);
        print("  -> Type: OLD FORMAT (guessed as TRANSFER)");
      } else {
        totalCashHandedOver += (h.actual_amount || 0);
        print("  -> Type: OLD FORMAT (guessed as CASH)");
      }
    }
  });
  
  print("\n=== CALCULATED TOTALS ===");
  print("Total cash handed over: " + totalCashHandedOver);
  print("Total transfer handed over: " + totalTransferHandedOver);
  
  // Calculate correct remaining amounts
  const correctRemainingCash = shift.current_cash - totalCashHandedOver;
  const correctRemainingTransfer = shift.transfer_revenue - totalTransferHandedOver;
  
  print("\n=== CORRECT VALUES ===");
  print("Correct remaining_cash: " + correctRemainingCash);
  print("Correct remaining_transfer: " + correctRemainingTransfer);
  
  // Update shift
  const updateResult = db.shifts.updateOne(
    { _id: shift._id },
    {
      $set: {
        remaining_cash: correctRemainingCash,
        remaining_transfer: correctRemainingTransfer,
        handed_over_cash: totalCashHandedOver,
        handed_over_transfer: totalTransferHandedOver,
        updated_at: new Date()
      }
    }
  );
  
  print("\nUpdate result: " + (updateResult.modifiedCount > 0 ? "SUCCESS" : "FAILED"));
  
  // Verify
  const updatedShift = db.shifts.findOne({ _id: shift._id });
  print("\nAFTER FIX:");
  print("  current_cash: " + (updatedShift.current_cash || 0));
  print("  remaining_cash: " + (updatedShift.remaining_cash || 0));
  print("  handed_over_cash: " + (updatedShift.handed_over_cash || 0));
  print("  transfer_revenue: " + (updatedShift.transfer_revenue || 0));
  print("  remaining_transfer: " + (updatedShift.remaining_transfer || 0));
  print("  handed_over_transfer: " + (updatedShift.handed_over_transfer || 0));
});

print("\n=== FIX COMPLETE ===");
