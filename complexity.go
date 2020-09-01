package complexity

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = "complexity is ..."

// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name: "complexity",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.FuncDecl:
			cnt := cntFuncBranch(n)
			pass.Reportf(n.Pos(), "branch cnt: %d", cnt)
		}
	})

	return nil, nil
}

func cntFuncBranch(decl *ast.FuncDecl) int {
	ast.Print(nil, decl)
	return cntBranch(decl.Body)
}

func cntBranch(stmt *ast.BlockStmt) int {
	if stmt == nil {
		return 0
	}

	cnt := 0
	for _, s := range stmt.List {
		switch s := s.(type) {
		case *ast.BlockStmt:
			cnt += cntBranch(s)
		case *ast.IfStmt:
			cnt += cntIfBranch(s)
		case *ast.ForStmt:
			cnt += cntForBranch(s)
		}
	}
	return cnt
}

func cntIfBranch(stmt *ast.IfStmt) int {
	cnt := 1
	switch s := stmt.Else.(type) {
	case *ast.IfStmt: // else if
		cnt += cntIfBranch(s)
	case *ast.BlockStmt: // only else
		cnt += 1 + cntBranch(s)
	}
	return cnt
}

func cntForBranch(stmt *ast.ForStmt) int {
	return 1 + cntBranch(stmt.Body)
}
