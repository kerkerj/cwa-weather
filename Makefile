.PHONY: build test lint sec check clean

build:
	go build -o bin/cwa-weather ./cmd/cwa-weather

test:
	go test -v ./...

lint:
	golangci-lint run ./...

sec:
	gosec ./...

check: test lint sec

clean:
	rm -rf bin/
