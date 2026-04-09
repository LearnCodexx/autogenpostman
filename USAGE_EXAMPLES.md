# Examples untuk menggunakan autogenpostman package

## 1. Setup di aplikasi target

### Option 1: Menggunakan scaffolding tool
```go
package main

import (
    "log"
    postmangen "github.com/learncodexx/autogenpostman"
)

func main() {
    out, err := postmangen.EnsurePostmanCommand(postmangen.ScaffoldConfig{
        WorkingDir:          ".",
        CommandPath:         "cmd/generate-postman/main.go", // Hindari konflik dengan cmd/postman/
        GeneratorImportPath: "github.com/learncodexx/autogenpostman",
        CollectionName:      "WuleAPI",
        OutputPath:          "cmd/postman/postman_collection.json",
        Force:               true,
    })
    if err != nil {
        log.Fatalf("generate postman command failed: %v", err)
    }
    log.Printf("Postman command scaffold ready: %s", out)
}
```

### Option 2: Direct usage (tanpa scaffold)
```go
package main

import (
    "context"
    "flag"
    "log"
    postmangen "github.com/learncodexx/autogenpostman"
)

func main() {
    var (
        swaggerInput = flag.String("swagger-input", "", "optional openapi/swagger file path")
        output       = flag.String("output", "cmd/postman/postman_collection.json", "output path")
        collection   = flag.String("collection-name", "WuleAPI", "collection name")
        mainFile     = flag.String("main-file", "cmd/server/main.go", "path to your main.go file")
        pretty       = flag.Bool("pretty", true, "pretty print output")
    )
    flag.Parse()

    cfg := postmangen.AutoConfig{
        WorkingDir:        ".",
        MainFile:          *mainFile,           // PENTING: sesuaikan dengan struktur aplikasi Anda
        SwaggerInputPath:  *swaggerInput,
        OutputPath:        *output,
        CollectionName:    *collection,
        Pretty:            *pretty,
        SwagOutputDir:     "cmd/postman",       // Output swagger ke cmd/postman
        Postman: postmangen.PostmanConfig{
            Options: map[string]string{
                "folderStrategy": "Tags",
            },
        },
    }

    if err := postmangen.GenerateAuto(context.Background(), cfg); err != nil {
        log.Fatalf("generate postman failed: %v", err)
    }

    log.Printf("Postman collection generated: %s", *output)
}
```

## 2. Pastikan main.go aplikasi Anda memiliki swagger annotations

```go
package main

import (
    "log"
    "github.com/gin-gonic/gin"
    _ "your-app/docs" // Akan di-generate oleh swag
)

// @title           Wule API
// @version         1.0
// @description     API untuk aplikasi Wule
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

func main() {
    r := gin.Default()
    
    // @Summary      Show a account
    // @Description  get string by ID
    // @Tags         accounts
    // @Accept       json
    // @Produce      json
    // @Param        id   path      int  true  "Account ID"
    // @Success      200  {object}  Account
    // @Failure      400  {object}  HTTPError
    // @Failure      404  {object}  HTTPError
    // @Failure      500  {object}  HTTPError
    // @Router       /accounts/{id} [get]
    r.GET("/api/v1/accounts/:id", GetAccount)
    
    log.Fatal(r.Run(":8080"))
}
```

## 3. Jalankan generator

```bash
go run your-generator-file.go --main-file=cmd/server/main.go
```

## Penyebab error Anda:
- Package mencoba mencari `cmd/main.go` di aplikasi target tapi tidak ada
- Path `/home/amazon/BRI/PROJECT/wule-api/cmd/main.go` tidak ditemukan
- Solusi: berikan parameter `--main-file` yang menunjuk ke file main.go yang benar di aplikasi Anda