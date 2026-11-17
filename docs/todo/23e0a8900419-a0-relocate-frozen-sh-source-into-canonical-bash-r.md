---
id: 23e0a8900419
kind: refactor
title: A0 Relocate frozen sh source into canonical Bash root module
seq: 1
status: todo
priority: p0
labels:
    - layer-separation
created: 2026-11-17T22:11:04.030724Z
sprint: 575
sprint_id: 65ce85f5-87c0-5863-bf55-62de86f5df36
sprint_title: 'Layer separation: Classic <- Bash# <- Yoke import-graph delineation'
---

Sprint575 foundation relocation slice ONLY; no engine extraction or behavior changes. Canonical root module decision github.com/qiangli/bash, with flat syntax/expand/interp/etc packages; existing coreutils remains its original qiangli/coreutils module/pin. Import the exact sh snapshot a4141771a3dc9b920f17864bcc37c75c69159701 into qiangli/bash root (existing thin repo currently pins sh/coreutils). Preserve all source/tests/fixtures byte-for-byte except mechanical mvdan.cc/sh/v3 -> github.com/qiangli/bash package/import references and module identity. Do not strip Bash# implementation in this slice; this is recoverable source relocation before the separately reviewed physical extraction. Remove only obsolete sh gitlink/.gitmodules entry, retain coreutils pin. Keep source history recoverable and record original repo+SHA in a migration note/commit. Preserve BSD3 sh license and existing MIT project license distinctly; no relicensing. Preserve independent nested modules/fixtures, updating package references only where needed; list any unresolved pin work rather than invent versions. Add minimal make-free DAG build/test front doors for root library module (official image only), no product CLI relocation yet because S3/design owns that boundary. Do not run builds/tests/foreign suites on dev box; allowed static DAG parse and source-copy/reference-only equivalence comparison. No dependency upgrades, Go behavior/API edits, image lock changes, test assertion weakening, old repo archiving, remote pushes or root umbrella edits. Do not commit generated KB.md. Commit isolated branch only with injected actual agent attribution and proper Sprint/Story/Story-ID. Manager reviews source equivalence, runs official-image gates and publishes. This slice does not claim final layer isolation or finished standalone binary home.
