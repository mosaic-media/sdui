#!/usr/bin/env bash
# Drift guard: regenerate the bindings and fail if the committed files are stale.
# Run in CI (and locally) so a schema change without a regenerate can't land.
set -euo pipefail
cd "$(dirname "$0")/.."

scripts/generate.sh >/dev/null

GENERATED=(sdui/contract/contract.gen.go ts/contract.gen.ts ui/components.gen.go ts/ui.ts)
if ! git diff --quiet -- "${GENERATED[@]}"; then
  echo "ERROR: generated bindings are stale. Run scripts/generate.sh and commit." >&2
  git --no-pager diff --stat -- "${GENERATED[@]}" >&2
  exit 1
fi
echo "generated bindings are up to date."
