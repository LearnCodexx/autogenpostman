package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/learncodexx/autogenpostman"
)

func main() {
	var (
		projectPath    = flag.String("project", ".", "path to Go project for route scanning")
		output         = flag.String("output", "docs/collection.json", "postman collection output path")
		collectionName = flag.String("collection-name", "Auto-Discovered API", "postman collection name")
		pretty         = flag.Bool("pretty", true, "pretty print output json")
		useDiscovery   = flag.Bool("discovery", false, "use automatic route discovery instead of swagger")
	)
	flag.Parse()

	if *useDiscovery {
		// Use route discovery mode
		cfg := autogenpostman.AutoConfig{
			WorkingDir:     ".",
			OutputPath:     *output,
			CollectionName: *collectionName,
			Pretty:         *pretty,
			RouteDiscovery: autogenpostman.RouteDiscoveryConfig{
				Enabled:     true,
				ProjectPath: *projectPath,
				IncludeAuth: true,
				TagStrategy: "path", // "path", "handler", or "file"
			},
			SwagOutputDir: "docs",
		}

		fmt.Printf("🔍 Using automatic route discovery mode\n")
		fmt.Printf("📂 Scanning project: %s\n", *projectPath)

		if err := autogenpostman.GenerateWithRouteDiscovery(context.Background(), cfg); err != nil {
			log.Fatalf("Route discovery generation failed: %v", err)
		}
	} else {
		// Use standard swagger mode
		cfg := autogenpostman.AutoConfig{
			WorkingDir:     ".",
			OutputPath:     *output,
			CollectionName: *collectionName,
			Pretty:         *pretty,
		}

		fmt.Printf("📄 Using standard swagger mode\n")

		if err := autogenpostman.GenerateAuto(context.Background(), cfg); err != nil {
			log.Fatalf("Standard generation failed: %v", err)
		}
	}

	fmt.Printf("✅ Postman collection generated: %s\n", *output)
}