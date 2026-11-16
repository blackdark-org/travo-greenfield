# 0013 Copied shadcn-svelte components

Status: accepted

Use upstream-derived shadcn-svelte primitives throughout the SPA, preserving the
app-specific wrapper contracts and semantic palette. We own the copied source
from huntabyte/shadcn-svelte `28f5a3bef5ea7934044a7fb60a3e6f201b3b090a`
(retrieved 2026-11-16). `web/src/lib/components/ui/LICENSE.md` is the upstream
MIT notice and stays with the copy; the same text ships in the built UI as
`web/public/shadcn-svelte-LICENSE.txt`. Upstream updates are reviewed diffs.
Manual installation avoids overwriting the app theme or introducing Kit.

The initial Button/Dialog probe at 95ecdf0 measured 63,878 B eager shell and
19,152 B WiFi chunk (baseline 37,851 B shell, 58,711 B total). The general variant
and Tailwind-merge engines are too expensive here: even variants/lite plus clsx
and a dynamically imported dialog measured 44,947 B shell, 15,144 B dialog.
Static lookup for the copied Button's finite variant table, clsx, and explicit
lazy library chunks reduced the probe to 39,757 B shell, 6,219 B largest lazy
chunk and 77,385 B total, before theme integration. All three limits remain
unchanged. Component-layer CSS supplies defaults; utility classes supply caller
overrides, so no runtime class conflict resolver is needed.

Dialogs load on demand. Vite's Rolldown chunk groups keep Bits UI primitives,
internal helpers and utilities lazy without pulling their dependencies into the
eager shell. Accessibility remains implemented by Bits UI. No router/device
behavior changes are part of adoption. The final inventory and measurements live
in [the migration evidence](../testing/shadcn-migration.md).
