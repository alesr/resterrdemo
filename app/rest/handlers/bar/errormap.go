package bar

import (
	"github.com/alesr/resterr"
)

// ErrMap is the mapping between business layer errors (services) and the JSON errors
// we want to send back from the REST API.
//
// All expected errors resulting from downstream processing should be mapped here.
// Errors that are not mapped are sent to the client as a 500 error without details.
var ErrMap = map[error]resterr.RESTErr{
	// In this case, we're choosing to not map any errors from the bar service which will
	// result to errors being translated as 500.
}
