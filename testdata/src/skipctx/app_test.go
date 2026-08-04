// Copyright (c) 2026 the go-widgets/mvvmlint authors. All rights reserved.
// Use of this source code is governed by a BSD-3-Clause license that can be
// found in the LICENSE file at the root of this repository.

package skipctx

import "github.com/go-widgets/toolkit"

// A test-only fixture mutating widget state directly: allowed by default
// (mvvmlint skips *_test.go unless -includetests).
func fixture() {
	e := toolkit.NewEntry()
	e.Text = "ad-hoc"
}
