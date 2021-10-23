SOURCE = ./...

.DEFAULT_GOAL := test

vet:
	go vet $(SOURCE)
.PHONY: vet

test-fmt:
	test -z $(shell go fmt $(SOURCE))
.PHONY: test-fmt

test: vet test-fmt
	go test -cover $(SOURCE) -count=1
.PHONY: test
