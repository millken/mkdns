package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/millken/golog"
	"github.com/millken/mkdns/server"
)

func helloHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hello\n")
}

func mainHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Welcome to the home page!\n")
}

type customHealthCheck struct {
	mu      sync.RWMutex
	healthy bool
}

func (h *customHealthCheck) CheckHealth() error {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if !h.healthy {
		return errors.New("not ready yet!")
	}
	return nil
}

func main() {
	addr := flag.String("listen", ":8080", "HTTP port to listen on")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", helloHandler)
	mux.HandleFunc("/", mainHandler)

	healthCheck := new(customHealthCheck)
	time.AfterFunc(10*time.Second, func() {
		healthCheck.mu.Lock()
		defer healthCheck.mu.Unlock()
		healthCheck.healthy = true
	})

	options := &server.Options{}

	s := server.New(mux, options)
	golog.Infof("Listening on %s", *addr)
	go func() {
		if err := s.ListenAndServe(*addr); err != nil {
			log.Fatal("failed to start server: ", err)
		}
	}()
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt, os.Kill)

	// wait INT or KILL
	<-stop
	golog.Info("shutting down ...")
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		golog.Error("shutting down err", err)
	}

}
