// Copyright (c) 2026 the go-widgets/mvvmlint authors. All rights reserved.
// Use of this source code is governed by a BSD-3-Clause license that can be
// found in the LICENSE file at the root of this repository.

//mvvmlint:allow

// Package allowdir opts out of the direct-mutation rule with a file directive,
// so its intentional widget-state writes are not flagged.
package allowdir

import (
	"github.com/go-widgets/mvvm"
	"github.com/go-widgets/toolkit"
)

func glue() {
	e := toolkit.NewEntry()
	e.Text = "intentional low-level glue" // exempt via //mvvmlint:allow
	_ = mvvm.NewObservable("x")
}
