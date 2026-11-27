package repl

import (
	"context"
	"fmt"
	"os"

	"dari/harness/api"
	herrors "dari/harness/errors"
	"dari/harness/repl/parser"
)

const prompt = "> "

func Run(ctx context.Context, baseURL string) error {
	scanner := parser.InitScanner()
	var transcript []api.Message
	var previousTotal int64
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

		response, err := api.Send(ctx, baseURL, api.SendRequest{Messages: transcript, Input: input})
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		transcript = response.Messages
		fmt.Println(response.Reply)
		fmt.Println(formatUsage(messageInput(response.Usage, previousTotal), response.Usage.Output, total(response.Usage)))
		previousTotal = total(response.Usage)
	}
}
