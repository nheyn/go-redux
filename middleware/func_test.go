package middleware

import (
	"context"
	"testing"
)

func TestComposeWillPassThroughItArguments(t *testing.T) {
	ctxs := []context.Context{}
	actions := []interface{}{}
	mws := []Func{}
	for i := 0; i < 3; i++ {
		mw := func(ctx context.Context, action interface{}, next Next) error {
			ctxs = append(ctxs, ctx)
			actions = append(actions, action)

			return next(ctx, action)
		}

		mws = append(mws, mw)
	}

	testCtx := context.Background()
	testAction := "test action"

	var finalCtx context.Context
	var finalAction interface{}
	composedMw := composeFuncs(mws...)
	composedMw(testCtx, testAction, func(ctx context.Context, action interface{}) error {
		finalCtx = ctx
		finalAction = action

		return nil
	})

	if len(ctxs) != len(mws) {
		t.Error(len(mws), "contexts should have been saved, but", len(ctxs), "where")
	}
	if len(actions) != len(mws) {
		t.Error(len(mws), "actions should have been saved, but", len(actions), "where")
	}

	for i, ctx := range ctxs {
		if ctx != testCtx {
			t.Error("The middleware.Func called at", i, "was called with the incorrect context")
		}
	}
	for i, action := range actions {
		if action != testAction {
			t.Error("The middleware.Func called at", i, "was called with the incorrect action")
		}
	}

	if finalCtx != testCtx {
		t.Error("The composed Func was called with the incorrect context")
	}
	if finalAction != testAction {
		t.Error("The composed Func was called with the incorrect actions")
	}
}

func TestComposeWillCallMiddlewareInOrder(t *testing.T) {
	called := []int{}
	mws := []Func{}
	for _i := 0; _i < 3; _i++ {
		i := _i
		mw := func(ctx context.Context, action interface{}, next Next) error {
			called = append(called, i)

			return next(ctx, action)
		}

		mws = append(mws, mw)
	}

	composedMw := composeFuncs(mws...)
	composedMw(nil, nil, func(_ context.Context, _ interface{}) error {
		return nil
	})

	if len(called) != len(mws) {
		t.Error(len(mws), "middleware.Funcs should have been called, but", len(called), "where")
	}

	for i, callIndex := range called {
		if callIndex != i {
			t.Error("The middleware.Func called at", i, ", but should have been called at", callIndex)
		}
	}
}
