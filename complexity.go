package complexity

import (
	"flag"
	"fmt"
	"go/ast"
	// "math"

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
	over int
)

func init() {
	flag.IntVar(&over, "over", 10, "show functions with Cyclomatic complexity > k")
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.FuncDecl:
			operators, operands := walkFuncDecl(n)
			fmt.Println(operators, operands)
			// dist_opt := len(operators) // distinct operators
			// dist_opd := len(operands)  // distrinct operands
			// var sum_opt, sum_opd int
			// for _, val := range operators {
			// 	sum_opt += val
			// }
			//
			// for _, val := range operands {
			// 	sum_opd += val
			// }
			//
			// n_vocab := dist_opt + dist_opd
			// length := sum_opt + sum_opd
			// volume := float64(length) * math.Log2(float64(n_vocab))
			// fmt.Println(operands)
			// difficulty := float64(dist_opt * sum_opd / (2 * dist_opd))
			//
			// pass.Reportf(n.Pos(), "Cyclomatic complexity: %f, %f", volume, difficulty)
		}
	})

	return nil, nil
}

type branchVisitor func(n ast.Node) (w ast.Visitor)

// Visit is ...
func (v branchVisitor) Visit(n ast.Node) (w ast.Visitor) {
	return v(n)
}

// walkFuncDecl counts Cyclomatic complexity
func walkFuncDecl(fd *ast.FuncDecl) (map[string]int, map[string]int) {
	operators, operands := map[string]int{}, map[string]int{}

	var v ast.Visitor
	v = branchVisitor(func(n ast.Node) (w ast.Visitor) {
		switch n := n.(type) {
		case *ast.FuncDecl:
			break
		case *ast.BasicLit:
			operators[n.Value]++
		case *ast.FuncLit:
			fmt.Println(n)
		case *ast.Ident:
			// if object type is 'var'
			if n.Obj != nil && n.Obj.Kind == 4 {
				operands[n.Name]++
			} else if n.Obj == nil {
				// add 'print' and 'println'
				operators[n.Name]++
				operators["()"]++
			}
		}
		return v
	})
	ast.Walk(v, fd)

	print(operators, operands)
	return operators, operands
}
