//go:build darwin

package powermanager

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// DarwinPowerManager implements PowerManager for macOS.
type DarwinPowerManager struct{}

// NewPowerManager returns a PowerManager for the current platform (macOS).
func NewPowerManager() PowerManager {
	return &DarwinPowerManager{}
}

// GetPowerStatus parses `pmset -g batt` output on macOS.
//
// Example output:
//
//	Now drawing from 'AC Power'
//	 -InternalBattery-0 (id=...)	85%; charged; 0:00 remaining
func (m *DarwinPowerManager) GetPowerStatus() (*PowerStatus, error) {
	out, err := exec.Command("pmset", "-g", "batt").Output()
	if err != nil {
		return nil, fmt.Errorf("running pmset: %w", err)
	}

	output := string(out)
	lines := strings.Split(output, "\n")

	ps := &PowerStatus{
		BatteryPercent: -1,
		ACConnected:    true,
		Status:         "ac_power",
	}

	for _, line := range lines {
		if strings.Contains(line, "'Battery Power'") {
			ps.OnBattery = true
			ps.ACConnected = false
			ps.Status = "battery"
		}
		// Look for battery percentage, e.g. "85%;"
		if strings.Contains(line, "InternalBattery") {
			parts := strings.Fields(line)
			for _, part := range parts {
				if strings.HasSuffix(part, "%;") {
					pctStr := strings.TrimSuffix(part, "%;")
					if pct, err := strconv.Atoi(pctStr); err == nil {
						ps.BatteryPercent = pct
					}
				}
			}
		}
	}

	return ps, nil
}
