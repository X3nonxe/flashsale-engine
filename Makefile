postgres:
	docker run --name flashsale_db -p 5432:5432 -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=password123 -d postgres:15-alpine

createdb:
	docker exec -it flashsale_db createdb --username=admin --owner=admin flashsale_db

dropdb:
	docker exec -it flashsale_db dropdb flashsale_db

migrateup:
	migrate -path db/migration -database "postgresql://admin:password123@localhost:5432/flashsale_db?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://admin:password123@localhost:5432/flashsale_db?sslmode=disable" -verbose down

test:
	go test -v ./...

run:
	go run cmd/api/main.go

.PHONY: postgres createdb dropdb migrateup migratedown test run