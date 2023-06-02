package main

import (
	"context"
	"flag"
	"fmt"
	"guess_my_word/app"
	"guess_my_word/internal/datastore"
	"guess_my_word/internal/sessions"
	"guess_my_word/internal/words"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
)

const defaultAddress = ":3000"

func main() {
	help := flag.Bool("help", false, "displays help text and exits")
	flag.Parse()

	if help != nil && *help {
		fmt.Fprintln(os.Stderr, "This application serves the Guess My Word app at "+defaultAddress)
		os.Exit(0)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	r := gin.Default()

	if err := setupStores(ctx, r); err != nil {
		log.Fatalf("Unable to set up datastore: %s", err)
	}

	// Load the HTML templates into gin
	if err := app.SetupTemplates(r); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load embedded templates: %s\n", err)
		os.Exit(1)
	}

	// And the static assets
	if err := app.SetupAssets(r); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load embedded assets: %s\n", err)
		os.Exit(1)
	}

	// Add all HTTP handlers
	if err := app.AddHandlers(r); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("Listening and serving HTTP on %s\n", defaultAddress)
	srv := http.Server{Addr: defaultAddress, Handler: r}

	done := make(chan interface{})
	go gracefulShutdown(ctx, &srv, done)
	_ = srv.ListenAndServe()
	<-done
}

func setupStores(ctx context.Context, r *gin.Engine) error {
	var dataClient datastore.Client
	var sessionClient sessions.Store
	var err error

	if addr, ok := os.LookupEnv("REDIS_ADDR"); ok {
		sessionClient, err = sessions.NewRedis(addr, []byte("secret"))
		if err != nil {
			return fmt.Errorf("redis setup failure: %w", err)
		}

		dataClient = datastore.NewRedis(addr)

	} else if host, ok := os.LookupEnv("REDIS_HOST"); ok {
		db := 0
		if dbParsed, err := strconv.ParseInt(os.Getenv("REDIS_DB"), 10, 64); err == nil {
			db = int(dbParsed)
		}

		port := os.Getenv("REDIS_PORT")
		user := os.Getenv("REDIS_USER")
		pass := os.Getenv("REDIS_PASSWORD")

		sessionClient, err = sessions.NewRedisSecure(host, port, user, pass, db, []byte("secret"))
		if err != nil {
			return fmt.Errorf("redis setup failure: %w", err)
		}

		dataClient = datastore.NewRedisSecure(
			host,
			port,
			user,
			pass,
			db,
		)
	} else {
		log.Println("WARNING: No REDIS_ADDR or REDIS_HOST env var set. Falling back upon in-memory store")
		sessionClient = memstore.NewStore([]byte("secret"))
		dataClient = datastore.NewMemory()
	}

	// Set up data storage
	app.SetupStores(
		words.NewListStore(dataClient),
		words.NewWordStore(dataClient),
	)

	// Set up session management
	sessions.Configure(r, sessionClient)

	return words.PopulateDefaultLists(ctx, dataClient)
}

func gracefulShutdown(ctx context.Context, srv *http.Server, done chan<- interface{}) {
	const drainTimeout = time.Minute

	// Wait for the process to be interrupted
	<-ctx.Done()

	// Gracefully drain all connections
	drainCtx, cancel := context.WithTimeout(context.Background(), drainTimeout)
	defer cancel()
	if err := srv.Shutdown(drainCtx); err != nil {
		fmt.Fprintf(os.Stderr, "failed to gracefully shut down: %s\n", err)
	}
	done <- true
}
