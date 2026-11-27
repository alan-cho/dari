# Plan: Token Usage

Status: `COMPLETE`

Date: 2026-11-27

## Goal

After every model reply, show the tokens of the user's message, the tokens of the reply, and the total size of the transcript.

## Out of Scope

- Displaying cache read and cache write tokens separately
- Cost in dollars
- Counting tokens before sending (`CountTokens`)
- Token limits, budgets, or warnings

## Context

- *Initial problem:* there is no way to see how large each message is or how fast the conversation's context is growing.
- *Currently built:* `agent.Send` in `harness/agent/conversation.go` receives an `anthropic.Message` and returns only the updated transcript and the reply text. `message.Usage` is thrown away.
- *Missing:* a harness-owned usage type, a way to derive per-message counts from it, and a printed usage line after each reply.

## Design

```
 Anthropic API / Ollama ──(1) returns Message with Usage──► agent.Send
                                                              │ (2) translates SDK Usage into
                                                              ▼
                                       agent.Usage {Prompt = input + cache read + cache write, Output}
                                                              │ (3) returned with transcript and reply to
                                                              ▼
                                                           main.go
                                                              │ (4) total = Total(turn)
                                                              │ (5) in = (total − previousTotal) − out
                                                              │ (6) previousTotal = total, on success only
                                                              ▼
 User ◄──(7) prints reply, then "[input 16 · output 4 | total 56]"
```

The API reports the size of the whole prompt, not of individual messages. The total is exact. "in" is derived from the difference between two totals, so it also includes chat-template tokens (role headers, end-of-turn markers) and any difference in how the previous reply re-tokenizes as input.

The prompt size adds the cache tokens because Ollama and Anthropic both report the cached prefix of the prompt in `cache_read_input_tokens` and leave it out of `input_tokens`.

## Pattern

[Anti-Corruption Layer](https://learn.microsoft.com/en-us/azure/architecture/patterns/anti-corruption-layer): the SDK's `anthropic.Usage` is translated into `agent.Usage` at the `agent` boundary, so `main` never depends on SDK types. This is the same approach plan 002 takes for errors.

## Accepted Decisions

- [2026-11-27] - Track only input and output tokens: `keep the first version minimal; cache tokens can be added once prompt caching is used`
- [2026-11-27] - Never calculate a dollar cost: `models usually run locally, where there is no per-token price, and token counts alone show context growth and speed`
- [2026-11-27] - Show per-turn usage plus a running session total: `per-turn shows the cost of one message, the total shows the cost of the conversation`
- [2026-11-27] - `agent.Usage` is a harness-owned type, built from `message.Usage`: `SDK types stay inside agent`
- [2026-11-27] - `main` owns the running total next to the transcript, and adds to it with a pure `AddUsage`: `conversation state stays explicit, the same as the transcript`
- [2026-11-27] - A failed turn leaves the total unchanged: `matches plan 001, where a turn either lands completely or not at all`
- [2026-11-27] - Replaces the running session total: "total" is the current transcript size (full prompt + reply), "in" is the user's message, "out" is the reply: `Ollama caches the prompt prefix, so summed input counts mean neither transcript size nor message size`
- [2026-11-27] - Prompt size is `input_tokens + cache_read_input_tokens + cache_creation_input_tokens`: `both Ollama and Anthropic leave the cached prefix out of input_tokens`
- [2026-11-27] - "in" is `(total − previousTotal) − out`, and `main` keeps `previousTotal` instead of a summed `Usage`, which removes `AddUsage`: `the API only reports whole-prompt sizes, so a message's size is the growth between totals`

## Work Breakdown

| #   | Task                                                        | Completed |
| --- | ----------------------------------------------------------- | --------- |
| 1   | Define `agent.Usage` and translate it from the SDK response | `YES`     |
| 2   | Return usage from `agent.Send`                              | `YES`     |
| 3   | Derive message input and total, keep `previousTotal` in `main` | `YES`  |
| 4   | Format and print the usage line after each reply            | `YES`     |
