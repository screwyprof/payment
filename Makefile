GOBIN?=${GOPATH}/bin
PACKAGES=$(shell go list ./... | grep -v ./test)

## all              : fetch deps, build, test and install the application
all: deps unit-test build install

## build            : build application
build:
	@echo "--> Building application"
	go build -race -o build/gopherpay ./cmd/gopherpay

## install          : install the application to $GOPATH/bin
install:
	@echo "--> Installing application"
	go install ./cmd/gopherpay

## deps             : fetch dependencies
deps:
	@echo "--> Running dep"
	GO111MODULE=on go mod download

## unit-test        : run unit tests
unit-test:
	@echo "--> Running unit tests"
	GOCACHE=off go test -v $(PACKAGES)

## unit-test-race   : run unit tests with -race flag
unit-test-race:
	@echo "--> Running go test --race"
	GOCACHE=off go test -v -race $(PACKAGES)

## integration-test : run integration tests
integration-test:
	@echo "--> Running unit tests"
	GOCACHE=off go test -race -v ./test/...

## fmt              : format go files
fmt:
	@echo "--> Formatting go files"
	go fmt ./...

## clean            : cleaning up
clean:
	@echo "--> Cleaning..."
	@rm -rf ./build

## docker-build     : building docker image
docker-build:
	@echo "--> Building docker image ..."
	docker build --build-arg GO_VERSION=1.11.2 -t gopherpayment .

## docker-run       : running service in docker in dev mode
docker-run:
	@echo "--> Running docker container ..."
	docker run --rm -p 8080:8080 --name gopherpay gopherpayment:latest

## help             : show this help screen
help : Makefile
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | sort

# To avoid unintended conflicts with file names, always add to .PHONY
# unless there is a reason not to.
# https://www.gnu.org/software/make/manual/html_node/Phony-Targets.html
.PHONY: all build install deps unit-test unit-test-race integration-test docker-build docker-run fmt clean
