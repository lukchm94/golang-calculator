package eventsDispatcher

import "errors"

type EventsDispatcherError error

var (
	NotImplementedError EventsDispatcherError = errors.New("[EventsDispatcherError] not implemented")
)
