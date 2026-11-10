# Chat Server

A real-time, multi-room chat server in Go, with fast message delivery and durable message history.

Status: **Active**

## Why
To build a chat backend from first principles: concurrent WebSocket clients, catching up on missed messages, and persistence that doesn't slow down delivery.

## How it works
```
Client ──ws──▶ Room (in-memory fan-out to connected clients)
                 │
                 ├─▶ Redis Stream per room ──▶ catch-up for reconnecting clients
                 └─▶ Redis persist stream ──▶ consumer group ──▶ PostgreSQL (batched)
```
- **Fan-out.** Each room broadcasts to its connected clients in memory, so delivery never waits on a database write.
- **Catch-up.** Each room keeps a Redis Stream of recent messages. Reconnecting clients read everything after their last seen ID (`XRANGE`).
- **Write-behind persistence.** Messages also go to a persist stream. A consumer group reads them in batches (`XREADGROUP`), writes them to PostgreSQL, then acknowledges them (`XACK`), so nothing is lost if the writer crashes.
- **Validation.** Payloads are checked with go-playground/validator.

## Key design decisions
| Decision | Why | Trade-off |
|---|---|---|
| Redis Streams instead of pub/sub | Streams keep messages and track acknowledgements; pub/sub drops anything missed | More moving parts |
| Write-behind to Postgres | Database speed doesn't affect chat delivery | History lags a moment behind live messages |

## Stack
Go 1.23 · gorilla/websocket · Redis Streams · PostgreSQL · Docker Compose · Air (hot reload)

## Running it
Runs on Docker Compose with PostgreSQL and Redis (see `compose.yml`), with hot reload through Air. Connections are authenticated against an external auth service (`AUTH_URL`), which this repo doesn't include.

## What I'd improve next
- Horizontal scaling: fan out across instances through the room streams
- Load tests to measure latency under many concurrent clients
