// Copyright (c) 2026 the go-widgets/mvvmlint authors. All rights reserved.
// Use of this source code is governed by a BSD-3-Clause license that can be
// found in the LICENSE file at the root of this repository.

// Package nomvvm imports the toolkit but ships no MVVM view-model, which the
// missing-view-model rule reports on the toolkit import.
package nomvvm

import (
	"github.com/go-widgets/toolkit" // want `widgets present but no MVVM view-model; wire state through github.com/go-widgets/mvvm`
)

func build() *toolkit.Entry {
	return toolkit.NewEntry()
}
