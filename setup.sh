#!/usr/bin/env bash
set -e

mise settings add idiomatic_version_file_enable_tools go
mise install

go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.2.2

go install mvdan.cc/gofumpt@latest

go install golang.org/x/tools/gopls@latest
