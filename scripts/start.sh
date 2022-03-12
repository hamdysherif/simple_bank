#!/bin/sh

set -e

echo "running migration"

/app/migrate -path=/app/db/migrate -database="$DB_SOURCE" -verbose up

exec "$@"
