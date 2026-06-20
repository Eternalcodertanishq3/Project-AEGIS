package powermanager

// PowerStatus holds the current power state of the device.
type PowerStatus struct {
	OnBattery      bool   `json:"on_battery"`
	BatteryPercent int    `json:"battery_percent"`
	ACConnected    bool   `json:"ac_connected"`
	Status         string `json:"status"`
}

// PowerManager detects the device's current power state.
type PowerManager interface {
	GetPowerStatus() (*PowerStatus, error)
}
