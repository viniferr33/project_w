include .env

clear:
	@if [ -f "${BINARY}" ]; then \
		rm ${BINARY}; \
	fi
	@echo "Deleted binary ${BINARY}";

compile:
	@echo "Building binary ... ${BINARY}";
	go build -o ${BINARY} cmd/wondrous_baazar.go

run: compile
	./${BINARY}

docker_build:
	docker build . -t "${IMAGE_NAME}"

compose_up:
	docker compose up -d

compose_down:
	docker compose down

migrate_create:
	docker compose --profile tools run --rm migrate create -ext sql -dir /migrations $(name)

migrate_up:
	docker compose --profile tools run --rm migrate up

migrate_down:
	docker compose --profile tools run --rm migrate down

test:
	go test -v -cover ./...

integration_test:
	docker-compose -f docker-compose.test.yml run test