DROP SCHEMA IF EXISTS partitions CASCADE;
CREATE SCHEMA partitions;

CREATE OR REPLACE FUNCTION partitions.create_partition_ddl(
  of_table regclass,
  part_start text,
  part_end text,
  suffix text
)
RETURNS text
LANGUAGE plpgsql
AS $FUNC$
DECLARE
  partition_name text;
  values_range text;
  ddl text;
BEGIN
partition_name := format('%s%s', of_table, suffix);
values_range := format('(%L) TO (%L)', part_start, part_end);
ddl := format('CREATE TABLE IF NOT EXISTS %s PARTITION OF %s FOR VALUES FROM %s', partition_name, of_table, values_range);
RETURN ddl;
END
$FUNC$
;
COMMENT ON FUNCTION partitions.create_partition_ddl IS 'Specification of time partitions.';

CREATE VIEW partitions.hourly AS
WITH 
limits AS (
  SELECT
    ts AS part_start,
    ts + '1 hour'::interval AS part_end,
    to_char(ts, '_YYYY_mm_dd_HH24_MI') AS suffix
  FROM generate_series(
    date_trunc('hour', CURRENT_TIMESTAMP), 
    date_trunc('hour', CURRENT_TIMESTAMP) + '2 hour'::interval, 
    '1 hour'::interval
  ) AS ts
)
SELECT
  part_start,
  part_end,
  suffix
FROM 
limits;

CREATE PROCEDURE partitions.update_partitions(base_table regclass)
LANGUAGE plpgsql
AS $PROC$
DECLARE
  rec record;
  create_partition_ddl text;
BEGIN
  FOR rec IN SELECT * FROM partitions.hourly ORDER BY part_start
  LOOP
    create_partition_ddl := partitions.create_partition_ddl(base_table, rec.part_start::text, rec.part_end::text, rec.suffix);
    RAISE NOTICE 'Creating table from % to %', rec.part_start, rec.part_end;  
    EXECUTE create_partition_ddl;
  END LOOP;
END
$PROC$
;

-- Create initial partitions
CALL partitions.update_partitions('api.events'::regclass);
-- Schedule partitions update for minute 0 every hour.
SELECT cron.schedule('partitions-update', '0 * * * *', 'CALL partitions.update_partitions($$api.events$$::regclass)');