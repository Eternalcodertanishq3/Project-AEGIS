//go:build darwin

package resourceprofiler

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

// DarwinProfiler implements Profiler for macOS.
type DarwinProfiler struct{}

// NewProfiler returns a Profiler for the current platform (macOS).
func NewProfiler() Profiler {
	return &DarwinProfiler{}
}

// DetectProfile detects the system hardware profile on macOS.
func (p *DarwinProfiler) DetectProfile() (*SystemProfile, error) {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	totalRAMMB, err := darwinRAM()
	if err != nil {
		return nil, fmt.Errorf("detecting RAM: %w", err)
	}

	profile := &SystemProfile{
		Tier:       ClassifyTier(totalRAMMB),
		OS:         runtime.GOOS,
		Arch:       runtime.GOARCH,
		TotalRAMMB: totalRAMMB,
		CPUCores:   runtime.NumCPU(),
		Hostname:   hostname,
	}
	return profile, nil
}

// darwinRAM uses sysctl hw.memsize to retrieve total physical RAM in MB.
func darwinRAM() (uint64, error) {
	out, err := exec.Command("sysctl", "-n", "hw.memsize").Output()
	if err != nil {
		return 0, fmt.Errorf("running sysctl hw.memsize: %w", err)
	}
	bytes, err := strconv.ParseUint(strings.TrimSpace(string(out)), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parsing sysctl output %q: %w", string(out), err)
	}
	return bytes / (1024 * 1024), nil
}
