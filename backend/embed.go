package backend

import (
	"embed"
	"io/fs"
	"log"
)

//go:embed all:frontend/dist
var embeddedFrontendRaw embed.FS

// EmbeddedFrontend is the frontend filesystem with the "frontend/dist" prefix stripped,
// so files are served from the root (e.g., index.html, not frontend/dist/index.html).
var EmbeddedFrontend fs.FS

func init() {
	var err error
	EmbeddedFrontend, err = fs.Sub(embeddedFrontendRaw, "frontend/dist")
	if err != nil {
		log.Fatalf("failed to create sub-filesystem for embedded frontend: %v", err)
	}
}
