mod release

set quiet := true
set shell := ["bash", "-cu", "-o", "pipefail"]

[private]
help:
    just --list --unsorted --list-submodules

build:
    unbuffer go build -trimpath -o dist/curry | gostack

fmt:
    go fmt ./...

lint:
    unbuffer go vet ./... | gostack
    unbuffer golangci-lint run --color never | gostack

fix:
    unbuffer golangci-lint run --color never --fix | gostack

test:
    unbuffer go test -cover ./... | gostack --test
