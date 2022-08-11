modules:
	cd src && go mod download
mock:
	cd src && go generate ./...
lint:
	cd src && golangci-lint --timeout 2m --config ../.golangci.yml run --exclude-use-default=false ./...
test:
	cd src && go test -tags test -parallel=7 ./...
test-coverage:
	cd src && go test -tags test -parallel=7 ./... -coverprofile ../cover.out -covermode=atomic
	go install github.com/boumenot/gocover-cobertura@latest
	cd src && ${GOPATH}/bin/gocover-cobertura < ../cover.out > ../coverage.xml
	cd src && go tool cover -func ../cover.out | grep total

dynamo-local:
	docker network inspect sgorecki-event-store > /dev/null 2>&1  || docker network create sgorecki-event-store
	docker run --net sgorecki-event-store -d --rm --name dynamodb -p 8000:8000 amazon/dynamodb-local
	docker run --net sgorecki-event-store -d --rm --name event-store-admin -p 8001:8001 -e AWS_REGION=local -e AWS_ACCESS_KEY_ID=local -e AWS_SECRET_ACCESS_KEY=local -e DYNAMO_ENDPOINT=http://dynamodb:8000 aaronshaf/dynamodb-admin

dynamo-stop:
	docker stop event-store-admin
	docker stop dynamodb
