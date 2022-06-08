FROM postgres:14
RUN apt update && apt install postgresql-14-cron
ADD postgresql.conf /etc/postgresql/postgresql.conf
ADD schema/* /docker-entrypoint-initdb.d/
