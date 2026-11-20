---
id: a16b05f9c1a8
kind: chore
title: Remove obsolete cmd/gosh proof-of-concept command
seq: 179
status: done
created: 2026-11-18T04:43:20.21526Z
assignee: s575-codex-manager
sprint: 575
sprint_id: 65ce85f5-87c0-5863-bf55-62de86f5df36
sprint_title: 'Layer separation: Classic <- Bash# <- Yoke import-graph delineation'
closed: 2026-11-20T12:26:21.015331Z
closed_by: s575-codex-manager
---

Operator requested removing cmd/gosh now that the full cmd/bash artifact builds and runs. Delete only the obsolete command package and its command-specific tests; remove its active build references. Shared shell implementation and POSIX certification tests remain untouched.

## Manager closeout — 2026-11-20
Operator-requested obsolete cmd/gosh deletion is published at 64a08b80; retained NO_GOSH_BUILD_PASS and subsequent independent builds/suites verify the remaining artifact routes. Shared engine and licensed test fixtures were preserved.
Existing receipts and source ancestry were reviewed; no new test/build ran during reconciliation. Evidence belongs to its exact historical source/artifact identities, not arbitrary later heads. S575 full-run evidence remains under posix-baseline/s575-2026-11-19 with its raw 4 FAIL, 54 UNRESOLVED and 9 INSPECT unchanged. RC38 is on operator hold.
