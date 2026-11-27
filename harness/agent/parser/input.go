package parser

import (
	"bufio"
	"io"
	"os"
	"strings"

	herrors "dari/harness/errors"
)

func InitScanner() *bufio.Scanner {
	return bufio.NewScanner(os.Stdin)
}

func ParseInput(scanner *bufio.Scanner) (string, error) {
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return "", herrors.New(herrors.InputRead, err)
		}
		return "", herrors.New(herrors.EndOfInput, io.EOF)
	}
	return strings.TrimSpace(scanner.Text()), nil
}
