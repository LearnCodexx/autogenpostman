package main

import (
	"fmt"

	postmangen "github.com/learncodexx/autogenpostman"
)

func main() {
    fmt.Println("🎯 TESTING FLEXIBLE FILE DETECTION")
    
    if output, err := postmangen.EasyGenerate(); err == nil {
        fmt.Printf("✅ Success! Collection: %s\n", output)
    } else {
        fmt.Printf("ℹ️  EasyGenerate result: %v\n", err)
    }
}