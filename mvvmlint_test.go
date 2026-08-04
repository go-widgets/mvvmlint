// Copyright (c) 2026 the go-widgets/mvvmlint authors. All rights reserved.
// Use of this source code is governed by a BSD-3-Clause license that can be
// found in the LICENSE file at the root of this repository.

package mvvmlint

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

// setFlag sets an analyzer flag for the duration of a subtest and restores it.
func setFlag(t *testing.T, name, value string) {
	t.Helper()
	f := Analyzer.Flags.Lookup(name)
	if f == nil {
		t.Fatalf("flag %q not found", name)
	}
	prev := f.Value.String()
	if err := Analyzer.Flags.Set(name, value); err != nil {
		t.Fatalf("set -%s=%s: %v", name, value, err)
	}
	t.Cleanup(func() { _ = Analyzer.Flags.Set(name, prev) })
}

// TestDefaults exercises both rules with the built-in configuration:
//   - direct:      flagged direct assignments (and un-flagged callback/composite/non-toolkit writes),
//   - bound:       the sanctioned &w.Field bind — no diagnostics,
//   - nomvvm:      the missing-view-model diagnostic,
//   - allowdir:    the //mvvmlint:allow directive exemption,
//   - bindingfile: the *_binding.go exemption is file-scoped, not package-scoped.
func TestDefaults(t *testing.T) {
	dir := analysistest.TestData()
	analysistest.Run(t, dir, Analyzer,
		"direct", "bound", "nomvvm", "allowdir", "bindingfile", "notoolkit",
		// *_test.go is skipped by default: skipctx (toolkit in prod, a mutation in
		// a test file) and tkonlytest (toolkit imported only from a test file) must
		// both emit nothing.
		"skipctx", "tkonlytest")
}

// TestIncludeTests opts into analyzing *_test.go files, so a direct mutation in a
// test fixture is flagged.
func TestIncludeTests(t *testing.T) {
	setFlag(t, "includetests", "true")
	analysistest.Run(t, analysistest.TestData(), Analyzer, "incltests")
}

// TestStateFieldsFlag overrides the allowlist so only Items/Selected are guarded
// and a Text write is allowed, exercising the -statefields flag path.
func TestStateFieldsFlag(t *testing.T) {
	setFlag(t, "statefields", "Items,,Selected")
	analysistest.Run(t, analysistest.TestData(), Analyzer, "statefieldsflag")
}

// TestRequireVMOff turns the missing-view-model rule off for a pure widget
// library, which must then produce no diagnostics.
func TestRequireVMOff(t *testing.T) {
	setFlag(t, "requirevm", "false")
	analysistest.Run(t, analysistest.TestData(), Analyzer, "purelib")
}

// TestExportedKnobs guards the exported constants and default allowlist against
// accidental drift.
func TestExportedKnobs(t *testing.T) {
	if ToolkitPath != "github.com/go-widgets/toolkit" {
		t.Errorf("ToolkitPath = %q", ToolkitPath)
	}
	if MVVMPath != "github.com/go-widgets/mvvm" {
		t.Errorf("MVVMPath = %q", MVVMPath)
	}
	want := map[string]bool{"Text": true, "Items": true, "Selected": true, "Rows": true, "Root": true}
	have := make(map[string]bool)
	for _, f := range DefaultStateFields {
		have[f] = true
	}
	for k := range want {
		if !have[k] {
			t.Errorf("DefaultStateFields missing %q", k)
		}
	}
}
