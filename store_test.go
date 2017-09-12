package store

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestGivenInitialState(t *testing.T) {
	initialState := State{
		"Updater 0": testUpdater{"Updater 0", nil},
		"Updater 1": testUpdater{"Updater 1", nil},
		"Updater 2": testUpdater{"Updater 2", nil},
	}

	st := New(initialState)

	currState := State{}
	st.Select(&currState)

	if len(initialState) != len(currState) {
		t.Error("State from the Store has ", len(initialState), " Updaters and the intial state has ", len(currState))
	}

	for key, initialData := range initialState {
		data, exists := currState[key]
		if !exists {
			t.Error("Missing ", key, " the State in the Store")
		}

		if data.(testUpdater).name != initialData.(testUpdater).name {
			t.Error("Inncorrect data the State in the Store in Updater ", key)
		}
	}
}

func TestStoreWillCallUpdate(t *testing.T) {
	state := State{
		"Updater 0": testUpdater{},
		"Updater 1": testUpdater{},
		"Updater 2": testUpdater{},
	}

	st := New(state)

	testAction := "Test action"
	err := st.Dispatch(context.Background(), testAction)
	if err != nil {
		t.Fatal(err)
		return
	}

	updatedState := State{}
	st.Select(&updatedState)

	for key, data := range updatedState {
		testData := data.(testUpdater)

		if !testData.didUpdate() {
			t.Error("Update method not called on", key)
		}

		if testData.actions[0] != testAction {
			t.Error("Update method called with incorrect action, ", testData.actions[0], ", on ", key)
		}
	}
}

func TestStoreWillQueueActions(t *testing.T) {
	state := State{
		"Updater 0": testUpdater{},
		"Updater 1": testUpdater{},
		"Updater 2": testUpdater{},
	}

	st := New(state)
	senders := 3
	actionPerSender := 10

	var wg sync.WaitGroup
	wg.Add(senders)
	for i := 0; i < senders; i++ {
		go func(groupdId int) {
			defer wg.Done()

			for j := 0; j < actionPerSender; j++ {
				st.Dispatch(context.Background(), groupdId)

				time.Sleep(time.Duration(groupdId) * time.Millisecond)
			}
		}(i)
	}

	wg.Wait()

	currState := State{}
	st.Select(&currState)

	for key, data := range currState {
		testData := data.(testUpdater)

		if !testData.didUpdate() {
			t.Error("Update method not called on", key)
			continue
		}

		if len(testData.actions) != senders*actionPerSender {
			t.Error(
				"Update method should have been called",
				senders*actionPerSender,
				"but was called",
				len(testData.actions),
				"on",
				key,
			)
		}
	}
}

func TestStoreWillTrackSubscribers(t *testing.T) {
	state := State{
		"Updater 0": testUpdater{},
		"Updater 1": testUpdater{},
		"Updater 2": testUpdater{},
	}

	testSubs := []subscriber{
		make(chan Store),
		make(chan Store),
		make(chan Store),
	}

	st := New(state)

	unsubs := []func() bool{}
	for _, testSub := range testSubs {
		unsub := st.Subscribe(testSub)

		unsubs = append(unsubs, unsub)
	}

	done := make(chan struct{})
	st.accessSubscribers <- func(subs *subscriberSet) {
		defer close(done)

		for i, testSub := range testSubs {
			if _, hasSub := (*subs)[testSub]; !hasSub {
				t.Error("The subscriber, at index", i, "did not subscribe")
			}
		}
	}

	for i, unsub := range unsubs {
		if didUnsub := unsub(); !didUnsub {
			t.Error("The subscriber, at index", i, "did not unsubscribe correctly")
		}
	}

	done = make(chan struct{})
	st.accessSubscribers <- func(subs *subscriberSet) {
		defer close(done)

		for i, testSub := range testSubs {
			if _, hasSub := (*subs)[testSub]; hasSub {
				t.Error("The subscriber, at index", i, "did not unsubscribe")
			}
		}
	}

	for i, unsub := range unsubs {
		if didUnsub := unsub(); didUnsub {
			t.Error("The subscriber, at index", i, "did not fail when unsubscribe more then once")
		}
	}
}

func TestStoreWillUpdateSubscribers(t *testing.T) {
	state := State{
		"Updater 0": testUpdater{},
		"Updater 1": testUpdater{},
		"Updater 2": testUpdater{},
	}

	st := New(state)
	subscribers := 3
	actionPerSubscriber := 10

	var wg sync.WaitGroup
	wg.Add(subscribers)
	for i := 0; i < subscribers; i++ {
		testSub := make(chan Store, actionPerSubscriber)
		st.Subscribe(testSub)

		go func() {
			defer wg.Done()

			for j := 0; j < actionPerSubscriber; j++ {
				currStore := <-testSub
				currState := State{}
				currStore.Select(&currState)

				for key, data := range currState {
					testData := data.(testUpdater)

					if len(testData.actions) != j+1 {
						t.Error(
							"After", j+1, "dispatches a subscriber",
							"(in", key, "at index", i, ")",
							"it already has", len(testData.actions), "actions",
						)
						t.Error(testData.actions)
					}

					if val, isInt := testData.actions[j].(int); !isInt || val != j {
						t.Error(
							"After", j+1, "dispatches a subscriber",
							"(in", key, "at index", i, ")",
							"has the incorrect action:", testData.actions[j],
						)
					}
				}
			}
		}()
	}

	for j := 0; j < actionPerSubscriber; j++ {
		st.Dispatch(context.Background(), j)

		//NOTE, if done without .Sleep(...), a tests action might be missed by the test code
		//TODO, add a way to guarantee .Select(...) can get every state
		time.Sleep(time.Millisecond)
	}

	wg.Wait()
}
