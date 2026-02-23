#!/bin/bash

# Fix MongoDB Transaction Error
# This script helps setup MongoDB replica set for transactions

echo "🔧 MongoDB Transaction Fix"
echo "=========================="
echo ""

echo "Vấn đề: Backend cần MongoDB replica set để chạy transactions"
echo ""

echo "Giải pháp:"
echo "1. Dừng MongoDB hiện tại"
echo "2. Khởi động MongoDB với replica set"
echo "3. Khởi tạo replica set"
echo ""

read -p "Bạn có muốn tiếp tục? (y/n) " -n 1 -r
echo ""

if [[ ! $REPLY =~ ^[Yy]$ ]]
then
    exit 1
fi

echo ""
echo "📋 Hướng dẫn thực hiện:"
echo ""

echo "Bước 1: Dừng MongoDB"
echo "  sudo systemctl stop mongod"
echo "  # Hoặc nếu chạy thủ công:"
echo "  # pkill mongod"
echo ""

echo "Bước 2: Khởi động MongoDB với replica set"
echo "  mongod --replSet rs0 --dbpath /path/to/data"
echo "  # Hoặc với systemd:"
echo "  # Sửa file /etc/mongod.conf:"
echo "  # replication:"
echo "  #   replSetName: rs0"
echo "  # Sau đó:"
echo "  # sudo systemctl start mongod"
echo ""

echo "Bước 3: Khởi tạo replica set"
echo "  mongosh"
echo "  rs.initiate()"
echo "  # Đợi vài giây"
echo "  rs.status()"
echo ""

echo "Bước 4: Khởi động lại backend"
echo "  cd backend"
echo "  go run main.go"
echo ""

echo "✅ Sau khi làm xong, test lại trên trình duyệt!"
