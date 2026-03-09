#!/bin/bash

# Script to backup MongoDB from Docker container
# Usage: ./backup-mongodb.sh [output_directory]

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
CONTAINER_NAME="cafe-pos-mongodb"
DB_NAME="cafe_pos"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
OUTPUT_DIR="${1:-./mongodb_backup_${TIMESTAMP}}"

echo -e "${GREEN}=== MongoDB Backup Script ===${NC}"
echo ""

# Check if container is running
if ! docker ps --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
    echo -e "${RED}Error: Container '$CONTAINER_NAME' is not running${NC}"
    exit 1
fi

# Create output directory
mkdir -p "$OUTPUT_DIR"

# Perform backup
echo "Backing up database '$DB_NAME' to '$OUTPUT_DIR'..."
docker exec $CONTAINER_NAME mongodump \
    --db $DB_NAME \
    --out /tmp/backup

# Copy backup from container
echo "Copying backup files from container..."
docker cp $CONTAINER_NAME:/tmp/backup/$DB_NAME "$OUTPUT_DIR/"

# Clean up container
echo "Cleaning up temporary files in container..."
docker exec $CONTAINER_NAME rm -rf /tmp/backup

# Create backup info file
cat > "$OUTPUT_DIR/backup_info.txt" << EOF
Backup Information
==================
Database: $DB_NAME
Container: $CONTAINER_NAME
Timestamp: $TIMESTAMP
Date: $(date)
EOF

echo ""
echo -e "${GREEN}=== Backup Complete ===${NC}"
echo ""
echo "Backup saved to: $OUTPUT_DIR"
echo ""
echo "To restore this backup, run:"
echo "  ./restore-mongodb.sh $OUTPUT_DIR/$DB_NAME"
