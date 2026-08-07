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

GEN_UUID = $(shell cat /proc/sys/kernel/random/uuid)
GEN_USER_ID = $(shell shuf -i 1-1000 -n 1)

add_event:
	@curl -i -X POST http://localhost:9090/events \
		-H "Content-Type: application/json" \
		-d '{"event_id": "'$(GEN_UUID)'", "user_id": '$(GEN_USER_ID)', "action": "click", "timestamp": "'$(shell date --iso-8601=seconds)'"}'
