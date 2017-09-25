package middleware

import (
	"context"
	"testing"
)

func TestGeneratedNextWillUseGivenNext(t *testing.T) {
	createNext := createNextGenerator(func(ctx context.Context, action interface{}, next Next) error {
		return next(ctx, action)
	})

	calledNext := false
	currNext := createNext(func(_ context.Context, _ interface{}) error {
		calledNext = true

		return nil
	})
	currNext(nil, nil)

	if !calledNext {
		t.Error("The Next func passed to the the generator function was not called")
	}
}

func TestGeneratedNextWillPassThroughItArguments(t *testing.T) {
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

	var finalCtx context.Context
	var finalAction interface{}
	createNext := createNextGenerator(mws...)
	currNext := createNext(func(ctx context.Context, action interface{}) error {
		finalCtx = ctx
		finalAction = action

		return nil
	})

	testCtx := context.Background()
	testAction := "test action"
	currNext(testCtx, testAction)

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
		t.Error("The base next function was called with the incorrect context")
	}
	if finalAction != testAction {
		t.Error("The base next function was called with the incorrect actions")
	}
}

func TestGeneratedNextWillCallMiddlewareInOrder(t *testing.T) {
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

	createNext := createNextGenerator(mws...)
	currNext := createNext(func(_ context.Context, _ interface{}) error {
		return nil
	})
	currNext(nil, nil)

	if len(called) != len(mws) {
		t.Error(len(mws), "middleware.Funcs should have been called, but", len(called), "where")
	}

	for i, callIndex := range called {
		if callIndex != i {
			t.Error("The middleware.Func called at", i, ", but should have been called at", callIndex)
		}
	}
}
