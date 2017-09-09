package store

// A Store keeps track of data in a State, and "attempts to make state mutations predictable".
type Store struct {
	state *State
}

func New(initialState State) Store {
	st := Store{
		state: &initialState,
	}

	return st
}
