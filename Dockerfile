# All-in-one image: Go backend + Next.js frontend + SQLite, single port 8080.
#
# Designed for one-click deploy on Railway / Render / Fly / a VPS / DockerHub.
# No extra services required (no Postgres, no Redis, no fetcher).
#
# The Go backend listens on :8080 (public). It serves /api/*, /healthz,
# /robots.txt, /sitemap.xml itself and PROXIES every other path to the
# Next.js process running on 127.0.0.1:3000 inside the same container.
# Database is a SQLite file under /data which Railway mounts as a volume.
#
# Build:   docker build -t videoxnx .
# Run:     docker run -p 8080:8080 -v videoxnx-data:/data videoxnx

# syntax=docker/dockerfile:1.7

############################
# 1. Build the Go backend  #
############################
FROM golang:1.22-alpine AS gobuild
WORKDIR /src
RUN apk add --no-cache git ca-certificates
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/server ./cmd/server

#####################################
# 2. Build the Next.js (standalone) #
#####################################
FROM node:20-alpine AS webbuild
WORKDIR /app
COPY frontend/package.json frontend/package-lock.json* ./
RUN npm install --no-audit --no-fund
COPY frontend/ ./
ENV NEXT_TELEMETRY_DISABLED=1 \
    NEXT_PUBLIC_API_URL="" \
    NEXT_PUBLIC_SITE_URL="" \
    NEXT_PUBLIC_DMCA_EMAIL=""
RUN npm run build

##########################
# 3. Runtime (alpine+node)
##########################
FROM node:20-alpine AS run
RUN apk add --no-cache ca-certificates tzdata curl \
 && addgroup -S app && adduser -S -G app -u 10001 app
WORKDIR /app

# Backend binary
COPY --from=gobuild /out/server /app/server

# Next.js standalone output (server.js + minimal node_modules)
COPY --from=webbuild /app/.next/standalone/ ./
COPY --from=webbuild /app/.next/static       ./.next/static
COPY --from=webbuild /app/public             ./public

# Start script
COPY scripts/start.sh /app/start.sh
RUN chmod +x /app/start.sh \
 && mkdir -p /data \
 && chown -R app:app /app /data

USER app
EXPOSE 8080

ENV NODE_ENV=production \
    NEXT_TELEMETRY_DISABLED=1 \
    PORT=8080 \
    FRONTEND_PORT=3000 \
    FRONTEND_URL=http://127.0.0.1:3000 \
    INTERNAL_API_URL=http://127.0.0.1:8080 \
    DATABASE_URL=sqlite:///data/videoxnx.db \
    APP_ENV=production \
    SCRAPER_ENABLED=true \
    SCRAPER_FETCHER_URL="" \
    SCRAPER_MAX_PAGES=200 \
    AGE_GATE_ENABLED=true

HEALTHCHECK --interval=30s --timeout=5s --start-period=20s --retries=3 \
  CMD curl -fsS http://127.0.0.1:8080/healthz || exit 1

CMD ["/app/start.sh"]

