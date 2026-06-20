# Platform Notes — Project AEGIS

> This document catalogs every known platform-specific caveat, driver requirement, and workaround. **Nothing is glossed over.** If a module requires additional setup on a given OS, it is documented here.

---

## Table of Contents

- [Windows](#windows)
- [macOS](#macos)
- [Linux (Desktop)](#linux-desktop)
- [Raspberry Pi / ARM64](#raspberry-pi--arm64)
- [General Notes](#general-notes)

---

## Windows

### Supported Versions

- Windows 10 (version 1809+) and Windows 11
- x86_64 (amd64) architecture

### RTL-SDR Requires Zadig WinUSB Driver Swap

> ⚠️ **This is the single most common setup issue on Windows.**

Windows does not natively support the RTL2832U chipset used by SDR dongles. The default Windows driver will claim the device, preventing `rtl_fm` / `rtl_power` from accessing it.

**Required one-time setup:**

1. Download [Zadig](https://zadig.akeo.ie/) (portable, no install required)
2. Plug in the RTL-SDR dongle
3. Open Zadig → **Options** → **List All Devices**
4. Select **Bulk-In, Interface (Interface 0)** from the dropdown
   - If you see "RTL2838UHIDIR" or similar, that's the correct device
5. Set the target driver to **WinUSB (v6.x.x.x)**
6. Click **Replace Driver**
7. Wait for "Driver installation successful"

**Reverting:** If you need to use the dongle with other software (e.g., SDR#), use Zadig to switch back to the original driver.

**AEGIS behavior when driver is not swapped:** The SDR Monitor module will report "No RTL-SDR device found" in the dashboard. It will not crash or block other modules.

### COM Port Access for LoRa / Meshtastic

Meshtastic-compatible LoRa devices appear as COM ports on Windows (e.g., `COM3`, `COM4`).

- **No special drivers needed** for most devices — Windows 10+ includes USB-CDC-ACM drivers
- Some devices (e.g., those using CH340/CH341 chips) may require a [manufacturer driver](http://www.wch-ic.com/downloads/CH341SER_EXE.html)
- AEGIS scans available COM ports automatically via `go.bug.st/serial`
- The user can also specify a port manually in the settings

**Permissions:** No admin rights are required to access COM ports on Windows.

### Windows Defender / SmartScreen

When running an unsigned binary downloaded from the internet, Windows SmartScreen may show a warning:

1. Click **"More info"**
2. Click **"Run anyway"**

This is a one-time prompt per binary. AEGIS does not require admin rights for core functionality.

### Firewall

Windows Firewall may prompt when AEGIS starts (especially if configured for LAN access). For localhost-only operation, no firewall rule is needed.

### Path Length Limitations

Windows has a default `MAX_PATH` limit of 260 characters. AEGIS keeps its data directory structure shallow to avoid issues. If running from a deeply nested USB path, consider enabling [long path support](https://learn.microsoft.com/en-us/windows/win32/fileio/maximum-file-path-limitation).

---

## macOS

### Supported Versions

- macOS 13 (Ventura) and later
- Apple Silicon (arm64) — primary target
- Intel Macs (amd64) — not in the default build matrix but can be compiled from source

### Gatekeeper May Block Unsigned Binaries

macOS Gatekeeper will quarantine unsigned binaries downloaded from the internet. You will see:

> **"aegis" can't be opened because Apple cannot check it for malicious software.**

**Workaround (one-time per binary):**

```bash
# Option 1: Remove the quarantine attribute
xattr -d com.apple.quarantine ./aegis

# Option 2: Right-click → Open in Finder, then click "Open" in the dialog
```

Alternatively, from **System Settings → Privacy & Security**, scroll down and click **"Open Anyway"** next to the blocked application notice.

**This also applies to sidecar binaries.** You may need to un-quarantine each sidecar binary individually:

```bash
xattr -rd com.apple.quarantine ./sidecars/
```

### Serial Port Naming

On macOS, serial ports use the `/dev/cu.*` naming convention:

| Device | Typical Port Name |
|--------|------------------|
| Meshtastic LoRa (USB-serial) | `/dev/cu.usbserial-*` or `/dev/cu.SLAB_USBtoUART` |
| Meshtastic LoRa (Bluetooth) | `/dev/cu.meshtastic-*` |
| RTL-SDR | Not a serial device — uses `libusb` |

AEGIS scans `/dev/cu.*` automatically. The `/dev/tty.*` counterparts are not used (they have different flow-control behavior on macOS).

### Homebrew Dependencies (Optional)

Core AEGIS requires **no Homebrew packages**. However, for development:

```bash
brew install go node
```

### RTL-SDR on macOS

RTL-SDR generally works out of the box on macOS with the bundled sidecar binaries. No driver swap is needed (unlike Windows). If issues arise:

```bash
# Check if the device is recognized
system_profiler SPUSBDataType | grep -i rtl
```

### macOS-Specific Power Detection

AEGIS uses `pmset -g batt` and IOKit framework queries for battery status detection. This works on both MacBooks and desktop Macs (which report "AC Power" with no battery).

---

## Linux (Desktop)

### Supported Distributions

- Debian 12+ / Ubuntu 22.04+
- Fedora 38+
- Other distributions with glibc 2.31+ should work but are not tested
- x86_64 (amd64) architecture

### udev Rules for USB Devices

Linux requires udev rules to allow non-root access to USB devices (RTL-SDR dongles, serial adapters).

**RTL-SDR udev rule:**

Create `/etc/udev/rules.d/20-rtlsdr.rules`:

```
# RTL-SDR devices
SUBSYSTEM=="usb", ATTRS{idVendor}=="0bda", ATTRS{idProduct}=="2832", MODE="0666", GROUP="plugdev"
SUBSYSTEM=="usb", ATTRS{idVendor}=="0bda", ATTRS{idProduct}=="2838", MODE="0666", GROUP="plugdev"
```

**Meshtastic / LoRa serial udev rule:**

Create `/etc/udev/rules.d/20-meshtastic.rules`:

```
# Silicon Labs CP210x (common Meshtastic USB-serial)
SUBSYSTEM=="tty", ATTRS{idVendor}=="10c4", ATTRS{idProduct}=="ea60", MODE="0666", GROUP="dialout"

# WCH CH340/CH341
SUBSYSTEM=="tty", ATTRS{idVendor}=="1a86", ATTRS{idProduct}=="7523", MODE="0666", GROUP="dialout"
```

**After creating rules, reload:**

```bash
sudo udevadm control --reload-rules
sudo udevadm trigger
```

### Serial Port Permissions (`/dev/ttyUSB*`)

Linux serial devices require the user to be in the `dialout` group (Debian/Ubuntu) or `uucp` group (Arch/Fedora):

```bash
# Debian / Ubuntu
sudo usermod -aG dialout $USER

# Arch / Fedora
sudo usermod -aG uucp $USER
```

**Log out and back in** (or reboot) for group changes to take effect.

Without this, AEGIS will report "Permission denied" when trying to access LoRa devices on `/dev/ttyUSB*` or `/dev/ttyACM*`.

### RTL-SDR Kernel Module Conflict

The Linux kernel includes a `dvb_usb_rtl28xxu` module that claims RTL-SDR dongles for DVB-T TV reception, preventing `rtl_fm` from accessing them.

**Blacklist the conflicting module:**

```bash
echo 'blacklist dvb_usb_rtl28xxu' | sudo tee /etc/modprobe.d/blacklist-rtlsdr.conf
sudo modprobe -r dvb_usb_rtl28xxu
```

This persists across reboots. AEGIS will detect when the kernel module has claimed the device and display a specific error message pointing to this documentation.

### Executable Permissions

After downloading or extracting AEGIS, the binary may not have execute permission:

```bash
chmod +x ./aegis
chmod +x ./sidecars/*/linux/*
```

### AppArmor / SELinux

On systems with AppArmor (Ubuntu) or SELinux (Fedora/RHEL), the security framework may restrict sidecar process spawning. If modules fail to start:

- **AppArmor:** Check `dmesg | grep apparmor` for denials
- **SELinux:** Check `ausearch -m AVC` for denials

Running AEGIS from a user home directory or USB mount typically avoids these issues.

---

## Raspberry Pi / ARM64

### Hardware Requirements

- **Architecture:** ARM64 (aarch64) is **required**. 32-bit ARMv7 (armhf) is **not supported**.
- **Minimum device:** Raspberry Pi Zero 2W (1 GB RAM)
- **Recommended:** Raspberry Pi 4 (4 GB+ RAM) or Pi 5

### Compatible Devices

| Device | RAM | Tier | Notes |
|--------|-----|------|-------|
| Pi Zero 2W | 1 GB | Minimum | Knowledge + maps + mesh only; no AI |
| Pi 3B+ | 1 GB | Minimum | Same as Zero 2W, slightly faster |
| Pi 4 (2 GB) | 2 GB | Standard | Can run small LLM models |
| Pi 4 (4–8 GB) | 4–8 GB | Standard | Comfortable for most modules |
| Pi 5 (8 GB) | 8 GB | Standard | Best Pi experience |

### OS Requirements

- **Raspberry Pi OS (64-bit)** — based on Debian, primary test target
- **Ubuntu Server 22.04+ (arm64)** — also supported
- 32-bit Raspberry Pi OS is **not supported** (no arm64 Go binary)

### Memory Constraints

On the Pi Zero 2W (1 GB RAM), AEGIS operates in **Minimum tier**:

- AI Assistant is **disabled** (llama-server requires ≥ 2 GB to load even tiny models)
- SDR Monitor is **disabled** (rtl_fm uses ~200 MB for wideband scanning)
- Plant ID is **disabled** (ONNX runtime needs ~500 MB)
- Reticulum P2P is **disabled** (rnsd + Python runtime overhead)

The following modules remain available:
- Knowledge Library (Kiwix — memory usage scales with ZIM file count)
- Offline Maps (MapLibre renders in the browser, not on the Pi)
- Notes
- Mesh Messaging (minimal memory footprint)
- Skill Trees
- Celestial Navigation
- Medical Triage
- Position Beacon

### Swap Configuration

For Pi devices with limited RAM, a swap file can extend capacity at the cost of SD card wear:

```bash
sudo dphys-swapfile swapoff
sudo sed -i 's/CONF_SWAPSIZE=100/CONF_SWAPSIZE=1024/' /etc/dphys-swapfile
sudo dphys-swapfile setup
sudo dphys-swapfile swapon
```

> ⚠️ Excessive swap on SD cards will reduce card lifespan. Consider using a USB SSD for the AEGIS data directory on long-running base stations.

### GPIO and Serial

On Raspberry Pi, the hardware UART is available at `/dev/ttyAMA0` (or `/dev/serial0`). If connecting a LoRa module directly via GPIO (without USB), disable the serial console:

```bash
sudo raspi-config
# → Interface Options → Serial Port
# → Login shell over serial: No
# → Serial port hardware enabled: Yes
```

### Headless Base Station

The `boot/pi-image/` directory contains configuration for running AEGIS as a dedicated always-on base station:

- Auto-starts AEGIS on boot via systemd
- Configures the Pi as a WiFi access point (hostapd) so other devices can connect
- Optimized for minimal resource usage

---

## General Notes

### Kolibri (Education Module) Is Not Bundled

Kolibri is a Python/Django application that provides educational content (Khan Academy, etc.). It is **not included in the core AEGIS binary** due to:

1. Python dependency conflicts with zero-install philosophy
2. Significant packaging weight (~500 MB+ with content)
3. Django application server adds resource overhead

**Available as an optional add-on via:**

- **Docker overlay** (Linux x86_64 only):
  ```bash
  docker run -d --name kolibri \
    -p 8081:8080 \
    -v aegis-data/kolibri:/root/.kolibri \
    learningequality/kolibri
  ```
- **Native install** (all platforms): Follow [Kolibri's install guide](https://kolibri.readthedocs.io/en/latest/install.html)

AEGIS can optionally proxy to a running Kolibri instance if configured in `config.json`.

### USB Drive Filesystem

For maximum cross-platform compatibility of the portable data directory:

| Filesystem | Windows | macOS | Linux | Max File Size | Recommended |
|-----------|:-------:|:-----:|:-----:|--------------|:-----------:|
| exFAT | ✅ | ✅ | ✅ | 128 PB | ✅ Yes |
| NTFS | ✅ | Read-only* | ✅ | 16 TB | ⚠️ No |
| FAT32 | ✅ | ✅ | ✅ | 4 GB | ❌ No (too small for ZIM/models) |
| ext4 | ❌ | ❌ | ✅ | 16 TB | ❌ No (Linux only) |

**Recommendation:** Format USB drives as **exFAT** for maximum compatibility across all three target operating systems.

*\* macOS can read NTFS natively but requires third-party software or kernel extensions to write.*

### Content Pack Sizes

Plan your USB drive capacity accordingly:

| Content | Approximate Size |
|---------|-----------------|
| AEGIS binary (single OS) | ~30–50 MB |
| Sidecars (single OS) | ~200–500 MB |
| Wikipedia ZIM (English, all) | ~95 GB |
| Wikipedia ZIM (English, top 100k) | ~12 GB |
| WikiMed (medical subset) | ~1.5 GB |
| Offline maps (regional, one country) | ~500 MB–5 GB |
| LLM model (TinyLlama 1.1B, Q4) | ~700 MB |
| LLM model (Mistral 7B, Q4) | ~4 GB |
| LLM model (Llama 3 8B, Q4) | ~5 GB |

### Network Security Reminder

AEGIS has **no authentication** by default. When binding to a LAN-accessible address:

- Other devices on the network can access all AEGIS functionality
- Do not expose port 8080 to the public internet
- Use OS-level firewall rules to restrict access if needed
- Consider binding to a specific interface rather than `0.0.0.0`

---

*This document is updated whenever a new platform-specific issue is discovered. If you encounter an undocumented caveat, please open an issue or submit a PR.*
