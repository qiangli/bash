---
id: 4635ac5d7ae2
kind: task
title: Provide standalone Bash artifact build and conformance DAG
seq: 181
status: done
created: 2026-11-18T04:53:36.477336Z
assignee: s575-codex-manager
sprint: 575
sprint_id: 65ce85f5-87c0-5863-bf55-62de86f5df36
sprint_title: 'Layer separation: Classic <- Bash# <- Yoke import-graph delineation'
closed: 2026-11-20T12:26:22.30891Z
closed_by: s575-codex-manager
---

Standalone repository build infrastructure only. Preserve engine and CLI implementation unchanged. Add Dockerfile and canonical DAG targets for fixed-image GOWORK=off build, native artifact, regression/smoke, sealed GNU Bash 5.3 fixture coverage. Licensed POSIX shell certification remains a distinct explicitly configured private harness gate. No umbrella or sibling source mounts.

## Manager closeout — 2026-11-20
Standalone DAG/container infrastructure published from 083ce721 through ed6bbe5. Original root/command Go suites qualified at af9c9f90, artifact/smoke at cdb1212, GNU86 mode0 at a6e766b7, and final invocation-only DAG/Docker checks at ed6bbe5. Recorded input-identity reuse is explicit; GNU86 is not licensed certification.
Existing receipts and source ancestry were reviewed; no new test/build ran during reconciliation. Evidence belongs to its exact historical source/artifact identities, not arbitrary later heads. S575 full-run evidence remains under posix-baseline/s575-2026-11-19 with its raw 4 FAIL, 54 UNRESOLVED and 9 INSPECT unchanged. RC38 is on operator hold.
