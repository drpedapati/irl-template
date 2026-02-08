# Minimal helpers for working with this static artifact.

SHELL := /bin/sh
.ONESHELL:

PORT ?= 8765
HOST ?= 127.0.0.1

.PHONY: help serve smoke

help:
	@echo "Targets:"
	@echo "  make serve        Start a local web server (Ctrl-C to stop)"
	@echo "  make smoke        Start server briefly and verify index.html loads"
	@echo ""
	@echo "Vars:"
	@echo "  PORT=8765 HOST=127.0.0.1"

serve:
	@echo "Serving on http://localhost:$(PORT) (from $$(pwd))"
	python3 -m http.server $(PORT)

smoke:
	set -eu
	python3 -m http.server $(PORT) >/tmp/irl-explainer-site-server.$(PORT).log 2>&1 &
	pid="$$!"
	trap 'kill "$$pid" 2>/dev/null || true' EXIT INT TERM
	sleep 0.6
	curl -fsS "http://$(HOST):$(PORT)/index.html" | grep -q "<title>IRL: Idempotent Research Loop</title>"
	echo "OK: index.html served and title matched"

