# Mosaic SDUI

The published **Server-Driven-UI contract** for [Mosaic](https://github.com/mosaic-media). It is to the interface what [`sdk`](https://github.com/mosaic-media/sdk) is to content: the single, language-neutral surface a **producer** (the Platform, or a Module) emits and a **client** (the Shell, a future Flutter client) renders — so nobody re-writes it per language.

## Single source of truth

**[`schema/sdui.schema.json`](schema/sdui.schema.json) is the source of truth.** The language bindings are **generated** from it — never hand-written:

```
schema/sdui.schema.json        ← the contract (JSON Schema 2020-12). Edit this.
        │  scripts/generate.sh (quicktype)
        ├──────────────► sdui/contract/contract.gen.go   (Go types)   — generated
        ├──────────────► ts/contract.gen.ts              (TS types)   — generated
        └──────────────► (dart, swift, … when a client needs them)
```

Two guards keep it honest:

- **Drift guard** — `scripts/check-generated.sh` regenerates and fails if the committed bindings are stale (run it in CI). Change the schema → regenerate → commit.
- **Conformance tests** (`sdui/conformance_test.go`) — validate what the hand-written builders *produce*, and every file in `definitions/`, against the schema. So even the ergonomic layer cannot drift from the contract.

**JSON Schema, not protobuf** — the node tree is open (props are an untyped bag by design), it rides GraphQL as JSON, and the definitions and tokens are JSON data. See [ADR 0025](https://github.com/mosaic-media/architecture/blob/main/docs/adr/0025-sdui-contract-repository.md).

## Layout

```
schema/         the single source of truth (JSON Schema)
sdui/
  contract/     GENERATED Go types (contract.gen.go) — do not edit
  sdui.go       aliases + constants over the generated types
  actions.go    Action constructors
  components.go standard-component builders (Screen, Section, PosterCard, …)
ts/
  contract.gen.ts  GENERATED TypeScript types — do not edit
definitions/    the standard component library, as data (a client registers these)
tokens/         design tokens (W3C DTCG) — compiled to CSS vars / Flutter theme
scripts/        generate.sh, check-generated.sh
```

Only the ergonomic builders are hand-written; they sit *on top of* the generated types and are held to the schema by the conformance tests. Generation is also wired to `go generate ./...`.

## Using it — a Go producer

```go
import "github.com/mosaic-media/sdui/sdui"

home := sdui.Screen(sdui.Child(
    sdui.HeroBanner("Spirited Away",
        sdui.Meta("2001", "Anime Film", "PG"),
        sdui.Slot("actions",
            sdui.Button("Play", "primary", sdui.Play("part-1")),
        ),
    ),
    sdui.Section("Continue watching", sdui.Child(
        sdui.Carousel(sdui.Child(
            sdui.PosterCard("Cowboy Bebop", "Anime Series",
                sdui.Progress(0.6),
                sdui.Act(sdui.Navigate("detail", map[string]any{"title": "Cowboy Bebop"})),
            ),
        )),
    )),
))
// json.Marshal(home) → exactly the payload the Shell renders.
```

Add it like any Go module:

```bash
go get github.com/mosaic-media/sdui@v0.1.0
```

For local work across the sibling repos, use a
`replace github.com/mosaic-media/sdui => ../sdui` in the consumer's
`go.mod` instead.

## Using it — a TypeScript client

Published to npm as **`@mosaic-media/sdui`**:

```bash
npm install @mosaic-media/sdui
```

```ts
import type { UINode, Action, ComponentDefinition } from "@mosaic-media/sdui";
import heroBanner from "@mosaic-media/sdui/definitions/hero-banner.json";
import tokens from "@mosaic-media/sdui/tokens.json";
```

The package is types + JSON data (no runtime code); it's meant for a bundler
(the Shell uses Vite). Until the first npm release lands you can install straight
from git: `npm install github:mosaic-media/sdui`.

## The standard definitions

The reusable components — `PosterCard`, `HeroBanner`, `Section`, `Badge`, … — live here as `ComponentDefinition` data, not per-client code. A client registers them; a producer emits `{ "type": "HeroBanner", … }` and it renders identically on every client, with the Module shipping **zero** UI code. A Module can ship its own definitions the same way. Only the irreducible **primitives** are per-client native code; definitions compose only those ([ADR 0024](https://github.com/mosaic-media/architecture/blob/main/docs/adr/0024-primitives-and-definitions.md)).

## Regenerating

```bash
scripts/generate.sh         # regenerate Go + TS from the schema
scripts/check-generated.sh  # fail if committed bindings are stale
go test ./...               # unit + schema-conformance tests
```

Requires `npx` (quicktype is fetched on demand) and `gofmt`.

## Releasing

One tag publishes every language:

```bash
# bump package.json + tag must match
git tag v0.1.0 && git push origin v0.1.0
```

- **Go** needs nothing more — the module proxy serves the tag; consumers `go get …@v0.1.0`.
- **npm** is published by `.github/workflows/release.yml` on the tag. It requires an `NPM_TOKEN` repository secret (an npm automation token for the `@mosaic-media` scope). Without it, publish manually: `npm publish --access public`.

`.github/workflows/verify.yml` runs on every push: it fails if the generated bindings are stale, runs the Go + conformance tests, and typechecks the TypeScript.

## Next

- Wire the Shell to import `@mosaic-media/sdui`, load `definitions/*.json`, and generate its CSS variables from the tokens — retiring its local copies.
- Migrate the remaining standard definitions from the Shell into `definitions/`.
- A tokens generator (DTCG → CSS + Dart) and the light theme.
- Add the Dart target to `generate.sh` when the Flutter client lands.

## Licence

**Apache-2.0** (see [`LICENSE`](LICENSE) and [`NOTICE`](NOTICE)). A contract surface must be permissive so a Module may build its UI against it under any licence, as the SDK is ([ADR 0022](https://github.com/mosaic-media/architecture/blob/main/docs/adr/0022-licensing.md), [ADR 0025](https://github.com/mosaic-media/architecture/blob/main/docs/adr/0025-sdui-contract-repository.md)).
