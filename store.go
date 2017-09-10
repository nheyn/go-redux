package store

// A Store keeps track of data in a State, and "attempts to make state mutations predictable".
type Store struct {
	accessState chan func(*State)
	actionQueue chan queuedAction
}

// Creates a new Store that start with the given state.
func New(initialState State) Store {
	s := Store{
		accessState: make(chan func(*State)),
		actionQueue: make(chan queuedAction, 8), //NOTE, 8 was randomly chosen - do test to see what works
	}

	go s.trackState(initialState)
	go s.listenForActions()

	return s
}

// Dispatches the given action to all of the Updaters in the state of the Store. If an error is
// returned, then the State will not not change (even for the Updaters that had already completed).
func (s Store) Dispatch(action interface{}) error {
	errChan := make(chan error)
	s.actionQueue <- queuedAction{action, errChan}

	return <-errChan
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
		done := make(chan struct{})

		s.accessState <- func(st *State) {
			//NOTE, started in it's own goroutine and using the s.accessState channel twice so others
			//      can read the state while performUpdates(...) is running.
			go func(orginalState State) {
				newState, err := performUpdates(orginalState, curr.action)
				if err != nil {
					curr.err <- err
					close(curr.err)
					close(done)
					return
				}

				s.accessState <- func(mutableSt *State) {
					for key, data := range newState {
						(*mutableSt)[key] = data
					}

					close(curr.err)
					close(done)
				}
			}(*st)
		}

		<-done
	}
}

// A method that will keep track of the sate, which can only be accessed throught the accessState
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
