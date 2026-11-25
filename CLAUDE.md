# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project purpose

`dari` is a from-scratch **LLM agent harness written in Go**, built as a learning project to understand how harnesses like Claude Code work internally. The goal is understanding, not shipping: prefer explicit, readable code over abstractions, and do not pull in agent frameworks. Use only the official Anthropic Go SDK (`github.com/anthropics/anthropic-sdk-go`) for API calls.

## Collaboration style

The user is learning by building. Scaffold routine code (module setup, types, API client wiring), but leave the core harness logic for the user to implement via a single `TODO(human)` marker — e.g. tool-call dispatch, appending tool results to the message history, loop exit conditions, permission checks, and context trimming.

## Planned build order

Each stage maps to one harness responsibility; build them incrementally:

1. Minimal agent loop: Messages API call → detect `tool_use` blocks → execute → append `tool_result` → repeat until the model replies with no tool calls. One tool (`read_file`) plus a REPL.
2. More tools: `bash`, `write_file`, `list_dir`.
3. Permission gate before risky tool calls (y/n prompt).
4. Streaming responses.
5. Context management: token counting, summarizing/truncating old messages.
6. Extras: system prompt loaded from a file, hooks, subagents, MCP client.

## Commands

Not yet scaffolded. Once a Go module exists, use the standard toolchain:

- Build: `go build ./...`
- Run: `go run .` (requires `ANTHROPIC_API_KEY` in the environment)
- Test all: `go test ./...`
- Single test: `go test ./path/to/pkg -run TestName`
- Lint/format: `go vet ./...` and `gofmt -l .`

Update this section and add an architecture section once the code layout exists.
