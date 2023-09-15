package osexitmain

import (
	"go/ast"
	"go/types"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

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

	// check only func main
	if !isMainPackage(pass.Pkg) {
		return nil, nil
	}

	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeTypes := []ast.Node{
		(*ast.BinaryExpr)(nil),
	}

	inspect.Preorder(nodeTypes, check(pass))

	return nil, nil
}

func isMainPackage(pkg *types.Package) bool {
	return pkg.Name() == "main"
}
