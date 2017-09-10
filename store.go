package store

// A Store keeps track of data in a State, and "attempts to make state mutations predictable".
type Store struct {
	accessState chan func(*State)
}

// Creates a new Store that start with the given state.
func New(initialState State) Store {
	s := Store{
		accessState: make(chan func(*State)),
	}

	go s.trackState(initialState)

	return s
}

// A method that will keep track of the sate, which can only be accessed throught the accessState
//  channel. The data in the given State will be put into the tracked State to start.
func (s Store) trackState(initialState State) {
	currState := State{}
	for key, data := range initialState {
		currState[key] = data
	}

	for accessFn := range s.accessState {
		accessFn(&currState)
	}
}

// Dispatches the given action to all of the Updaters in the state of the Store. If an error is
// returned, then the State will not not change (even for the Updaters that had already completed).
func (s Store) Dispatch(action interface{}) error {
	errChan := make(chan error)

	s.accessState <- func(st *State) {
		newState, err := performUpdates(*st, action)
		if err != nil {
			errChan <- err
			return
		}

		for key, data := range newState {
			(*st)[key] = data
		}

		close(errChan)
	}

	return <-errChan
}
