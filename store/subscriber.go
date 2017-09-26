package store

// A subscriber is a channel that will send the Store on updates.
type subscriber chan<- *Store

// A map that repecents a set of subscribers.
type subscriberSet map[subscriber]struct{}

// Adds the given subscriber to the set.
func (subs *subscriberSet) add(sub subscriber) {
	(*subs)[sub] = struct{}{}
}

// Removes the given subscriber from the set. Returns false if the given subscription is
// not in the set to remove.
func (subs *subscriberSet) remove(sub subscriber) bool {
	_, hasSub := (*subs)[sub]
	if !hasSub {
		return false
	}

	close(sub)
	delete(*subs, sub)
	return true
}

// Sends the given Store to all of the set's subscribers.
func (subs *subscriberSet) publish(st *Store) {
	for sub, _ := range *subs {
		sub <- st
	}
}
