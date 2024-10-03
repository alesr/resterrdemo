package bar

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type barServiceMock struct {
	fetchFunc func() error
}

func (m *barServiceMock) Fetch() error {
	return m.fetchFunc()
}

type errHandlerMock struct {
	handleFunc func(ctx context.Context, w http.ResponseWriter, err error)
}

func (e *errHandlerMock) Handle(ctx context.Context, w http.ResponseWriter, err error) {
	e.handleFunc(ctx, w, err)
}

func TestNewHandler(t *testing.T) {
	t.Parallel()

	svc := &barServiceMock{}
	errHandler := &errHandlerMock{}

	handler, err := NewHandler(slog.Default(), svc, errHandler)

	require.NoError(t, err)
	require.NotNil(t, handler)
	assert.NotNil(t, handler.logger)
	assert.NotNil(t, handler.barSvc)
	assert.NotNil(t, handler.errHandler)
}

func TestBarHandler_Get(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name            string
		serviceError    error
		expectedStatus  int
		wasErrorHandled bool
	}{
		{
			name:            "service returns no error",
			serviceError:    nil,
			expectedStatus:  http.StatusTeapot,
			wasErrorHandled: false,
		},
		{
			name:            "service returns error",
			serviceError:    assert.AnError,
			expectedStatus:  http.StatusInternalServerError,
			wasErrorHandled: true,
		},
	}

	for _, tc := range testCases {
		tc := tc // capture the range variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var errorHandled bool

			svc := barServiceMock{
				fetchFunc: func() error {
					return tc.serviceError
				},
			}

			errHandler := errHandlerMock{
				handleFunc: func(ctx context.Context, w http.ResponseWriter, err error) {
					errorHandled = true
					assert.ErrorIs(t, err, tc.serviceError)
					w.WriteHeader(http.StatusInternalServerError)
				},
			}

			handler, err := NewHandler(slog.Default(), &svc, &errHandler)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodGet, "/bar", nil)
			w := httptest.NewRecorder()

			handler.Get(w, req)

			assert.Equal(t, tc.expectedStatus, w.Result().StatusCode)
			assert.Equal(t, tc.wasErrorHandled, errorHandled)
		})
	}
}
