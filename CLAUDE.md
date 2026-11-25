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

- Run against Claude: `go run ./harness` (needs `ANTHROPIC_API_KEY`; model defaults to `claude-opus-5-5`, override with `DARI_MODEL`)
- Run against a local model (free, for testing): `ANTHROPIC_BASE_URL=http://localhost:11434 ANTHROPIC_API_KEY=ollama DARI_MODEL=<ollama-model> go run ./harness` — Ollama (`brew services start ollama`) serves an Anthropic-compatible `/v1/messages`, so the same SDK code works unchanged.
- Build: `go build -o dari ./harness`
- Test all: `go test ./...`; single test: `go test ./harness/tools -run TestName`
- Lint/format: `go vet ./...` and `gofmt -l .`

## Architecture

- `harness/main.go` — REPL: reads a line, calls `Agent.Send`, prints errors. Wires the client, tool registry, and system prompt together.
- `harness/agent` — the harness loop. `Agent` owns the full `messages` history (the API is stateless; everything is resent each call). `Send` appends the user turn, then loops: call the model → append `resp.ToParam()` → run every `tool_use` block → send all `tool_result`s back in **one** user message → repeat until `StopReason` is not `tool_use`. Capped by `maxIterations`.
- `harness/tools` — `Tool` interface (`Name`, `Definition`, `Run`) plus a `Registry` that dispatches by name and produces the `Tools` param in stable order (keeps the prompt prefix cacheable). Tool errors are returned to the model as `is_error` results, not treated as fatal. To add a tool: implement `Tool` and register it in `harness/main.go`.

Uses the Messages API directly (manual loop), not the SDK's `BetaToolRunner` — owning the loop is the point of the project. Model: `claude-opus-5-5`.
