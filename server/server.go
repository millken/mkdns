package server

import (
	"context"
	"net/http"
	"sync"

	"github.com/millken/mkdns/drivers"
	"github.com/millken/mkdns/drivers/xdp"
	"github.com/millken/mkdns/internal/scheduler"
	"github.com/pkg/errors"
)

// Options is the set of optional parameters.
type Options struct {
	Driver drivers.Driver
}

type Server struct {
	handler   http.Handler
	scheduler *scheduler.Scheduler
	Server    http.Server
	driver    drivers.Driver
	once      sync.Once
}

func New(h http.Handler, opts *Options) *Server {
	srv := &Server{
		handler:   h,
		scheduler: scheduler.New(),
	}
	if opts != nil {
		srv.driver = opts.Driver
	}
	return srv
}

func (srv *Server) init() {
	srv.once.Do(func() {
		if srv.driver == nil {
			srv.driver = xdp.New(srv.scheduler)
		}
		if srv.handler == nil {
			srv.handler = http.DefaultServeMux
		}
		srv.scheduler.Start()
	})
}

func (srv *Server) ListenAndServe(addr string) error {
	srv.init()
	srv.Server.Addr = addr
	srv.Server.Handler = srv.handler
	if err := srv.driver.Start(context.Background()); err != nil {
		return errors.Wrapf(err, "failed to start driver")
	}
	return srv.Server.ListenAndServe()
}

// Shutdown gracefully shuts down the server without interrupting any active connections.
func (srv *Server) Shutdown(ctx context.Context) error {
	defer srv.scheduler.Stop()
	if err := srv.driver.Stop(ctx); err != nil {
		return err
	}
	return nil
}
