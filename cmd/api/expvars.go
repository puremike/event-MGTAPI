package main

import (
	"expvar"
	"net/http"
	"runtime"
	"time"
)

var (
	version      = "1.0.0"
	startTime    = time.Now()
	requestCount = expvar.NewInt("http_requests_total")
)

func init() {
	// Set static metrics
	expvar.NewString("version").Set(version)

	// Dynamic metrics
	expvar.Publish("memory", expvar.Func(func() any {
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)
		return memStats
	}))

	expvar.Publish("cpu", expvar.Func(func() any {
		return runtime.NumCPU()
	}))

	expvar.Publish("uptime", expvar.Func(func() any {
		return time.Since(startTime).String()
	}))

	expvar.Publish("gc_pause_ns", expvar.Func(func() any {
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)
		return memStats.PauseNs
	}))

	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))
}

func (app *application) expvars(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount.Add(1)
		next.ServeHTTP(w, r)
	})
}
