package middleware

import (
	"context"
	"github.com/nheyn/go-redux/store"
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
		t.Error("The generated next function was called with the incorrect context")
	}
	if finalAction != testAction {
		t.Error("The generated next function was called with the incorrect actions")
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

func TestBaseNextCallsDispatchWithContextAndAction(t *testing.T) {
	var dispatchedCtx context.Context
	var dispatchedAction interface{}
	baseNext := createBaseNext(func(ctx context.Context, s store.State, action interface{}) (store.State, error) {
		dispatchedCtx = ctx
		dispatchedAction = action

		return s, nil
	}, nil, nil)

	testCtx := context.Background()
	testAction := "test action"
	err := baseNext(testCtx, testAction)
	if err != nil {
		t.Error(err)
	}

	if dispatchedCtx != testCtx {
		t.Error("The base next function was called with the incorrect context")
	}

	if dispatchedAction != testAction {
		t.Error("The base next function was called with the incorrect action")
	}
}

func TestBaseNextUpdatesTheState(t *testing.T) {
	testState := store.State{"testKey": testUpdater("testUpdater")}
	updatedState := store.State{}
	baseNext := createBaseNext(func(ctx context.Context, s store.State, action interface{}) (store.State, error) {
		return testState, nil
	}, store.State{}, &updatedState)

	err := baseNext(nil, nil)
	if err != nil {
		t.Error(err)
	}

	if len(testState) != len(updatedState) {
		t.Error("The updated state has", len(updatedState), "Updaters but should have", len(testState))
	}

	for key, testUpdater := range testState {
		currUpdater, exits := updatedState[key]
		if !exits {
			t.Error("The Updater with key", key, "is not in the updated state")
			continue
		}

		if currUpdater != testUpdater {
			t.Error("The Updater with key", key, "is no the updater that was set during dispatch")
		}
	}
}

type testUpdater string

func (t testUpdater) Update(_ context.Context, _ interface{}) (store.Updater, error) {
	return t, nil
}
