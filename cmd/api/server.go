package main

import (
	"net/http"
	"time"
)

func (app *application) server(mux http.Handler) error {
	server := &http.Server{
		Addr:         ":" + app.config.port,
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.logger.Infow("starting server on port", "port", app.config.port, "environment", app.config.env)
	return server.ListenAndServe()

}
