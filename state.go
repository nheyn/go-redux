package store

// A State is a map that contains the current data for a Store.
type State map[interface {}]Updater
