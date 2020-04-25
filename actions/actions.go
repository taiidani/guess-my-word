package actions

import (
	"html/template"
	"io/ioutil"
	"log"
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

	r.GET("/", HomeHandler)
	r.GET("/guess", GuessHandler)
}

func loadTemplate() (*template.Template, error) {
	t := template.New("")

	err := pkger.Walk("/templates", func(name string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		log.Println("Walking templates; found ", name, " with name ", info.Name())
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
