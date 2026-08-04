// Copyright (c) 2026 the go-widgets/mvvmlint authors. All rights reserved.
// Use of this source code is governed by a BSD-3-Clause license that can be
// found in the LICENSE file at the root of this repository.

// Package statefieldsflag is checked with a custom -statefields allowlist
// ("Items,,Selected"): only those fields are guarded, so a Text write is
// allowed while Items/Selected writes are flagged.
package statefieldsflag

import (
	"github.com/go-widgets/mvvm"
	"github.com/go-widgets/toolkit"
)

func mutate() {
	e := toolkit.NewEntry()
	e.Text = "not in the custom allowlist" // allowed by -statefields override

	lb := toolkit.NewListBox()
	lb.Items = nil // want `direct mutation of toolkit widget state field "Items"`
	lb.Selected = 2 // want `direct mutation of toolkit widget state field "Selected"`

	_ = mvvm.NewObservable("x")
}
