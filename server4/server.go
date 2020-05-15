package server4

import (
	// "compress/flate"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

var (
	debug = false

	HelloWorldMessage = map[string]interface{}{
		"message": "Hello World",
	}
)

type Handler struct{}

func (hdlr *Handler) helloWorld(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	je := json.NewEncoder(w)
	if debug {
		je.SetIndent("", "  ")
	}
	if err := je.Encode(HelloWorldMessage); err != nil {
		log.Printf("warn: failed to write json response: %v", err)
	}
}

func ListenAndServe(ctx context.Context, stdout io.Writer, addr string) error {
	hdlr := &Handler{}

	r := chi.NewRouter()
	// r.Use(middleware.NewCompressor(flate.DefaultCompression, "application/json").Handler())
	r.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{
		Logger: log.New(stdout, "", log.LstdFlags),
	}))
	r.Get("/", hdlr.helloWorld)

	srv := &http.Server{
		Addr:        addr,
		Handler:     r,
		IdleTimeout: 10 * time.Second,
	}
	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.ListenAndServe()
	}()
	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("warn: failed to stop the server gracefully: %v", err)
		}
		return ctx.Err()
	}
}
