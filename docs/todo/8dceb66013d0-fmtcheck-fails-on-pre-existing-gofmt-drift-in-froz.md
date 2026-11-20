---
id: 8dceb66013d0
kind: chore
title: fmtcheck fails on pre-existing gofmt drift in frozen interpreter files (interp/api.go, interp/bashpp_concurrency_unix_test.go, moreinterp/coreutils/coreutils.go)
seq: 184
status: todo
priority: p3
created: 2026-11-20T20:34:19.461847Z
---

Found by S660 R1 (muse, 2026-11-20): scripts/fmtcheck.sh lists 3 files on the untouched base with go1.27.2. Pre-existing; not caused by R1. The interpreter is frozen for certification: a gofmt-only commit needs a recorded decision (whitespace-only, no semantic change) before it lands. Until then bashy dag fmtcheck in this repo fails on these files.
