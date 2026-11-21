# bash

**A pure-Go GNU Bash 5.3 and POSIX.1-2016 userland that behaves the same on Linux, macOS and Windows, with no cgo and no host tools.**

It is the *one system* in [bashy](https://github.com/qiangli/bashy#bashy)'s pitch: one system (bash), one language ([Bash#](https://github.com/bashsharp/bashsharp)), one app (software shaping), one agent (genie).

This repository holds the interpreter source in its root module and pins the
userland, so the shell can be found, cloned and measured in one place:

| Part | What it is | Module |
|---|---|---|
| [`syntax/`](./syntax), [`expand/`](./expand), [`interp/`](./interp) and siblings | The interpreter: a fork of [mvdan/sh](https://github.com/mvdan/sh) carrying the Bash 5.3 and POSIX patches (parser, expansion, builtins, job control, traps), relocated from [`qiangli/sh`](https://github.com/qiangli/sh) | `github.com/qiangli/bash` |
| [`coreutils/`](https://github.com/qiangli/coreutils) (pinned submodule) | The userland: the POSIX-required utilities plus the GNU coreutils programs, as pure-Go applets | `github.com/qiangli/coreutils` |

`coreutils` stays independent: it builds, tests and versions on its own,
and this repository only pins a known-good revision.

## Built for POSIX certification; not certified

**No POSIX certification is claimed.** The work is aimed at certification
against POSIX.1-2016 (IEEE Std 1003.1). These are the measured results so far,
and each one names its corpus:

| Measure | Result | Corpus and evidence |
|---|---|---|
| GNU Bash 5.3's own test suite | 86/86 runnable fixtures on Linux, macOS and Windows; 86/86 in both dialect modes in the sealed Linux container gate (2026-11-10) | [claims ledger](https://github.com/bashsharp/bashsharp/blob/main/docs/claims.md) · `bashy dag test-bash-gate` in bashy |
| Required POSIX utility profile | 102 mandatory VSC utility targets served by 92 Go applets and 10 shell commands, with zero external utility providers | [coreutils README](https://github.com/qiangli/coreutils#posix-certification-boundary) |
| Integrated licensed run (S561, Profile D) | 9,270 of 9,337 PASS-group outcomes; 67 non-passes (4 FAIL, 54 UNRESOLVED, 9 INSPECT) | [status](https://github.com/qiangli/coreutils/blob/main/docs/posix-vsc-pcts-status.md) |

The remaining non-passes and the submission decisions are tracked as
certification-readiness work. A result applies to the exact candidate it
measured, not to every future commit.

## Use it

Most people get this shell inside [bashy](https://github.com/qiangli/bashy),
a single binary that is a drop-in `bash`/`sh` plus this userland:

```sh
curl -fsSL https://dhnt.io/install.sh | sh -s -- --no-service   # macOS / Linux
bashy --posix -c 'echo posix ok'
```

To embed it in a Go program, run scripts with this repository's `interp` package and serve
commands from `coreutils` through the interpreter's `ExecHandler`. The result is a
portable Unix-like environment inside one Go process. bashy's lean `cmd/bash` is
the reference wiring.

## Clone

```sh
git clone --recurse-submodules https://github.com/qiangli/bash
```

## Why it exists

An AI agent writes one script and expects it to mean the same thing everywhere.
A fixed, measured shell and userland is what makes that true. It is also the
base under [software shaping](https://github.com/qiangli/bashy/blob/main/docs/software-shaping.md):

> **Software shaping lets anyone tell an agent, in plain words, how a running tool should work, and get their own verified, reversible version, while everyone else keeps theirs.**

Shaping never changes this layer. Because it stays fixed, a check that verifies a shaped change means the same thing on every machine.

## License

This project is MIT licensed ([LICENSE](./LICENSE)). The in-root interpreter
tree keeps its BSD-3-Clause license ([LICENSE.sh](./LICENSE.sh)), and the
pinned `coreutils` submodule keeps its MIT license.
