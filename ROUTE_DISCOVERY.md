# Route Discovery Feature

## Overview

Fitur **Route Discovery** memungkinkan `autogenpostman` untuk secara otomatis menemukan dan menganalisis routes dari source code Go Anda tanpa perlu swagger annotations manual.

## Cara Kerja

1. **🔍 Scan Source Code**: Mencari pattern routing di file `.go`
2. **📦 Extract Models**: Menemukan struct definitions untuk request/response
3. **📄 Generate Swagger**: Membuat OpenAPI spec dari discovered routes
4. **⚡ Convert to Postman**: Mengkonversi ke Postman collection

## Supported Frameworks

- ✅ **Fiber** (`app.Post("/path", handler)`)
- ✅ **Gin** (`router.GET("/path", handler)`)
- ✅ **Echo** (`e.POST("/path", handler)`)
- ✅ **Gorilla Mux** (`HandleFunc("/path", handler).Methods("POST")`)
- ✅ **Generic Groups** (`users.Post("/signin", handler)`)

## Usage Examples

### 1. Command Line Tool

```bash
# Run route discovery on current project
go run cmd/route-discovery/main.go -discovery=true

# Specify project path
go run cmd/route-discovery/main.go -discovery=true -project=./my-api

# Custom output and collection name
go run cmd/route-discovery/main.go -discovery=true \
  -output=api-collection.json \
  -collection-name="My API"
```

### 2. Programmatic Usage

```go
package main

import (
    "context"
    "github.com/learncodexx/autogenpostman"
)

func main() {
    cfg := autogenpostman.AutoConfig{
        WorkingDir:     ".",
        OutputPath:     "docs/collection.json",
        CollectionName: "Auto-Discovered API",
        Pretty:         true,
        RouteDiscovery: autogenpostman.RouteDiscoveryConfig{
            Enabled:     true,
            ProjectPath: "./my-api-project",
            IncludeAuth: true,
            TagStrategy: "path",
        },
    }

    err := autogenpostman.GenerateWithRouteDiscovery(context.Background(), cfg)
    if err != nil {
        panic(err)
    }
}
```

### 3. Integration with Existing Generator

```go
// Fallback: try route discovery if swagger fails
cfg := autogenpostman.AutoConfig{
    WorkingDir:     ".",
    OutputPath:     "collection.json",
    CollectionName: "API",
    RouteDiscovery: autogenpostman.RouteDiscoveryConfig{
        Enabled: true, // Enable as fallback
    },
}

// Try standard generation first
err := autogenpostman.GenerateAuto(context.Background(), cfg)
if err != nil {
    fmt.Println("⚠️  Standard generation failed, trying route discovery...")
    err = autogenpostman.GenerateWithRouteDiscovery(context.Background(), cfg)
}
```

## Configuration Options

### RouteDiscoveryConfig

```go
type RouteDiscoveryConfig struct {
    Enabled     bool   // Enable route discovery
    ProjectPath string // Path to scan (default: working dir)
    IncludeAuth bool   // Detect and include auth requirements
    TagStrategy string // How to generate tags: "path", "handler", "file"
}
```

### TagStrategy Options

- **"path"**: Generate tags from URL path (e.g., `/users/signin` → `Users`)
- **"handler"**: Generate tags from handler name (e.g., `UserHandler.SignIn` → `UserHandler`)
- **"file"**: Generate tags from file location (e.g., `handlers/user.go` → `User`)

## Route Patterns Detected

### Fiber Routes

```go
app.Post("/users/signin", ua.SignIn)
users := app.Group("/users")
users.Post("/signup", ua.SignUp)
```

### Gin Routes

```go
r.POST("/api/users", createUser)
router.GET("/api/users/:id", getUser)
```

### Echo Routes

```go
e.POST("/users", createUser)
echo.GET("/users/:id", getUser)
```

### Generic Groups

```go
users := app.Group("/users")
users.Post("/signin", handler)    // → POST /users/signin
users.Get("/profile", handler)    // → GET /users/profile
```

## Model Detection

Package secara otomatis mendeteksi struct definitions untuk API models:

```go
// ✅ Detected as API model (contains "Dto")
type UserDto struct {
    Name     string `json:"name" validate:"required"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=6"`
}

// ✅ Detected as API model (in dto/ folder)
type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}
```

### Model Detection Rules

Models are detected if:

- ✅ Struct name contains: `Dto`, `Request`, `Response`, `Model`, `Data`
- ✅ File is in directories: `dto/`, `model/`, `domain/`, `types/`
- ✅ Struct is public (starts with capital letter)

## Authentication Detection

Route discovery dapat mendeteksi requirements authentication:

```go
// ✅ Detected as authenticated route
userGet := users.Group("/get")
userGet.Use(JWTProtected)  // ← Auth middleware detected
userGet.Post("/all", ua.GetAllUser)

// ✅ Auth keywords detected in context
withAuth := users.Group("/")
withAuth.Use(JWTProtected)
withAuth.Post("/signout", ua.SignOut)
```

## Generated Output

### Example Swagger Generated

```json
{
  "swagger": "2.0",
  "info": {
    "title": "Auto-Discovered API",
    "version": "1.0.0",
    "description": "API automatically generated from 5 discovered routes"
  },
  "paths": {
    "/users/signin": {
      "post": {
        "summary": "Post Users",
        "description": "POST operation for /users/signin endpoint",
        "tags": ["Users"],
        "responses": {
          "200": { "description": "Success" },
          "400": { "description": "Bad Request" }
        }
      }
    }
  }
}
```

### Example Postman Collection

- 📁 **Users**
  - 📄 POST Sign in user (`/users/signin`)
  - 📄 POST User registration (`/users/signup`)
  - 📄 POST User sign out (`/users/signout`) 🔒
  - 📄 POST Get user by email (`/users/get/byemail`) 🔒
  - 📄 POST Get all users (`/users/get/all`) 🔒

## Benefits

### ✅ **No Annotations Required**

- Tidak perlu menambahkan `// @Summary`, `// @Router`, dll
- Bekerja dengan kode yang sudah ada

### ✅ **Framework Agnostic**

- Support multiple web frameworks
- Extensible untuk framework baru

### ✅ **Smart Detection**

- Auto-detect authentication requirements
- Group routes by logical structure
- Extract models from existing structs

### ✅ **Fallback Compatible**

- Bisa digunakan sebagai fallback jika swagger generation gagal
- Tidak mengganggu workflow yang sudah ada

## Troubleshooting

### No Routes Discovered

```bash
❌ Error: no routes discovered - ensure your Go files contain route definitions
```

**Solutions:**

1. Pastikan file menggunakan supported framework patterns
2. Check project path benar: `-project=./path/to/api`
3. Ensure routes defined in `.go` files (not `.md` or config)

### Missing Models

```bash
⚠️  No models found for request/response types
```

**Solutions:**

1. Ensure struct names follow convention (`UserDto`, `LoginRequest`)
2. Place models in dedicated folders (`dto/`, `model/`)
3. Make structs public (start with capital letter)

### Incomplete Route Information

```bash
⚠️ Route detected but missing handler details
```

**Solutions:**

1. Use clear handler references: `ua.SignIn` instead of anonymous functions
2. Ensure handlers are properly named and accessible

## Roadmap

### Planned Features

- [ ] **Smart Request/Response Inference**: Detect parameter types from handler signatures
- [ ] **Middleware Chain Analysis**: Extract middleware information and requirements
- [ ] **Comment Extraction**: Use comments near routes as descriptions
- [ ] **OpenAPI 3.0 Support**: Generate OpenAPI 3.0 instead of Swagger 2.0
- [ ] **Custom Templates**: Allow custom swagger generation templates
- [ ] **Interactive Mode**: CLI prompts for configuration options

### Framework Support

- [ ] **Chi Router**
- [ ] **FastHTTP**
- [ ] **Revel**
- [ ] **Beego**
- [ ] **Buffalo**
