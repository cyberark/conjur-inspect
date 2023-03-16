build:
	go build -o ./dev/tmp/ ./cmd/conjur-inspect

test:
	go test -count=1 -coverpkg=./... -coverprofile=c.out -v ./...
	go tool cover -func c.out

install:
	go install ./cmd/conjur-inspect

release:
	goreleaser release --snapshot --rm-dist

# Example usage of run: make run ARGS="variable get -i path/to/var"
run:
	go run ./cmd/conjur-inspect $(ARGS)
