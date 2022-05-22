up:
	docker compose up --detach
down:
	docker compose down --volumes
bench: down up
	docker compose run --rm bench
psql: up
	docker compose exec db psql --user track
