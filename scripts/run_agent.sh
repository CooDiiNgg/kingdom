#!/usr/bin/env bash

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

BOOTSTRAP=$(ls -t "$ROOT"/agent_*_bootstrap.sh 2>/dev/null | head -n 1 || true)

if [[ -z "$BOOTSTRAP" ]]; then
  echo "[run_agent] No agent bootstrap script found in project root. Create an agent via the TUI first." >&2
  exit 1
fi

echo "[run_agent] Executing $BOOTSTRAP …"
exec bash "$BOOTSTRAP" "$@"