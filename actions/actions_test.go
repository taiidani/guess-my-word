package actions

import (
	"context"
	"time"

	"guess_my_word/internal/datastore"
	"guess_my_word/internal/model"
	"guess_my_word/internal/words"

	"github.com/gin-gonic/gin"
)

type mockWordStore struct {
	mockValidate  func(context.Context, string) bool
	mockGetForDay func(context.Context, time.Time, string) (model.Word, error)
}

func init() {
	dataStore = datastore.NewMemory()
	wordStore = words.NewWordStore(dataStore)
}

func (m *mockWordStore) Validate(ctx context.Context, word string) bool {
	return m.mockValidate(ctx, word)
}
func (m *mockWordStore) GetForDay(ctx context.Context, tm time.Time, mode string) (model.Word, error) {
	return m.mockGetForDay(ctx, tm, mode)
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	AddHandlers(r)
	return r
}
