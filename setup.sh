#!/usr/bin/env bash
set -e

mise install

go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.2.2

go install mvdan.cc/gofumpt@latest

go install golang.org/x/tools/gopls@latest
