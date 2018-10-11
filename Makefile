.PHONY: setup
setup:
	@if [ -z `which dep 2> /dev/null` ]; then \
		go get -u github.com/golang/dep/cmd/dep;\
	fi
	@dep ensure --vendor-only

.PHONY: build
build: cmd/nouha.go
	go build -o bin/nouha cmd/nouha.go

.PHONY: run
run: bin/nouha
	./bin/nouha

.PHONY: install
install: build
	@go install

.PHONY: lint
lint: govet gofmt golint goimports

.PHONY: govet
govet:
	@go list ./... | xargs go vet

.PHONY: gofmt
gofmt:
	@echo "$(GO_FILES)" | xargs -Ix gofmt -w x

.PHONY: golint
golint:
	@if [ -z `which golint 2> /dev/null` ]; then \
		go get -u github.com/golang/lint/golint; \
	fi
	@go list ./... | xargs golint -set_exit_status

.PHONY: goimports
goimports:
	@if [ -z `which goimports 2> /dev/null` ]; then \
		go get golang.org/x/tools/cmd/goimports; \
	fi
	@echo "$(GO_FILES)" | xargs -Ix goimports -w x
