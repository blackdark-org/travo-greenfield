# Settings pages for optional catalog services

The Services page installs, removes, starts and stops optional OpenWrt services (SQM, Cloudflare Tunnel, vnStat, Watchcat, speedtestcpp) and sets whether they start at boot. It does not configure them. Each card links the service's LuCI app (`:8080`), which the catalog installs next to the package.

## Why this is out of scope

Each of these packages ships a LuCI app that owns its UCI config. A Travo settings page per service would repeat that UI and need its own config model, transaction rules and device fixtures for every package. The catalog stays one data table (`internal/catalog`), so adding a service costs one entry.

A service that travellers use daily and that can cut connectivity gets its own feature instead, as WireGuard, Tailscale and AdGuard Home did.

## Prior requests

- #268 built the catalog and left settings to LuCI.
