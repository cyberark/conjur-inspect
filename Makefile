build:
	go build -o ./dev/tmp/ ./cmd/conjur-preflight

test:
	go test -count=1 -coverprofile=c.out -v ./...

install:
	go install ./cmd/conjur-preflight

release:
	goreleaser release --snapshot --rm-dist

# Example usage of run: make run ARGS="variable get -i path/to/var"
run:
	go run ./cmd/conjur-preflight $(ARGS)
