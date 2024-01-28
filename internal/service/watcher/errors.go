package watcher

import "fmt"

// WatcherFailedInitError represents an error when trying to init a Watcher (fsnotify.Watcher)
type WatcherFailedInitError struct {
	Err error
}

// Allow WatcherFailedInitError to satisfy error interface.
func (e *WatcherFailedInitError) Error() string {
	return fmt.Sprintf("Watcher failed to initialize: %q", e.Err)
}

// EventError represents an error when trying to catch events in the Watcher (fsnotify.Watcher)
type EventError struct {
	Err error
}

// Allow EventError to satisfy error interface.
func (e *EventError) Error() string {
	return fmt.Sprintf("Event error: %q", e.Err)
}
