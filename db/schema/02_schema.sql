CREATE VIEW api.searches AS
SELECT
	id,
	created_at AS ts,
	message->'search'->>'category' as category, 
	message->'search'->>'query' as query
FROM api.events 
WHERE message->>'message_type' = 'search';

CREATE VIEW api.searches_top_50_by_minute AS
SELECT
	date_trunc('minute', ts) AS ts,
	category,
	query,
	count(*) as seaches_count
FROM api.searches
GROUP BY 1, 2, 3
ORDER BY 4 DESC
LIMIT 50;