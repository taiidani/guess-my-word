package test

import "embed"

//go:embed templates/*.html
var Templates embed.FS

//go:embed assets
var Assets embed.FS
