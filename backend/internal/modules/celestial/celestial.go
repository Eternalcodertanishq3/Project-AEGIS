package celestial

import (
	"fmt"
	"math"
	"time"
)

// CelestialResult holds the computed position data.
type CelestialResult struct {
	DateTime      string  `json:"date_time"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	SunAzimuth    float64 `json:"sun_azimuth"`
	SunAltitude   float64 `json:"sun_altitude"`
	SunRise       string  `json:"sun_rise"`
	SunSet        string  `json:"sun_set"`
	MoonPhase     string  `json:"moon_phase"`
	MoonPhaseAngle float64 `json:"moon_phase_angle"`
	PolarisAlt    float64 `json:"polaris_altitude"`
	TrueNorth     string  `json:"true_north_hint"`
	DayLength     string  `json:"day_length"`
}

// StarInfo describes a navigation star.
type StarInfo struct {
	Name          string  `json:"name"`
	Constellation string  `json:"constellation"`
	Description   string  `json:"description"`
	Magnitude     float64 `json:"magnitude"`
	Usage         string  `json:"usage"`
}

// GetNavigationStars returns the key stars used for celestial navigation.
func GetNavigationStars() []StarInfo {
	return []StarInfo{
		{Name: "Polaris", Constellation: "Ursa Minor", Magnitude: 1.98,
			Description: "The North Star. Located at the tip of the Little Dipper's handle.",
			Usage: "Polaris sits within 1° of true north. Its altitude above the horizon equals your latitude."},
		{Name: "Sigma Octantis", Constellation: "Octans", Magnitude: 5.47,
			Description: "The Southern Hemisphere's pole star. Very faint and difficult to see.",
			Usage: "In the Southern Hemisphere, use the Southern Cross to find the South Celestial Pole instead."},
		{Name: "Sirius", Constellation: "Canis Major", Magnitude: -1.46,
			Description: "The brightest star in the night sky. Brilliant blue-white.",
			Usage: "Rises nearly due east and sets nearly due west. Visible from almost everywhere on Earth."},
		{Name: "Vega", Constellation: "Lyra", Magnitude: 0.03,
			Description: "Bright blue-white star, part of the Summer Triangle.",
			Usage: "Rises in the northeast, sets in the northwest. Dominant summer star in Northern Hemisphere."},
		{Name: "Canopus", Constellation: "Carina", Magnitude: -0.74,
			Description: "Second brightest star. Visible mostly from Southern Hemisphere.",
			Usage: "Used as reference for southern navigation when the Southern Cross is not visible."},
		{Name: "Betelgeuse", Constellation: "Orion", Magnitude: 0.42,
			Description: "Red supergiant forming Orion's left shoulder.",
			Usage: "Orion's belt rises due east and sets due west — one of the most reliable direction indicators."},
		{Name: "Rigel", Constellation: "Orion", Magnitude: 0.13,
			Description: "Blue supergiant at Orion's right knee.",
			Usage: "Along with Betelgeuse, makes Orion easy to identify for east-west orientation."},
	}
}

// NavigationTechnique describes a method for celestial navigation.
type NavigationTechnique struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Hemisphere  string   `json:"hemisphere"` // north, south, both
	Description string   `json:"description"`
	Steps       []string `json:"steps"`
}

// GetTechniques returns celestial navigation techniques.
func GetTechniques() []NavigationTechnique {
	return []NavigationTechnique{
		{
			ID: "polaris-latitude", Name: "Latitude from Polaris", Hemisphere: "north",
			Description: "Determine your latitude by measuring Polaris's altitude above the horizon.",
			Steps: []string{
				"Locate the Big Dipper (Ursa Major). Find the two 'pointer stars' on the front of the cup.",
				"Extend a line through the pointer stars about 5× the distance between them. This leads to Polaris.",
				"Measure Polaris's angle above the horizon. Use a fist at arm's length ≈ 10°.",
				"That angle IS your latitude. Polaris at 45° altitude = you're at 45°N latitude.",
			},
		},
		{
			ID: "southern-cross", Name: "South from Southern Cross", Hemisphere: "south",
			Description: "Find true south using the Southern Cross (Crux) constellation.",
			Steps: []string{
				"Identify the Southern Cross — four bright stars forming a cross shape.",
				"Extend the long axis of the cross 4.5× its length toward the horizon.",
				"That point is approximately the South Celestial Pole.",
				"Drop a vertical line from that point to the horizon — that's due south.",
			},
		},
		{
			ID: "orion-east-west", Name: "East-West from Orion's Belt", Hemisphere: "both",
			Description: "Use Orion's belt to find east and west directions.",
			Steps: []string{
				"Identify Orion by his three belt stars in a straight line.",
				"Orion's belt rises almost exactly due east and sets almost exactly due west.",
				"When the belt is rising (tilted), it points toward the east.",
				"When the belt is setting (tilted the other way), it points toward the west.",
			},
		},
		{
			ID: "crescent-moon", Name: "Direction from Crescent Moon", Hemisphere: "both",
			Description: "Use the crescent moon to find approximate south (or north).",
			Steps: []string{
				"Find a crescent moon in the sky.",
				"Draw an imaginary line connecting the two tips (horns) of the crescent.",
				"Extend that line down to the horizon.",
				"In the Northern Hemisphere, this point is approximately south.",
				"In the Southern Hemisphere, this point is approximately north.",
			},
		},
	}
}

// Calculate computes sun position, moon phase, and navigation data for a given time and location.
func Calculate(lat, lon float64, t time.Time) *CelestialResult {
	// Julian Day Number
	jd := julianDay(t)
	// Sun position
	azimuth, altitude := sunPosition(lat, lon, jd, t)
	// Sunrise / Sunset
	rise, set, dayLen := sunRiseSet(lat, lon, t)
	// Moon phase
	phase, phaseAngle := moonPhase(jd)
	// Polaris altitude (roughly equals latitude for Northern Hemisphere)
	polarisAlt := lat
	if lat < 0 {
		polarisAlt = 0 // Not visible from Southern Hemisphere
	}

	northHint := "Face Polaris (North Star) for true north"
	if lat < 0 {
		northHint = "Use the Southern Cross to find south, then face opposite for north"
	}

	return &CelestialResult{
		DateTime:       t.Format(time.RFC3339),
		Latitude:       lat,
		Longitude:      lon,
		SunAzimuth:     math.Round(azimuth*100) / 100,
		SunAltitude:    math.Round(altitude*100) / 100,
		SunRise:        rise,
		SunSet:         set,
		MoonPhase:      phase,
		MoonPhaseAngle: math.Round(phaseAngle*100) / 100,
		PolarisAlt:     math.Round(polarisAlt*100) / 100,
		TrueNorth:      northHint,
		DayLength:      dayLen,
	}
}

func julianDay(t time.Time) float64 {
	y := float64(t.Year())
	m := float64(t.Month())
	d := float64(t.Day()) + float64(t.Hour())/24.0 + float64(t.Minute())/1440.0
	if m <= 2 {
		y--
		m += 12
	}
	A := math.Floor(y / 100)
	B := 2 - A + math.Floor(A/4)
	return math.Floor(365.25*(y+4716)) + math.Floor(30.6001*(m+1)) + d + B - 1524.5
}

func sunPosition(lat, lon, jd float64, t time.Time) (azimuth, altitude float64) {
	// Days since J2000.0
	n := jd - 2451545.0
	// Mean longitude & anomaly
	L := math.Mod(280.46+0.9856474*n, 360)
	g := math.Mod(357.528+0.9856003*n, 360) * math.Pi / 180
	// Ecliptic longitude
	lambda := (L + 1.915*math.Sin(g) + 0.020*math.Sin(2*g)) * math.Pi / 180
	// Obliquity
	epsilon := 23.439 * math.Pi / 180
	// Right ascension & declination
	sinRA := math.Cos(epsilon) * math.Sin(lambda)
	cosRA := math.Cos(lambda)
	ra := math.Atan2(sinRA, cosRA)
	dec := math.Asin(math.Sin(epsilon) * math.Sin(lambda))
	// Hour angle
	gmst := math.Mod(280.46061837+360.98564736629*n, 360)
	lmst := gmst + lon
	ha := (lmst*math.Pi/180 - ra)

	latRad := lat * math.Pi / 180
	sinAlt := math.Sin(latRad)*math.Sin(dec) + math.Cos(latRad)*math.Cos(dec)*math.Cos(ha)
	altitude = math.Asin(sinAlt) * 180 / math.Pi

	cosAz := (math.Sin(dec) - math.Sin(latRad)*sinAlt) / (math.Cos(latRad) * math.Cos(math.Asin(sinAlt)))
	cosAz = math.Max(-1, math.Min(1, cosAz))
	azimuth = math.Acos(cosAz) * 180 / math.Pi
	if math.Sin(ha) > 0 {
		azimuth = 360 - azimuth
	}
	return
}

func sunRiseSet(lat, lon float64, t time.Time) (rise, set, dayLength string) {
	jd := julianDay(time.Date(t.Year(), t.Month(), t.Day(), 12, 0, 0, 0, time.UTC))
	n := jd - 2451545.0
	g := math.Mod(357.528+0.9856003*n, 360) * math.Pi / 180
	L := math.Mod(280.46+0.9856474*n, 360)
	lambda := (L + 1.915*math.Sin(g) + 0.020*math.Sin(2*g)) * math.Pi / 180
	epsilon := 23.439 * math.Pi / 180
	dec := math.Asin(math.Sin(epsilon) * math.Sin(lambda))

	latRad := lat * math.Pi / 180
	cosH := (math.Sin(-0.8333*math.Pi/180) - math.Sin(latRad)*math.Sin(dec)) / (math.Cos(latRad) * math.Cos(dec))

	if cosH > 1 {
		return "none (polar night)", "none (polar night)", "0h 0m"
	}
	if cosH < -1 {
		return "none (midnight sun)", "none (midnight sun)", "24h 0m"
	}

	H := math.Acos(cosH) * 180 / math.Pi
	dayMinutes := 2 * H / 15 * 60

	// Approximate sunrise/sunset in UTC
	gmst := math.Mod(280.46061837+360.98564736629*n, 360)
	sinRA := math.Cos(epsilon) * math.Sin(lambda)
	cosRA := math.Cos(lambda)
	ra := math.Atan2(sinRA, cosRA) * 180 / math.Pi

	transit := (ra - gmst - lon) / 360 * 24
	for transit < 0 {
		transit += 24
	}
	for transit > 24 {
		transit -= 24
	}
	riseHours := transit - H/15
	setHours := transit + H/15

	riseH := int(riseHours)
	riseM := int((riseHours - float64(riseH)) * 60)
	if riseH < 0 { riseH += 24 }
	setH := int(setHours)
	setM := int((setHours - float64(setH)) * 60)
	if setH >= 24 { setH -= 24 }

	dh := int(dayMinutes) / 60
	dm := int(dayMinutes) % 60

	rise = formatTime(riseH, riseM)
	set = formatTime(setH, setM)
	dayLength = formatDuration(dh, dm)
	return
}

func formatTime(h, m int) string {
	if m < 0 { m = 0 }
	if m > 59 { m = 59 }
	ampm := "AM"
	h12 := h
	if h12 >= 12 { ampm = "PM" }
	if h12 > 12 { h12 -= 12 }
	if h12 == 0 { h12 = 12 }
	return fmt.Sprintf("%d:%02d %s UTC", h12, m, ampm)
}

func formatDuration(h, m int) string {
	return fmt.Sprintf("%dh %dm", h, m)
}

func moonPhase(jd float64) (string, float64) {
	// Days since known new moon (Jan 6, 2000 18:14 UTC)
	daysSinceNew := jd - 2451550.1
	synodic := 29.53058867
	phase := math.Mod(daysSinceNew, synodic)
	if phase < 0 {
		phase += synodic
	}
	angle := phase / synodic * 360

	var name string
	switch {
	case phase < 1.85:
		name = "New Moon"
	case phase < 7.38:
		name = "Waxing Crescent"
	case phase < 9.23:
		name = "First Quarter"
	case phase < 14.77:
		name = "Waxing Gibbous"
	case phase < 16.61:
		name = "Full Moon"
	case phase < 22.15:
		name = "Waning Gibbous"
	case phase < 23.99:
		name = "Last Quarter"
	default:
		name = "Waning Crescent"
	}
	return name, angle
}
