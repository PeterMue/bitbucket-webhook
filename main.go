package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/PeterMue/bitbucket-webhook/config"
	"github.com/PeterMue/bitbucket-webhook/header"
)

// Dispatch webhook events to configured Handler if any
func dispatchEvent(w http.ResponseWriter, r *http.Request) {
	h := header.New(r.Header)

	// Diagnostics works without signature check
	if h.EventKey == "diagnostics:ping" {
		if _, err := w.Write([]byte(`{ "status" : "OK" }`)); err != nil {
			log.Printf("Failed to dispatch %s: %v", h.EventKey, err)
		}
		return
	}

	// Request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{ "status" : "Error" }`))
		return
	}

	// Config

	cfg, ok := r.Context().Value("config").(*config.Config)
	if !ok {
		log.Panicf("Unable to get config from request context")
	}

	// Everything else needs a valid signature
	if valid, err := h.Signature.Validate(body, cfg.Secret); !valid {
		log.Printf("Signature verification falied: %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid signature"))
		return
	}

	// Go, find event handler
	handlers, err := cfg.Handler(h.EventKey)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No handler for given event key configured"))
		return
	}

	// load body and run handler
	failed := 0
	for _, handler := range handlers {
		if err := handler.Run(h, body); err != nil {
			failed++
		}
	}

	if failed > 0 {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Failed to complete webhooks (%d of %d failed)", failed, len(handlers))))
		return
	}

}

// Load configuration and validate or fail with os.Exit(1)
func loadConfig() *config.Config {
	cfg, err := config.ParseFlags(os.Args)
	if err != nil {
		log.Fatalf("Invalid config: %s", err)
	}
	return cfg
}

// Start http server
func startServer(cfg *config.Config) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/webhook", dispatchEvent)

	srv := &http.Server{
		Handler: mux,
		Addr:    cfg.Listen,
		BaseContext: func(l net.Listener) context.Context {
			return context.WithValue(context.Background(), "config", cfg)
		},
	}

	go func() {
		log.Println("starting server")
		err := srv.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			log.Println("server shutdown")
		} else {
			fmt.Printf("server err: %v\n", err)
			os.Exit(1)
		}
	}()

	return srv
}

// Stop http server
func stopServer(srv *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	defer cancel()
	srv.Shutdown(ctx)
}

// Simple webhook listener for Bitbucket webhooks that executes configurable shell commands when a webhook is triggered.
func main() {
	cfg := loadConfig()

	// start server
	srv := startServer(cfg)

	// subscribe for sytem signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM)
	defer close(sigs)

	// handle system signals
	for {
		select {
		case sig := <-sigs:
			switch sig {
			case syscall.SIGHUP: // Reload
				stopServer(srv)
				cfg = loadConfig()
				srv = startServer(cfg)
			case syscall.SIGINT: // Graceful Shutdown
				stopServer(srv)
				os.Exit(0)
			}
		}
	}
}
