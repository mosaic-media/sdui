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

**JSON Schema, not protobuf, for the *authoring* layer** — the node tree is open (props are an untyped bag by design) and the definitions and tokens are JSON data. See [ADR 0025](https://github.com/mosaic-media/architecture/blob/main/docs/adr/0025-sdui-contract-repository.md). The *wire* is protobuf end to end: `UINode` is generated as a message too ([ADR 0044](https://github.com/mosaic-media/architecture/blob/main/docs/adr/0044-contracts-protobuf-workspace.md)), and since [ADR 0061](https://github.com/mosaic-media/architecture/blob/main/docs/adr/0061-one-client-transport.md) protobuf/Connect is the *only* client transport — there is no GraphQL surface left to ride.

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

Screens are authored with the declarative `ui` layer — a widget tree where
children, props and slots intermix. `Build()` compiles to the protobuf `UINode`
the transport carries.

```go
import "github.com/mosaic-media/sdui/ui"

home := ui.Screen(
    ui.Hero("Spirited Away",
        ui.Meta("2001", "Anime Film", "PG"),
        ui.Actions(
            ui.Button("Play", "primary", ui.OnTap(ui.Play("part-1"))),
        ),
    ),
    ui.Section("Continue watching",
        ui.Carousel(
            ui.PosterCard("Cowboy Bebop", "Anime Series",
                ui.Progress(0.6),
                ui.OnTap(ui.Navigate("detail", map[string]any{"title": "Cowboy Bebop"})),
            ),
        ),
    ),
).Build() // → exactly the UINode the Shell renders.
```

The `ui` constructors are **generated** from [`ui.spec.json`](ui.spec.json) by
[`tools/genui`](tools/genui), which also emits the TypeScript twin and lints the
spec against the standard definitions (see [Regenerating](#regenerating)). The
`github.com/mosaic-media/sdui/sdui` package keeps the shared types (`Node`,
`Action`, the action constructors, tone/type constants).

Add it like any Go module:

```bash
go get github.com/mosaic-media/sdui/ui@latest
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

// The generated protobuf client transport — the service descriptors + message
// schemas a Connect-Web client speaks, and the wire UINode message. Pair with
// @connectrpc/connect + @connectrpc/connect-web.
//
//   /auth    — AuthService: sign in, get a session (ADR 0061).
//   /session — SessionService: the two-lane live session (ADR 0041).
//
// Together they are the whole client surface; there is no second transport.
import { AuthService } from "@mosaic-media/sdui/auth";
import { SessionService, RegionUpdate_Op } from "@mosaic-media/sdui/session";
import { UINodeSchema } from "@mosaic-media/sdui/sdui-pb";
```

The package is types + JSON data (no runtime code); it's meant for a bundler
(the Shell uses Vite). Until the first npm release lands you can install straight
from git: `npm install github:mosaic-media/sdui`.

### Widget-style authoring (mocks & fixtures)

The `@mosaic-media/sdui/ui` subpath is a declarative authoring layer — the
TypeScript twin of the Go `sdui/ui` package — for hand-writing screens (the
Shell's mock payloads, storybook stories) as a widget tree instead of raw
`UINode` JSON. `build()` returns exactly that JSON. The API mirrors the Go one
name-for-name, so a screen transliterates between the two:

```ts
import {
  Screen, Section, Carousel, Hero, PosterCard, Button,
  Actions, Meta, Progress, OnTap, Navigate, Play,
} from "@mosaic-media/sdui/ui";

const home = Screen(
  Hero("Spirited Away",
    Meta("2001", "Anime Film", "PG"),
    Actions(Button("Play", "primary", OnTap(Play("part-1")))),
  ),
  Section("Continue watching",
    Carousel(
      PosterCard("Cowboy Bebop", "Anime Series",
        Progress(0.6), OnTap(Navigate("detail", { title: "Cowboy Bebop" }))),
    ),
  ),
).build(); // → the same UINode payload the Shell renders
```

## The standard definitions

The reusable components — `PosterCard`, `HeroBanner`, `Section`, `Badge`, … — live here as `ComponentDefinition` data, not per-client code. A client registers them; a producer emits `{ "type": "HeroBanner", … }` and it renders identically on every client, with the Module shipping **zero** UI code. A Module can ship its own definitions the same way. Only the irreducible **primitives** are per-client native code; definitions compose only those ([ADR 0024](https://github.com/mosaic-media/architecture/blob/main/docs/adr/0024-primitives-and-definitions.md)).

## Regenerating

Two generators, both driven by a single source of truth:

- **The contract** (`schema/sdui.schema.json`) → `sdui/contract/contract.gen.go`,
  `ts/contract.gen.ts` via quicktype.
- **The `ui` authoring layer** (`ui.spec.json`) → `ui/components.gen.go`,
  `ts/ui.ts` via [`tools/genui`](tools/genui). The same tool **lints** the spec
  against `definitions/*.json`: every definition must have a component, and every
  prop a definition's template binds must be exposed by some helper — so a new
  component that nothing authors fails the build, and Go/TS can never drift.

**Everything runs in a container; nothing is generated or tested on the host.**
The full gate — version check, drift guard, Go tests, TypeScript typecheck — is
one command:

```bash
docker compose -f docker-compose.test.yml run --rm test
```

Regenerating, or running one step, is the same command with the step named:

```bash
docker compose -f docker-compose.test.yml run --rm test bash scripts/generate.sh        # regenerate (schema + ui spec), then lint
docker compose -f docker-compose.test.yml run --rm test bash scripts/check-generated.sh  # fail if any generated file is stale or lint fails
docker compose -f docker-compose.test.yml run --rm test go run ./tools/genui -lint       # just lint the ui spec against the definitions
docker compose -f docker-compose.test.yml run --rm test go test ./...                    # unit + schema-conformance tests
```

Editing the `ui` layer means editing `ui.spec.json` (and adding a component there
when you add a `definitions/*.json`), never the generated files. This is the one
repository that needs **two** toolchains at once — the drift guard regenerates
Go *and* TypeScript and needs `go`, `gofmt`, `node` and `npx` together — which is
why it has a `Dockerfile.test` of its own that pins all of them. A host with
three of the four produces a check that passes by not running.

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
