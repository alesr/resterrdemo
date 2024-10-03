package foo

import (
	"errors"
)

var (
	// Enumerate service errors.
	// These errors represent issues that can occur
	// during the processing of business logic.
	ErrGetFaleid = errors.New("could not get foo")

	// All errors in this list may or may not be included in the error map used by the error handler in the transport layer.
	// If an error is not mapped to a JSON representation, it implies that the error should not be exposed to the client.
)
