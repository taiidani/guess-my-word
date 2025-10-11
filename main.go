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
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/getsentry/sentry-go"
	sentryslog "github.com/getsentry/sentry-go/slog"
	"github.com/go-chi/chi/v5"
	gsessions "github.com/gorilla/sessions"
	"github.com/quasoft/memstore"
)

func main() {
	port := flag.Int("port", 3000, "port to listen on")
	help := flag.Bool("help", false, "displays help text and exits")
	flag.Parse()
	bind := fmt.Sprintf(":%d", *port)

	if help != nil && *help {
		fmt.Fprintf(os.Stderr, "This application serves the Guess My Word app at %s\n", bind)
		os.Exit(0)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Set up Sentry
	err := sentry.Init(sentry.ClientOptions{
		SampleRate:       1.0,
		EnableTracing:    true,
		TracesSampleRate: 1.0,
		EnableLogs:       true,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	defer sentry.Flush(2 * time.Second)

	// Set up the structured logger
	initLogging(ctx)

	// Set up the app
	r := chi.NewRouter()

	if err := setupStores(ctx, r); err != nil {
		log.Fatalf("Unable to set up datastore: %s", err)
	}

	// Add all HTTP handlers
	if err := app.AddHandlers(r); err != nil {
		sentry.CaptureException(err)
		os.Exit(1)
	}

	slog.Info("Listening and serving HTTP", "port", bind)
	srv := http.Server{Addr: bind, Handler: r}

	done := make(chan any)
	go gracefulShutdown(ctx, &srv, done)
	_ = srv.ListenAndServe()
	<-done
}

func initLogging(ctx context.Context) {
	var logger *slog.Logger

	switch os.Getenv("SENTRY_ENVIRONMENT") {
	case "prod", "production":
		handler := sentryslog.Option{
			// Explicitly specify the levels that you want to be captured.
			EventLevel: []slog.Level{slog.LevelError},                                 // Captures only [slog.LevelError] as error events.
			LogLevel:   []slog.Level{slog.LevelWarn, slog.LevelInfo, slog.LevelDebug}, // Captures remaining items as log entries.
		}.NewSentryHandler(ctx)
		logger = slog.New(handler)
	default:
		var level slog.Level
		switch os.Getenv("LOG_LEVEL") {
		case "error":
			level = slog.LevelError
		case "warn":
			level = slog.LevelWarn
		case "debug":
			level = slog.LevelDebug
		default:
			level = slog.LevelInfo
		}

		handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: level,
		})
		logger = slog.New(handler)
	}

	slog.SetDefault(logger)
}

func setupStores(ctx context.Context, r chi.Router) error {
	var dataClient datastore.Client
	var sessionClient gsessions.Store
	var err error

	db, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	if addr, ok := os.LookupEnv("REDIS_ADDR"); ok {
		sessionClient, err = sessions.NewRedis(addr, db, []byte("secret"))
		if err != nil {
			return fmt.Errorf("redis setup failure: %w", err)
		}

		dataClient = datastore.NewRedis(addr, db)
	} else if host, ok := os.LookupEnv("REDIS_HOST"); ok {
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
		slog.Warn("No REDIS_ADDR or REDIS_HOST env var set. Falling back upon in-memory store")
		sessionClient = memstore.NewMemStore([]byte("secret"))
		dataClient = datastore.NewMemory()
	}

	// Set up data storage
	app.SetupStores(
		words.NewListStore(dataClient),
		words.NewWordStore(dataClient),
	)

	// Set up session management
	sessions.Configure(r, "guessmyword", sessionClient)

	return words.PopulateDefaultLists(ctx, dataClient)
}

func gracefulShutdown(ctx context.Context, srv *http.Server, done chan<- any) {
	const drainTimeout = time.Minute

	// Wait for the process to be interrupted
	<-ctx.Done()
	fmt.Fprintf(os.Stderr, "Interrupted. Shutting down\n")

	// Gracefully drain all connections
	drainCtx, cancel := context.WithTimeout(context.Background(), drainTimeout)
	defer cancel()
	if err := srv.Shutdown(drainCtx); err != nil {
		fmt.Fprintf(os.Stderr, "failed to gracefully shut down: %s\n", err)
	}
	done <- true
}
