package actions

import (
	"context"
	"fmt"
	"guess_my_word/internal/datastore"
	"guess_my_word/internal/model"
	"guess_my_word/internal/words"
	"html/template"
	"io"
	"io/fs"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type wordClient interface {
	Validate(context.Context, string) bool
	GetForDay(context.Context, time.Time, string) (model.Word, error)
}

var dataStore words.Store
var wordStore wordClient

func init() {
	dataStore = datastore.NewRedis(os.Getenv("REDIS_ADDR"))
	wordStore = words.NewWordStore(dataStore)
}

// AddHandlers will add the application handlers to the HTTP server
func AddHandlers(r *gin.Engine, templates fs.FS, assets fs.FS) (err error) {
	if gin.IsDebugging() {
		err = addHandlersStaticPreProduction(r)
	} else {
		err = addHandlersStaticProduction(r, templates, assets)
	}
	if err != nil {
		return fmt.Errorf("could not register static handlers: %w", err)
	}

	r.Use(middlewareStandardHeaders())
	r.GET("/", HomeHandler)
	r.GET("/reveal", RevealHandler)
	r.GET("/ping", PingHandler)
	r.GET("/guess", GuessHandler)
	r.GET("/hint", HintHandler)
	return nil
}

func loadTemplate(templates fs.FS) (*template.Template, error) {
	t := template.New("")

	err := fs.WalkDir(templates, "templates", func(name string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := templates.Open(name)
		if err != nil {
			return err
		}
		defer file.Close()

		h, err := io.ReadAll(file)
		if err != nil {
			return err
		}
		t, err = t.New(info.Name()).Parse(string(h))
		if err != nil {
			return err
		}

		return nil
	})

	return t, err
}

func addHandlersStaticProduction(r *gin.Engine, templates, assets fs.FS) error {
	t, err := loadTemplate(templates)
	if err != nil {
		return fmt.Errorf("could not load templates: %w", err)
	}
	r.SetHTMLTemplate(t)

	r.GET("/assets/*filepath", newStaticHandler(assets))
	return nil
}

func addHandlersStaticPreProduction(r *gin.Engine) error {
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./assets")
	return nil
}

// convertUTCToLocal will take a given time in UTC and convert it to a given user's timezone
// TZ for PDT (-7:00) is a positive 420, so SUBTRACT that from the unix timestamp
func convertUTCToUser(t time.Time, tz int) time.Time {
	ret := t.In(time.FixedZone("User", tz*-1))
	ret = ret.Add(time.Minute * -1 * time.Duration(tz))
	return ret
}
