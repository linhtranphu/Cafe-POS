// Find shift opened at 11:40
db = db.getSiblingDB('cafe_pos');

print("=== FINDING SHIFT OPENED AT 11:40 ===\n");

// Get today's date
const today = new Date();
today.setHours(0, 0, 0, 0);

const tomorrow = new Date(today);
tomorrow.setDate(tomorrow.getDate() + 1);

// Find shifts created today
const shiftsToday = db.shifts.find({
  role: "WAITER",
  created_at: {
    $gte: today,
    $lt: tomorrow
  }
}).sort({ created_at: -1 }).toArray();

print("Found " + shiftsToday.length + " waiter shifts today\n");

shiftsToday.forEach(shift => {
  const createdTime = new Date(shift.created_at);
  const hours = createdTime.getHours();
  const minutes = createdTime.getMinutes();
  const timeStr = hours + ":" + (minutes < 10 ? "0" + minutes : minutes);
  
  print("\n=== Shift " + shift._id + " ===");
  print("User: " + shift.user_name);
  print("Status: " + shift.status);
  print("Created at: " + createdTime + " (" + timeStr + ")");
  print("Start Cash: " + (shift.start_cash || 0));
  print("Current Cash: " + (shift.current_cash || 0));
  print("Remaining Cash: " + (shift.remaining_cash || 0));
  print("Handed Over Cash: " + (shift.handed_over_cash || 0));
  print("Transfer Revenue: " + (shift.transfer_revenue || 0));
  print("Remaining Transfer: " + (shift.remaining_transfer || 0));
  print("Handed Over Transfer: " + (shift.handed_over_transfer || 0));
  
  // Check if this is the 11:40 shift
  if (hours === 11 && minutes === 40) {
    print("\n*** THIS IS THE 11:40 SHIFT ***");
    
    // Get all handovers for this shift
    print("\n=== HANDOVERS FOR THIS SHIFT ===");
    const handovers = db.cash_handovers.find({
      waiter_shift_id: shift._id
    }).sort({ created_at: 1 }).toArray();
    
    print("Total handovers: " + handovers.length + "\n");
    
    handovers.forEach((h, i) => {
      print("Handover #" + (i + 1) + " (" + h._id + "):");
      print("  Status: " + h.status);
      print("  Type: " + h.handover_type);
      print("  Cash Declared: " + (h.cash_declared_amount || 0));
      print("  Transfer Declared: " + (h.transfer_declared_amount || 0));
      print("  Declared (old): " + (h.declared_amount || 0));
      print("  Cash Actual: " + (h.cash_actual_amount || 0));
      print("  Transfer Actual: " + (h.transfer_actual_amount || 0));
      print("  Actual (old): " + (h.actual_amount || 0));
      print("  Created: " + h.created_at);
      if (h.confirmed_at) {
        print("  Confirmed: " + h.confirmed_at);
      }
      print("");
    });
  }
});

// Also check for OPEN shifts regardless of time
print("\n=== ALL OPEN WAITER SHIFTS ===\n");
const openShifts = db.shifts.find({
  role: "WAITER",
  status: "OPEN"
}).toArray();

print("Found " + openShifts.length + " open waiter shifts\n");

openShifts.forEach(shift => {
  const createdTime = new Date(shift.created_at);
  const hours = createdTime.getHours();
  const minutes = createdTime.getMinutes();
  const timeStr = hours + ":" + (minutes < 10 ? "0" + minutes : minutes);
  
  print("\n=== OPEN Shift " + shift._id + " ===");
  print("User: " + shift.user_name);
  print("Created at: " + timeStr);
  print("Start Cash: " + (shift.start_cash || 0));
  print("Current Cash: " + (shift.current_cash || 0));
  print("Remaining Cash: " + (shift.remaining_cash || 0));
  print("Handed Over Cash: " + (shift.handed_over_cash || 0));
  print("Transfer Revenue: " + (shift.transfer_revenue || 0));
  print("Remaining Transfer: " + (shift.remaining_transfer || 0));
  print("Handed Over Transfer: " + (shift.handed_over_transfer || 0));
});
