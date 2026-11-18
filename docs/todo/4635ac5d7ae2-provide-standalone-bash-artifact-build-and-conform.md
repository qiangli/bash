---
id: 4635ac5d7ae2
kind: task
title: Provide standalone Bash artifact build and conformance DAG
seq: 181
status: assigned
created: 2026-11-18T04:53:36.477336Z
assignee: s575-consumer-refs
sprint: 575
sprint_id: 65ce85f5-87c0-5863-bf55-62de86f5df36
sprint_title: 'Layer separation: Classic <- Bash# <- Yoke import-graph delineation'
---

Standalone repository build infrastructure only. Preserve engine and CLI implementation unchanged. Add Dockerfile and canonical DAG targets for fixed-image GOWORK=off build, native artifact, regression/smoke, sealed GNU Bash 5.3 fixture coverage. Licensed POSIX shell certification remains a distinct explicitly configured private harness gate. No umbrella or sibling source mounts.
