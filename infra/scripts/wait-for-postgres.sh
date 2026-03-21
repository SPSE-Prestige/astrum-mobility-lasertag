#!/bin/sh
set -e

HOST="${1:-postgres}"
PORT="${2:-5432}"
MAX_RETRIES="${3:-30}"
RETRY_INTERVAL="${4:-2}"

echo "Waiting for PostgreSQL at ${HOST}:${PORT}..."

retries=0
until pg_isready -h "$HOST" -p "$PORT" -q 2>/dev/null; do
  retries=$((retries + 1))
  if [ "$retries" -ge "$MAX_RETRIES" ]; then
    echo "ERROR: PostgreSQL not ready after ${MAX_RETRIES} attempts. Exiting."
    exit 1
  fi
  echo "Attempt ${retries}/${MAX_RETRIES} - PostgreSQL not ready, retrying in ${RETRY_INTERVAL}s..."
  sleep "$RETRY_INTERVAL"
done

echo "PostgreSQL is ready."
