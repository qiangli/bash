---
id: 98fc4a609fcf
kind: chore
title: Preserve the imported release engine without deferred hook extractions
seq: 182
status: assigned
created: 2026-11-18T05:22:02.627017Z
assignee: s575-codex-manager
sprint: 575
sprint_id: 65ce85f5-87c0-5863-bf55-62de86f5df36
sprint_title: 'Layer separation: Classic <- Bash# <- Yoke import-graph delineation'
---

Operator holds the internal sh split and requires engine implementation unchanged. The packaging main branch inherited earlier A1/A2 hook/error-detail additions. Restore the affected engine trees exactly from accepted imported release snapshot adb3b4e, preserving all baseline implementation/tests and all command/build/reference changes outside those trees. Earlier additions remain in published Git history for separately scoped future work.
