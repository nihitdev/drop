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

```sh
go build ./cmd/drop
```

## Roadmap

- Receive/upload mode
- Multiple input files
- Password protection
- Transfer progress
- mDNS or a friendly local hostname
