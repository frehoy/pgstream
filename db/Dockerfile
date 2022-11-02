FROM postgres:14
RUN apt update && apt install --yes postgresql-14-cron && rm -rf /var/lib/apt/lists/*
ADD postgresql.conf /etc/postgresql/postgresql.conf
ADD schema/* /docker-entrypoint-initdb.d/
