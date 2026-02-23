// Check handover and shift data to diagnose the issue
db = db.getSiblingDB('cafe_pos');

print("=== CHECKING HANDOVER DATA ===\n");

// Find the most recent handover
const recentHandover = db.cash_handovers.findOne(
  {},
  { sort: { created_at: -1 } }
);

if (recentHandover) {
  print("Most Recent Handover:");
  print("ID: " + recentHandover._id);
  print("Waiter Shift ID: " + recentHandover.waiter_shift_id);
  print("Status: " + recentHandover.status);
  print("\nDeclared Amounts:");
  print("  cash_declared_amount: " + (recentHandover.cash_declared_amount || 0));
  print("  transfer_declared_amount: " + (recentHandover.transfer_declared_amount || 0));
  print("  declared_amount (deprecated): " + (recentHandover.declared_amount || 0));
  print("\nActual Amounts:");
  print("  cash_actual_amount: " + (recentHandover.cash_actual_amount || 0));
  print("  transfer_actual_amount: " + (recentHandover.transfer_actual_amount || 0));
  print("  actual_amount (deprecated): " + (recentHandover.actual_amount || 0));
  print("\nDiscrepancies:");
  print("  cash_discrepancy: " + (recentHandover.cash_discrepancy || 0));
  print("  transfer_discrepancy: " + (recentHandover.transfer_discrepancy || 0));
  print("  discrepancy (deprecated): " + (recentHandover.discrepancy || 0));
  
  // Check the corresponding shift
  print("\n=== CHECKING SHIFT DATA ===\n");
  const shift = db.shifts.findOne({ _id: recentHandover.waiter_shift_id });
  
  if (shift) {
    print("Shift ID: " + shift._id);
    print("User: " + shift.user_name);
    print("Status: " + shift.status);
    print("\nCash Amounts:");
    print("  start_cash: " + (shift.start_cash || 0));
    print("  current_cash: " + (shift.current_cash || 0));
    print("  remaining_cash: " + (shift.remaining_cash || 0));
    print("  handed_over_cash: " + (shift.handed_over_cash || 0));
    print("\nTransfer Amounts:");
    print("  transfer_revenue: " + (shift.transfer_revenue || 0));
    print("  remaining_transfer: " + (shift.remaining_transfer || 0));
    print("  handed_over_transfer: " + (shift.handed_over_transfer || 0));
    print("\nTotals:");
    print("  total_revenue: " + (shift.total_revenue || 0));
    print("  handover_count: " + (shift.handover_count || 0));
  } else {
    print("Shift not found!");
  }
} else {
  print("No handovers found!");
}

print("\n=== ALL CONFIRMED HANDOVERS FOR THIS SHIFT ===\n");
if (recentHandover) {
  const allHandovers = db.cash_handovers.find({
    waiter_shift_id: recentHandover.waiter_shift_id,
    status: "CONFIRMED"
  }).toArray();
  
  print("Total confirmed handovers: " + allHandovers.length);
  allHandovers.forEach((h, i) => {
    print("\nHandover #" + (i + 1) + ":");
    print("  ID: " + h._id);
    print("  Cash declared: " + (h.cash_declared_amount || 0));
    print("  Transfer declared: " + (h.transfer_declared_amount || 0));
    print("  Cash actual: " + (h.cash_actual_amount || 0));
    print("  Transfer actual: " + (h.transfer_actual_amount || 0));
    print("  Confirmed at: " + h.confirmed_at);
  });
}
