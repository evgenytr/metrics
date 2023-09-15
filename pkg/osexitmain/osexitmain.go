// Package osexitmain contains checker to detect os.Exit direct calls in main function of package main
package osexitmain

import (
	"go/ast"
	"go/types"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"regexp"
)

var generatedCodeRegexp = regexp.MustCompile(`^// Code generated .* DO NOT EDIT\.$`)
var Analyzer = &analysis.Analyzer{
	Name:     "osexitmain",
	Doc:      "checker to detect os.Exit direct calls in main function of package main",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	// check only main package
	if !isMainPackage(pass.Pkg) {
		return nil, nil
	}

	for _, file := range pass.Files {
		//skip generated files
		skipFile := false
		for _, cg := range file.Comments {
			for _, comment := range cg.List {
				if generatedCodeRegexp.MatchString(comment.Text) {
					skipFile = true
				}
			}
		}

		if skipFile {
			continue
		}

		ast.Inspect(file, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.FuncDecl:
				if !isMainFunc(x) {
					return false
				}
			case *ast.SelectorExpr:
				if isOsExitCall(x) {
					pass.Reportf(x.Pos(), "os.Exit direct call detected in main function of main package")
				}
			}
			return true
		})
	}

	return nil, nil
}

func isMainPackage(pkg *types.Package) bool {
	return pkg.Name() == "main"
}

func isMainFunc(currFunc *ast.FuncDecl) bool {
	return currFunc.Name.Name == "main"
}

func isOsExitCall(currExpr *ast.SelectorExpr) bool {
	if currExpr.X == nil {
		return false
	}

	if ident, ok := currExpr.X.(*ast.Ident); ok {
		if ident.Name == "os" && currExpr.Sel.Name == "Exit" {
			return true
		}
	}

	return false
}
