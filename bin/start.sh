set -e
psql -v ON_ERROR_STOP=1 --host=127.0.0.1 --file=schema.sql
PGRST_SERVER_PORT=3000 DYLD_LIBRARY_PATH=/Applications/Postgres.app/Contents/Versions/14/lib bin/postgrest postgrest.conf 
