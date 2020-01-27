package actions

import "github.com/gobuffalo/buffalo"

import "time"

// HomeHandler is a default handler to serve up
// a home page.
func HomeHandler(c buffalo.Context) error {
	yesterday, _ := generateWord(time.Now().UTC().AddDate(0, 0, -1))
	c.Set("yesterday", yesterday)

	return c.Render(200, r.HTML("index.html"))
}
