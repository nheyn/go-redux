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

// A store config that will use make the Store use the default state update function.
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

// Creates a function that will get the updated version of the given updater.
func getPerformUpdateFor(
	ctx context.Context,
	action interface{},
	outChan chan<- keyedData,
	errChan chan<- error,
) func(keyedData) {
	return func(in keyedData) {
		ctxWithKey := context.WithValue(ctx, keyKey, in.key)

		newData, err := in.data.Update(ctxWithKey, action)
		if err != nil {
			errChan <- err
			return
		}

		outChan <- keyedData{key: in.key, data: newData}
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

// KeyFrom, when called in .Update(...) method of an Updater, will get the key that was the
// current Updater was registered to the store under.
func KeyFrom(ctx context.Context) (interface{}, bool) {
	key := ctx.Value(keyKey)
	if key == nil {
		return nil, false
	}

	return key, true
}
