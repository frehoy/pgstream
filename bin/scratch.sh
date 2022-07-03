#!/usr/bin/env sh

# shellcheck disable=SC1091
. ./writer/.env

# curl "localhost:3000" \
# 	-H "Content-Type: application/json" \
# 	-H "Authorization: Bearer $TOKEN"

# Reload schema cache
docker-compose kill -s SIGUSR1 rest

# Create a job
printf "\npost to jobs:\n"
curl "localhost:3000/jobs" \
	-X POST \
	-d '{"work":{ "direct": "insert"}}' \
	-H "Content-Type: application/json" \
	-H "Authorization: Bearer $TOKEN"
printf "\npost to jobs:\n"
curl "localhost:3000/jobs" \
	-X POST \
	-d '{"work":{ "direct": "insert"}}' \
	-H "Content-Type: application/json" \
	-H "Authorization: Bearer $TOKEN"
printf "\npost to jobs:\n"
curl "localhost:3000/jobs" \
	-X POST \
	-d '{"work":{ "direct": "insert"}}' \
	-H "Content-Type: application/json" \
	-H "Authorization: Bearer $TOKEN"

printf "\nstart_job:\n"
start_job_resp=$(curl "localhost:3000/rpc/start_job" \
	-X POST \
	-d '{ "worker_id": 1 }' \
	-H "Content-Type: application/json" \
	-H "Authorization: Bearer $TOKEN"
)
job_a=$(echo "$start_job_resp" | jq .job_id)
echo "Got job_id: $job_a"

printf "\nstart_job:\n"
start_job_resp=$(curl "localhost:3000/rpc/start_job" \
	-X POST \
	-d '{ "worker_id": 2 }' \
	-H "Content-Type: application/json" \
	-H "Authorization: Bearer $TOKEN"
)
job_b=$(echo "$start_job_resp" | jq .job_id)
echo "Got job_id: $job_a"

printf "\nstart_job:\n"
start_job_resp=$(curl "localhost:3000/rpc/start_job" \
	-X POST \
	-d '{ "worker_id": 3 }' \
	-H "Content-Type: application/json" \
	-H "Authorization: Bearer $TOKEN"
)
job_c=$(echo "$start_job_resp" | jq .job_id)
echo "Got job_id: $job_a"


# Doing "work"
sleep 3

curl "localhost:3000/jobs?id=eq.$job_a" \
	-X PATCH \
	-d '{"status": "finished"}' \
	-H "Content-Type: application/json" \
	-H "Authorization: Bearer $TOKEN"

# Doing "work"
sleep 3

curl "localhost:3000/jobs?id=eq.$job_b" \
	-X PATCH \
	-d '{"status": "finished"}' \
	-H "Content-Type: application/json" \
	-H "Authorization: Bearer $TOKEN"

# Doing "work"
sleep 3

curl "localhost:3000/jobs?id=eq.$job_c" \
	-X PATCH \
	-d '{"status": "finished"}' \
	-H "Content-Type: application/json" \
	-H "Authorization: Bearer $TOKEN"
