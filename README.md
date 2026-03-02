# AgentOps Platform — Autonomous DevOps Agent

An autonomous DevOps agent built in **Go** that **monitors**, **decides**, and **acts** on infrastructure events without human intervention, with a full audit trail.

## The Problem

Engineering teams waste 30–40% of their time on repetitive infrastructure work: incident response, deployment pipelines, cost optimization, and security patching. Existing tools require a human at every step.

## How It Works

**Core agent loop:** **Monitor → Analyze → Plan → Execute → Verify → Learn**

| Phase    | Description |
|----------|-------------|
| **Monitor** | Streams events from AWS / GCP / K8s via Go goroutines (Kafka, EventBridge, etc.) |
| **Analyze** | Uses LLM APIs to classify and prioritize issues and produce a remediation plan |
| **Plan**    | Multi-step remediation plan (scale, restart, patch, rollback, etc.) |
| **Execute** | Runs actions via Go SDK clients (AWS SDK, kubectl, GCP client) |
| **Verify**  | Confirms resolution and triggers rollback if needed |
| **Learn**   | Stores outcomes in a vector DB for future decisions and full audit |

## Features

- **Auto-Incident Response** — Detects alerts, diagnoses root cause, applies fixes autonomously  
- **Cost Guardian** — Rightsizes resources, terminates idle instances  
- **Security Sentinel** — Auto-patches CVEs, rotates secrets, enforces policies  
- **Deployment Surgeon** — Monitors deploys, auto-rolls back on anomaly detection  
- **Audit Trail** — Every action logged with reasoning for compliance  

## Tech Stack

- **Agent core** — Go (goroutines, channels, context)
- **LLM** — OpenAI / Anthropic-compatible API
- **Event streams** — Kafka / AWS EventBridge (stub included for dev)
- **State** — Redis + PostgreSQL (optional); in-memory stub for dev
- **Vector memory** — Qdrant / Weaviate (planned); in-memory store for dev
- **Observability** — OpenTelemetry + Prometheus (planned)
- **Deployment** — Single Go binary or Docker sidecar  

## Quick Start

```bash
# Resolve dependencies (generates go.sum)
go mod tidy

# Copy env and set your LLM API key (optional for stub-only run)
cp .env.example .env

# Run (stub stream + optional LLM)
go run ./cmd/agent
```

With an LLM API key set, the agent will analyze each event and produce plans; with stub executor/verifier it won’t change real infrastructure.

## Build

```bash
go build -o agent ./cmd/agent
```

## Docker

```bash
docker build -t agentops/agent .
docker run --env-file .env agentops/agent
```

## Project Layout

```
cmd/agent/          # Main entrypoint
internal/
  agent/            # Core loop: Monitor → Analyze → Plan → Execute → Verify → Learn
  monitor/          # Event streams (Stream interface + stub)
  analyze/          # LLM analyzer + noop
  execute/          # Executor interface + stub
  verify/           # Verifier interface + stub
  learn/            # Outcome store (memory stub; vector DB later)
  types/            # Event, Plan, Action, AuditEntry
pkg/config/        # Config from env
```

## Why Go?

- **Concurrency** — Thousands of event streams via goroutines/channels  
- **Performance** — Sub-millisecond decision loops at scale  
- **Single binary** — Easy deployment as sidecar or standalone  
- **Strong typing** — Safer agentic actions on critical infrastructure  

## License

MIT
