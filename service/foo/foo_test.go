package foo

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type repoMock struct {
	fetchFunc func() error
}

func (m *repoMock) Fetch() error {
	return m.fetchFunc()
}

func TestNew(t *testing.T) {
	t.Parallel()

	got := New(&repoMock{})

	require.NotNil(t, got)

	assert.IsType(t, &Service{}, got)
	assert.NotNil(t, got.repo)
}

func TestService_Fetch(t *testing.T) {
	t.Parallel()

	var fetchWasCalled bool

	repo := repoMock{
		fetchFunc: func() error {
			fetchWasCalled = true
			return assert.AnError
		},
	}

	svc := Service{repo: &repo}

	got := svc.Fetch()

	require.True(t, fetchWasCalled)
	assert.ErrorIs(t, got, ErrGetFaleid)
}
