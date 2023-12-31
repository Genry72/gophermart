.PHONY: test

test:
	go build -o ./cmd/gophermart/ ./cmd/gophermart/
	chmod +x ./cmd/gophermart
	./gophermarttest \
  	   -test.v -test.run=^TestGophermart$ \
  	   -gophermart-binary-path=cmd/gophermart/gophermart \
  	   -gophermart-host=localhost \
  	   -gophermart-port=8080 \
  	   -gophermart-database-uri="postgres://postgres:pass@localhost:5432/gophermart?sslmode=disable" \
  	   -accrual-binary-path=cmd/accrual/accrual_linux_amd64 \
  	   -accrual-host=localhost \
  	   -accrual-port=8083 \
  	   -accrual-database-uri="postgres://postgres:pass@localhost:5432/gophermart?sslmode=disable"

.PHONY: run
run:
	go run ./cmd/gophermart -d "postgres://postgres:pass@localhost:5432/gophermart?sslmode=disable"

.PHONY: accrual
accrual:
	./cmd/accrual/accrual_linux_amd64 \
	-d "postgres://postgres:pass@localhost:5433/gophermart?sslmode=disable" \
	-a "localhost:8083"

.PHONY: gen
gen:
	mockgen -source=internal/repositories/repositories.go \
    -destination=internal/repositories/mocks/mock_repository.go && \
    mockgen -source=internal/usecases/usecase.go \
    -destination=internal/usecases/mocks/mock_usecase.go && \
    mockgen -source=internal/handlers/jwtauth/auth.go \
    -destination=internal/handlers/jwtauth/mocks/mock_auth.go

.PHONY: mytest
mytest:
	go test -v -count=1 ./...

.PHONY: cover
cover:
	go test -short -count=1 -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out