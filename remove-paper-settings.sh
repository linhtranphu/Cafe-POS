#!/bin/bash

echo "🔧 Removing paper_width and label_size from shop_settings..."
echo "These settings belong to printer configuration, not shop settings."
echo ""

cd backend/cmd/migrate
go run remove_paper_settings_from_shop.go
