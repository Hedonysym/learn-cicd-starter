#!/usr/bin/env bash
set -euo pipefail
goose -dir ./sql/schema postgres "$DATABASE_URL" up