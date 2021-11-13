package actions

import (
	"context"
	"time"

	"guess_my_word/actions/test"
	"guess_my_word/internal/datastore"
	"guess_my_word/internal/words"

	"github.com/gin-gonic/gin"
)

type mockWordStore struct {
	mockValidate  func(context.Context, string) bool
	mockGetForDay func(context.Context, time.Time, string) (string, error)
}

func init() {
	wordStore = words.NewWordStore(datastore.NewMemory())
}

func (m *mockWordStore) Validate(ctx context.Context, word string) bool {
	return m.mockValidate(ctx, word)
}
func (m *mockWordStore) GetForDay(ctx context.Context, tm time.Time, mode string) (string, error) {
	return m.mockGetForDay(ctx, tm, mode)
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	AddHandlers(r, test.Templates, test.Assets)
	return r
}
