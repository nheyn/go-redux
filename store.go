package store

// A Store keeps track of data in a State, and "attempts to make state mutations predictable".
type Store struct {
	actionQueue       chan queuedAction
	accessState       chan func(*State)
	accessSubscribers chan func(*subscriberSet)
}

// Creates a new Store that start with the given state.
func New(initialState State) Store {
	s := Store{
		actionQueue:       make(chan queuedAction, 8), //NOTE, 8 was randomly chosen - do test to see what works
		accessState:       make(chan func(*State)),
		accessSubscribers: make(chan func(*subscriberSet)),
	}

	go s.trackState(initialState)
	go s.listenForActions()
	go s.trackSubscribers()

	return s
}

// Dispatches the given action to all of the Updaters in the state of the Store. If an error is
// returned, then the State will not not change (even for the Updaters that had already completed).
func (s Store) Dispatch(action interface{}) error {
	errChan := make(chan error)
	s.actionQueue <- queuedAction{action, errChan}

	return <-errChan
}

// Select allows the given selector to pull its required data from the current State of the Store.
func (s Store) Select(sel Selector) {
	done := make(chan struct{})
	s.accessState <- func(st *State) {
		defer close(done)

		sel.SelectFrom(st)
	}
	<-done
}

// Send a refrence to the Store to the given subscriber every time the State is updated.
func (s Store) Subscribe(sub subscriber) func() bool {
	s.accessSubscribers <- func(subs *subscriberSet) {
		subs.add(sub)
	}

	return func() bool {
		didUnsub := make(chan bool)
		s.accessSubscribers <- func(subs *subscriberSet) {
			didUnsub <- subs.remove(sub)
		}
		return <-didUnsub
	}
}

// A struct that contains an action wating to be dispatched to the Updaters. It also includes a channel
// send any errors that occur, and is closed when the action is complete.
type queuedAction struct {
	action interface{}
	err    chan error
}

// A method that will list for actions in the action queue, and start a new goroutine to peform them when
// the state is aviable.
func (s Store) listenForActions() {
	for curr := range s.actionQueue {
		currState := State{}
		s.Select(&currState)

		newState, err := performUpdates(currState, curr.action)
		if err != nil {
			curr.err <- err
			close(curr.err)

			continue
		}

		done := make(chan struct{})
		s.accessState <- func(mutableSt *State) {
			for key, data := range newState {
				(*mutableSt)[key] = data
			}

			close(curr.err)
			close(done)
		}
		<-done
	}
}

// A method that will keep track of the subscriberSet, which can only be accessed throught the
// accessSubscribers channel.
func (s Store) trackSubscribers() {
	subs := subscriberSet{}

	for accessFn := range s.accessSubscribers {
		accessFn(&subs)
	}
}

// A method that will keep track of the state, which can only be accessed throught the accessState
// channel. The data in the given State will be put into the tracked State to start.
func (s Store) trackState(initialState State) {
	currState := State{}
	for key, data := range initialState {
		currState[key] = data
	}

	for accessFn := range s.accessState {
		accessFn(&currState)
	}
}
