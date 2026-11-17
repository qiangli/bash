package shellexec

import (
	"context"

	"github.com/qiangli/bash/interp"
	"github.com/qiangli/bash/lower/shellrt"
	"github.com/qiangli/bash/syntax"
)

func (sh *shell) taskPolicyContext(ctx context.Context) context.Context {
	if sh.cfg.lang == syntax.LangBashPP && shellrt.InTask(ctx) {
		return interp.WithTaskPolicy(ctx)
	}
	return ctx
}
