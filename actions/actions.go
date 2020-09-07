package actions

import (
	"html/template"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/markbates/pkger"
)

// AddHandlers will add the application handlers to the HTTP server
func AddHandlers(r *gin.Engine) {
	if gin.IsDebugging() {
		addHandlersStaticPreProduction(r)
	} else {
		addHandlersStaticProduction(r)
	}

	r.Use(middlewareStandardHeaders())
	r.GET("/", HomeHandler)
	r.GET("/ping", PingHandler)
	r.GET("/guess", GuessHandler)
	r.GET("/hint", HintHandler)
}

func loadTemplate() (*template.Template, error) {
	t := template.New("")

	err := pkger.Walk("/templates", func(name string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.HasSuffix(name, ".html") {
			return nil
		}

		file, err := pkger.Open(name)
		if err != nil {
			return err
		}
		defer file.Close()

		h, err := ioutil.ReadAll(file)
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

func addHandlersStaticProduction(r *gin.Engine) {
	pkger.Include("/templates")
	t, err := loadTemplate()
	if err != nil {
		panic(err)
	}
	r.SetHTMLTemplate(t)

	pkger.Include("/assets")
	r.GET("/assets/*filepath", StaticHandler)
}

func addHandlersStaticPreProduction(r *gin.Engine) {
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./assets")
}
