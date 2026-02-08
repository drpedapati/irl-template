# Minimal helpers for working with this static artifact.

SHELL := /bin/sh
.ONESHELL:

PORT ?= 8765
HOST ?= 127.0.0.1

.PHONY: help serve smoke
.PHONY: ips zt

help:
	@echo "Targets:"
	@echo "  make serve        Start a local web server (Ctrl-C to stop)"
	@echo "  make smoke        Start server briefly and verify index.html loads"
	@echo "  make ips          Print likely LAN IPs (best-effort)"
	@echo "  make zt           Print ZeroTier info (if zerotier-cli is installed)"
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

ips:
	@echo "Likely IPs (best-effort):"
	@echo "  localhost: 127.0.0.1"
	@ipconfig getifaddr en0 2>/dev/null | sed 's/^/  en0: /' || true
	@ipconfig getifaddr en1 2>/dev/null | sed 's/^/  en1: /' || true
	@ipconfig getifaddr bridge0 2>/dev/null | sed 's/^/  bridge0: /' || true
	@ipconfig getifaddr awdl0 2>/dev/null | sed 's/^/  awdl0: /' || true
	@echo "  (If these are blank, run: ifconfig | head -n 80)"

zt:
	@command -v zerotier-cli >/dev/null 2>&1 || { echo "zerotier-cli not found"; exit 0; }
	@echo "ZeroTier:"
	@zerotier-cli info || true
	@zerotier-cli listnetworks || true
