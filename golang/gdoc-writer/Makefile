.PHONY: lint test test-coverage mock modules build all

lint:
	cd src && golangci-lint --timeout 3m --config ../.golangci.yml run --exclude-use-default=false ./...

test:
	cd src && go test -tags test ./...

test-coverage:
	cd src && go test -tags test ./... -covermode=atomic -coverpkg=./... -coverprofile ../cover.out && go tool cover -func ../cover.out | grep total
	go install github.com/boumenot/gocover-cobertura@latest
	cd src && ${GOPATH}/bin/gocover-cobertura < ../cover.out > ../coverage.xml

mock:
	cd src && go generate ./...

modules:
	cd src && go mod download

build:
	cd src && go build -o gdoc-writer

all:
	make mock
	make lint
	make test
	make build