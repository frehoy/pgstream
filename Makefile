up:
	docker compose up --detach
down:
	docker compose down --volumes
re: down up 

bench: up
	docker compose run --rm bench
psql: up
	docker compose exec db psql --user track
