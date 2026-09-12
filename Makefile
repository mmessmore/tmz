APP := tmz
SOURCES := $(shell find . -name '*.go' | grep -v /vendor/)
VERSION := $(shell sed -n 's/.*version = "\(.*\)"/\1/p' info.go)

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

dist/$(APP)_$(VERSION)_%.deb: dist/$(APP).linux.%
	rm -rf dist/deb/$*
	mkdir -p dist/deb/$*/DEBIAN dist/deb/$*/usr/bin
	install -m 755 dist/$(APP).linux.$* dist/deb/$*/usr/bin/$(APP)
	printf 'Package: $(APP)\nVersion: $(VERSION)\nSection: utils\nPriority: optional\nArchitecture: $*\nMaintainer: Mike Messmore <mike@messmore.org>\nDescription: Thin tmux wrapper that also understands screen'"'"'s attach flags\n' > dist/deb/$*/DEBIAN/control
	dpkg-deb --build --root-owner-group dist/deb/$* $@
	rm -rf dist/deb/$*

.PHONY: deb
deb: dist/$(APP)_$(VERSION)_amd64.deb dist/$(APP)_$(VERSION)_arm64.deb
	rm -f dist/$(APP).linux.amd64 dist/$(APP).linux.arm64
	@rmdir dist/deb 2>/dev/null || true

.PHONY: test
test: $(SOURCES) go.mod
	go test

.PHONY: clean

clean:
	rm -fr $(APP) dist/
