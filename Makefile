protoc:
	protoc -I. --go-grpc_out=. --go_out=. pkg/proto/$(file).proto

lint:
	golangci-lint run --config=./.golangci.yml

migrate-new:
	migrate create -ext sql -dir db/migration -seq $(name)

migrate-up:
	migrate -path db/migration \
	-database "postgresql://root:pass@127.0.0.1:5436/appointments?sslmode=disable&application_name=appointment-service" \
	-verbose up

migrate-down:
	migrate -path db/migration \
  -database "postgresql://root:pass@127.0.0.1:5436/appointments?sslmode=disable&application_name=appointment-service" \
  -verbose down

mockgen gateway:
	mockgen --source=./internal/gateway/$(file).go -destination=./internal/gateway/mock/$(file).go -package=mock

test:
	go test ./...

test-verbose:
	go test -v ./... -cover

test-cover:
	go test -coverprofile=coverage.out ./... &&  go tool cover -html=coverage.out

up:
	docker compose up -d