package main

import (
	"go/ast"
	"go/token"
	"strconv"
	"strings"
)

func scanP0(f *ast.File, imports map[string]string) map[string]bool {
	found := map[string]bool{}
	ast.Inspect(f, func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.GenDecl:
			if n.Tok == token.CONST {
				found["const"] = true
				ast.Inspect(n, func(node ast.Node) bool {
					if id, ok := node.(*ast.Ident); ok && id.Name == "iota" {
						found["iota"] = true
					}
					return true
				})
			}
		case *ast.DeferStmt:
			found["defer"] = true
			if lit, ok := n.Call.Fun.(*ast.FuncLit); ok {
				ast.Inspect(lit.Body, func(node ast.Node) bool {
					if call, ok := node.(*ast.CallExpr); ok {
						if id, ok := call.Fun.(*ast.Ident); ok && id.Name == "recover" && id.Obj == nil {
							found["recover"] = true
						}
					}
					return true
				})
			}
		case *ast.RangeStmt:
			if stringExpr(n.X) {
				found["rune"] = true
			}
		case *ast.SelectorExpr:
			if pkg, ok := n.X.(*ast.Ident); ok && packageName(imports, pkg) == "utf8" {
				found["rune"] = true
			}
		case *ast.StructType:
			for _, field := range n.Fields.List {
				if len(field.Names) == 0 {
					found["embed"] = true
				}
			}
		case *ast.TypeSpec:
			if _, ok := n.Type.(*ast.InterfaceType); ok {
				found["iface"] = true
			}
		case *ast.Field:
			if _, ok := n.Type.(*ast.FuncType); ok {
				found["funcvalue"] = true
			}
			if _, ok := n.Type.(*ast.Ellipsis); ok {
				found["variadic"] = true
			}
		case *ast.MapType:
			if _, ok := n.Value.(*ast.FuncType); ok {
				found["funcvalue"] = true
			}
		case *ast.AssignStmt:
			for _, rhs := range n.Rhs {
				if _, ok := rhs.(*ast.FuncLit); ok {
					found["funcvalue"] = true
				}
			}
		case *ast.ValueSpec:
			for _, rhs := range n.Values {
				if _, ok := rhs.(*ast.FuncLit); ok {
					found["funcvalue"] = true
				}
			}
		case *ast.FuncDecl:
			scanFuncP0(n, imports, found)
		case *ast.CallExpr:
			if at, ok := n.Fun.(*ast.ArrayType); ok && at.Len == nil {
				if id, ok := at.Elt.(*ast.Ident); ok && id.Name == "rune" {
					found["rune"] = true
				}
			}
			if id, ok := n.Fun.(*ast.Ident); ok && id.Obj == nil {
				if id.Name == "close" {
					found["close"] = true
				}
				if id.Name == "make" && len(n.Args) == 2 {
					if _, ok := n.Args[0].(*ast.ChanType); ok {
						found["bufchan"] = true
					}
				}
			}
			if sel, ok := n.Fun.(*ast.SelectorExpr); ok {
				if pkg, ok := sel.X.(*ast.Ident); ok {
					name := packageName(imports, pkg)
					if name == "fmt" && sel.Sel.Name == "Errorf" && len(n.Args) > 0 {
						if lit, ok := n.Args[0].(*ast.BasicLit); ok {
							format, _ := strconv.Unquote(lit.Value)
							if hasWrapVerb(format) {
								found["wrap"] = true
							}
						}
					}
					if name == "errors" && (sel.Sel.Name == "Is" || sel.Sel.Name == "As") {
						found["unwrap"] = true
					}
					if sel.Sel.Name == "Run" && testingPointer(pkg, "T", imports) {
						found["subtest"] = true
					}
				}
			}
		}
		return true
	})
	found["errwrap"] = found["wrap"] && found["unwrap"]
	return found
}

func scanFuncP0(fn *ast.FuncDecl, imports map[string]string, found map[string]bool) {
	if fn.Recv != nil && fn.Name.Name == "Error" && fn.Type.Params.NumFields() == 0 && fn.Type.Results.NumFields() == 1 {
		if id, ok := fn.Type.Results.List[0].Type.(*ast.Ident); ok && id.Name == "string" {
			found["errtype"] = true
		}
	}
	if strings.HasPrefix(fn.Name.Name, "Fuzz") && fn.Type.Params.NumFields() == 1 {
		if testingType(fn.Type.Params.List[0].Type, "F", imports) {
			found["fuzz"] = true
		}
	}
	returnsFunc := false
	if fn.Type.Results != nil {
		for _, field := range fn.Type.Results.List {
			if _, ok := field.Type.(*ast.FuncType); ok {
				returnsFunc = true
			}
		}
	}
	ast.Inspect(fn.Body, func(node ast.Node) bool {
		if call, ok := node.(*ast.CallExpr); ok {
			if id, ok := call.Fun.(*ast.Ident); ok && id.Obj != nil && id.Obj == fn.Name.Obj {
				found["recursion"] = true
			}
		}
		if _, ok := node.(*ast.FuncLit); ok && returnsFunc {
			found["closure"] = true
		}
		return true
	})
}

func testingPointer(id *ast.Ident, name string, imports map[string]string) bool {
	if id.Obj == nil {
		return false
	}
	field, ok := id.Obj.Decl.(*ast.Field)
	return ok && testingType(field.Type, name, imports)
}

func testingType(expr ast.Expr, name string, imports map[string]string) bool {
	star, ok := expr.(*ast.StarExpr)
	if !ok {
		return false
	}
	sel, ok := star.X.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	pkg, ok := sel.X.(*ast.Ident)
	return ok && packageName(imports, pkg) == "testing" && sel.Sel.Name == name
}

func stringExpr(expr ast.Expr) bool {
	if lit, ok := expr.(*ast.BasicLit); ok {
		return lit.Kind == token.STRING
	}
	if call, ok := expr.(*ast.CallExpr); ok {
		id, ok := call.Fun.(*ast.Ident)
		return ok && id.Name == "string"
	}
	id, ok := expr.(*ast.Ident)
	if !ok || id.Obj == nil {
		return false
	}
	switch decl := id.Obj.Decl.(type) {
	case *ast.Field:
		typ, ok := decl.Type.(*ast.Ident)
		return ok && typ.Name == "string"
	case *ast.ValueSpec:
		if typ, ok := decl.Type.(*ast.Ident); ok && typ.Name == "string" {
			return true
		}
		for i, n := range decl.Names {
			if n.Name == id.Name && i < len(decl.Values) {
				return stringExpr(decl.Values[i])
			}
		}
	case *ast.AssignStmt:
		for i, lhs := range decl.Lhs {
			if n, ok := lhs.(*ast.Ident); ok && n.Name == id.Name && i < len(decl.Rhs) {
				if rhs, ok := decl.Rhs[i].(*ast.BasicLit); ok {
					return rhs.Kind == token.STRING
				}
			}
		}
	}
	return false
}

func hasWrapVerb(format string) bool {
	for i := 0; i < len(format)-1; i++ {
		if format[i] == '%' {
			i++
			if format[i] == 'w' {
				return true
			}
		}
	}
	return false
}
