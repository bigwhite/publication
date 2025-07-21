package checkpubfuncname

import (
	"go/ast"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const Doc = `check for top-level function names that do not start with an uppercase letter (not exported) in non-main packages.

This analyzer helps enforce a convention that all top-level functions in library packages
should be exported if they are intended for external use, or kept unexported (lowercase)
if they are internal helpers. This specific check flags functions that might have been
intended to be package-private but were accidentally named with a non-uppercase first letter,
or vice-versa if the policy was to export all top-level funcs.`

// Analyzer is the instance of our custom analyzer.
var Analyzer = &analysis.Analyzer{
	Name:     "checkpubfuncname",
	Doc:      Doc,
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer}, // We need the inspector pass
	// ResultType: // Not producing any facts or results for other analyzers
	// FactTypes:  // Not using facts
}

func run(pass *analysis.Pass) (interface{}, error) {
	// Skip "main" package, as main.main is an exception, and other funcs might be internal.
	// This is a simplistic check; real linters have more sophisticated ways to handle package types.
	if pass.Pkg.Name() == "main" {
		return nil, nil
	}

	// Get the inspector. This is provided by the inspect.Analyzer requirement.
	inspectResult := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// We are interested in function declarations (ast.FuncDecl)
	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	inspectResult.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl, ok := n.(*ast.FuncDecl)
		if !ok {
			return
		}

		// We are interested in top-level functions in the package.
		if funcDecl.Recv == nil { // It's a function, not a method
			funcName := funcDecl.Name.Name

			// Skip special function "init"
			if funcName == "init" {
				return
			}

			// Skip test functions (TestXxx, BenchmarkXxx, ExampleXxx)
			if strings.HasPrefix(funcName, "Test") ||
				strings.HasPrefix(funcName, "Benchmark") ||
				strings.HasPrefix(funcName, "Example") {
				return
			}

			if len(funcName) > 0 {
				firstChar := rune(funcName[0])
				if !unicode.IsUpper(firstChar) {
					// This is a top-level function in a non-main package,
					// and its name starts with a lowercase letter.
					pass.Reportf(funcDecl.Pos(), "top-level function '%s' in package '%s' is not exported (name starts with lowercase)", funcName, pass.Pkg.Name())
				}
			}
		}
	})

	return nil, nil // No result for other analyzers, no error
}
