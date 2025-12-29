#!/bin/bash
# Migration helper script

set -e

MIGRATIONS_DIR="./migrations"
DATABASE_URL="${DATABASE_URL:-postgresql://postgres:dev@localhost:5432/content_analyzer?sslmode=disable}"

case "$1" in
  up)
    echo "Running migrations up..."
    migrate -path $MIGRATIONS_DIR -database "$DATABASE_URL" up
    ;;
  down)
    echo "Running migrations down..."
    migrate -path $MIGRATIONS_DIR -database "$DATABASE_URL" down
    ;;
  force)
    echo "Forcing version $2..."
    migrate -path $MIGRATIONS_DIR -database "$DATABASE_URL" force "$2"
    ;;
  version)
    echo "Current migration version:"
    migrate -path $MIGRATIONS_DIR -database "$DATABASE_URL" version
    ;;
  create)
    echo "Creating new migration: $2"
    migrate create -ext sql -dir $MIGRATIONS_DIR -seq "$2"
    ;;
  *)
    echo "Usage: $0 {up|down|force|version|create} [args]"
    exit 1
    ;;
esac
