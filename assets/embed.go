package assets


import (
  "embed"
	"io/fs"
)


//go:embed migrations
var MigrationsFS embed.FS

//go:embed static
var staticFS embed.FS

//go:embed templates
var TemplatesFS embed.FS


func StaticAssets() fs.FS {
	s, _ := fs.Sub(staticFS, "static")
	return s
}
