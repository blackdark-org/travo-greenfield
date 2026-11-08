# 0011 Caller classification for lockout protection

Status: accepted (payloads via spike S3)

Source IP → neighbour table (ARP) → MAC → AP station dumps (`hostapd.*` `get_clients`) → ethernet / wifi-ap / wifi-client / unknown, because `br-lan` mixes wired and WiFi clients. One pure function, table-driven tests on verbatim captured payloads. Unknown fails closed. Refusal: 409 `wifi_lockout_risk`, overridable only with `acknowledge_lockout: true`.
