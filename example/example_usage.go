package main

import (
	"context"
	"fmt"
	
	postmangen "github.com/learncodexx/autogenpostman"
)

func main() {
	fmt.Println("🚀 AutoGenPostman - Ultra Accommodating Example")
	
	// Example 1: One-liner generation (handles everything automatically)
	fmt.Println("\n📦 Example 1: QuickGenerate (simplest)")
	if err := postmangen.QuickGenerate("."); err != nil {
		fmt.Printf("❌ QuickGenerate failed: %v\n", err)
	} else {
		fmt.Println("✅ QuickGenerate succeeded!")
	}
	
	// Example 2: With custom collection name
	fmt.Println("\n📦 Example 2: QuickGenerateNamed")
	if err := postmangen.QuickGenerateNamed(".", "My Awesome API"); err != nil {
		fmt.Printf("❌ QuickGenerateNamed failed: %v\n", err)
	} else {
		fmt.Println("✅ QuickGenerateNamed succeeded!")
	}
	
	// Example 3: Advanced configuration (for power users)
	fmt.Println("\n📦 Example 3: Advanced GenerateAuto")
	ctx := context.Background()
	cfg := postmangen.AutoConfig{
		WorkingDir:     ".",
		MainFile:       "",  // Auto-detect
		OutputPath:     "custom/my-collection.json",
		CollectionName: "Advanced API",
		Pretty:         true,
	}
	
	if err := postmangen.GenerateAuto(ctx, cfg); err != nil {
		fmt.Printf("❌ Advanced generation failed: %v\n", err)
	} else {
		fmt.Println("✅ Advanced generation succeeded!")
	}
	
	fmt.Println("\n🎉 All examples completed!")
	fmt.Println("\n💡 Package features:")
	fmt.Println("✅ Auto-detects main.go in 10+ locations") 
	fmt.Println("✅ Auto-creates basic main.go if missing")
	fmt.Println("✅ Lenient validation (warns instead of failing)")
	fmt.Println("✅ Auto-scaffolds project structure")
	fmt.Println("✅ One-command generation")
	fmt.Println("✅ Handles missing swagger annotations gracefully")
	fmt.Println("✅ Comprehensive error messages with solutions")
}