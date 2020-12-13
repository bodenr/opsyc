package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/justinas/alice"
	"github.com/rs/zerolog/hlog"

	"github.com/bodenr/opsyc/api"
	"github.com/bodenr/opsyc/util"
)

type App struct {
	UIHandler *api.UIHandler
}

func (app *App) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var resource string

	resource, req.URL.Path = util.ShiftPath(req.URL.Path)
	if resource == "ui" {
		app.UIHandler.ServeHTTP(res, req)
		return
	}
	http.Error(res, "Not Found", http.StatusNotFound)
}

func newServer() *http.Server {

	middleware := alice.New()
	middleware = middleware.Append(hlog.NewHandler(util.Log))

	middleware = middleware.Append(hlog.AccessHandler(
		func(r *http.Request, status, size int, duration time.Duration) {
			hlog.FromRequest(r).Info().
				Str("method", r.Method).
				Str("url", r.URL.String()).
				Int("status", status).
				Int("size", size).
				Dur("duration", duration).
				Msg("")
		}))

	middleware = middleware.Append(hlog.RemoteAddrHandler("ip"))
	middleware = middleware.Append(hlog.UserAgentHandler("user_agent"))
	middleware = middleware.Append(hlog.RefererHandler("referer"))
	middleware = middleware.Append(hlog.RequestIDHandler("req_id", "Request-Id"))

	app := middleware.Then(&App{UIHandler: new(api.UIHandler)})

	return &http.Server{
		Addr:         ":8080",
		Handler:      app,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
}

func httpShutdown(server *http.Server, terminate <-chan os.Signal, done chan bool) {
	termSig := <-terminate
	util.Log.Info().Str("lifecycle", "shutdown").Str("component", "main").Str(
		"signal", termSig.String()).Msg("")

	util.Log.Info().Str("component", "http").Msg("shutting down http server")

	ctx, halt := context.WithTimeout(context.Background(), 20*time.Second)
	defer halt()

	server.SetKeepAlivesEnabled(false)
	if err := server.Shutdown(ctx); err != nil {
		util.Log.Fatal().Str("lifecycle", "stop").Str("component", "http").Msg(
			"failed to gracefully stop http server")
	} else {
		util.Log.Info().Str("lifecycle", "stop").Str("component", "http").Msg(
			"graceful shutdown of http complete")
	}

	close(done)
}

func main() {
	util.Log.Info().Str("lifecycle", "start").Str("component", "main").Msg("")

	serverStop := make(chan bool, 1)
	sigStop := make(chan os.Signal)
	signal.Notify(sigStop, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT)

	server := newServer()

	go httpShutdown(server, sigStop, serverStop)

	util.Log.Info().Str("lifecycle", "start").Str("component", "http").Msg(
		"starting http server on port " + server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		util.Log.Fatal().Str("component", "http").Msg(
			"failed to start http server due to " + err.Error())
	}

	<-serverStop
	util.Log.Info().Str("lifecycle", "stop").Str("component", "main").Msg("")
}
