#!/bin/bash

echo "🔄 Closing All Open Waiter Shifts (Direct MongoDB)"
echo "==================================================="
echo ""

# Try mongosh first, fallback to mongo
if command -v mongosh &> /dev/null; then
    mongosh cafe_pos -u admin -p password123 --authenticationDatabase admin < close-all-shifts-mongo.js
elif command -v mongo &> /dev/null; then
    mongo cafe_pos -u admin -p password123 --authenticationDatabase admin < close-all-shifts-mongo.js
else
    echo "❌ Error: Neither mongosh nor mongo command found"
    echo "Please install MongoDB shell"
    exit 1
fi

echo ""
echo "Done!"
