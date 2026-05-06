# Root-level Dockerfile so Railway / any platform that does NOT have a
# Root Directory configured can still deploy the project.
#
# This builds and runs the BACKEND service. For the frontend and fetcher,
# create separate Railway services with Root Directory set to "frontend"
# and "fetcher" respectively.
#
# We just delegate to backend/Dockerfile by re-implementing the same steps
# with the build context narrowed to ./backend.

# syntax=docker/dockerfile:1.7
FROM golang:1.22-alpine AS build
WORKDIR /src
RUN apk add --no-cache git ca-certificates
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/server ./cmd/server

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata && adduser -D -u 10001 app
WORKDIR /app
COPY --from=build /out/server /app/server
RUN mkdir -p /app/data && chown -R app:app /app/data
USER app
EXPOSE 8080
ENV PORT=8080 \
    DATA_DIR=/app/data \
    DATABASE_URL=sqlite:///app/data/videoxnx.db
CMD ["/app/server"]
