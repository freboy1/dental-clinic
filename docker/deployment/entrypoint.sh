#!/bin/sh

set -e  # Exit immediately if a command exits with a non-zero status

echo "Running migrations..."
goose -dir /bin/migrations postgres "$DB_DSN" up

echo "Starting Go app..."
exec /bin/app
