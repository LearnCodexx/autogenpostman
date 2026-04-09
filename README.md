# AutoGenPostman

Go package for automatically generating Postman collections from OpenAPI/Swagger with support for multiple project structures and **automatic route discovery**.

## 🎯 Why Use This?

- **🔍 Route Discovery**: Auto scan Go code for API routes (NEW!)
- **📄 No Manual Work**: Auto-convert Swagger → Postman
- **🏗️ Multi Structure**: Support various Go project layouts
- **🤖 Smart Detection**: Auto-find main.go in different locations
- **🚀 Production Ready**: Used in production applications
- **🔌 Easy Integration**: Plug & play to existing projects

## ✨ **NEW: Automatic Route Discovery**

Generate Postman collections directly from your Go source code **without needing swagger annotations!**

```bash
# Auto-discover routes from your Go project
go run cmd/route-discovery/main.go -discovery=true

# Scan specific project
go run cmd/route-discovery/main.go -discovery=true -project=./my-api
```

**Supported Frameworks**: Fiber, Gin, Echo, Gorilla Mux, and more!

👉 **[Full Route Discovery Documentation](ROUTE_DISCOVERY.md)**

## 🚨 WARNING: Common Issues & Solutions

### ❌ **"inconsistent vendoring" Error**

```bash
# SOLUTION 1: Sync vendor
go mod tidy && go mod vendor

# SOLUTION 2: Skip vendor (recommended)
export GOFLAGS=-mod=mod
go run cmd/postman/main.go

# SOLUTION 3: Delete vendor folder
rm -rf vendor/
go mod tidy
```

### ❌ **"swag command not found"**

```bash
# Install swag
go install github.com/swaggo/swag/cmd/swag@latest

# Add to PATH (add to ~/.bashrc or ~/.zshrc)
export PATH=$PATH:$(go env GOPATH)/bin

# Verify installation
swag --version
```

### ❌ **"npx command not found"**

```bash
# Ubuntu/Debian
sudo apt update && sudo apt install nodejs npm

# macOS
brew install node

# CentOS/RHEL
sudo yum install nodejs npm

# Verify
node --version && npm --version
```

### ❌ **"no Go files found"**

```bash
# Specify correct main.go location
go run cmd/postman/main.go --main-file="path/to/your/actual/main.go"

# Common locations:
# --main-file="main.go"              # Root
# --main-file="cmd/api/main.go"      # API service
# --main-file="cmd/server/main.go"   # Server
# --main-file="app/main.go"          # App folder
```

### ❌ **"cannot find package docs"**

```go
// Add this import to your main.go
import _ "yourapp/docs"  // Replace 'yourapp' with your module name
```

### ❌ **Empty Postman Collection Generated**

Your main.go needs proper swagger annotations:

```go
// @title           Your API Name          // REQUIRED
// @version         1.0                   // REQUIRED
// @description     API description       // REQUIRED
// @host            localhost:8080        // REQUIRED
// @BasePath        /api/v1              // REQUIRED

// For each endpoint:
// @Summary      Description
// @Tags         category
// @Success      200  {object}  ResponseType
// @Router       /path [method]           // REQUIRED
```

## 📦 Installation & Setup

### Step 1: Prerequisites (CRUCIAL!)

```bash
# 1. Install Go tools
go install github.com/swaggo/swag/cmd/swag@latest

# 2. Install Node.js (for OpenAPI → Postman conversion)
# Ubuntu/Debian: sudo apt install nodejs npm
# macOS: brew install node
# Windows: Download from nodejs.org

# 3. Fix vendor issues (if any)
export GOFLAGS=-mod=mod  # Add to ~/.bashrc for permanent fix
```

### Step 2: Add to Your Project

```bash
# In your Go project directory
go get github.com/learncodexx/autogenpostman@latest
```

### Step 3: Choose Setup Method

#### 🚀 Method A: One-Liner (Fastest)

```bash
go run github.com/learncodexx/autogenpostman/cmd/postman@latest \
  -collection-name "My API" \
  -output "docs/collection.json"
```

#### 🔧 Method B: Custom Setup (Flexible)

```go
// cmd/setup-postman/main.go
package main

import (
    "log"
    postmangen "github.com/learncodexx/autogenpostman"
)

func main() {
    out, err := postmangen.EnsurePostmanCommand(postmangen.ScaffoldConfig{
        WorkingDir:          ".",
        CommandPath:         "cmd/gen-postman/main.go",    // Where to create generator
        GeneratorImportPath: "github.com/learncodexx/autogenpostman",
        CollectionName:      "My API",                     // Postman collection name
        OutputPath:          "api/postman_collection.json", // Output location
        Force:               true,                         // Overwrite if exists
    })
    if err != nil {
        log.Fatalf("❌ Setup failed: %v", err)
    }
    log.Printf("✅ Generator ready: %s", out)
}
```

Run setup:

```bash
go run cmd/setup-postman/main.go
```

### Step 4: Ensure Swagger Annotations

Your main.go MUST have proper annotations:

```go
package main

import (
    "github.com/gin-gonic/gin"
    _ "yourproject/docs"  // ⚠️ CRITICAL: Replace 'yourproject' with your module name
)

// ⚠️ CRITICAL: These annotations are REQUIRED
// @title           My API
// @version         1.0
// @description     My API description
// @termsOfService  http://swagger.io/terms/
// @contact.name    API Support
// @contact.email   support@mycompany.com
// @license.name    MIT
// @license.url     https://opensource.org/licenses/MIT
// @host            localhost:8080
// @BasePath        /api/v1

func main() {
    r := gin.Default()

    // ⚠️ CRITICAL: Each endpoint needs annotations
    // @Summary      Health check
    // @Description  Check if API is running
    // @Tags         health
    // @Accept       json
    // @Produce      json
    // @Success      200  {object}  map[string]string
    // @Router       /health [get]
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })

    // @Summary      Get users
    // @Description  Get all users
    // @Tags         users
    // @Accept       json
    // @Produce      json
    // @Success      200  {array}   User
    // @Failure      500  {object}  ErrorResponse
    // @Router       /users [get]
    r.GET("/api/v1/users", getUsers)

    r.Run(":8080")
}

// Define response types (REQUIRED for proper swagger generation)
type User struct {
    ID    int    `json:"id" example:"1"`
    Name  string `json:"name" example:"John Doe"`
    Email string `json:"email" example:"john@example.com"`
}

type ErrorResponse struct {
    Error string `json:"error" example:"Something went wrong"`
}
```

## 🏃‍♂️ Generate Postman Collection

### Auto Mode (Recommended)

```bash
# Let the tool auto-detect everything
go run cmd/gen-postman/main.go

# With custom collection name
go run cmd/gen-postman/main.go --collection-name="Production API"
```

### Manual Mode (If Auto Fails)

```bash
# Specify main.go location
go run cmd/gen-postman/main.go \
  --main-file="cmd/api/main.go" \
  --output="build/postman.json" \
  --collection-name="My API"

# Use existing OpenAPI file
go run cmd/gen-postman/main.go \
  --swagger-input="api/openapi.yaml" \
  --output="postman.json"
```

### Debug Mode (For Troubleshooting)

```bash
# See what the tool is doing
export DEBUG=1
go run cmd/gen-postman/main.go --main-file="your/main.go"
```

## 📁 Supported Project Structures

### ✅ Standard Go Project

```
myproject/
├── go.mod
├── cmd/
│   └── server/main.go           # ← Auto-detected
├── internal/
├── api/
└── docs/
    └── postman_collection.json  # ← Generated here
```

### ✅ Simple Web App

```
myproject/
├── go.mod
├── main.go                      # ← Auto-detected
├── handlers/
└── docs/
    └── postman_collection.json  # ← Generated here
```

### ✅ Microservice Pattern

```
myproject/
├── go.mod
├── cmd/
│   ├── api/main.go             # ← Auto-detected
│   └── worker/main.go
├── services/
├── api/
│   └── openapi.yml             # ← Auto-detected if exists
└── docs/                       # ← Generated here
```

### ✅ Monorepo Structure

```
company-apis/
├── user-service/
│   ├── go.mod
│   ├── cmd/main.go             # ← Use --main-file flag
│   └── docs/
├── order-service/
│   ├── go.mod
│   ├── app/main.go             # ← Use --main-file flag
│   └── docs/
```

## 🆘 Comprehensive Troubleshooting

### Problem: "exit status 1" with swag

**Diagnosis:**

```bash
# Test swag manually
swag init -g main.go
```

**Solutions:**

```bash
# 1. Fix import path in main.go
import _ "your-actual-module-name/docs"  # Check go.mod for correct name

# 2. Verify swagger annotations syntax
# Must have @title, @version, @host, @BasePath

# 3. Install latest swag
go install github.com/swaggo/swag/cmd/swag@latest

# 4. Use absolute path
go run cmd/gen-postman/main.go --main-file="/full/path/to/main.go"
```

### Problem: "Cannot parse source files"

**Solutions:**

```bash
# 1. Check file exists
ls -la your/main.go

# 2. Verify Go syntax
go build your/main.go

# 3. Use correct module path
cd /path/to/your/project
go run cmd/gen-postman/main.go --main-file="./main.go"
```

### Problem: Generated Collection is Empty

**Root Causes & Fixes:**

1. **Missing swagger annotations** → Add @Router tags to ALL endpoints
2. **Wrong import path** → Fix `_ "yourmod/docs"` import
3. **No @title/@host/@BasePath** → Add required top-level annotations
4. **Syntax errors in annotations** → Check @ symbols and proper format

### Problem: "Module not found" Errors

**Solutions:**

```bash
# 1. Update dependencies
go mod tidy
go mod download

# 2. Clear module cache
go clean -modcache
go mod download

# 3. Use replace directive (for local development)
# In go.mod:
# replace github.com/learncodexx/autogenpostman => /path/to/local/copy
```

## 🔍 Debug & Validation Steps

### 1. Test Swagger Generation Manually

```bash
# This should work first before using autogenpostman
swag init -g main.go -o docs/
ls docs/  # Should see: docs.go, swagger.json, swagger.yaml
```

### 2. Validate OpenAPI File

```bash
# Check if swagger.json is valid
cat docs/swagger.json | jq .  # Should be valid JSON
```

### 3. Test Postman Conversion Manually

```bash
# This should work if swagger.json is valid
npx openapi-to-postmanv2 -s docs/swagger.json -o test-collection.json
```

### 4. Check Generator Configuration

```bash
# See what the generator is trying to do
go run cmd/gen-postman/main.go --help
```

## 📊 Performance & Best Practices

### ⚡ Optimization Tips

```bash
# Use module mode to avoid vendor issues
export GOFLAGS=-mod=mod

# Cache swag binary location
export PATH=$PATH:$(go env GOPATH)/bin

# Pre-install dependencies in CI
RUN go install github.com/swaggo/swag/cmd/swag@latest
```

### 🔒 Production Usage

```yaml
# .github/workflows/api-docs.yml
name: Generate API Documentation
on: [push, pull_request]

jobs:
  generate-docs:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v4
        with:
          go-version: "1.21"
      - uses: actions/setup-node@v4
        with:
          node-version: "18"

      - name: Install dependencies
        run: |
          go install github.com/swaggo/swag/cmd/swag@latest
          go mod download

      - name: Generate documentation
        run: |
          export GOFLAGS=-mod=mod
          go run github.com/learncodexx/autogenpostman/cmd/postman@latest \
            -collection-name="${{ github.repository }} API"
          go run cmd/postman/main.go

      - name: Upload artifacts
        uses: actions/upload-artifact@v4
        with:
          name: api-documentation
          path: |
            docs/swagger.json
            docs/postman_collection.json
```

## 🧪 Testing Your Setup

### Quick Verification Script

```go
// test-setup.go
package main

import (
    "context"
    "log"
    postmangen "github.com/learncodexx/autogenpostman"
)

func main() {
    cfg := postmangen.AutoConfig{
        WorkingDir:       ".",
        MainFile:         "main.go",  // Adjust this
        OutputPath:       "test-collection.json",
        CollectionName:   "Test API",
        Pretty:           true,
    }

    if err := postmangen.GenerateAuto(context.Background(), cfg); err != nil {
        log.Printf("❌ Generation failed: %v", err)
        log.Println("💡 Try: go run cmd/gen-postman/main.go --main-file='path/to/your/main.go'")
        return
    }

    log.Println("✅ Success! Check test-collection.json")
}
```

Run test:

```bash
go run test-setup.go
```

## 🤝 Getting Help

### 🐛 Report Issues

- **GitHub**: [Create Issue](https://github.com/learncodexx/autogenpostman/issues)
- **Include**: Go version, Node.js version, error message, project structure
- **Sample**: Minimal reproducible example

### 📚 More Examples

- **Complex API**: See `example/` folder
- **Custom Structures**: See documentation above
- **CI/CD Integration**: See workflow examples

### 💬 Community

- **Discussions**: [GitHub Discussions](https://github.com/learncodexx/autogenpostman/discussions)
- **Stack Overflow**: Tag with `autogenpostman` `go` `swagger`

---

## 📄 License

MIT License - feel free to use in your projects!

## ⭐ Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 🙏 Acknowledgments

- [Swaggo](https://github.com/swaggo/swag) for Swagger generation
- [openapi-to-postmanv2](https://www.npmjs.com/package/openapi-to-postmanv2) for conversion
- All contributors and users of this project
  MIT License - See [LICENSE](LICENSE) file.

### One-call (recommended for use from other applications)

Use `GenerateAuto` if you want the package to handle everything automatically.

Automatic behavior:

- package will try to generate Swagger first (using `swag` binary if available, or fallback to `go run github.com/swaggo/swag/cmd/swag@latest`)
- if Swagger generation fails, package will fallback to existing OpenAPI files in `docs/`
- if neither exists, only then return error

```go
package main

import (
    "context"
    "log"

    postmangen "github.com/learncodexx/autogenpostman"
)

func main() {
    err := postmangen.GenerateAuto(context.Background(), postmangen.AutoConfig{
        WorkingDir:     ".",
        OutputPath:     "docs/postman_collection.json",
        CollectionName: "User Service API",
        Pretty:         true,
    })
    if err != nil {
        log.Fatal(err)
    }
}
```

If your OpenAPI is in a custom location, set `SwaggerInputPath`:

```go
err := postmangen.GenerateAuto(context.Background(), postmangen.AutoConfig{
    WorkingDir:       ".",
    SwaggerInputPath: "specs/openapi.yaml",
    OutputPath:       "docs/postman_collection.json",
})
```

### Advanced / low-level config

```go
package main

import (
    "context"
    "log"

    postmangen "github.com/learncodexx/autogenpostman"
)

func main() {
    err := postmangen.Generate(context.Background(), postmangen.Config{
        WorkingDir:     ".",
        OutputPath:     "docs/postman_collection.json",
        CollectionName: "User Service API",
        Pretty:         true,
        Swag: postmangen.SwagConfig{
            Enabled:         true,
            MainFile:        "cmd/main.go",
            OutputDir:       "docs",
            ParseDependency: true,
            ParseInternal:   true,
        },
        Postman: postmangen.PostmanConfig{
            Options: map[string]string{
                "folderStrategy": "Tags",
            },
        },
    })
    if err != nil {
        log.Fatal(err)
    }
}
```

## Mode Without Swag

If swagger file already exists, disable `Swag.Enabled` and set `SwaggerInputPath`.

```go
err := postmangen.Generate(context.Background(), postmangen.Config{
    WorkingDir:       ".",
    SwaggerInputPath: "docs/openapi.yaml",
    OutputPath:       "docs/postman_collection.json",
    Pretty:           true,
})
```

## Converter Options

`Postman.Options` will be passed to the converter's `-O` flag.
Example:

```go
Postman: postmangen.PostmanConfig{
    Options: map[string]string{
        "folderStrategy": "Tags",
        "requestNameSource": "Fallback",
    },
}
```

## Notes

If you want to use local converter binary (`openapi2postmanv2`) without `npx`:

```go
Postman: postmangen.PostmanConfig{
    UseLocalCLI: true,
    CLIPath:     "openapi2postmanv2",
}
```
