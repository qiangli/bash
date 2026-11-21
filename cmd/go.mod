module github.com/qiangli/bash/cmd

go 1.27.0

toolchain go1.27.4

require (
	github.com/kr/pretty v0.3.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	golang.org/x/sync v0.24.0 // indirect
)

tool golang.org/x/tools/cmd/stringer

replace github.com/ergochat/readline => github.com/qiangli/readline v0.1.4-0.20261013004858-df59d4e258e4

require (
	github.com/bashsharp/bashsharp v0.0.0-20261119110900-4141cc07cfdc
	github.com/creack/pty v1.1.24
	github.com/ergochat/readline v0.1.3
	github.com/go-quicktest/qt v1.101.0
	github.com/google/go-cmp v0.7.0
	github.com/google/renameio/v2 v2.0.2
	github.com/qiangli/bash v0.0.0-20261118061051-cdb1212b069f
	github.com/rogpeppe/go-internal v1.14.1
	golang.org/x/mod v0.41.0
	golang.org/x/sys v0.49.0
	golang.org/x/term v0.47.0
	golang.org/x/text v0.43.0
	golang.org/x/tools v0.51.0
	mvdan.cc/editorconfig v0.3.0
)

replace github.com/qiangli/bash => ..

replace golang.org/x/sys v0.49.0 => golang.org/x/sys v0.48.0

replace golang.org/x/term v0.47.0 => golang.org/x/term v0.46.0

replace golang.org/x/text v0.43.0 => golang.org/x/text v0.42.0
