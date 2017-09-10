package store

// A Selector can select values out of a State.
type Selector interface {
	// Pulls the required data from the given State into the Selector.
	// NOTE: DO NOT mutate the State in this method, only read from it.
	SelectFrom(st *State)
}
