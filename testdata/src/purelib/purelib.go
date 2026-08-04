// Copyright (c) 2026 the go-widgets/mvvmlint authors. All rights reserved.
// Use of this source code is governed by a BSD-3-Clause license that can be
// found in the LICENSE file at the root of this repository.

// Package purelib is a pure widget library: it imports the toolkit and ships no
// view-model. Checked with -requirevm=false, the missing-view-model rule is off,
// so a package that only constructs widgets produces no diagnostics.
package purelib

import "github.com/go-widgets/toolkit"

// Build returns a fresh entry; it never mutates a guarded state field.
func Build() *toolkit.Entry {
	return toolkit.NewEntry()
}
