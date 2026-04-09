package main

import (
	"context"
	"fmt"

	postmangen "github.com/learncodexx/autogenpostman"
)

func main() {
    fmt.Println("🎯 TESTING FLEXIBLE FILE DETECTION")
    
    // Test the enhanced auto-detection directly
    cfg := postmangen.AutoConfig{
        WorkingDir:     ".",
        OutputPath:     "flexible_collection.json", 
        CollectionName: "Flexible API",
        Pretty:         true,
    }
    
    ctx := context.Background()
    if err := postmangen.GenerateAuto(ctx, cfg); err != nil {
        fmt.Printf("❌ Error: %v\n", err)
    } else {
        fmt.Printf("✅ Success! Check flexible_collection.json\n")
    }
}