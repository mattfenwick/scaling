#!/usr/bin/env sh

set -xv
set -eu

# Variables needed:
# PGHOST
# PGPASSWORD
# PGUSER
# SCALING_DATABASE


pg_isready

BASEDIR=${BASEDIR:-$(dirname "$0")}

psql --command "create database \"$SCALING_DATABASE\" encoding UTF8" || true

PGDATABASE="$SCALING_DATABASE" \
  psql --file "$BASEDIR"/schema.sql
