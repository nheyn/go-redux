package store

// A State is a map that contains the current data for a Store.
type State map[interface{}]Updater

// Selects a shallow copy of the the given state.
func (selSt *State) SelectFrom(currSt *State) {
	for key, data := range *currSt {
		(*selSt)[key] = data
	}
}

// Gets the updated version of the updaters in the given state, after the action is performed on them.
func performUpdates(s State, action interface{}) (State, error) {
	updateChan := make(chan keyedData, len(s))
	errChan := make(chan error, len(s))
	for key, data := range s {
		go performUpdate(keyedData{key, data}, action, updateChan, errChan)
	}

	newState := State{}
	for len(newState) != len(s) {
		select {
		case update := <-updateChan:
			newState[update.key] = update.data
		case err := <-errChan:
			return nil, err
		}
	}

	return newState, nil
}

// A data Updater from a State, with it's assigned key.
type keyedData struct {
	key  interface{}
	data Updater
}

// Gets the updated version of the given updater.
func performUpdate(in keyedData, action interface{}, outChan chan<- keyedData, errChan chan<- error) {
	newData, err := in.data.Update(action)
	if err != nil {
		errChan <- err
		return
	}

	outChan <- keyedData{key: in.key, data: newData}
}
