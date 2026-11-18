---
id: a16b05f9c1a8
kind: chore
title: Remove obsolete cmd/gosh proof-of-concept command
seq: 179
status: assigned
created: 2026-11-18T04:43:20.21526Z
assignee: s575-codex-manager
sprint: 575
sprint_id: 65ce85f5-87c0-5863-bf55-62de86f5df36
sprint_title: 'Layer separation: Classic <- Bash# <- Yoke import-graph delineation'
---

Operator requested removing cmd/gosh now that the full cmd/bash artifact builds and runs. Delete only the obsolete command package and its command-specific tests; remove its active build references. Shared shell implementation and POSIX certification tests remain untouched.
