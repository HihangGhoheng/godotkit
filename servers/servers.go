package gdk_server

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type HttpServerImpl interface {
	ListenAndServe() error
	SetKeepAlivesEnabled(v bool)
	Shutdown(ctx context.Context) error
}

type HttpServer struct {
	HttpServerImpl
	Log  *logrus.Logger
	Addr string
}

func (srv *HttpServer) gracefullShutdown(ctx context.Context, quit chan os.Signal) {
	signal.Notify(
		quit,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	<-quit
	srv.Log.Warnf("Got signal: %v, shutting down server", quit)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	srv.SetKeepAlivesEnabled(false)
	if err := srv.Shutdown(ctx); err != nil {
		srv.Log.Errorf("Got error when shutting down server: %+v", err)
	}

	srv.Log.Print("Server exiting!")
}

func MakeServer(port string, log *logrus.Logger, apihandler http.Handler) *HttpServer {
	srv := http.Server{
		Addr:    ":" + port,
		Handler: apihandler,
	}

	return &HttpServer{
		&srv,
		log,
		":" + port,
	}
}

func Run(ctx context.Context, srv *HttpServer, quit chan os.Signal) {
	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			srv.Log.Errorf("Failed to up service. Got error %+v", err)
		}
	}()

	srv.Log.Infof("Service is up! Listen on %s", srv.Addr)
	srv.gracefullShutdown(ctx, quit)
}
