//go:build windows

package powermanager

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

// systemPowerStatus mirrors the Windows SYSTEM_POWER_STATUS structure.
type systemPowerStatus struct {
	ACLineStatus        byte
	BatteryFlag         byte
	BatteryLifePercent  byte
	SystemStatusFlag    byte
	BatteryLifeTime     uint32
	BatteryFullLifeTime uint32
}

// WindowsPowerManager implements PowerManager for Windows.
type WindowsPowerManager struct{}

// NewPowerManager returns a PowerManager for the current platform (Windows).
func NewPowerManager() PowerManager {
	return &WindowsPowerManager{}
}

// GetPowerStatus retrieves the current power state on Windows via GetSystemPowerStatus.
func (m *WindowsPowerManager) GetPowerStatus() (*PowerStatus, error) {
	kernel32 := windows.NewLazySystemDLL("kernel32.dll")
	proc := kernel32.NewProc("GetSystemPowerStatus")

	var sps systemPowerStatus
	ret, _, err := proc.Call(uintptr(unsafe.Pointer(&sps)))
	if ret == 0 {
		return nil, fmt.Errorf("GetSystemPowerStatus failed: %w", err)
	}

	acConnected := sps.ACLineStatus == 1
	batteryPercent := int(sps.BatteryLifePercent)
	if batteryPercent > 100 {
		// 255 means unknown
		batteryPercent = -1
	}

	status := "ac_power"
	onBattery := false
	if !acConnected {
		onBattery = true
		status = "battery"
	}

	return &PowerStatus{
		OnBattery:      onBattery,
		BatteryPercent: batteryPercent,
		ACConnected:    acConnected,
		Status:         status,
	}, nil
}
