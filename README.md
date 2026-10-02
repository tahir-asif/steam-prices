# Steam Price Monitor

A full‑stack web application that tracks and visualizes historical Steam game prices. Search for any Steam game, view its price trend over time to see the best time to buy.

**Live Demo:** [steam-prices.vercel.app](https://steam-prices.vercel.app)  

> Note: This used to be hosted on AWS, but its now been moved to Vercel, Render & Neon due to cost reasons.

## Features

- 🔍 **Instant Search** – Debounced search against the Steam store with dropdown results.
- 📈 **Interactive Price History** – Recharts line chart showing price changes over time, with tooltips and time‑range filtering.
- 🗄️ **On‑Demand Tracking** – Viewing a new game automatically adds it to the database and starts tracking its price.
- ⚡ **Automated Price Sync** – GitHub Actions scheduled workflow runs hourly, fetches current Steam prices and stores only changed records.
- 🚀 **Seeded Popular Games** – Database pre‑populated with the 100 most popular Steam games via the Steam Spy API.
- 🧪 **Comprehensive Testing** – Integration tests for the backend (testcontainers‑go) and component tests for the frontend (Vitest + React Testing Library).
- 🚢 **CI/CD & Auto‑Deploy** – GitHub Actions runs tests; Vercel and Render auto‑deploy on push.
- 🔒 **Secure by Design** – No direct internet access to the API or database; all traffic goes through CloudFront with locked‑down security groups.
- 🐳 **Dockerized Backend** – Go backend containerized with Docker for consistent local and production runs.

## Tech Stack

| Layer | Current | Legacy (AWS) |
|:------|:--------------------------------|:-------------|
| **Frontend** | Vercel (static hosting) | S3 + CloudFront CDN |
| **Backend** | Render Web Service (Go) | EC2 with Docker, behind CloudFront VPC Origin |
| **Database** | Neon (serverless PostgreSQL) | Amazon RDS for PostgreSQL |
| **Worker** | GitHub Actions scheduled workflow | AWS Lambda (Go) + EventBridge cron |
| **Container Registry** | – (direct Render build) | Amazon ECR |
| **CI/CD** | Vercel & Render auto‑deploy on push | GitHub Actions |
| **Testing** | testcontainers‑go, Vitest, React Testing Library | – (same as current) |

## Local Development (Quick Start)

### Prerequisites
- Git
- Go
- Docker (+ Docker Desktop on MacOS or Windows for docker compose)
- Node.JS v and npm

```bash
# Clone the repository
git clone https://github.com/tahir-asif/steam-prices.git
cd steam-prices
```

```bash
# Start PostgreSQL and Adminer (Docker)
docker compose up -d
```

```bash
# Backend
cd backend
go mod download
go run cmd/api/main.go   # API server on :3000
```

```bash

# Frontend
cd frontend
npm install
npm run dev              # Dev server on :5173
```

The app will be available at http://localhost:5173.
You may copy `.env.example` to `.env` and adjust `DATABASE_URL` if needed; the default fallback works with the Docker Compose PostgreSQL.

## Future Improvements
- User accounts and wishlist tracking (with Steam OAuth)
- Email / push notifications on price drops

## Licence
[LICENSE](LICENSE)
