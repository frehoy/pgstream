DROP SCHEMA IF EXISTS api CASCADE;
CREATE SCHEMA api;

CREATE TABLE api.events (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    message TEXT
);

CREATE INDEX events_created_at ON api.events USING BRIN (created_at);

-- web_anon
DROP ROLE IF EXISTS web_anon;
CREATE ROLE web_anon nologin;
GRANT usage ON SCHEMA api TO web_anon;
GRANT SELECT ON api.events TO web_anon;

-- authenticator
DROP ROLE IF EXISTS authenticator;
CREATE ROLE authenticator noinherit LOGIN PASSWORD 'pw';
GRANT web_anon TO authenticator;

-- api_user
DROP ROLE IF EXISTS api_user;
CREATE ROLE api_user nologin;
GRANT api_user TO authenticator;
GRANT usage ON SCHEMA api TO api_user;
GRANT ALL ON api.events TO api_user;

