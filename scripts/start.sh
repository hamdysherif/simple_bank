#!/bin/sh

set -e

echo "running migration"
source /app/app.env
/app/migrate -path=/app/db/migrate -database="$DB_SOURCE" -verbose up

exec "$@"
