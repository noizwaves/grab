#!/usr/bin/env bash
set -eou pipefail

if type grab 2>/dev/null; then
  echo "grab already on path, skipping installation"
  exit 0
fi

DESTINATION="$HOME/.local/bin/grab"

if [ -f "$DESTINATION" ]; then
  echo "grab already exists at $DESTINATION, skipping installation"
  exit 0
fi

function get_arch() {
  # `arch` doesn't exist in cachyos
  # `uname -i` returns unknown on cachyos
  local arch=`uname --machine`
  case "$arch" in
    x86_64)
      echo "amd64"
      ;;
    aarch64)
      echo "arm64"
      ;;
    *)
      >&2 echo "Unsupported architecture: $arch"
      exit 1
  esac
}

function get_platform() {
  local linux="linux"
  local darwin="darwin"

  # https://stackoverflow.com/a/8597411
  if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    echo "linux"
  elif [[ "$OSTYPE" == "darwin"* ]]; then
    echo "darwin"
  else
    >&2 echo "Unsupported platform: $OSTYPE"
    exit 1
  fi
}

mkdir -p $(dirname "$DESTINATION")

echo "Downloading latest version of grab..."
GRAB_ARCH=get_arch
GRAB_PLATFORM=get_platform
SOURCE="https://github.com/noizwaves/grab/releases/download/latest/grab-$(get_platform)-$(get_arch)"
curl -L --output "$DESTINATION" --silent "$SOURCE"
chmod +x "$DESTINATION"

echo "Installed grab to $DESTINATION"
