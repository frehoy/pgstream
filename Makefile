up:
	docker compose up --detach
down:
	docker compose down --volumes
build:
	docker compose build writer db
re: build down up

bench: up build
	docker compose run --rm writer
psql: up
	docker compose exec db psql --user track

test:
	cd writer; go test
