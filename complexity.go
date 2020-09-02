package complexity

import (
	"flag"
	"fmt"
	"go/ast"
	"math"
	"reflect"

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
			dist_opt := len(operators) // distinct operators
			dist_opd := len(operands)  // distrinct operands
			var sum_opt, sum_opd int
			for _, val := range operators {
				sum_opt += val
			}

			for _, val := range operands {
				sum_opd += val
			}

			n_vocab := dist_opt + dist_opd
			length := sum_opt + sum_opd
			fmt.Println("n1", dist_opt, "n2", dist_opd)
			fmt.Println("N1", sum_opt, "N2", sum_opd)
			volume := float64(length) * math.Log2(float64(n_vocab))
			difficulty := float64(dist_opt*sum_opd) / float64(2*dist_opd)

			pass.Reportf(n.Pos(), "Halstead complexity: %f:3, %f", volume, difficulty)
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
		walkOps(n, operators, operands)
		return v
	})
	ast.Walk(v, fd)

	return operators, operands
}

func walkOps(n ast.Node, opt map[string]int, opd map[string]int) {
	// if n != nil {
	// 	fmt.Println(reflect.ValueOf(n).Elem(), reflect.ValueOf(n).Elem().Type())
	// }
	switch n := n.(type) {
	case *ast.FuncDecl:
		if n.Recv == nil {
			opt["func"]++
			opt[n.Name.Name]++
			opt["()"]++
			opt["{}"]++
		}
	case *ast.AssignStmt:
		if n.Tok.IsOperator() {
			opt[n.Tok.String()]++
		}
		for _, exp := range n.Lhs {
			walkExpr(exp, opt, opd)
		}
	case *ast.ExprStmt:
		walkExpr(n.X, opt, opd)
	case *ast.BasicLit:
		if n.Kind.IsLiteral() {
			opd[n.Value]++
		} else {
			opt[n.Value]++
		}

	}
}

func walkExpr(exp ast.Expr, opt map[string]int, opd map[string]int) {
	switch exp := exp.(type) {
	case *ast.Ident:
		if exp.Obj == nil {
			opt[exp.Name]++
		} else {
			opd[exp.Name]++
		}
	case *ast.CallExpr:
		walkExpr(exp.Fun, opt, opd)
		if exp.Lparen.IsValid() && exp.Rparen.IsValid() {
			opt["()"]++
		}
		fmt.Println(exp.Args)
	}
	// switch n := n.(type) {
	//   case
	//
	// }
}
