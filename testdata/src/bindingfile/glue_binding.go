// Copyright (c) 2026 the go-widgets/mvvmlint authors. All rights reserved.
// Use of this source code is governed by a BSD-3-Clause license that can be
// found in the LICENSE file at the root of this repository.

package bindingfile

import "github.com/go-widgets/toolkit"

// This file is named *_binding.go, so its widget-state writes are exempt.
func glue() {
	e := toolkit.NewEntry()
	e.Text = "sanctioned binding-file glue" // exempt: *_binding.go
}
