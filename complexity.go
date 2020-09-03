package complexity

import (
	"flag"
	"fmt"
	"math"

	"go/ast"
	"go/token"

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

			volume := calcHalstComp(n)

			loc := countLOC(pass.Fset, n)
			maintIdx := calcMaintIndex(volume, cycloComp, loc)
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

func calcHalstComp(fd *ast.FuncDecl) float64 {
	operators, operands := map[string]int{}, map[string]int{}

	var v ast.Visitor
	v = branchVisitor(func(n ast.Node) (w ast.Visitor) {
		walkStmt(n, operators, operands)
		return v
	})
	ast.Walk(v, fd)

	fmt.Println(operators, operands)
	distOpt := len(operators) // distinct operators
	distOpd := len(operands)  // distrinct operands
	var sumOpt, sumOpd int
	for _, val := range operators {
		sumOpt += val
	}

	for _, val := range operands {
		sumOpd += val
	}

	nVocab := distOpt + distOpd
	length := sumOpt + sumOpd
	volume := float64(length) * math.Log2(float64(nVocab))
	difficulty := float64(distOpt*sumOpd) / float64(2*distOpd)
	fmt.Println("difficulty", difficulty)

	return volume
}

// counts lines of a function
func countLOC(fs *token.FileSet, n *ast.FuncDecl) int {
	f := fs.File(n.Pos())
	startLine := f.Line(n.Pos())
	endLine := f.Line(n.End())
	return endLine - startLine + 1
}

// calcMaintComp calculates the maintainability index
// source: https://docs.microsoft.com/en-us/archive/blogs/codeanalysis/maintainability-index-range-and-meaning
func calcMaintIndex(halstComp float64, cycloComp, loc int) int {
	origVal := 171.0 - 5.2*math.Log(halstComp) - 0.23*float64(cycloComp) - 16.2*math.Log(float64(loc))
	normVal := int(math.Max(0.0, origVal*100.0/171.0))
	return normVal
}

func walkStmt(n ast.Node, opt map[string]int, opd map[string]int) {
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
		for _, exp := range n.Rhs {
			walkExpr(exp, opt, opd)
		}
	case *ast.ExprStmt:
		walkExpr(n.X, opt, opd)
	case *ast.IfStmt:
		if n.If.IsValid() {
			opt["if"]++
			opt["{}"]++
		}
		if n.Init != nil {
			walkStmt(n.Init, opt, opd)
		}
		walkExpr(n.Cond, opt, opd)
		walkStmt(n.Body, opt, opd)
		if n.Else != nil {
			opt["else"]++
			opt["{}"]++
			walkStmt(n.Else, opt, opd)
		}
	case *ast.ForStmt:
		if n.For.IsValid() {
			opt["for"]++
			opt["{}"]++
		}
		if n.Init != nil {
			walkStmt(n.Init, opt, opd)
		}
		if n.Cond != nil {
			walkExpr(n.Cond, opt, opd)
		}
		if n.Post != nil {
			walkStmt(n.Post, opt, opd)
		}
		walkStmt(n.Body, opt, opd)
	case *ast.SwitchStmt:
		if n.Switch.IsValid() {
			opt["switch"]++
		}
		if n.Init != nil {
			walkStmt(n.Init, opt, opd)
		}
		if n.Tag != nil {
			walkExpr(n.Tag, opt, opd)
		}
		walkStmt(n.Body, opt, opd)
	case *ast.CaseClause:
		if n.List == nil {
			opt["default"]++
		} else {
			for _, c := range n.List {
				walkExpr(c, opt, opd)
			}
		}
		if n.Colon.IsValid() {
			opt[":"]++
		}
		if n.Body != nil {
			for _, b := range n.Body {
				walkStmt(b, opt, opd)
			}

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
	case *ast.BasicLit:
		if exp.Kind.IsLiteral() {
			opd[exp.Value]++
		} else {
			opt[exp.Value]++
		}
	case *ast.BinaryExpr:
		walkExpr(exp.X, opt, opd)
		opt[exp.Op.String()]++
		walkExpr(exp.Y, opt, opd)
	case *ast.ParenExpr:
		appendValidParen(exp.Lparen.IsValid(), exp.Rparen.IsValid(), opt)
		walkExpr(exp.X, opt, opd)
	case *ast.CallExpr:
		walkExpr(exp.Fun, opt, opd)
		appendValidParen(exp.Lparen.IsValid(), exp.Rparen.IsValid(), opt)
		for _, ea := range exp.Args {
			walkExpr(ea, opt, opd)
		}
	}
}

func appendValidParen(lvalid bool, rvalid bool, opt map[string]int) {
	if lvalid && rvalid {
		opt["()"]++
	}
}
