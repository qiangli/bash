---
name: bash
description: Build/test/lint targets for the sh fork as a bashy dag pipeline (agent-first; the Makefile is generated from this file and only forwards here)
---

# Bash — standalone DAG task file

The `github.com/qiangli/bash` repository (Bash 5.3 interpreter patches). The Makefile in
this directory is GENERATED from this file (`bashy dag --makefile`) and only forwards to
the targets below — do not edit it. This DAG file is the agent-first task list, runnable
with `bashy dag`:

```bash
bashy dag --list            # available targets
bashy dag test              # the suite (skips the Docker-only bash oracle)
bashy dag --json test       # machine-readable envelope for an agent
```

This checkout builds with fixed public module pins and `GOWORK=off`; no sibling or umbrella checkout is needed. The official image identity is in `tools/bashy-build/image.lock`.

## Tasks

### build
Build the library and the bash and shfmt commands.
Sources: cmd/, interp/, expand/, syntax/, pattern/, shell/, go.mod, go.sum
Effects: read

```bash
set -e
GOWORK=off go build -mod=readonly ./...
(cd cmd && GOWORK=off go build -mod=readonly ./...)
```

### test
Run the full suite. Skips `TestRunnerRunConfirm`/`TestParseConfirm` — the
behaves-like-real-Bash oracle, which needs Docker + bash 5.3 (see `confirm`).
Effects: read

```bash
set -e
GOWORK=off go test -mod=readonly -count=1 -timeout=30m ./... -skip 'TestRunnerRunConfirm|TestParseConfirm'
(cd cmd && GOWORK=off go test -mod=readonly -count=1 -timeout=30m ./... -skip 'TestRunnerRunConfirm|TestParseConfirm')
```

### test-quick
The quick tier — what push CI runs: build + unit/mock tests (`go test -short`,
without the `//go:build full` evaluator files, < 3 min). Ported from the old
Makefile, whose `test` was this target's alias; here `test` above stays the
full suite, so the tier keeps its own name.
Effects: exec,read

```bash
/bin/bash ./scripts/test-quick.sh
```

### test-moreinterp
`moreinterp/` is a separate Go module (u-root-backed coreutils ExecHandler);
test it independently.
Effects: read

```bash
cd moreinterp && go test ./...
```

### test-race
Race detector over the suite (CI runs this on Linux).
Effects: read

```bash
go test -race ./... -skip 'TestRunnerRunConfirm|TestParseConfirm'
```

### bashpp-race-gate
Bash++ race/lifecycle gate. Records toolchain, package and test discovery,
the broad unsafe-oracle audit, focused schedule runs, and the complete race
suite in `artifacts/bashpp-race-gate.txt`. The real Bash compatibility corpus
remains the separate `confirm` task because it needs the external Bash oracle.
The gate defaults `GOMEMLIMIT` to 2 GiB, gives each `go test` an explicit
timeout plus an outer watchdog, and records portable `/usr/bin/time` peak-RSS
evidence. `GOMEMLIMIT` bounds the Go heap, not total process-tree RSS. CI does
not impose a separate RSS ceiling; the sprint manager supplies the hard 3 GiB
process-tree ceiling for the final acceptance run. A 15-minute global gate
deadline leaves headroom inside the CI job's 30-minute limit for other steps.
Use `scripts/bashpp-race-gate.sh --discovery-only` to validate package and
focused-test discovery without launching race tests. Run
`scripts/bashpp-race-gate-signal-self-test.sh` to send TERM during discovery
and assert the gate and its captured descendants are reaped. The CI job gives
this bounded gate enough time to emit its own diagnostics.
Effects: read, write

```bash
/bin/bash ./scripts/bashpp-race-gate.sh
```

### confirm
The canonical "behaves like real Bash 5.3" oracle: the same script table is run
by the interpreter and by a dockerized Bash, and outputs/exit codes are diffed. Needs
Docker + `mvdan.cc/dockexec`.
Effects: read, net

```bash
CGO_ENABLED=0 go test -run TestRunnerRunConfirm -exec 'dockexec bash:5.3' ./interp
```

### vet
Static checks: gofmt + go vet.
Effects: read

```bash
set -e
gofmt -s -d .
go vet ./...
```

### tidy
go mod tidy + gofmt -s -w . + go vet.
Effects: write

```bash
set -e
go mod tidy
gofmt -s -w .
go vet ./...
```

### clean
go clean.
Effects: destroy

```bash
go clean ./...
```

### fmtcheck
The formatting gate (reports unformatted files, never rewrites them), plus a
check that the generated Makefile still matches this file — regenerate with
`bashy dag --makefile` instead of editing it.
Effects: exec,read

```bash
set -e
./scripts/fmtcheck.sh
bashy dag --makefile --check
```

### hooks
Install the pre-push formatting gate.
Effects: exec,write

```bash
git config core.hooksPath scripts/hooks
echo "hooks installed: core.hooksPath=scripts/hooks"
```

### build-bash
Build the Bash artifact with the released product Go overlay and native signal launcher.
Effects: exec, net, read, write

```bash
set -eu
export GOWORK=off
bin_dir="${BIN_DIR:-bin}"
mkdir -p "$bin_dir"
bin_dir="$(cd "$bin_dir" && pwd -P)"
BASHY_EXE="${BASHY:-bashy}"
go_cmd="${GO:-go}"
goos="$($go_cmd env GOOS)"
goarch="$($go_cmd env GOARCH)"
hostos="$($go_cmd env GOHOSTOS)"
hostarch="$($go_cmd env GOHOSTARCH)"
ext=
[ "$goos" != windows ] || ext=.exe
if [ "$goos" = "$hostos" ] && [ "$goarch" = "$hostarch" ] && [ "$(BASHY="$BASHY_EXE" "$BASHY_EXE" scripts/launcher-wanted.sh build-bash)" = 1 ]; then
  (cd cmd && BASHY_EXE= "$BASHY_EXE" ../scripts/go-product.sh build -mod=readonly -trimpath -o "$bin_dir/bash.real" ./bash)
  cc -x c -std=c11 -O2 -Wall -Wextra -Werror -o "$bin_dir/bash" native/siglaunch.c.in
else
  (cd cmd && BASHY_EXE= "$BASHY_EXE" ../scripts/go-product.sh build -mod=readonly -trimpath -o "$bin_dir/bash$ext" ./bash)
fi
```

### test-go
Run the preserved Go regression suite.
Requires: test
Effects: read

```bash
:
```

### test-full
Run the existing full-tag regression suite without changing test selections.
Effects: exec, net, read, write

```bash
set -e
GOWORK=off go test -mod=readonly -count=1 -tags full -timeout=30m ./...
(cd cmd && GOWORK=off go test -mod=readonly -count=1 -tags full -timeout=30m ./...)
```

### smoke
Smoke-test the actual built Bash artifact.
Requires: build-bash
Effects: exec, read

```bash
set -eu
[ "$(bin/bash -c 'printf "%s" "$((6 * 7))"')" = 42 ]
bin/bash --version
```

### official-build
Build this checkout in the locked official image.
Effects: exec, net, read, write

```bash
BASHY_BUILD_CACHE_VOLUME="${BASHY_BUILD_CACHE_VOLUME:-bash-build-cache}" BASHY_BUILD_CPUS="${BASHY_BUILD_CPUS:-8}" tools/bashy-build/run.sh "$(pwd -P)" -- 'GOWORK=off /bashy dag --no-journal --force --file DAG.md build-bash smoke'
```

### official-test
Run preserved Go regressions in the same locked official image.
Effects: exec, net, read, write

```bash
BASHY_BUILD_CACHE_VOLUME="${BASHY_BUILD_CACHE_VOLUME:-bash-build-cache}" BASHY_BUILD_CPUS="${BASHY_BUILD_CPUS:-8}" tools/bashy-build/run.sh "$(pwd -P)" -- 'GOWORK=off /bashy dag --no-journal --force --file DAG.md test-go'
```

### test-bash-fixtures
Fetch and SHA-256 verify the unchanged GNU Bash 5.3 fixtures (gitignored).
Effects: exec, net, read, write

```bash
"${BASHY:-bashy}" scripts/fetch-bash53-fixtures.sh
```

### test-bash-gate
Sealed GNU Bash 5.3 Classic gate: exactly the existing 86 named fixtures, mode 0.
This is GNU Bash compatibility coverage, not licensed POSIX certification.
Requires: official-build
Effects: exec, net, read, write

```bash
set -eu
export GOWORK=off
export GO="$(pwd -P)/tools/bashy-build/go.sh"
export HOST_GO="$GO"
export BASHY_BUILD_CACHE_VOLUME="${BASHY_BUILD_CACHE_VOLUME:-bash-build-cache}"
export BASHY_BUILD_CPUS="${BASHY_BUILD_CPUS:-8}"
BASH53_ARTIFACT=A BASH53_BASHPP=0 BASH53_GATE=1 BASH53_TESTEE_BIN="$(pwd -P)/bin/bash" BASH53_SOURCE_COMMIT="$(git rev-parse HEAD)" BASHY_TEST_LANE="${BASHY_TEST_LANE:-bash-classic}" BASH53_BASE_IMAGE="${BASH53_BASE_IMAGE:-$(sed -n 's/^image=//p' tools/bashy-build/image.lock)}" BASH53_OCI="${BASH53_OCI:-${BASHY:-bashy} oci}" "${BASHY:-bashy}" scripts/test-bash-container.sh
```

### docker-build
Build this repository Dockerfile after verifying the locked base image digest.
Effects: exec, net, read, write

```bash
set -eu
oci="${BASHY:-bashy} oci"
image=$(sed -n 's/^image=//p' tools/bashy-build/image.lock)
digest=$(sed -n 's/^image_digest=//p' tools/bashy-build/image.lock)
[ "$($oci image inspect --format '{{.Digest}}' "$image")" = "$digest" ]
$oci build --pull=never --build-arg "BUILD_IMAGE=$image" -t "${BASH_IMAGE:-localhost/bash-standalone:dev}" -f Dockerfile .
```

### test-posix-shell
Run the licensed POSIX Profile B **shell-only** release gate, explicitly and
only from a clean, pushed Bashy checkout at the supervisor-approved commit.
The private harness and approval remain outside this repository. Its supervisor
issues a 10-field `.vsc-phase-gate` file (gitignored, mode 0600) in the private
harness checkout; this target never issues or edits approval. Configure a new
unique arm and new private/public evidence roots outside Git, then run:

    cd /absolute/path/to/bashy
    export VSC_PCTS_HARNESS_DIR=/absolute/path/to/private/harness
    export VSC_DO_STATE_DIR=/absolute/path/to/do2-state-root
    export VSC_DO_STATE="$VSC_DO_STATE_DIR/state"
    export VSC_SHELL_ARM=shell-only-unique-arm-name
    export VSC_PRIVATE_EVIDENCE_DIR=/absolute/path/to/new-private-evidence
    export VSC_PUBLIC_EVIDENCE_DIR=/absolute/path/to/new-public-evidence
    bashy dag test-posix-shell

The DAG passes its absolute source root as `BASHY_SOURCE_DIR`. The private
wrapper checks its clean HEAD and sibling pins against the approval before
transfer, builds `cmd/sh` on DO2, wires it as the pure SUT, runs preflight and
phase gate, then dispatches **one shell set/493 TPs** under the original
envelope with external GNU/system providers and no utility sets. It retrieves
sealed evidence into the configured roots and must return success only for
493 certification-pass-group TPs, zero blockers, caps, and runner failures.
Record the arm, exact approved commits, public ledger/runner hashes, and
retrieval verdict in the release evidence. Missing harness/approval, source
drift, incomplete receipt, or any red TP exits nonzero and blocks the Bashy
POSIX release gate. The local strict `bin/sh{,.real}` build uses the canonical `build-sh`
DAG target and command module; the private wrapper must invoke this same
target when building the transferred checkout on its runner. The GNU-compatible
`bin/bash` 86-fixture gate remains separate.
Effects: cred, exec, net, read, write

```bash
set -eu
BASHY_EXE="${BASHY:-bashy}"
: "${VSC_PCTS_HARNESS_DIR:?set VSC_PCTS_HARNESS_DIR to the licensed harness checkout}"
: "${VSC_DO_STATE_DIR:?set VSC_DO_STATE_DIR to approved DO2 state}"
: "${VSC_DO_STATE:?set VSC_DO_STATE to approved DO2 state name}"
: "${VSC_SHELL_ARM:?set VSC_SHELL_ARM to a new shell-only arm name}"
: "${VSC_PRIVATE_EVIDENCE_DIR:?set VSC_PRIVATE_EVIDENCE_DIR outside Git}"
: "${VSC_PUBLIC_EVIDENCE_DIR:?set VSC_PUBLIC_EVIDENCE_DIR outside Git}"
gate="$VSC_PCTS_HARNESS_DIR/scripts/release-posix-shell-gate.sh"
if [ ! -x "$gate" ]; then
  echo "test-posix-shell: licensed gate wrapper missing or not executable" >&2
  exit 1
fi
export BASHY_SOURCE_DIR="$(pwd -P)"
"$BASHY_EXE" dag build-sh
env "$gate"
```


### build-sh
Build the strict POSIX shell artifact with the released product Go overlay and native signal launcher.
Effects: exec, net, read, write

```bash
set -eu
export GOWORK=off
bin_dir="${BIN_DIR:-bin}"
mkdir -p "$bin_dir"
bin_dir="$(cd "$bin_dir" && pwd -P)"
BASHY_EXE="${BASHY:-bashy}"
go_cmd="${GO:-go}"
goos="$($go_cmd env GOOS)"
goarch="$($go_cmd env GOARCH)"
hostos="$($go_cmd env GOHOSTOS)"
hostarch="$($go_cmd env GOHOSTARCH)"
ext=
[ "$goos" != windows ] || ext=.exe
if [ "$goos" = "$hostos" ] && [ "$goarch" = "$hostarch" ] && [ "$(BASHY="$BASHY_EXE" "$BASHY_EXE" scripts/launcher-wanted.sh build-sh)" = 1 ]; then
  (cd cmd && BASHY_EXE= "$BASHY_EXE" ../scripts/go-product.sh build -mod=readonly -trimpath -o "$bin_dir/sh.real" ./sh)
  cc -x c -std=c11 -O2 -Wall -Wextra -Werror -o "$bin_dir/sh" native/siglaunch.c.in
else
  (cd cmd && BASHY_EXE= "$BASHY_EXE" ../scripts/go-product.sh build -mod=readonly -trimpath -o "$bin_dir/sh$ext" ./sh)
fi
```

