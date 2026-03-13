#!/bin/bash
# Lock all orders with status CREATED in MongoDB (bulk data fix)
# Usage: ./cancel-created-orders.sh [--dry-run] [--uri "mongodb://..."]

set -euo pipefail

# ── Config ─────────────────────────────────────────────────────────────────────
MONGO_URI="mongodb://admin:108trannhatduat@localhost:27017/cafe_pos?directConnection=true&authSource=admin"
DB_NAME="cafe_pos"
COLLECTION="orders"
TARGET_STATUS="CREATED"
NEW_STATUS="LOCKED"
DRY_RUN=false

# ── Parse args ─────────────────────────────────────────────────────────────────
while [[ $# -gt 0 ]]; do
  case $1 in
    --dry-run)
      DRY_RUN=true
      shift
      ;;
    --uri)
      MONGO_URI="$2"
      shift 2
      ;;
    *)
      echo "Usage: $0 [--dry-run] [--uri \"mongodb://...\"]"
      exit 1
      ;;
  esac
done

# ── Detect mongo client ────────────────────────────────────────────────────────
if command -v mongosh &>/dev/null; then
  MONGO_CMD="mongosh"
elif command -v mongo &>/dev/null; then
  MONGO_CMD="mongo"
elif command -v python3 &>/dev/null && python3 -c "import pymongo" &>/dev/null 2>&1; then
  MONGO_CMD="python3"
else
  echo "❌ No MongoDB client found (mongosh, mongo, or python3+pymongo)."
  echo ""
  echo "To install mongosh on Ubuntu/Debian:"
  echo "  curl -fsSL https://www.mongodb.org/static/pgp/server-7.0.asc | sudo gpg -o /usr/share/keyrings/mongodb-server-7.0.gpg --dearmor"
  echo "  echo 'deb [ arch=amd64,arm64 signed-by=/usr/share/keyrings/mongodb-server-7.0.gpg ] https://repo.mongodb.org/apt/ubuntu jammy/mongodb-org/7.0 multiverse' | sudo tee /etc/apt/sources.list.d/mongodb-org-7.0.list"
  echo "  sudo apt-get update && sudo apt-get install -y mongodb-mongosh"
  echo ""
  echo "Or install pymongo (usually faster on Ubuntu):"
  echo "  pip3 install pymongo"
  exit 1
fi

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🔒  Lock Orders — Status: $TARGET_STATUS → $NEW_STATUS"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "DB:   $DB_NAME"
[[ $DRY_RUN == true ]] && echo "MODE: DRY RUN (no changes)" || echo "MODE: LIVE"
echo ""

# ── Helper: run mongo JS or python ────────────────────────────────────────────
run_mongo() {
  local js="$1"
  if [[ "$MONGO_CMD" == "python3" ]]; then
    python3 - <<PYEOF
import pymongo, datetime, sys
client = pymongo.MongoClient("$MONGO_URI")
db = client["$DB_NAME"]
$js
PYEOF
  else
    $MONGO_CMD "$MONGO_URI" --quiet --eval "$js"
  fi
}

# ── Count first ────────────────────────────────────────────────────────────────
if [[ "$MONGO_CMD" == "python3" ]]; then
  COUNT=$(python3 - <<PYEOF
import pymongo
client = pymongo.MongoClient("$MONGO_URI")
db = client["$DB_NAME"]
print(db["$COLLECTION"].count_documents({"status": "$TARGET_STATUS"}))
PYEOF
)
else
  COUNT=$($MONGO_CMD "$MONGO_URI" --quiet --eval "
    db = db.getSiblingDB('$DB_NAME');
    db.$COLLECTION.countDocuments({ status: '$TARGET_STATUS' });
  ")
fi

echo "📋 Found $COUNT order(s) with status '$TARGET_STATUS'"
echo ""

if [[ "$COUNT" -eq 0 ]]; then
  echo "✅ Nothing to lock."
  exit 0
fi

# ── Preview ────────────────────────────────────────────────────────────────────
echo "Preview (first 10):"
if [[ "$MONGO_CMD" == "python3" ]]; then
  python3 - <<PYEOF
import pymongo
client = pymongo.MongoClient("$MONGO_URI")
db = client["$DB_NAME"]
for o in db["$COLLECTION"].find({"status": "$TARGET_STATUS"}, {"order_number":1,"waiter_name":1,"total":1,"created_at":1}).limit(10):
    print(f"  #{o.get('order_number','')} | {o.get('waiter_name','')} | {o.get('total',0)} VND | {o.get('created_at','')}")
PYEOF
else
  $MONGO_CMD "$MONGO_URI" --quiet --eval "
    db = db.getSiblingDB('$DB_NAME');
    db.$COLLECTION.find(
      { status: '$TARGET_STATUS' },
      { order_number: 1, waiter_name: 1, total: 1, created_at: 1, status: 1 }
    ).limit(10).forEach(o => {
      print('  #' + o.order_number + ' | ' + o.waiter_name + ' | ' + o.total + ' VND | ' + o.created_at);
    });
  "
fi
echo ""

# ── Dry run exit ───────────────────────────────────────────────────────────────
if [[ $DRY_RUN == true ]]; then
  echo "🔍 Dry run complete. $COUNT order(s) would be locked."
  echo "   Run without --dry-run to apply changes."
  exit 0
fi

# ── Confirm ────────────────────────────────────────────────────────────────────
read -r -p "⚠️  Lock $COUNT order(s) as LOCKED? [y/N] " CONFIRM
if [[ "$(echo "$CONFIRM" | tr '[:upper:]' '[:lower:]')" != "y" ]]; then
  echo "Aborted."
  exit 0
fi

# ── Execute ────────────────────────────────────────────────────────────────────
NOW=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

if [[ "$MONGO_CMD" == "python3" ]]; then
  RESULT=$(python3 - <<PYEOF
import pymongo, datetime
client = pymongo.MongoClient("$MONGO_URI")
db = client["$DB_NAME"]
now = datetime.datetime.utcnow()
result = db["$COLLECTION"].update_many(
    {"status": "$TARGET_STATUS"},
    {"\$set": {"status": "$NEW_STATUS", "locked_at": now, "updated_at": now}}
)
print(result.modified_count)
PYEOF
)
else
  RESULT=$($MONGO_CMD "$MONGO_URI" --quiet --eval "
    db = db.getSiblingDB('$DB_NAME');
    const result = db.$COLLECTION.updateMany(
      { status: '$TARGET_STATUS' },
      {
        \$set: {
          status: '$NEW_STATUS',
          locked_at: new Date('$NOW'),
          updated_at: new Date('$NOW')
        }
      }
    );
    print(result.modifiedCount);
  ")
fi

echo ""
echo "✅ Locked $RESULT order(s)."
echo "   status: $TARGET_STATUS → $NEW_STATUS"
echo "   locked_at: $NOW"
