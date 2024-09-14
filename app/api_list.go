package app

import (
	"guess_my_word/internal/model"
	"net/http"
)

type listReply struct {
	Error string `json:"error,omitempty"`
}

// ListHandler is an API handler to provide a list to a user.
func ListHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// case http.MethodGet:
	// 	if list, err := listGet(r); err != nil {
	// 		renderJson(w, 400, listReply{Error: err.Error()})
	// 	} else {
	// 		renderJson(w, 200, list)
	// 	}
	// case http.MethodPost:
	// 	if err := listUpdate(r); err != nil {
	// 		renderJson(w, 400, listReply{Error: err.Error()})
	// 	}
	// case http.MethodPut:
	// 	if err := listCreate(r); err != nil {
	// 		renderJson(w, 400, listReply{Error: err.Error()})
	// 	}
	// case http.MethodDelete:
	// 	if err := listDelete(r); err != nil {
	// 		renderJson(w, 400, listReply{Error: err.Error()})
	// 	}
	default:
		reply := listReply{}
		reply.Error = "unexpected HTTP method"
		renderJson(w, 400, reply)
	}
}

// ListsHandler is an API handler to enumerate all lists.
func ListsHandler(w http.ResponseWriter, r *http.Request) {
	lists, err := listStore.GetLists(r.Context())
	if err != nil {
		renderJson(w, 500, listReply{Error: err.Error()})
		return
	}

	ret := []model.List{}
	for _, list := range lists {
		add, err := listStore.GetList(r.Context(), list)
		if err != nil {
			renderJson(w, 500, listReply{Error: err.Error()})
			return
		}

		// Clear the words out to reduce network load
		add.Words = []string{}
		ret = append(ret, add)
	}

	renderJson(w, 200, ret)
}

// func listCreate(r *http.Request) error {
// 	type list struct {
// 		List model.List `form:"list"`
// 	}

// 	request := list{}
// 	if err := c.ShouldBind(&request); err != nil {
// 		return fmt.Errorf("invalid request received: %w", err)
// 	}

// 	return listStore.CreateList(r.Context(), request.List.Name, request.List)
// }

// func listUpdate(r *http.Request) error {
// 	type list struct {
// 		List model.List `form:"list"`
// 	}

// 	request := list{}
// 	if err := c.ShouldBind(&request); err != nil {
// 		return fmt.Errorf("invalid request received: %w", err)
// 	}

// 	return listStore.UpdateList(r.Context(), request.List.Name, request.List)
// }

// func listDelete(r *http.Request) error {
// 	type list struct {
// 		Name string `form:"name"`
// 	}

// 	request := list{}
// 	if err := c.ShouldBind(&request); err != nil {
// 		return fmt.Errorf("invalid request received: %w", err)
// 	}

// 	return listStore.DeleteList(r.Context(), request.Name)
// }

// func listGet(r *http.Request) (model.List, error) {
// 	type list struct {
// 		Name string `form:"name"`
// 	}

// 	request := list{}
// 	if err := c.ShouldBind(&request); err != nil {
// 		return model.List{}, fmt.Errorf("invalid request received: %w", err)
// 	}

// 	return listStore.GetList(r.Context(), request.Name)
// }
