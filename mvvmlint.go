// Copyright (c) 2026 the go-widgets/mvvmlint authors. All rights reserved.
// Use of this source code is governed by a BSD-3-Clause license that can be
// found in the LICENSE file at the root of this repository.

// Package mvvmlint provides a go/analysis analyzer that makes MVVM usage
// mandatory for applications built on the go-widgets pixel/cell toolkit.
//
// The toolkit exposes mutable value fields on its widgets (a text entry's
// Text, a list's Items, a table's Rows, and so on). Mutating those fields ad
// hoc scatters view state across imperative call sites and defeats the point of
// a ViewModel layer. The sanctioned way to touch widget state is through the
// go-widgets/mvvm binders, which take a POINTER into the field
// (BindField(obs, &w.Text, &w.OnChange, invalidate)); the address-of form is a
// binding, not a mutation, and is never flagged.
//
// The analyzer reports two rules, but only in packages that import the toolkit:
//
//  1. Direct widget-state mutation. An assignment whose left-hand side is a
//     selector X.Field, where the static type of X is (a pointer to) a toolkit
//     widget type and Field is in the configurable state-field allowlist. This
//     is the core "mutate ad hoc" violation. Composite literals and the
//     address-of form &w.Field are not assignments to a selector and are never
//     flagged. Assignments in files named *_binding.go, and any file carrying a
//     //mvvmlint:allow directive, are exempt so intentional glue can opt out.
//
//  2. Missing view-model. A package that imports the toolkit but does not
//     import go-widgets/mvvm gets a single diagnostic asking that state be wired
//     through a view-model. This rule can be turned off (via -requirevm=false)
//     for pure widget libraries that legitimately ship no ViewModel.
package mvvmlint

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// ToolkitPath is the import path of the go-widgets pixel toolkit.
const ToolkitPath = "github.com/go-widgets/toolkit"

// MVVMPath is the import path of the go-widgets MVVM library.
const MVVMPath = "github.com/go-widgets/mvvm"

// allowDirective, when present as a comment in a source file, exempts that file
// from the direct-mutation rule.
const allowDirective = "//mvvmlint:allow"

// DefaultStateFields is the built-in allowlist of toolkit widget value fields
// that carry display or data state. Assigning to one of these on a toolkit
// widget is the "mutate ad hoc" violation the analyzer exists to catch.
//
// The set is intentionally curated to value/state fields (not callback slots
// such as OnChange, and not one-time construction options): it is derived from
// the mutable state fields of the toolkit widgets — an entry's Text, a list's
// Items and Selected, a table's Rows/Columns, a tree's Root, a slider's
// Value/Low/High, a check's Checked, a switch's On, and so on.
var DefaultStateFields = []string{
	"Active",
	"Checked",
	"Columns",
	"Content",
	"Current",
	"Cursor",
	"Data",
	"Expanded",
	"Fraction",
	"High",
	"Items",
	"Low",
	"Max",
	"Min",
	"Nodes",
	"On",
	"Options",
	"Page",
	"Progress",
	"Range",
	"Root",
	"Rows",
	"ScrollRow",
	"Selected",
	"Text",
	"Value",
	"Values",
}

// config holds the analyzer's tunable behaviour. It is stored on the Analyzer so
// that flags and programmatic construction share one code path.
type config struct {
	stateFieldsFlag string // raw -statefields flag value ("" => defaults)
	requireVM       bool   // emit the missing-view-model rule
	toolkitPath     string // toolkit import path (overridable for testing)
	mvvmPath        string // mvvm import path (overridable for testing)
}

// Analyzer is the mvvmlint go/analysis analyzer. Use it with singlechecker, or
// as a go vet tool: go vet -vettool=$(which mvvmlint) ./...
var Analyzer = newAnalyzer()

func newAnalyzer() *analysis.Analyzer {
	cfg := &config{requireVM: true, toolkitPath: ToolkitPath, mvvmPath: MVVMPath}
	a := &analysis.Analyzer{
		Name:     "mvvmlint",
		Doc:      "enforce that go-widgets apps route widget state through go-widgets/mvvm instead of mutating widget fields ad hoc",
		URL:      "https://github.com/go-widgets/mvvmlint",
		Requires: []*analysis.Analyzer{inspect.Analyzer},
		Run:      func(pass *analysis.Pass) (any, error) { return run(pass, cfg) },
	}
	a.Flags.StringVar(&cfg.stateFieldsFlag, "statefields", "",
		"comma-separated allowlist of widget state fields to guard (empty uses the built-in default)")
	a.Flags.BoolVar(&cfg.requireVM, "requirevm", true,
		"report packages that import the toolkit but no MVVM view-model (set false for pure widget libraries)")
	a.Flags.StringVar(&cfg.toolkitPath, "toolkitpath", ToolkitPath,
		"import path treated as the widget toolkit")
	a.Flags.StringVar(&cfg.mvvmPath, "mvvmpath", MVVMPath,
		"import path treated as the MVVM library")
	return a
}

// stateFieldSet resolves the effective allowlist for a pass: the -statefields
// flag when set, otherwise the built-in default.
func (c *config) stateFieldSet() map[string]bool {
	names := DefaultStateFields
	if strings.TrimSpace(c.stateFieldsFlag) != "" {
		names = splitFields(c.stateFieldsFlag)
	}
	set := make(map[string]bool, len(names))
	for _, n := range names {
		set[n] = true
	}
	return set
}

// splitFields splits a comma-separated flag value into trimmed, non-empty names.
func splitFields(s string) []string {
	var out []string
	for _, part := range strings.Split(s, ",") {
		if p := strings.TrimSpace(part); p != "" {
			out = append(out, p)
		}
	}
	return out
}

func run(pass *analysis.Pass, cfg *config) (any, error) {
	// The analyzer only concerns packages that use the widget toolkit. The
	// toolkit import spec's position doubles as the "is it imported?" signal and
	// as the anchor for the missing-view-model diagnostic (a single, stable
	// source location).
	tkPos := toolkitImportPos(pass, cfg.toolkitPath)
	if tkPos == token.NoPos {
		return nil, nil
	}

	// Rule 2: toolkit present but no MVVM view-model.
	if cfg.requireVM && !importsPath(pass.Pkg, cfg.mvvmPath) {
		pass.Reportf(tkPos,
			"widgets present but no MVVM view-model; wire state through %s", cfg.mvvmPath)
	}

	// Rule 1: direct widget-state mutation.
	fields := cfg.stateFieldSet()
	exempt := exemptFiles(pass)
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	insp.Preorder([]ast.Node{(*ast.AssignStmt)(nil)}, func(n ast.Node) {
		assign := n.(*ast.AssignStmt)
		if exempt[pass.Fset.File(assign.Pos()).Name()] {
			return
		}
		for _, lhs := range assign.Lhs {
			sel, ok := lhs.(*ast.SelectorExpr)
			if !ok {
				continue
			}
			if !fields[sel.Sel.Name] {
				continue
			}
			if !isToolkitWidget(pass.TypesInfo.TypeOf(sel.X), cfg.toolkitPath) {
				continue
			}
			pass.Reportf(sel.Pos(),
				"direct mutation of toolkit widget state field %q; bind it through %s (e.g. BindField(obs, &w.%s, ...)) instead",
				sel.Sel.Name, cfg.mvvmPath, sel.Sel.Name)
		}
	})
	return nil, nil
}

// importsPath reports whether pkg directly imports the given path.
func importsPath(pkg *types.Package, path string) bool {
	for _, imp := range pkg.Imports() {
		if imp.Path() == path {
			return true
		}
	}
	return false
}

// toolkitImportPos returns the position of the first toolkit import spec across
// the pass's files, or token.NoPos if the package does not import the toolkit.
// pass.Files is ordered by position, so the first match is also the earliest,
// giving the missing-view-model diagnostic a single, stable, source-anchored
// location.
func toolkitImportPos(pass *analysis.Pass, toolkitPath string) token.Pos {
	for _, f := range pass.Files {
		for _, imp := range f.Imports {
			if strings.Trim(imp.Path.Value, `"`) == toolkitPath {
				return imp.Pos()
			}
		}
	}
	return token.NoPos
}

// exemptFiles returns the set of file names (as reported by token.File.Name)
// that are exempt from the direct-mutation rule: files named *_binding.go, and
// files carrying a //mvvmlint:allow directive.
func exemptFiles(pass *analysis.Pass) map[string]bool {
	exempt := make(map[string]bool)
	for _, f := range pass.Files {
		name := pass.Fset.File(f.Pos()).Name()
		if strings.HasSuffix(name, "_binding.go") {
			exempt[name] = true
			continue
		}
		if hasAllowDirective(f) {
			exempt[name] = true
		}
	}
	return exempt
}

// hasAllowDirective reports whether a file contains a //mvvmlint:allow comment.
func hasAllowDirective(f *ast.File) bool {
	for _, group := range f.Comments {
		for _, c := range group.List {
			if strings.TrimSpace(c.Text) == allowDirective {
				return true
			}
		}
	}
	return false
}

// isToolkitWidget reports whether t is (a pointer to) a named type declared in
// the toolkit package — i.e. a widget whose state fields the rule guards.
func isToolkitWidget(t types.Type, toolkitPath string) bool {
	// A nil type asserts to false below, so no explicit nil guard is needed.
	if ptr, ok := t.(*types.Pointer); ok {
		t = ptr.Elem()
	}
	named, ok := t.(*types.Named)
	if !ok {
		return false
	}
	// named.Obj() is never nil; its package is nil only for universe types,
	// which have no assignable state fields — pkg != nil short-circuits those.
	pkg := named.Obj().Pkg()
	return pkg != nil && pkg.Path() == toolkitPath
}
