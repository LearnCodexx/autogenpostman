package main

import (
	"fmt"
	"log"
	
	postmangen "github.com/learncodexx/autogenpostman"
)

func main() {
	fmt.Println("🚀 Super Simple Postman Generator")
	
	// Method 1: Easiest - auto-find swagger and generate
	if output, err := postmangen.EasyGenerate(); err == nil {
		fmt.Printf("✅ Success! Collection: %s\n", output)
		return
	}
	
	// Method 2: Simple - specify files  
	if err := postmangen.SimpleGenerate("docs/swagger.json", "my_collection.json"); err == nil {
		fmt.Println("✅ Success! Collection: my_collection.json")
		return
	}
	
	// If both fail, show help
	fmt.Println("❌ No swagger file found")
	fmt.Println("\n💡 Quick Setup:")
	fmt.Println("1. Add annotations to main.go:")
	fmt.Println("   // @title My API")
	fmt.Println("   // @version 1.0")
	fmt.Println("   // @host localhost:8080")
	fmt.Println("")
	fmt.Println("2. Generate swagger:")
	fmt.Println("   swag init")
	fmt.Println("")
	fmt.Println("3. Run this again!")
}