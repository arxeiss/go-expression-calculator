// +build tools

package main

import (
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/onsi/ginkgo/ginkgo"
	// _ "golang.org/x/lint/golint"
	// _ "golang.org/x/tools/cmd/goimports"
	// _ "github.com/securego/gosec/v2/cmd/gosec"
	// _ "honnef.co/go/tools/cmd/staticcheck"
)
