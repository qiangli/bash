---
id: 98fc4a609fcf
kind: chore
title: Preserve the imported release engine without deferred hook extractions
seq: 182
status: done
created: 2026-11-18T05:22:02.627017Z
assignee: s575-codex-manager
sprint: 575
sprint_id: 65ce85f5-87c0-5863-bf55-62de86f5df36
sprint_title: 'Layer separation: Classic <- Bash# <- Yoke import-graph delineation'
closed: 2026-11-20T12:26:22.908716Z
closed_by: s575-codex-manager
---

Operator holds the internal sh split and requires engine implementation unchanged. The packaging main branch inherited earlier A1/A2 hook/error-detail additions. Restore the affected engine trees exactly from accepted imported release snapshot adb3b4e, preserving all baseline implementation/tests and all command/build/reference changes outside those trees. Earlier additions remain in published Git history for separately scoped future work.

## Manager closeout — 2026-11-20
The 16be6198 restoration preserves the accepted imported engine snapshot adb3b4e; later standalone qualification covers the restored engine with separate original dependency profiles. Earlier hook/extraction attempts remain in history and are not part of this accepted delivery.
Existing receipts and source ancestry were reviewed; no new test/build ran during reconciliation. Evidence belongs to its exact historical source/artifact identities, not arbitrary later heads. S575 full-run evidence remains under posix-baseline/s575-2026-11-19 with its raw 4 FAIL, 54 UNRESOLVED and 9 INSPECT unchanged. RC38 is on operator hold.
