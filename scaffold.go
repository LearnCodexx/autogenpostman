package generatepostmanfile

import (
	"fmt"
	"os"
	"path/filepath"
)

// ScaffoldConfig controls auto-generation for cmd/postman/main.go.
type ScaffoldConfig struct {
	WorkingDir          string
	CommandPath         string
	GeneratorImportPath string
	CollectionName      string
	OutputPath          string
	Force               bool
}

// EnsurePostmanCommand creates cmd/postman/main.go scaffold for one-call Postman generation.
func EnsurePostmanCommand(cfg ScaffoldConfig) (string, error) {
	workingDir, err := resolveWorkingDir(cfg.WorkingDir)
	if err != nil {
		return "", err
	}

	if cfg.CommandPath == "" {
		cfg.CommandPath = filepath.Join("cmd", "postman", "main.go")
	}
	if cfg.GeneratorImportPath == "" {
		cfg.GeneratorImportPath = "github.com/learncodexx/autogenpostman"
	}
	if cfg.CollectionName == "" {
		cfg.CollectionName = "API"
	}
	if cfg.OutputPath == "" {
		cfg.OutputPath = "docs/postman_collection.json"
	}

	target := cfg.CommandPath
	if !filepath.IsAbs(target) {
		target = filepath.Join(workingDir, target)
	}

	if !cfg.Force {
		if _, statErr := os.Stat(target); statErr == nil {
			return target, nil
		}
	}

	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return "", fmt.Errorf("create command directory: %w", err)
	}

	content := fmt.Sprintf(`package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	postmangen "%s"
)

func main() {
	var (
		swaggerInput = flag.String("swagger-input", "", "optional openapi/swagger file path; when empty will use auto mode")
		output       = flag.String("output", "%s", "postman collection output path")
		collection   = flag.String("collection-name", "%s", "postman collection name")
		pretty       = flag.Bool("pretty", true, "pretty print output json")
	)
	flag.Parse()

	cfg := postmangen.AutoConfig{
		WorkingDir:       ".",
		SwaggerInputPath: *swaggerInput,
		OutputPath:       *output,
		CollectionName:   *collection,
		Pretty:           *pretty,
		Postman: postmangen.PostmanConfig{
			Options: map[string]string{
				"folderStrategy": "Tags",
			},
		},
	}

	if err := postmangen.GenerateAuto(context.Background(), cfg); err != nil {
		log.Fatalf("generate postman failed: %%v", err)
	}

	fmt.Printf("Postman collection generated: %%s\n", *output)
}
`, cfg.GeneratorImportPath, cfg.OutputPath, cfg.CollectionName)

	if err := os.WriteFile(target, []byte(content), 0o644); err != nil {
		return "", fmt.Errorf("write scaffold file: %w", err)
	}

	return target, nil
}
