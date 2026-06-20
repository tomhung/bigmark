BINARY := bigmark
PKG    := ./cmd/bigmark

.PHONY: build install test clean fmt vet

build:
	go build -o $(BINARY) $(PKG)

install:
	go install $(PKG)

fmt:
	go fmt ./...

vet:
	go vet ./...

test: build
	@echo "smoke test: each mode renders without error"
	@./$(BINARY) "TEST" "smoke" >/dev/null && echo "  tier1 ok"
	@./$(BINARY) -2 "TEST" >/dev/null && echo "  tier2 ok"
	@./$(BINARY) -3 "test" >/dev/null && echo "  tier3 ok"
	@./$(BINARY) -r "TT" >/dev/null && echo "  rotated ok"
	@./$(BINARY) --canvas --seed 1 "TT" >/dev/null && echo "  canvas ok"

clean:
	rm -f $(BINARY)
	rm -rf dist
