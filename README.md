# zenith
![zenith](assets/banner_transparent.svg)

> External multihack for Counter-Strike 1.6 (Build 8684) written in Go

![Go](https://img.shields.io/badge/Go-1.21-blue)
![Platform](https://img.shields.io/badge/Platform-Windows-blue)
![Build](https://img.shields.io/badge/Build-8684-green)

![demo](assets/demo.gif)

---

## Features

### ESP
- Player boxes with team colors
- Health bars
- Player names
- Distance display
- Configurable hotkey and colors

### Aimbot
- FOV-based target selection
- Smooth aiming
- Lock-on target tracking
- Configurable FOV radius and smoothness

### Triggerbot
- Automatic fire when crosshair is over an enemy
- Configurable activation key
- Random delay to reduce detection

### Overlay
- Built with [veil](https://github.com/Cameliuu/veil) — a transparent GDI overlay library for Windows
- Click-through, always-on-top
- 60fps repaint loop

---

## Prerequisites
- Counter-Strike 1.6 Steam (Build 8684)
- Windows 10/11
- Go 1.21+

## Building
```bash
git clone https://github.com/Cameliuu/zenith
cd zenith
go build -ldflags="-H windowsgui" -o zenith.exe .
```

## Usage
1. Launch CS 1.6
2. Run `zenith.exe`
3. Use hotkeys to toggle features

## Default Hotkeys
| Key | Feature |
|-----|---------|
| F2  | Toggle ESP |
| F3  | Toggle Aimbot |
| F4  | Toggle Triggerbot |

---

## Architecture

```
zenith/
├── engine/         — game loop, entity reading, feature logic
│   ├── aimbot/     — aimbot logic
│   ├── config/     — feature configuration
│   ├── entity/     — entity reading from memory
│   ├── overlay/    — W2S, coordinate utils
│   └── offsets/    — memory offsets for build 8684
├── memory/         — ReadProcessMemory/WriteProcessMemory wrappers
└── input/          — key binding utilities
```

Built on top of [veil](https://github.com/Cameliuu/veil) — a transparent Windows overlay library also written in Go.

---

## Disclaimer
For educational purposes only.
