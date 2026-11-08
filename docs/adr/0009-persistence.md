# 0009 Persistence without a KV store

Status: accepted

UCI is the source of truth (incl. `/etc/config/<name>` for Travo settings). Derived or small state: atomic JSON in `/etc/travo/`. Guards in a separate `/etc/travo-guards/`. No bbolt in 0.1: revocation uses the epoch file, stats history stays in RAM (no per-sample flash writes). History becomes a later task with a retention rule.
