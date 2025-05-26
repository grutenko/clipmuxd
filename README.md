# clipmuxd

**clipmuxd** is a background service that receives and sends clipboard data to other devices.

⚠️ **Work in Progress**
This project is currently under active development. Features and APIs may change.

## How it works

- **clipmuxd** runs in the background managing clipboard data synchronization between devices.
- UI applications connect to the local clipmuxd service to send commands and receive current clipboard state information.

## Features

- On first run, SSL certificates for the client (UI) and clipmuxd service are generated automatically.
- These certificates are used to establish a secure connection from the UI to the local clipmuxd via Unix sockets or Named Pipes using mTLS (mutual TLS authentication).
- Certificate exchange happens during the handshake phase when connections are established.
- The same clipmuxd key is used to establish secure connections between different devices.

## Default ports

- `51523` — gRPC handshake (secure connection establishment between devices)
- `63482` — main gRPC channel for clipboard data and commands
- `49876` — device discovery on the network (UDP)

- For local UI connections:
  - Unix socket: `/run/clipmuxd.sock` (Linux)
  - Named Pipe: `\\.\pipe\clipmuxd` (Windows)

---

clipmuxd provides secure clipboard exchange and data synchronization between local applications and remote devices.
