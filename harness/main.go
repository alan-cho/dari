package main

import (
	"os"

	"github.com/anthropics/anthropic-sdk-go/v2"
)


func main() {
	client := anthropic.NewClient()
	model := os.Getenv("MODEL")
}
