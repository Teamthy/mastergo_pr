# MasterGo

A production-grade API platform for cryptocurrency payments with secure API key management, wallet functionality, and real-time webhook delivery.

**Live Demo:** [https://mastergo-pr-1.onrender.com](https://mastergo-pr-1.onrender.com/)


## Overview

MasterGo is a comprehensive solution for developers and businesses looking to integrate cryptocurrency payments into their applications. The platform provides a secure, auditable, and developer-friendly API for managing crypto transactions, user authentication, and automated webhook notifications.

### Key Features

- **Secure Authentication** — JWT-based authentication with email verification and two-factor authentication
- **API Key Management** — Industry-standard API key system with key rotation and revocation
- **Wallet Management** — Native cryptocurrency wallet creation and transaction handling
- **Webhook System** — Reliable event delivery with exponential backoff retry strategy
- **Audit Logging** — Complete audit trail of all user actions and API operations
- **Analytics Dashboard** — Real-time insights into API usage and transaction metrics
- **Rate Limiting** — Built-in rate limiting to prevent abuse and ensure fair usage


## Tech Stack

### Backend
- **Language:** Go 1.25
- **Framework:** Chi Router (lightweight HTTP router)
- **Database:** PostgreSQL (data persistence)
- **Cache:** Redis (session management and rate limiting)
- **Blockchain:** Ethereum via go-ethereum library
- **Email:** SendGrid for transactional emails
- **Authentication:** JWT tokens

### Frontend
- **Framework:** Next.js 16
- **Language:** TypeScript
- **Styling:** Tailwind CSS
- **UI Library:** Lucide React for icons
- **State Management:** XState for complex UI state
- **Web3:** ethers.js for blockchain interaction

### Infrastructure
- **Containerization:** Docker & Docker Compose
- **Deployment:** Render (cloud platform)



## Project Structure

```
.
├── backend/                     # Go API server
│   ├── cmd/api/                 # Application entry point
│   ├── internal/
│   │   ├── config/              # Configuration management
│   │   ├── handler/             # HTTP request handlers
│   │   ├── middleware/          # Custom middleware (auth, rate limiting)
│   │   ├── service/             # Business logic layer
│   │   ├── repository/          # Data access layer
│   │   ├── models/              # Data structures
│   │   ├── database/            # DB connections
│   │   ├── crypto/              # Encryption utilities
│   │   ├── routes/              # API route definitions
│   │   └── utils/               # Helper functions
│   ├── migrations/              # Database schema
│   ├── Dockerfile
│   ├── go.mod & go.sum
│   └── main.go
├── frontend/                    # Next.js web application
│   ├── app/                     # Next.js app router
│   │   ├── auth/                # Authentication pages
│   │   ├── dashboard/           # User dashboard
│   │   └── api/                 # API routes (proxy)
│   ├── components/              # Reusable React components
│   ├── lib/                     # Utility functions and contexts
│   ├── styles/                  # Global styles
│   ├── public/                  # Static assets
│   ├── Dockerfile
│   ├── package.json
│   └── tsconfig.json
├── docker-compose.yml           # Multi-container setup
└── migrations/                  # Database initialization scripts
```


## Getting Started

### Prerequisites

- Docker & Docker Compose (recommended)
- Go 1.25+ (for backend development)
- Node.js 18+ (for frontend development)
- PostgreSQL 14+ (if running locally)
- Redis 7+ (if running locally)

### Quick Start with Docker

```bash
# Clone the repository
git clone https://github.com/mastergo/mastergo-pr.git
cd mastergo-pr

# Start all services
docker-compose up -d

# Backend runs on http://localhost:8080
# Frontend runs on http://localhost:3000
```

### Local Development Setup

#### Backend

```bash
cd backend

# Install dependencies
go mod download

# Set up environment variables
cp .env.example .env
# Edit .env with your configuration

# Run database migrations
# (Handled automatically on startup)

# Start the server
go run cmd/api/main.go
```

#### Frontend

```bash
cd frontend

# Install dependencies
npm install

# Create environment file
cp .env.example .env
# Edit .env with API URL

# Start development server
npm run dev
```



## API Documentation

### Authentication

All API requests to protected endpoints require a valid JWT token in the `Authorization` header.

```
Authorization: Bearer <your_jwt_token>
```

### Core Endpoints

#### Authentication
- `POST /auth/signup` — Create a new account
- `POST /auth/login` — Authenticate and receive JWT
- `POST /auth/verify-email` — Verify email with OTP
- `POST /auth/resend-otp` — Request new verification code
- `GET /auth/password-strength` — Check password strength
- `GET /auth/email-available` — Validate email availability

#### Wallet
- `POST /api/v1/wallet/create` — Create a new wallet
- `GET /api/v1/wallet/balance` — Get wallet balance
- `GET /api/v1/wallet/transactions` — List transactions
- `POST /api/v1/wallet/withdraw` — Initiate withdrawal

#### API Keys
- `POST /api/v1/apikeys` — Generate new API key
- `GET /api/v1/apikeys` — List your API keys
- `POST /api/v1/apikeys/{id}/regenerate` — Rotate API key
- `DELETE /api/v1/apikeys/{id}` — Revoke API key

### Response Format

All responses are JSON:

```json
{
  "success": true,
  "data": { },
  "error": null
}
```



## Security Features

### Authentication & Authorization
- JWT tokens with configurable expiration
- Email verification for new accounts
- Two-factor authentication support
- Secure password hashing with bcrypt

### API Security
- API key authentication for programmatic access
- CORS protection with configurable allowed origins
- Rate limiting per IP and API key
- Request validation and sanitization

### Data Protection
- Master key encryption for sensitive data
- Audit logging for compliance
- Secure wallet key storage
- Password reset tokens with TTL

### Infrastructure
- HTTPS/TLS for all connections
- Environment variable configuration for secrets
- Docker containerization for isolation
- Database connection pooling




## Database Schema

The system uses PostgreSQL with the following main tables:

- **users** — User accounts and authentication
- **wallets** — Cryptocurrency wallets and balances
- **api_keys** — API credentials and permissions
- **transactions** — Transaction history and status
- **webhook_events** — Outbound webhook events and delivery status
- **audit_logs** — Complete audit trail of system events



## Deployment

### Production Deployment

The project is configured for deployment on Render:

```bash
# Build
docker build -t mastergo-backend ./backend
docker build -t mastergo-frontend ./frontend

# Push to container registry
# Configure on Render dashboard
```

### Environment Configuration

Set these environment variables in your deployment platform:
- `PORT` — API server port
- `DBURL` — PostgreSQL connection string
- `REDIS_URL` — Redis connection string
- `JWT_SECRET` — JWT signing secret
- `MASTER_KEY` — Master encryption key
- `ETH_RPC_URL` — Ethereum RPC endpoint

---

## Development Workflow

### Running Tests

```bash
cd backend
go test ./...

cd ../frontend
npm run lint
```

### Code Quality

- Backend: Follow Go conventions, use `gofmt` and `golint`
- Frontend: ESLint configuration provided in `eslint.config.mjs`

### Database Migrations

Migrations are located in `backend/migrations/` and run automatically on application startup.

---

## Performance & Scaling

- **Caching** — Redis for session and rate limit data
- **Connection Pooling** — PostgreSQL with pgx for efficient connections
- **Rate Limiting** — Token bucket algorithm
- **Async Processing** — Goroutines for webhook delivery and transaction watching

---



## Support & Contributing

For issues and questions:
- Review the project documentation in the root directory
- Check existing API documentation files
- Submit issues through the repository



## License

This project is proprietary. All rights reserved.


## Deployment Links

- **Production:** https://mastergo-pr-1.onrender.com
- **API Documentation:** Available at `/health` endpoint



**Last Updated:** May 2026
**Developed by:** Olusanya Timothy (teamthy)


