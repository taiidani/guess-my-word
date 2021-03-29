package actions

import (
	"context"
	"time"

	"guess_my_word/actions/test"

	"github.com/gin-gonic/gin"
)

type mockWordStore struct {
	mockValidate  func(context.Context, string) bool
	mockGetForDay func(context.Context, time.Time, string) (string, error)
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
