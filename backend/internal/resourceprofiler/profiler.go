package resourceprofiler

import "fmt"

// HardwareTier classifies the device's hardware capability.
type HardwareTier string

const (
	TierMinimum  HardwareTier = "minimum"  // ≤1 GB RAM
	TierStandard HardwareTier = "standard" // 2–8 GB RAM
	TierOptimal  HardwareTier = "optimal"  // >8 GB RAM
)

// SystemProfile holds detected hardware and OS information.
type SystemProfile struct {
	Tier       HardwareTier `json:"tier"`
	OS         string       `json:"os"`
	Arch       string       `json:"arch"`
	TotalRAMMB uint64       `json:"total_ram_mb"`
	CPUCores   int          `json:"cpu_cores"`
	Hostname   string       `json:"hostname"`
}

// Profiler detects the system's hardware profile.
type Profiler interface {
	DetectProfile() (*SystemProfile, error)
}

// ClassifyTier returns the hardware tier based on total RAM in megabytes.
func ClassifyTier(ramMB uint64) HardwareTier {
	switch {
	case ramMB <= 1024:
		return TierMinimum
	case ramMB <= 8192:
		return TierStandard
	default:
		return TierOptimal
	}
}

// String returns a human-readable description of the profile.
func (p *SystemProfile) String() string {
	return fmt.Sprintf("[%s] %s/%s — %d MB RAM, %d cores (%s)",
		p.Tier, p.OS, p.Arch, p.TotalRAMMB, p.CPUCores, p.Hostname)
}
