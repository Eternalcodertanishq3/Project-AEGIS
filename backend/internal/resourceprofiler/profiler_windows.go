//go:build windows

package resourceprofiler

import (
	"fmt"
	"os"
	"runtime"
	"unsafe"

	"golang.org/x/sys/windows"
)

// memoryStatusEx corresponds to the Windows MEMORYSTATUSEX structure.
type memoryStatusEx struct {
	Length               uint32
	MemoryLoad           uint32
	TotalPhys            uint64
	AvailPhys            uint64
	TotalPageFile        uint64
	AvailPageFile        uint64
	TotalVirtual         uint64
	AvailVirtual         uint64
	AvailExtendedVirtual uint64
}

// WindowsProfiler implements Profiler for Windows.
type WindowsProfiler struct{}

// NewProfiler returns a Profiler for the current platform (Windows).
func NewProfiler() Profiler {
	return &WindowsProfiler{}
}

// DetectProfile detects the system hardware profile on Windows.
func (p *WindowsProfiler) DetectProfile() (*SystemProfile, error) {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	totalRAMMB, err := getWindowsRAM()
	if err != nil {
		// Fallback: estimate from Go's view of memory (unreliable but better than zero).
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		totalRAMMB = m.Sys / (1024 * 1024)
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

// getWindowsRAM calls GlobalMemoryStatusEx to retrieve total physical RAM.
func getWindowsRAM() (uint64, error) {
	kernel32 := windows.NewLazySystemDLL("kernel32.dll")
	globalMemoryStatusEx := kernel32.NewProc("GlobalMemoryStatusEx")

	var memStatus memoryStatusEx
	memStatus.Length = uint32(unsafe.Sizeof(memStatus))

	ret, _, err := globalMemoryStatusEx.Call(uintptr(unsafe.Pointer(&memStatus)))
	if ret == 0 {
		return 0, fmt.Errorf("GlobalMemoryStatusEx failed: %w", err)
	}
	return memStatus.TotalPhys / (1024 * 1024), nil
}
