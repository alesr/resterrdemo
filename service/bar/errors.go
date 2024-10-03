package bar

import "errors"

var (
	// Enumerate repository errors.
	// These errors are used by the repository layer to represent failed
	// operations when interacting with the database.
	// We define them here, not in the repository package, to avoid creating
	// a dependency between the domain layer and the storage layer.
	// By defining the errors in the domain layer, we ensure that changes to the repository implementation
	// won't require changes to the business logic. Instead, the new implementation will need to adapt
	// to the domain layer, not the other way around.
	ErrBarNotFound = errors.New("bar not found")

	// Enumerate service errors.
	// These errors represent issues that can occur
	// during the processing of business logic.
	ErrBarUnavailable = errors.New("bar is unavailable at the moment")

	// All errors in this list may or may not be included in the error map used by the error handler in the transport layer.
	// If an error is not mapped to a JSON representation, it implies that the error should not be exposed to the client.
)
