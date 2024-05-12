DOCKER_COMPOSE_FILE ?= docker-compose.yml

migrate-down:
	docker-compose -f ${DOCKER_COMPOSE_FILE} --profile tools run --rm migrate down

migrate-up:
	docker-compose -f ${DOCKER_COMPOSE_FILE} --profile tools run --rm migrate up

build: migrate-up
	docker-compose up
