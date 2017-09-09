package store

import "testing"

type testUpdater struct {
	didUpdate bool
	action    interface{}
}

func (u testUpdater) Update(action interface{}) (Updater, error) {
	return testUpdater{true, action}, nil
}

func TestStateWillCallUpdate(t *testing.T) {
	state := State{
		"Updater 0": testUpdater{false, nil},
		"Updater 1": testUpdater{false, nil},
		"Updater 2": testUpdater{false, nil},
	}

	testAction := "Test action"
	updatedState, err := performUpdates(state, testAction)
	if err != nil {
		t.Fatal(err)
		return
	}

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
