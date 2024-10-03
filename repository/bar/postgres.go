package bar

import (
	domain "github.com/alesr/resterrdemo/service/bar"
)

// Postgresql carries the connection to the database,
// and the methods to interact with it.
type Postgresql struct{}

// NewPostgres instantiates a new Postgresql struct.
func NewPostgres() *Postgresql { return &Postgresql{} }

// Fetch fetches foo entities from the database.
// In our example, we simulate that we couldn't find a record.
func (p Postgresql) Fetch() error { return domain.ErrBarNotFound }
