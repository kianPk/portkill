# portkill

**Free a busy TCP port in one command.**

When `EADDRINUSE` ruins your flow — find and stop whatever is listening.

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

## Install

```bash
go install github.com/kianPk/portkill@latest
```

Or clone and build:

```bash
git clone https://github.com/kianPk/portkill.git
cd portkill && go build -o portkill .
```

## Usage

```bash
portkill 3000          # stop listener on :3000
portkill -l 8080       # list only
portkill -f 5173       # force kill (SIGKILL)
```

## Why

Dev servers, Vite, Next, Docker side-cars — ports get stuck. `portkill` wraps `lsof` + signal into a tiny CLI you can remember.

## Requirements

- Linux or macOS
- `lsof` available on `PATH`

## License

MIT © [kianPk](https://github.com/kianPk)
