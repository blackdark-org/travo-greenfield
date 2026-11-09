# 0004 Frontend stack

Status: accepted

Svelte 5 + plain Vite static SPA, shadcn-svelte (only used components), hash routing, pre-gzipped assets, `size-limit` budget 50 KB gzip, system font stack, own SVG charts. Follows agent-presets `stack/embedded-ui.md`.

Deviation from the preset: Paraglide is used with `en` as the only locale for now (preset default kept; English-only is a locale list, not a rewrite). Record in `AGENTS.md`.

Implementation (`web/`):
- Pinned exactly in `web/package.json`; TypeScript 6.0.3. TS 7.0 (native port) has no classic JS compiler API: `@hey-api/openapi-ts` 0.99 crashes (`ts.SyntaxKind` undefined; fixed only in `@next` prereleases, hey-api#4235) and `svelte-check` 4.7 needs TS 6 installed next to TS 7 via an npm alias (sveltejs/language-tools#3063). Move to TS 7 when hey-api ships the TypeScript-free generator as `latest` and svelte-check drops the TS 6 requirement (planned with TS 7.1 APIs).
- pnpm's default `minimumReleaseAge` (24 h) is kept with no exclusions: pin versions older than a day.
- Generated: `web/src/lib/api/gen/**` (hey-api) is committed so the UI builds and reviews show API changes; `mise run check:gen` fails on drift, including new untracked files. `web/openapi.json` is only an intermediate written by `go run ./cmd/openapi` and git-ignored: the Go handlers are the single source, and the daemon serves the spec at runtime. Paraglide output is git-ignored (its default) and compiled in `postinstall`; the inlang message-format plugin is an npm devDependency, never fetched from a CDN.
- Dev: `travo-greenfield --mock` on `:8081`, `pnpm -C web dev`; `DAEMON_URL` overrides the proxy target.
- Below `md` the sidebar is a drawer: `inert` when closed, Escape closes, focus moves to Close and back to Menu.
- Biome overrides for `.svelte` files switch off `noUnusedImports`, `noUnusedVariables`, `useConst` and `useImportType`: Biome only parses the `<script>` block, so it cannot see template usage or runes (`$state`, `$props`) and would flag or rewrite code that is in fact used.
- Nav groups with collapse state and `GroupHeading` usage arrive with the first real grouped pages; the sidebar is flat until then.
- Pages are code-split per route (`web/src/lib/pages/lazy.ts`): each chunk loads on first visit behind the `LoadState` skeleton, a failed import shows the error with retry. The size-limit glob covers every emitted chunk.
- `mise run check` enforces the 50 KB gzip budget after the build.
