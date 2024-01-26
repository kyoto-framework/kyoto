package component

// Future is a component state getter.
// Under the hood it waits for async.Future,
// gets resulting state and completes it with metadata.
type Future func() State
