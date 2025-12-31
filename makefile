run:
	go run .

test:
	go test ./...

gen:
	go tool sqlc generate

.PHONY: run test gen
