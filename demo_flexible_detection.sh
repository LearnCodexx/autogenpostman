#!/bin/bash
# Demo of flexible file detection

echo "🚀 FLEXIBLE FILE DETECTION DEMO"
echo "==============================="

echo -e "\n📁 Creating test structures..."

# Test 1: Dedicated API file (best practice)
mkdir -p test-flexible/api
cat > test-flexible/api/routes.go << 'EOF'
package main

// @title Flexible API
// @version 1.0
// @host localhost:8080

import "github.com/gin-gonic/gin"

func main() {
    r := gin.Default()
    
    // @Router /users [get]
    r.GET("/users", getUsers)
    
    r.Run(":8080")
}

// @Summary Get users
// @Router /users [get]
func getUsers() {}
EOF

# Test 2: Main.go with mixed business logic (traditional)
cat > test-flexible/main.go << 'EOF'
package main

import "github.com/gin-gonic/gin"

// Business logic mixed with API
func main() {
    // Complex business initialization
    setupDatabase()
    startBackgroundJobs()
    
    r := gin.Default()
    r.GET("/health", healthCheck)  // Only one simple endpoint
    r.Run(":8080")
}

func setupDatabase() {}
func startBackgroundJobs() {}
func healthCheck() {}
EOF

# Test 3: Postman-specific file (ideal)
cat > test-flexible/postman.go << 'EOF'
package main

// @title Clean API for Postman
// @version 2.0
// @host localhost:8080
// @BasePath /api/v1

import "github.com/gin-gonic/gin"

// @Summary Health check
// @Tags health
// @Router /api/v1/health [get]
func healthEndpoint() {}

// @Summary Get users  
// @Tags users
// @Router /api/v1/users [get]
func getUsersEndpoint() {}

// @Summary Create user
// @Tags users
// @Router /api/v1/users [post] 
func createUserEndpoint() {}
EOF

echo "✅ Test files created!"

echo -e "\n🔍 Testing flexible detection..."
cd test-flexible

echo -e "\n=== Files with swagger annotations ==="
grep -l "@" *.go */*.go 2>/dev/null || echo "None found"

echo -e "\n=== Annotation counts ==="
for f in $(find . -name "*.go"); do
    count=$(grep -c "@title\|@Router\|@Summary" "$f" 2>/dev/null || echo "0")
    echo "$f: $count annotations"
done

echo -e "\n💡 EXPECTED PRIORITY:"
echo "1. ✅ postman.go (6 annotations) - BEST (dedicated file)"
echo "2. ✅ api/routes.go (4 annotations) - GOOD (dedicated API)"  
echo "3. ❌ main.go (0 annotations) - SKIP (business logic mixed)"

echo -e "\n🎯 This approach separates concerns:"
echo "- main.go = App business logic"
echo "- postman.go = API documentation"
echo "- api/routes.go = API definitions"

echo -e "\n✅ MUCH CLEANER & MORE FLEXIBLE!"

# Cleanup
cd ..
rm -rf test-flexible