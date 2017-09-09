package store

import "testing"

func TestGivenInitialState(t *testing.T) {
	initialState := State{
		"Updater 0": testUpdater{},
		"Updater 1": testUpdater{},
		"Updater 2": testUpdater{},
	}

	st := New(initialState)
	currState := st.state

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

func TestStoreWillCallUpdate(t *testing.T) {
	state := State{
		"Updater 0": testUpdater{false, nil},
		"Updater 1": testUpdater{false, nil},
		"Updater 2": testUpdater{false, nil},
	}

	st := New(state)

	testAction := "Test action"
	err := st.Dispatch(testAction)
	if err != nil {
		t.Fatal(err)
		return
	}
	updatedState := st.state

	for key, data := range updatedState {
		testData := data.(testUpdater)

		if !testData.didUpdate {
			t.Error("Update method not called on", key)
		}

		if testData.action != testAction {
			t.Error("Update method called with incorrect action, ", testData.action, ", on ", key)
		}
	}
}
