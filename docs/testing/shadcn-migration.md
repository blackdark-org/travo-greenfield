# shadcn-svelte migration evidence (#119)

Base: `95ecdf0ea9feec6a518be38c84a72272f2ec0166`. Budgets remain shell 40,960 B,
each lazy JavaScript chunk 8,192 B, total 102,400 B, gzip level 9.

Early production probes (before migrating consumers):

| Configuration | Shell B | Largest lazy B | Total B |
| --- | ---: | ---: | ---: |
| Baseline | 37,851 | 5,271 | 58,711 |
| Default upstream Button + Dialog on WiFi route | 63,878 | 19,152 | 98,879 |
| variants/lite, clsx, on-demand Dialog | 44,947 | 15,144 | 80,729 |
| Static variants, explicit lazy library chunks | 39,757 | 6,219 | 77,385 |

The optimized probe still uses the upstream Button template and Bits UI Dialog.
No accessibility functionality is removed. The final inventory, browser checks
and measurements will accompany the consumer migration.

Before screenshots: `web/screenshots/before/` (ignored local artifacts), both
390 and 1280 px, light/dark, login/setup/Dashboard/WiFi/VPN/System, WiFi join and
mobile drawer. Captured against the mock daemon on isolated ports 18181/18192.
Preview snapshots failed twice because no screenshot frames were produced;
the repository's headless Chrome helper was used instead.
