// Copyright (c) 2026 the go-widgets/mvvmlint authors. All rights reserved.
// Use of this source code is governed by a BSD-3-Clause license that can be
// found in the LICENSE file at the root of this repository.

// Package bindingfile proves the *_binding.go exemption is file-scoped: an
// ordinary file in the same package is still guarded.
package bindingfile

import (
	"github.com/go-widgets/mvvm"
	"github.com/go-widgets/toolkit"
)

func app() {
	e := toolkit.NewEntry()
	e.Text = "mutated in an ordinary file" // want `direct mutation of toolkit widget state field "Text"`
	_ = mvvm.NewObservable("x")
}
