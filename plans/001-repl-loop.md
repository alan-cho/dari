# Plan: REPL Loop

Status: `DRAFT`

Date: 2026-11-26

## Goal

Run a terminal chat in which each user line is sent to the model with the full conversation so far, and the model's reply is printed back, so the model remembers earlier turns.

## Out of Scope

- Tools and tool-call loops (`tool_use` handling beyond rejecting it)
- Streaming responses
- Persisting conversations across process restarts
- Context compaction or truncation
- System prompt tuning, thinking, and effort configuration
- Retry logic beyond the SDK's built-in retries

## Context

- *Initial problem:* `dari` needs a working conversation with the model before any agent behavior (tools, loops over tool results) can be layered on top.
- *Currently built:* `harness/main.go` constructs an `anthropic.Client` and reads the model ID from the `MODEL` environment variable. It does not compile yet (both values are unused). `harness/agent/` and `harness/tools/` exist but are empty.
- *Missing:* reading user input in a loop, holding the conversation transcript across turns, sending it to the model, handling the response's stop reason, and printing the reply.

## Design

```
 User ──(1) types a line──► REPL (main.go)
                              │
                              │ (2) passes input to
                              ▼
                         Conversation
                              │ (3) appends user message to transcript
                              │ (4) sends transcript to
                              ▼
                         Anthropic API (Messages.New)
                              │
                              │ (5) returns reply + stop_reason to
                              ▼
                         Conversation
                              │ (6) checks stop_reason
                              │ (7) appends reply to transcript
                              │ (8) returns reply text or error to
                              ▼
 User ◄──(9) prints reply── REPL (main.go)
                              │
                              └── (10) waits for next line, or exits on EOF
```

## Accepted Decisions

- [2026-11-26] - Start with the REPL loop only, without tools: `keep the first milestone small enough to understand end to end before adding agent behavior`
- [2026-11-26] - The user handwrites the implementation and Claude reviews it: `the project's goal is learning how a harness works, not shipping code`
- [2026-11-26] - The terminal input reader is separate from SDK conversion and returns a plain string: `keep terminal I/O independent of SDK types; slash commands and richer input types are deferred until needed`
- [2026-11-26] - Defer deciding which package owns SDK conversion: `decide when the conversion logic exists and its shape is known`
- [2026-11-27] - Claude writes the implementation instead of reviewing handwritten code: `the user changed their mind about the collaboration style`
- [2026-11-27] - The SDK call and transcript live in `agent` as pure functions, and `main` owns the transcript slice: `resolves the deferred ownership decision; keeps parser free of SDK types and keeps conversation state explicit`
- [2026-11-27] - A failed turn returns the previous transcript unchanged: `a turn either lands in the transcript completely or not at all, so a retry never leaves a dangling user message`

## Work Breakdown

| #   | Task                                                    | Completed |
| --- | ------------------------------------------------------- | --------- |
| 1   | Read user input line by line, exit cleanly on EOF       | `YES`     |
| 2   | Hold the conversation transcript across turns           | `YES`     |
| 3   | Send the transcript to the model and receive the reply  | `YES`     |
| 4   | Handle stop reasons and API errors                      | `YES`     |
| 5   | Print the reply and return to the prompt                | `YES`     |
