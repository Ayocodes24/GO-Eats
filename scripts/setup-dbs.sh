#!/usr/bin/env bash
# Creates the three separate databases — one per microservice.
# Run once before starting services for the first time.
set -e

PG="PGPASSWORD=yourpassword psql -U postgres -h localhost -p 5432"

for DB in goeats_user goeats_restaurant goeats_order; do
  if $PG -lqt | cut -d'|' -f1 | grep -qw "$DB"; then
    echo "DB '$DB' already exists — skipping"
  else
    $PG -c "CREATE DATABASE $DB;"
    echo "Created database: $DB"
  fi
done

echo "All databases ready."
