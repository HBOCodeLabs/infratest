//go:build tools
// +build tools

package tools

import (
	_ "github.com/golang/mock/mockgen"
	_ "sigs.k8s.io/kind"
	_ "golang.org/x/tools/cmd/godoc"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
)
