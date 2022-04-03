package actions

import (
	"fmt"
	"guess_my_word/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type listReply struct {
	Error string `json:"error,omitempty"`
}

// ListHandler is an API handler to provide a list to a user.
func ListHandler(c *gin.Context) {
	switch c.Request.Method {
	case http.MethodGet:
		if list, err := listGet(c); err != nil {
			c.JSON(400, listReply{Error: err.Error()})
		} else {
			c.JSON(200, list)
		}
	case http.MethodPost:
		if err := listUpdate(c); err != nil {
			c.JSON(400, listReply{Error: err.Error()})
		}
	case http.MethodPut:
		if err := listCreate(c); err != nil {
			c.JSON(400, listReply{Error: err.Error()})
		}
	case http.MethodDelete:
		if err := listDelete(c); err != nil {
			c.JSON(400, listReply{Error: err.Error()})
		}
	default:
		reply := listReply{}
		reply.Error = "unexpected HTTP method"
		c.JSON(400, reply)
	}
}

// ListsHandler is an API handler to enumerate all lists.
func ListsHandler(c *gin.Context) {
	lists, err := listStore.GetLists(c)
	if err != nil {
		c.JSON(500, listReply{Error: err.Error()})
		return
	}

	ret := []model.List{}
	for _, list := range lists {
		add, err := listStore.GetList(c, list)
		if err != nil {
			c.JSON(500, listReply{Error: err.Error()})
			return
		}

		// Clear the words out to reduce network load
		add.Words = []string{}
		ret = append(ret, add)
	}

	c.JSON(200, ret)
}

func listCreate(c *gin.Context) error {
	type list struct {
		List model.List `form:"list"`
	}

	request := list{}
	if err := c.ShouldBind(&request); err != nil {
		return fmt.Errorf("invalid request received: %w", err)
	}

	return listStore.CreateList(c, request.List.Name, request.List)
}

func listUpdate(c *gin.Context) error {
	type list struct {
		List model.List `form:"list"`
	}

	request := list{}
	if err := c.ShouldBind(&request); err != nil {
		return fmt.Errorf("invalid request received: %w", err)
	}

	return listStore.UpdateList(c, request.List.Name, request.List)
}

func listDelete(c *gin.Context) error {
	type list struct {
		Name string `form:"name"`
	}

	request := list{}
	if err := c.ShouldBind(&request); err != nil {
		return fmt.Errorf("invalid request received: %w", err)
	}

	return listStore.DeleteList(c, request.Name)
}

func listGet(c *gin.Context) (model.List, error) {
	type list struct {
		Name string `form:"name"`
	}

	request := list{}
	if err := c.ShouldBind(&request); err != nil {
		return model.List{}, fmt.Errorf("invalid request received: %w", err)
	}

	return listStore.GetList(c, request.Name)
}
