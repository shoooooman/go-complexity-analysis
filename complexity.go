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
				// fmt.Println("cyclo", cycloComp, pass.Pkg.Name(), n.Name)
				pass.Reportf(n.Pos(), "Cyclomatic complexity: %d", cycloComp)
			}

			volume := calcHalstComp(n)

			loc := countLOC(pass.Fset, n)
			maintIdx := calcMaintIndex(volume, cycloComp, loc)
			if maintIdx < maintunder {
				// fmt.Println("maint", maintIdx, pass.Pkg.Name(), n.Name)
				pass.Reportf(n.Pos(), "Maintainability index: %d", maintIdx)
			}

			// pass.Reportf(n.Pos(), "Cyclomatic complexity: %d", cycloComp)
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

	walkDecl(fd, operators, operands)

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

func walkDecl(n ast.Node, opt map[string]int, opd map[string]int) {
	switch n := n.(type) {
	case *ast.GenDecl:
		appendValidSymb(n.Lparen.IsValid(), n.Rparen.IsValid(), opt, "()")

		if n.Tok.IsOperator() {
			opt[n.Tok.String()]++
		} else {
			opd[n.Tok.String()]++
		}
		for _, s := range n.Specs {
			walkSpec(s, opt, opd)
		}
	case *ast.FuncDecl:
		if n.Recv == nil {
			opt["func"]++
			opt[n.Name.Name]++
			opt["()"]++
		}
		walkStmt(n.Body, opt, opd)
	}
}

func walkStmt(n ast.Node, opt map[string]int, opd map[string]int) {
	switch n := n.(type) {
	case *ast.DeclStmt:
		walkDecl(n.Decl, opt, opd)
	// TODO(sff1019): Add EmptyStmt
	// TODO(sff1019): Add LabeledStmt
	case *ast.ExprStmt:
		walkExpr(n.X, opt, opd)
	case *ast.SendStmt:
		walkExpr(n.Chan, opt, opd)
		if n.Arrow.IsValid() {
			opt["<-"]++
		}
		walkExpr(n.Value, opt, opd)
	case *ast.IncDecStmt:
		walkExpr(n.X, opt, opd)
		if n.Tok.IsOperator() {
			opt[n.Tok.String()]++
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
	case *ast.GoStmt:
		if n.Go.IsValid() {
			opt["go"]++
		}
		walkExpr(n.Call, opt, opd)
	case *ast.DeferStmt:
		if n.Defer.IsValid() {
			opt["defer"]++
		}
		walkExpr(n.Call, opt, opd)
	case *ast.ReturnStmt:
		if n.Return.IsValid() {
			opt["return"]++
		}
		for _, e := range n.Results {
			walkExpr(e, opt, opd)
		}
	case *ast.BranchStmt:
		if n.Tok.IsOperator() {
			opt[n.Tok.String()]++
		} else {
			opd[n.Tok.String()]++
		}
		if n.Label != nil {
			walkExpr(n.Label, opt, opd)
		}
	case *ast.BlockStmt:
		appendValidSymb(n.Lbrace.IsValid(), n.Rbrace.IsValid(), opt, "{}")
		for _, s := range n.List {
			walkStmt(s, opt, opd)
		}
	case *ast.IfStmt:
		if n.If.IsValid() {
			opt["if"]++
		}
		if n.Init != nil {
			walkStmt(n.Init, opt, opd)
		}
		walkExpr(n.Cond, opt, opd)
		walkStmt(n.Body, opt, opd)
		if n.Else != nil {
			opt["else"]++
			walkStmt(n.Else, opt, opd)
		}
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
	case *ast.SelectStmt:
		if n.Select.IsValid() {
			opt["select"]++
		}
		walkStmt(n.Body, opt, opd)
	case *ast.ForStmt:
		if n.For.IsValid() {
			opt["for"]++
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
	case *ast.RangeStmt:
		if n.For.IsValid() {
			opt["for"]++
		}
		if n.Key != nil {
			walkExpr(n.Key, opt, opd)
			if n.Tok.IsOperator() {
				opt[n.Tok.String()]++
			} else {
				opd[n.Tok.String()]++
			}
		}
		if n.Value != nil {
			walkExpr(n.Value, opt, opd)
		}
		opt["range"]++
		walkExpr(n.X, opt, opd)
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

func walkSpec(spec ast.Spec, opt map[string]int, opd map[string]int) {
	switch spec := spec.(type) {
	case *ast.ValueSpec:
		for _, n := range spec.Names {
			walkExpr(n, opt, opd)
			if spec.Type != nil {
				walkExpr(spec.Type, opt, opd)
			}
			if spec.Values != nil {
				for _, v := range spec.Values {
					walkExpr(v, opt, opd)
				}
			}
		}
	}
}

func walkExpr(exp ast.Expr, opt map[string]int, opd map[string]int) {
	switch exp := exp.(type) {
	case *ast.ParenExpr:
		appendValidSymb(exp.Lparen.IsValid(), exp.Rparen.IsValid(), opt, "()")
		walkExpr(exp.X, opt, opd)
	case *ast.SelectorExpr:
		walkExpr(exp.X, opt, opd)
		walkExpr(exp.Sel, opt, opd)
	case *ast.IndexExpr:
		walkExpr(exp.X, opt, opd)
		appendValidSymb(exp.Lbrack.IsValid(), exp.Rbrack.IsValid(), opt, "{}")
		walkExpr(exp.Index, opt, opd)
	case *ast.SliceExpr:
		walkExpr(exp.X, opt, opd)
		appendValidSymb(exp.Lbrack.IsValid(), exp.Rbrack.IsValid(), opt, "[]")
		if exp.Low != nil {
			walkExpr(exp.Low, opt, opd)
		}
		if exp.High != nil {
			walkExpr(exp.High, opt, opd)
		}
		if exp.Max != nil {
			walkExpr(exp.Max, opt, opd)
		}
	case *ast.TypeAssertExpr:
		walkExpr(exp.X, opt, opd)
		appendValidSymb(exp.Lparen.IsValid(), exp.Rparen.IsValid(), opt, "()")
		if exp.Type != nil {
			walkExpr(exp.Type, opt, opd)
		}
	case *ast.CallExpr:
		walkExpr(exp.Fun, opt, opd)
		appendValidSymb(exp.Lparen.IsValid(), exp.Rparen.IsValid(), opt, "()")
		if exp.Ellipsis != 0 {
			opt["..."]++
		}
		for _, a := range exp.Args {
			walkExpr(a, opt, opd)
		}
	case *ast.StarExpr:
		if exp.Star.IsValid() {
			opt["*"]++
		}
		walkExpr(exp.X, opt, opd)
	case *ast.UnaryExpr:
		if exp.Op.IsOperator() {
			opt[exp.Op.String()]++
		} else {
			opd[exp.Op.String()]++
		}
		walkExpr(exp.X, opt, opd)
	case *ast.BinaryExpr:
		walkExpr(exp.X, opt, opd)
		opt[exp.Op.String()]++
		walkExpr(exp.Y, opt, opd)
	case *ast.KeyValueExpr:
		walkExpr(exp.Key, opt, opd)
		if exp.Colon.IsValid() {
			opt[":"]++
		}
		walkExpr(exp.Value, opt, opd)
	case *ast.BasicLit:
		if exp.Kind.IsLiteral() {
			opd[exp.Value]++
		} else {
			opt[exp.Value]++
		}
	case *ast.FuncLit:
		walkStmt(exp.Body, opt, opd)
	case *ast.CompositeLit:
		appendValidSymb(exp.Lbrace.IsValid(), exp.Rbrace.IsValid(), opt, "{}")
		if exp.Type != nil {
			walkExpr(exp.Type, opt, opd)
		}
		for _, e := range exp.Elts {
			walkExpr(e, opt, opd)
		}
	case *ast.Ident:
		if exp.Obj == nil {
			opt[exp.Name]++
		} else {
			opd[exp.Name]++
		}
	case *ast.Ellipsis:
		if exp.Ellipsis.IsValid() {
			opt["..."]++
		}
		if exp.Elt != nil {
			walkExpr(exp.Elt, opt, opd)
		}
	}
}

func appendValidSymb(lvalid bool, rvalid bool, opt map[string]int, symb string) {
	if lvalid && rvalid {
		opt[symb]++
	}

	n := 0
	switch n {
	case 0:
	case 1:
	case 2:
	case 3:
	case 4:
	case 5:
	case 6:
	case 7:
	case 8:
	case 9:
	case 10:
	case 11:
	}
}
