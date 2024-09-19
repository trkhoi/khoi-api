dev:
	go run cmd/api/main.go
test:
	go test -p 1 -mod=mod -coverprofile=c.out -failfast -timeout 5m ./...