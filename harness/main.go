package main

import (
	"os"

	"github.com/anthropics/anthropic-sdk-go/v2"
)


func main() {
	client := anthropic.NewClient()
	model := os.Getenv("MODEL")

	/*
	While Loop
	1. Take user input
	2. Add to transcript
	3. Send to API
	4. Check stop_reason
	5. Determine next steps based on stop_reason
	*/
	for {

	}
}
