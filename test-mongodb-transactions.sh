#!/bin/bash

echo "🔍 Testing MongoDB Transactions Support..."
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test MongoDB connection and transactions
mongosh --quiet --eval "
use cafe_pos;

// Test 1: Basic connection
try {
    db.runCommand({ ping: 1 });
    print('${GREEN}✓${NC} Connected to MongoDB');
} catch(e) {
    print('${RED}✗${NC} Cannot connect to MongoDB: ' + e);
    quit(1);
}

// Test 2: Replica set status
try {
    var status = rs.status();
    if (status.ok === 1) {
        print('${GREEN}✓${NC} Replica set is active');
        print('  - Set name: ${YELLOW}' + status.set + '${NC}');
        var primary = status.members.find(m => m.stateStr === 'PRIMARY');
        if (primary) {
            print('  - Primary: ${YELLOW}' + primary.name + '${NC}');
        }
    } else {
        print('${RED}✗${NC} Replica set not configured properly');
        quit(1);
    }
} catch(e) {
    print('${RED}✗${NC} Replica set not initialized: ' + e);
    print('');
    print('${YELLOW}ℹ${NC}  MongoDB is running in standalone mode.');
    print('${YELLOW}ℹ${NC}  Transactions require replica set mode.');
    print('');
    print('To fix this, see: MONGODB_REPLICA_SET_SETUP.md');
    quit(1);
}

// Test 3: Transaction support
try {
    var session = db.getMongo().startSession();
    session.startTransaction();
    
    session.getDatabase('cafe_pos').test_transactions.insertOne({
        test: 'transaction_test',
        timestamp: new Date()
    });
    
    session.commitTransaction();
    session.endSession();
    
    print('${GREEN}✓${NC} Transactions are working!');
    
    // Cleanup
    db.test_transactions.drop();
} catch(e) {
    print('${RED}✗${NC} Transaction failed: ' + e);
    quit(1);
}

// Test 4: Multi-document transaction
try {
    var session = db.getMongo().startSession();
    session.startTransaction();
    
    var testDb = session.getDatabase('cafe_pos');
    testDb.test_multi_1.insertOne({doc: 1});
    testDb.test_multi_2.insertOne({doc: 2});
    
    session.commitTransaction();
    session.endSession();
    
    print('${GREEN}✓${NC} Multi-document transactions working!');
    
    // Cleanup
    db.test_multi_1.drop();
    db.test_multi_2.drop();
} catch(e) {
    print('${RED}✗${NC} Multi-document transaction failed: ' + e);
    quit(1);
}

print('');
print('${GREEN}✅ All tests passed!${NC}');
print('${GREEN}✅ MongoDB is ready for batch management system.${NC}');
print('');
print('You can now run backend tests:');
print('  ${YELLOW}cd backend && go test -v -run=\"Batch\" ./application/services/...${NC}');
"

exit_code=$?

if [ $exit_code -ne 0 ]; then
    echo ""
    echo "❌ MongoDB is not configured for transactions."
    echo ""
    echo "📖 Please follow the setup guide:"
    echo "   cat MONGODB_REPLICA_SET_SETUP.md"
    echo ""
    echo "🚀 Quick fix (Docker):"
    echo "   1. Update docker-compose.yml to add replica set config"
    echo "   2. docker-compose down && docker-compose up -d mongodb"
    echo "   3. Wait 15 seconds, then run this script again"
fi

exit $exit_code
