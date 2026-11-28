# CLAUDE.md

## Project

`dari` is a from-scratch **LLM agent harness written in Go**, built as a learning project to understand how harnesses like Claude Code work internally. The goal is understanding, not shipping: prefer explicit, readable code over abstractions, and do not pull in agent frameworks. Use only the official Anthropic Go SDK (`github.com/anthropics/anthropic-sdk-go`) for API calls.

## Code Style

Modular, pure functions are most preferred. Never comment, if they are not maintained they are misleading. Rather, the code must be explicit and clear to be self-documenting. Before implementing a new solution, consider existing patterns within the repository to either refactor or follow an intended design. 

## Collaboration Style

The user is learning by building. Do not uplift or build without interaction from the user.

Whenever there is an opportunity to implement or refactor code into a popular system design, then always propose relevant architectural designs when applicable. Prioritize systems that can be explained from reliable sources like: [https://learn.microsoft.com/en-us/azure/architecture/](https://learn.microsoft.com/en-us/azure/architecture/). The domain is irrelevant, if the design pattern is suitable for the specific use-case, propose it. This can range from APIs, databases, software design, etc.

## Project Structure

TBD

## Documentation

Plans will live with `/plans` directory. It contains architectural or implementation decisions over various features that are added to the repository. When implementing a feature, refer to the `/plans` directory to search for an already decided design. Especially if it is explicitly mentioned by the user. Once all the work listed under `WORK BREAKDOWN` are completed, then mark the plan's status as `COMPLETE`. 