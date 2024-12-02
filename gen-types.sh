#!/bin/bash

HERE="$(dirname "$(realpath "$0")")"
OUT="$HERE/go/pkg/models/db"
CONNECTION_STR="postgres://postgres:postgres@127.0.0.1:54322/postgres?sslmode=disable"

# create the output directory if it does not exist
mkdir -p "$OUT"

xo schema "$CONNECTION_STR" \
  --schema=public \
  --out="$OUT" \
  --exclude="environment_variables_after_actions" \
  --exclude="environment_variables_before_actions" \
  --exclude="environments_after_actions" \
  --exclude="users_after_actions" \
  --exclude="variables_after_actions"