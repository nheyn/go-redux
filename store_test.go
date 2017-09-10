package store

import "testing"

func TestGivenInitialState(t *testing.T) {
	initialState := State{
		"Updater 0": testUpdater{"Updater 0", nil},
		"Updater 1": testUpdater{"Updater 1", nil},
		"Updater 2": testUpdater{"Updater 2", nil},
	}

	st := New(initialState)

	wait := make(chan struct{})
	st.accessState <- func(currState *State) {
		defer close(wait)

		if len(initialState) != len(*currState) {
			t.Error("State from the Store has ", len(initialState), " Updaters and the intial state has ", len(*currState))
		}

		for key, initialData := range initialState {
			s := *currState
			data, exists := s[key]
			if !exists {
				t.Error("Missing ", key, " the State in the Store")
			}

			if data.(testUpdater).name != initialData.(testUpdater).name {
				t.Error("Inncorrect data the State in the Store in Updater ", key)
			}
		}
	}
	<-wait
}

func TestStoreWillCallUpdate(t *testing.T) {
	state := State{
		"Updater 0": testUpdater{},
		"Updater 1": testUpdater{},
		"Updater 2": testUpdater{},
	}

	st := New(state)

	testAction := "Test action"
	err := st.Dispatch(testAction)
	if err != nil {
		t.Fatal(err)
		return
	}

	wait := make(chan struct{})
	st.accessState <- func(updatedState *State) {
		defer close(wait)

		for key, data := range *updatedState {
			testData := data.(testUpdater)

			if !testData.didUpdate() {
				t.Error("Update method not called on", key)
			}

			if testData.actions[0] != testAction {
				t.Error("Update method called with incorrect action, ", testData.actions[0], ", on ", key)
			}
		}
	}
	<-wait
}
