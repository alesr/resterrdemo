package bar

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
)

type barService interface {
	Fetch() error
}

type errHandler interface {
	Handle(ctx context.Context, w http.ResponseWriter, err error)
}

// BarHandler implements HTTP handlers and processes requests related to the bar resource.
type BarHandler struct {
	logger     *slog.Logger
	barSvc     barService
	errHandler errHandler
}

// NewHandler instantiates a new BarHandler struct.
func NewHandler(logger *slog.Logger, barSvc barService, errHandler errHandler) (*BarHandler, error) {
	return &BarHandler{
		logger:     logger.WithGroup("bar-rest-handler"),
		barSvc:     barSvc,
		errHandler: errHandler,
	}, nil
}

// Get mimics an HTTP handler for fetching a bar resource.
func (bh *BarHandler) Get(w http.ResponseWriter, r *http.Request) {
	if err := bh.barSvc.Fetch(); err != nil {
		bh.errHandler.Handle(r.Context(), w, fmt.Errorf("could not get bar from service: %w", err))
		return
	}
	w.WriteHeader(http.StatusTeapot)
}
