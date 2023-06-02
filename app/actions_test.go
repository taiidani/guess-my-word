package app

import (
	"guess_my_word/internal/datastore"
	"guess_my_word/internal/words"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
)

func init() {
	client := datastore.NewMemory()
	listStore = words.NewListStore(client)
	wordStore = words.NewWordStore(client)
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	if err := SetupTemplates(r); err != nil {
		panic(err)
	}

	sessionClient := memstore.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("guessmyword", sessionClient))

	_ = AddHandlers(r)
	return r
}
