DROP SCHEMA IF EXISTS queue CASCADE;
CREATE SCHEMA queue;
GRANT usage ON SCHEMA queue TO api_user;
ALTER DEFAULT PRIVILEGES 
FOR ROLE track
IN SCHEMA queue
GRANT SELECT, INSERT, UPDATE ON TABLES
TO api_user;

CREATE TYPE queue.job_status AS ENUM(
  'waiting',
  'in_progress',
  'finished',
  'failed'
);

CREATE TABLE queue.jobs(
  id BIGINT GENERATED ALWAYS AS IDENTITY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  scheduled_for TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  started_at TIMESTAMP WITH TIME ZONE,
  finished_at TIMESTAMP WITH TIME ZONE,
  status queue.job_status DEFAULT 'waiting',
  work JSONB,
  worker_id integer
)
PARTITION BY LIST (status);
CREATE TABLE queue.jobs_waiting 
PARTITION OF queue.jobs FOR VALUES IN ('waiting');
CREATE TABLE queue.jobs_in_progress 
PARTITION OF queue.jobs FOR VALUES IN ('in_progress');
CREATE TABLE queue.jobs_finished 
PARTITION OF queue.jobs FOR VALUES IN ('finished');
CREATE TABLE queue.jobs_failed 
PARTITION OF queue.jobs FOR VALUES IN ('failed');

CREATE VIEW api.jobs AS SELECT * FROM queue.jobs;
GRANT ALL ON api.jobs TO api_user;

CREATE OR REPLACE FUNCTION notify_new_job()
  RETURNS trigger AS $$
BEGIN
  PERFORM pg_notify('queue_jobs', 'new job');
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER notify_new_job
  AFTER INSERT ON queue.jobs
  EXECUTE PROCEDURE notify_new_job();

CREATE OR REPLACE FUNCTION set_job_finished_at()
  RETURNS trigger AS $FUNC$
BEGIN
    IF new.status = 'finished' THEN
      new.finished_at = CURRENT_TIMESTAMP;
    END if;
    RETURN new;
END
$FUNC$
LANGUAGE plpgsql;  
CREATE TRIGGER set_job_finished_at
  BEFORE UPDATE ON queue.jobs
  FOR EACH ROW
  EXECUTE PROCEDURE set_job_finished_at();

DROP TYPE IF EXISTS api.job_payload CASCADE;
CREATE TYPE api.job_payload AS (
  job_id BIGINT,
  work JSONB
);

DROP FUNCTION IF EXISTS api.start_job CASCADE;
CREATE FUNCTION api.start_job(worker_id integer)
RETURNS api.job_payload AS 
$FUNC$
UPDATE queue.jobs 
SET 
  status='in_progress', 
  started_at=CURRENT_TIMESTAMP,
  worker_id = $1
WHERE id IN (
  SELECT id 
  FROM queue.jobs 
  WHERE status='waiting' 
  FOR UPDATE SKIP LOCKED 
  LIMIT 1
) RETURNING id, work;
$FUNC$ 
LANGUAGE SQL;

NOTIFY pgrst, 'reload schema';