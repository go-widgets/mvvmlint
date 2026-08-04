// Copyright (c) 2026 the go-widgets/mvvmlint authors. All rights reserved.
// Use of this source code is governed by a BSD-3-Clause license that can be
// found in the LICENSE file at the root of this repository.

// Package bound wires widget state the sanctioned way: through the mvvm binders
// with an address-of into the field. It must produce no diagnostics.
package bound

import (
	"github.com/go-widgets/mvvm"
	"github.com/go-widgets/toolkit"
)

func wire() {
	e := toolkit.NewEntry()
	obs := mvvm.NewObservable("hello")

	// The address-of form &e.Text is a binding, not a mutation: never flagged.
	mvvm.BindField(obs, &e.Text, &e.OnChange, nil)

	// A composite literal that seeds a field at construction is not an
	// assignment to a selector and is never flagged.
	_ = toolkit.Entry{Text: "seed"}
}
