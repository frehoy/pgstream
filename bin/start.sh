#!/usr/bin/env sh

set -e

psql -v ON_ERROR_STOP=1 --host=127.0.0.1 --file=schema.sql
bin/postgrest postgrest.conf
