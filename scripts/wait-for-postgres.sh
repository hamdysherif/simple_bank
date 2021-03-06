#!/bin/sh
# wait-for-postgres.sh

set -e

echo "starting postgres"
until psql -h postgres -d simple_bank -U root -c "\q"; do
  >&2 echo "Postgres is unavailable - sleeping"
  sleep 1
done

echo "Postgres is up - executing command"

exec "$@"
