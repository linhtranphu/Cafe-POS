#!/bin/bash

echo "🔨 Rebuilding frontend..."

cd frontend

# Clean old build
echo "🧹 Cleaning old build..."
rm -rf dist

# Install dependencies (if needed)
echo "📦 Checking dependencies..."
npm install

# Build
echo "🏗️ Building..."
npm run build

if [ $? -eq 0 ]; then
    echo "✅ Frontend build successful!"
    echo ""
    echo "📋 Next steps:"
    echo "1. Deploy to server: scp -r dist/* user@tacafe.store:/path/to/web/root/"
    echo "2. Or if using docker: docker-compose up -d --build frontend"
    echo "3. Clear browser cache: Ctrl+Shift+R or Cmd+Shift+R"
else
    echo "❌ Build failed!"
    exit 1
fi
