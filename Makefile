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

tools:
	echo "Installing tools from tools.go"
	cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %
.PHONY: tools

mock: tools
	mockgen -source pkg/aws/dax.go -destination mock/dax.go -package mock
	mockgen -source pkg/aws/ec2.go -destination mock/ec2.go -package mock
.PHONY: mock
