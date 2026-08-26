# drop

Share a file or directory with another device on your local network—straight
from the terminal.

`drop` starts a temporary HTTP server and prints a URL and QR code. Open either
one on the receiving device to download the shared item. Directories are packed
into a ZIP automatically, and every share uses a cryptographically random URL
token.

## Features

- Share files and directories over your LAN
- Scan a terminal QR code instead of typing the URL
- Stop automatically after the first successful download
- Allow repeated downloads with `--keep`
- Set an optional expiry time
- Run on Windows, Linux, and macOS

## Usage

```sh
drop file.zip
drop ./folder
drop --keep file.zip
drop --expires 10m file.zip
drop --port 9000 file.zip
```

Options:

| Option | Description | Default |
| --- | --- | --- |
| `--port <number>` | TCP port to listen on; `0` chooses a free port | `0` |
| `--keep` | Keep serving after a successful download | off |
| `--expires <duration>` | Stop after a Go duration such as `30s`, `10m`, or `1h` | no expiry |

By default, the server stops after one completed download. Press Ctrl+C to stop
it at any time.

## Install

### Windows

Run in PowerShell:

```powershell
irm https://raw.githubusercontent.com/nihitdev/drop/main/install.ps1 | iex
```

The installer selects the AMD64 or ARM64 executable from the latest GitHub
release, verifies its SHA-256 checksum, installs it for the current user, and
adds it to the user `PATH`. Open a new terminal, then check the installation:

```powershell
drop --help
```

Windows Firewall may ask for permission the first time `drop.exe` listens on
the network. Allow private-network access so other devices on your LAN can
connect.

### Build from source

[Go 1.27 or newer](https://go.dev/dl/) is required.

```sh
git clone https://github.com/nihitdev/drop.git
cd drop
go build -o drop ./cmd/drop
```

On Windows, use `go build -o drop.exe ./cmd/drop` instead.

To install directly from the Go toolchain:

```sh
go install github.com/nihitdev/drop/cmd/drop@latest
```

Make sure Go's binary directory (usually `~/go/bin`) is on your `PATH`.

## How it works

1. `drop` prepares the selected file. A directory becomes a temporary ZIP.
2. It listens on an available local port and prints a tokenized download URL.
3. The receiving device downloads the item through its web browser.
4. The server stops after one download unless `--keep` is enabled, and any
   temporary archive is removed.

The URL token helps prevent accidental discovery, but `drop` does not encrypt
traffic or require authentication. Use it only on networks you trust, and avoid
sharing sensitive data over public Wi-Fi.

## Development

```sh
go test ./...
go vet ./...
```

Tags matching `v*` run the release workflow and publish checksum-verified
Windows AMD64 and ARM64 executables to GitHub Releases.

## Roadmap

- Receive/upload mode
- Multiple input files
- Password protection
- Transfer progress
- mDNS or a friendly local hostname
