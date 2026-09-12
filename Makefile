APP := tmz
SOURCES := $(shell find . -name '*.go' | grep -v /vendor/)

$(APP): $(SOURCES) go.mod
	go fmt *.go
	go build -o $(APP)

dist/$(APP).darwin: export GOOS = darwin
dist/$(APP).darwin: $(SOURCES) go.mod
	go build -o dist/$(APP).darwin

dist/$(APP).exe: export GOOS = windows
dist/$(APP).exe: $(SOURCES) go.mod
	go build -o dist/$(APP).exe

dist/$(APP).linux.amd64: export GOOS = linux
dist/$(APP).linux.amd64: export GOARCH = amd64
dist/$(APP).linux.amd64: $(SOURCES) go.mod
	go build  -o dist/$(APP).linux.amd64

dist/$(APP).linux.arm64: export GOOS = linux
dist/$(APP).linux.arm64: export GOARCH = arm64
dist/$(APP).linux.arm64: $(SOURCES) go.mod
	go build  -o dist/$(APP).linux.arm64

.PHONY: dist
dist: dist/$(APP).darwin dist/$(APP).exe dist/$(APP).linux.amd64 dist/$(APP).linux.arm64

# this is too much of a pain to do inline
.PHONY: install
install: dist
	@./install.sh

.PHONY: test
test: $(SOURCES) go.mod
	go test

.PHONY: clean

clean:
	rm -fr $(APP) dist/
