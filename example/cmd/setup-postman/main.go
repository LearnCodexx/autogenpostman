package main

import (
	"log"

	postmangen "github.com/learncodexx/autogenpostman"
)

func main() {
	out, err := postmangen.EnsurePostmanCommand(postmangen.ScaffoldConfig{
		WorkingDir:          ".",
		CommandPath:         "cmd/postman-generator/main.go",
		GeneratorImportPath: "github.com/learncodexx/autogenpostman",
		CollectionName:      "Example API",
		OutputPath:          "docs/postman_collection.json",
		Force:               true,
	})
	if err != nil {
		log.Fatalf("Setup postman generator failed: %v", err)
	}

	log.Printf("Postman generator ready: %s", out)
	log.Println("")
	log.Println("Next steps:")
	log.Println("1. go run cmd/postman-generator/main.go")
	log.Println("2. Check docs/postman_collection.json")
}