#!/usr/bin/env bash

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BIN_DIR="$ROOT/bin"

mkdir -p "$BIN_DIR"

pushd "$ROOT" >/dev/null

go mod tidy

echo "==> Building kingdom_server …"
go build -o "$BIN_DIR/kingdom-server" ./cmd/kingdom/

echo "==> Building kingdom_client …"
go build -o "$BIN_DIR/kingdom-client" ./ui/tui/

echo "==> Building kingdom_agent  …"
go build -o "$BIN_DIR/kingdom-agent" ./cmd/agent/

echo "All binaries built in $BIN_DIR"

popd >/dev/null
