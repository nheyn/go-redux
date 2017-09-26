package store

import "context"

// A PerformDispatch function is used to dispatch the given action to given State.
type PerformDispatch func(context.Context, State, interface{}) (State, error)

// A Store keeps track of data in a State, and "attempts to make state mutations predictable".
type Store struct {
	PerformDispatch
	actionQueue       chan queuedAction
	accessState       chan func(*State)
	accessSubscribers chan func(*subscriberSet)
}

// Creates a new Store that start with the given state.
func New(initialState State, configs ...func(*Store)) *Store {
	s := &Store{
		PerformDispatch:   nil,
		actionQueue:       make(chan queuedAction),
		accessState:       make(chan func(*State)),
		accessSubscribers: make(chan func(*subscriberSet)),
	}

	// Configure store
	defaultConfigs := []func(*Store){defaultPeformDispatchConfig}
	for _, config := range append(defaultConfigs, configs...) {
		config(s)
	}

	// Start store
	go s.trackState(initialState)
	go s.listenForActions()
	go s.trackSubscribers()

	return s
}

// Dispatches the given action to all of the Updaters in the state of the Store. If an error is
// returned, then the State will not not change (even for the Updaters that had already completed).
func (s *Store) Dispatch(ctx context.Context, action interface{}) error {
	errChan := make(chan error)
	s.actionQueue <- queuedAction{ctx, action, errChan}

	return <-errChan
}

// Select allows the given selector to pull its required data from the current State of the Store.
func (s *Store) Select(sel Selector) {
	done := make(chan struct{})
	s.accessState <- func(st *State) {
		defer close(done)

		sel.SelectFrom(st)
	}
	<-done
}

// Send a refrence to the Store to the given subscriber every time the State is updated.
func (s *Store) Subscribe(sub subscriber) func() bool {
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
	ctx    context.Context
	action interface{}
	err    chan error
}

// A method that will list for actions in the action queue, and start a new goroutine to peform them when
// the state is aviable.
func (s *Store) listenForActions() {
	for curr := range s.actionQueue {
		err := s.performAction(curr.ctx, curr.action)
		if err != nil {
			curr.err <- err
		}

		close(curr.err)
	}
}

// Perform the give action on the current State of the Store. It an error is returned,
// the State will not be updated.
func (s *Store) performAction(ctx context.Context, action interface{}) error {
	// Perform the action on the current state
	currState := State{}
	s.Select(&currState)

	newState, err := s.PerformDispatch(ctx, currState, action)
	if err != nil {
		return err
	}

	// Update the store with the updated state
	done := make(chan struct{})
	s.accessState <- func(mutableSt *State) {
		defer close(done)

		for key, data := range newState {
			(*mutableSt)[key] = data
		}
	}
	<-done

	// Tell subscribers about the change
	s.accessSubscribers <- func(subs *subscriberSet) {
		subs.publish(s)
	}

	return nil
}

// A method that will keep track of the state, which can only be accessed throught the accessState
// channel. The data in the given State will be put into the tracked State to start.
func (s *Store) trackState(initialState State) {
	currState := State{}
	for key, data := range initialState {
		currState[key] = data
	}

	for accessFn := range s.accessState {
		accessFn(&currState)
	}
}

// A method that will keep track of the subscriberSet, which can only be accessed throught the
// accessSubscribers channel.
func (s *Store) trackSubscribers() {
	subs := subscriberSet{}

	for accessFn := range s.accessSubscribers {
		accessFn(&subs)
	}
}
