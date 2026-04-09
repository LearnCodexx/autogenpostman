package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	postmangen "github.com/learncodexx/autogenpostman"
)

func main() {
	var (
		swaggerInput = flag.String("swagger-input", "", "optional openapi/swagger file path; when empty will use auto mode")
		output       = flag.String("output", "docs/postman_collection.json", "postman collection output path")
		collection   = flag.String("collection-name", "WULE", "postman collection name")
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
		log.Fatalf("generate postman failed: %v", err)
	}

	fmt.Printf("Postman collection generated: %s\n", *output)
}
