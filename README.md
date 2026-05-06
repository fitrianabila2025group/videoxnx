# VideoXNX — Aggregator Mirror (Full Stack)

A production-ready, full-stack content aggregator that mirrors a WordPress-style video blog
into your own database, exposes a REST API, and serves a modern Next.js frontend with an
admin panel.

> **Important — Legal & Safety**
>
> - Use this software **only** with explicit written permission from the source website operator.
> - The scraper respects `robots.txt`, rate-limits requests, rotates user-agents, and never bypasses
>   logins, paywalls, Cloudflare, or anti-bot systems.
> - Video files are **not** rehosted. Only metadata, thumbnails, descriptions, source URLs, and
>   embed iframes are stored.
> - A keyword safety filter blocks any post that hints at minors or non-consensual content.
> - The site enforces an **18+ age gate** and provides DMCA / takedown / report mechanisms.

---

## Stack

| Layer       | Tech |
|------------ |------|
| Backend     | Go 1.22 + Fiber v2 + GORM + PostgreSQL |
| Scraper     | goquery + WP REST API + `robotstxt` + cron via gocron |
| Auth        | JWT (HS256) + bcrypt + HTTP-only cookies |
| Frontend    | Next.js 14 (App Router) + Tailwind CSS |
| Cache       | Redis (optional) |
| Container   | Docker / docker-compose |

---

## Project structure

```
backend/
  cmd/server/main.go
  internal/
    config/      # env loader
    controllers/ # HTTP handlers (public + admin)
    database/    # gorm connect + AutoMigrate
    middleware/  # JWT, secure headers
    models/      # GORM entities
    routes/      # fiber router
    scraper/     # http client, parser, scheduler
    services/    # auth, posts upsert, safety filter
    utils/       # slug, jwt, sanitizer
  migrations/    # optional manual SQL (GORM handles AutoMigrate)
  Dockerfile
  go.mod

frontend/
  src/
    app/         # Next App Router pages (public + /admin)
    components/  # AgeGate, PostCard, ReportButton, ...
    lib/         # api client
    styles/      # Tailwind globals
  Dockerfile
  package.json

docker-compose.yml
.env.example
```

---

## Quick start (Docker)

1. Copy env file and edit values:

   ```bash
   cp .env.example .env
   # edit JWT_SECRET, ADMIN_EMAIL, ADMIN_PASSWORD, DMCA_EMAIL, SOURCE_BASE_URL, etc.
   ```

2. Boot everything:

   ```bash
   docker compose up -d --build
   ```

3. Open:

   - Public site: http://localhost:3000
   - API healthcheck: http://localhost:8080/healthz
   - Admin login: http://localhost:3000/admin/login

The first boot will run `AutoMigrate` and create the admin user from `ADMIN_EMAIL` /
`ADMIN_PASSWORD`. The scheduler will start a scrape pass immediately if `SCRAPER_ENABLED=true`.

---

## Local dev (without Docker)

### Backend

```bash
cd backend
cp ../.env.example ../.env   # if not done
export $(grep -v '^#' ../.env | xargs)
go run ./cmd/server
```

### Frontend

```bash
cd frontend
npm install
NEXT_PUBLIC_API_URL=http://localhost:8080 npm run dev
```

---

## Environment variables

| Variable | Purpose |
|---|---|
| `APP_ENV` | `development` or `production` |
| `APP_URL` | Public base URL (used in sitemap + cookies) |
| `PORT` | API port (default `8080`) |
| `DATABASE_URL` | Postgres DSN |
| `REDIS_URL` | Redis DSN (optional) |
| `JWT_SECRET` | HS256 secret for admin tokens |
| `ADMIN_EMAIL` / `ADMIN_PASSWORD` | Created on first boot if absent |
| `SOURCE_BASE_URL` | Source site to mirror, e.g. `https://indoxvx.cam` |
| `SCRAPER_ENABLED` | `true` to enable cron scheduler |
| `SCRAPER_INTERVAL_MINUTES` | Cron interval |
| `SCRAPER_MAX_PAGES` | Pagination depth per pass |
| `SCRAPER_RATE_LIMIT_MS` | Min delay between requests |
| `SCRAPER_RESPECT_ROBOTS` | Honor `robots.txt` (recommended `true`) |
| `SCRAPER_USER_AGENT` | Identify your bot honestly |
| `AGE_GATE_ENABLED` | Toggle 18+ gate |
| `DMCA_EMAIL` | Address shown on DMCA page |
| `CORS_ALLOWED_ORIGINS` | Comma-separated list, or `*` |

---

## Scraper behavior

The scraper performs the following steps each pass:

1. Try the **WordPress REST API** (`/wp-json/wp/v2/posts?_embed=1&page=N`). If reachable
   and returning JSON, paginate through up to `SCRAPER_MAX_PAGES` and upsert posts.
2. If the REST API is disabled or unavailable, fall back to **HTML scraping** with
   configurable CSS selectors (see `internal/scraper/selectors.go`). Listing pages are
   crawled and each post page is parsed individually.
3. Each candidate post is run through the **safety filter** (`internal/services/safety.go`):
   - `blocked` — hard-block keywords (minors, non-consensual, illegal acts) → status `blocked`.
   - `review` — ambiguous keywords → status `hidden`, flagged for admin.
   - `safe` → status `published`.
4. Posts are upserted by `source_url`. Existing posts are updated, but admin moderation
   choices are preserved.
5. A row is appended to `scrape_logs` with status, counts, duration, and any errors.
6. Errors during one page never crash the pass.

Manual scrape: `Admin → Scraper → Run scraper now` (also `POST /api/admin/scraper/run`).

### Customizing selectors

Open `backend/internal/scraper/selectors.go` and adjust CSS selectors for the source theme.
The defaults cover most WordPress video themes.

---

## API

### Public

```
GET  /api/posts?page=&per_page=
GET  /api/posts/:slug
GET  /api/categories
GET  /api/categories/:slug/posts
GET  /api/tags
GET  /api/tags/:slug/posts
GET  /api/search?q=
GET  /api/trending
POST /api/reports
GET  /api/site
GET  /robots.txt
GET  /sitemap.xml
```

### Admin (`Authorization: Bearer <token>` or HTTP-only cookie)

```
POST   /api/admin/login
POST   /api/admin/logout
GET    /api/admin/dashboard
GET    /api/admin/posts?status=&safety_status=&q=
PUT    /api/admin/posts/:id
DELETE /api/admin/posts/:id
PATCH  /api/admin/posts/:id/status

GET    /api/admin/categories
POST   /api/admin/categories
PUT    /api/admin/categories/:id
DELETE /api/admin/categories/:id

GET    /api/admin/tags
POST   /api/admin/tags
PUT    /api/admin/tags/:id
DELETE /api/admin/tags/:id

POST   /api/admin/scraper/run
GET    /api/admin/scraper/logs

GET    /api/admin/reports
PATCH  /api/admin/reports/:id

GET    /api/admin/settings
PUT    /api/admin/settings
```

---

## Security

- Bcrypt for password hashing.
- JWT (HS256) tokens, 12h expiry, available as `Authorization: Bearer` or HTTP-only cookie.
- Rate limiting middleware (120 req/min/IP) and Fiber's `recover`.
- Secure response headers.
- All scraped HTML is run through `bluemonday` (allowlisted iframe/video tags) before storage.
- Iframes render with `sandbox` + `referrerPolicy=no-referrer`.
- All DB access via GORM (parameterized).
- Secrets read from environment — none hardcoded.

---

## Deployment

### Railway

Each service has its own `railway.json` so Railway will auto-detect the Dockerfile.

1. Push this repo to GitHub.
2. In Railway create a new project from the repo and add **Postgres** plugin.
3. Add three services from the same repo with these **Root Directory** settings:
   - `backend/` — set vars from `.env.example` plus `DATABASE_URL` (from Postgres plugin),
     `JWT_SECRET`, `ADMIN_EMAIL`, `ADMIN_PASSWORD`, `SCRAPER_FETCHER_URL=http://${{fetcher.RAILWAY_PRIVATE_DOMAIN}}:9090`,
     `CORS_ALLOWED_ORIGINS=https://<frontend-domain>`. Generate a public domain.
   - `fetcher/` — no extra env. Internal port 9090. (Don't expose publicly.)
   - `frontend/` — set `NEXT_PUBLIC_API_URL=https://<backend-domain>` as a **Build & Deploy variable**.
     Generate a public domain.
4. Health checks: backend `/healthz`, fetcher `/healthz`.
5. The scraper auto-runs on boot and every `SCRAPER_INTERVAL_MINUTES` (default 60).

### VPS (single host)

```bash
git clone <this-repo> /opt/videoxnx
cd /opt/videoxnx
cp .env.example .env  # then edit
docker compose up -d --build
```

Put Caddy in front for TLS:

```caddyfile
example.com { reverse_proxy localhost:3000 }
api.example.com { reverse_proxy localhost:8080 }
```

### EasyPanel

Use the same Dockerfiles. Add Postgres + Redis services from EasyPanel templates and
inject the env vars.

---

## Operations

- **Default admin** is created on first boot from `ADMIN_EMAIL` / `ADMIN_PASSWORD`.
- **Migrations** run automatically via `AutoMigrate` on boot.
- **Logs**: `docker compose logs -f backend`.
- **Manual scrape**: from the Admin → Scraper page or via the API.

---

## Disclaimer

This codebase is provided as a technical scaffold. You are responsible for:

- Obtaining and maintaining permission from the source operator.
- Complying with all applicable laws (DMCA, GDPR, age-verification, jurisdictional rules
  for adult content).
- Operating the safety filter and reviewing flagged content promptly.