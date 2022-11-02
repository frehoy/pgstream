# pgstream

Streaming event tracking but it's just postgresql, postgREST and pg_cron. 

`make bench` will start the writer which sends randomly generated json events of faked "searches" with a category and a random query to the PostgREST api. Pretend that they're users submitting some tracking data.
PostgREST writes them to the `api.events` table which is partitioned by minute. 
When a minute is over the searches are counted and inserted into `agg.searches_per_minute`. 
It does 3000 events per second with two writing threads on my laptop. It's fast and cpu-cheap because the aggregation query only has to read a single partiton (the last uncounted minute). 

TODO: 
* Figure out if the writer works properly with multiple threads, what's a goroutine? 
* Roll up minutely counts into hourly, daily, weekly, monthly and total counts. 
* Drop old minute partitions from `api.events`. Put them in a daily partitioned archive table maybe. 
