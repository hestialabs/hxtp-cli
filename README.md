# HxTP CLI (hxtp-cli)

The official Command Line Interface for the **Hestia Labs Cross-Platform Trust Protocol (HxTP)**. Built for developers to manage smart hardware, secure device handshakes, and interact with the Hestia Cloud ecosystem.

[![SLSA 3](https://slsa.dev/images/badge-level3.svg)](https://slsa.dev)
[![Go Version](https://img.shields.io/github/go-mod/go-version/hestialabs/hxtp-cli)](https://go.dev)

## Features

- **Zero-Trust Security**: Native HxTP/3.1 signing and verification for every command.
- **Multi-Transport**: Seamless switching between REST, MQTT, and WebSocket layers.
- **Cross-Platform**: Native binaries for Linux, macOS, and Windows.
- **SLSA Level 3**: High-integrity builds with verifiable provenance.
- **Developer First**: Interactive TUI wizards powered by `huh` and `lipgloss`.

## Installation

### Shell (Linux & macOS)

```bash
curl -fsSL https://hestialabs.in/install.sh | bash
```

### Windows

Download the latest `.zip` from the [Releases](https://github.com/hestialabs/hxtp-cli/releases) page and add the binary to your PATH.

## Quick Start

### 1. Login
Securely authenticate the CLI with your Hestia Cloud account.
```bash
hxtp-cli login
```

### 2. Discover Devices
List all devices registered to your Smart Spaces.
```bash
hxtp-cli device list
```

### 3. Send Commands
Control your hardware directly from the terminal.
```bash
hxtp-cli send <device_id> <action> --param "brightness=80"
```

### 4. Interactive Device Registration
Add new hardware to a space using the setup wizard.
```bash
hxtp-cli device add
```

## Configuration

The CLI stores its configuration in `~/.hxtp/config.json`. Sensitive tokens are stored securely in the system keychain (Keyring).

## Development

Requires [Go 1.24+](https://go.dev).

```bash
# Clone the repo
git clone https://github.com/hestialabs/hxtp-cli.git
cd hxtp-cli

# Build locally
go build -o hxtp-cli ./cmd/hxtp-cli

# Run tests
go test ./...
```

## Security

HxTP/3.1 mandates pipe-separated canonical framing for all signatures. This CLI implements bit-perfect parity with the `hxtp-go` and `hxtp-js` SDKs.

## License

MIT © [HestiaLabs](https://hestialabs.in)
