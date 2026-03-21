#!/bin/sh
set -e

host="$1"
port="$2"
shift 2
cmd="$@"

until pg_isready -h "$host" -p "$port" -q 2>/dev/null; do
  echo "⏳ Waiting for database at $host:$port..."
  sleep 2
done

echo "✅ Database is ready!"
exec $cmd
