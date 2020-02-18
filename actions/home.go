package actions

import (
	"guess_my_word/internal/data"
	"time"

	"github.com/gobuffalo/buffalo"
)

// HomeHandler is a default handler to serve up
// a home page.
func HomeHandler(c buffalo.Context) error {
	var yesterday *data.Date
	yesterDate := time.Now().UTC().AddDate(0, 0, -1)
	yesterday, err := data.LoadDate(yesterDate)
	if err != nil {
		c.Logger().Infof("Could not load yesterday from data store. Defaulting to temporary new date: %s", err)
		yesterday = data.NewDate(yesterDate)
		c.Logger().Debugf("Generated temporary yesterDate word '%s' from '%s'", yesterday.Word, yesterday.ID)
	}
	c.Set("yesterday", yesterday)

	return c.Render(200, r.HTML("index.html"))
}
