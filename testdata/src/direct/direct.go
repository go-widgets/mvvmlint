// Copyright (c) 2026 the go-widgets/mvvmlint authors. All rights reserved.
// Use of this source code is governed by a BSD-3-Clause license that can be
// found in the LICENSE file at the root of this repository.

// Package direct imports both the toolkit and mvvm (so the missing-view-model
// rule stays silent) and mutates widget state directly — the core violation.
package direct

import (
	"github.com/go-widgets/mvvm"
	"github.com/go-widgets/toolkit"
)

func mutate() {
	e := toolkit.NewEntry()
	e.Text = "typed ad hoc" // want `direct mutation of toolkit widget state field "Text"`

	lb := toolkit.NewListBox()
	lb.Items = []string{"a", "b"} // want `direct mutation of toolkit widget state field "Items"`
	lb.Selected = 1               // want `direct mutation of toolkit widget state field "Selected"`

	// A compound assignment is still a mutation.
	var v toolkit.Entry
	v.Text += "!" // want `direct mutation of toolkit widget state field "Text"`

	// Assigning a callback slot is NOT a state field — this is how composition
	// works and must never be flagged.
	e.OnChange = func(string) {}

	// A non-toolkit type with a same-named field must never be flagged.
	other := struct{ Text string }{}
	other.Text = "fine"

	_ = mvvm.NewObservable("x")
}
