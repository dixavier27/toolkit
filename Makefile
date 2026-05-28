.PHONY: build test run install sync-tupa clean

# Versão injetada via ldflags. Default: git describe.
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)

build:
	go build -ldflags "-X main.version=$(VERSION)" -o bin/eco ./cmd/eco

test:
	go test ./...

run: build
	./bin/eco

install:
	go install -ldflags "-X main.version=$(VERSION)" ./cmd/eco

# Atualiza o snapshot do tupa-go vendored em internal/tupavendor/source/.
# Use REF=v0.1.0 para puxar uma tag específica (default: main).
sync-tupa:
	@bash scripts/sync-tupa-vendor.sh $(REF)

clean:
	rm -rf bin/ dist/ release/
