package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"aegis/backend"
	"aegis/backend/internal/api"
	"aegis/backend/internal/modules/ai"
	"aegis/backend/internal/modules/knowledge"
	"aegis/backend/internal/modules/maps"
	"aegis/backend/internal/orchestrator"
	"aegis/backend/internal/powermanager"
	"aegis/backend/internal/resourceprofiler"
	"aegis/backend/internal/store"
)

const banner = `
   ╔═══════════════════════════════════════════════════╗
   ║     _    _____ ____ ___ ____                      ║
   ║    / \  | ____/ ___|_ _/ ___|                     ║
   ║   / _ \ |  _|| |  _ | |\___ \                     ║
   ║  / ___ \| |__| |_| || | ___) |                    ║
   ║ /_/   \_\_____\____|___|____/                     ║
   ║                                                   ║
   ║  Survival Computer · Phase 0 · v0.1.0             ║
   ╚═══════════════════════════════════════════════════╝
`

func main() {
	// ─── Flags ────────────────────────────────────────────────────────
	port := flag.Int("port", 8080, "HTTP server port")
	dataDir := flag.String("data-dir", "./aegis-data", "Path to the data directory")
	flag.Parse()

	fmt.Print(banner)

	// ─── Database ────────────────────────────────────────────────────
	db, err := store.New(*dataDir)
	if err != nil {
		log.Fatalf("database init failed: %v", err)
	}
	defer db.Close()
	log.Println("✓ Database initialized")

	// ─── Resource Profiler ───────────────────────────────────────────
	profiler := resourceprofiler.NewProfiler()
	profile, err := profiler.DetectProfile()
	if err != nil {
		log.Fatalf("hardware detection failed: %v", err)
	}
	log.Printf("✓ Hardware profile: %s", profile)

	// ─── Power Manager ───────────────────────────────────────────────
	pm := powermanager.NewPowerManager()
	powerStatus, err := pm.GetPowerStatus()
	if err != nil {
		log.Printf("⚠ Power detection failed (non-fatal): %v", err)
	} else {
		log.Printf("✓ Power status: %s (battery: %d%%)", powerStatus.Status, powerStatus.BatteryPercent)
	}

	// ─── Orchestrator ────────────────────────────────────────────────
	orch := orchestrator.New()
	modules := orch.ListModules()
	log.Printf("✓ Module registry: %d modules registered", len(modules))

	// ─── Knowledge Module ────────────────────────────────────────────
	kiwixMgr := knowledge.NewKiwixManager(*dataDir, 9080)
	kiwixMgr.DiscoverZIMFiles()
	zimFiles := kiwixMgr.GetZIMFiles()
	log.Printf("✓ Knowledge module: %d ZIM files discovered", len(zimFiles))
	
	if _, err := kiwixMgr.FindSidecar(); err == nil {
		if err := kiwixMgr.Start(); err != nil {
			log.Printf("⚠ Knowledge module: failed to start kiwix-serve: %v", err)
		}
	} else {
		log.Printf("⚠ Knowledge module: %v", err)
	}
	
	knowledgeHandlers := knowledge.NewHandlers(kiwixMgr)

	// ─── Maps Module ─────────────────────────────────────────────────
	mapMgr := maps.NewMapManager(*dataDir)
	mapMgr.DiscoverMapFiles()
	mapFiles := mapMgr.GetMapFiles()
	log.Printf("✓ Maps module: %d PMTiles files discovered", len(mapFiles))
	mapsHandlers := maps.NewHandlers(mapMgr)

	// ─── AI Module ───────────────────────────────────────────────────
	aiMgr := ai.NewAIManager(*dataDir, 8081)
	aiMgr.DiscoverModels()
	aiModels := aiMgr.GetModels()
	log.Printf("✓ AI module: %d GGUF models discovered", len(aiModels))

	if _, err := aiMgr.FindSidecar(); err == nil {
		log.Printf("✓ AI module: llama-server found")
	} else {
		log.Printf("⚠ AI module: %v", err)
	}

	aiHandlers := ai.NewHandlers(aiMgr)

	// ─── HTTP Router ─────────────────────────────────────────────────
	deps := &api.Deps{
		Profiler:          profiler,
		PowerManager:      pm,
		Orchestrator:      orch,
		KnowledgeHandlers: knowledgeHandlers,
		MapsHandlers:      mapsHandlers,
		AIHandlers:        aiHandlers,
	}
	handler := api.NewRouter(deps, backend.EmbeddedFrontend)

	// ─── Server ──────────────────────────────────────────────────────
	addr := fmt.Sprintf(":%d", *port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// ─── Startup Summary ─────────────────────────────────────────────
	fmt.Println("─────────────────────────────────────────────────────")
	fmt.Printf("  Port:           %d\n", *port)
	fmt.Printf("  Data Dir:       %s\n", *dataDir)
	fmt.Printf("  OS:             %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("  Hardware Tier:  %s\n", profile.Tier)
	fmt.Printf("  RAM:            %d MB\n", profile.TotalRAMMB)
	fmt.Printf("  CPU Cores:      %d\n", profile.CPUCores)
	fmt.Println("─────────────────────────────────────────────────────")
	log.Printf("🌐 AEGIS listening on http://localhost:%d", *port)

	// ─── Graceful Shutdown ───────────────────────────────────────────
	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.ListenAndServe()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-quit:
		log.Printf("Received signal %v, shutting down gracefully…", sig)
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	kiwixMgr.Stop()
	aiMgr.Stop()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("forced shutdown: %v", err)
	}
	log.Println("AEGIS shut down cleanly. Stay alive out there. 🏕️")
}
