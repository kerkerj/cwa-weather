.PHONY: build test lint sec cover check clean

build:
	go build -o bin/cwa-weather ./cmd/cwa-weather

test:
	go test -v ./...

lint:
	golangci-lint run ./...

sec:
	gosec ./...

cover:
	@go test -coverprofile=coverage.out ./... > /dev/null 2>&1
	@COVERAGE=$$(go tool cover -func=coverage.out | grep total | awk '{print $$3}' | tr -d '%'); \
	echo "Coverage: $${COVERAGE}%"; \
	if [ $$(echo "$${COVERAGE} < 75" | bc) -eq 1 ]; then \
		echo "FAIL: coverage $${COVERAGE}% < 75% threshold"; \
		exit 1; \
	fi
	@rm -f coverage.out

check: test cover lint sec

clean:
	rm -rf bin/
