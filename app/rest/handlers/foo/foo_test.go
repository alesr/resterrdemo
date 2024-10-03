package foo

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alesr/resterr"
	"github.com/alesr/resterrdemo/service/foo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var noopLogger = slog.New(slog.NewJSONHandler(io.Discard, nil))

type serviceMock struct {
	fetchFunc func() error
}

func (m *serviceMock) Fetch() error {
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

	svc := &serviceMock{}

	errHandler := &errHandlerMock{}

	handler, err := NewHandler(noopLogger, svc, errHandler)

	require.NoError(t, err)
	require.NotNil(t, handler)
	assert.NotNil(t, handler.logger)
	assert.NotNil(t, handler.fooSvc)
	assert.NotNil(t, handler.errHandler)
}

func TestFooHandler_Get(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name            string
		serviceError    error
		expectedStatus  int
		wasErrorHandled bool
	}{
		{
			name:           "service returns no error",
			expectedStatus: http.StatusOK,
		},
		{
			name:            "service returns error",
			serviceError:    assert.AnError,
			expectedStatus:  http.StatusInternalServerError,
			wasErrorHandled: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var errorHandled bool
			svc := serviceMock{
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

			handler, err := NewHandler(noopLogger, &svc, &errHandler)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodGet, "/foo", nil)
			w := httptest.NewRecorder()

			handler.Get(w, req)

			assert.Equal(t, tc.expectedStatus, w.Result().StatusCode)
			assert.Equal(t, tc.wasErrorHandled, errorHandled)
		})
	}
}

func TestFooHandler_ErrMap(t *testing.T) {
	t.Parallel()

	errHandler, err := resterr.NewHandler(noopLogger, ErrMap)
	require.NoError(t, err)

	testCases := []struct {
		name  string
		given error
		want  resterr.RESTErr
	}{
		{
			name:  "unmapped error is returned as internal server error",
			given: assert.AnError,
			want: resterr.RESTErr{
				StatusCode: http.StatusInternalServerError,
				Message:    "something went wrong",
			},
		},
		{
			name:  "mapped error is returned as the equivalent JSON error",
			given: foo.ErrGetFaleid,
			want: resterr.RESTErr{
				StatusCode: http.StatusTeapot,
				Message:    "could not perform the get foo operation",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			errHandler.Handle(context.TODO(), w, tc.given)

			defer w.Result().Body.Close()

			var result resterr.RESTErr

			err := json.NewDecoder(w.Result().Body).Decode(&result)
			require.NoError(t, err)

			assert.Equal(t, tc.want.StatusCode, result.StatusCode)
			assert.Equal(t, tc.want.Message, result.Message)
		})
	}
}
