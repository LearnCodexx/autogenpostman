# generate_postman_file

Package Go untuk generate Postman collection otomatis dari OpenAPI/Swagger.

## Pendekatan

Package ini tidak melakukan mapping OpenAPI ke Postman secara custom.
Ia hanya menjadi orchestration layer untuk tool yang sudah mature:

- `swag` untuk generate swagger docs (opsional)
- `openapi-to-postmanv2` (via `npx`) untuk convert ke Postman Collection v2.1

Dengan ini, package lebih stabil untuk dipakai di banyak aplikasi.

## Install

Jika module ini dipublish di git, dari project lain:

```bash
go get learncodexx/point_of_sale/generate_postman_file@latest
```

## Auto Create `cmd/postman/main.go` (setelah `go get`)

Go tidak menyediakan hook otomatis yang berjalan tepat setelah `go get`.
Sebagai gantinya, jalankan 1 command ini setelah install package:

```bash
go run learncodexx/point_of_sale/generate_postman_file/cmd/postmaninit@latest \
  -working-dir . \
  -command-path cmd/postman/main.go \
  -collection-name "User Service API"
```

Setelah itu file `cmd/postman/main.go` akan dibuat otomatis.

## Prasyarat Runtime

- Go 1.24+
- `swag` tersedia di PATH (kalau `Swag.Enabled=true`)
- Node.js + `npx` (default converter mode)

Install `swag`:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

## Contoh Pemakaian

### One-call (direkomendasikan untuk dipakai dari aplikasi lain)

Gunakan `GenerateAuto` jika Anda ingin package yang mengatur semuanya.

Perilaku otomatis:

- package akan mencoba generate Swagger dulu (pakai `swag` binary kalau ada, atau fallback `go run github.com/swaggo/swag/cmd/swag@latest`)
- jika generate Swagger gagal, package akan fallback ke file OpenAPI yang sudah ada di `docs/`
- kalau keduanya tidak ada, baru return error

```go
package main

import (
    "context"
    "log"

    postmangen "learncodexx/point_of_sale/generate_postman_file"
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

Jika OpenAPI Anda ada di lokasi custom, isi `SwaggerInputPath`:

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

    postmangen "learncodexx/point_of_sale/generate_postman_file"
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

## Mode Tanpa Swag

Kalau swagger file sudah ada, nonaktifkan `Swag.Enabled` dan isi `SwaggerInputPath`.

```go
err := postmangen.Generate(context.Background(), postmangen.Config{
    WorkingDir:       ".",
    SwaggerInputPath: "docs/openapi.yaml",
    OutputPath:       "docs/postman_collection.json",
    Pretty:           true,
})
```

## Opsi Converter

`Postman.Options` akan diteruskan ke flag `-O` converter.
Contoh:

```go
Postman: postmangen.PostmanConfig{
    Options: map[string]string{
        "folderStrategy": "Tags",
        "requestNameSource": "Fallback",
    },
}
```

## Catatan

Kalau ingin pakai binary converter lokal (`openapi2postmanv2`) tanpa `npx`:

```go
Postman: postmangen.PostmanConfig{
    UseLocalCLI: true,
    CLIPath:     "openapi2postmanv2",
}
```
