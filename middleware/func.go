package middleware

import "context"

// A middleware.Func is function that can act as middleware for a go-redux Store.
type Func func(context.Context, interface{}, Next) error

// Combines the given middleware Funcs to a single one.
func composeFuncs(mws ...Func) Func {
	getNextFor := createNextGenerator(mws...)

	return func(ctx context.Context, action interface{}, next Next) error {
		composedNext := getNextFor(next)

		return composedNext(ctx, action)
	}
}
