# Run the desktop shell and the Datastar server.
#
#   make run     start the server, then the shell; stop the server when the shell exits
#   make server  server only
#   make dev     shell only, aimed at that server
#   make build   production binary
#   make test    unit tests
#
# Overrides: make run TOKEN=secret ADDR=127.0.0.1:9090

ADDR  ?= 127.0.0.1:8080
TOKEN ?= dev
API   ?= http://$(ADDR)
WAILS ?= wails

.DEFAULT_GOAL := help

.PHONY: help run server dev build test

help:
	@printf '%s\n' \
		'make run     start the server and the shell' \
		'make server  start the server on $(ADDR)' \
		'make dev     start the shell against $(API)' \
		'make build   wails build' \
		'make test    go test' \
		'' \
		'TOKEN=$(TOKEN)  ADDR=$(ADDR)'

run:
	DESKTOP_ADDR='$(ADDR)' DESKTOP_TOKEN='$(TOKEN)' go run ./cmd/server & pid=$$!; \
	trap 'pkill -P $$pid 2>/dev/null || true; kill $$pid 2>/dev/null || true' EXIT INT TERM; \
	ready=0; i=0; \
	while [ $$i -lt 50 ]; do \
		if curl -s -o /dev/null --max-time 1 '$(API)/desktop/bootstrap'; then ready=1; break; fi; \
		i=$$((i + 1)); sleep 0.2; \
	done; \
	if [ $$ready -ne 1 ]; then echo 'server did not start on $(ADDR)' >&2; exit 1; fi; \
	DESKTOP_API_BASE='$(API)' DESKTOP_DEV_TOKEN='$(TOKEN)' $(WAILS) dev

server:
	DESKTOP_ADDR='$(ADDR)' DESKTOP_TOKEN='$(TOKEN)' go run ./cmd/server

dev:
	DESKTOP_API_BASE='$(API)' DESKTOP_DEV_TOKEN='$(TOKEN)' $(WAILS) dev

build:
	$(WAILS) build

test:
	go test ./internal/...
	go test -tags production ./internal/...
