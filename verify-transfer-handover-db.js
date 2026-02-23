// Verify Transfer Handover in Database
// This script checks the database to ensure transfer handovers are recorded correctly

db = db.getSiblingDB('cafe_pos');

print("=================================");
print("🔍 VERIFY: Transfer Handover Data");
print("=================================\n");

// Find the most recent transfer handover
const recentHandover = db.cash_handovers.findOne(
  {
    transfer_declared_amount: {$gt: 0},
    cash_declared_amount: 0
  },
  {sort: {created_at: -1}}
);

if (!recentHandover) {
  print("❌ No transfer-only handovers found");
  quit(1);
}

print("📋 Most Recent Transfer Handover:");
print("   ID: " + recentHandover._id);
print("   Status: " + recentHandover.status);
print("   Type: " + recentHandover.handover_type);
print("");

print("💰 Declared Amounts:");
print("   Cash: " + (recentHandover.cash_declared_amount || 0) + " VND");
print("   Transfer: " + (recentHandover.transfer_declared_amount || 0) + " VND");
print("   Total (old): " + (recentHandover.declared_amount || 0) + " VND");
print("");

print("✅ Actual Amounts:");
print("   Cash: " + (recentHandover.cash_actual_amount || 0) + " VND");
print("   Transfer: " + (recentHandover.transfer_actual_amount || 0) + " VND");
print("   Total (old): " + (recentHandover.actual_amount || 0) + " VND");
print("");

// Verify handover data
let handoverPass = true;

// Check 1: Cash declared should be 0
if (recentHandover.cash_declared_amount !== 0) {
  print("❌ FAIL: Cash declared should be 0");
  handoverPass = false;
} else {
  print("✅ PASS: Cash declared is 0");
}

// Check 2: Transfer declared should be > 0
if (recentHandover.transfer_declared_amount <= 0) {
  print("❌ FAIL: Transfer declared should be > 0");
  handoverPass = false;
} else {
  print("✅ PASS: Transfer declared is " + recentHandover.transfer_declared_amount);
}

// Check 3: If confirmed, cash actual should be 0
if (recentHandover.status === "CONFIRMED") {
  if (recentHandover.cash_actual_amount !== 0) {
    print("❌ FAIL: Cash actual should be 0");
    handoverPass = false;
  } else {
    print("✅ PASS: Cash actual is 0");
  }
  
  // Check 4: Transfer actual should match declared
  if (recentHandover.transfer_actual_amount !== recentHandover.transfer_declared_amount) {
    print("❌ FAIL: Transfer actual (" + recentHandover.transfer_actual_amount + 
          ") should match declared (" + recentHandover.transfer_declared_amount + ")");
    handoverPass = false;
  } else {
    print("✅ PASS: Transfer actual matches declared");
  }
}

print("");

// Check the corresponding shift
print("📊 Checking Shift Data:");
const shift = db.shifts.findOne({_id: recentHandover.waiter_shift_id});

if (!shift) {
  print("❌ Shift not found!");
  quit(1);
}

print("   Shift ID: " + shift._id);
print("   User: " + shift.user_name);
print("   Status: " + shift.status);
print("");

print("💵 Cash Amounts:");
print("   Current Cash: " + (shift.current_cash || 0) + " VND");
print("   Remaining Cash: " + (shift.remaining_cash || 0) + " VND");
print("   Handed Over Cash: " + (shift.handed_over_cash || 0) + " VND");
print("");

print("💳 Transfer Amounts:");
print("   Transfer Revenue: " + (shift.transfer_revenue || 0) + " VND");
print("   Remaining Transfer: " + (shift.remaining_transfer || 0) + " VND");
print("   Handed Over Transfer: " + (shift.handed_over_transfer || 0) + " VND");
print("");

// Verify shift data
let shiftPass = true;

// Check 5: Handed over cash should be 0 (no cash handover)
if (shift.handed_over_cash > 0) {
  // Check if there are any cash handovers
  const cashHandovers = db.cash_handovers.countDocuments({
    waiter_shift_id: shift._id,
    cash_declared_amount: {$gt: 0},
    status: "CONFIRMED"
  });
  
  if (cashHandovers === 0) {
    print("❌ FAIL: Handed over cash is " + shift.handed_over_cash + " but no cash handovers found");
    shiftPass = false;
  } else {
    print("✅ PASS: Handed over cash is from actual cash handovers");
  }
} else {
  print("✅ PASS: No cash handed over (0 VND)");
}

// Check 6: Handed over transfer should match handover amount
if (recentHandover.status === "CONFIRMED") {
  const expectedHandedOverTransfer = recentHandover.transfer_actual_amount || 0;
  
  // Get all confirmed transfer handovers for this shift
  const allTransferHandovers = db.cash_handovers.aggregate([
    {
      $match: {
        waiter_shift_id: shift._id,
        status: "CONFIRMED"
      }
    },
    {
      $group: {
        _id: null,
        totalTransfer: {$sum: "$transfer_actual_amount"}
      }
    }
  ]).toArray();
  
  const totalTransferHandedOver = allTransferHandovers.length > 0 ? allTransferHandovers[0].totalTransfer : 0;
  
  if (shift.handed_over_transfer !== totalTransferHandedOver) {
    print("❌ FAIL: Handed over transfer (" + shift.handed_over_transfer + 
          ") doesn't match sum of handovers (" + totalTransferHandedOver + ")");
    shiftPass = false;
  } else {
    print("✅ PASS: Handed over transfer matches handovers (" + totalTransferHandedOver + " VND)");
  }
}

// Check 7: Remaining transfer should be transfer_revenue - handed_over_transfer
const expectedRemainingTransfer = (shift.transfer_revenue || 0) - (shift.handed_over_transfer || 0);
if (shift.remaining_transfer !== expectedRemainingTransfer) {
  print("❌ FAIL: Remaining transfer (" + shift.remaining_transfer + 
        ") should be " + expectedRemainingTransfer);
  shiftPass = false;
} else {
  print("✅ PASS: Remaining transfer is correct (" + shift.remaining_transfer + " VND)");
}

// Check 8: Current cash should not be affected by transfer handover
// (This is harder to verify without knowing the initial state, but we can check it's not negative)
if (shift.current_cash < 0) {
  print("❌ FAIL: Current cash is negative (" + shift.current_cash + ")");
  shiftPass = false;
} else {
  print("✅ PASS: Current cash is not negative (" + shift.current_cash + " VND)");
}

// Check 9: Remaining cash should not be negative
if (shift.remaining_cash < 0) {
  print("❌ FAIL: Remaining cash is negative (" + shift.remaining_cash + ")");
  shiftPass = false;
} else {
  print("✅ PASS: Remaining cash is not negative (" + shift.remaining_cash + " VND)");
}

print("");
print("=================================");

if (handoverPass && shiftPass) {
  print("✅ ALL CHECKS PASSED!");
  print("=================================");
  quit(0);
} else {
  print("❌ SOME CHECKS FAILED!");
  print("=================================");
  quit(1);
}
