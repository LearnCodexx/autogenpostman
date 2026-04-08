# generate_postman_file

Go package for automatically generating Postman collections from OpenAPI/Swagger.

## Approach

This package does not perform custom OpenAPI to Postman mapping.
It only serves as an orchestration layer for mature tools:

- `swag` for generating swagger docs (optional)
- `openapi-to-postmanv2` (via `npx`) for converting to Postman Collection v2.1

This approach makes the package more stable for use across many applications.

## Installation

If this module is published to git, from other projects:

```bash
go get github.com/learncodexx/autogenpostman@latest
```

## Auto Create `cmd/postman/main.go` (after `go get`)

Go does not provide automatic hooks that run right after `go get`.
Instead, run this command after installing the package:

```bash
go run github.com/learncodexx/autogenpostman/cmd/postmaninit@latest \
  -working-dir . \
  -command-path cmd/postman/main.go \
  -collection-name "User Service API"
```

After that, the file `cmd/postman/main.go` will be created automatically.

## Runtime Requirements

- Go 1.24+
- `swag` available in PATH (if `Swag.Enabled=true`)
- Node.js + `npx` (default converter mode)

Install `swag`:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

## Usage Examples

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
