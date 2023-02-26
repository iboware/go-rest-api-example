.PHONY: test
test.unit:
	echo "=> Running Tests"
	go test -tags=unit -v ./...

.PHONY: build
build:
	echo "=> Building..."
	CGO_ENABLED=0 go build -a -ldflags '-w -s' -o bin/go-rest-api-example

.PHONY: run
run:
	./bin/go-rest-api-example

.PHONY: docs
docs:
	swag i
	
docs.fmt:
	swag fmt