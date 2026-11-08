# 0001 Clean-room greenfield, named travo-greenfield

Status: accepted

Built in parallel to the existing Travo repo. Reuses its lessons and captured device payloads, not its code. If it proves better, it replaces or moves into `travo`.

The project, binary, init script and UCI package are named `travo-greenfield` so both can coexist on one device while comparing. `invariants.md` writes the state paths as `/etc/travo/` and `/etc/travo-guards/`; they take the same `travo-greenfield` suffix, defined once as a constant at implementation time.
