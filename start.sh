#!/bin/sh

set -e

echo "Run database migrations"
. /app/app.env
/app/migrate -path /app/migrations -database "$DB_SOURCE" -verbose up

echo "Start the application"
exec "$@"