# Communications Protocols — Project AEGIS

> Overview of the off-grid communication systems integrated into AEGIS.

---

## Table of Contents

- [1. Overview](#1-overview)
- [2. Meshtastic / LoRa Mesh Messaging](#2-meshtastic--lora-mesh-messaging)
- [3. Reticulum Encrypted P2P](#3-reticulum-encrypted-p2p)
- [4. SDR Monitoring (Receive-Only)](#4-sdr-monitoring-receive-only)
- [5. libp2p Local Peer Sync](#5-libp2p-local-peer-sync)
- [6. APRS-Style Position Beacon](#6-aprs-style-position-beacon)
- [7. Protocol Comparison](#7-protocol-comparison)

---

## 1. Overview

AEGIS integrates five complementary communication systems, each serving a distinct purpose. No single protocol covers all scenarios — together, they provide layered communication from close-range mesh to wide-area monitoring.

```
┌─────────────────────────────────────────────────────────────────┐
│                    AEGIS Communications Stack                    │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌─────────────┐  ┌─────────────┐  ┌──────────────────────────┐ │
│  │  Meshtastic  │  │ Reticulum   │  │     libp2p Peer Sync    │ │
│  │  LoRa Mesh   │  │ Encrypted   │  │  (WiFi / Ethernet LAN)  │ │
│  │  (915 MHz)   │  │ P2P Network │  │                          │ │
│  │              │  │             │  │  mDNS discovery +        │ │
│  │  Range: km   │  │ Range: any  │  │  content replication     │ │
│  │  Speed: slow │  │ transport   │  │                          │ │
│  │  No infra    │  │             │  │  Range: LAN              │ │
│  └──────────────┘  └─────────────┘  └──────────────────────────┘ │
│                                                                  │
│  ┌──────────────────────────┐  ┌────────────────────────────┐   │
│  │   SDR Monitor            │  │   APRS Position Beacon     │   │
│  │   (RTL-SDR, RX only)     │  │   (over LoRa mesh)         │   │
│  │                          │  │                             │   │
│  │   Passive scanning of    │  │   Periodic lat/lon/alt      │   │
│  │   radio spectrum         │  │   broadcast for mutual      │   │
│  │   24 MHz – 1.766 GHz     │  │   awareness                 │   │
│  └──────────────────────────┘  └────────────────────────────┘   │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Hardware Requirements

| Protocol | Required Hardware | Min. Tier |
|----------|------------------|-----------|
| Meshtastic mesh | LoRa radio (T-Beam, T-Echo, etc.) + USB cable | Minimum |
| Reticulum P2P | Any network interface (LoRa, WiFi, serial) | Optimal |
| SDR Monitor | RTL-SDR dongle (RTL2832U-based) | Standard |
| libp2p sync | WiFi or Ethernet (existing LAN) | Standard |
| APRS beacon | LoRa radio (same as Meshtastic) | Minimum |

---

## 2. Meshtastic / LoRa Mesh Messaging

### What Is Meshtastic?

[Meshtastic](https://meshtastic.org/) is an open-source LoRa mesh networking project. It enables long-range, low-power text messaging between devices with **no cellular or internet infrastructure**. Messages hop through intermediate nodes to extend range across the mesh.

### How AEGIS Integrates Meshtastic

AEGIS implements a **native Go client** that communicates with a connected Meshtastic radio device over serial (USB). This avoids the need for the official Meshtastic Python CLI or Android/iOS app.

```
┌──────────────┐      Serial (USB)      ┌──────────────┐
│  AEGIS       │ ◄────────────────────► │  LoRa Radio  │
│  Go Backend  │    Protobuf messages   │  (T-Beam,    │
│  (meshmsg    │                        │   T-Echo,    │
│   module)    │                        │   etc.)      │
└──────────────┘                        └──────┬───────┘
                                               │
                                          LoRa RF
                                         (868/915 MHz)
                                               │
                                        ┌──────▼───────┐
                                        │  Other Mesh  │
                                        │  Nodes       │
                                        └──────────────┘
```

### Protocol Details

| Parameter | Value |
|-----------|-------|
| Frequency | 868 MHz (EU), 915 MHz (US), 923 MHz (AU/NZ) |
| Modulation | LoRa (CSS — Chirp Spread Spectrum) |
| Data rate | ~1 kbps (varies by spreading factor) |
| Range | 1–10+ km (line of sight), 0.5–3 km (urban) |
| Max message size | ~228 bytes (per Meshtastic packet) |
| Power | Typically 100 mW – 1 W |
| Encryption | AES-256 (Meshtastic default channel encryption) |
| Routing | Flooding-based mesh with hop limits |

### Serial Communication

AEGIS uses [`go.bug.st/serial`](https://pkg.go.dev/go.bug.st/serial) for cross-platform serial access:

| OS | Port Pattern | Example |
|----|-------------|---------|
| Windows | `COMx` | `COM3` |
| Linux | `/dev/ttyUSB*`, `/dev/ttyACM*` | `/dev/ttyUSB0` |
| macOS | `/dev/cu.*` | `/dev/cu.usbserial-0001` |

### Protobuf Schema

AEGIS communicates with the radio using the [Meshtastic protobuf schema](https://github.com/meshtastic/protobufs). Key message types:

```protobuf
message MeshPacket {
    uint32 from = 1;
    uint32 to = 2;
    oneof payload_variant {
        Data decoded = 3;
        bytes encrypted = 4;
    }
    uint32 channel = 7;
    uint32 hop_limit = 10;
}

message Data {
    PortNum portnum = 1;
    bytes payload = 2;
}
```

### API Endpoints

```
GET  /api/mesh/nodes             # List discovered mesh nodes
POST /api/mesh/messages          # Send a message
WS   /api/mesh/messages/stream   # Real-time message stream
```

### Capabilities and Limitations

✅ **What Meshtastic does well:**
- Long-range communication without infrastructure
- Low power consumption (battery-powered radios last days)
- Mesh routing extends effective range
- Works in areas with zero cell coverage

⚠️ **Limitations:**
- Very low data rate (~1 kbps) — text only, no images/files
- Messages are small (228 bytes max per packet)
- Higher latency than IP-based networks (seconds to minutes for multi-hop)
- Requires a compatible LoRa radio device

---

## 3. Reticulum Encrypted P2P

### What Is Reticulum?

[Reticulum](https://reticulum.network/) (RNS) is an encrypted networking stack designed for unstable, low-bandwidth, and intermittent connectivity. It operates over any transport — LoRa, serial, WiFi, TCP/IP, or even audio modems — and provides end-to-end encryption by default.

### How AEGIS Integrates Reticulum

Reticulum's reference implementation is in Python. Since AEGIS is a Go application, Reticulum runs as a **sidecar process** — a PyInstaller-frozen `rnsd` binary that requires no Python installation on the user's machine.

```
┌──────────────┐   HTTP Bridge   ┌──────────────┐
│  AEGIS       │ ◄──────────────►│    rnsd      │
│  Go Backend  │   localhost     │  (Reticulum  │
│  (reticulum  │                 │   Daemon)    │
│   module)    │                 │              │
└──────────────┘                 └──────┬───────┘
                                        │
                                   Any Transport
                                  (LoRa, WiFi,
                                   Serial, TCP)
                                        │
                                 ┌──────▼───────┐
                                 │  Other RNS   │
                                 │  Nodes       │
                                 └──────────────┘
```

### Key Properties

| Property | Detail |
|----------|--------|
| Encryption | End-to-end, built-in (Curve25519/AES/HMAC) |
| Identity | Cryptographic identity (public key hash) |
| Transport-agnostic | Works over any byte stream |
| Delay-tolerant | Designed for intermittent links |
| No central server | Fully decentralized |
| Interoperability | Compatible with Nomad Network, Sideband app |

### Network Topology

Reticulum does not require a specific topology. Nodes discover each other through **announces** — cryptographically signed identity broadcasts.

```
Node A ──LoRa──► Node B ──WiFi──► Node C
                   │
                 Serial
                   │
                   ▼
                Node D
```

Each node maintains a routing table built from received announces. Messages are forwarded hop-by-hop using the destination's hash address.

### Identity Management

AEGIS stores Reticulum identity keys in `aegis-data/identity/`:

```
aegis-data/identity/
├── identity            # Reticulum identity key (Curve25519 private key)
└── known_destinations  # Cache of discovered peer identities
```

> ⚠️ **Security note:** The identity file is the cryptographic root of trust. Protect this file — anyone with access can impersonate the node.

### Interoperability

AEGIS Reticulum nodes are interoperable with the broader Reticulum ecosystem:

- **Nomad Network** — Terminal-based P2P communication tool
- **Sideband** — Android/Desktop Reticulum messenger
- **LXMF** — Delay-tolerant message transfer protocol built on Reticulum

### Hardware Tier

Reticulum requires **Optimal tier** (≥ 16 GB RAM) because:
- The `rnsd` PyInstaller binary has significant memory overhead (~300–500 MB)
- Additional memory needed for the transport layer
- Not practical on RAM-constrained devices like Pi Zero 2W

---

## 4. SDR Monitoring (Receive-Only)

### What Is SDR?

Software-Defined Radio (SDR) uses software to process radio signals that would traditionally require dedicated hardware circuits. An RTL-SDR dongle (based on the RTL2832U chipset) is a ~$25 USB device that can receive radio signals from **24 MHz to 1.766 GHz**.

### AEGIS SDR Architecture

AEGIS uses the `rtl_fm` and `rtl_power` command-line tools as **sidecar processes** for signal reception. These are **receive-only** — AEGIS never transmits, avoiding all amateur radio licensing requirements.

```
┌──────────────┐    stdout (JSON)   ┌──────────────┐     USB     ┌───────────┐
│  AEGIS       │ ◄─────────────────│  rtl_fm /    │ ◄──────────│  RTL-SDR  │
│  Go Backend  │                    │  rtl_power   │             │  Dongle   │
│  (sdrmonitor │                    │  (sidecar)   │             │           │
│   module)    │                    └──────────────┘             └───────────┘
└──────────────┘                                                      │
                                                                  Antenna
                                                                 (24 MHz –
                                                                  1.766 GHz)
```

### Supported Monitoring Modes

| Mode | Tool | Description |
|------|------|-------------|
| **Narrowband FM** | `rtl_fm` | Demodulate a single FM channel (weather radio, emergency broadcasts) |
| **Wideband power scan** | `rtl_power` | Scan a frequency range and measure signal power — spectrum activity overview |
| **AM demodulation** | `rtl_fm -M am` | Receive AM broadcasts (shortwave, aviation) |

### Frequency Ranges of Interest

| Band | Frequency | Use |
|------|-----------|-----|
| FM broadcast | 87.5–108 MHz | Commercial radio stations |
| NOAA Weather | 162.400–162.550 MHz | US weather alerts |
| Marine VHF | 156–162 MHz | Maritime emergency (Ch. 16: 156.800 MHz) |
| Aviation | 118–136 MHz (AM) | Air traffic control |
| FRS/GMRS | 462–467 MHz | Family/general mobile radio |
| Amateur 2m | 144–148 MHz | Ham radio |
| Amateur 70cm | 420–450 MHz | Ham radio |
| ISM (LoRa) | 868/915 MHz | LoRa signals (metadata only) |

### Platform-Specific Setup

| OS | Setup Required |
|----|---------------|
| **Windows** | ⚠️ Zadig WinUSB driver swap (see [PLATFORM_NOTES.md](PLATFORM_NOTES.md#rtl-sdr-requires-zadig-winusb-driver-swap)) |
| **Linux** | Blacklist `dvb_usb_rtl28xxu` kernel module; add udev rule |
| **macOS** | Generally works out of the box |

### API Endpoints

```
GET /api/comms/sdr/scan?band=noaa     # Scan a predefined band
GET /api/comms/sdr/scan?freq=162.4M   # Scan a specific frequency
GET /api/comms/sdr/spectrum           # Get current spectrum power data
```

### Important Notes

- 📡 **Receive-only by design.** AEGIS never transmits RF signals. No radio license is required.
- 🔇 **Audio demodulation** is processed server-side; the browser receives decoded data (not raw audio).
- ⚡ SDR monitoring requires **Standard tier** hardware due to `rtl_fm` memory usage (~200 MB for wideband).

---

## 5. libp2p Local Peer Sync

### What Is libp2p?

[libp2p](https://libp2p.io/) is a modular peer-to-peer networking library — the same library that powers IPFS. AEGIS uses it for **local network peer discovery and content synchronization** without requiring any internet connection or custom kernel modules.

### Why libp2p Instead of batman-adv?

The original design considered `batman-adv` (Better Approach To Mobile Adhoc Networking), which is a Linux kernel module that creates a layer-2 mesh network. However:

- ❌ `batman-adv` is a **Linux kernel module** — it cannot run on Windows or macOS
- ❌ Requires root/admin privileges to load
- ❌ Not available on Raspberry Pi OS by default

libp2p solves all of these issues:

- ✅ Pure userspace — no kernel modules
- ✅ Cross-platform (Go implementation)
- ✅ mDNS-based discovery works on any local network
- ✅ No elevated privileges required

### Architecture

```
AEGIS Instance A                           AEGIS Instance B
┌──────────────┐                          ┌──────────────┐
│  peersync     │                          │  peersync     │
│  module       │                          │  module       │
│               │  ◄── mDNS discovery ──► │               │
│  ┌──────────┐ │                          │ ┌──────────┐ │
│  │ libp2p   │ │  ◄── content sync ────► │ │ libp2p   │ │
│  │ host     │ │      (WiFi / Ethernet)  │ │ host     │ │
│  └──────────┘ │                          │ └──────────┘ │
└──────────────┘                          └──────────────┘
```

### Discovery

AEGIS uses **mDNS (Multicast DNS)** to discover other AEGIS instances on the same local network:

1. Each AEGIS instance announces itself via mDNS with a service type (`_aegis._tcp`)
2. When a peer is discovered, AEGIS establishes a libp2p connection
3. Peers exchange metadata (available content packs, versions)
4. Differential sync begins for new or updated content

### Content Sync Protocol

```
Peer A                                    Peer B
  │                                         │
  │  1. mDNS announce (_aegis._tcp)         │
  │ ──────────────────────────────────────► │
  │                                         │
  │  2. Connect (libp2p stream)             │
  │ ◄─────────────────────────────────────  │
  │                                         │
  │  3. Exchange content manifest           │
  │     (list of packs + hashes)            │
  │ ◄────────────────────────────────────►  │
  │                                         │
  │  4. Differential sync                   │
  │     (only transfer missing/updated)     │
  │ ◄────────────────────────────────────►  │
  │                                         │
  │  5. Verify integrity (hash check)       │
  │ ◄────────────────────────────────────►  │
  │                                         │
```

### Syncable Content

| Content Type | Sync Behavior |
|-------------|---------------|
| ZIM content packs | Hash-based diff, full file transfer |
| Map tiles (.pmtiles) | Hash-based diff, full file transfer |
| Notes database | Row-level CRDT sync (conflict-free) |
| Mesh message history | Append-only merge |
| AI model files | Hash-based, user-initiated only (large files) |

### API Endpoints

```
GET  /api/sync/peers    # List discovered peers
POST /api/sync/push     # Initiate sync with a specific peer
```

### Network Requirements

- Both instances must be on the **same local network** (WiFi or Ethernet)
- No internet connection required
- mDNS must not be blocked by the router/firewall (multicast on port 5353)
- Works with consumer WiFi routers, direct WiFi (ad-hoc), or direct Ethernet

---

## 6. APRS-Style Position Beacon

### What Is APRS?

APRS (Automatic Packet Reporting System) is a tactical real-time communication system used in amateur radio for position reporting, weather stations, and short messaging. AEGIS implements an **APRS-style** beacon protocol over the LoRa mesh network.

### How It Works

The beacon module periodically broadcasts a compact position report over the Meshtastic mesh network. This allows all AEGIS nodes in the mesh to maintain mutual position awareness — critical for search-and-rescue, group coordination, and asset tracking.

```
┌──────────────┐    Beacon Packet    ┌──────────────┐
│  AEGIS Node  │ ──────────────────► │  LoRa Mesh   │
│  (beacon     │    (via meshmsg     │  Network     │
│   module)    │     module)         │              │
└──────────────┘                     └──────┬───────┘
                                            │
                                     Received by all
                                     nodes in range
                                            │
                                     ┌──────▼───────┐
                                     │  Other AEGIS  │
                                     │  Nodes        │
                                     │  (map display)│
                                     └──────────────┘
```

### Beacon Packet Format

AEGIS uses a compact binary format for position beacons, transmitted as Meshtastic data packets:

| Field | Size | Description |
|-------|------|-------------|
| Node ID | 4 bytes | Unique node identifier |
| Timestamp | 4 bytes | Unix timestamp (seconds) |
| Latitude | 4 bytes | Signed fixed-point (1e-7 degrees) |
| Longitude | 4 bytes | Signed fixed-point (1e-7 degrees) |
| Altitude | 2 bytes | Meters above sea level |
| Speed | 1 byte | km/h (0–255) |
| Heading | 1 byte | Degrees / 2 (0–179 → 0°–358°) |
| Battery | 1 byte | Percentage (0–100) |
| Status | 1 byte | Bitfield (moving, stationary, emergency) |
| **Total** | **22 bytes** | Fits within single Meshtastic packet |

### Position Sources

The beacon module can obtain position from multiple sources (in priority order):

1. **Manual entry** — User enters coordinates via the UI
2. **GPS via Meshtastic** — Many LoRa devices have built-in GPS (T-Beam)
3. **Celestial navigation module** — Calculated position from manual observations
4. **Last known position** — Cached from any previous source

### Beacon Intervals

| Mode | Default Interval | Configurable |
|------|-----------------|:------------:|
| Stationary | 15 minutes | ✅ |
| Moving | 2 minutes | ✅ |
| Emergency | 30 seconds | ✅ |

### Map Integration

Position beacons from all discovered nodes are plotted on the offline maps module:

- Real-time position updates on the map
- Historical track lines (stored in SQLite)
- Node status indicators (battery, signal strength)
- Emergency alert highlighting

### Differences from Traditional APRS

| Feature | Traditional APRS | AEGIS Beacon |
|---------|-----------------|--------------|
| Transport | AX.25 / VHF radio (144.39 MHz) | LoRa mesh (868/915 MHz) |
| License | Amateur radio license required | No license required (ISM band) |
| Infrastructure | Digipeaters, iGates, internet backbone | Fully self-contained mesh |
| Data format | ASCII text | Compact binary |
| Range | Line-of-sight VHF | LoRa mesh (multi-hop) |

---

## 7. Protocol Comparison

### Decision Matrix

| | Meshtastic | Reticulum | SDR Monitor | libp2p Sync | APRS Beacon |
|---|:---------:|:---------:|:-----------:|:-----------:|:-----------:|
| **Direction** | Bidirectional | Bidirectional | Receive only | Bidirectional | Broadcast |
| **Range** | km (LoRa) | Any transport | 24 MHz–1.7 GHz | LAN only | km (LoRa) |
| **Speed** | ~1 kbps | Varies by transport | N/A | LAN speed | ~1 kbps |
| **Encryption** | AES-256 | E2E (Curve25519) | N/A | TLS 1.3 | None |
| **Infrastructure** | None | None | None | WiFi/Ethernet | None |
| **Extra hardware** | LoRa radio | None* | RTL-SDR dongle | None | LoRa radio |
| **License** | None (ISM band) | None | None (RX only) | None | None (ISM) |
| **Min. tier** | Minimum | Optimal | Standard | Standard | Minimum |
| **Use case** | Off-grid messaging | Secure P2P comms | Signal monitoring | Content sharing | Position tracking |

*\* Reticulum can use any available transport, including the LoRa radio shared with Meshtastic.*

### When to Use What

| Scenario | Recommended Protocol |
|----------|---------------------|
| Send a text message to someone 5 km away, no internet | **Meshtastic** |
| Exchange encrypted messages over any available link | **Reticulum** |
| Listen for emergency weather broadcasts | **SDR Monitor** |
| Share a ZIM content pack with another AEGIS node on the same WiFi | **libp2p Sync** |
| Broadcast your position to nearby AEGIS nodes | **APRS Beacon** |
| Monitor radio activity in your area | **SDR Monitor** |
| Communicate with Sideband or Nomad Network users | **Reticulum** |

---

*This document covers the communication protocols as designed for AEGIS v1. Protocol implementations will be refined during Phase 4 (Comms) of the development roadmap.*
