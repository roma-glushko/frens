package ui

import (
	"embed"
	"io/fs"
)

//go:embed all:dist
var distFS embed.FS

// Assets returns the embedded UI assets filesystem.
// The filesystem is rooted at the dist directory.
func Assets() (fs.FS, error) {
	return fs.Sub(distFS, "dist")
}
