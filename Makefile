###########################################################################
##  GENERATED FILE — DO NOT EDIT.                                        ##
##  Produced from dag.md by `bashy dag --makefile`. Edit dag.md instead, ##
##  then regenerate; manual changes are overwritten and fail --check.    ##
###########################################################################
# Tasks live in dag.md (a bashy dag file); this Makefile only forwards.
#   bashy dag <target>           run a target (its dependencies first)
#   bashy dag --list             list targets
#   bashy dag --dry-run <target> print the plan, run nothing
#   bashy dag --skill DIR        export these tasks as an Agent Skill
# Regenerate with: bashy dag --makefile  (verify: bashy dag --makefile --check)

# List the dag targets, then how to run them natively.
help: FORCE
	@command -v bashy >/dev/null || { echo "requires bashy: https://github.com/qiangli/bashy"; exit 127; }; bashy dag --list; echo "Run: bashy dag <target> | skill: bashy dag --skill DIR | more: bashy help"

# Build the library and the bash and shfmt commands.
build: FORCE
	@command -v bashy >/dev/null || { echo "requires bashy: https://github.com/qiangli/bashy"; exit 127; }; bashy dag build

# Run the full suite.
test: FORCE
	@command -v bashy >/dev/null || { echo "requires bashy: https://github.com/qiangli/bashy"; exit 127; }; bashy dag test

# The quick tier — what push CI runs: build + unit/mock tests (`go test -short`, without the `//go:build full` evaluator files, < 3 min).
test-quick: FORCE
	@command -v bashy >/dev/null || { echo "requires bashy: https://github.com/qiangli/bashy"; exit 127; }; bashy dag test-quick

# `moreinterp/` is a separate Go module (u-root-backed coreutils ExecHandler); test it independently.
test-moreinterp: FORCE
	@command -v bashy >/dev/null || { echo "requires bashy: https://github.com/qiangli/bashy"; exit 127; }; bashy dag test-moreinterp

# Race detector over the suite (CI runs this on Linux).
test-race: FORCE
	@command -v bashy >/dev/null || { echo "requires bashy: https://github.com/qiangli/bashy"; exit 127; }; bashy dag test-race

# Bash++ race/lifecycle gate.
bashpp-race-gate: FORCE
	@command -v bashy >/dev/null || { echo "requires bashy: https://github.com/qiangli/bashy"; exit 127; }; bashy dag bashpp-race-gate

# The canonical "behaves like real Bash 5.3" oracle: the same script table is run by the interpreter and by a dockerized Bash, and outputs/exit codes are diffed.
confirm: FORCE
	@command -v bashy >/dev/null || { echo "requires bashy: https://github.com/qiangli/bashy"; exit 127; }; bashy dag confirm

# Static checks: gofmt + go vet.
vet: FORCE
	@command -v bashy >/dev/null || { echo "requires bashy: https://github.com/qiangli/bashy"; exit 127; }; bashy dag vet

# go mod tidy + gofmt -s -w .
tidy: FORCE
	@command -v bashy >/dev/null || { echo "requires bashy: https://github.com/qiangli/bashy"; exit 127; }; bashy dag tidy

# go clean.
clean: FORCE
	@command -v bashy >/dev/null || { echo "requires bashy: https://github.com/qiangli/bashy"; exit 127; }; bashy dag clean

# The formatting gate (reports unformatted files, never rewrites them), plus a check that the generated Makefile still matches this file — regenerate with `bashy…
fmtcheck: FORCE
	@command -v bashy >/dev/null || { echo "requires bashy: https://github.com/qiangli/bashy"; exit 127; }; bashy dag fmtcheck

# Install the pre-push formatting gate.
hooks: FORCE
	@command -v bashy >/dev/null || { echo "requires bashy: https://github.com/qiangli/bashy"; exit 127; }; bashy dag hooks

# Build the Bash artifact with the released product Go overlay and native signal launcher.
build-bash: FORCE
	@command -v bashy >/dev/null || { echo "requires bashy: https://github.com/qiangli/bashy"; exit 127; }; bashy dag build-bash

# Run the preserved Go regression suite.
test-go: FORCE
	@command -v bashy >/dev/null || { echo "requires bashy: https://github.com/qiangli/bashy"; exit 127; }; bashy dag test-go

# Run the existing full-tag regression suite without changing test selections.
test-full: FORCE
	@command -v bashy >/dev/null || { echo "requires bashy: https://github.com/qiangli/bashy"; exit 127; }; bashy dag test-full

# Smoke-test the actual built Bash artifact.
smoke: FORCE
	@command -v bashy >/dev/null || { echo "requires bashy: https://github.com/qiangli/bashy"; exit 127; }; bashy dag smoke

# Build this checkout in the locked official image.
official-build: FORCE
	@command -v bashy >/dev/null || { echo "requires bashy: https://github.com/qiangli/bashy"; exit 127; }; bashy dag official-build

# Run preserved Go regressions in the same locked official image.
official-test: FORCE
	@command -v bashy >/dev/null || { echo "requires bashy: https://github.com/qiangli/bashy"; exit 127; }; bashy dag official-test

# Fetch and SHA-256 verify the unchanged GNU Bash 5.3 fixtures (gitignored).
test-bash-fixtures: FORCE
	@command -v bashy >/dev/null || { echo "requires bashy: https://github.com/qiangli/bashy"; exit 127; }; bashy dag test-bash-fixtures

# Sealed GNU Bash 5.3 Classic gate: exactly the existing 86 named fixtures, mode 0.
test-bash-gate: FORCE
	@command -v bashy >/dev/null || { echo "requires bashy: https://github.com/qiangli/bashy"; exit 127; }; bashy dag test-bash-gate

# Build this repository Dockerfile after verifying the locked base image digest.
docker-build: FORCE
	@command -v bashy >/dev/null || { echo "requires bashy: https://github.com/qiangli/bashy"; exit 127; }; bashy dag docker-build

# Run the licensed POSIX Profile B **shell-only** release gate, explicitly and only from a clean, pushed Bashy checkout at the supervisor-approved commit.
test-posix-shell: FORCE
	@command -v bashy >/dev/null || { echo "requires bashy: https://github.com/qiangli/bashy"; exit 127; }; bashy dag test-posix-shell

# Build the strict POSIX shell artifact with the released product Go overlay and native signal launcher.
build-sh: FORCE
	@command -v bashy >/dev/null || { echo "requires bashy: https://github.com/qiangli/bashy"; exit 127; }; bashy dag build-sh

.DEFAULT: ; @command -v bashy >/dev/null || { echo "requires bashy: https://github.com/qiangli/bashy"; exit 127; }; bashy dag $@

# FORCE names no file and has no recipe: as the one prerequisite it
# keeps every forward running even when a file or directory on disk
# shares a target's name. Dependency order stays in dag, not here.
FORCE:
