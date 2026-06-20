# Architecture — Project AEGIS

> Last updated: 2024-12-01 • Phase 0

This document describes the internal architecture of Project AEGIS. For platform-specific caveats, see [PLATFORM_NOTES.md](PLATFORM_NOTES.md). For the plugin development guide, see [PLUGIN_DEVELOPMENT.md](PLUGIN_DEVELOPMENT.md).

---

## Table of Contents

- [1. System Overview](#1-system-overview)
- [2. Backend Architecture](#2-backend-architecture)
- [3. Frontend Architecture](#3-frontend-architecture)
- [4. Sidecar Process Model](#4-sidecar-process-model)
- [5. Plugin System](#5-plugin-system)
- [6. Data Flow](#6-data-flow)
- [7. Cross-Platform Strategy](#7-cross-platform-strategy)
- [8. Hardware Tier Detection](#8-hardware-tier-detection)
- [9. Security Model](#9-security-model)

---

## 1. System Overview

Project AEGIS is a single self-contained binary that serves a browser-based UI and orchestrates a constellation of sidecar processes for heavy-lifting tasks. The design optimizes for three constraints simultaneously: **offline operation**, **cross-platform portability**, and **low-spec hardware support**.

### High-Level Diagram

```
┌──────────────────────────────────────────────────────────────────────────┐
│                            USER'S BROWSER                                │
│                                                                          │
│   ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐ │
│   │Dashboard │  │Knowledge │  │AI Chat   │  │Mesh Msgs │  │  Maps    │ │
│   └──────────┘  └──────────┘  └──────────┘  └──────────┘  └──────────┘ │
│                   React + TypeScript + Tailwind + shadcn/ui              │
└──────────────────────────────┬───────────────────────────────────────────┘
                               │
                     HTTP (REST) / WebSocket
                               │
┌──────────────────────────────▼───────────────────────────────────────────┐
│                         GO BACKEND (aegis binary)                        │
│                                                                          │
│  ┌─────────────────────────────────────────────────────────────────────┐ │
│  │                          API Layer                                  │ │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌───────────────────┐  │ │
│  │  │REST Mux  │  │WebSocket │  │Static FS │  │Plugin Route Mux   │  │ │
│  │  │          │  │Hub       │  │(embed)   │  │                   │  │ │
│  │  └──────────┘  └──────────┘  └──────────┘  └───────────────────┘  │ │
│  └─────────────────────────────────────────────────────────────────────┘ │
│                                                                          │
│  ┌─────────────────────────────────────────────────────────────────────┐ │
│  │                       Core Services                                 │ │
│  │  ┌──────────────┐  ┌───────────────┐  ┌──────────────────────────┐ │ │
│  │  │ Orchestrator │  │ Module        │  │ Resource Profiler        │ │ │
│  │  │ (sidecar     │  │ Registry      │  │ (OS-specific via build   │ │ │
│  │  │  lifecycle)  │  │               │  │  tags)                   │ │ │
│  │  └──────────────┘  └───────────────┘  └──────────────────────────┘ │ │
│  │  ┌──────────────┐  ┌───────────────┐  ┌──────────────────────────┐ │ │
│  │  │ Plugin       │  │ Power Budget  │  │ SQLite Store             │ │ │
│  │  │ Loader       │  │ Manager       │  │ (modernc.org/sqlite)     │ │ │
│  │  └──────────────┘  └───────────────┘  └──────────────────────────┘ │ │
│  └─────────────────────────────────────────────────────────────────────┘ │
│                                                                          │
│  ┌─────────────────────────────────────────────────────────────────────┐ │
│  │                         Modules                                     │ │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ │ │
│  │  │knowledge │ │aiengine  │ │meshmsg   │ │sdrmonitor│ │notes     │ │ │
│  │  │          │ │          │ │          │ │          │ │          │ │ │
│  │  └──────────┘ └──────────┘ └──────────┘ └──────────┘ └──────────┘ │ │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ │ │
│  │  │maps      │ │medical   │ │reticulum │ │peersync  │ │beacon    │ │ │
│  │  └──────────┘ └──────────┘ └──────────┘ └──────────┘ └──────────┘ │ │
│  └─────────────────────────────────────────────────────────────────────┘ │
└──────────┬───────────────┬───────────────┬───────────────┬──────────────┘
           │               │               │               │
    ┌──────▼──────┐ ┌──────▼──────┐ ┌──────▼──────┐ ┌─────▼──────┐
    │ kiwix-serve │ │llama-server │ │  rtl_fm /   │ │   rnsd     │
    │             │ │             │ │  rtl_power  │ │ (Reticulum)│
    │ (per-OS     │ │ (per-OS     │ │ (per-OS     │ │ (per-OS    │
    │  prebuilt)  │ │  prebuilt)  │ │  prebuilt)  │ │  prebuilt) │
    └─────────────┘ └─────────────┘ └─────────────┘ └────────────┘
           ▲               ▲               ▲               ▲
           └───────────────┴───────────────┴───────────────┘
                        Sidecar Processes
                   (managed as child processes)
```

### Data Directory Layout

All persistent state resides in a single portable directory:

```
aegis-data/
├── config.json              # User preferences, module toggles
├── aegis.db                 # SQLite database (notes, settings, vectors)
├── models/                  # Downloaded LLM .gguf files
├── content-packs/           # ZIM files, .pmtiles maps, medical DB
├── identity/                # Reticulum identity keys, mesh node config
├── plugins/                 # User-installed plugins
└── logs/                    # Rotating log files
```

---

## 2. Backend Architecture

### Language & Toolchain

- **Go 1.22+** with standard `cmd/` / `internal/` layout
- **No cgo** — all dependencies are pure Go to enable clean cross-compilation
- **`modernc.org/sqlite`** for database access (pure-Go SQLite implementation)

### Dependency Injection

AEGIS uses **constructor-based dependency injection** — no package-level mutable state. All services are constructed and wired together in `main.go`:

```go
func main() {
    // 1. Detect hardware
    profiler := resourceprofiler.New()
    tier := profiler.DetectTier()

    // 2. Open data store
    store, err := store.Open(dataDir)

    // 3. Initialize core services
    orchestrator := orchestrator.New(tier, sidecarDir)
    powerMgr := powermanager.New(profiler)
    pluginLoader := orchestrator.PluginLoader(pluginDir)

    // 4. Register modules (each checks tier compatibility)
    modules := []module.Module{
        knowledge.New(store, orchestrator),
        aiengine.New(store, orchestrator, tier),
        meshmsg.New(store),
        // ...
    }

    // 5. Build API router and serve
    api := api.New(store, modules, pluginLoader)
    http.ListenAndServe(":8080", api.Handler())
}
```

### Module Interface

Every module implements a common interface:

```go
type Module interface {
    // Metadata
    ID() string
    Name() string
    Domain() string               // "knowledge", "survival", "comms", "ai", "power"
    MinTier() resourceprofiler.Tier

    // Lifecycle
    Init(ctx context.Context) error
    Start(ctx context.Context) error
    Stop(ctx context.Context) error
    Status() ModuleStatus

    // API
    RegisterRoutes(mux *http.ServeMux)
}
```

### Error Handling

All errors are wrapped with context using `fmt.Errorf("...: %w", err)` to preserve error chains. The top-level API layer translates errors into appropriate HTTP status codes.

### Coding Conventions

- Table-driven tests for any logic with more than two branches
- No package-level mutable state — dependencies injected via constructors
- `internal/` packages prevent leaking implementation details
- Platform-divergent code isolated behind interfaces with build-tag-suffixed files

---

## 3. Frontend Architecture

### Stack

| Technology | Purpose |
|-----------|---------|
| React 18+ | Component framework |
| TypeScript | Type safety |
| Tailwind CSS | Utility-first styling |
| shadcn/ui | Accessible UI primitives |
| Vite | Build toolchain |

### Embedding

The frontend is compiled to static files (`frontend/dist/`) and embedded into the Go binary using `go:embed`:

```go
//go:embed all:frontend/dist
var frontendFS embed.FS
```

This means the final AEGIS binary is **entirely self-contained** — no external HTML/CSS/JS files to distribute.

### Module Mirroring

The frontend mirrors the backend module structure:

```
frontend/src/
├── modules/
│   ├── knowledge/     # Knowledge library UI
│   ├── maps/          # Offline maps viewer
│   ├── aiengine/      # AI chat interface
│   ├── meshmsg/       # Mesh messaging UI
│   ├── notes/         # Notes editor
│   ├── medical/       # Medical triage (with disclaimers)
│   ├── sdrmonitor/    # SDR spectrum display
│   └── ...
├── components/        # Shared UI components
├── hooks/             # Shared React hooks
└── lib/               # Utility functions
```

Each frontend module folder corresponds 1:1 with a backend module of the same name.

### UI Design Principles

- **Minimal, premium, sentence case** — this is a tool for use under stress, not a marketing site
- **No unnecessary animation** — functional transitions only
- **Responsive** — works on desktop and tablet screen sizes
- **Offline-first** — all UI assets embedded, no CDN calls

---

## 4. Sidecar Process Model

### Why Not Docker?

Docker Desktop is a heavy, license-encumbered dependency on Windows/macOS and is unavailable on Raspberry Pi Zero–class hardware. AEGIS uses **sidecar processes** instead — native per-OS prebuilt binaries managed as child processes.

### Sidecar Lifecycle

```
┌──────────┐     spawn      ┌──────────┐
│          │ ──────────────► │          │
│  AEGIS   │                 │ Sidecar  │
│  Backend │ ◄────────────── │ Process  │
│          │   health check  │          │
│          │ ──────────────► │          │
│          │     teardown    │          │
└──────────┘                 └──────────┘
```

The **Orchestrator** manages sidecar lifecycles:

1. **Discovery** — Locates the correct per-OS binary in `sidecars/<tool>/<os>/`
2. **Spawn** — Starts the sidecar as a child process with configured arguments
3. **Health Check** — Polls the sidecar's health endpoint (HTTP) or process status
4. **Restart** — Automatically restarts crashed sidecars (with backoff)
5. **Teardown** — Gracefully stops all sidecars on AEGIS shutdown (SIGTERM → SIGKILL timeout)

### Current Sidecars

| Sidecar | Purpose | Communication |
|---------|---------|---------------|
| `kiwix-serve` | ZIM file serving | HTTP proxy |
| `llama-server` | LLM inference + embeddings | HTTP (OpenAI-compatible API) |
| `rtl_fm` / `rtl_power` | SDR signal reception | stdout pipe (JSON output) |
| `rnsd` | Reticulum encrypted P2P | Local HTTP bridge |

### Sidecar Directory Layout

```
sidecars/
├── kiwix-serve/
│   ├── windows/kiwix-serve.exe
│   ├── linux/kiwix-serve
│   └── macos/kiwix-serve
├── llama-server/
│   ├── windows/llama-server.exe
│   ├── linux/llama-server
│   └── macos/llama-server
├── rtl-sdr/
│   ├── windows/rtl_fm.exe
│   ├── linux/rtl_fm
│   └── macos/rtl_fm
└── rnsd/
    ├── windows/rnsd.exe
    ├── linux/rnsd
    └── macos/rnsd
```

---

## 5. Plugin System

### Overview

AEGIS supports a **manifest-driven plugin system** that allows adding new modules without modifying core code. Plugins are discovered at startup from the `aegis-data/plugins/` directory.

### Plugin Structure

```
my-plugin/
├── manifest.json        # Plugin metadata and configuration
├── bin/                 # Per-OS binaries (if sidecar-based)
│   ├── windows/
│   ├── linux/
│   └── macos/
├── frontend/            # Optional frontend bundle
│   └── index.js
└── data/                # Plugin-specific data files
```

### Manifest Schema

Each plugin declares its requirements via `manifest.json`:

```json
{
  "id": "my-plugin",
  "name": "My Custom Plugin",
  "version": "1.0.0",
  "domain": "utility",
  "description": "A custom plugin for AEGIS",
  "entrypoint": {
    "embedded": false,
    "windows": "bin/windows/my-plugin.exe",
    "linux": "bin/linux/my-plugin",
    "darwin": "bin/macos/my-plugin"
  },
  "hardware_tier_min": "standard",
  "api_routes": ["/api/plugins/my-plugin/*"],
  "permissions": ["filesystem"],
  "resource_limits": {
    "max_memory_mb": 256,
    "max_cpu_percent": 25
  },
  "ui_module": "frontend/index.js"
}
```

### Plugin Discovery Flow

```
Startup
  │
  ├─► Scan aegis-data/plugins/
  │     │
  │     ├─► For each manifest.json:
  │     │     ├─► Validate against schema
  │     │     ├─► Check hardware_tier_min vs detected tier
  │     │     ├─► Check permissions against policy
  │     │     └─► Register API routes
  │     │
  │     └─► Report invalid/incompatible plugins in logs
  │
  └─► Plugins available via /api/plugins
```

> 📖 See [PLUGIN_DEVELOPMENT.md](PLUGIN_DEVELOPMENT.md) for the full development guide.

---

## 6. Data Flow

### Request Lifecycle (REST)

```
Browser                     Go Backend                  Sidecar/DB
  │                            │                            │
  │  GET /api/knowledge/search │                            │
  │ ──────────────────────────►│                            │
  │                            │  proxy to kiwix-serve      │
  │                            │ ──────────────────────────►│
  │                            │                            │
  │                            │  ◄─── search results ──── │
  │                            │                            │
  │  ◄── JSON response ─────  │                            │
  │                            │                            │
```

### Request Lifecycle (WebSocket — AI Chat)

```
Browser                     Go Backend                  llama-server
  │                            │                            │
  │  WS /api/ai/chat           │                            │
  │ ──────────────────────────►│                            │
  │                            │                            │
  │  { "message": "..." }      │                            │
  │ ──────────────────────────►│                            │
  │                            │  1. Generate embeddings    │
  │                            │ ──────────────────────────►│
  │                            │  ◄── embedding vector ──── │
  │                            │                            │
  │                            │  2. Vector search (SQLite) │
  │                            │  → retrieve context chunks │
  │                            │                            │
  │                            │  3. LLM completion w/      │
  │                            │     context                │
  │                            │ ──────────────────────────►│
  │                            │  ◄── stream tokens ─────── │
  │                            │                            │
  │  ◄── stream tokens ─────  │                            │
  │                            │                            │
```

### Mesh Message Flow

```
LoRa Radio ◄──► Serial Port ◄──► Go Backend ◄──► Browser
                                     │
                                     ▼
                              SQLite (store)
```

---

## 7. Cross-Platform Strategy

### Build Tags

Platform-specific code is **never** scattered through business logic. Instead, each OS-divergent feature uses Go build tags:

```
internal/resourceprofiler/
├── profiler.go             # Shared interface + types
├── profiler_windows.go     # //go:build windows
├── profiler_linux.go       # //go:build linux
└── profiler_darwin.go      # //go:build darwin
```

The shared interface is defined once:

```go
// profiler.go
type Profiler interface {
    DetectTier() Tier
    CPUCores() int
    TotalMemoryMB() int
    AvailableMemoryMB() int
    DiskFreeGB() float64
}
```

Each `_<os>.go` file provides the platform-specific implementation.

### Build Matrix

| Target | GOOS | GOARCH | Notes |
|--------|------|--------|-------|
| Windows desktop | `windows` | `amd64` | Primary desktop target |
| Linux desktop | `linux` | `amd64` | Debian/Ubuntu-class |
| Linux ARM (Pi) | `linux` | `arm64` | Raspberry Pi OS, Pi 4+, Pi Zero 2W |
| macOS (Apple Silicon) | `darwin` | `arm64` | macOS 13+ |

### Platform-Specific Implementations

| Feature | Windows | Linux | macOS |
|---------|---------|-------|-------|
| Hardware detection | WMI queries | `/sys/class/`, `/proc/` | `sysctl`, IOKit |
| Battery status | WMI `Win32_Battery` | `/sys/class/power_supply/` | `pmset -g batt` |
| Serial ports | COM ports (`COM3`, etc.) | `/dev/ttyUSB*`, `/dev/ttyACM*` | `/dev/cu.*` |
| Process management | `os/exec` + Win32 API | `os/exec` + signals | `os/exec` + signals |

---

## 8. Hardware Tier Detection

### Tiers

AEGIS classifies hardware into three tiers and adjusts module availability accordingly:

| Tier | RAM | CPU | Use Case |
|------|-----|-----|----------|
| **Minimum** | < 2 GB | ≤ 2 cores | RPi Zero 2W, old netbooks |
| **Standard** | 2–12 GB | 2–8 cores | Mid-range laptops, RPi 4 |
| **Optimal** | > 12 GB | > 8 cores | Modern laptops/desktops |

### Detection Algorithm

```
1. Query total physical RAM
2. Query CPU core count
3. Query available disk space
4. Classify into tier:
   - RAM < 2GB OR cores ≤ 2        → Minimum
   - RAM 2-12GB AND cores 2-8      → Standard
   - RAM > 12GB AND cores > 8      → Optimal
5. Override: if available RAM < 512MB at runtime → downgrade one tier
```

### Module Availability by Tier

| Module | Minimum | Standard | Optimal |
|--------|:-------:|:--------:|:-------:|
| Knowledge Library | ✅ | ✅ | ✅ |
| Offline Maps | ✅ | ✅ | ✅ |
| Notes | ✅ | ✅ | ✅ |
| Skill Trees | ✅ | ✅ | ✅ |
| Celestial Nav | ✅ | ✅ | ✅ |
| Mesh Messaging | ✅ | ✅ | ✅ |
| Position Beacon | ✅ | ✅ | ✅ |
| Medical Triage | ✅ | ✅ | ✅ |
| Data Tools | ✅ | ✅ | ✅ |
| AI Assistant | ❌ | ✅ | ✅ |
| SDR Monitor | ❌ | ✅ | ✅ |
| Plant / Fungi ID | ❌ | ✅ | ✅ |
| Local Peer Sync | ❌ | ✅ | ✅ |
| Encrypted P2P | ❌ | ❌ | ✅ |
| Agent Orchestrator | ❌ | ❌ | ✅ |

When a module is unavailable due to tier restrictions, the dashboard shows it as **disabled** with a clear reason rather than hiding it or failing silently.

---

## 9. Security Model

### Design Philosophy

AEGIS is a **single-user offline tool**. There is no authentication or authorization system in v1, matching the philosophy of the tools it builds upon.

### Network Exposure

- By default, AEGIS binds to `localhost:8080` — not accessible from other machines
- If configured for LAN access, the operator should understand the implications
- No TLS in v1 (would require certificate management that conflicts with zero-install)
- Document network controls in deployment guide rather than building auth

### Sidecar Isolation

- Sidecars inherit the OS-level permissions of the AEGIS process
- Each sidecar has configurable resource limits (`max_memory_mb`, `max_cpu_percent`)
- Plugins declare required permissions (`network`, `serial`, `filesystem`, `subprocess`)
- The orchestrator enforces declared resource limits via OS-level controls where available

### Data at Rest

- SQLite database is not encrypted by default (single-user device)
- Reticulum identity keys stored in `aegis-data/identity/` — the user is responsible for protecting this directory
- No telemetry, no phone-home, no analytics — ever

---

*This document is a living reference. It will be updated as the architecture evolves through each implementation phase.*
