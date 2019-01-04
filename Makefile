CONTAINER_NAME=luppiter-test
PG_PASSWORD=rootpass
PG_DB=luppiter_test
PG_PORT=15432

test:
	@docker stop $(CONTAINER_NAME) > /dev/null || true
	@docker rm $(CONTAINER_NAME) > /dev/null || true
	@docker run --name $(CONTAINER_NAME) -e 'POSTGRES_PASSWORD=$(PG_PASSWORD)' -e 'POSTGRES_DB=$(PG_DB)' -p $(PG_PORT):5432 -d postgres:9
	@sleep 10

	@DB_HOST=127.0.0.1 DB_USERNAME=postgres DB_PASSWORD=$(PG_PASSWORD) DB_NAME=$(PG_DB) DB_PORT=$(PG_PORT) go test -v -cover ./... || true

	@docker stop $(CONTAINER_NAME) > /dev/null || true
	@docker rm $(CONTAINER_NAME) > /dev/null || true
