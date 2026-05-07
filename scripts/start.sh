#!/bin/sh
# Start Next.js (background) then the Go backend (foreground).
# Backend proxies non-/api requests to Next.js on 127.0.0.1:$FRONTEND_PORT.
set -e

: "${FRONTEND_PORT:=3000}"
: "${PORT:=8080}"

echo "[start] launching Next.js on 127.0.0.1:${FRONTEND_PORT}"
HOSTNAME=127.0.0.1 PORT="${FRONTEND_PORT}" node /app/server.js &
WEB_PID=$!

# Forward signals so Docker stop is graceful
trap 'echo "[start] stopping..."; kill -TERM "$WEB_PID" 2>/dev/null; kill -TERM "$API_PID" 2>/dev/null; wait' INT TERM

echo "[start] launching Go backend on :${PORT}"
/app/server &
API_PID=$!

wait -n "$WEB_PID" "$API_PID"
EXIT=$?
echo "[start] a child exited with code ${EXIT}, shutting down siblings"
kill -TERM "$WEB_PID" "$API_PID" 2>/dev/null || true
wait
exit "$EXIT"
