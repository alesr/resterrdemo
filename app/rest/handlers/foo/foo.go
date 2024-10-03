package foo

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
)

type fooService interface {
	Fetch() error
}

type errHandler interface {
	Handle(ctx context.Context, w http.ResponseWriter, err error)
}

// FooHandler implements HTTP handlers and processes requests related to the foo resource.
type FooHandler struct {
	logger     *slog.Logger
	fooSvc     fooService
	errHandler errHandler
}

// NewHandler instantiates a new FooHandler struct.
func NewHandler(logger *slog.Logger, fooSvc fooService, errHandler errHandler) (*FooHandler, error) {
	return &FooHandler{
		logger:     logger.WithGroup("foo-rest-handler"),
		fooSvc:     fooSvc,
		errHandler: errHandler,
	}, nil
}

// Get mimics an HTTP handler for fetching a foo resource.
func (fh *FooHandler) Get(w http.ResponseWriter, r *http.Request) {
	if err := fh.fooSvc.Fetch(); err != nil {
		fh.errHandler.Handle(r.Context(), w, fmt.Errorf("could not get foo from service: %w", err))
		return
	}
	w.WriteHeader(http.StatusOK)
}
