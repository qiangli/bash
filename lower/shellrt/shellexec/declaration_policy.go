package shellexec

import (
	"context"
	"github.com/qiangli/bash/interp"
	"github.com/qiangli/bash/lower/shellrt"
	"github.com/qiangli/bash/syntax"
)

func (sh *shell) declarationContext(ctx context.Context) context.Context {
	p := shellrt.DeclarationsFromContext(ctx)
	if p == nil || sh.cfg.lang != syntax.LangBashPP {
		return ctx
	}
	declarations := make([]interp.BashPPDeclaration, 0, len(p.Declarations))
	for _, d := range p.Declarations {
		declarations = append(declarations, interp.BashPPDeclaration{ID: d.ID, Name: d.Name, Type: d.Type, Constant: d.Constant})
	}
	return interp.WithDeclarations(ctx, declarations, p.RefuseAssignment)
}
