package store

import "context"

// A State is a map that contains the current data for a Store.
type State map[interface{}]Updater

// Selects a shallow copy of the the given state.
func (selSt *State) SelectFrom(currSt *State) {
	for key, data := range *currSt {
		(*selSt)[key] = data
	}
}

// A config function that will use make the given Store use the default state update function, which
// is returned from getPerformUpdateFor(...).
func defaultPeformDispatchConfig(s *Store) {
	s.PerformDispatch = func(ctx context.Context, st State, action interface{}) (State, error) {
		cancelableCtx, cancel := context.WithCancel(ctx)
		defer cancel()

		updateChan := make(chan keyedData, len(st))
		errChan := make(chan error, len(st))

		performUpdate := getPerformUpdateFor(cancelableCtx, action, updateChan, errChan)

		for key, data := range st {
			go performUpdate(keyedData{key, data})
		}

		newState := State{}
		for len(newState) != len(st) {
			select {
			case update := <-updateChan:
				newState[update.key] = update.data
			case <-ctx.Done():
				return nil, ctx.Err()
			case err := <-errChan:
				return nil, err
			}
		}

		return newState, nil
	}
}

// Creates a function that will pefrom the update for the given action with the given context. The
// given channeles will return all data and/or errors, so the returned function should be called
// on a seperate goroutine.
func getPerformUpdateFor(
	ctx context.Context,
	action interface{},
	updateChan chan<- keyedData,
	errChan chan<- error,
) func(keyedData) {
	return func(inital keyedData) {
		ctxWithKey := contextWithKey(ctx, inital.key)

		updatedData, err := inital.data.Update(ctxWithKey, action)
		if err != nil {
			errChan <- err
			return
		}

		updateChan <- keyedData{key: inital.key, data: updatedData}
	}
}

// A data Updater from a State, with it's assigned key.
type keyedData struct {
	key  interface{}
	data Updater
}

// The key for a Value that should be added to the ctx
type contextKey int

const keyKey contextKey = 0

// Will add the given key to the context, so it can be accessed by KeyFrom(...).
func contextWithKey(ctx context.Context, key interface{}) context.Context {
	return context.WithValue(ctx, keyKey, key)
}

// When called in the .Update(...) method of an Updater, it will get the Store key for the
// current Updater.
func KeyFrom(ctx context.Context) (interface{}, bool) {
	key := ctx.Value(keyKey)
	if key == nil {
		return nil, false
	}

	return key, true
}
