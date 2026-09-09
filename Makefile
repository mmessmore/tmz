APP     := tmz
SOURCES := $(shell find . -name '*.go' | grep -v /vendor/)
DEST    ?= $(shell for d in "$$HOME/bin" "$$HOME/.local/bin" /usr/local/bin; do [ -d "$$d" ] && echo "$$d" && break; done)

# Static, reproducible binaries for anything we ship.
GO_DIST := CGO_ENABLED=0 go build -trimpath -ldflags "-s -w"

.DELETE_ON_ERROR:
.PHONY: all build dist install test lint fmt clean

all: build

build: $(APP)

$(APP): $(SOURCES) go.mod $(wildcard go.sum)
	go build -trimpath -o $(APP) .

dist/$(APP).darwin.arm64: export GOOS = darwin
dist/$(APP).darwin.arm64: export GOARCH = arm64
dist/$(APP).darwin.amd64: export GOOS = darwin
dist/$(APP).darwin.amd64: export GOARCH = amd64
dist/$(APP).linux.amd64:  export GOOS = linux
dist/$(APP).linux.amd64:  export GOARCH = amd64
dist/$(APP).linux.arm64:  export GOOS = linux
dist/$(APP).linux.arm64:  export GOARCH = arm64
dist/$(APP).windows.amd64.exe: export GOOS = windows
dist/$(APP).windows.amd64.exe: export GOARCH = amd64

dist/$(APP).%: $(SOURCES) go.mod $(wildcard go.sum)
	@mkdir -p $(@D)
	$(GO_DIST) -o $@ .

DIST := \
	dist/$(APP).darwin.arm64 \
	dist/$(APP).darwin.amd64 \
	dist/$(APP).linux.amd64 \
	dist/$(APP).linux.arm64 \
	dist/$(APP).windows.amd64.exe

dist: $(DIST)

install: $(APP)
	@test -n "$(DEST)" || { echo "ERROR: no install dir found; pass DEST=<dir>"; exit 1; }
	@test -d "$(DEST)" || { echo "ERROR: destination $(DEST) does not exist"; exit 1; }
	install -m 755 $(APP) "$(DEST)"
	@echo "installed $(APP) to $(DEST)"

test: $(SOURCES) go.mod $(wildcard go.sum)
	go test ./...

lint:
	@test -z "$$(gofmt -l .)" || { echo "unformatted files:"; gofmt -l .; exit 1; }
	go vet ./...

fmt:
	go fmt ./...

clean:
	rm -rf $(APP) dist/
