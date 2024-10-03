package foo

import (
	"errors"
)

var errNetworkKaput = errors.New("network kaput")

// Postgresql carries the connection to the database,
// and the methods to interact with it.
type Postgresql struct{}

// NewPostgres instantiates a new Postgresql struct.
func NewPostgres() *Postgresql { return &Postgresql{} }

// Fetch fetches foo entities from the database.
// In our example, we simulate an network error which we would not
// like to expose on HTTP responses.
func (p Postgresql) Fetch() error { return errNetworkKaput }
