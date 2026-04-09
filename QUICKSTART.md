# 🚀 Quick Start Guide

## 5 Menit Setup AutoGenPostman

### ⚡ Super Quick (Kalau sudah ada swagger annotations)

```bash
# 1. Install dependencies  
go install github.com/swaggo/swag/cmd/swag@latest

# 2. Add package
go get github.com/learncodexx/autogenpostman@latest

# 3. One-liner setup & generate (alternative)
go run github.com/learncodexx/autogenpostman/cmd/postman@latest \
  -collection-name "Your API Name"

# 4. Generate postman
go run cmd/postman/main.go
```

Done! Check `docs/postman_collection.json` 🎉

---

### 📝 Kalau Belum Ada Swagger Annotations

**1. Install tools:**
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

**2. Add annotations ke main.go:**
```go
// main.go
package main

import (
    "github.com/gin-gonic/gin"
    _ "yourapp/docs"  // ADD THIS
)

// ADD THESE ANNOTATIONS:
// @title My API
// @version 1.0  
// @host localhost:8080
// @BasePath /api/v1

func main() {
    r := gin.Default()
    
    // ADD ANNOTATIONS TO EACH ROUTE:
    // @Summary Get users
    // @Tags users  
    // @Success 200 {array} map[string]interface{}
    // @Router /api/v1/users [get]
    r.GET("/api/v1/users", getUsers)
    
    r.Run(":8080")
}
```

**3. Generate postman:**
```bash
go get github.com/learncodexx/autogenpostman@latest
go run github.com/learncodexx/autogenpostman/cmd/postman@latest
go run cmd/postman-generator/main.go
```

**Result:** `docs/postman_collection.json` ✅

---

### 🆘 Fix Common Errors

**"inconsistent vendoring":**
```bash
go mod tidy && export GOFLAGS=-mod=mod
```

**"swag command not found":**
```bash
go install github.com/swaggo/swag/cmd/swag@latest
export PATH=$PATH:$(go env GOPATH)/bin
```

**"npx command not found":**
```bash
# Ubuntu: sudo apt install nodejs npm
# macOS: brew install node
```

**"cannot find main.go":**
```bash
go run cmd/postman-generator/main.go --main-file="path/to/your/main.go"
```

---

### 📁 Works With Any Structure

✅ `main.go` (root)  
✅ `cmd/main.go`  
✅ `cmd/api/main.go`  
✅ `cmd/server/main.go`  
✅ Custom path (via `--main-file` flag)

### 📄 Need More Help?

- **Full docs:** [README.md](README.md)
- **Examples:** [IMPROVEMENTS.md](IMPROVEMENTS.md)
- **Issues:** [GitHub](https://github.com/learncodexx/autogenpostman/issues)

---

⭐ **Pro Tip:** Kalau project sudah jalan tapi error, coba dulu:
```bash
go run cmd/postman-generator/main.go --swagger-input="docs/openapi.yaml"
```