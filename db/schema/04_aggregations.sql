-- Counts searches from api.searches.

DROP SCHEMA IF EXISTS agg CASCADE;
CREATE SCHEMA agg;

CREATE TABLE agg.searches_per_minute(
  minute timestamp with time zone NOT NULL,
  category text  NOT NULL,
  query text  NOT NULL,
  c bigint  NOT NULL
);
CREATE INDEX ON agg.searches_per_minute USING BTREE (minute);
COMMENT ON TABLE agg.searches_per_minute IS 'Number of searches per minute by category and query';

CREATE VIEW agg.searches_per_minute_to_count AS
SELECT
  date_trunc('minute', ts) AS minute,
  category,
  query,
  count(*) as c
FROM api.searches
WHERE 
  -- Not in the current minute (count searches only when the minute is over).
  ts < date_trunc('minute', CURRENT_TIMESTAMP)
  AND 
  -- After last counted minute
  ts >= COALESCE(
    -- The start of the first uncounted minute
    (SELECT MAX(minute)  + '1 minute'::interval FROM agg.searches_per_minute),
    -- If we have no counted minutes use a date long ago
    '1990-01-01'::timestamp
    )
GROUP BY 1, 2, 3
;
COMMENT ON VIEW agg.searches_per_minute_to_count IS 'Uncounted searches. Insert from this into agg.searches_per_minute';

-- Schedule update_searches_per_minute every minute
SELECT cron.schedule(
  'update-searches-per-minute', 
  '* * * * *', 
  'INSERT INTO agg.searches_per_minute SELECT * FROM agg.searches_per_minute_to_count'
);
