# Plan: Interface API

Status: `COMPLETE`

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
 CLI client (harness/repl)
   │ (6) on success: replaces its transcript, sets previousTotal, prints reply and usage
   │     on failure: keeps both unchanged, prints the error
```

- `harness/api` holds the wire types: `Message{Role, Text}`, `SendRequest`, `SendResponse`, `Usage`, `ErrorResponse{Code, Description}`. The server and the client both import them. `api.Send` is the HTTP client for that contract, and `api.MessagesPath` is shared with the server's route.
- `agent` converts `[]api.Message` to and from `[]anthropic.MessageParam`, so SDK types never cross the wire.
- `harness/server` uses only `net/http`. A pure `statusFor(herrors.Code) int` maps codes to statuses, e.g. `Internal` → 500, `APIRequest` → 502, `Truncated`/`Refused` → 422, `BadRequest` → 400.
- `harness/repl` reuses `parser.InitScanner`/`ParseInput` and `formatUsage` from plan 003. `parser` moves from `harness/agent/parser` to `harness/repl/parser`, since only the CLI reads the terminal.
- The single `harness/main.go` becomes two binaries: `harness/cmd/server` and `harness/cmd/cli`. Both read `SERVER_URL`: the server listens on its `host:port`, and the CLI sends requests to it with routes like `/messages` appended. `harness/config` reads required env vars; `MODEL`, `MAX_TOKENS` and `SERVER_URL` have no defaults. The Makefile loads and exports `.env`, and `.env.example` lists every variable. `make server` and `make cli` replace `make run`.
- New error codes: `ServerRequest` (the CLI cannot reach the server), `ServerResponse` (the server answered with an undecodable body), `BadRequest` (malformed request body) and `Internal` (the server hit an error without a harness code).
- `api.Usage` replaces `agent.Usage`. `Total` and `MessageInput` move to `harness/repl`, which owns `previousTotal`.

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
- [2026-11-27] - `api.Usage` replaces `agent.Usage`, and `Total`/`MessageInput` move to `harness/repl`: `the stateless server cannot know previousTotal, and the CLI never imports agent or the SDK`
- [2026-11-27] - The terminal client package is `repl`, not `cli`: `cmd/cli is the program and repl is the code it runs, so the names no longer look like duplicates`
- [2026-11-27] - `Send` lives in `api` next to the wire types: `the contract and the client that speaks it change together, like an SDK; the REPL only handles the terminal`
- [2026-11-27] - One `SERVER_URL` for both programs; the server listens on its host and port: `everything runs on one machine for now; split into a listen address and a base URL once a proxy, HTTPS or a container sits between client and server`
- [2026-11-27] - Required env vars have no defaults and the program fails at startup with `MissingEnv`/`InvalidEnv`, which replace `MissingModel`: `configuration is explicit, and a missing value fails loudly instead of silently using a guess`
- [2026-11-27] - `maxTokens` moves from a constant to `MAX_TOKENS`: `it is a model setting like MODEL`
- [2026-11-27] - Add a `ServerResponse` code for undecodable server replies: `keeps an unreachable server distinct from a garbled answer`

## Work Breakdown

| #   | Task                                                      | Completed |
| --- | --------------------------------------------------------- | --------- |
| 1   | Define the `api` wire types                               | `YES`     |
| 2   | Convert between `api` messages and SDK messages in `agent` | `YES`     |
| 3   | Build the HTTP server and the error-to-status mapping     | `YES`     |
| 4   | Build the CLI client on top of the existing parser        | `YES`     |
| 5   | Split into server and CLI binaries and update the Makefile | `YES`     |
