package main

import (
	"guess_my_word/actions"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	actions.AddHandlers(r)
	r.Run(":3000")
}
