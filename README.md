# Steam Price Monitor

A full‑stack web application that tracks and visualizes historical Steam game prices. Search for any Steam game, view its price trend over time.

**Live Demo:** [steam-prices.vercel.app](https://steam-prices.vercel.app)

## Features

- **Instant Search** – Debounced search against the Steam store.
- **Price History Charts** – Interactive line chart displaying price changes over time.
- **Hourly Updates** – Background worker fetches current prices via GitHub Actions, recording only changed prices.
- **On‑Demand Tracking** – Viewing a new game automatically adds it to the database and starts tracking.
- **Tested** – Integration tests for the backend and component tests for the frontend.

## Tech Stack

| Layer | Technologies |
|:------|:-------------|
| **Frontend** | React, TypeScript, Vite, React Router, Recharts |
| **Backend** | Go, Gin, Steam API client |
| **Database** | PostgreSQL, versioned migrations with golang‑migrate |
| **DevOps** | Docker, GitHub Actions, Vercel (frontend), Render (API), Neon (DB) |
| **Testing** | testcontainers‑go, Vitest, React Testing Library |

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
