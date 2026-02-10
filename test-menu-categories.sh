#!/bin/bash

# Test Menu Categories API
# Make sure backend is running on port 8080

BASE_URL="http://localhost:8080/api"
TOKEN=""

echo "🧪 Testing Menu Categories API"
echo "================================"
echo ""

# Function to make authenticated requests
auth_request() {
    if [ -z "$TOKEN" ]; then
        echo "❌ No token set. Please login first."
        return 1
    fi
    curl -s -H "Authorization: Bearer $TOKEN" "$@"
}

# 1. Login as admin to get token
echo "1️⃣ Logging in as admin..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
    -H "Content-Type: application/json" \
    -d '{"username":"admin","password":"admin123"}')

TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    echo "❌ Login failed!"
    echo "Response: $LOGIN_RESPONSE"
    exit 1
fi

echo "✅ Login successful!"
echo ""

# 2. Get all categories (should be empty initially)
echo "2️⃣ Getting all categories..."
CATEGORIES=$(auth_request "$BASE_URL/manager/menu-categories")
echo "Response: $CATEGORIES"
echo ""

# 3. Create a new category
echo "3️⃣ Creating new category 'Cà phê'..."
CREATE_RESPONSE=$(auth_request -X POST "$BASE_URL/manager/menu-categories" \
    -H "Content-Type: application/json" \
    -d '{"name":"Cà phê"}')
echo "Response: $CREATE_RESPONSE"
CATEGORY_ID=$(echo $CREATE_RESPONSE | grep -o '"id":"[^"]*' | cut -d'"' -f4)
echo "Category ID: $CATEGORY_ID"
echo ""

# 4. Create another category
echo "4️⃣ Creating new category 'Trà'..."
CREATE_RESPONSE2=$(auth_request -X POST "$BASE_URL/manager/menu-categories" \
    -H "Content-Type: application/json" \
    -d '{"name":"Trà"}')
echo "Response: $CREATE_RESPONSE2"
echo ""

# 5. Try to create duplicate (should fail)
echo "5️⃣ Trying to create duplicate 'Cà phê' (should fail)..."
DUPLICATE_RESPONSE=$(auth_request -X POST "$BASE_URL/manager/menu-categories" \
    -H "Content-Type: application/json" \
    -d '{"name":"Cà phê"}')
echo "Response: $DUPLICATE_RESPONSE"
echo ""

# 6. Get all categories again
echo "6️⃣ Getting all categories again..."
CATEGORIES=$(auth_request "$BASE_URL/manager/menu-categories")
echo "Response: $CATEGORIES"
echo ""

# 7. Get single category
if [ ! -z "$CATEGORY_ID" ]; then
    echo "7️⃣ Getting single category by ID..."
    SINGLE_CATEGORY=$(auth_request "$BASE_URL/manager/menu-categories/$CATEGORY_ID")
    echo "Response: $SINGLE_CATEGORY"
    echo ""
fi

# 8. Update category
if [ ! -z "$CATEGORY_ID" ]; then
    echo "8️⃣ Updating category name to 'Cà phê sữa'..."
    UPDATE_RESPONSE=$(auth_request -X PUT "$BASE_URL/manager/menu-categories/$CATEGORY_ID" \
        -H "Content-Type: application/json" \
        -d '{"name":"Cà phê sữa"}')
    echo "Response: $UPDATE_RESPONSE"
    echo ""
fi

# 9. Try to delete category (should succeed if no menu items)
if [ ! -z "$CATEGORY_ID" ]; then
    echo "9️⃣ Deleting category..."
    DELETE_RESPONSE=$(auth_request -X DELETE "$BASE_URL/manager/menu-categories/$CATEGORY_ID")
    echo "Response: $DELETE_RESPONSE"
    echo ""
fi

# 10. Get all categories final check
echo "🔟 Final check - Getting all categories..."
CATEGORIES=$(auth_request "$BASE_URL/manager/menu-categories")
echo "Response: $CATEGORIES"
echo ""

echo "================================"
echo "✅ Test completed!"
