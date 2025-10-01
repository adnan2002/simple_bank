#!/bin/sh
set -e

echo "Waiting for database to be ready..."
/app/wait-for.sh db:5432 -t 60 -- echo "Database is ready!"

echo "Running database migrations..."
source /app/app.env
make migrateup

echo "Starting application..."
exec /app/main