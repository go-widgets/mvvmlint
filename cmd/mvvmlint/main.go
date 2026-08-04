// Copyright (c) 2026 the go-widgets/mvvmlint authors. All rights reserved.
// Use of this source code is governed by a BSD-3-Clause license that can be
// found in the LICENSE file at the root of this repository.

// Command mvvmlint is a standalone driver for the mvvmlint analyzer. It doubles
// as a go vet tool:
//
//	go install github.com/go-widgets/mvvmlint/cmd/mvvmlint@latest
//	go vet -vettool=$(which mvvmlint) ./...
package main

import (
	"github.com/go-widgets/mvvmlint"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(mvvmlint.Analyzer)
}
