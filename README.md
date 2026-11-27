# 🐱🖍️ CrayonKat

Multi-region uptime and incident monitoring with a highly available control plane, written in Go.

> 🚧 **Early development.** Design is done, and I'm now writing code: **Phase 1, networking**, building the probe that checks whether a site is up.

## Progress

| Phase | What | Status |
|---|---|---|
| 0. Design | Requirements, check rules, alerting policy ([doc](docs/design/000-requirements.md)) | ✅ Done |
| 1. Networking | `probe check` CLI: HTTP checks with timeouts, then TCP/DNS/TLS | 🔨 In progress |
| 2. Control plane | API and PostgreSQL storage for monitors and results | ⏳ Planned |
| 3+ | Shipping, multi-region probes, high availability | ⏳ Planned |

## Why

A monitoring script on one machine isn't reliable. If the machine sleeps, crashes or loses internet, monitoring stops, and a brief network blip can look like an outage. CrayonKat runs probes from several regions and declares a site **DOWN** only when most of them agree, so one probe's bad network can't trigger a false alert.

## What works today

The `probe check` command sends one HTTP request to a URL and reports the result:

```bash
go run ./cmd/probe check https://example.com
# UP 200 https://example.com
```

| Situation | Result | Exit code |
|---|---|---|
| Status 200 to 399 | `UP` | `0` |
| Status 4xx or 5xx | `DOWN` | `1` |
| No response (DNS failure, refused connection, timeout after 10s) | `DOWN` | `1` |
| Wrong arguments | usage message | `1` |

Results go to stdout and errors to stderr, so you can script it:

```bash
go run ./cmd/probe check https://example.com || echo "site is down"
```

## Coming next

- Classify *why* a check failed (DNS, refused, timeout, TLS)
- TCP, DNS and TLS certificate checks
- A control plane that stores results in PostgreSQL and opens incidents
- Probes in 3 regions with a majority vote
- Alerts by Telegram and email

## Requirements

- Go 1.25+

## Docs

- [Requirements](docs/design/000-requirements.md): what CrayonKat must do and why
- [Architecture decisions](docs/adr/): records of key design choices
