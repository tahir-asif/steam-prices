# Steam Price Monitor

A full‑stack web application that tracks and visualizes historical Steam game prices. Search for any Steam game, view its price trend over time, and let the backend keep the data fresh every hour—now fully hosted on AWS.

**Live Demo (AWS):** [d35qn74h6do2xm.cloudfront.net](https://d35qn74h6do2xm.cloudfront.net)

> **Legacy Demo (Vercel / Render / Neon):** [steam-prices.vercel.app](https://steam-prices.vercel.app)  
> *This older deployment remains available temporarily for reference, but will be retired soon.*

## Features

- 🔍 **Instant Search** – Debounced search against the Steam store with dropdown results.
- 📈 **Price History Charts** – Interactive line chart (Recharts) showing price changes over time.
- ⚡ **Serverless Price Worker** – An AWS Lambda function, triggered by EventBridge every hour, fetches current Steam prices and stores only changed records.
- 🗄️ **On‑Demand Tracking** – Viewing a new game automatically adds it to the database and starts tracking.
- 🚀 **Seeded Initial Data** – Database pre‑populated with the 100 most popular Steam games (via Steam Spy API).
- 🔒 **Secure by Design** – No direct internet access to the API or database; all traffic goes through CloudFront with locked‑down security groups.
- 🧪 **Tested** – Integration tests for the backend (testcontainers‑go) and component tests for the frontend (Vitest).

## Tech Stack

| Layer | Current | Legacy |
|:------|:--------------|:--------------------------------|
| **Frontend** | S3 + CloudFront CDN | Vercel (static hosting) |
| **Backend** | EC2 with Docker, behind CloudFront VPC Origin | Render Web Service (Go) |
| **Database** | Amazon RDS for PostgreSQL | Neon (serverless PostgreSQL) |
| **Worker** | AWS Lambda (Go) + EventBridge cron | GitHub Actions scheduled workflow |
| **Container Registry** | Amazon ECR | – (direct Render build) |
| **CI/CD** | GitHub Actions | Vercel & Render auto‑deploy on push |
| **Testing** | testcontainers‑go, Vitest, React Testing Library | – (same as current) |

## Local Development (Quick Start)

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
- Retrieving historic price data
- Improved landing page with recent searches
- CloudFront caching policies for better performance
