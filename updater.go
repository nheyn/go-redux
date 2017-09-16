package store

import "context"

// A value that can create an updated version of itself based on a given action.
// NOTE: All Updaters should be immutable
type Updater interface {

	// Creates a new Updater with the updates for the given action. The given context contains the
	// key that was used to register this Updater in a State. It can be accessed key using the
	// store.KeyFrom(...) function.
	// NOTE: Because Updaters should be immutable, always create a new version when an update occures.
	Update(ctx context.Context, ac interface{}) (Updater, error)
}
