package actions

import (
	"guess_my_word/internal/data"

	"github.com/dgraph-io/badger/v2"
)

func (as *ActionSuite) Test_HomeHandler() {
	opts := badger.DefaultOptions("")
	opts.InMemory = true
	db, teardown := data.NewBadgerBackend(&opts)
	data.SetBackend(db)
	defer teardown()

	res := as.HTML("/").Get()
	as.Equal(200, res.Code)
	as.Contains(res.Body.String(), "Guess My Word")
}
