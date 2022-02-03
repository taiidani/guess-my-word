package main

import (
	"context"
	"fmt"
	"guess_my_word/actions"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
)

const defaultAddress = ":3000"

func main() {
	r := gin.Default()

	if err := actions.AddHandlers(r); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("Listening and serving HTTP on %s\n", defaultAddress)
	srv := http.Server{Addr: defaultAddress, Handler: r}

	done := make(chan interface{})
	go gracefulShutdown(&srv, done)
	srv.ListenAndServe()
	<-done
}

func gracefulShutdown(srv *http.Server, done chan<- interface{}) {
	const drainTimeout = time.Minute

	// Wait for the process to be interrupted
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	<-ctx.Done()

	// Gracefully drain all connections
	drainCtx, cancel := context.WithTimeout(context.Background(), drainTimeout)
	defer cancel()
	if err := srv.Shutdown(drainCtx); err != nil {
		fmt.Fprintf(os.Stderr, "failed to gracefully shut down: %s\n", err)
	}
	done <- true
}
