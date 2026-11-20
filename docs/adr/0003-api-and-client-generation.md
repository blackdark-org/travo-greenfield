# 0003 Generated API spec and client

Status: accepted

Huma generates OpenAPI 3.1 from typed Go handlers; hey-api generates the TypeScript client (and optionally TanStack Query bindings) from that spec. Nothing is hand-maintained. CI regenerates both and fails on a diff. Spec served publicly at `/api/openapi.json`.

If a UI boundary not covered by the generated client needs runtime validation, use Zod Mini or Valibot, per the preset size rules. None has needed it so far, so neither is a dependency.
