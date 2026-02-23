#!/bin/bash

echo "📋 Listing All Open Shifts"
echo "=========================="
echo ""

cd backend
go run cmd/list-open-shifts/main.go
