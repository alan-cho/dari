OLLAMA_URL := http://localhost:11434
OLLAMA_LOG := .ollama.log
OLLAMA_PID := .ollama.pid
SERVER_BIN := .bin/server
LOG_DIR := .logs
SERVER_LOG := $(LOG_DIR)/.$(shell date +%Y%m%d-%H%M%S).server.log
SERVER_ENV := ANTHROPIC_BASE_URL=$(OLLAMA_URL) ANTHROPIC_API_KEY=ollama
-include .env
export

.PHONY: run\:cli server cli model serve stop require-model

run\:cli: model
	@go build -o $(SERVER_BIN) ./harness/cmd/server
	@mkdir -p $(LOG_DIR)
	@! curl -s -o /dev/null $(SERVER_URL) || { echo "$(SERVER_URL) is already in use, stop that process first"; exit 1; }
	@$(SERVER_ENV) $(SERVER_BIN) >$(SERVER_LOG) 2>&1 & \
		trap "kill $$! 2>/dev/null" EXIT; \
		for attempt in $$(seq 20); do curl -s -o /dev/null $(SERVER_URL) && break; sleep 0.25; done; \
		curl -s -o /dev/null $(SERVER_URL) || { echo "dari server did not start, see $(SERVER_LOG)"; exit 1; }; \
		go run ./harness/cmd/cli

server: model
	$(SERVER_ENV) go run ./harness/cmd/server

cli:
	go run ./harness/cmd/cli

model: require-model serve
	@ollama show $(MODEL) >/dev/null 2>&1 || ollama pull $(MODEL)

serve:
	@curl -sf $(OLLAMA_URL)/api/version >/dev/null || { ollama serve >$(OLLAMA_LOG) 2>&1 & echo $$! >$(OLLAMA_PID); }
	@for attempt in $$(seq 20); do curl -sf $(OLLAMA_URL)/api/version >/dev/null && exit 0; sleep 0.5; done; \
		echo "ollama server did not start, see $(OLLAMA_LOG)"; exit 1

stop:
	@test -f $(OLLAMA_PID) && kill $$(cat $(OLLAMA_PID)) && rm $(OLLAMA_PID) || echo "no ollama server started by make"

require-model:
	@test -n "$(MODEL)" || { echo "set MODEL in .env or run: make server MODEL=<ollama model name>"; exit 1; }
