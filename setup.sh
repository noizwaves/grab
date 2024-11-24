#!/usr/bin/env bash
set -e

mise install

go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2

go install mvdan.cc/gofumpt@latest

go install golang.org/x/tools/gopls@latest
