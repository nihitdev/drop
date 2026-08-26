# drop

Fast local file sharing from the terminal.

`drop` serves one file—or a temporary ZIP of one directory—to another device on
your local network. Each share gets a cryptographically random URL token and a
compact terminal QR code.

## Usage

```sh
drop file.zip
drop ./folder
drop --keep file.zip
drop --expires 10m file.zip
drop --port 9000 file.zip
```

Transfers stay on the local network. By default, the server stops after one
completed download. `--keep` allows multiple downloads, while `--expires`
limits the share lifetime in either mode. Press Ctrl+C to stop the server at any
time.

Directories are ZIP-compressed into a temporary file, which is removed when the
server stops. The project is currently focused on sending files.

## Build

Linux and macOS:

```sh
go build ./cmd/drop
```

Windows (PowerShell):

```powershell
go build -o drop.exe ./cmd/drop
.\drop.exe file.zip
```

To cross-compile a Windows binary from Linux or macOS:

```sh
GOOS=windows GOARCH=amd64 go build -o drop.exe ./cmd/drop
```

The CLI uses the same flags and behavior on Windows, Linux, and macOS. Windows
Firewall may ask for permission the first time `drop.exe` listens on the local
network; allow access on private networks for LAN sharing.

## Install on Windows

From PowerShell:

```powershell
irm https://raw.githubusercontent.com/nihitdev/drop/main/install.ps1 | iex
```

The installer downloads the latest executable for AMD64 or ARM64 Windows,
verifies its SHA-256 checksum, installs it under the current user's local app
directory, and adds that directory to the user PATH. Open a new terminal after
installation, then run `drop <file>`.

Version tags matching `v*` trigger the release workflow, which tests the project
and publishes both Windows executables and their checksums to GitHub Releases.

## Roadmap

- Receive/upload mode
- Multiple input files
- Password protection
- Transfer progress
- mDNS or a friendly local hostname
