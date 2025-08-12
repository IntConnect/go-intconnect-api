.PHONY: migrate-fresh
migrate-fresh:
	migrate -path migrations -database "postgres://postgres:@127.0.0.1:5432/go_intconnect_api?sslmode=disable" down
	migrate -path migrations -database "postgres://postgres:@127.0.0.1:5432/go_intconnect_api?sslmode=disable" up

.PHONY: migrate-up
migrate-up:
	migrate -path migrations -database "postgres://postgres:@127.0.0.1:5432/go_intconnect_api?sslmode=disable" up 1

migrate:
	migrate -path migrations -database "postgres://postgres:@127.0.0.1:5432/go_intconnect_api?sslmode=disable" up


.PHONY: migrate-down
migrate-down:
	migrate -path migrations -database "postgres://postgres:@127.0.0.1:5432/go_intconnect_api?sslmode=disable" down 1

migrate-force:
	migrate -path migrations -database "postgres://postgres:@127.0.0.1:5432/go_intconnect_api?sslmode=disable" force $(version)

inject:
	wire gen ./cmd/injection/injector.go

