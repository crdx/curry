set quiet := true

[private]
help:
    just --list --unsorted --list-submodules

build:
    #!/bin/bash
    set -eo pipefail
    unbuffer go build -trimpath -o dist/curry | gostack

fmt:
    just --fmt
    find . -name '*.just' -print0 | xargs -0 -I{} just --fmt -f {}
    go fmt ./...

lint:
    #!/bin/bash
    set -eo pipefail
    unbuffer go vet ./... | gostack
    unbuffer golangci-lint --color never run | gostack

fix:
    #!/bin/bash
    set -eo pipefail
    unbuffer golangci-lint --color never run --fix | gostack

test:
    #!/bin/bash
    set -eo pipefail
    unbuffer go test -cover ./... | gostack --test
