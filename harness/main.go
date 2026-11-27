package main

import (
	"context"
	"fmt"
	"os"

	"dari/harness/agent"
	"dari/harness/agent/parser"
	herrors "dari/harness/errors"

	"github.com/anthropics/anthropic-sdk-go/v2"
)

const prompt = "> "

func main() {
	if err := run(context.Background()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	client, err := agent.CreateClient()
	if err != nil {
		return err
	}
	scanner := parser.InitScanner()
	var transcript []anthropic.MessageParam
	for {
		fmt.Print(prompt)
		input, err := parser.ParseInput(scanner)
		if herrors.HasCode(err, herrors.EndOfInput) {
			fmt.Println()
			return nil
		}
		if err != nil {
			return err
		}
		if input == "" {
			continue
		}

		updated, reply, err := agent.Send(ctx, client, transcript, input)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		transcript = updated
		fmt.Println(reply)
	}
}
