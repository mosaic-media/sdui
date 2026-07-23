#!/usr/bin/env bash
# Generate the per-language contract bindings from the single source of truth,
# schema/sdui.schema.json. Do NOT hand-edit the generated files.
#
#   Go   -> sdui/contract/contract.gen.go   (package contract)
#   TS   -> ts/contract.gen.ts
#   Dart -> add when a Dart client lands (quicktype --lang dart)
#
# Run: scripts/generate.sh   (requires npx + gofmt on PATH)
set -euo pipefail
cd "$(dirname "$0")/.."

SCHEMA="schema/sdui.schema.json"
# quicktype is PINNED. Unpinned (`npx --yes quicktype`) fetches the latest,
# which makes the generator — and therefore the drift guard that regenerates
# and diffs — non-deterministic: a new quicktype release silently restyles the
# output, and the guard then fails on a "stale" binding that is only a different
# generator. 25.1.0 is the last release before 26.0.0, whose TypeScript renderer
# stopped annotating the open `props`/`params` maps; pinning here keeps those
# annotations (the type contract) and freezes the formatting. Bump deliberately,
# regenerating and reviewing the diff in the same change.
QT=(npx --yes quicktype@25.1.0 -s schema --top-level MosaicSDUI "$SCHEMA")

echo "generating Go -> sdui/contract/contract.gen.go"
"${QT[@]}" --lang go --package contract --just-types-and-package -o sdui/contract/contract.gen.go
# The consolidated schema has a wrapper root (so the generator reaches every
# top-level type). Strip that unused wrapper struct, then mark as generated.
sed -i '/^type MosaicSDUI struct {/,/^}/d' sdui/contract/contract.gen.go
sed -i '1s|^|// Code generated from schema/sdui.schema.json by quicktype. DO NOT EDIT.\n|' sdui/contract/contract.gen.go
gofmt -w sdui/contract/contract.gen.go

echo "generating TypeScript -> ts/contract.gen.ts"
"${QT[@]}" --lang typescript --just-types -o ts/contract.gen.ts
sed -i '/^export interface MosaicSDUI {/,/^}/d' ts/contract.gen.ts
sed -i '1s|^|// Code generated from schema/sdui.schema.json by quicktype. DO NOT EDIT.\n|' ts/contract.gen.ts

# The ui authoring layer (ui/components.gen.go, ts/ui.ts) is generated from
# ui.spec.json by tools/genui, which also lints the spec against definitions/.
echo "generating ui layer -> ui/components.gen.go, ts/ui.ts (+ spec lint)"
go run ./tools/genui

echo "done."
