package store

// A value that can create an updated version of itself based on a given action.
// NOTE: All Updaters should be immutable
type Updater interface {

	// Creates a new Updater with the updates for the given action.
	// NOTE: Because Updaters should be immutable, always create a new version when an update occures.
	Update(ac interface{}) (Updater, error)
}

// A State is a map that contains the current data for a Store.
type State map[interface{}]Updater
