run:
	@go run cmd/main.go

migrate_up:
	migrate -path db/migrations -database "postgres://postgres:secretpassword@localhost:5432/events_db?sslmode=disable" up

migrate_down:
	migrate -path db/migrations -database "postgres://postgres:secretpassword@localhost:5432/events_db?sslmode=disable" down

migrate_up1:
	migrate -path db/migrations -database "postgres://postgres:secretpassword@localhost:5432/events_db?sslmode=disable" up 1

migrate_down1:
	migrate -path db/migrations -database "postgres://postgres:secretpassword@localhost:5432/events_db?sslmode=disable" down 1
