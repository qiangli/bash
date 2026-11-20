---
id: 78e7f5cf4f6b
kind: refactor
title: Move unchanged Bash CLI entrypoint from bashy into Bash repository
seq: 180
status: done
created: 2026-11-18T04:52:55.506037Z
assignee: s575-codex-manager
sprint: 575
sprint_id: 65ce85f5-87c0-5863-bf55-62de86f5df36
sprint_title: 'Layer separation: Classic <- Bash# <- Yoke import-graph delineation'
closed: 2026-11-20T12:26:21.644506Z
closed_by: s575-codex-manager
---

Operator: just move bashy/cmd/bash into bash/. Both command source files remain unchanged; reference the pinned existing shared CLI in Bashsharp and released Readline fork. Engine code and behavior unchanged. Verify standalone official-image go build and binary execution.

## Manager closeout — 2026-11-20
Unchanged command entrypoints relocated in a9c50210 and 05671b76. Retained official-image standalone builds and executable smoke receipts passed. This accepts command ownership/reference relocation, not the deferred shared-engine split.
Existing receipts and source ancestry were reviewed; no new test/build ran during reconciliation. Evidence belongs to its exact historical source/artifact identities, not arbitrary later heads. S575 full-run evidence remains under posix-baseline/s575-2026-11-19 with its raw 4 FAIL, 54 UNRESOLVED and 9 INSPECT unchanged. RC38 is on operator hold.
