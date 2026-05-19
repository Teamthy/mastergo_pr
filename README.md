MasterGo

Scalable crypto payments infrastructure built with Go, Next.js, PostgreSQL, Redis, and Ethereum.

MasterGo provides:

- progressive authentication workflows
- secure API key infrastructure
- Ethereum wallet management
- internal transaction ledgering
- webhook delivery pipelines
- audit logging and rate limiting



---

Architecture

Client (Next.js)
        ↓
Go API Layer
        ↓
--------------------------------
Auth Service
Wallet Service
API Key Service
Webhook Service
--------------------------------
        ↓
PostgreSQL + Redis
        ↓
Ethereum Network

---

Core Features

Authentication

- JWT authentication
- OTP email verification
- progressive onboarding flow
- rate-limited auth endpoints

API Key Infrastructure

- hashed secret storage
- key rotation and revocation
- scoped developer access

Wallet Infrastructure

- Ethereum wallet generation
- encrypted private key storage
- internal balance ledger
- transaction history tracking

Platform Reliability

- Redis-backed caching
- webhook retry handling
- audit logging
- request validation
- containerized deployment

---

Tech Stack

Backend

- Go
- Chi Router
- PostgreSQL
- Redis
- go-ethereum
- JWT

Frontend

- Next.js
- TypeScript
- Tailwind CSS
- XState
- ethers.js

Infrastructure

- Docker
- Docker Compose
- Render

---

API Surface

Auth

- POST "/auth/signup"
- POST "/auth/login"
- POST "/auth/verify-email"

Wallet

- POST "/wallet/create"
- GET "/wallet/balance"
- POST "/wallet/withdraw"

API Keys

- POST "/apikeys"
- GET "/apikeys"
- DELETE "/apikeys/:id"

---

Local Development

git clone

docker-compose up --build

Backend: "localhost:8080"
Frontend: "localhost:3000"

---

Security

- bcrypt password hashing
- encrypted wallet key storage
- hashed API secrets
- JWT middleware protection
- Redis rate limiting
- DB transaction safety for financial operations

---

Future Improvements

- event-driven transaction processing
- websocket wallet updates
- observability with Prometheus/Grafana
- multi-chain wallet support
- Kubernetes deployment

---

License

Proprietary — All rights reserved.



**Last Updated:** May 2026
**Developed by:** Olusanya Timothy (teamthy)


