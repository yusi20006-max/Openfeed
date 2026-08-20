#!/data/data/com.termux/files/usr/bin/bash
set -euo pipefail

# OpenFeed Termux bootstrap.
# This script only prepares Go dependencies; it does not change runtime configuration.

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
cd "$ROOT_DIR"

command -v go >/dev/null 2>&1 || {
  printf '%s\n' "Error: Go is not installed. Install it in Termux with: pkg install golang"
  exit 1
}

GO_VERSION="$(go version | awk '{print $3}')"
printf 'Using %s\n' "$GO_VERSION"

MIN_MAJOR=1
MIN_MINOR=25
VERSION_NUMBER="${GO_VERSION#go}"
MAJOR="${VERSION_NUMBER%%.*}"
REST="${VERSION_NUMBER#*.}"
MINOR="${REST%%.*}"

if [ "$MAJOR" -lt "$MIN_MAJOR" ] || { [ "$MAJOR" -eq "$MIN_MAJOR" ] && [ "$MINOR" -lt "$MIN_MINOR" ]; }; then
  printf '%s\n' "Error: OpenFeed requires Go 1.25 or newer."
  exit 1
fi

# Go's comma fallback only retries on 404/410. A pipe fallback also retries on
# transient/network errors such as 403, which is important on filtered networks.
# Prefer the mirror documented for restricted networks, then try the official
# proxy, and finally direct module downloads.
export GOPROXY="https://goproxy.cn|https://proxy.golang.org|direct"

# The checksum database can itself be unreachable on filtered networks. The
# repository already records go.sum; disabling the remote checksum lookup here
# keeps bootstrap usable when the checksum service is inaccessible.
export GOSUMDB=off

printf '%s\n' "Downloading Go modules..."
if ! go mod download; then
  printf '%s\n' "Error: Go module download failed."
  printf '%s\n' "Try again after checking network access, or run: GOPROXY=direct GOSUMDB=off go mod download"
  exit 1
fi

printf '%s\n' "Building OpenFeed..."
go build -o openfeed ./cmd/server

printf '%s\n' "OpenFeed is ready. Start it with: ./openfeed"
