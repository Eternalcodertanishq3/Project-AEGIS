<p align="center">
  <h1 align="center">🛡️ Project AEGIS</h1>
  <p align="center"><strong>A self-contained offline survival and resilience computer</strong></p>
  <p align="center">
    <a href="#quick-start">Quick Start</a> •
    <a href="#architecture">Architecture</a> •
    <a href="#modules">Modules</a> •
    <a href="docs/ARCHITECTURE.md">Full Docs</a> •
    <a href="#license">License</a>
  </p>
</p>

---

## Overview

Project AEGIS is a portable, cross-platform offline survival computer that runs identically on **Windows**, **Linux**, and **macOS** from a single binary — no installation and no admin rights required. It extends proven offline-knowledge systems with **off-grid communications** (LoRa mesh, encrypted P2P, SDR monitoring) and **true cross-platform, low-spec portability** down to a Raspberry Pi Zero 2W.

AEGIS is designed for scenarios where infrastructure cannot be trusted or has ceased to exist: natural disasters, remote expeditions, communication blackouts, or simply operating in areas with no internet coverage.

## Key Features

- **🔌 Zero Install** — Double-click and run. All state lives in a single portable data directory.
- **✈️ Fully Offline** — No internet required after initial content/model download.
- **🖥️ Cross-Platform** — Single codebase targeting Windows 10+, Linux (x64/ARM64), and macOS 13+.
- **📡 Off-Grid Comms** — LoRa mesh messaging (Meshtastic), encrypted P2P (Reticulum), SDR monitoring.
- **🤖 Local AI** — Offline LLM assistant with RAG over your knowledge library, powered by llama.cpp.
- **🗺️ Offline Maps** — ProtoMaps `.pmtiles` with MapLibre GL JS — full vector maps with no tile server.
- **📚 Knowledge Library** — Kiwix-powered ZIM file browser with full-text search.
- **🔋 Hardware-Aware** — Auto-detects hardware tier and gracefully disables modules it can't support.
- **🧩 Plugin System** — Manifest-driven extensibility — add modules without touching core code.
- **🔒 USB-Portable** — Copy the data directory between machines, run from a thumb drive.

## Quick Start

### Prerequisites

- [Go 1.22+](https://go.dev/dl/)
- [Node.js 20+](https://nodejs.org/) (for frontend build only)

### Build & Run

```bash
# Clone the repository
git clone https://github.com/project-aegis/aegis.git
cd aegis

# Build the frontend
cd frontend && npm ci && npm run build && cd ..

# Build and run the backend (embeds the frontend)
go run ./backend/cmd/aegis
```

Or use the Makefile:

```bash
make build   # Build frontend + backend
make run     # Build and run
```

AEGIS will start on `http://localhost:8080`. Open it in any modern browser.

### Cross-Compile for All Targets

```bash
make build-all
```

This produces binaries in `dist/` for:
- `windows/amd64`
- `linux/amd64`
- `linux/arm64`
- `darwin/arm64`

## Architecture

AEGIS uses a **Go backend** with an **embedded React frontend** and **sidecar processes** for heavy-lifting tools.

```
┌─────────────────────────────────────────────────┐
│                   Browser UI                     │
│          React + TypeScript + Tailwind           │
└──────────────────────┬──────────────────────────┘
                       │ HTTP / WebSocket
┌──────────────────────▼──────────────────────────┐
│               Go Backend (aegis)                 │
│  ┌──────────┐ ┌──────────┐ ┌──────────────────┐ │
│  │ REST API │ │ WS Hub   │ │ Module Registry  │ │
│  └──────────┘ └──────────┘ └──────────────────┘ │
│  ┌──────────┐ ┌──────────┐ ┌──────────────────┐ │
│  │ SQLite   │ │ Profiler │ │ Plugin Loader    │ │
│  └──────────┘ └──────────┘ └──────────────────┘ │
└──────┬───────────┬───────────┬──────────────────┘
       │           │           │
  ┌────▼────┐ ┌────▼────┐ ┌───▼─────┐
  │ kiwix   │ │ llama   │ │ rtl-sdr │  ... (sidecar processes)
  │ -serve  │ │ -server │ │ tools   │
  └─────────┘ └─────────┘ └─────────┘
```

**Key design decisions:**
- **No Docker required.** Sidecars are native per-OS prebuilt binaries managed as child processes.
- **No cgo.** Pure-Go SQLite (`modernc.org/sqlite`) enables clean cross-compilation.
- **No authentication.** Single-user offline tool. Network-level controls documented instead.
- **Build-tag isolation.** Platform-specific code uses `_windows.go` / `_linux.go` / `_darwin.go` files behind shared interfaces.

> 📖 See [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) for the full architecture deep-dive.

## Hardware Requirements

AEGIS auto-detects hardware and adjusts available modules accordingly.

| Tier | Example Device | RAM | Storage | Capabilities |
|------|---------------|-----|---------|-------------|
| **Minimum** | Raspberry Pi Zero 2W | 1 GB | 16 GB SD | Knowledge library, notes, maps, mesh messaging, skill trees, celestial nav |
| **Standard** | Laptop (2018+), RPi 4 | 4–8 GB | 64 GB+ | All Minimum + AI assistant (small models), SDR monitoring, plant ID |
| **Optimal** | Modern laptop / desktop | 16+ GB | 256 GB+ | All Standard + large AI models, full RAG, multi-agent orchestrator, Reticulum P2P |

## Modules

| Module | Domain | Tech | Min. Tier |
|--------|--------|------|-----------|
| Knowledge Library | Knowledge | Kiwix (ZIM files) | Minimum |
| Offline Maps | Knowledge | ProtoMaps + MapLibre GL JS | Minimum |
| AI Assistant + RAG | AI | llama.cpp `llama-server` | Standard |
| Notes | Knowledge | Native Go + SQLite | Minimum |
| Data Tools | Knowledge | CyberChef (vendored static) | Minimum |
| Medical Triage | Survival | Go rules engine + formulary DB | Minimum |
| Plant / Fungi ID | Survival | ONNX vision model | Standard |
| Skill Trees | Survival | Static Markdown / JSON | Minimum |
| Celestial Nav & Weather | Survival | Pure math (Go) | Minimum |
| Mesh Messaging | Comms | Meshtastic protobufs over serial | Minimum |
| Encrypted P2P | Comms | Reticulum (RNS) sidecar | Optimal |
| SDR Monitor | Comms | rtl_fm / rtl_power (receive-only) | Standard |
| Local Peer Sync | Comms | libp2p (Go) | Standard |
| Position Beacon | Comms | APRS-style over LoRa | Minimum |
| Agent Orchestrator | AI | Local LLM intent router | Optimal |
| Plugin SDK | Extensibility | JSON Schema manifest + scaffold CLI | Minimum |
| Resource Profiler | Power | OS-specific hardware probes | Minimum |
| Power Budget Manager | Power | Battery / AC state polling | Minimum |
| System Benchmark | Power | Native Go benchmark suite | Minimum |

## Development Setup

### Prerequisites

| Tool | Version | Purpose |
|------|---------|---------|
| Go | ≥ 1.22 | Backend compilation |
| Node.js | ≥ 20 | Frontend build toolchain |
| npm | ≥ 10 | Frontend dependency management |
| Make | any | Build automation (optional) |

### Getting Started

```bash
# 1. Clone the repo
git clone https://github.com/project-aegis/aegis.git
cd aegis

# 2. Install frontend dependencies
cd frontend && npm ci && cd ..

# 3. Run in development mode
make run
```

### Project Structure

```
project-aegis/
├── backend/          # Go backend (cmd/, internal/)
├── frontend/         # React + TypeScript frontend
├── sidecars/         # Per-OS prebuilt third-party binaries
├── plugin-sdk/       # Plugin manifest schema + scaffold tools
├── content-packs/    # ZIM files, maps, medical DB
├── scripts/          # Build and packaging scripts
├── docs/             # Architecture, platform notes, guides
├── boot/             # Optional RPi base-station image
└── tests/            # Test suites
```

### Useful Commands

```bash
make frontend     # Build React frontend
make backend      # Build Go binary for current OS
make build        # Build frontend + backend
make run          # Build and run
make build-all    # Cross-compile for all targets
make test         # Run Go test suite
make clean        # Remove build artifacts
```

## Documentation

| Document | Description |
|----------|-------------|
| [ARCHITECTURE.md](docs/ARCHITECTURE.md) | System architecture deep-dive |
| [PLATFORM_NOTES.md](docs/PLATFORM_NOTES.md) | OS-specific caveats and driver requirements |
| [PLUGIN_DEVELOPMENT.md](docs/PLUGIN_DEVELOPMENT.md) | Guide to building AEGIS plugins |
| [COMMS_PROTOCOLS.md](docs/COMMS_PROTOCOLS.md) | Communications protocol reference |
| [MASTER_BUILD_PROMPT.md](MASTER_BUILD_PROMPT.md) | Full build specification |

## Contributing

Project AEGIS is in active development. See the [phased roadmap](MASTER_BUILD_PROMPT.md#10-phased-implementation-roadmap) for current priorities.

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/my-module`)
3. Follow the [coding conventions](MASTER_BUILD_PROMPT.md#11-coding-conventions)
4. Write table-driven tests for any logic with more than two branches
5. Submit a pull request

## License

Project AEGIS is released under the [MIT License](LICENSE).

```
MIT License — Copyright (c) 2024-2025 Project AEGIS Contributors
```
