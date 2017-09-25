package store

import (
	"sync"
	"testing"
)

func TestSubscriberSetCanSubscribe(t *testing.T) {
	testSubs := []subscriber{
		make(chan *Store),
		make(chan *Store),
		make(chan *Store),
	}

	subs := &subscriberSet{}
	for _, sub := range testSubs {
		subs.add(sub)
	}

	if len(*subs) != len(testSubs) {
		t.Error("Store should have", len(testSubs), "subs but has", len(*subs))
	}

	for i, testSub := range testSubs {
		_, hasSubs := (*subs)[testSub]

		if !hasSubs {
			t.Error("Store subs missing, at index", i)
			break
		}
	}
}

func TestSubscriberSetCanUnsubscribe(t *testing.T) {
	testSubs := []subscriber{
		make(chan *Store),
		make(chan *Store),
		make(chan *Store),
	}

	subs := &subscriberSet{}
	for _, sub := range testSubs {
		subs.add(sub)
	}

	for i, testSub := range testSubs {
		if ok := subs.remove(testSub); !ok {
			t.Error("Subscribers was unable to unsubscribe, at index", i)
		}

		if _, hasSub := (*subs)[testSub]; hasSub {
			t.Error("Subscribers remaining after unsubscribed, at index", i)
			break
		}

		if ok := subs.remove(testSub); ok {
			t.Error("Subscribers was `able` to unsubscribe even though the subscriber was already removed, at index", i)
		}
	}
}

func TestSubscriberSetCanPublish(t *testing.T) {
	state := State{
		"Updater 0": testUpdater{"Updater 0", nil},
	}

	subs := &subscriberSet{}
	subscribers := 3
	updatesPreSubscribers := 10

	var wg sync.WaitGroup
	wg.Add(subscribers)
	for i := 0; i < subscribers; i++ {
		sub := make(chan *Store)
		subs.add(sub)

		go func() {
			defer wg.Done()

			for j := 0; j < updatesPreSubscribers; j++ {
				currStore := <-sub
				currState := State{}
				currStore.Select(&currState)

				if currState["Updater 0"].(testUpdater).name != state["Updater 0"].(testUpdater).name {
					t.Error(
						"After", j+1, "dispatches a subscriber, ",
						"at index", i, ", was given a store with the inccorect data",
					)
				}
			}
		}()
	}

	st := New(state)
	for j := 0; j < updatesPreSubscribers; j++ {
		subs.publish(st)
	}
	wg.Wait()
}
