up:
	docker compose up --detach
down:
	docker compose down --volumes
build:
	docker compose build
re: build down up

bench: up
	docker compose run --rm writer
psql: up
	docker compose exec db psql --user track
