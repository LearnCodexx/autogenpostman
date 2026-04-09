#!/bin/bash
# Debug conversion issue

echo "🔍 DEBUGGING POSTMAN CONVERSION FAILURE"
echo "======================================"

echo -e "\n📁 Current Directory:"
pwd

echo -e "\n📦 Check Node.js Tools:"
node --version 2>/dev/null || echo "❌ Node.js not found!"
npm --version 2>/dev/null || echo "❌ NPM not found!"
npx --version 2>/dev/null || echo "❌ NPX not found!"

echo -e "\n📄 Check Generated Swagger:"
if [ -f "docs/swagger.json" ]; then
    echo "✅ swagger.json exists"
    echo "📊 File size: $(wc -c < docs/swagger.json) bytes"
    
    echo -e "\n📊 Swagger paths:"
    cat docs/swagger.json | jq '.paths | keys[]' 2>/dev/null || echo "❌ Invalid JSON or no paths!"
    
    echo -e "\n📊 Swagger info:"
    cat docs/swagger.json | jq '.info' 2>/dev/null || echo "❌ No info section!"
    
    echo -e "\n🧪 Testing Manual Conversion:"
    npx openapi-to-postmanv2 -s docs/swagger.json -o docs/manual-test.json 2>&1
    
    if [ -f "docs/manual-test.json" ]; then
        echo "✅ Manual conversion succeeded!"
        echo "📊 Collection items:"
        cat docs/manual-test.json | jq '.item | length' 2>/dev/null || echo "❌ Invalid collection!"
    else
        echo "❌ Manual conversion failed!"
    fi
else
    echo "❌ swagger.json not found!"
fi

echo -e "\n📂 Directory Contents:"
ls -la docs/ 2>/dev/null || echo "❌ docs/ directory not found!"

echo -e "\n🔍 Check api/main.go:"
if [ -f "api/main.go" ]; then
    echo "✅ api/main.go exists"
    echo -e "\n📊 Swagger annotations:"
    grep -n "@" api/main.go || echo "❌ No swagger annotations found!"
    
    echo -e "\n📊 Route definitions:"
    grep -n "GET\\|POST\\|PUT\\|DELETE" api/main.go || echo "❌ No routes found!"
else
    echo "❌ api/main.go not found!"
fi

echo -e "\n🚀 SUGGESTED FIXES:"
echo "1. Add swagger annotations to all handlers"
echo "2. Install Node.js: sudo apt install nodejs npm"
echo "3. Add docs import: import _ \"yourmodule/docs\""
echo "4. Use explicit path: --swagger-input=docs/swagger.json"