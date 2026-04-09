package main

import (
	"context"
	"fmt"
	"log"
	"os"
	
	postmangen "github.com/learncodexx/autogenpostman"
)

func main() {
	fmt.Println("🚀 AutoGenPostman - SIMPLE Mode")
	
	// Step 1: Check if swagger.json exists
	if _, err := os.Stat("docs/swagger.json"); os.IsNotExist(err) {
		fmt.Println("❌ No swagger.json found. Please run: swag init")
		os.Exit(1)
	}
	
	// Step 2: Simple conversion (no auto-detection complexity)
	cfg := postmangen.AutoConfig{
		WorkingDir:       ".",
		SwaggerInputPath: "docs/swagger.json",  // Use existing swagger
		OutputPath:       "postman_collection.json",
		CollectionName:   "My API",
		Pretty:           true,
	}
	
	if err := postmangen.GenerateAuto(context.Background(), cfg); err != nil {
		log.Printf("❌ Failed: %v", err)
		fmt.Println("\n💡 Try:")
		fmt.Println("1. swag init -g main.go")
		fmt.Println("2. npx openapi-to-postmanv2 -s docs/swagger.json -o collection.json")
		os.Exit(1)
	}
	
	fmt.Println("✅ Success! Check postman_collection.json")
}