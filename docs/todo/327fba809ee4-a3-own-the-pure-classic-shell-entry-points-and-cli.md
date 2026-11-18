---
id: 327fba809ee4
kind: refactor
title: A3 Own the pure Classic shell entry points and CLI
seq: 177
status: assigned
priority: p0
created: 2026-11-18T01:05:50.953516Z
weave: 4
assignee: s575-classic-cli
sprint: 575
sprint_id: 65ce85f5-87c0-5863-bf55-62de86f5df36
sprint_title: 'Layer separation: Classic <- Bash# <- Yoke import-graph delineation'
---

Sprint575 minimum-change CLI ownership slice, from the accepted runtime boundary and independent CLI source/design reviews. Foundation main is accepted72c4bd210f74. Copy Classic CLI implementation from frozen bashy27255b0d internal/cli into this module internal/cli; add cmd/bash and cmd/sh independent entry points. Repoint engine imports to this module. Preserve startup/login/profile/rc, flags/argv0/sh/POSIX parsing, prompts/readline, carrier/process groups/job control/signals, streaming/parser retries, inherited FDs, environment, umask, ownedexec frame/adoption and existing Classic assertions. No lower-to-upper imports in source or tests: no bashy, Bashsharp, coreutils or Yoke dependencies.

Remove only Bashsharp activation and Go-source/fenced-runtime dispatch from copied CLI: front/polyglot imports; dialect/source globals/resolution/dispatch and lang selectors; gosource_startup. Pure Bash ignores inherited B-only activation environment, so a parent exporting it cannot make valid Bash fail. Explicit unsupported B invocation options (including off aliases and --source) refuse status2 with a pointer to bashy/bashsharp. Preserve existing shell set/shopt unknown-option semantics: do not add special CLI refusals there. No B-named false-return stubs. Do not strip the foundation engine in this story; its separate ownership story follows, so this slice does not yet claim zero internal polyglot linkage.

Copy package-local Classic tests and helpers beside A CLI. Preserve the POSIX parse/execution main_test cases and startup/carrier/sh/flags tests, splitting mixed cases by assertion. Exclude B activation/Go-source tests and every test importing or building upper product/utility modules; leave original C coverage in place for later coordinated move, do not delete it here. For main_test heredoc diagnostics, preserve script/input/output/warnings but replace its coreutils cat dependency with the reviewed test-local cat exec handler. The in-process coreutils touch/umask test stays with coreutils/product integration, not A. Preserve Classic halves of mixed selector tests; replace only assertions intentionally changed by rejecting B invocation options. Add meaningful selector-env-inert and explicit-refusal cases if absent. Keep private helpers/flag state private; no public test-support API.

Preserve the existing product startup signal contract for standalone artifacts: copy neutral native/siglaunch.c.in, scripts/go-product.sh, scripts/launcher-wanted.sh and tools/productgo/main.go plus GO-LICENSE from the same frozen C tree. Copy only their directly relevant neutral tests (main/managed/native signal), leaving C bootstrap/artifact tests with C. Add DAG build-bash/build-sh targets using those exact native launcher and pinned runtime overlay mechanisms, updating module/linker references only; no sibling files or go.work needed. Existing root foundation build/test targets stay. No launcher/runtime logic changes or dependency version upgrades; direct-require promotions only at existing graph versions if necessary and documented. Retain licenses and modes.

Scope excludes engine stripping, Bashsharp/Yoke/C consumer rewiring, unrelated feature/bug fixes or cleanup. No mutation outside isolated workspace, no pushes or branch switches. ALL Go builds/tests and foreign suites run only in manager's digest-locked official image; static formatting/metadata inspection locally, no Go build/test/run. No skips/assertion weakening. Commit with correct Sprint575/Story/Story-ID and actual agent. Report exact copied file inventory, intentionally omitted upper tests with their preserved original home, structural removals and final SHA. Manager checks source equivalence, standalone build/CLI tests and full named Bash fixture gate before final independence claim.
