DROP SCHEMA IF EXISTS api CASCADE;
CREATE SCHEMA api;

CREATE TABLE api.events(
	id integer primary key generated always as identity,
	created_at TIMESTAMP WITH TIME ZONE default current_timestamp,
	message TEXT
);

DROP ROLE IF EXISTS web_anon;
create role web_anon nologin;

grant usage on schema api to web_anon;
grant select on api.events to web_anon;

DROP ROLE IF EXISTS authenticator;
create role authenticator noinherit login password 'pw';
grant web_anon to authenticator;

DROP ROLE IF EXISTS api_user;
create role api_user nologin;
grant api_user to authenticator;

grant usage on schema api to api_user;
grant all on api.events to api_user;