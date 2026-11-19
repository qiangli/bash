---
id: 199ef70fd68c
kind: refactor
title: A4 Physically remove the owned Bashsharp engine from Classic
seq: 178
status: todo
priority: p0
labels:
    - superseded-plan
created: 2026-11-18T02:38:52.618234Z
sprint: 575
sprint_id: 65ce85f5-87c0-5863-bf55-62de86f5df36
sprint_title: 'Layer separation: Classic <- Bash# <- Yoke import-graph delineation'
---

Sprint575 full physical separation, after accepted/published B0 source14417523a728 (immutable engine snapshot in Bashsharp). Implement the independently reviewed Classic strip plan in umbrella .agents/s575/classic-strip-plan.log, corrected by docs/s575-runtime-boundary.md. Foundation main has accepted A1/A2 contracts72c4bd2; A3 CLI is a separate disjoint slice and is not yours to edit. B owns its complete superset engine; remove B implementation from A without public parser/Runner plugin redesign or external upstream replacement.

Work only in A syntax/expand/interp/internal/shell and directly required Classic cmd/gosh/shfmt/tests plus root module metadata if strictly necessary. Remove Bashsharp AST/parser/lexer/printer/walk branches, LangBashPP activation, interpreter typed/GoSource/polyglot/fence/task runtime, B-specific state/options, and B package trees lower/gosource/polyglot/grammar. Trace shared files carefully and inline the EXISTING Classic branch at removed B hook calls. No B-named false-return stubs. Preserve all Classic Bash/POSIX/mksh/Bats/zsh behavior, variant bit values, flags/diagnostics, variables/arrays/namerefs, command substitutions, jobs/traps/signals, FDs/redirections/exec, strict POSIX/startup and ownership/frame protocols. Never replace shared files with external upstream versions or lose fork fixes.

CRITICAL correction to the initial read-only strip plan: retain neutral expand ValueKind/Variable/Environ/WriteEnviron, ALL current value-kind constants including Object/KeepValue/Unknown/Unset, generic Object/Variable/serialization helpers and their signatures. B now aliases those identities. Remove A language activation/evaluation branches, NOT the shared neutral data contract. Keep accepted interp HandlerContext public fields/methods/HandlerKind/HandlerHooks/WithHandlerContext/AdaptHandlerContext, ExitStatus/SignaledStatus/compat helpers, ValueField/ValueDescription and errdetail stdlib leaf. Adapted views still never fall back to A Runner; ordinary A views retain Classic behavior. A DescribeValue may retain its neutral Classic projection; typed B reflection stays in B. Stat bare-context behavior remains. Keep A execbudget/fileutil/pathconv/pattern/winmode/interp/ownedexec stable and neutral.

Handler context construction currently always calls bashPPEnv/bashPPWriter even in Classic. Restore direct Classic environment/writers while preserving layered locals/temp assignment, live env and script-stdin override. New/Reset/Run also allocate or enter B state on normal paths: remove only B state/hooks, retaining all ordinary state transitions. Preserve A1/A2 API shapes exactly so B1 aliases/bridge remain compatible without editing B or raising its lower pin to hide failures.

Test ownership: B already preserves frozen syntax/expand/interp/internal/shell and all B runtime/compiler fixture trees. Remove A B-only tests/rows that now live there, split mixed files BY ASSERTION, preserve every Classic input/expectation. Keep referenced semantic evidence tests (history_test, interp_test TestRunnerPosixStdinArgv0, issue7_command/jobcontrol interface, jobcontrol_issue7_unix, startup_env, strictposix, time_issue7) and their Classic test IDs. Do not skip, weaken or rebaseline. Record exact removed tests and their existing B homes; if a B-only command/formatter test is not represented in B, report it for coordinated preservation instead of silently deleting it. A3 copied Classic CLI/tests belong to its separate worker; do not touch internal/cli/cmd/bash/cmd/sh/native/productgo/scripts or their source.

No unrelated behavior fixes, dependency upgrades, broad cleanups, generated KB.md or upper/sibling/root umbrella edits. All Go builds/tests and foreign suites run only in manager's locked official image. Static rg/AST/gofmt/metadata inspection locally only: no go run/build/test. Commit coherent increments to the assigned isolated branch with Sprint575/Story/Story-ID and actual agent; report exact removed implementation, preserved neutral contracts/test ownership and compile questions. Manager will run syntax/expand/interp/full suites plus named86 Classic gate, compare baseline and new B suite, and return named failures. Never fix upper failures by modifying accepted lower contracts.
