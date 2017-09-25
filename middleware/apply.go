package middleware

import (
	"context"
	"github.com/nheyn/go-redux/store"
)

// A middleware.Func is function that can act as middleware for a go-redux Store.
type Func func(context.Context, interface{}, Next) error

// Apply returns the configuration function that can be passed to store.New(...). It will wrap
// the Store's dispatch calls with the given middleware.
// Ex)
//	someMiddleware := func(s *Store) middleware.Func {
//		return func(ctx context.Context, action interface{}, next middleware.Next) error {
//			// ...middleware code goes here...
//
//			return next(ctx, action)
//		}
//	}
func Apply(mwGens ...func(s *store.Store) Func) func(*store.Store) {
	return func(s *store.Store) {
		//TODO, replace the .PerformDispatch with one that calls the given middleware
	}
}
