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

// Selects a shallow copy of the the given state.
func (selSt *State) SelectFrom(currSt *State) {
	for key, data := range *currSt {
		(*selSt)[key] = data
	}
}

// Gets the update version of the updaters in the given state, after the action is performed on them.
func performUpdates(s State, action interface{}) (State, error) {
	newState := State{}
	for name, data := range s {
		newData, err := data.Update(action)
		if err != nil {
			return nil, err
		}

		newState[name] = newData
	}

	return newState, nil
}
