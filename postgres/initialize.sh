#!/usr/bin/env sh

set -xv
set -eu

# Variables needed:
# PGHOST
# PGPASSWORD
# PGUSER
# SCALING_DATABASE


pg_isready

psql --command "create database \"$SCALING_DATABASE\" encoding UTF8"

PGDATABASE="$SCALING_DATABASE" \
  psql --file ./schema.sql
