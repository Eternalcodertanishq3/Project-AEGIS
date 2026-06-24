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

	// ─── Notes Module Routes ──────────────────────────────────────────
	if deps.NotesHandlers != nil {
		deps.NotesHandlers.RegisterRoutes(mux)
	}

	// ─── Medical Module Routes ────────────────────────────────────────
	if deps.MedicalHandlers != nil {
		deps.MedicalHandlers.RegisterRoutes(mux)
	}

	// ─── Data Tools Module Routes ─────────────────────────────────────
	if deps.DataToolsHandlers != nil {
		deps.DataToolsHandlers.RegisterRoutes(mux)
	}

	// ─── Skill Trees Module Routes ────────────────────────────────────
	if deps.SkillTreesHandlers != nil {
		deps.SkillTreesHandlers.RegisterRoutes(mux)
	}

	// ─── Celestial Navigation Module Routes ───────────────────────────
	if deps.CelestialHandlers != nil {
		deps.CelestialHandlers.RegisterRoutes(mux)
	}

	// ─── Plant/Fungi ID Module Routes ─────────────────────────────────
	if deps.PlantIDHandlers != nil {
		deps.PlantIDHandlers.RegisterRoutes(mux)
	}

	// ─── Mesh Messaging Module Routes ─────────────────────────────────
	if deps.MeshHandlers != nil {
		deps.MeshHandlers.RegisterRoutes(mux)
	}

	// ─── Encrypted P2P Module Routes ──────────────────────────────────
	if deps.P2PHandlers != nil {
		deps.P2PHandlers.RegisterRoutes(mux)
	}

	// ─── SDR Monitor Module Routes ────────────────────────────────────
	if deps.SDRHandlers != nil {
		deps.SDRHandlers.RegisterRoutes(mux)
	}

	// ─── Local Peer Sync Module Routes ────────────────────────────────
	if deps.PeerSyncHandlers != nil {
		deps.PeerSyncHandlers.RegisterRoutes(mux)
	}

	// ─── Position Beacon Module Routes ────────────────────────────────
	if deps.BeaconHandlers != nil {
		deps.BeaconHandlers.RegisterRoutes(mux)
	}

	// ─── Embedded Frontend ───────────────────────────────────────────
	fileServer := http.FileServerFS(frontendFS)
	mux.Handle("/", fileServer)

	return mux
}

