//go:build linux

package powermanager

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const powerSupplyPath = "/sys/class/power_supply"

// LinuxPowerManager implements PowerManager for Linux.
type LinuxPowerManager struct{}

// NewPowerManager returns a PowerManager for the current platform (Linux).
func NewPowerManager() PowerManager {
	return &LinuxPowerManager{}
}

// GetPowerStatus reads power information from /sys/class/power_supply/.
func (m *LinuxPowerManager) GetPowerStatus() (*PowerStatus, error) {
	entries, err := os.ReadDir(powerSupplyPath)
	if err != nil {
		// No power_supply directory likely means a desktop without battery.
		return &PowerStatus{
			OnBattery:      false,
			BatteryPercent: -1,
			ACConnected:    true,
			Status:         "ac_power",
		}, nil
	}

	ps := &PowerStatus{
		BatteryPercent: -1,
		ACConnected:    true,
		Status:         "ac_power",
	}

	for _, entry := range entries {
		supplyType := readSysFile(filepath.Join(powerSupplyPath, entry.Name(), "type"))

		switch strings.TrimSpace(supplyType) {
		case "Battery":
			capStr := readSysFile(filepath.Join(powerSupplyPath, entry.Name(), "capacity"))
			if cap, err := strconv.Atoi(strings.TrimSpace(capStr)); err == nil {
				ps.BatteryPercent = cap
			}
			statusStr := readSysFile(filepath.Join(powerSupplyPath, entry.Name(), "status"))
			statusStr = strings.TrimSpace(statusStr)
			if statusStr == "Discharging" {
				ps.OnBattery = true
				ps.ACConnected = false
				ps.Status = "battery"
			}
		case "Mains":
			onlineStr := readSysFile(filepath.Join(powerSupplyPath, entry.Name(), "online"))
			if strings.TrimSpace(onlineStr) == "0" {
				ps.ACConnected = false
			}
		}
	}

	return ps, nil
}

// readSysFile reads a sysfs file, returning empty string on any error.
func readSysFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(data)
}

// ensure LinuxPowerManager satisfies the interface at compile time.
var _ PowerManager = (*LinuxPowerManager)(nil)

// Provide a descriptive error wrapper helper.
func wrapErr(context string, err error) error {
	return fmt.Errorf("%s: %w", context, err)
}
