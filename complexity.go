package complexity

import (
	"flag"
	"fmt"
	"go/ast"
	"go/token"
	"math"

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

var (
	cycloover  int
	maintunder int
)

func init() {
	flag.IntVar(&cycloover, "cycloover", 10, "show functions with the Cyclomatic complexity > N")
	flag.IntVar(&maintunder, "maintunder", 20, "show functions with the Maintainability index < N")
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.FuncDecl:
			cycloComp := calcCycloComp(n)
			if cycloComp > cycloover {
				fmt.Println("cyclo", cycloComp, pass.Pkg.Name(), n.Name)
			}

			// FIXME: mock
			halstComp := calcHalstComp()

			loc := countLOC(pass.Fset, n)
			maintIdx := calcMaintIndex(halstComp, cycloComp, loc)
			if maintIdx < maintunder {
				fmt.Println("maint", maintIdx, pass.Pkg.Name(), n.Name)
			}

			pass.Reportf(n.Pos(), "Cyclomatic complexity: %d", cycloComp)
		}
	})

	return nil, nil
}

type branchVisitor func(n ast.Node) (w ast.Visitor)

// Visit is ...
func (v branchVisitor) Visit(n ast.Node) (w ast.Visitor) {
	return v(n)
}

// calcCycloComp calculates the Cyclomatic complexity
func calcCycloComp(fd *ast.FuncDecl) int {
	// ast.Print(nil, fd)

	comp := 1
	var v ast.Visitor
	v = branchVisitor(func(n ast.Node) (w ast.Visitor) {
		switch n := n.(type) {
		case *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt, *ast.CaseClause, *ast.CommClause:
			comp++
		case *ast.BinaryExpr:
			if n.Op == token.LAND || n.Op == token.LOR {
				comp++
			}
		}
		return v
	})
	ast.Walk(v, fd)

	return comp
}

// FIXME: mock
func calcHalstComp() int {
	return 1
}

// counts lines of a function
func countLOC(fs *token.FileSet, n *ast.FuncDecl) int {
	f := fs.File(n.Pos())
	startLine := f.Line(n.Pos())
	endLine := f.Line(n.End())
	return endLine - startLine + 1
}

// calcMaintComp calculates the maintainability index
// source: https://docs.microsoft.com/ja-jp/archive/blogs/codeanalysis/maintainability-index-range-and-meaning
func calcMaintIndex(halstComp, cycloComp, loc int) int {
	origVal := 171.0 - 5.2*math.Log(float64(halstComp)) - 0.23*float64(cycloComp) - 16.2*math.Log(float64(loc))
	normVal := int(math.Max(0.0, origVal*100.0/171.0))
	return normVal
}
