package main

import (
	"context"
	"fmt"
	"os"

	"dari/harness/config"
	"dari/harness/repl"
)

const serverURLEnv = "SERVER_URL"

func main() {
	if err := run(context.Background()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	serverURL, err := config.RequireURL(serverURLEnv)
	if err != nil {
		return err
	}
	return repl.Run(ctx, serverURL.String())
}
