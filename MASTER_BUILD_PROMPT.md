# Project AEGIS — Master Build Prompt

**How to use this document:** this is a complete build specification written to be handed directly to a coding agent (Claude Code or equivalent). It is structured so the agent can start at Phase 0 and proceed phase by phase without needing clarification on naming, tech stack, or license — those are fixed below. Where a third-party dependency lacks support for a target OS, the agent should document the gap in `docs/PLATFORM_NOTES.md` and degrade that module gracefully rather than blocking the build.

---

## 1. Mission

Build **Project AEGIS**, a self-contained offline survival and resilience computer that runs identically on Windows, Linux, and macOS from a single portable binary, with no installation and no admin rights required for its core functionality. It extends the proven offline-knowledge model of Project N.O.M.A.D. (Kiwix, local AI, offline maps) with two capabilities NOMAD does not attempt: off-grid communication (LoRa mesh, encrypted P2P, SDR monitoring) and true cross-platform, low-spec portability.

## 2. Non-negotiable constraints

- **Runs on Windows 10+, Linux (Debian/Ubuntu-class, incl. Raspberry Pi OS arm64), and macOS 13+** from one codebase.
- **No mandatory installer.** The core path (knowledge library, notes, AI assistant, mesh messaging) must run by double-clicking a downloaded/USB-resident binary. Heavier optional modules (Education/Kolibri, full Docker-based NOMAD parity mode) may require a separate install — document this explicitly, don't silently assume it.
- **No mandatory internet connection** after initial content/model download. The system must detect and report offline status rather than failing silently.
- **Low-spec compatible** down to a Raspberry Pi Zero 2W (1GB RAM). The system must auto-detect hardware tier and disable modules it can't support rather than crashing or hanging.
- **USB-portable.** All persistent state lives in one data directory that can be copied between machines.
- **Modular and plugin-extensible.** New modules are addable without modifying core code, via a manifest-driven plugin system.
- **No authentication by default**, matching NOMAD's philosophy — this is a single-user offline tool. Document network-level controls (don't expose port 8080 beyond localhost/LAN without understanding the implications) rather than building auth into v1.

## 3. Naming

Working name: **Project AEGIS** (binary/package name: `aegis`). This is a placeholder — rename via find-replace on `aegis`/`AEGIS` throughout the codebase if a different name is chosen later. Do not hardcode the name in ways that resist renaming (e.g., avoid baking it into database schema field names).

## 4. Architecture decisions and rationale

| Decision | Choice | Why |
|---|---|---|
| Backend language | **Go** | Cross-compiles natively to windows/linux/darwin × amd64/arm64 from one machine, no per-OS build toolchain needed, produces a single static binary, `embed.FS` lets the frontend ship inside the binary |
| Database | **SQLite via `modernc.org/sqlite`** (pure Go, no cgo) | cgo-based SQLite drivers complicate cross-compilation; pure-Go avoids needing a C cross-compiler per target |
| Frontend | **React + TypeScript + Tailwind + shadcn/ui**, built to static files and embedded in the Go binary | Browser-based UI is OS-agnostic by construction; one build artifact, no platform-specific UI code |
| Module execution model | **Sidecar processes**, not Docker | Docker Desktop is a heavy, license-encumbered dependency on Windows/macOS and unavailable on Pi Zero–class hardware. Each heavy tool (Kiwix, llama.cpp, RTL-SDR) ships as a per-OS prebuilt binary that the Go process spawns, monitors, and tears down — no container runtime required. Docker remains available as an *optional* power-user overlay on Linux/x86 for users who want full NOMAD parity (Kolibri, Qdrant), but it is never required for the core build. |
| Local AI inference | **llama.cpp (`llama-server`)**, not Ollama | Ollama requires its own separate installer/background service, which breaks "no install, runs from USB." `llama-server` is a single portable binary that loads a `.gguf` file directly. An optional setting to point at an existing Ollama or OpenAI-compatible host is still supported, mirroring NOMAD's own v1.31 "remote host" feature, for users who already run one. |
| Vector search (RAG) | Embeddings via `llama-server`'s `/embedding` endpoint, stored in SQLite, brute-force cosine search | Avoids a Qdrant/Docker dependency. Brute-force search is adequate at the corpus sizes this device handles (thousands, not millions, of chunks); revisit only if profiling shows it's a bottleneck. |
| LoRa mesh | **Native Go client** against the published [Meshtastic protobuf schema](https://github.com/meshtastic/protobufs), transport via `go.bug.st/serial` | Keeps mesh messaging inside the single Go binary with no Python dependency; `go.bug.st/serial` is a maintained, genuinely cross-platform (Windows COM ports, Linux `/dev/ttyUSB*`, macOS `/dev/cu.*`) serial library. |
| Encrypted P2P | **Reticulum (RNS)** via a bundled per-OS sidecar | RNS's reference implementation is Python. There is no mature Go reimplementation, so this is the one module that ships a bundled, PyInstaller-frozen `rnsd` binary per OS (end users still never install Python themselves) talking to the Go backend over a local HTTP bridge. This also gives interoperability with existing Reticulum apps (Nomad Network, Sideband) for free. |
| SDR monitoring | **rtl-sdr tools (`rtl_fm`, `rtl_power`)**, receive-only, bundled per-OS | Official builds exist for all three target OSes. Receive-only by design avoids transmission licensing issues. **Windows caveat:** the RTL2832U dongle needs a one-time Zadig WinUSB driver swap — document this in `docs/PLATFORM_NOTES.md`, do not gloss over it. |
| Local peer/content sync | **libp2p (Go)**, replacing the originally-considered `batman-adv` | `batman-adv` is a Linux kernel module — it cannot run on Windows or macOS and is a hard blocker for the cross-platform requirement. libp2p is pure userspace, the same library IPFS is built on, and gives mDNS-based local peer discovery plus a transport for syncing content packs over an existing local WiFi/Ethernet network without any kernel-level mesh networking. |
| Education module (Kolibri) | **Optional Tier-2 add-on**, not bundled in the core image | Kolibri is a Python/Django application; bundling it per-OS adds real packaging weight for a module that's NOMAD-parity rather than core to the survival use case. Document it as a separate optional install or Docker overlay rather than overpromising it in the portable core. |
| Power/hardware detection | OS-specific code behind a common Go interface, using build tags (`_windows.go`, `_linux.go`, `_darwin.go`) | Windows: WMI queries. Linux: `/sys/class/power_supply`. macOS: `pmset`/IOKit. Idiomatic Go pattern for OS-divergent logic without runtime branching. |

## 5. Directory structure

```
project-aegis/
├── README.md
├── LICENSE
├── MASTER_BUILD_PROMPT.md
├── go.mod
├── Makefile
├── backend/
│   ├── cmd/aegis/main.go
│   ├── embed.go                      # go:embed of frontend/dist
│   └── internal/
│       ├── api/                      # REST + WebSocket handlers
│       ├── orchestrator/             # sidecar process lifecycle, plugin loader
│       ├── resourceprofiler/
│       │   ├── profiler.go           # shared interface
│       │   ├── profiler_windows.go
│       │   ├── profiler_linux.go
│       │   └── profiler_darwin.go
│       ├── powermanager/             # same per-OS pattern as resourceprofiler
│       ├── store/                    # SQLite data layer
│       └── modules/
│           ├── knowledge/            # Kiwix sidecar wrapper
│           ├── maps/                 # static .pmtiles serving
│           ├── aiengine/             # llama-server lifecycle, RAG, agent router
│           ├── notes/                # native Go+SQLite
│           ├── datatools/            # CyberChef static vendor
│           ├── medical/              # triage rules engine
│           ├── plantid/              # ONNX vision inference
│           ├── skilltrees/           # static content loader
│           ├── celestialnav/         # pure-math nav/weather calc
│           ├── meshmsg/              # Meshtastic protobuf-over-serial client
│           ├── reticulum/            # bridge to rnsd sidecar
│           ├── sdrmonitor/           # rtl-sdr sidecar wrapper
│           ├── peersync/             # libp2p local discovery + content sync
│           └── beacon/               # APRS-style position broadcast
├── frontend/
│   ├── src/
│   │   ├── modules/                  # one folder per backend module, mirrored
│   │   ├── components/
│   │   ├── hooks/
│   │   └── lib/
│   └── package.json
├── sidecars/                         # per-OS prebuilt third-party binaries
│   ├── kiwix-serve/{windows,linux,macos}/
│   ├── llama-server/{windows,linux,macos}/
│   ├── rtl-sdr/{windows,linux,macos}/
│   └── rnsd/{windows,linux,macos}/
├── plugin-sdk/
│   ├── manifest.schema.json
│   ├── scaffold.sh
│   ├── scaffold.ps1
│   └── examples/
├── content-packs/
│   ├── zim-survival-pack/
│   ├── maps-regional/                # .pmtiles files
│   └── medical-db/
├── scripts/
│   ├── fetch-sidecars.sh             # downloads pinned sidecar binary versions
│   ├── build-windows.sh
│   ├── build-linux.sh
│   ├── build-macos.sh
│   └── package-usb.sh                # assembles the final USB folder layout
├── boot/
│   └── pi-image/                     # optional: dedicated always-on base-station image
├── docs/
│   ├── ARCHITECTURE.md
│   ├── PLATFORM_NOTES.md             # Windows driver caveats, macOS Gatekeeper, etc.
│   ├── PLUGIN_DEVELOPMENT.md
│   └── COMMS_PROTOCOLS.md
├── .github/workflows/build.yml
└── tests/
```

## 6. Module specifications

| Module | Domain | Core tech | Cross-platform approach |
|---|---|---|---|
| Knowledge library | Knowledge | Kiwix (ZIM files) | Bundled `kiwix-serve` sidecar per OS, managed as a child process |
| Offline maps | Knowledge | ProtoMaps (`.pmtiles`) + MapLibre GL JS | Served as static files directly from the Go backend, no sidecar |
| AI assistant + RAG | AI | llama.cpp `llama-server` | Bundled sidecar per OS; embeddings via its API; vectors in SQLite |
| Education | Knowledge | Kolibri | Optional Tier-2 add-on, documented separately, not in core image |
| Notes | Knowledge | Native Go + React + SQLite | No external dependency |
| Data tools | Knowledge | CyberChef (vendored static build) | Pure client-side JS, no backend dependency |
| Medical triage | Survival | Rules engine (Go) + offline formulary DB | Native Go, bundled SQLite dataset, heavily disclaimed UI |
| Plant/fungi ID | Survival | ONNX vision model | Bundled per-OS ONNX runtime sidecar; disclaimed, not authoritative ID |
| Skill trees | Survival | Static Markdown/JSON content | No dependency |
| Celestial nav & weather | Survival | Pure math (Go) | No dependency, no GPS required |
| Mesh messaging | Comms | Meshtastic protobufs over serial | Native Go, `go.bug.st/serial` |
| Encrypted P2P | Comms | Reticulum (RNS) | Bundled per-OS `rnsd` sidecar, local HTTP bridge |
| SDR monitor | Comms | `rtl_fm` / `rtl_power`, receive-only | Bundled per-OS sidecars; Windows needs Zadig driver swap |
| Local peer/content sync | Comms | libp2p (Go) | mDNS discovery + diff/sync over existing local WiFi/Ethernet |
| Position beacon | Comms | APRS-style over LoRa | Built on the mesh messaging module |
| Multi-agent orchestrator | AI | Local-LLM intent router (Go) | Routes a query to the right module's API; falls back to keyword routing on Minimum-tier hardware where running a router model isn't affordable |
| Plugin SDK | Extensibility | JSON Schema manifest + scaffold CLI | `scaffold.sh` (bash) and `scaffold.ps1` (PowerShell) cover all three OSes |
| Resource profiler | Power | OS-specific hardware probes | Go build tags, one file per OS behind a shared interface |
| Power budget manager | Power | Battery/AC state polling | Same per-OS pattern; throttles/disables sidecars under low battery |
| System benchmark | Power | Native Go benchmark suite | CPU/RAM/disk scoring, optional LoRa range test, SDR sensitivity check |

## 7. Command Center API contract

```
GET    /api/health
GET    /api/system/profile              # hardware tier, detected capabilities
GET    /api/system/power                # battery/AC status
GET    /api/modules                     # list installed/available modules + status
POST   /api/modules/:id/enable
POST   /api/modules/:id/disable

GET    /api/knowledge/search?q=
GET    /api/knowledge/article/:zimId/*path

WS     /api/ai/chat                     # streaming chat
GET    /api/ai/models
POST   /api/ai/models/:id/download

GET    /api/mesh/nodes
POST   /api/mesh/messages
WS     /api/mesh/messages/stream

GET    /api/comms/sdr/scan?band=

GET    /api/sync/peers
POST   /api/sync/push

GET    /api/plugins
POST   /api/plugins/install
```

## 8. Plugin manifest schema

`plugin-sdk/manifest.schema.json`:

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "title": "AEGIS Plugin Manifest",
  "type": "object",
  "required": ["id", "name", "version", "domain", "entrypoint", "hardware_tier_min"],
  "properties": {
    "id": { "type": "string", "pattern": "^[a-z0-9-]+$" },
    "name": { "type": "string" },
    "version": { "type": "string", "pattern": "^\\d+\\.\\d+\\.\\d+$" },
    "domain": { "type": "string", "enum": ["knowledge", "survival", "comms", "ai", "power", "utility"] },
    "description": { "type": "string", "maxLength": 280 },
    "entrypoint": {
      "type": "object",
      "description": "Per-OS binary path, or embedded: true for native Go modules",
      "properties": {
        "windows": { "type": "string" },
        "linux": { "type": "string" },
        "darwin": { "type": "string" },
        "embedded": { "type": "boolean" }
      }
    },
    "hardware_tier_min": { "type": "string", "enum": ["minimum", "standard", "optimal"] },
    "api_routes": { "type": "array", "items": { "type": "string" } },
    "permissions": {
      "type": "array",
      "items": { "type": "string", "enum": ["network", "serial", "filesystem", "subprocess"] }
    },
    "resource_limits": {
      "type": "object",
      "properties": {
        "max_memory_mb": { "type": "integer" },
        "max_cpu_percent": { "type": "integer" }
      }
    },
    "ui_module": { "type": "string", "description": "Path to frontend bundle entry, relative to module dir" }
  }
}
```

## 9. Cross-platform build pipeline

`.github/workflows/build.yml` (matrix build, no per-OS source forks):

```yaml
name: build
on: [push]
jobs:
  build:
    strategy:
      matrix:
        include:
          - os: windows-latest
            goos: windows
            goarch: amd64
          - os: ubuntu-latest
            goos: linux
            goarch: amd64
          - os: ubuntu-latest
            goos: linux
            goarch: arm64
          - os: macos-latest
            goos: darwin
            goarch: arm64
    runs-on: ${{ matrix.os }}
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with: { go-version: '1.22' }
      - run: cd frontend && npm ci && npm run build
      - run: GOOS=${{ matrix.goos }} GOARCH=${{ matrix.goarch }} go build -o dist/aegis-${{ matrix.goos }}-${{ matrix.goarch }} ./backend/cmd/aegis
```

`scripts/package-usb.sh` assembles the final portable layout: the three OS binaries, `sidecars/` for the host's OS only (to keep the drive small — or all three if optimizing for "any computer" over drive size), `content-packs/`, and an empty `aegis-data/` directory, onto a target USB drive path passed as an argument.

## 10. Phased implementation roadmap

| Phase | Deliverable | Definition of done |
|---|---|---|
| 0 — Scaffold | Go backend serving the embedded (empty) React frontend; SQLite store initialized; hardware-tier stub | `go run ./backend/cmd/aegis` opens a working dashboard identically via `GOOS=windows/linux/darwin go build` |
| 1 — Knowledge core | Kiwix sidecar integration, search + article view in frontend; ProtoMaps static serving + map UI | Search a bundled test ZIM file and render an article and a map tile, on all three OSes |
| 2 — AI assistant | `llama-server` sidecar lifecycle, model download manager, streaming chat UI, embedding-based RAG over a test ZIM | A chat query returns a streamed response grounded in retrieved chunks |
| 3 — Agent orchestrator | Query router directing to Knowledge vs. AI vs. (stubbed) Comms modules | A natural-language query is routed to the correct module's handler |
| 4 — Comms | Meshtastic Go client (serial), `rnsd` sidecar + bridge, libp2p local peer sync, RTL-SDR sidecar wrapper | Two physical LoRa nodes exchange a message; two AEGIS instances on the same LAN discover each other and sync a content pack |
| 5 — Survival modules | Medical triage rules engine, plant-ID inference, skill trees, celestial nav calculator | Each module returns correct output against a fixed test-case set; medical and plant-ID UIs carry explicit non-authoritative disclaimers |
| 6 — Power & resilience | Per-OS resource profiler and power manager, system benchmark | On a simulated Minimum-tier profile, AI assistant and SDR modules auto-disable; benchmark produces a comparable score across OSes |
| 7 — Plugin SDK & polish | Manifest loader, `scaffold.sh`/`scaffold.ps1`, full docs set, CI matrix build, USB packaging script | A scaffolded "hello world" plugin is discovered and enabled without restarting the binary |

## 11. Coding conventions

- **Go:** standard `cmd/`/`internal/` layout; errors wrapped with `fmt.Errorf("...: %w", err)`; no package-level mutable state — dependencies constructed and injected in `main.go`; table-driven tests for anything with more than two branches.
- **React:** functional components and hooks only; Tailwind utility classes; shadcn/ui for interactive primitives; one frontend module folder per backend module, same name.
- **UI tone:** minimal, premium, sentence case throughout, no unnecessary animation — this is a tool someone may be using under stress, not a marketing site.
- **Platform-divergent code** always goes behind a shared interface with Go build-tag-suffixed files (`_windows.go`/`_linux.go`/`_darwin.go`), never inline OS branches scattered through business logic.

## 12. v1.0 acceptance criteria

- [ ] Core binary (knowledge + notes + AI assistant) runs unmodified on Windows 10/11, Debian 12/Ubuntu 22.04+, and macOS 13+
- [ ] Runs from a USB drive on a machine with no prior install and no admin rights, for the core path
- [ ] Mesh messaging verified between two physical LoRa-equipped devices
- [ ] Local peer sync verified between two AEGIS instances on the same network, no internet required
- [ ] Resource profiler correctly detects a Minimum-tier device and disables AI/SDR modules rather than failing
- [ ] Zero outbound network calls after initial content/model download, outside explicit user-triggered updates
- [ ] A scaffolded plugin goes from `scaffold.sh my-module` to a running, discoverable module in under 5 minutes
- [ ] `docs/PLATFORM_NOTES.md` documents every known per-OS caveat (Windows RTL-SDR driver, Kolibri non-bundling, etc.) — none are silently glossed over

## 13. Directive

Begin at Phase 0. Work phase by phase in the order above; each phase's definition of done must pass before starting the next. Naming, license, and core tech stack are fixed by this document — do not re-litigate them. Where a target-OS gap is discovered that isn't already called out in Section 4, add it to `docs/PLATFORM_NOTES.md` and choose graceful degradation over blocking the build.
