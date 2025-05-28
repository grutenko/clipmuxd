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

```
+------------------+                                         +---------------------+
| Инициатор (A)    |                                         | Принимающий (B)      |
+------------------+                                         +---------------------+
        |                                                        |
        | --- InitHandshake(PubKeyA, UUID_A, Nonce) ------------>|
        |                                                        |
        | <--------- InitHandshake(PubKeyB, UUID_B) -------------|
        |                                                        |
        |           Проверка UUID_A в B:                          |
        |           Если есть — B отменяет сессию                 |
        |                                                        |
        |           Проверка UUID_B в A:                          |
        |           Если есть — A отправляет CancelSession(UUID_B) |
        |                                                        |
        | --- Генерация shared_key = DH(PrivA, PubB) ------------>|
        |                                                        |
        | --- Отправка Code(nonce) — сигнал показать код -------->|
        |                                                        |
        |                              Отображение кода:         |
        |                              code_B = SHA256(shared_key)[:N]  |
        |                                                        |
        |                      Пользователь вводит код на A      |
        |                                                        |
        | --- Локальная генерация code_A = SHA256(shared_key)[:N] |
        |                                                        |
        | --- Сравнение введённого кода с code_A:                |
        |         - если не совпадает — отмена сессии             |
        |         - если совпадает — продолжаем                    |
        |                                                        |
        | --- A отправляет зашифрованные JWT_A и CertA ---------->|
        |          (шифрование с использованием shared_key)       |
        |                                                        |
        |           B расшифровывает и добавляет JWT_A и CertA    |
        |           в списки доверенных                           |
        |                                                        |
        | <--- B отправляет зашифрованные JWT_B и CertB ----------|
        |          (шифрование с использованием shared_key)       |
        |                                                        |
        |           A расшифровывает и добавляет JWT_B и CertB    |
        |           в списки доверенных                           |
        |                                                        |
        |                Доверие установлено, готово к mTLS      |
+------------------+                                         +---------------------+

```
