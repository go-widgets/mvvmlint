// Copyright (c) 2026 the go-widgets/mvvmlint authors. All rights reserved.
// Use of this source code is governed by a BSD-3-Clause license that can be
// found in the LICENSE file at the root of this repository.

// Package mvvm is a minimal stand-in for github.com/go-widgets/mvvm used only by
// the analyzer's analysistest fixtures. It provides just enough of the binder
// surface (an Observable and BindField) for a fixture to bind through a pointer.
package mvvm

// Observable is a trivial value holder.
type Observable[T any] struct{ v T }

// NewObservable seeds an Observable.
func NewObservable[T any](v T) *Observable[T] { return &Observable[T]{v: v} }

// Get returns the current value.
func (o *Observable[T]) Get() T { return o.v }

// Set stores a value.
func (o *Observable[T]) Set(v T) { o.v = v }

// BindField binds obs two-way to a widget value field and its callback slot via
// pointers — the sanctioned way to touch widget state.
func BindField[T any](obs *Observable[T], field *T, hook *func(T), invalidate func()) {
	*field = obs.Get()
	_ = hook
	_ = invalidate
}
