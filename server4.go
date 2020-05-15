package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/comail/colog"
	"github.com/taisukeyamashita/test/server4"
)

func run4(stdout, stderr io.Writer, args []string) int {
	colog.Register()
	colog.SetOutput(stderr)
	colog.SetMinLevel(colog.LDebug)

	var (
		addr string
	)
	flgs := flag.NewFlagSet("server", flag.ExitOnError)
	flgs.StringVar(&addr, "l", "127.0.0.1:8080", "listen address")
	if err := flgs.Parse(args[1:]); err != nil {
		fmt.Fprintln(stderr, err)
		return 128
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	errCh := make(chan error, 1)
	go func() {
		log.Printf("info: starting server on %s ...", addr)
		errCh <- server4.ListenAndServe(ctx, stdout, addr)
	}()

	select {
	case err := <-errCh:
		if err != nil {
			log.Printf("info: server encounters an error: %v", err)
		} else {
			log.Print("info: server has stopped unexpectedly")
		}
		return 1
	case sig := <-sigCh:
		cancel()
		if err := <-errCh; errors.Is(err, context.Canceled) {
			return 0
		} else if err != nil {
			log.Printf("info: failed to stop server by receiving a signal %s: %v", sig, err)
			return 1
		} else {
			return 0
		}
	case <-ctx.Done():
		return 0
	}
}

func main4() {
	os.Exit(run4(os.Stdout, os.Stderr, os.Args))
}
