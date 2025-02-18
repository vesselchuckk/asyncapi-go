db_login:
	psql postgresql://admin:qwerty@localhost:5432/postgres?sslmode=disable
db_create_migration:
	migrate create -ext sql -dir migrations -seq init_schema
db_migrate:
	migrate -database postgresql://admin:qwerty@localhost:5432/postgres?sslmode=disable -path migrations up
start:
	go run cmd/main.go