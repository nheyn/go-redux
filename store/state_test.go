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

func TestStateCanErrorDuringUpdate(t *testing.T) {
	state := State{
		"Updater 0": testUpdater{},
		"Updater 1": testUpdaterError{},
		"Updater 2": testUpdater{},
	}

	testAction := "Test action"
	updatedState, err := performUpdates(context.Background(), state, testAction)
	if err == nil {
		t.Fatal("Updater should have retuned an error")
		return
	}

	if updatedState != nil {
		t.Fatal("Updater should not have retuned an updater")
		return
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

type testUpdaterError struct {
	errorOn int
	actions []interface{}
}

func (u testUpdaterError) Update(_ context.Context, action interface{}) (Updater, error) {
	var newActions []interface{}
	if u.actions == nil {
		newActions = []interface{}{action}
	} else {
		newActions = append(u.actions, action)
	}

	newU := testUpdaterError{u.errorOn, newActions}
	if len(newU.actions)-1 == u.errorOn {
		return nil, u
	}

	return u, nil
}

func (u testUpdaterError) Error() string {
	return "The testUpdaterError correctly return this error"
}
