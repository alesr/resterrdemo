package foo

import (
	"net/http"

	"github.com/alesr/resterr"
	"github.com/alesr/resterrdemo/service/foo"
)

// ErrMap is the mapping between business layer errors (services) and the JSON errors
// we want to send back from the REST API.
//
// All expected errors resulting from downstream processing should be mapped here.
// Errors that are not mapped are sent to the client as a 500 error without details.
var ErrMap = map[error]resterr.RESTErr{
	foo.ErrGetFaleid: {
		StatusCode: http.StatusTeapot,
		Message:    "could not perform the get foo operation",
	},
}
