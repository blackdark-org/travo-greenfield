# Disabling a radio (`wifi-device`)

Travo enables and disables access points, not radios. "Per radio" on the WiFi page means the AP section (`wifi-iface`) on each radio. Travo never writes the `disabled` option of a `wifi-device` section.

## Why this is out of scope

One radio usually carries several roles at once: the uplink STA, the user's AP and the guest AP can share it, and a join can move the AP to the other radio. Disabling the `wifi-device` drops all of them together, including the uplink the router is online through, so the lockout checks that guard a single AP change could not hold. Travo's WiFi changes stay on interface sections, applied with `reload_config` and never `wifi` ([invariants §1](../docs/invariants.md)).

To silence a radio, disable each AP on it, or use LuCI.

## Prior requests

- #74 decided "per-radio" as AP sections.
