#!/bin/bash
# Debug script untuk troubleshoot empty collection

echo "🔍 DEBUGGING EMPTY POSTMAN COLLECTION"
echo "====================================="

echo -e "\n📁 Current Directory:"
pwd

echo -e "\n📝 Found main.go files:"
find . -name "main.go" -type f

echo -e "\n🎯 Checking api/main.go structure:"
if [ -f "api/main.go" ]; then
    echo "✅ api/main.go exists"
    echo -e "\n📊 File size:"
    wc -l api/main.go
    
    echo -e "\n🔍 Swagger annotations found:"
    grep -n "@" api/main.go || echo "❌ No swagger annotations found!"
    
    echo -e "\n🔍 Route definitions:"
    grep -n "GET\|POST\|PUT\|DELETE" api/main.go || echo "❌ No route definitions found!"
else
    echo "❌ api/main.go not found!"
fi

echo -e "\n📦 Checking generated swagger.json:"
if [ -f "docs/swagger.json" ]; then
    echo "✅ swagger.json exists"
    echo -e "\n📊 Swagger paths:"
    cat docs/swagger.json | jq '.paths | keys[]' 2>/dev/null || echo "❌ Invalid JSON or no paths!"
    
    echo -e "\n📊 Swagger info:"
    cat docs/swagger.json | jq '.info' 2>/dev/null || echo "❌ No info section!"
else
    echo "❌ swagger.json not found!"
fi

echo -e "\n📦 Checking collection.json:"
if [ -f "docs/collection.json" ]; then
    echo "✅ collection.json exists"
    echo -e "\n📊 Collection items count:"
    cat docs/collection.json | jq '.item | length' 2>/dev/null || echo "❌ Invalid JSON!"
    
    echo -e "\n📊 Collection info:"
    cat docs/collection.json | jq '.info.name, .info.description.content' 2>/dev/null
else
    echo "❌ collection.json not found!"
fi

echo -e "\n🚀 SOLUTION STEPS:"
echo "1. Check api/main.go has proper swagger annotations"
echo "2. Ensure @Router paths match actual routes"  
echo "3. Add docs import: import _ \"yourmodule/docs\""
echo "4. Regenerate: go run cmd/postman-generator/main.go"