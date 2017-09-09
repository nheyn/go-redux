package store

import "testing"

func TestGivenInitialState(t *testing.T) {
	initialState := State{
		"Updater 0": testUpdater{},
		"Updater 1": testUpdater{},
		"Updater 2": testUpdater{},
	}

	st := New(initialState)
	currState := *st.state

	if len(initialState) != len(currState) {
		t.Error("State from the Store has ", len(initialState), " Updaters and the intial state has ", len(currState))
	}

	for key, initialData := range initialState {
		data, exists := currState[key]
		if !exists {
			t.Error("Missing ", key, " the State in the Store")
		}

		if data != initialData {
			t.Error("Inncorrect data the State in the Store in Updater ", key)
		}
	}
}
