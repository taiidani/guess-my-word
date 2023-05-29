package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"guess_my_word/actions"
	"guess_my_word/internal/datastore"
	"guess_my_word/internal/words"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

const defaultAddress = ":3000"

//go:embed templates
var templates embed.FS

//go:embed assets
var assets embed.FS

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
	t, err := template.ParseFS(templates, "templates/**")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load embedded templates: %s\n", err)
		os.Exit(1)
	}
	r.SetHTMLTemplate(t)

	// And the static assets
	sub, err := fs.Sub(assets, "assets")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load embedded assets: %s\n", err)
		os.Exit(1)
	}
	r.StaticFS("/assets", http.FS(sub))

	// Add all HTTP handlers
	if err := actions.AddHandlers(r); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("Listening and serving HTTP on %s\n", defaultAddress)
	srv := http.Server{Addr: defaultAddress, Handler: r}

	done := make(chan interface{})
	go gracefulShutdown(ctx, &srv, done)
	srv.ListenAndServe()
	<-done
}

func setupStores(ctx context.Context, r *gin.Engine) error {
	var dataClient datastore.Client
	var sessionClient sessions.Store
	var err error

	if addr, ok := os.LookupEnv("REDIS_ADDR"); ok {
		sessionClient, err = redis.NewStore(10, "tcp", addr, "", []byte("secret"))
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

		sessionClient, err = redis.NewStore(
			10,
			"tcp",
			fmt.Sprintf("%s:%s", host, port),
			pass,
			[]byte("secret"),
		)
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
	actions.SetupStores(
		words.NewListStore(dataClient),
		words.NewWordStore(dataClient),
	)

	// Set up session management
	r.Use(sessions.Sessions("guessmyword", sessionClient))

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
