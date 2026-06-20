package orchestrator

import (
	"context"
	"fmt"
	"os/exec"
	"sync"
)

// SidecarState represents the lifecycle state of a sidecar process.
type SidecarState string

const (
	SidecarStopped  SidecarState = "stopped"
	SidecarStarting SidecarState = "starting"
	SidecarRunning  SidecarState = "running"
	SidecarStopping SidecarState = "stopping"
	SidecarFailed   SidecarState = "failed"
)

// Sidecar wraps an exec.Cmd and manages its lifecycle.
// In Phase 0 this is a stub — no real processes are launched.
type Sidecar struct {
	mu       sync.Mutex
	ModuleID string
	Binary   string
	Args     []string
	State    SidecarState
	cmd      *exec.Cmd
	cancel   context.CancelFunc
}

// NewSidecar creates a new Sidecar configuration for a module.
func NewSidecar(moduleID, binary string, args ...string) *Sidecar {
	return &Sidecar{
		ModuleID: moduleID,
		Binary:   binary,
		Args:     args,
		State:    SidecarStopped,
	}
}

// Start launches the sidecar process.
// Phase 0 stub: sets state to Running without actually launching a process.
func (s *Sidecar) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.State == SidecarRunning {
		return fmt.Errorf("sidecar %q is already running", s.ModuleID)
	}

	// Phase 0: We don't actually start a process yet.
	// When real binaries are ready, this will use exec.CommandContext.
	s.State = SidecarRunning
	return nil
}

// Stop terminates the sidecar process.
// Phase 0 stub: sets state to Stopped.
func (s *Sidecar) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.State != SidecarRunning {
		return fmt.Errorf("sidecar %q is not running (state: %s)", s.ModuleID, s.State)
	}

	if s.cancel != nil {
		s.cancel()
		s.cancel = nil
	}

	s.State = SidecarStopped
	return nil
}

// HealthCheck returns nil if the sidecar is healthy.
// Phase 0 stub: always returns healthy if running.
func (s *Sidecar) HealthCheck() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.State != SidecarRunning {
		return fmt.Errorf("sidecar %q is not running (state: %s)", s.ModuleID, s.State)
	}

	// Phase 0: assume healthy if running.
	return nil
}

// GetState returns the current sidecar state.
func (s *Sidecar) GetState() SidecarState {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.State
}
