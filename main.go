package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/alesr/resterr"
	"github.com/alesr/resterrdemo/app/rest"
	barhandler "github.com/alesr/resterrdemo/app/rest/handlers/bar"
	foohandler "github.com/alesr/resterrdemo/app/rest/handlers/foo"
	barrepo "github.com/alesr/resterrdemo/repository/bar"
	foorepo "github.com/alesr/resterrdemo/repository/foo"
	"github.com/alesr/resterrdemo/service/bar"
	"github.com/alesr/resterrdemo/service/foo"
)

const addr = ":8080"

func main() {
	logger := slog.Default()

	// Initialize foo storage, service (business) and transport error handler.

	fooRepo := foorepo.NewPostgres()
	fooSvc := foo.New(fooRepo)

	fooErrHandler, err := resterr.NewHandler(logger, foohandler.ErrMap)
	if err != nil {
		logger.Error("Failed to initialize foo error handler.", errAttr(err))
		os.Exit(1)
	}

	fooHandler, err := foohandler.NewHandler(logger, fooSvc, fooErrHandler)
	if err != nil {
		logger.Error("Failed to initialize foo handler.", errAttr(err))
		os.Exit(2)
	}

	// Initialize bar storage, service (business) and transport error handler.

	barRepo := barrepo.NewPostgres()
	barSvc := bar.New(barRepo)

	barErrHandler, err := resterr.NewHandler(logger, barhandler.ErrMap)
	if err != nil {
		logger.Error("Failed to initialize bar error handler.", errAttr(err))
		os.Exit(3)
	}

	barHandler, err := barhandler.NewHandler(logger, barSvc, barErrHandler)
	if err != nil {
		logger.Error("Failed to initialize bar handler.", errAttr(err))
		os.Exit(4)
	}

	// Inject handles on our REST transport layer.

	restApp, err := rest.NewApp(logger, addr, fooHandler, barHandler)
	if err != nil {
		logger.Error("Failed to initialize REST APP.", errAttr(err))
		os.Exit(5)
	}

	go func() {
		if err := restApp.Run(); err != nil {
			logger.Error("Failed to run REST APP.", slog.String("addr", addr), errAttr(err))
			os.Exit(6)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	<-ctx.Done()

	if err := restApp.Shutdown(ctx); err != nil {
		logger.Error("Failed to shutdown REST APP.", errAttr(err))
		os.Exit(7)
	}
}

func errAttr(err error) slog.Attr {
	return slog.String("error", err.Error())
}

// There's a benefit in explicitely declaring the application's dependencies directly on
// our main function.
// However, this can lead us to having extensive and hard to read main files. As you see fit,
// consider spliting the main logic into multiple files under the same main package for
// better readability and reusability accross applications.
// Note that, if choose doing so, you need to build the binary running `go build .` instead of `go build main.go`.
