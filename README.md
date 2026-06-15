# 🛡️ fraud-shield-go

Anti-fraud microservice built in Go — real-time risk analysis via REST API, fixed-rules score engine and event publishing via Apache Kafka.

## Tech Stack

- **Go 1.22** — core language
- **Chi** — lightweight HTTP router
- **confluent-kafka-go** — Kafka producer
- **Apache Kafka** — event streaming
- **Docker Compose** — local infrastructure

## Fraud Rules

| Rule | Condition | Score |
|------|-----------|-------|
| High amount | `> $10,000` | +40 |
| Medium amount | `> $5,000` | +20 |
| High-risk country | KP, IR, CU, SY | +40 |
| Suspicious hours | 00h–05h UTC | +20 |
| Missing merchant | `merchant_id` empty | +30 |

**Decision:**
- `score >= 70` → `BLOCKED` + `HIGH`
- `score 40–69` → `APPROVED` + `MEDIUM`
- `score < 40` → `APPROVED` + `LOW`

## API Reference

### POST /transactions

Analyzes a transaction and publishes the result to Kafka.

**Request**
```json
{
  "account_id": "acc-001",
  "amount": 150.00,
  "merchant_id": "merch-42",
  "merchant_country": "BR"
}
```

**Response — Approved (200)**
```json
{
  "transaction": {
    "id": "e3b0c442-...",
    "account_id": "acc-001",
    "amount": 150,
    "merchant_id": "merch-42",
    "merchant_country": "BR",
    "created_at": "2026-06-15T13:31:43Z"
  },
  "result": {
    "transaction_id": "e3b0c442-...",
    "score": 0,
    "risk_level": "LOW",
    "decision": "APPROVED",
    "reasons": [],
    "analyzed_at": "2026-06-15T13:31:43Z"
  }
}
```

**Response — Blocked (403)**
```json
{
  "transaction": {
    "id": "f7a1d903-...",
    "account_id": "acc-002",
    "amount": 12000,
    "merchant_id": "merch-99",
    "merchant_country": "KP",
    "created_at": "2026-06-15T13:31:43Z"
  },
  "result": {
    "transaction_id": "f7a1d903-...",
    "score": 80,
    "risk_level": "HIGH",
    "decision": "BLOCKED",
    "reasons": [
      "transaction amount exceeds $10,000",
      "merchant located in high-risk country"
    ],
    "analyzed_at": "2026-06-15T13:31:43Z"
  }
}
```

### GET /health

```json
{ "status": "ok" }
```

## Makefile Commands

| Command | Description |
|---------|-------------|
| `make run` | Run the server |
| `make build` | Build binary to `bin/` |
| `make infra` | Start Kafka + Zookeeper + Kafka UI |
| `make infra-down` | Stop infrastructure |
| `make tidy` | Run `go mod tidy` |

## Kafka Event

Every analyzed transaction publishes a `TransactionAnalyzedEvent` to the `transaction.analyzed` topic:

```json
{
  "transaction": { ... },
  "result": { ... }
}
```

The message key is the `transaction_id`, ensuring partition consistency per transaction.

---

Built with ☕ and Go by [joaodddev](https://github.com/joaodddev)