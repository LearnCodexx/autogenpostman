# AutoGenPostman

Go package untuk otomatis generate Postman collection dari OpenAPI/Swagger.

## 🚀 Features

🔍 **Smart Path Detection**: Otomatis cari main.go di lokasi umum  
📁 **Flexible Structure**: Support berbagai struktur project  
⚙️ **Zero Config**: Langsung jalan dengan default settings  
🛠️ **Customizable**: Override path via command-line flags  
🔄 **Fallback Support**: Auto-generate swagger atau pakai OpenAPI yang ada  

## 📋 Prerequisites (WAJIB!)

### 1. Install Go Tools
```bash
# Install swag untuk generate swagger
go install github.com/swaggo/swag/cmd/swag@latest

# Verify installation
swag --version
```

### 2. Install Node.js & npm
```bash
# Ubuntu/Debian
sudo apt update && sudo apt install nodejs npm

# macOS  
brew install node

# Verify installation
node --version
npm --version
```

### 3. Fix Vendor Issues (Jika Ada)
Jika dapat error "inconsistent vendoring":
```bash
# Di project target Anda
go mod tidy
go mod vendor

# Atau skip vendor
export GOFLAGS=-mod=mod
```

## 📦 Installation & Setup

### Step 1: Install Package
Di project target Anda:
```bash
go get github.com/learncodexx/autogenpostman@latest
```

### Step 2: Setup Generator Command
Buat file untuk setup generator (cukup 1x saja):

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
        CommandPath:         "cmd/postman-generator/main.go",  // Output generator
        GeneratorImportPath: "github.com/learncodexx/autogenpostman", 
        CollectionName:      "My API",                        // Nama collection
        OutputPath:          "docs/postman_collection.json",  // Output postman
        Force:               true,
    })
    if err != nil {
        log.Fatalf("Setup failed: %v", err)
    }
    log.Printf("Generator ready: %s", out)
}
```

Jalankan setup:
```bash
go run cmd/setup-postman/main.go
```

### Step 3: Setup API Annotations
Pastikan main.go Anda punya swagger annotations:

```go
// main.go atau cmd/api/main.go
package main

import (
    "github.com/gin-gonic/gin"
    _ "yourapp/docs"  // Will be generated
)

// @title           Your API Name
// @version         1.0
// @description     API description here
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@yourcompany.com

// @host      localhost:8080
// @BasePath  /api/v1

func main() {
    r := gin.Default()
    
    // @Summary      Health check
    // @Description  Check if API is running
    // @Tags         health
    // @Accept       json
    // @Produce      json
    // @Success      200  {object}  map[string]string
    // @Router       /health [get]
    r.GET("/health", healthCheck)
    
    r.Run(":8080")
}
```

## 🏃‍♂️ Usage

### Method 1: Auto-Detection (Recommended)
```bash
# Generator otomatis cari main.go dan generate
go run cmd/postman-generator/main.go
```

### Method 2: Specify Custom Paths  
```bash
# Jika main.go di lokasi custom
go run cmd/postman-generator/main.go \
  --swagger-input="" \
  --main-file="path/to/your/main.go" \
  --output="custom/output.json" \
  --collection-name="Custom API"
```

### Method 3: Pakai OpenAPI File yang Ada
```bash
# Jika sudah punya openapi.yaml
go run cmd/postman-generator/main.go \
  --swagger-input="api/openapi.yaml" \
  --output="docs/postman.json"
```

## 📁 Supported Project Structures

### ✅ Standard Go Project
```
project/
├── cmd/
│   └── api/main.go          # ← Auto-detected
├── docs/
│   └── postman_collection.json  # ← Generated
```

### ✅ Simple Web App  
```
project/
├── main.go                  # ← Auto-detected
├── docs/
│   └── postman_collection.json  # ← Generated
```

### ✅ Microservice
```
project/
├── cmd/
│   └── server/main.go       # ← Auto-detected  
├── api/
│   └── openapi.yaml        # ← Auto-detected
└── docs/
    └── postman_collection.json  # ← Generated
```

## 🔧 Step-by-Step Example

### Contoh Project Baru:

1. **Init project:**
```bash
mkdir my-api && cd my-api
go mod init my-api
```

2. **Install dependencies:**
```bash
go get github.com/gin-gonic/gin
go get github.com/learncodexx/autogenpostman@latest
go install github.com/swaggo/swag/cmd/swag@latest
```

3. **Buat main.go dengan annotations:**
```go
// main.go
package main

import (
    "github.com/gin-gonic/gin"
    _ "my-api/docs"
)

// @title My API
// @version 1.0  
// @host localhost:8080
// @BasePath /api/v1
func main() {
    r := gin.Default()
    
    // @Summary Get users
    // @Tags users
    // @Success 200 {array} map[string]interface{}
    // @Router /api/v1/users [get]
    r.GET("/api/v1/users", func(c *gin.Context) {
        c.JSON(200, []map[string]interface{}{
            {"id": 1, "name": "John"},
        })
    })
    
    r.Run(":8080")
}
```

4. **Setup generator:**
```bash
mkdir -p cmd/setup-postman
# Copy setup code dari atas ke cmd/setup-postman/main.go 
go run cmd/setup-postman/main.go
```

5. **Generate postman:**
```bash
go run cmd/postman-generator/main.go
```

6. **Result:**
```
docs/
├── postman_collection.json  # ← Your Postman collection!
├── swagger.json            # ← Generated swagger
└── docs.go                 # ← Generated by swag
```

## 🆘 Troubleshooting

### Error: "inconsistent vendoring"
```bash
# Solution 1: Fix vendor
go mod tidy && go mod vendor

# Solution 2: Skip vendor  
go run -mod=mod cmd/postman-generator/main.go
```

### Error: "no Go files found"
```bash
# Specify correct main file location
go run cmd/postman-generator/main.go --main-file="your/main.go"
```

### Error: "swag command not found"
```bash
# Install swag
go install github.com/swaggo/swag/cmd/swag@latest

# Add to PATH (add to ~/.bashrc)
export PATH=$PATH:$(go env GOPATH)/bin
```

### Error: "npx command not found"  
```bash
# Install Node.js
# Ubuntu: sudo apt install nodejs npm
# macOS: brew install node
```

### Error: "cannot find package docs"
Tambahkan blank import ke main.go:
```go
import _ "yourapp/docs"
```

### Generate Empty Collection
Main.go perlu swagger annotations yang proper:
```go
// @title API Name        # ← WAJIB
// @version 1.0          # ← WAJIB  
// @host localhost:8080  # ← WAJIB
// @BasePath /api        # ← WAJIB

// @Summary endpoint description    # ← Per endpoint
// @Router /path [method]           # ← Per endpoint
```

## 📚 Advanced Usage

### Custom Configuration
```go
cfg := postmangen.AutoConfig{
    WorkingDir:        ".",
    MainFile:          "cmd/server/main.go",    # Custom main path
    SwaggerInputPath:  "api/spec.yaml",        # Custom OpenAPI
    OutputPath:        "build/postman.json",   # Custom output
    CollectionName:    "Production API",       # Custom name
    Pretty:            true,                   # Pretty JSON
    SwagOutputDir:     "docs",                 # Swagger output
    Postman: postmangen.PostmanConfig{
        Options: map[string]string{
            "folderStrategy": "Tags",           # Group by tags
        },
    },
}

postmangen.GenerateAuto(context.Background(), cfg)
```

### CI/CD Integration
```yaml
# .github/workflows/postman.yml
name: Generate Postman Collection
on: [push]
jobs:
  postman:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    - uses: actions/setup-go@v3
    - uses: actions/setup-node@v3
    - run: go install github.com/swaggo/swag/cmd/swag@latest  
    - run: go run cmd/postman-generator/main.go
    - uses: actions/upload-artifact@v3
      with:
        name: postman-collection
        path: docs/postman_collection.json
```

## 🤝 Contributing

Issues dan PR welcome di [GitHub repository](https://github.com/learncodexx/autogenpostman).

## 📄 License

MIT License

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
