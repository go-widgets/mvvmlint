// Copyright (c) 2026 the go-widgets/mvvmlint authors. All rights reserved.
// Use of this source code is governed by a BSD-3-Clause license that can be
// found in the LICENSE file at the root of this repository.

// Package incltests exercises -includetests: its *_test.go mutation is flagged.
package incltests

import (
	"github.com/go-widgets/mvvm"
	"github.com/go-widgets/toolkit"
)

func wire() {
	e := toolkit.NewEntry()
	mvvm.BindField(mvvm.NewObservable("x"), &e.Text, &e.OnChange, nil)
}
