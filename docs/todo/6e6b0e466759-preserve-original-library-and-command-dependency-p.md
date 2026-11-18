---
id: 6e6b0e466759
kind: task
title: Preserve original library and command dependency profiles in separate modules
seq: 183
status: assigned
created: 2026-11-18T05:54:30.486483Z
assignee: s575-consumer-refs
sprint: 575
sprint_id: 65ce85f5-87c0-5863-bf55-62de86f5df36
sprint_title: 'Layer separation: Classic <- Bash# <- Yoke import-graph delineation'
---

Official-image root cause proven by manager: TestFcS88EditorSessionThroughInteractivePTY fails under the relocated CLI Readline fork profile but identical engine source passes with its original upstream Readline profile. Preserve both original profiles by placing existing command sources in a nested cmd module with its original CLI dependencies and own-parent Bash replacement; restore library go.mod/go.sum exactly to accepted A0. Move added native launcher test unchanged except relocation paths. Build/DAG/Makefile contexts cover both modules without changing engine implementation or test assertions. Manager owns official metadata reconciliation and verification.
