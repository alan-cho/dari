# Plan: Interface API

Status: `DRAFT`

Date: 2026-11-27

## Goal

Expose the agent as a stateless HTTP API and turn the terminal REPL into a client of it, so every frontend (CLI now, a webapp later) uses the same contract.

## Out of Scope

- Authentication and authorization
- Streaming responses (SSE or WebSockets)
- Server-side sessions or persistence
- Building the webapp
- Separate backends per frontend (Backends for Frontends)
- Tool use

## Context

- *Initial problem:* `harness/main.go` reads stdin and calls `agent.Send` in the same process, so the only way to use the agent is the terminal.
- *Currently built:* `agent.Send` is a pure function over `[]anthropic.MessageParam`, `parser` reads terminal input, and `main` owns the transcript. Plan 003 adds `agent.Usage`.
- *Missing:* wire types that don't depend on the SDK, an HTTP server around `agent.Send`, a CLI that talks to that server, and a mapping from error codes to HTTP statuses.

## Design

```
 CLI client ──(1) POST /messages {messages, input}──► HTTP server (harness/server)
                                                        │ (2) decodes api types and calls
                                                        ▼
                                                     agent.Send ──(3) sends transcript to──► Anthropic API
                                                        │ (4) returns transcript, reply, usage to
                                                        ▼
                                                   HTTP server
                                                        │ (5) encodes 200 {messages, reply, usage}
                                                        │     or 4xx/5xx {code, description}
                                                        ▼
 CLI client (harness/cli)
   │ (6) on success: replaces its transcript, adds usage to its total, prints reply and usage
   │     on failure: keeps both unchanged, prints the error
```

- `harness/api` holds the wire types: `Message{Role, Text}`, `SendRequest`, `SendResponse`, `Usage`, `ErrorResponse{Code, Description}`. The server and the client both import them.
- `agent` converts `[]api.Message` to and from `[]anthropic.MessageParam`, so SDK types never cross the wire.
- `harness/server` uses only `net/http`. A pure `statusFor(herrors.Code) int` maps codes to statuses, e.g. `MissingModel` → 500, `APIRequest` → 502, `Truncated`/`Refused` → 422, `BadRequest` → 400.
- `harness/cli` reuses `parser.InitScanner`/`ParseInput` and `formatUsage` from plan 003.
- The single `harness/main.go` becomes two binaries: `harness/cmd/server` and `harness/cmd/cli`. `SERVER_ADDR` sets the address and defaults to `localhost:8080`.
- New error codes: `ServerRequest` (the CLI cannot reach the server) and `BadRequest` (malformed request body).

## Pattern

- [Web API design best practices](https://learn.microsoft.com/en-us/azure/architecture/best-practices/api-design): resource naming, status codes, and a consistent error body.
- [Anti-Corruption Layer](https://learn.microsoft.com/en-us/azure/architecture/patterns/anti-corruption-layer): `api` types are translated into SDK types at the `agent` boundary.
- The stateless request and response shape mirrors the Anthropic Messages API. [Backends for Frontends](https://learn.microsoft.com/en-us/azure/architecture/patterns/backends-for-frontends) is the option to take if the CLI and the webapp ever need different contracts.

## Accepted Decisions

- [2026-11-27] - The agent runs as an HTTP server and the CLI is an HTTP client: `any frontend can use the same contract, and the client/server split is how real harnesses are structured`
- [2026-11-27] - The server is stateless and the client sends the full transcript on every request: `mirrors the Anthropic Messages API; no session storage to manage yet`
- [2026-11-27] - Wire types live in their own `api` package and are separate from SDK types: `the contract stays stable if the SDK or the model provider changes`
- [2026-11-27] - Error codes map to HTTP statuses through a single pure function: `one source of truth, like the descriptions table in plan 002`
- [2026-11-27] - Builds on plan 003, so responses carry `Usage`: `token counts reach every frontend`

## Work Breakdown

| #   | Task                                                      | Completed |
| --- | --------------------------------------------------------- | --------- |
| 1   | Define the `api` wire types                               | `NO`      |
| 2   | Convert between `api` messages and SDK messages in `agent` | `NO`      |
| 3   | Build the HTTP server and the error-to-status mapping     | `NO`      |
| 4   | Build the CLI client on top of the existing parser        | `NO`      |
| 5   | Split into server and CLI binaries and update the Makefile | `NO`      |
