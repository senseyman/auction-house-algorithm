GOLINT := golangci-lint

PACKAGES_FOR_TEST := $(shell go list ./... | grep -v model | grep -v "mock" | grep -v "model")


all: dep gen-mock lint vet test

dep:
	go mod tidy
	go mod download

dep-update:
	go get -t -u ./...

test:
	@go test -cover -race -tags=unit -parallel 10 -count=1 -v $(PACKAGES_FOR_TEST)

vet:
	go vet ./...

check-lint:
	@which $(GOLINT) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin v1.57.2

lint: dep check-lint ## Lint the files local env
	$(GOLINT) run --timeout=5m -c .golangci.yml

check-mockgen:
	@which mockgen || go install go.uber.org/mock/mockgen@latest

gen-mock:
	mockgen -package mock -source service/auction/interface.go -destination service/auction/mock/interface.go