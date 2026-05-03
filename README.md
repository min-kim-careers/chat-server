# High-Performance Chat Server

A real-time chat server built in Go, designed as the messaging backbone for the PC Parts Trading Platform. Supports multiple rooms, persistent message history, and low-latency delivery under concurrent load.

---

## Technical Highlights

**Hub pattern** — a central hub manages all active WebSocket connections and coordinates message routing across rooms. Each room maintains its own isolated broadcast group, keeping message fanout efficient regardless of total connection count.

**Dual-layer storage strategy** — Redis handles active session caching and pub/sub for low-latency message delivery. PostgreSQL persists message history, ensuring no data loss across reconnects or server restarts. The separation keeps the hot path fast without sacrificing durability.

**Dual-server architecture** — WebSocket and REST API run as separate servers on independent ports. This decouples real-time traffic from administrative operations (room management, message retrieval) and allows each to scale independently.

**Pluggable authentication** — designed to link directly with an external auth system, allowing the marketplace platform to gate room access without coupling auth logic to the chat server itself.

---

## Stack

| Layer | Technology |
|---|---|
| Language | Go |
| WebSocket | net/http |
| REST API | Gin |
| Cache / Pub-Sub | Redis |
| Persistence | PostgreSQL |
| Containerisation | Docker |

---

## Context

Built as a standalone service to integrate with the [PC Parts Trading Platform](protected.pcmarket.app), providing real-time buyer-seller communication. Designed from the start as a separate deployable rather than an embedded feature, reflecting a service-oriented approach to system design.
