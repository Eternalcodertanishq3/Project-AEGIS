//go:build linux

package resourceprofiler

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
)

// LinuxProfiler implements Profiler for Linux.
type LinuxProfiler struct{}

// NewProfiler returns a Profiler for the current platform (Linux).
func NewProfiler() Profiler {
	return &LinuxProfiler{}
}

// DetectProfile detects the system hardware profile on Linux.
func (p *LinuxProfiler) DetectProfile() (*SystemProfile, error) {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	totalRAMMB, err := readMemInfo()
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

// readMemInfo parses /proc/meminfo to extract total RAM in MB.
func readMemInfo() (uint64, error) {
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return 0, fmt.Errorf("opening /proc/meminfo: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "MemTotal:") {
			var totalKB uint64
			_, err := fmt.Sscanf(line, "MemTotal: %d kB", &totalKB)
			if err != nil {
				return 0, fmt.Errorf("parsing MemTotal line %q: %w", line, err)
			}
			return totalKB / 1024, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("reading /proc/meminfo: %w", err)
	}
	return 0, fmt.Errorf("MemTotal not found in /proc/meminfo")
}
