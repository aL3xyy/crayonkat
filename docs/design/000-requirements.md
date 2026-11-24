# 000 — Requirements

Status: Draft · Phase 0

## Problem

CrayonKat is something I'm building for me and my friends to keep an eye on our websites and APIs and know when they go down. A simple script running on one machine isn't always reliable because it might stop working if the machine crashes, goes to sleep, or loses internet. It can also mistake temporary network issues for actual outages, leading to unnecessary alerts or missed problems.

## Users

- **Me:** I monitor my personal portfolio and side projects.
- **Friends:** they run their **own copy** (self-hosted). There are no shared accounts, and their data never touches my system.
- I'm usually away from my PC, so alerts must reach my **phone**.

## Deployment modes

| Mode | Who | How |
|---|---|---|
| **Single-node** | Friends, anyone who downloads it | One command (`docker compose up`) + one config file |
| **Multi-region HA** | My own deployment | 3 regions, control plane with no single point of failure |

Single-node gets built first. Multi-region builds on top of it.

## What gets monitored

- My portfolio website
- Friends' personal projects (real URLs are needed for testing)
- HTTP(S) endpoints in v1

### Check rules

| Situation | Result |
|---|---|
| Response in under 10s with an expected status | **UP** |
| Response in 10–15s | **DEGRADED** |
| No response within 15s (total timeout) | **DOWN**, reason `TIMEOUT` |
| DNS failure, connection refused | **DOWN** |
| HTTP 5xx | **DOWN** |
| HTTP 4xx | **DOWN** by default; a monitor can list expected codes (e.g. `401` for an auth-protected API) |
| HTTP 429 | **UNKNOWN**, reason `RATE_LIMITED`. The server responded, so it's alive; it is rejecting *our monitor*. |
| Redirects | Follow them (max 10, Go's default). If the final host isn't the same domain → **DOWN**, reason `UNEXPECTED_REDIRECT` |

### Multi-region verdict (majority vote)

| Regional results | Overall |
|---|---|
| 3 UP / 2 UP + 1 DOWN | **UP** |
| 2 DOWN + 1 UP / 3 DOWN | **DOWN** |
| Probe stops reporting | That region counts as **UNKNOWN**, not DOWN |
| Not enough reports for a majority | **UNKNOWN** |

Why a majority: the probes sit in different networks, so their failures are independent. If the site is really down, every probe sees it. If only one does, the problem is most likely near that probe.

Regional failures stay visible on their own, even when the overall status is UP.

**Principle:** never confuse the target failing with our monitoring failing. Silence from a probe tells us about the probe, not the target.

## Alerting

- **Primary:** Telegram bot (free, a simple HTTPS API, no approval process)
- **Secondary:** Email (an independent channel, so alerts still arrive if Telegram is down)
- **Check interval:** 5 min by default, configurable per monitor (shorter for important ones)
- **Confirm before alerting:** 2 failures in a row → about 5–10 min from an outage to the alert
- **Notify on state transitions only**, not on every check
- DEGRADED → warning · DOWN → critical

### SSL certificate expiry

- Warn at **14, 7, 3, 1** days before expiry
- Store sent warnings in Postgres, unique on `(monitor_id, cert_fingerprint, threshold)`. That gives no duplicate alerts, and a renewed certificate has a new fingerprint, so its schedule starts fresh.

## Non-goals (v1)

- User accounts, logins, hosting other people's data
- WhatsApp (needs Meta business verification, templates and per-message cost)
- Configurable redirect allowlists (the same-domain rule is enough for v1)
- Non-HTTP checks in the multi-region verdict (TCP/DNS/TLS come in Phase 1 as probe primitives only)
- A web UI

## Open questions

- How do regions vote when some say **DEGRADED** (e.g. 1 UP, 1 DEGRADED, 1 DOWN)?
- Who gets alerted when a monitor stays **UNKNOWN** for a long time? Silent blindness is dangerous.
- How do we stop **flapping** (a site toggling UP/DOWN and spamming alerts)? Look into hysteresis.
- Where does my multi-region deployment run? → Decide in **Phase 4** (ADR).
- The roadmap's Phase 3 says "multi-tenant schema," which conflicts with the self-hosted decision. → Revisit in **Phase 3**.
