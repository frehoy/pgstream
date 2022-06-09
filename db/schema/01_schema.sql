DROP SCHEMA IF EXISTS api CASCADE;
CREATE SCHEMA api;

CREATE TABLE api.events (
    id BIGINT GENERATED ALWAYS AS IDENTITY,
    ts TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    message JSONB
) PARTITION BY RANGE (ts);

CREATE VIEW api.events_per_second AS
SELECT
    date_trunc('second', ts) AS event_ts_second,
    COUNT(*)
FROM api.events
GROUP BY 1
ORDER BY 1 DESC
LIMIT 20;

CREATE INDEX events_created_at ON api.events USING BRIN (ts);

-- anon
DROP ROLE IF EXISTS anon;
CREATE ROLE anon nologin;
GRANT usage ON SCHEMA api TO anon;
GRANT SELECT ON api.events TO anon;

-- authenticator
DROP ROLE IF EXISTS authenticator;
CREATE ROLE authenticator noinherit LOGIN PASSWORD 'authenticator_password';
GRANT anon TO authenticator;

-- api_user
DROP ROLE IF EXISTS api_user;
CREATE ROLE api_user nologin;
GRANT api_user TO authenticator;
GRANT usage ON SCHEMA api TO api_user;
GRANT ALL ON api.events TO api_user;

