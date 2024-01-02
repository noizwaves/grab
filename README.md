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
