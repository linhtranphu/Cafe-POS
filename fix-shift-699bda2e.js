// Fix shift 699bda2ea31585ce7ad4c47c and its handover
db = db.getSiblingDB('cafe_pos');

const shiftId = ObjectId("699bda2ea31585ce7ad4c47c");
const handoverId = ObjectId("699bda5ca31585ce7ad4c484");

print("=== FIXING HANDOVER AND SHIFT ===\n");

// Step 1: Fix handover - set transfer_actual_amount correctly
print("Step 1: Fixing handover " + handoverId);
const handover = db.cash_handovers.findOne({_id: handoverId});

print("BEFORE:");
print("  cash_declared_amount: " + (handover.cash_declared_amount || 0));
print("  transfer_declared_amount: " + (handover.transfer_declared_amount || 0));
print("  cash_actual_amount: " + (handover.cash_actual_amount || 0));
print("  transfer_actual_amount: " + (handover.transfer_actual_amount || 0));
print("  actual_amount (old): " + (handover.actual_amount || 0));

// Update handover: move actual_amount to transfer_actual_amount
db.cash_handovers.updateOne(
  {_id: handoverId},
  {
    $set: {
      transfer_actual_amount: handover.actual_amount,
      cash_actual_amount: 0,
      transfer_discrepancy: 0,
      cash_discrepancy: 0,
      updated_at: new Date()
    }
  }
);

const updatedHandover = db.cash_handovers.findOne({_id: handoverId});
print("\nAFTER:");
print("  cash_actual_amount: " + (updatedHandover.cash_actual_amount || 0));
print("  transfer_actual_amount: " + (updatedHandover.transfer_actual_amount || 0));

// Step 2: Fix shift - move handed_over_cash to handed_over_transfer
print("\n\nStep 2: Fixing shift " + shiftId);
const shift = db.shifts.findOne({_id: shiftId});

print("BEFORE:");
print("  current_cash: " + (shift.current_cash || 0));
print("  remaining_cash: " + (shift.remaining_cash || 0));
print("  handed_over_cash: " + (shift.handed_over_cash || 0));
print("  transfer_revenue: " + (shift.transfer_revenue || 0));
print("  remaining_transfer: " + (shift.remaining_transfer || 0));
print("  handed_over_transfer: " + (shift.handed_over_transfer || 0));

// Calculate correct values
// The handover was 30,000 transfer, not cash
// So we need to:
// 1. Add back 30,000 to remaining_cash (it was wrongly deducted)
// 2. Deduct 30,000 from remaining_transfer (it should have been deducted)
// 3. Move 30,000 from handed_over_cash to handed_over_transfer

const correctRemainingCash = shift.remaining_cash + 30000;  // Add back wrongly deducted
const correctRemainingTransfer = shift.remaining_transfer - 30000;  // Deduct correctly
const correctHandedOverCash = shift.handed_over_cash - 30000;  // Remove wrong amount
const correctHandedOverTransfer = (shift.handed_over_transfer || 0) + 30000;  // Add correct amount

print("\nCALCULATED CORRECT VALUES:");
print("  remaining_cash: " + shift.remaining_cash + " + 30000 = " + correctRemainingCash);
print("  remaining_transfer: " + shift.remaining_transfer + " - 30000 = " + correctRemainingTransfer);
print("  handed_over_cash: " + shift.handed_over_cash + " - 30000 = " + correctHandedOverCash);
print("  handed_over_transfer: " + (shift.handed_over_transfer || 0) + " + 30000 = " + correctHandedOverTransfer);

// Update shift
db.shifts.updateOne(
  {_id: shiftId},
  {
    $set: {
      remaining_cash: correctRemainingCash,
      remaining_transfer: correctRemainingTransfer,
      handed_over_cash: correctHandedOverCash,
      handed_over_transfer: correctHandedOverTransfer,
      updated_at: new Date()
    }
  }
);

const updatedShift = db.shifts.findOne({_id: shiftId});
print("\nAFTER:");
print("  current_cash: " + (updatedShift.current_cash || 0));
print("  remaining_cash: " + (updatedShift.remaining_cash || 0));
print("  handed_over_cash: " + (updatedShift.handed_over_cash || 0));
print("  transfer_revenue: " + (updatedShift.transfer_revenue || 0));
print("  remaining_transfer: " + (updatedShift.remaining_transfer || 0));
print("  handed_over_transfer: " + (updatedShift.handed_over_transfer || 0));

print("\n=== FIX COMPLETE ===");
print("\nExpected result in frontend:");
print("  💵 Tiền mặt: " + updatedShift.remaining_cash + " ₫ (should be 72,000)");
print("  💳 Tiền CK: " + updatedShift.remaining_transfer + " ₫ (should be 0)");
print("  Đã bàn giao: " + updatedShift.handed_over_transfer + " ₫ (should be 30,000)");
