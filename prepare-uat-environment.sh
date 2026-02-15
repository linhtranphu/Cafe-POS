#!/bin/bash

# Batch UAT Environment Setup Script
# This script prepares the environment for User Acceptance Testing

set -e

echo "========================================="
echo "Batch UAT Environment Setup"
echo "========================================="
echo ""

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Function to print colored output
print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

# Step 1: Check if MongoDB is running
echo "Step 1: Checking MongoDB..."
if docker ps | grep -q mongodb; then
    print_success "MongoDB is running"
else
    print_warning "MongoDB not running. Starting MongoDB..."
    docker-compose up -d mongodb
    sleep 5
    print_success "MongoDB started"
fi
echo ""

# Step 2: Check if backend is running
echo "Step 2: Checking Backend..."
if curl -s http://localhost:8080/health > /dev/null 2>&1; then
    print_success "Backend is running"
else
    print_error "Backend is not running"
    echo "Please start backend manually:"
    echo "  cd backend && go run main.go"
    exit 1
fi
echo ""

# Step 3: Check if frontend is running
echo "Step 3: Checking Frontend..."
if curl -s http://localhost:5173 > /dev/null 2>&1; then
    print_success "Frontend is running"
else
    print_error "Frontend is not running"
    echo "Please start frontend manually:"
    echo "  cd frontend && npm run dev"
    exit 1
fi
echo ""

# Step 4: Seed test data
echo "Step 4: Seeding test data..."
echo "This will create:"
echo "  - Admin user (manager@test.com / password123)"
echo "  - Barista user (barista@test.com / password123)"
echo "  - Sample ingredients"
echo "  - Sample menu items"
echo ""
read -p "Do you want to seed test data? (y/n) " -n 1 -r
echo ""
if [[ $REPLY =~ ^[Yy]$ ]]; then
    cd backend
    
    echo "Creating admin user..."
    go run cmd/seed-admin/main.go
    print_success "Admin user created"
    
    echo "Creating sample ingredients..."
    go run cmd/seed/main.go
    print_success "Sample ingredients created"
    
    echo "Creating sample menu items..."
    go run cmd/seed-menu/main.go
    print_success "Sample menu items created"
    
    cd ..
else
    print_warning "Skipped test data seeding"
fi
echo ""

# Step 5: Create sample batch definitions
echo "Step 5: Creating sample batch definitions..."
read -p "Do you want to create sample batch definitions? (y/n) " -n 1 -r
echo ""
if [[ $REPLY =~ ^[Yy]$ ]]; then
    # Get auth token
    echo "Getting auth token..."
    TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
        -H "Content-Type: application/json" \
        -d '{"email":"manager@test.com","password":"password123"}' \
        | grep -o '"token":"[^"]*' | cut -d'"' -f4)
    
    if [ -z "$TOKEN" ]; then
        print_error "Failed to get auth token"
        exit 1
    fi
    print_success "Auth token obtained"
    
    # Get first ingredient ID
    INGREDIENT_ID=$(curl -s -X GET http://localhost:8080/api/ingredients \
        -H "Authorization: Bearer $TOKEN" \
        | grep -o '"id":"[^"]*' | head -1 | cut -d'"' -f4)
    
    if [ -z "$INGREDIENT_ID" ]; then
        print_error "No ingredients found. Please seed ingredients first."
        exit 1
    fi
    
    # Create batch definition
    echo "Creating 'Cà Phê Concentrate' batch definition..."
    curl -s -X POST http://localhost:8080/api/batch-definitions \
        -H "Authorization: Bearer $TOKEN" \
        -H "Content-Type: application/json" \
        -d "{
            \"name\": \"Cà Phê Concentrate\",
            \"unit\": \"ml\",
            \"shelf_life_hours\": 24,
            \"conversion_rates\": [
                {
                    \"source_ingredient_id\": \"$INGREDIENT_ID\",
                    \"source_quantity\": 100,
                    \"source_unit\": \"g\",
                    \"batch_quantity\": 500,
                    \"wastage_rate\": 0.1
                }
            ],
            \"low_stock_threshold\": 500,
            \"expiry_warning_hours\": 4
        }" > /dev/null
    
    print_success "Sample batch definition created"
else
    print_warning "Skipped batch definition creation"
fi
echo ""

# Step 6: Summary
echo "========================================="
echo "Setup Complete!"
echo "========================================="
echo ""
echo "Test Accounts:"
echo "  Manager: manager@test.com / password123"
echo "  Barista: barista@test.com / password123"
echo ""
echo "URLs:"
echo "  Backend API: http://localhost:8080"
echo "  Frontend: http://localhost:5173"
echo ""
echo "Next Steps:"
echo "  1. Open BATCH_UAT_GUIDE.md for detailed test cases"
echo "  2. Open BATCH_UAT_CHECKLIST.md for quick checklist"
echo "  3. Start testing!"
echo ""
print_success "Ready for UAT!"

