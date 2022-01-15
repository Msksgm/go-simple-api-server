run:
	go run main.go
up:
	docker compose up -d
down:
	docker compose down
create-migration:
	migrate create -ext sql -dir postgres/migrations -seq ${file}
run-migration:
	migrate -database $(POSTGRESQL_URL) -path postgres/migrations up