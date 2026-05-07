# golings — build / install / uninstall

BINARY  ?= golings
PREFIX  ?= $(HOME)/.local
BINDIR  ?= $(PREFIX)/bin
SHAREDIR?= $(PREFIX)/share/golings

GO ?= go

.PHONY: all build install uninstall reinstall verify run clean help

all: build

## Build the golings binary in the project root.
build:
	$(GO) build -o $(BINARY) ./cmd/golings

## Install the binary to $(BINDIR) and copy exercises/, solutions/, info.json
## to $(SHAREDIR). The installed binary will run from $(SHAREDIR) so `golings`
## works from any directory.
install: build
	@install -d $(BINDIR) $(SHAREDIR)
	@cp -r exercises solutions info.json $(SHAREDIR)/
	@install -m 0755 $(BINARY) $(SHAREDIR)/$(BINARY).bin
	@# Wrapper cd's into SHAREDIR so info.json + exercises/ resolve.
	@printf '#!/bin/sh\ncd "%s" && exec ./$(BINARY).bin "$$@"\n' "$(SHAREDIR)" > $(BINDIR)/$(BINARY)
	@chmod 0755 $(BINDIR)/$(BINARY)
	@echo "Installed $(BINARY) to $(BINDIR)/$(BINARY)"
	@echo "Data dir:  $(SHAREDIR)"
	@case ":$$PATH:" in *":$(BINDIR):"*) ;; *) \
	  echo; echo "Note: $(BINDIR) is not on your PATH."; \
	  echo "  fish:   fish_add_path $(BINDIR)"; \
	  echo "  bash:   echo 'export PATH=\"$(BINDIR):\$$PATH\"' >> ~/.bashrc"; \
	;; esac

## Remove the installed binary, wrapper, and data directory.
uninstall:
	@rm -f $(BINDIR)/$(BINARY)
	@rm -rf $(SHAREDIR)
	@echo "Removed $(BINDIR)/$(BINARY) and $(SHAREDIR)"

## Rebuild and reinstall.
reinstall: uninstall install

## Run `golings verify` against the working tree (does not require install).
verify: build
	./$(BINARY) verify

## Run `golings watch` against the working tree.
run: build
	./$(BINARY) watch

## Remove build artifacts from the working tree.
clean:
	@rm -f $(BINARY)
	@echo "Cleaned local build artifacts."

help:
	@echo "Targets:"
	@echo "  make build       — compile ./$(BINARY)"
	@echo "  make install     — install to $(BINDIR) (override with PREFIX=...)"
	@echo "  make uninstall   — remove installed files"
	@echo "  make reinstall   — uninstall + install"
	@echo "  make verify      — build and run all exercises once"
	@echo "  make run         — build and start the watch loop"
	@echo "  make clean       — remove the local binary"
