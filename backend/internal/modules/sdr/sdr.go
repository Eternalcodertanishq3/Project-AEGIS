package sdr

// FrequencyEntry represents a single radio frequency entry.
type FrequencyEntry struct {
	ID          string  `json:"id"`
	Frequency   string  `json:"frequency"`
	FreqMHz     float64 `json:"freq_mhz"`
	Name        string  `json:"name"`
	Mode        string  `json:"mode"` // AM, FM, NFM, USB, LSB, CW, Digital
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Priority    string  `json:"priority,omitempty"` // critical, important, useful, reference
}

// BandPlan represents a radio frequency band allocation.
type BandPlan struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	StartMHz    float64 `json:"start_mhz"`
	EndMHz      float64 `json:"end_mhz"`
	Allocation  string `json:"allocation"`
	Description string `json:"description"`
}

// SDRStatus describes the current SDR hardware status.
type SDRStatus struct {
	HardwareDetected bool   `json:"hardware_detected"`
	DeviceType       string `json:"device_type"`
	Mode             string `json:"mode"` // "reference-only", "monitoring", "scanning"
	FrequencyCount   int    `json:"frequency_count"`
	BandPlanCount    int    `json:"band_plan_count"`
}

// FrequencyGroup groups frequencies by category.
type FrequencyGroup struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Icon        string           `json:"icon"`
	Description string           `json:"description"`
	Frequencies []FrequencyEntry `json:"frequencies"`
}

// SDRDatabase holds the offline frequency reference data.
type SDRDatabase struct {
	Groups    []FrequencyGroup
	BandPlans []BandPlan
}

// NewSDRDatabase creates the frequency reference database.
func NewSDRDatabase() *SDRDatabase {
	return &SDRDatabase{
		Groups: []FrequencyGroup{
			{
				ID: "emergency", Name: "Emergency Frequencies", Icon: "alert-triangle",
				Description: "Critical emergency and distress frequencies — memorize these",
				Frequencies: []FrequencyEntry{
					{ID: "vhf-distress", Frequency: "156.800 MHz", FreqMHz: 156.800, Name: "VHF Channel 16 — Maritime Distress", Mode: "FM", Description: "International maritime distress, safety, and calling frequency. Monitored by coast guards worldwide.", Category: "emergency", Priority: "critical"},
					{ID: "guard-freq", Frequency: "121.500 MHz", FreqMHz: 121.500, Name: "VHF Guard — Aviation Emergency", Mode: "AM", Description: "International aeronautical emergency frequency. All aircraft monitor this. ELT beacons transmit here.", Category: "emergency", Priority: "critical"},
					{ID: "uhf-guard", Frequency: "243.000 MHz", FreqMHz: 243.000, Name: "UHF Guard — Military Emergency", Mode: "AM", Description: "Military emergency frequency. 2nd harmonic of 121.5 MHz.", Category: "emergency", Priority: "critical"},
					{ID: "cb-ch9", Frequency: "27.065 MHz", FreqMHz: 27.065, Name: "CB Channel 9 — Emergency", Mode: "AM", Description: "Citizens Band emergency channel. Monitored by REACT teams and truckers in the US.", Category: "emergency", Priority: "critical"},
					{ID: "cb-ch19", Frequency: "27.185 MHz", FreqMHz: 27.185, Name: "CB Channel 19 — Highway", Mode: "AM", Description: "Most popular CB channel. Truckers and travelers. De facto calling channel on highways.", Category: "emergency", Priority: "important"},
					{ID: "frs-ch1", Frequency: "462.5625 MHz", FreqMHz: 462.5625, Name: "FRS Channel 1 — Family Radio", Mode: "NFM", Description: "First FRS channel — often used as a default calling channel on consumer walkie-talkies.", Category: "emergency", Priority: "important"},
					{ID: "gmrs-calling", Frequency: "462.5625 MHz", FreqMHz: 462.5625, Name: "GMRS Channel 1", Mode: "NFM", Description: "GMRS calling channel. Higher power than FRS. License required in US but not enforced in emergencies.", Category: "emergency", Priority: "important"},
				},
			},
			{
				ID: "weather", Name: "Weather Stations", Icon: "cloud",
				Description: "NOAA Weather Radio and weather broadcasting frequencies",
				Frequencies: []FrequencyEntry{
					{ID: "noaa-1", Frequency: "162.400 MHz", FreqMHz: 162.400, Name: "NOAA Weather 1 (WX1)", Mode: "NFM", Description: "NOAA Weather Radio — continuous weather broadcasts. Most widely used WX frequency.", Category: "weather", Priority: "important"},
					{ID: "noaa-2", Frequency: "162.425 MHz", FreqMHz: 162.425, Name: "NOAA Weather 2 (WX2)", Mode: "NFM", Description: "NOAA Weather Radio — alternate frequency.", Category: "weather"},
					{ID: "noaa-3", Frequency: "162.450 MHz", FreqMHz: 162.450, Name: "NOAA Weather 3 (WX3)", Mode: "NFM", Description: "NOAA Weather Radio — alternate frequency.", Category: "weather"},
					{ID: "noaa-4", Frequency: "162.475 MHz", FreqMHz: 162.475, Name: "NOAA Weather 4 (WX4)", Mode: "NFM", Description: "NOAA Weather Radio — alternate frequency.", Category: "weather"},
					{ID: "noaa-5", Frequency: "162.500 MHz", FreqMHz: 162.500, Name: "NOAA Weather 5 (WX5)", Mode: "NFM", Description: "NOAA Weather Radio — alternate frequency.", Category: "weather"},
					{ID: "noaa-6", Frequency: "162.525 MHz", FreqMHz: 162.525, Name: "NOAA Weather 6 (WX6)", Mode: "NFM", Description: "NOAA Weather Radio — alternate frequency.", Category: "weather"},
					{ID: "noaa-7", Frequency: "162.550 MHz", FreqMHz: 162.550, Name: "NOAA Weather 7 (WX7)", Mode: "NFM", Description: "NOAA Weather Radio — alternate frequency.", Category: "weather"},
				},
			},
			{
				ID: "aviation", Name: "Aviation", Icon: "plane",
				Description: "Aircraft communication and navigation frequencies",
				Frequencies: []FrequencyEntry{
					{ID: "atis", Frequency: "118.000-136.975 MHz", FreqMHz: 118.0, Name: "VHF Airband (AM)", Mode: "AM", Description: "Entire VHF aviation communication band. Includes tower, ground, approach, departure, center, and ATIS.", Category: "aviation"},
					{ID: "unicom", Frequency: "122.800 MHz", FreqMHz: 122.800, Name: "UNICOM — Uncontrolled Airports", Mode: "AM", Description: "Common frequency at airports without control towers. Pilots self-announce positions.", Category: "aviation"},
					{ID: "multicom", Frequency: "122.900 MHz", FreqMHz: 122.900, Name: "MULTICOM", Mode: "AM", Description: "Self-announce frequency for airports with no UNICOM or CTAF. Also used for aerial surveys.", Category: "aviation"},
					{ID: "sar-air", Frequency: "123.100 MHz", FreqMHz: 123.100, Name: "SAR — Search and Rescue", Mode: "AM", Description: "Search and rescue operations, civil aviation.", Category: "aviation", Priority: "important"},
				},
			},
			{
				ID: "marine", Name: "Marine VHF", Icon: "anchor",
				Description: "Maritime communication channels",
				Frequencies: []FrequencyEntry{
					{ID: "marine-ch16", Frequency: "156.800 MHz", FreqMHz: 156.800, Name: "Channel 16 — Distress/Calling", Mode: "FM", Description: "International distress and calling channel. Always monitor this at sea.", Category: "marine", Priority: "critical"},
					{ID: "marine-ch06", Frequency: "156.300 MHz", FreqMHz: 156.300, Name: "Channel 6 — Ship to Ship Safety", Mode: "FM", Description: "Intership safety communications.", Category: "marine", Priority: "important"},
					{ID: "marine-ch13", Frequency: "156.650 MHz", FreqMHz: 156.650, Name: "Channel 13 — Bridge to Bridge", Mode: "FM", Description: "Navigation safety. Used by ships in close proximity.", Category: "marine"},
					{ID: "marine-ch22a", Frequency: "157.100 MHz", FreqMHz: 157.100, Name: "Channel 22A — Coast Guard Broadcast", Mode: "FM", Description: "US Coast Guard liaison and marine safety broadcasts.", Category: "marine", Priority: "important"},
				},
			},
			{
				ID: "amateur", Name: "Amateur Radio (Ham)", Icon: "radio",
				Description: "Ham radio frequencies — license required to transmit",
				Frequencies: []FrequencyEntry{
					{ID: "2m-calling", Frequency: "146.520 MHz", FreqMHz: 146.520, Name: "2m Simplex Calling", Mode: "FM", Description: "National 2-meter simplex calling frequency. Best first frequency to try for local ham contacts.", Category: "amateur", Priority: "important"},
					{ID: "70cm-calling", Frequency: "446.000 MHz", FreqMHz: 446.000, Name: "70cm Simplex Calling", Mode: "FM", Description: "National 70-centimeter simplex calling frequency.", Category: "amateur"},
					{ID: "2m-emer", Frequency: "146.550 MHz", FreqMHz: 146.550, Name: "2m Emergency Simplex", Mode: "FM", Description: "Secondary emergency simplex frequency. Used by ARES/RACES groups.", Category: "amateur", Priority: "important"},
					{ID: "hf-20m", Frequency: "14.300 MHz", FreqMHz: 14.300, Name: "20m Emergency Net", Mode: "USB", Description: "Maritime Mobile Service Net and Intercontinental Traffic Net. HF long-distance emergency traffic.", Category: "amateur", Priority: "important"},
					{ID: "hf-40m", Frequency: "7.290 MHz", FreqMHz: 7.290, Name: "40m Traffic Net", Mode: "LSB", Description: "East Coast traffic net and regional emergency communications. Best at night.", Category: "amateur"},
				},
			},
			{
				ID: "lora-ism", Name: "LoRa / ISM Bands", Icon: "wifi",
				Description: "License-free IoT and mesh networking frequencies",
				Frequencies: []FrequencyEntry{
					{ID: "lora-us", Frequency: "902-928 MHz", FreqMHz: 915.0, Name: "LoRa US (ISM 915)", Mode: "Digital", Description: "US/Australia ISM band used by Meshtastic, LoRa, and other IoT devices.", Category: "lora"},
					{ID: "lora-eu", Frequency: "863-870 MHz", FreqMHz: 868.0, Name: "LoRa EU (ISM 868)", Mode: "Digital", Description: "EU ISM band used by LoRa and Meshtastic devices in Europe.", Category: "lora"},
					{ID: "meshtastic-default", Frequency: "906.875 MHz", FreqMHz: 906.875, Name: "Meshtastic Default (US)", Mode: "Digital", Description: "Default Meshtastic long-range channel preset. LoRa spread factor 11.", Category: "lora", Priority: "important"},
				},
			},
		},
		BandPlans: []BandPlan{
			{ID: "lf", Name: "Low Frequency (LF)", StartMHz: 0.03, EndMHz: 0.3, Allocation: "Navigation, Time Signals", Description: "LORAN-C navigation, time signal stations (WWVB at 60 kHz)."},
			{ID: "mf", Name: "Medium Frequency (MF)", StartMHz: 0.3, EndMHz: 3.0, Allocation: "AM Broadcast, Maritime", Description: "AM radio (530-1700 kHz), maritime coastal stations, 160m ham band."},
			{ID: "hf", Name: "High Frequency (HF)", StartMHz: 3.0, EndMHz: 30.0, Allocation: "Shortwave, Ham, Military", Description: "Shortwave broadcast, amateur radio (80m-10m), military, aviation HF. Global propagation via ionosphere."},
			{ID: "vhf", Name: "Very High Frequency (VHF)", StartMHz: 30.0, EndMHz: 300.0, Allocation: "FM, TV, Aviation, Marine, Ham", Description: "FM radio (88-108 MHz), aviation (118-137 MHz), marine VHF (156-163 MHz), 2m ham (144-148 MHz), NOAA weather."},
			{ID: "uhf", Name: "Ultra High Frequency (UHF)", StartMHz: 300.0, EndMHz: 3000.0, Allocation: "TV, Cellular, GMRS, Ham, LoRa", Description: "UHF TV, 70cm ham (420-450 MHz), FRS/GMRS (462-467 MHz), LoRa/ISM (902-928 MHz), cellular."},
		},
	}
}

// GetGroups returns all frequency groups.
func (db *SDRDatabase) GetGroups() []FrequencyGroup {
	return db.Groups
}

// GetGroup returns a single frequency group.
func (db *SDRDatabase) GetGroup(id string) *FrequencyGroup {
	for _, g := range db.Groups {
		if g.ID == id {
			return &g
		}
	}
	return nil
}

// GetBandPlans returns all band plans.
func (db *SDRDatabase) GetBandPlans() []BandPlan {
	return db.BandPlans
}

// SearchFrequencies searches frequencies by name or description.
func (db *SDRDatabase) SearchFrequencies(query string) []FrequencyEntry {
	var results []FrequencyEntry
	q := sdrToLower(query)
	for _, g := range db.Groups {
		for _, f := range g.Frequencies {
			if sdrContains(sdrToLower(f.Name), q) ||
				sdrContains(sdrToLower(f.Description), q) ||
				sdrContains(sdrToLower(f.Frequency), q) ||
				sdrContains(sdrToLower(f.Category), q) {
				results = append(results, f)
			}
		}
	}
	if results == nil {
		results = make([]FrequencyEntry, 0)
	}
	return results
}

// GetStatus returns SDR module status.
func (db *SDRDatabase) GetStatus() *SDRStatus {
	total := 0
	for _, g := range db.Groups {
		total += len(g.Frequencies)
	}
	return &SDRStatus{
		HardwareDetected: false,
		DeviceType:       "none",
		Mode:             "reference-only",
		FrequencyCount:   total,
		BandPlanCount:    len(db.BandPlans),
	}
}

func sdrToLower(s string) string {
	b := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		b[i] = c
	}
	return string(b)
}

func sdrContains(s, sub string) bool {
	if len(sub) == 0 {
		return true
	}
	if len(sub) > len(s) {
		return false
	}
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
