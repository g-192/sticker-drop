# StickerDrop - Promo API

A showcase project for Golang with Postgres and Redis.

## Project Overview
This is an example project to showcase the concept of a Giveaway integrated into your platform. When launching a promo (for example: "Free Stickers"), the backend needs to handle heavy, concurrent traffic spikes without overselling (race conditions) or allowing users to claim multiple items.

![Sticker Drop Card](image.png)

## Tech Stack
* **Backend:** Go (Golang)
* **API:** GraphQL (Schema-First using `gqlgen`)
* **Database:** PostgreSQL (via `pgxpool`)
* **Caching & Rate Limiting:** Redis
* **Frontend:** React, TypeScript, Vite, Apollo Client 4
* **Infrastructure:** Docker & Docker-Compose

## Architecture

### 1. Redis (Spam Protection)
To prevent the same user from claiming a sticker multiple times by spamming the frontend button, I implemented a fast in-memory check using Redis `SetNX`. It locks the user's email instantly, shielding the Postgres database from unnecessary load and stopping double-claims in milliseconds.

### 2. Atomic Database Updates
Instead of loading the current claimed count into Go memory and calculating `+1` (which causes race conditions during high concurrency), the SQL query directly handles the constraint:
`UPDATE drops SET claimed = claimed + 1 WHERE id = $1 AND claimed < total_available;`
Postgres perfectly queues these updates, guaranteeing zero overselling even if 5,000 requests hit the server at the exact same time.

### 3. Dependency Injection & Type Safety
* The Go backend uses Dependency Injection to pass the `pgxpool` and `redis.Client` down to the GraphQL resolvers, avoiding global states.
* The frontend uses strict TypeScript interfaces (`DropData`) alongside Apollo Client's `useQuery` and type-guards to prevent runtime exceptions if the API returns undefined data.

## How to run this project locally

**1. Start the Databases (Docker)**
```bash
docker compose up -d
```
*(This starts Postgres and Redis in the background).*

**2. Start the Go Backend (GraphQL API)**
```bash
go run main.go
```
*(Runs on http://localhost:8080. The first run automatically provisions the DB table).*

**3. Start the React Frontend**
Open a new terminal window:
```bash
cd frontend
npm install
npm run dev
```
*(Runs on http://localhost:5173).*