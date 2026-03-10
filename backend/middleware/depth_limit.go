package middleware

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// QueryDepthLimit is a gqlgen extension that rejects queries exceeding maxDepth
// levels of field nesting. This prevents infinite recursion through fields such
// as relatedMedia → [Media] → relatedMedia → ...
type QueryDepthLimit struct {
	MaxDepth int
}

var _ interface {
	graphql.HandlerExtension
	graphql.OperationInterceptor
} = QueryDepthLimit{}

func (QueryDepthLimit) ExtensionName() string                          { return "QueryDepthLimit" }
func (QueryDepthLimit) Validate(schema graphql.ExecutableSchema) error { return nil }

func (d QueryDepthLimit) InterceptOperation(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
	rc := graphql.GetOperationContext(ctx)
	if rc != nil && rc.Doc != nil {
		for _, op := range rc.Doc.Operations {
			if depth := selectionDepth(op.SelectionSet, 0); depth > d.MaxDepth {
				msg := fmt.Sprintf("query depth %d exceeds maximum allowed depth of %d", depth, d.MaxDepth)
				return func(_ context.Context) *graphql.Response {
					return &graphql.Response{
						Errors: gqlerror.List{{Message: msg}},
					}
				}
			}
		}
	}
	return next(ctx)
}

// selectionDepth recursively computes the maximum nesting depth of a selection set.
func selectionDepth(ss ast.SelectionSet, current int) int {
	if len(ss) == 0 {
		return current
	}
	max := current
	for _, sel := range ss {
		var childDepth int
		switch s := sel.(type) {
		case *ast.Field:
			childDepth = selectionDepth(s.SelectionSet, current+1)
		case *ast.InlineFragment:
			childDepth = selectionDepth(s.SelectionSet, current)
		case *ast.FragmentSpread:
			// Fragment spreads cannot be resolved without the full document; skip
			// (complexity limit handles abuse via fragment explosion separately).
			childDepth = current + 1
		}
		if childDepth > max {
			max = childDepth
		}
	}
	return max
}
