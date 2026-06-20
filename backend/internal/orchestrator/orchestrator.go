package orchestrator

import (
	"fmt"
	"sync"

	"aegis/backend/internal/resourceprofiler"
)

// Module represents a registered AEGIS module.
type Module struct {
	ID             string                       `json:"id"`
	Name           string                       `json:"name"`
	Domain         string                       `json:"domain"`
	Description    string                       `json:"description"`
	Status         string                       `json:"status"`
	HardwareMinTier resourceprofiler.HardwareTier `json:"hardware_min_tier"`
	Enabled        bool                         `json:"enabled"`
}

// Orchestrator manages the module registry and sidecar lifecycle.
type Orchestrator struct {
	mu      sync.RWMutex
	modules map[string]*Module
}

// New creates a new Orchestrator with the 14 Phase 0 module stubs pre-registered.
func New() *Orchestrator {
	o := &Orchestrator{
		modules: make(map[string]*Module),
	}
	o.registerDefaults()
	return o
}

// RegisterModule registers a module. Returns an error if the ID is already taken.
func (o *Orchestrator) RegisterModule(m Module) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	if _, exists := o.modules[m.ID]; exists {
		return fmt.Errorf("module %q already registered", m.ID)
	}
	o.modules[m.ID] = &m
	return nil
}

// ListModules returns a snapshot of all registered modules.
func (o *Orchestrator) ListModules() []Module {
	o.mu.RLock()
	defer o.mu.RUnlock()

	result := make([]Module, 0, len(o.modules))
	for _, m := range o.modules {
		result = append(result, *m)
	}
	return result
}

// EnableModule enables a module by ID.
func (o *Orchestrator) EnableModule(id string) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	m, ok := o.modules[id]
	if !ok {
		return fmt.Errorf("module %q not found", id)
	}
	m.Enabled = true
	m.Status = "ready"
	return nil
}

// DisableModule disables a module by ID.
func (o *Orchestrator) DisableModule(id string) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	m, ok := o.modules[id]
	if !ok {
		return fmt.Errorf("module %q not found", id)
	}
	m.Enabled = false
	m.Status = "stopped"
	return nil
}

// GetModule returns a single module by ID, or an error if not found.
func (o *Orchestrator) GetModule(id string) (*Module, error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	m, ok := o.modules[id]
	if !ok {
		return nil, fmt.Errorf("module %q not found", id)
	}
	cpy := *m
	return &cpy, nil
}

// registerDefaults pre-registers the 14 Phase 0 module stubs.
func (o *Orchestrator) registerDefaults() {
	defaults := []Module{
		{
			ID:              "nav-gps",
			Name:            "GPS Navigation",
			Domain:          "navigation",
			Description:     "Offline GPS positioning and waypoint tracking",
			Status:          "stopped",
			HardwareMinTier: resourceprofiler.TierMinimum,
		},
		{
			ID:              "nav-compass",
			Name:            "Digital Compass",
			Domain:          "navigation",
			Description:     "Magnetometer-based heading and bearing",
			Status:          "stopped",
			HardwareMinTier: resourceprofiler.TierMinimum,
		},
		{
			ID:              "nav-maps",
			Name:            "Offline Maps",
			Domain:          "navigation",
			Description:     "Tiled offline map rendering with MBTiles support",
			Status:          "stopped",
			HardwareMinTier: resourceprofiler.TierStandard,
		},
		{
			ID:              "env-weather",
			Name:            "Weather Station",
			Domain:          "environment",
			Description:     "Local barometric pressure, temperature, and forecast",
			Status:          "stopped",
			HardwareMinTier: resourceprofiler.TierMinimum,
		},
		{
			ID:              "env-solar",
			Name:            "Solar Tracker",
			Domain:          "environment",
			Description:     "Sun/moon position calculator and golden-hour alerts",
			Status:          "stopped",
			HardwareMinTier: resourceprofiler.TierMinimum,
		},
		{
			ID:              "med-firstaid",
			Name:            "First Aid Guide",
			Domain:          "medical",
			Description:     "Offline first-aid procedures and triage decision trees",
			Status:          "stopped",
			HardwareMinTier: resourceprofiler.TierMinimum,
		},
		{
			ID:              "med-vitals",
			Name:            "Vitals Monitor",
			Domain:          "medical",
			Description:     "Heart rate, SpO2, and temperature logging (sensor-dependent)",
			Status:          "stopped",
			HardwareMinTier: resourceprofiler.TierStandard,
		},
		{
			ID:              "comm-mesh",
			Name:            "Mesh Radio",
			Domain:          "communication",
			Description:     "LoRa/Meshtastic mesh networking integration",
			Status:          "stopped",
			HardwareMinTier: resourceprofiler.TierMinimum,
		},
		{
			ID:              "comm-beacon",
			Name:            "Emergency Beacon",
			Domain:          "communication",
			Description:     "SOS beacon and distress signal broadcaster",
			Status:          "stopped",
			HardwareMinTier: resourceprofiler.TierMinimum,
		},
		{
			ID:              "ref-flora",
			Name:            "Flora Database",
			Domain:          "reference",
			Description:     "Offline plant identification and foraging guide",
			Status:          "stopped",
			HardwareMinTier: resourceprofiler.TierStandard,
		},
		{
			ID:              "ref-fauna",
			Name:            "Fauna Database",
			Domain:          "reference",
			Description:     "Offline animal identification and behavior reference",
			Status:          "stopped",
			HardwareMinTier: resourceprofiler.TierStandard,
		},
		{
			ID:              "ref-knots",
			Name:            "Knot Reference",
			Domain:          "reference",
			Description:     "Animated knot-tying instructions for survival scenarios",
			Status:          "stopped",
			HardwareMinTier: resourceprofiler.TierMinimum,
		},
		{
			ID:              "tools-notepad",
			Name:            "Field Notepad",
			Domain:          "tools",
			Description:     "Markdown notepad with GPS-stamped entries",
			Status:          "stopped",
			HardwareMinTier: resourceprofiler.TierMinimum,
		},
		{
			ID:              "tools-checklist",
			Name:            "Survival Checklist",
			Domain:          "tools",
			Description:     "Configurable gear and preparation checklists",
			Status:          "stopped",
			HardwareMinTier: resourceprofiler.TierMinimum,
		},
	}

	for _, m := range defaults {
		o.modules[m.ID] = func(mod Module) *Module { return &mod }(m)
	}
}
