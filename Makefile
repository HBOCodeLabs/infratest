SOURCE = ./...

.DEFAULT_GOAL := test

export SHELL:=/bin/bash
export SHELLOPTS:=$(if $(SHELLOPTS),$(SHELLOPTS):)pipefail:errexit

export K8S_VERSION:=1.21.1

.ONESHELL:

vet:
	go vet $(SOURCE)
.PHONY: vet

test-fmt:
	test -z $(shell go fmt $(SOURCE))
.PHONY: test-fmt

test: vet test-fmt
	go test -cover ./pkg/... -count=1
.PHONY: test

tools:
	echo "Installing tools from tools.go"
	cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %
.PHONY: tools

mock: tools
	mockgen -source pkg/aws/dax.go -destination mock/dax.go -package mock
	mockgen -source pkg/aws/ec2.go -destination mock/ec2.go -package mock
	mockgen -source pkg/k8s/jobs.go -destination mock/k8s_jobs.go -package mock
.PHONY: 

.PHONY: integration-test
integration-test:
	echo K8S_VERSION: $(K8S_VERSION)
	go test -v -timeout 10m -count 1 ./integration/...