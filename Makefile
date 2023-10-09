.PHONY: test

test:
	go build -o ./cmd/gophermart/ ./cmd/gophermart/
	chmod +x ./cmd/gophermart
	./gophermarttest \
  	   -test.v -test.run=^TestGophermart/TestEndToEnd/register_user \
  	   -gophermart-binary-path=cmd/gophermart/gophermart \
  	   -gophermart-host=localhost \
  	   -gophermart-port=8080 \
  	   -gophermart-database-uri="postgres://postgres:pass@localhost:5432/gophermart?sslmode=disable" \
  	   -accrual-binary-path=cmd/accrual/accrual_linux_amd64 \
  	   -accrual-host=localhost \
  	   -accrual-port=8083 \
  	   -accrual-database-uri="postgres://postgres:pass@localhost:5432/gophermart?sslmode=disable"
.PHONY: runAccural
runAccural:
	./cmd/accrual/accrual_linux_amd64

.PHONY: run
run:
	go run ./cmd/gophermart -d "postgres://postgres:pass@localhost:5432/gophermart?sslmode=disable"