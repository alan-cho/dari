package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"dari/harness/agent"
	"dari/harness/config"
	"dari/harness/server"
)

const serverURLEnv = "SERVER_URL"

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	client, err := agent.CreateClient()
	if err != nil {
		return err
	}
	serverURL, err := config.RequireURL(serverURLEnv)
	if err != nil {
		return err
	}
	listener, err := net.Listen("tcp", serverURL.Host)
	if err != nil {
		return err
	}
	log.Printf("listening on %s", listener.Addr())
	return http.Serve(listener, server.NewHandler(client))
}
