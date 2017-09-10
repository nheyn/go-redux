package store

import "testing"

func TestSubscriberSetCanSubscribe(t *testing.T) {
	testSubs := []subscriber{
		make(chan Store),
		make(chan Store),
		make(chan Store),
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
		make(chan Store),
		make(chan Store),
		make(chan Store),
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
