# Plugin Development Guide — Project AEGIS

> Build and distribute custom modules for AEGIS without modifying core code.

---

## Table of Contents

- [1. Overview](#1-overview)
- [2. Plugin Manifest Schema](#2-plugin-manifest-schema)
- [3. Scaffolding a New Plugin](#3-scaffolding-a-new-plugin)
- [4. Plugin Lifecycle](#4-plugin-lifecycle)
- [5. API Route Registration](#5-api-route-registration)
- [6. Frontend Module Integration](#6-frontend-module-integration)
- [7. Resource Limits and Permissions](#7-resource-limits-and-permissions)
- [8. Testing Plugins](#8-testing-plugins)
- [9. Distribution](#9-distribution)

---

## 1. Overview

AEGIS plugins are self-contained modules that extend the system with new capabilities. The plugin system is **manifest-driven** — each plugin declares its metadata, requirements, and capabilities in a `manifest.json` file, and AEGIS discovers and loads plugins at startup.

### Design Principles

- **No core modification required** — plugins are discovered from the `aegis-data/plugins/` directory
- **Declarative configuration** — the manifest defines everything AEGIS needs to know
- **Hardware-aware** — plugins declare a minimum hardware tier; AEGIS won't load plugins the device can't support
- **Cross-platform** — plugins with sidecar binaries must provide per-OS builds
- **Sandboxed** — plugins declare required permissions; the orchestrator enforces resource limits

### Plugin Types

| Type | Description | Example |
|------|-------------|---------|
| **Embedded (Go)** | Compiled into the AEGIS binary as a Go module | Notes, Celestial Nav |
| **Sidecar** | External binary managed as a child process | Kiwix, llama-server |
| **Frontend-only** | Pure client-side JavaScript module | Data Tools (CyberChef) |
| **Hybrid** | Sidecar backend + frontend UI | SDR Monitor, AI Chat |

---

## 2. Plugin Manifest Schema

Every plugin must include a `manifest.json` file at its root. The manifest conforms to the JSON Schema defined in `plugin-sdk/manifest.schema.json`.

### Complete Manifest Example

```json
{
  "id": "weather-station",
  "name": "Weather Station",
  "version": "1.0.0",
  "domain": "survival",
  "description": "Reads data from a connected USB weather station and displays current conditions, trends, and forecasts.",
  "entrypoint": {
    "embedded": false,
    "windows": "bin/windows/weather-station.exe",
    "linux": "bin/linux/weather-station",
    "darwin": "bin/macos/weather-station"
  },
  "hardware_tier_min": "minimum",
  "api_routes": [
    "/api/plugins/weather-station/current",
    "/api/plugins/weather-station/history",
    "/api/plugins/weather-station/forecast"
  ],
  "permissions": [
    "serial",
    "filesystem"
  ],
  "resource_limits": {
    "max_memory_mb": 128,
    "max_cpu_percent": 10
  },
  "ui_module": "frontend/index.js"
}
```

### Field Reference

| Field | Type | Required | Description |
|-------|------|:--------:|-------------|
| `id` | `string` | ✅ | Unique plugin identifier. Lowercase alphanumeric + hyphens only (`^[a-z0-9-]+$`). |
| `name` | `string` | ✅ | Human-readable display name. |
| `version` | `string` | ✅ | Semantic version (`MAJOR.MINOR.PATCH`). |
| `domain` | `enum` | ✅ | One of: `knowledge`, `survival`, `comms`, `ai`, `power`, `utility`. |
| `description` | `string` | ❌ | Short description, max 280 characters. |
| `entrypoint` | `object` | ✅ | Defines how the plugin runs. See below. |
| `hardware_tier_min` | `enum` | ✅ | Minimum required tier: `minimum`, `standard`, or `optimal`. |
| `api_routes` | `string[]` | ❌ | List of API route patterns the plugin will register. |
| `permissions` | `enum[]` | ❌ | Required permissions: `network`, `serial`, `filesystem`, `subprocess`. |
| `resource_limits` | `object` | ❌ | Resource constraints for the plugin process. |
| `ui_module` | `string` | ❌ | Path to the frontend bundle entry point, relative to the plugin root. |

### Entrypoint Object

The `entrypoint` field defines how AEGIS starts the plugin:

```json
// Sidecar plugin — per-OS binary paths
{
  "entrypoint": {
    "embedded": false,
    "windows": "bin/windows/my-tool.exe",
    "linux": "bin/linux/my-tool",
    "darwin": "bin/macos/my-tool"
  }
}
```

```json
// Embedded plugin — compiled into the AEGIS binary
{
  "entrypoint": {
    "embedded": true
  }
}
```

```json
// Frontend-only plugin — no backend process
{
  "entrypoint": {
    "embedded": false
  }
}
```

---

## 3. Scaffolding a New Plugin

### Using the Scaffold CLI

AEGIS provides scaffold scripts to generate the plugin boilerplate:

**Linux / macOS:**

```bash
./plugin-sdk/scaffold.sh my-plugin --domain survival --tier minimum
```

**Windows (PowerShell):**

```powershell
.\plugin-sdk\scaffold.ps1 my-plugin -Domain survival -Tier minimum
```

### Generated Structure

```
aegis-data/plugins/my-plugin/
├── manifest.json          # Pre-filled manifest
├── bin/                   # Per-OS binaries (empty, ready for your builds)
│   ├── windows/
│   ├── linux/
│   └── macos/
├── frontend/              # Frontend module stub
│   ├── index.js
│   ├── MyPlugin.tsx       # React component template
│   └── styles.css
├── data/                  # Plugin-specific data files
├── README.md              # Plugin documentation template
└── Makefile               # Build targets for the plugin
```

### Manual Setup

If you prefer not to use the scaffold, create the directory structure manually:

1. Create a directory under `aegis-data/plugins/<your-plugin-id>/`
2. Create `manifest.json` (see schema above)
3. Add your binary/frontend assets
4. Restart AEGIS — the plugin will be discovered automatically

---

## 4. Plugin Lifecycle

### Discovery

On startup, AEGIS scans `aegis-data/plugins/` for directories containing a `manifest.json`:

```
┌──────────────┐
│   Startup    │
└──────┬───────┘
       │
       ▼
┌──────────────┐     ┌─────────────────────────┐
│ Scan plugins │────►│ For each manifest.json:  │
│   directory  │     │  1. Parse & validate     │
└──────────────┘     │  2. Check hardware tier  │
                     │  3. Verify permissions   │
                     │  4. Register routes      │
                     └─────────┬───────────────┘
                               │
                    ┌──────────┴──────────┐
                    │                     │
                    ▼                     ▼
             ┌────────────┐        ┌────────────┐
             │  Enabled   │        │  Disabled  │
             │ (tier OK)  │        │ (tier too  │
             │            │        │   low or   │
             │            │        │   invalid) │
             └────────────┘        └────────────┘
```

### State Machine

Each plugin progresses through these states:

```
  Discovered → Validated → Enabled → Running → Stopped
                  │                     │
                  ▼                     ▼
               Invalid              Crashed → Restarting
                                                  │
                                                  ▼ (max retries)
                                               Failed
```

| State | Description |
|-------|-------------|
| **Discovered** | Manifest found in plugins directory |
| **Validated** | Manifest passes schema validation |
| **Invalid** | Manifest fails validation — logged, plugin skipped |
| **Enabled** | Hardware tier check passed, plugin is eligible to run |
| **Disabled** | Hardware tier insufficient or user-disabled |
| **Running** | Sidecar process is active and healthy |
| **Stopped** | Gracefully stopped (user action or shutdown) |
| **Crashed** | Sidecar process exited unexpectedly |
| **Restarting** | Orchestrator is restarting the sidecar (with backoff) |
| **Failed** | Max restart attempts exceeded — manual intervention needed |

### Lifecycle Hooks

For embedded (Go) plugins, the module interface provides lifecycle hooks:

```go
type Module interface {
    Init(ctx context.Context) error    // Called once during startup
    Start(ctx context.Context) error   // Called to begin operation
    Stop(ctx context.Context) error    // Called for graceful shutdown
    Status() ModuleStatus              // Polled periodically by the dashboard
}
```

For sidecar plugins, the orchestrator manages the process lifecycle:

1. **Init** — Verify the binary exists for the current OS
2. **Start** — Spawn the process with configured arguments
3. **Health Check** — Poll the sidecar's health endpoint (configurable interval)
4. **Restart** — On crash, restart with exponential backoff (1s, 2s, 4s, 8s, max 60s)
5. **Stop** — Send SIGTERM (Unix) or `TerminateProcess` (Windows), wait 5s, then SIGKILL

---

## 5. API Route Registration

### Route Namespacing

All plugin API routes are namespaced under `/api/plugins/<plugin-id>/`:

```
/api/plugins/weather-station/current
/api/plugins/weather-station/history
/api/plugins/weather-station/forecast
```

This prevents collisions with core AEGIS routes and other plugins.

### Sidecar Route Proxying

For sidecar-based plugins, AEGIS acts as a **reverse proxy**. The sidecar runs its own HTTP server on a dynamically assigned localhost port, and AEGIS routes matching requests to it:

```
Browser → AEGIS (:8080) → Plugin Sidecar (:random_port)
```

The orchestrator passes the assigned port to the sidecar via environment variable:

```
AEGIS_PLUGIN_PORT=52341
```

The sidecar should bind to `127.0.0.1:$AEGIS_PLUGIN_PORT`.

### Embedded Route Registration

For embedded Go plugins, routes are registered directly on the HTTP mux:

```go
func (m *WeatherStation) RegisterRoutes(mux *http.ServeMux) {
    mux.HandleFunc("/api/plugins/weather-station/current", m.handleCurrent)
    mux.HandleFunc("/api/plugins/weather-station/history", m.handleHistory)
    mux.HandleFunc("/api/plugins/weather-station/forecast", m.handleForecast)
}
```

### WebSocket Support

Plugins can register WebSocket endpoints for real-time streaming:

```json
{
  "api_routes": [
    "/api/plugins/my-plugin/ws"
  ]
}
```

AEGIS upgrades matching connections and proxies the WebSocket to the sidecar.

---

## 6. Frontend Module Integration

### Module Registration

Plugins with a `ui_module` field in their manifest are automatically registered in the AEGIS frontend. The UI module is loaded dynamically and rendered as a tab/page in the dashboard.

### Frontend Bundle Structure

```
frontend/
├── index.js              # Entry point — exports the module definition
├── MyPlugin.tsx           # Main React component
├── components/            # Plugin-specific components
│   ├── DataChart.tsx
│   └── StatusBadge.tsx
├── hooks/
│   └── usePluginData.ts
└── styles.css             # Plugin-specific styles
```

### Module Definition (index.js)

The entry point must export a module definition object:

```javascript
import MyPlugin from './MyPlugin';

export default {
  // Unique identifier — must match manifest.json "id"
  id: 'weather-station',

  // Display name for the navigation
  name: 'Weather Station',

  // Icon (Lucide icon name)
  icon: 'cloud-sun',

  // Domain category for grouping in the sidebar
  domain: 'survival',

  // The main React component
  component: MyPlugin,

  // Optional: navigation sub-items
  routes: [
    { path: '/weather', label: 'Current', component: MyPlugin },
    { path: '/weather/history', label: 'History', component: HistoryView },
  ],
};
```

### Using Core UI Components

Plugins can import shared UI components from the AEGIS frontend library:

```tsx
import { Card, Badge, StatusIndicator } from '@aegis/components';
import { useModuleStatus, useWebSocket } from '@aegis/hooks';

export default function WeatherStation() {
  const status = useModuleStatus('weather-station');
  const { data } = useWebSocket('/api/plugins/weather-station/ws');

  return (
    <Card>
      <StatusIndicator status={status} />
      <h2>Current conditions</h2>
      {/* ... */}
    </Card>
  );
}
```

### Styling Guidelines

- Use Tailwind utility classes consistent with the AEGIS design system
- Follow **sentence case** for all UI text
- **No unnecessary animation** — functional transitions only
- Medical and identification modules must include explicit **non-authoritative disclaimers**

---

## 7. Resource Limits and Permissions

### Permissions

Plugins must declare all required permissions in the manifest. The orchestrator validates these at startup:

| Permission | Grants | Example Use |
|-----------|--------|-------------|
| `network` | Outbound network access | Fetching updates, remote API |
| `serial` | Access to serial ports (COM / ttyUSB) | LoRa radio, weather sensor |
| `filesystem` | Read/write to plugin data directory | Local data storage, logs |
| `subprocess` | Spawn child processes | Running helper scripts |

Plugins are **denied** any permission not explicitly declared.

### Resource Limits

```json
{
  "resource_limits": {
    "max_memory_mb": 256,
    "max_cpu_percent": 25
  }
}
```

- **`max_memory_mb`** — The orchestrator monitors the sidecar's RSS memory. If exceeded, a warning is logged. Persistent overuse (>120% for >60 seconds) triggers a restart.
- **`max_cpu_percent`** — Enforced via OS-level controls where available (cgroups on Linux, Job Objects on Windows). On macOS, this is advisory.

### Power Budget Integration

When the Power Budget Manager detects low battery:

1. Plugins with `hardware_tier_min: "optimal"` are stopped first
2. Then `standard` tier plugins
3. `minimum` tier plugins are kept running as long as possible

Plugins can subscribe to power state changes to gracefully reduce their own resource usage before being stopped.

---

## 8. Testing Plugins

### Local Development

1. Create your plugin in `aegis-data/plugins/<your-plugin-id>/`
2. Start AEGIS — your plugin will be discovered
3. Check `GET /api/plugins` to verify your plugin appears
4. Monitor `aegis-data/logs/` for errors

### Simulating Hardware Tiers

You can force a specific tier for testing:

```bash
AEGIS_FORCE_TIER=minimum ./aegis
```

This lets you verify your plugin correctly handles being disabled on lower tiers.

### Manifest Validation

Validate your manifest against the schema before testing:

```bash
# Using ajv-cli (Node.js)
npx ajv validate -s plugin-sdk/manifest.schema.json -d aegis-data/plugins/my-plugin/manifest.json
```

### Integration Test Checklist

- [ ] Plugin appears in `GET /api/plugins` with correct metadata
- [ ] Plugin starts successfully on all target platforms
- [ ] API routes return expected responses
- [ ] Frontend module renders in the dashboard
- [ ] Plugin handles graceful shutdown (SIGTERM)
- [ ] Plugin is correctly disabled on lower hardware tiers
- [ ] Resource limits are respected under load
- [ ] Plugin recovers from sidecar crash (if applicable)

---

## 9. Distribution

### Packaging

Distribute your plugin as a zip archive containing the full plugin directory:

```
weather-station-1.0.0.zip
└── weather-station/
    ├── manifest.json
    ├── bin/
    │   ├── windows/weather-station.exe
    │   ├── linux/weather-station
    │   └── macos/weather-station
    ├── frontend/
    │   └── index.js
    └── data/
```

### Installation

Users install plugins by extracting the archive into `aegis-data/plugins/` and restarting AEGIS:

```bash
unzip weather-station-1.0.0.zip -d aegis-data/plugins/
```

### Future: Hot Reload

> 🚧 Planned for Phase 7: Plugins will be discoverable and loadable without restarting the AEGIS binary. The `POST /api/plugins/install` endpoint will handle extraction, validation, and activation in a single step.

---

*For questions about plugin development, see the [ARCHITECTURE.md](ARCHITECTURE.md) document for system-level context, or open an issue on the project repository.*
