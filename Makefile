OLLAMA_URL := http://localhost:11434
OLLAMA_LOG := .ollama.log
OLLAMA_PID := .ollama.pid
MODEL ?= $(shell sed -n 's/^MODEL=//p' .env 2>/dev/null)

.PHONY: run serve stop require-model

run: require-model serve
	@ollama show $(MODEL) >/dev/null 2>&1 || ollama pull $(MODEL)
	ANTHROPIC_BASE_URL=$(OLLAMA_URL) ANTHROPIC_API_KEY=ollama MODEL=$(MODEL) go run ./harness

serve:
	@curl -sf $(OLLAMA_URL)/api/version >/dev/null || { ollama serve >$(OLLAMA_LOG) 2>&1 & echo $$! >$(OLLAMA_PID); }
	@for attempt in $$(seq 20); do curl -sf $(OLLAMA_URL)/api/version >/dev/null && exit 0; sleep 0.5; done; \
		echo "ollama server did not start, see $(OLLAMA_LOG)"; exit 1

stop:
	@test -f $(OLLAMA_PID) && kill $$(cat $(OLLAMA_PID)) && rm $(OLLAMA_PID) || echo "no ollama server started by make"

require-model:
	@test -n "$(MODEL)" || { echo "set MODEL in .env or run: make run MODEL=<ollama model name>"; exit 1; }
