SOURCE = ./...

.DEFAULT_GOAL := test

export SHELL:=/bin/bash
export SHELLOPTS:=$(if $(SHELLOPTS),$(SHELLOPTS):)pipefail:errexit

export K8S_VERSION:=1.21.1
export VAULT_VERSION:=1.10.3

.ONESHELL:

lint: tools
	golangci-lint run --new-from-rev=main
.PHONY: lint

lint-full:
	golangci-lint run
.PHONY: lint

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
	mockgen -source pkg/aws/iam.go -destination mock/iam.go -package mock
	mockgen -source pkg/aws/eks.go -destination mock/eks.go -package mock
	mockgen -source pkg/k8s/jobs.go -destination mock/k8s_jobs.go -package mock
	mockgen -source pkg/k8s/util.go -destination mock/k8s_util.go -package mock
	go generate ./...

.PHONY: mock

.PHONY: k8s-integration-test
k8s-integration-test:
	echo K8S_VERSION: $(K8S_VERSION)
	go test -v -timeout 10m -count 1 ./integration/k8s_test.go

.PHONY: vault-integration-test
vault-integration-test:
	echo VAULT_VERSION: $(VAULT_VERSION)
	go test -v -timeout 10m -count 1 ./integration/vault_test.go