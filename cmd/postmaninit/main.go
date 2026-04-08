package main

import (
	"flag"
	"fmt"
	"log"

	postmangen "learncodexx/point_of_sale/generate_postman_file"
)

func main() {
	var (
		workingDir  = flag.String("working-dir", ".", "project root directory")
		commandPath = flag.String("command-path", "cmd/postman/main.go", "target path for generated command")
		importPath  = flag.String("import-path", "learncodexx/point_of_sale/generate_postman_file", "go import path for generator package")
		outputPath  = flag.String("output", "docs/postman_collection.json", "default postman collection output path")
		collection  = flag.String("collection-name", "API", "default collection name")
		force       = flag.Bool("force", false, "overwrite file if already exists")
	)
	flag.Parse()

	out, err := postmangen.EnsurePostmanCommand(postmangen.ScaffoldConfig{
		WorkingDir:          *workingDir,
		CommandPath:         *commandPath,
		GeneratorImportPath: *importPath,
		CollectionName:      *collection,
		OutputPath:          *outputPath,
		Force:               *force,
	})
	if err != nil {
		log.Fatalf("generate postman command failed: %v", err)
	}

	fmt.Printf("Postman command scaffold ready: %s\n", out)
}
