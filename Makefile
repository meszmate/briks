.PHONY: build run clean vet release snapshot

build:
	go build -o briks ./cmd/briks/

run: build
	./briks

clean:
	rm -f briks
	rm -rf dist/

vet:
	go vet ./...

release:
	goreleaser release --clean

snapshot:
	goreleaser release --snapshot --clean
