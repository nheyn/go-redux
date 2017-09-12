package store

import (
	"context"
	"testing"
)

func TestStateWillCallUpdate(t *testing.T) {
	state := State{
		"Updater 0": testUpdater{},
		"Updater 1": testUpdater{},
		"Updater 2": testUpdater{},
	}

	testAction := "Test action"
	updatedState, err := performUpdates(context.Background(), state, testAction)
	if err != nil {
		t.Fatal(err)
		return
	}

	for key, data := range updatedState {
		testData := data.(testUpdater)

		if !testData.didUpdate() {
			t.Error("Update method not called on", key)
			continue
		}

		if testData.actions[0] != testAction {
			t.Error("Update method called with incorrect action, ", testData.actions[0], ", on ", key)
		}
	}
}

type testUpdater struct {
	name    string
	actions []interface{}
}

func (u testUpdater) Update(_ context.Context, action interface{}) (Updater, error) {
	var newActions []interface{}
	if u.actions == nil {
		newActions = []interface{}{action}
	} else {
		newActions = append(u.actions, action)
	}

	return testUpdater{u.name, newActions}, nil
}

func (u testUpdater) didUpdate() bool {
	return u.actions != nil && len(u.actions) > 0
}
