BINARY := applesoda
MODULE := github.com/pr0100pr0111/AppleSoda
VERSION := 0.1.0

GO := go
GOFLAGS := -trimpath
LDFLAGS := -s -w -X main.version=$(VERSION)

.PHONY: build build-nocgo tapspike run demo test vet fmt verify clean

build:
	CGO_ENABLED=1 $(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BINARY) ./cmd/applesoda

build-nocgo:
	CGO_ENABLED=0 $(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -tags nocgo -o $(BINARY) ./cmd/applesoda

tapspike:
	CGO_ENABLED=1 $(GO) build $(GOFLAGS) -o tapspike ./cmd/tapspike

run: build
	./$(BINARY)

demo: build-nocgo
	./$(BINARY) --demo

test:
	$(GO) test -race -count=1 ./...

vet:
	$(GO) vet ./...

fmt:
	gofmt -l -w .
	goimports -l -w . 2>/dev/null || true

verify: fmt vet test build build-nocgo
	@echo "All checks passed."

clean:
	rm -f $(BINARY) tapspike
	$(GO) clean -cache -testcache
