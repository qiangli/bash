---
id: 1033a71952d2
kind: refactor
title: A2 Extract the neutral structured-error data contract
seq: 176
status: todo
priority: p0
created: 2026-11-18T00:17:19.883155Z
sprint: 575
sprint_id: 65ce85f5-87c0-5863-bf55-62de86f5df36
sprint_title: 'Layer separation: Classic <- Bash# <- Yoke import-graph delineation'
---

Sprint575 A2: implement only the independently reviewed neutral structured-error boundary (umbrella docs/s575-runtime-boundary.md). Start from accepted frozen Bash source. Add a standard-library-only errdetail leaf owning ErrorDetail with EXACT existing JSON fields/tags and UnmarshalJSON behavior, a deep Clone(*ErrorDetail)*ErrorDetail helper, and ForeignErrorDetail(error)(*ErrorDetail,bool) using errors.As against an implicit interface ForeignErrorDetail()*ErrorDetail. In current polyglot, alias ErrorDetail to that leaf; keep private foreignError, workerErr, constructors, error rendering and unwrap behavior in polyglot. Replace private clone/err receiver methods (illegal on non-local alias) with package-local helper functions, preserving their logic/callers. Add a carrier method returning a clone; keep the existing public polyglot wrapper and its existing behavior. No interpreter/parser/HandlerContext/CLI, dependency or feature changes; A1 proceeds separately in interp. Do not move the worker implementation into the foundation or import coreutils/Yoke/Bashsharp. Preserve legacy plain-string/null/structured JSON, direct/wrapped error rendering, empty/deep causes, copy isolation and unwrap identity. Add only meaningful boundary tests; reuse existing polyglot tests unchanged where possible. All builds/tests/foreign suites run only in manager's locked official image; static formatting and Go metadata parsing allowed locally, no go run/build/test. No skips/assertion weakening, unrelated repairs, KB.md, pushes, branch switches, root umbrella or sibling writes. Commit isolated workspace with proper Sprint575 Story/Story-ID and actual agent. Manager checks source-preserving extraction and original functional suites. Coreutils consumer rewrite and B engine alias adoption follow in their own stories; do not implement them here.
