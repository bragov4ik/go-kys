package kys

import (
	"go/ast"

	"github.com/k0kubun/pp"
)

type Info struct {
	FuncDecl   uint
	FuncLit    uint
	ReturnStmt uint
	CallExpr   uint
	AssignStmt uint
}

func parseNode(n ast.Node, info *Info) {
	switch v := n.(type) {
	case *ast.FuncDecl:
		info.FuncDecl++
	case *ast.FuncLit:
		info.FuncLit++
	case *ast.ReturnStmt:
		info.ReturnStmt++
	case *ast.CallExpr:
		info.CallExpr++
	case *ast.AssignStmt:
		info.AssignStmt++

	case *ast.BlockStmt:
		// Should we parse blocks?
	case *ast.Ident:
	case *ast.FieldList:
	case *ast.Field:
	case *ast.FuncType:
	case *ast.SelectorExpr:
	case *ast.ExprStmt:
	case *ast.ImportSpec:
	case *ast.BasicLit:
	case *ast.File:
	default:
		if n != nil {
			pp.Printf("Unhandled type: %T %v\n", v, n)
		}
	}
}

func GetInfo(file *ast.File, info *Info) {
	ast.Inspect(file, func(n ast.Node) bool {
		parseNode(n, info)
		return true
	})
}
