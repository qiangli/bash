---
id: 78e7f5cf4f6b
kind: refactor
title: Move unchanged Bash CLI entrypoint from bashy into Bash repository
seq: 180
status: assigned
created: 2026-11-18T04:52:55.506037Z
assignee: s575-codex-manager
sprint: 575
sprint_id: 65ce85f5-87c0-5863-bf55-62de86f5df36
sprint_title: 'Layer separation: Classic <- Bash# <- Yoke import-graph delineation'
---

Operator: just move bashy/cmd/bash into bash/. Both command source files remain unchanged; reference the pinned existing shared CLI in Bashsharp and released Readline fork. Engine code and behavior unchanged. Verify standalone official-image go build and binary execution.
