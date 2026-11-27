# Plan: Error Library

Status: `DRAFT`

Date: 2026-11-27

## Goal

Define one error type for the harness so every package reports failures with a stable code, a human-readable description, and the operation where it happened, and callers branch on codes instead of message text.

## Out of Scope

- Stack traces
- Logging or telemetry
- Retry policies driven by error codes
- Translating every SDK error code one to one

## Context

- *Initial problem:* errors from `bufio`, `io`, and the Anthropic SDK reach callers in their original shapes, so `main` would need to know about each library to decide how to react.
- *Currently built:* `harness/agent/parser/input.go` prints read errors directly and returns raw `error` values. There is no shared error type.
- *Missing:* a shared error type, a fixed set of codes, and a convention for creating errors at package boundaries.

## Design

```
 bufio / io / SDK ──(1) returns raw error──► package boundary (parser, agent)
                                                │
                                                │ (2) translates into
                                                ▼
                                           errs.Error {Code, Description, Op, Cause}
                                                │
                                                │ (3) returned up to
                                                ▼
                                             main.go
                                                │ (4) branches on Code
                                                ▼
                                 exit / print and continue / warn
```

## Pattern

[Anti-Corruption Layer](https://learn.microsoft.com/en-us/azure/architecture/patterns/anti-corruption-layer): outside error types are translated into the harness's own vocabulary at the boundary where they enter, so callers depend only on `errs`.

## Accepted Decisions

- [2026-11-27] - An error has a Code, a Description, and an Op: `keep the shape minimal; codes for programs, descriptions for people, ops for locating the failure`
- [2026-11-27] - Codes are a named string type (`type Code string`): `codes stay readable when printed and are safe to add or reorder`
- [2026-11-27] - Descriptions are fixed per code in a single lookup table: `one source of truth keeps messages consistent across call sites`
- [2026-11-27] - The original error is kept as a `Cause` field and exposed through `Unwrap()`: `translate errors without swallowing them, so details like HTTP status remain available through errors.Is and errors.As`
- [2026-11-27] - Op is derived automatically from the caller of `New` via `runtime.Caller`, shortened to `package.Function`: `removes a hand-written op string from every function; op is a property of the code, not the request, so it does not belong in context.Context`

## Work Breakdown

| #   | Task                                              | Completed |
| --- | ------------------------------------------------- | --------- |
| 1   | Define the error type and code type               | `YES`     |
| 2   | Define the initial set of codes and descriptions  | `YES`     |
| 3   | Translate input reading errors in the parser      | `YES`     |
| 4   | Branch on codes in `main.go`                      | `YES`     |
