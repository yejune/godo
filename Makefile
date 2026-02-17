VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS = -X main.version=$(VERSION)
BINARY = godo
DIST = dist

.PHONY: build test assemble clean

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
