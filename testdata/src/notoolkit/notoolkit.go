// Copyright (c) 2026 the go-widgets/mvvmlint authors. All rights reserved.
// Use of this source code is governed by a BSD-3-Clause license that can be
// found in the LICENSE file at the root of this repository.

// Package notoolkit imports no widget toolkit, so the analyzer ignores it
// entirely — even a same-named Text field on an unrelated type is untouched.
package notoolkit

type box struct{ Text string }

func run() {
	b := box{}
	b.Text = "not a widget"
	_ = b
}
