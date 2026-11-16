# Copied shadcn-svelte source

Source: https://github.com/huntabyte/shadcn-svelte/tree/28f5a3bef5ea7934044a7fb60a3e6f201b3b090a/docs/src/lib/registry/ui

Retrieved 2026-11-16. Svelte 5 / Tailwind 4 registry; MIT notice in LICENSE.md.
`src/lib/utils.ts` derives its types from registry/lib/utils.ts and its `cn` signature
from the official manual installation instructions. These files are owned and
maintained here; updates are reviewed diffs, never a registry overwrite.

Local adaptations: `$lib` aliases, existing Lucide icon imports, Paraglide close
labels, semantic palette and spacing. Static variant lookup replaces the general
`tailwind-variants` engine; `cn` uses clsx without runtime Tailwind merging. Put
base appearance in component-layer CSS and caller overrides in utilities, rather
than depending on contradictory utility strings. See ADR 0013 for measured costs.
