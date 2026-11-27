package repl

import (
	"fmt"
	"strconv"
	"strings"

	"dari/harness/api"
)

func total(usage api.Usage) int64 {
	return usage.Prompt + usage.Output
}

func messageInput(turn api.Usage, previousTotal int64) int64 {
	return (total(turn) - previousTotal) - turn.Output
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
