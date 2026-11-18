# Standalone Bash build

The source checkout needs no umbrella or sibling repositories. Go resolves the
fixed public dependencies in `go.mod` with `GOWORK=off`.

The library root retains its original dependency profile. Existing command
sources live in the separate `cmd` module, which retains the CLI's Readline
profile and replaces only the Bash library with its parent directory in this
same checkout. Build and regression DAG targets cover both modules. This keeps
the library and command dependency selections distinct without changing either
implementation.

The existing Bashy executable and OCI backend drive these targets:

```sh
bashy dag official-build     # bin/bash and native bin/bash.real, then smoke
bashy dag official-test      # preserved default Go regression suite
bashy dag test-bash-gate     # sealed GNU Bash 5.3 Classic fixture gate
bashy dag docker-build      # repository Dockerfile, locked base identity
```

`image.lock` fixes the existing official image and digest. `run.sh` verifies the
local image against that identity before building. Provision this image on the
runner first; the target does not substitute another image. The Dockerfile uses
the same image tag; `docker-build` verifies its digest before invoking the build.

Build support (`native/siglaunch.c.in`, `tools/productgo/main.go`, and
`scripts/go-product.sh`) is retained from the released artifact recipe. It keeps
the pre-Go signal launcher and runtime overlay rather than changing shell code.

The copied sealed GNU harness lives in separate tool modules under
`tools/bash53suite` and `tools/bash53evidence`. Their test-tool dependencies do not
enter the Bash product module. The fixture download uses the released GNU Bash
5.3 URL and SHA-256 and remains gitignored. `test-bash-gate` runs artifact A in
Classic mode 0 with the unchanged named 86-fixture inventory. It is Bash
compatibility coverage, **not POSIX certification**.

`test-posix-shell` retains the existing, separately configured licensed VSC shell
gate. It requires the private harness, its existing supervisor-issued approval,
and the strict `cmd/sh` artifact. Licensed fixtures and private approval data
remain outside this public repository. Missing configuration fails closed.

The default Bash build cache is `bash-build-cache`; container lane labels default
to `bash-classic`. `BASHY_BUILD_CPUS` defaults to four and can be set by the runner
to share its capacity with independent Bash# builds.
