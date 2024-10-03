package rest

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type handlerMock struct {
	getFunc func(w http.ResponseWriter, r *http.Request)
}

func (h *handlerMock) Get(w http.ResponseWriter, r *http.Request) {
	h.getFunc(w, r)
}

func TestNewApp(t *testing.T) {
	logger := noopLogger()
	port := "dummy-port"
	fooHandler := &handlerMock{}
	barHandler := &handlerMock{}

	app, err := NewApp(logger, port, fooHandler, barHandler)
	require.NoError(t, err)
	require.NotNil(t, app)

	assert.NotNil(t, app.logger)
	assert.Equal(t, port, app.server.Addr)
	assert.Equal(t, fooHandler, app.fooHandler)
	assert.Equal(t, barHandler, app.barHandler)
}

func TestApp_Run_Shutdown(t *testing.T) {
	logger := noopLogger()

	fooHandler := handlerMock{
		getFunc: func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		},
	}

	barHandler := handlerMock{
		getFunc: func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusTeapot)
		},
	}

	app, err := NewApp(logger, ":8081", &fooHandler, &barHandler)
	require.NoError(t, err)
	require.NotNil(t, app)

	// Run the server in a separate goroutine
	go func() {
		err := app.Run()
		assert.NoError(t, err)
	}()

	// Wait for the server to start
	time.Sleep(100 * time.Millisecond)

	// Test /foo route
	req := httptest.NewRequest(http.MethodGet, "/foo", nil)
	w := httptest.NewRecorder()
	app.server.Handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

	// Test /bar route
	req = httptest.NewRequest(http.MethodGet, "/bar", nil)
	w = httptest.NewRecorder()
	app.server.Handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusTeapot, w.Result().StatusCode)

	// Shutdown the server
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err = app.Shutdown(ctx)
	assert.NoError(t, err)
}

func noopLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(io.Discard, nil))
}
