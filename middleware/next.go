package middleware

import (
	"context"
	"github.com/nheyn/go-redux/store"
)

// A middleware.Func is function that can act as middleware for a go-redux Store.
type Next func(context.Context, interface{}) error

// Creates a function that will generate a Next function that will call the given
// middleware Func.
func createNextGenerator(mws ...Func) func(Next) Next {
	if len(mws) == 0 {
		return func(next Next) Next { return next }
	}
	currMw := mws[0]
	getCurrNext := func(next Next) Next {
		return func(ctx context.Context, action interface{}) error {
			return currMw(ctx, action, next)
		}
	}

	if len(mws) == 1 {
		return getCurrNext
	}
	remainingMws := mws[1:]
	getRemainingNext := createNextGenerator(remainingMws...)

	return func(next Next) Next {
		return getCurrNext(getRemainingNext(next))
	}
}

// Returns a Next function that will perform the default dispatch(...) behavoir when called by
// a middleware Func. The new will be saved in updatedState after the returned next func is
// called, unless it returned an error.
func createBaseNext(dispatch store.PerformDispatch, st store.State, updatedSt *store.State) Next {
	return func(ctx context.Context, action interface{}) error {
		newSt, err := dispatch(ctx, st, action)
		if err != nil {
			return err
		}

		updatedSt.SelectFrom(&newSt)
		return nil
	}
}
