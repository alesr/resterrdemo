package bar

import (
	"testing"

	domain "github.com/alesr/resterrdemo/service/bar"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPostgres(t *testing.T) {
	t.Parallel()

	got := NewPostgres()

	require.NotNil(t, got)
	assert.IsType(t, &Postgresql{}, got)
}

func TestPostgresql_Fetch(t *testing.T) {
	t.Parallel()

	pg := Postgresql{}
	got := pg.Fetch()

	assert.Equal(t, domain.ErrBarNotFound, got)
}
