#!/bin/bash
set -e

echo "🔧 Setup Production Environment"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Check if .env.production exists
if [ ! -f .env.production ]; then
    echo "❌ File .env.production không tồn tại!"
    echo ""
    echo "Tạo file .env.production với nội dung:"
    echo ""
    cat << 'EOF'
# MongoDB Configuration
MONGO_INITDB_ROOT_USERNAME=admin
MONGO_INITDB_ROOT_PASSWORD=108trannhatduat
MONGO_INITDB_DATABASE=cafe_pos

# Backend Configuration
MONGODB_URI=mongodb://admin:108trannhatduat@mongodb:27017/cafe_pos?replicaSet=rs0&authSource=admin
MONGODB_DATABASE=cafe_pos
JWT_SECRET=your_jwt_secret_key_change_this_min_32_chars
PORT=3000

# Frontend Configuration
VITE_API_URL=http://localhost/api
EOF
    exit 1
fi

# Backup existing .env if exists
if [ -f .env ]; then
    echo "📦 Backup .env hiện tại..."
    cp .env .env.backup.$(date +%Y%m%d-%H%M%S)
    echo "✅ Đã backup: .env.backup.$(date +%Y%m%d-%H%M%S)"
    echo ""
fi

# Copy production env
echo "📝 Copy .env.production -> .env"
cp .env.production .env

echo ""
echo "✅ Setup hoàn tất!"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📊 Environment variables:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
cat .env | grep -v '^#' | grep -v '^$'
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "⚠️  LƯU Ý BẢO MẬT:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "1. Đảm bảo .env và .env.production KHÔNG được commit vào git"
echo "2. Thay đổi JWT_SECRET trước khi deploy production"
echo "3. Giữ file .env an toàn, không share cho người khác"
echo ""
echo "Tạo JWT secret mới:"
echo "  openssl rand -base64 32"
echo ""
