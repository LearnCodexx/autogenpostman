package main

import (
	postmangen "github.com/learncodexx/autogenpostman"
)

func main() {
	// Test 1: Auto-scaffold when no main.go exists
	err := postmangen.QuickGenerate(".")
	if err != nil {
		println("Error:", err.Error())
	}
}