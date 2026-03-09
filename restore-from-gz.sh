#!/bin/bash

# Script to restore MongoDB from .gz backup file
# Usage: ./restore-from-gz.sh <backup_file.gz>

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
CONTAINER_NAME="cafe-pos-mongodb"
DB_NAME="cafe_pos"
BACKUP_FILE="${1:-./mongodb_file/mongodb-backup-2026-03-07_01-00-05.gz}"

echo -e "${GREEN}=== MongoDB Restore from .gz Script ===${NC}"
echo ""

# Check if backup file exists
if [ ! -f "$BACKUP_FILE" ]; then
    echo -e "${RED}Error: Backup file '$BACKUP_FILE' not found${NC}"
    echo "Usage: ./restore-from-gz.sh <backup_file.gz>"
    exit 1
fi

echo "Backup file: $BACKUP_FILE"
echo ""

# Check if container exists
if ! docker ps -a --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
    echo -e "${YELLOW}Container '$CONTAINER_NAME' does not exist. Creating...${NC}"
    docker run -d --name $CONTAINER_NAME \
        -p 27017:27017 \
        -e MONGO_INITDB_ROOT_USERNAME=admin \
        -e MONGO_INITDB_ROOT_PASSWORD=admin123 \
        mongo:7.0 --replSet rs0 --bind_ip_all
    echo "Waiting for MongoDB to start..."
    sleep 10
fi

# Start container if not running
if ! docker ps --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
    echo -e "${YELLOW}Starting container '$CONTAINER_NAME'...${NC}"
    docker start $CONTAINER_NAME
    sleep 5
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

# Copy backup file to container
echo ""
echo "Copying backup file to container..."
docker cp "$BACKUP_FILE" $CONTAINER_NAME:/tmp/backup.gz

# Restore using mongorestore with gzip archive
echo ""
echo "Restoring database from gzip archive..."
docker exec $CONTAINER_NAME mongorestore --host localhost:27017 --username admin --password 108trannhatduat --authenticationDatabase admin --gzip --archive=/tmp/backup.gz --drop

# Clean up
echo ""
echo "Cleaning up temporary files..."
docker exec $CONTAINER_NAME rm -f /tmp/backup.gz

# Initialize replica set if needed
echo ""
echo "Checking replica set status..."
RS_STATUS=$(docker exec $CONTAINER_NAME mongosh --quiet --eval "rs.status().ok" 2>/dev/null || echo "0")

if [ "$RS_STATUS" = "0" ]; then
    echo -e "${YELLOW}Replica set not initialized. Initializing...${NC}"
    docker exec $CONTAINER_NAME mongosh --eval "rs.initiate({
        _id: 'rs0',
        members: [{ _id: 0, host: 'localhost:27017' }]
    })" > /dev/null 2>&1
    echo "Waiting for replica set to initialize..."
    sleep 5
    echo -e "${GREEN}Replica set initialized!${NC}"
else
    echo -e "${GREEN}Replica set already initialized${NC}"
fi

# Verify restore
echo ""
echo "Verifying restore..."
DBS=$(docker exec $CONTAINER_NAME mongosh --quiet --eval "db.adminCommand('listDatabases').databases.map(d => d.name).join(', ')")
echo -e "${GREEN}Databases found: $DBS${NC}"

if echo "$DBS" | grep -q "$DB_NAME"; then
    COLLECTIONS=$(docker exec $CONTAINER_NAME mongosh $DB_NAME --quiet --eval "db.getCollectionNames().length")
    echo -e "${GREEN}Database '$DB_NAME' restored successfully!${NC}"
    echo "Collections found: $COLLECTIONS"
else
    echo -e "${YELLOW}Warning: Database '$DB_NAME' not found. Available databases: $DBS${NC}"
fi

echo ""
echo -e "${GREEN}=== Restore Complete ===${NC}"
echo ""
echo "MongoDB is running on: localhost:27017"
echo "You can now start your backend application"
