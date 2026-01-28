.PHONY: build run clean vet

build:
	go build -o briks ./cmd/briks/

run: build
	./briks

clean:
	rm -f briks

vet:
	go vet ./...
