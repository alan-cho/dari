package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

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

		updated, reply, turn, err := agent.Send(ctx, client, transcript, input)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		transcript = updated
		fmt.Println(reply)
		fmt.Println(formatUsage(agent.MessageInput(turn, previousTotal), turn.Output, agent.Total(turn)))
		previousTotal = agent.Total(turn)
	}
}

func formatUsage(input, output, total int64) string {
	return fmt.Sprintf("[input %s · output %s | total %s]", groupDigits(input), groupDigits(output), groupDigits(total))
}

func groupDigits(n int64) string {
	digits := strconv.FormatInt(n, 10)
	var grouped strings.Builder
	for i, digit := range digits {
		if i > 0 && (len(digits)-i)%3 == 0 {
			grouped.WriteByte(',')
		}
		grouped.WriteRune(digit)
	}
	return grouped.String()
}
