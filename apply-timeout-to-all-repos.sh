#!/bin/bash

# Script to apply timeout and graceful handling to all MongoDB repositories
# This script adds WithQueryTimeout() and IsCollectionNotFoundError() to all Find methods

set -e

echo "=========================================="
echo "🔧 Applying Timeout to All Repositories"
echo "=========================================="
echo ""

REPO_DIR="backend/infrastructure/mongodb"

# List of repository files to update
REPOS=(
    "ingredient_repository.go"
    "user_repository.go"
    "shift_repository.go"
    "cashier_shift_repository.go"
    "batch_definition_repository.go"
    "batch_record_repository.go"
    "batch_usage_log_repository.go"
    "order_item_repository.go"
    "printer_config_repository.go"
    "print_template_repository.go"
    "shop_settings_repository.go"
    "expense_repository.go"
    "fund_transaction_repository.go"
    "cash_handover_repository.go"
    "stock_history_repository.go"
    "menu_category_repository.go"
    "facility_repository.go"
    "operating_expense_repository.go"
    "print_notification_repository.go"
    "cash_discrepancy_repository.go"
    "fund_handover_repository.go"
    "payment_audit_repository.go"
    "payment_discrepancy_repository.go"
    "cash_reconciliation_repository.go"
)

echo "📋 Repositories to update: ${#REPOS[@]}"
echo ""

# Backup original files
echo "💾 Creating backups..."
for repo in "${REPOS[@]}"; do
    if [ -f "$REPO_DIR/$repo" ]; then
        cp "$REPO_DIR/$repo" "$REPO_DIR/$repo.backup"
        echo "  ✓ Backed up $repo"
    fi
done

echo ""
echo "✅ Backups created"
echo ""
echo "⚠️  Manual update required for each repository:"
echo ""
echo "For each Find method that returns array:"
echo "  1. Add: ctx, cancel := WithQueryTimeout(ctx); defer cancel()"
echo "  2. Add graceful error handling for IsCollectionNotFoundError"
echo "  3. Ensure nil arrays are converted to empty arrays"
echo ""
echo "For each FindOne/Create/Update/Delete method:"
echo "  1. Add: ctx, cancel := WithQueryTimeout(ctx); defer cancel()"
echo ""
echo "📝 See MONGODB_GRACEFUL_HANDLING_SUMMARY.md for patterns"
echo ""
echo "To restore backups if needed:"
echo "  for f in $REPO_DIR/*.backup; do mv \"\$f\" \"\${f%.backup}\"; done"
echo ""
