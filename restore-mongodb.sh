#!/bin/bash

# Script to restore MongoDB backup to Docker container
# Usage: ./restore-mongodb.sh [backup_directory]

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
CONTAINER_NAME="cafe-pos-mongodb"
DB_NAME="cafe_pos"
BACKUP_DIR="${1:-./mongodb_backup}"

echo -e "${GREEN}=== MongoDB Restore Script ===${NC}"
echo ""

# Check if backup directory exists
if [ ! -d "$BACKUP_DIR" ]; then
    echo -e "${RED}Error: Backup directory '$BACKUP_DIR' not found${NC}"
    echo "Usage: ./restore-mongodb.sh [backup_directory]"
    exit 1
fi

# Check if container is running
if ! docker ps --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
    echo -e "${YELLOW}Container '$CONTAINER_NAME' is not running. Starting it...${NC}"
    
    # Try to start existing container
    if docker ps -a --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
        docker start $CONTAINER_NAME
        echo "Waiting for MongoDB to be ready..."
        sleep 10
    else
        echo -e "${RED}Error: Container '$CONTAINER_NAME' does not exist${NC}"
        echo "Please create the container first using docker-compose"
        exit 1
    fi
fi

# Wait for MongoDB to be ready
echo "Checking MongoDB connection..."
MAX_RETRIES=30
RETRY_COUNT=0

while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
    if docker exec $CONTAINER_NAME mongosh --quiet --eval "db.adminCommand('ping')" > /dev/null 2>&1; then
        echo -e "${GREEN}MongoDB is ready!${NC}"
        break
    fi
    
    RETRY_COUNT=$((RETRY_COUNT + 1))
    echo "Waiting for MongoDB... ($RETRY_COUNT/$MAX_RETRIES)"
    sleep 2
done

if [ $RETRY_COUNT -eq $MAX_RETRIES ]; then
    echo -e "${RED}Error: MongoDB did not become ready in time${NC}"
    exit 1
fi

# Copy backup to container
echo ""
echo "Copying backup files to container..."
docker cp "$BACKUP_DIR" $CONTAINER_NAME:/tmp/backup

# Restore the backup
echo ""
echo "Restoring database '$DB_NAME'..."
docker exec $CONTAINER_NAME mongorestore \
    --db $DB_NAME \
    --drop \
    /tmp/backup/$DB_NAME

# Clean up
echo ""
echo "Cleaning up temporary files..."
docker exec $CONTAINER_NAME rm -rf /tmp/backup

# Initialize replica set if needed
echo ""
echo "Checking replica set status..."
RS_STATUS=$(docker exec $CONTAINER_NAME mongosh --quiet --eval "rs.status().ok" 2>/dev/null || echo "0")

if [ "$RS_STATUS" = "0" ]; then
    echo -e "${YELLOW}Replica set not initialized. Initializing...${NC}"
    docker exec $CONTAINER_NAME mongosh --eval "rs.initiate({
        _id: 'rs0',
        members: [{ _id: 0, host: 'localhost:27017' }]
    })"
    echo "Waiting for replica set to initialize..."
    sleep 5
    echo -e "${GREEN}Replica set initialized!${NC}"
else
    echo -e "${GREEN}Replica set already initialized${NC}"
fi

# Verify restore
echo ""
echo "Verifying restore..."
COLLECTIONS=$(docker exec $CONTAINER_NAME mongosh $DB_NAME --quiet --eval "db.getCollectionNames().length")
echo -e "${GREEN}Database restored successfully!${NC}"
echo "Collections found: $COLLECTIONS"

echo ""
echo -e "${GREEN}=== Restore Complete ===${NC}"
echo ""
echo "You can now start your backend application"
