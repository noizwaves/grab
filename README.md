# dotlocalbin

Userspace centric dotfile dependency manager, aiming for compatibility across macOS/Linux/containerized platforms.

## Usage

1. Create a config file at `~/.dotlocalbin.yml` like:
```yaml
binaries:
  - name: fzf
    source: https://github.com/junegunn/fzf/releases/download/{{ .Version }}/fzf-{{ .Version }}-linux_amd64.tar.gz
    version: "0.45.0"
...
```

2. Run `go run main.go install`

## Configuration Format

- `binaries[*].source` (required): a Go template compatible string with access to variables:
  - `Version`
  - `Platform`: defaults to current `GOOS`
  - `Architecture`: defaults to current `GOARCH`
- `binaries[*].platforms` (optional): per platform (`GOOS`)/architecture (`GOARCH`) overrides for `Platform` and `Architecture` values in source URL. Takes the form:
  ```yaml
  linux:
    amd64: [unknown-linux-musl, x86_64]
    arm64: [unknown-linux-gnu, aarch64]
  darwin:
    amd64: [apple-darwin, x86_64]
    arm64: [apple-darwin, aarch64]
  ```
