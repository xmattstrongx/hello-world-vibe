APP := hello-go

.PHONY: fmt test build run run-live clean

fmt:
	gofmt -w ./cmd ./internal ./main.go

test:
	go test ./...

build:
	mkdir -p bin
	go build -o bin/$(APP) ./cmd/hello-go

run:
	go run .

run-live:
	go run . --live

clean:
	rm -rf bin
