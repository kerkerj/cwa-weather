.PHONY: build test lint clean

build:
	go build -o bin/cwa-weather ./cmd/cwa-weather

test:
	go test -v ./...

lint:
	golangci-lint run ./...

clean:
	rm -rf bin/
