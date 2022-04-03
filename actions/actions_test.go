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
	client := datastore.NewMemory()
	listStore = words.NewListStore(client)
	wordStore = words.NewWordStore(client)
}

func (m *mockWordStore) Validate(ctx context.Context, word string) bool {
	return m.mockValidate(ctx, word)
}
func (m *mockWordStore) GetForDay(ctx context.Context, tm time.Time, mode string) (model.Word, error) {
	return m.mockGetForDay(ctx, tm, mode)
}
func (m *mockWordStore) GetWord(ctx context.Context, key string) (model.Word, error) {
	return m.GetWord(ctx, key)
}
func (m *mockWordStore) SetWord(ctx context.Context, key string, word model.Word) error {
	return m.SetWord(ctx, key, word)
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	AddHandlers(r)
	return r
}
