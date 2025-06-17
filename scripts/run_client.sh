#!/usr/bin/env bash


set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
: "${C2_URL:=http://127.0.0.1:8000}"

export C2_URL
exec "$ROOT/bin/kingdom-client" "$@"