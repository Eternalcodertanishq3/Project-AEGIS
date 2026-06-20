package api

import (
	"io/fs"
	"net/http"
)

// NewRouter creates the HTTP router with all API routes and the embedded frontend.
// frontendFS should be the sub-filesystem rooted at the frontend dist directory.
func NewRouter(deps *Deps, frontendFS fs.FS) http.Handler {
	mux := http.NewServeMux()

	// ─── API Routes ──────────────────────────────────────────────────
	mux.HandleFunc("GET /api/health", healthHandler)
	mux.HandleFunc("GET /api/system/profile", systemProfileHandler(deps))
	mux.HandleFunc("GET /api/system/power", systemPowerHandler(deps))
	mux.HandleFunc("GET /api/modules", modulesListHandler(deps))
	mux.HandleFunc("POST /api/modules/{id}/enable", moduleEnableHandler(deps))
	mux.HandleFunc("POST /api/modules/{id}/disable", moduleDisableHandler(deps))
	mux.HandleFunc("GET /api/plugins", pluginsListHandler)
	mux.HandleFunc("POST /api/plugins/install", pluginsInstallHandler)

	// ─── Knowledge Module Routes ─────────────────────────────────────
	if deps.KnowledgeHandlers != nil {
		deps.KnowledgeHandlers.RegisterRoutes(mux)
	}

	// ─── Maps Module Routes ──────────────────────────────────────────
	if deps.MapsHandlers != nil {
		deps.MapsHandlers.RegisterRoutes(mux)
	}

	// ─── AI Module Routes ─────────────────────────────────────────────
	if deps.AIHandlers != nil {
		deps.AIHandlers.RegisterRoutes(mux)
	}

	// ─── Embedded Frontend ───────────────────────────────────────────
	fileServer := http.FileServerFS(frontendFS)
	mux.Handle("/", fileServer)

	return mux
}

