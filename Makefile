VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS = -X main.version=$(VERSION)
BINARY = godo
DIST = dist

.PHONY: build test assemble clean dev

build: $(BINARY) assemble

$(BINARY):
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) ./cmd/godo/

test:
	go test ./...

assemble: $(BINARY)
	rm -rf $(DIST)
	./$(BINARY) assemble --core ./core --persona ./personas/do/manifest.yaml --out $(DIST)

clean:
	rm -f $(BINARY)
	rm -rf $(DIST)

dev: build
	rm -rf .claude
	mkdir -p .claude
	@# Copy all assembled dirs/files into .claude/
	@for item in agents commands rules skills styles characters spinners; do \
		[ -d $(DIST)/$$item ] && cp -r $(DIST)/$$item .claude/ || true; \
	done
	@[ -f $(DIST)/settings.json ] && cp $(DIST)/settings.json .claude/ || true
	@[ -f $(DIST)/registry.yaml ] && cp $(DIST)/registry.yaml .claude/ || true
	@# CLAUDE.md goes to project root
	@[ -f $(DIST)/CLAUDE.md ] && cp $(DIST)/CLAUDE.md ./CLAUDE.md || true
	@echo "Dev environment ready: .claude/ + CLAUDE.md"
