package app

import "github.com/gin-gonic/gin"

type errorBag struct {
	baseBag
	Message error
}

func errorResponse(c *gin.Context, code int, err error) {
	data := errorBag{
		Message: err,
	}

	c.HTML(code, "error.gohtml", data)
}
