.PHONY: build run test lint clean debug

build:
	go build -o amaze ./cmd/amaze

run: build
	./amaze config.txt

test:
	go test ./...

lint:
	go vet ./...

clean:
	rm -f amaze

debug:
	go run ./cmd/amaze -v config.txt
