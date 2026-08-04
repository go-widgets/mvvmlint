// Copyright (c) 2026 the go-widgets/mvvmlint authors. All rights reserved.
// Use of this source code is governed by a BSD-3-Clause license that can be
// found in the LICENSE file at the root of this repository.

// Package toolkit is a minimal stand-in for github.com/go-widgets/toolkit used
// only by the analyzer's analysistest fixtures. It declares just enough widget
// shape (a couple of state fields and a callback slot) to exercise the rules.
package toolkit

// Entry is a single-line text widget.
type Entry struct {
	Text     string
	OnChange func(string)
}

// ListBox is a vertical list widget.
type ListBox struct {
	Items    []string
	Selected int
}

// NewEntry constructs an Entry, mirroring the real toolkit's constructor style.
func NewEntry() *Entry { return &Entry{} }

// NewListBox constructs a ListBox.
func NewListBox() *ListBox { return &ListBox{} }
