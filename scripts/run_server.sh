#!/usr/bin/env bash

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

exec "$ROOT/bin/kingdom-server" \
     -addr ":8000" \
     -db   "$ROOT/kingdom.db" \
     "$@"