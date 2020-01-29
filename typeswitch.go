package typeswitch

import (
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name: "typeswitch",
	Doc:  Doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

const Doc = "typeswitch finds a type which implement an interfaces which are used in type-switch but the type does not appear in any cases of the type-switch"

type enum struct {
	Interface  *types.Interface
	Implements []types.Object
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// find enum-like interfaces
	enums := getEnums(pass.Pkg)
	for _, p := range pass.Pkg.Imports() {
		for k, v := range getEnums(p) {
			enums[k] = v
		}
	}

	nodeFilter := []ast.Node{
		(*ast.TypeSwitchStmt)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		sw, ok := n.(*ast.TypeSwitchStmt)
		if !ok {
			return
		}

		var typ types.Type
		switch stmt := sw.Assign.(type) {
		case *ast.ExprStmt:
			if expr, ok := stmt.X.(*ast.TypeAssertExpr); ok {
				typ = pass.TypesInfo.TypeOf(expr.X).Underlying()
			}
		case *ast.AssignStmt:
			if expr, ok := stmt.Rhs[0].(*ast.TypeAssertExpr); ok {
				typ = pass.TypesInfo.TypeOf(expr.X).Underlying()
			}
		default:
			panic("unexpected type")
		}

		e := enums[typ]
		if e == nil {
			return
		}

		var ids []string
		for _, obj := range e.Implements {
			if !hasCase(pass, obj.Type(), sw) {
				ids = append(ids, obj.Id())
			}
		}

		if len(ids) != 0 {
			pass.Reportf(sw.Pos(), "type %s does not appear in any cases", strings.Join(ids, ","))
		}
	})

	return nil, nil
}

func getEnums(pkg *types.Package) map[types.Type]*enum {
	var itfs []*types.Interface
	var typs []types.Object

	// find interfaces
	for _, n := range pkg.Scope().Names() {
		obj, ok := pkg.Scope().Lookup(n).(*types.TypeName)
		if !ok {
			continue
		}

		if itf, ok := obj.Type().Underlying().(*types.Interface); ok {
			itfs = append(itfs, itf)
		} else {
			typs = append(typs, obj)
		}
	}

	// find implements
	enums := map[types.Type]*enum{}
	for _, itf := range itfs {
		e := &enum{
			Interface: itf,
		}

		for _, typ := range typs {
			if types.Implements(typ.Type(), itf) {
				e.Implements = append(e.Implements, typ)
			}
		}

		if len(e.Implements) >= 2 {
			enums[itf] = e
		}
	}

	return enums
}

func hasCase(pass *analysis.Pass, t types.Type, sw *ast.TypeSwitchStmt) bool {
	for _, s := range sw.Body.List {
		c, ok := s.(*ast.CaseClause)
		if !ok {
			continue
		}

		// default
		if c.List == nil {
			return true
		}

		for _, expr := range c.List {
			if types.Identical(t, pass.TypesInfo.TypeOf(expr)) {
				return true
			}
		}
	}
	return false
}
