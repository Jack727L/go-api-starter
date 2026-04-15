#!/usr/bin/env bash
# Regenerate Swagger docs from handler annotations.
# Run this before committing when handlers change.
set -e

command -v swag >/dev/null 2>&1 || {
    echo "swag not found — install it: go install github.com/swaggo/swag/cmd/swag@latest"
    exit 1
}

echo "Generating Swagger docs …"
swag init --parseDependency --parseInternal
echo "Done — docs/ updated."
