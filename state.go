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
	performUpdate := getPerformUpdateFor(action, updateChan, errChan)

	for key, data := range s {
		go performUpdate(keyedData{key, data})
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

// Creates a function that will get the updated version of the given updater.
func getPerformUpdateFor(action interface{}, outChan chan<- keyedData, errChan chan<- error) func(keyedData) {
	return func(in keyedData) {
		newData, err := in.data.Update(action)
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
