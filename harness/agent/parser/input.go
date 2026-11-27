package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func InitScanner() *bufio.Scanner {
	return bufio.NewScanner(os.Stdin)
}

func ParseInput(scanner *bufio.Scanner) (input string, err error) {
	input = strings.TrimSpace(scanner.Text())
	if !scanner.Scan() {
		if err = scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "Read Error: ", err)
		}
	}
	return input, err
}
