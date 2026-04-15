#!/usr/bin/env bash
# Regenerate db/sqlc/ from db/queries/*.sql and db/sqitch/deploy/*.sql
# Run this after modifying any SQL file.
set -e

command -v sqlc >/dev/null 2>&1 || {
    echo "sqlc not found — install it: https://docs.sqlc.dev/en/latest/overview/install.html"
    exit 1
}

echo "Generating SQLC code …"
sqlc generate
echo "Done — db/sqlc/ updated."
