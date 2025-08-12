.PHONY: migrate-fresh
migrate-fresh:
	migrate -path migrations -database "postgres://postgres:@127.0.0.1:5432/go_intconnect_system?sslmode=disable" down
	migrate -path migrations -database "postgres://postgres:@127.0.0.1:5432/go_intconnect_system?sslmode=disable" up

.PHONY: migrate-up
migrate-up:
	migrate -path migrations -database "postgres://postgres:@127.0.0.1:5432/go_intconnect_system?sslmode=disable" up 1

migrate:
	migrate -path migrations -database "postgres://postgres:@127.0.0.1:5432/go_intconnect_system?sslmode=disable" up


.PHONY: migrate-down
migrate-down:
	migrate -path migrations -database "postgres://postgres:@127.0.0.1:5432/go_intconnect_system?sslmode=disable" down 1

migrate-force:
	migrate -path migrations -database "postgres://postgres:@127.0.0.1:5432/go_intconnect_system?sslmode=disable" force $(version)

inject:
	wire gen ./cmd/injection/injector.go

