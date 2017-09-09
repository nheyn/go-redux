package store

// A Store keeps track of data in a State, and "attempts to make state mutations predictable".
type Store struct {
	State
}
