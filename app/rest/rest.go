// rest package is the application's REST API.
// It implements an HTTP server and handlers to process requests
// by calling the appropriate service for processing.
package rest

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
)

type handler interface {
	Get(w http.ResponseWriter, r *http.Request)
}

// App implements the transport layer by running an HTTP server.
type App struct {
	logger     *slog.Logger
	server     *http.Server
	fooHandler handler
	barHandler handler
}

// NewApp instantiates a new App struct.
func NewApp(logger *slog.Logger, addr string, fooHdler, barHdler handler) (*App, error) {
	app := App{
		logger:     logger.WithGroup("rest-app"),
		fooHandler: fooHdler,
		barHandler: barHdler,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /foo", app.fooHandler.Get)
	mux.HandleFunc("GET /bar", app.barHandler.Get)

	app.server = &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	return &app, nil
}

// Run starts the application, serving on the specified address and port as provided in the configuration.
func (app *App) Run() error {
	app.logger.Info("Starting REST demo app.", slog.String("addr", app.server.Addr))

	if err := app.server.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			app.logger.Warn("Server closed", slog.String("addr", app.server.Addr))
			return nil
		}
		return fmt.Errorf("could not listen and serve: %w", err)
	}
	return nil
}

// Shutdown gracefully shuts down the server.
// The shutdown timeout can be controlled by the context passed as an argument.
func (app *App) Shutdown(ctx context.Context) error {
	app.logger.InfoContext(ctx, "Server shutting down")
	if err := app.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("could not shut down server: %w", err)
	}
	return nil
}
