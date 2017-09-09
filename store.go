package store

// A Store keeps track of data in a State, and "attempts to make state mutations predictable".
type Store struct {
	state State
}

// Creates a new Store that start with the given state.
func New(initialState State) *Store {
	st := &Store{
		state: initialState,
	}

	return st
}

// Dispatches the given action to all of the Updaters in the state of the Store. If an error is
// returned, then the State will not not change (even for the Updaters that had already completed).
func (s *Store) Dispatch(action interface{}) error {
	newState, err := performUpdates(s.state, action)
	if err != nil {
		return err
	}

	s.state = newState
	return nil
}
