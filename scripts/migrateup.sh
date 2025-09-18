#!/usr/bin/env bash
set -euo pipefail

if [ -z "${DATABASE_URL:-}" ]; then
  echo "DATABASE_URL is not set"
  exit 1
fi

"$(go env GOPATH)/bin/goose" -dir ./sql/schema sqlite "$DATABASE_URL" up