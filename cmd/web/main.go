package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

// define strut to hold application-wide dependencies
type application struct {
	logger *slog.Logger
}

func main() {
	// define new command-line flag
	// saved as a pointer
	addr := flag.String("addr", ":4000", "HTTP network address")

	// parse command-line flag
	// NOTE: must be called before using the variable
	flag.Parse()

	// initialize new structured logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	/* 	// DEBUG LOGGING
		logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		AddSource: true,
	})) */

	app := &application{logger: logger,}

	// dereference pointer before using flag
	/* logger.Info("starting server", "addr", *addr) */
	logger.Info("starting server", slog.String("addr", *addr))

	// start new web server
	// params: TCP network address, servermux
	err := http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}
