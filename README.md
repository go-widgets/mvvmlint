# mvvmlint

[![ci](https://github.com/go-widgets/mvvmlint/actions/workflows/ci.yml/badge.svg)](https://github.com/go-widgets/mvvmlint/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/go-widgets/mvvmlint.svg)](https://pkg.go.dev/github.com/go-widgets/mvvmlint)
[![License: BSD-3-Clause](https://img.shields.io/badge/License-BSD--3--Clause-blue.svg)](LICENSE)

A [`go/analysis`](https://pkg.go.dev/golang.org/x/tools/go/analysis) analyzer that
makes **MVVM usage mandatory** for applications built on the go-widgets toolkit.

The toolkit exposes mutable value fields on its widgets — a text entry's `Text`,
a list's `Items`, a table's `Rows`, a tree's `Root`. Mutating those fields ad hoc
scatters view state across imperative call sites and defeats the point of a
ViewModel layer. The sanctioned way to touch widget state is through the
[`go-widgets/mvvm`](https://github.com/go-widgets/mvvm) binders, which take a
**pointer into the field**:

```go
mvvm.BindField(obs, &entry.Text, &entry.OnChange, invalidate) // sanctioned
entry.Text = "typed ad hoc"                                   // flagged
```

Run it in each app's CI as a **required status check** so non-MVVM code fails to
merge.

## Rules

The analyzer only fires in packages that import `github.com/go-widgets/toolkit`.

1. **Direct widget-state mutation.** An assignment whose left-hand side is a
   selector `X.Field`, where the static type of `X` is (a pointer to) a toolkit
   widget type **and** `Field` is in the state-field allowlist. This is the core
   "mutate ad hoc" violation.
   - The address-of form `&w.Field` (what the binders take) is **not** an
     assignment to a selector and is never flagged.
   - A composite literal `toolkit.Entry{Text: "x"}` seeds a field at
     construction and is never flagged.
   - Callback slots (`OnChange`, `OnClick`, …) are **not** state fields, so
     composing a handler is never flagged.
2. **Missing view-model.** A package that imports the toolkit but **not**
   `go-widgets/mvvm` earns one diagnostic, anchored on the toolkit import:
   *"widgets present but no MVVM view-model; wire state through
   github.com/go-widgets/mvvm"*. Turn it off with `-requirevm=false` for pure
   widget libraries that legitimately ship no ViewModel.

### Default state-field allowlist

Derived from the mutable state fields of the toolkit widgets:

```text
Active  Checked  Columns  Content  Current  Cursor  Data  Expanded  Fraction
High    Items    Low      Max      Min      Nodes   On    Options   Page
Progress Range   Root     Rows     ScrollRow Selected Text Value    Values
```

Override it with `-statefields=Text,Items,Selected,...` (comma-separated). When
set, the flag **replaces** the default set.

### Exemptions (opt out of rule 1)

- Files named `*_binding.go` — the conventional home for hand-written binding
  glue — are exempt.
- Any file containing a `//mvvmlint:allow` directive is exempt.

## Install & run locally

```sh
go install github.com/go-widgets/mvvmlint/cmd/mvvmlint@latest
go vet -vettool="$(go env GOPATH)/bin/mvvmlint" ./...
```

Flags are passed through `go vet` (unprefixed, since this is a single-analyzer
vet tool):

```sh
go vet -vettool="$(which mvvmlint)" -requirevm=false ./...
go vet -vettool="$(which mvvmlint)" -statefields=Text,Items ./...
```

## CI gate

### Reusable workflow (recommended)

Reference the shipped reusable workflow from your app repo and mark the resulting
check as **required** in branch protection:

```yaml
# .github/workflows/mvvm.yml in your app repo
name: mvvm
on:
  pull_request:
  push:
    branches: [main]
jobs:
  mvvm:
    uses: go-widgets/mvvmlint/.github/workflows/mvvmlint.yml@main
    # optional overrides:
    # with:
    #   go-version: "1.26.4"
    #   packages: "./..."
    #   version: "latest"
```

### Standalone job (copy-paste)

For repos that prefer inlining the steps:

```yaml
jobs:
  mvvm:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.26.4"
      - run: go install github.com/go-widgets/mvvmlint/cmd/mvvmlint@latest
      - run: go vet -vettool="$(go env GOPATH)/bin/mvvmlint" ./...
```

## False-positive posture

The rule is deliberately conservative — it fires only when **both** the static
receiver type resolves to a toolkit widget **and** the field is on the allowlist,
so a same-named field on an unrelated type (`myConfig.Text = ...`) is never
touched. The remaining intentional edges — low-level binding glue, or a widget
subsystem that must poke fields directly — are handled by the `*_binding.go`
naming convention and the `//mvvmlint:allow` directive, and rule 2 can be
disabled wholesale for pure widget libraries with `-requirevm=false`.

## License

BSD-3-Clause — see [LICENSE](LICENSE). Copyright (c) 2026 the go-widgets/mvvmlint
authors.
