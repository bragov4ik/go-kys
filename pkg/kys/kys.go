package kys

import (
	"go/ast"
	"go/token"

	"github.com/k0kubun/pp"
)

type Info struct {
	CycloComp  uint
	FuncLit    uint
	ReturnStmt uint
	CallExpr   uint
	AssignStmt uint
}

func parseNode(n ast.Node, info *Info) {
	switch v := n.(type) {
	case *ast.FuncDecl:
		info.CycloComp += calcCycloComp(v)
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

type branchVisitor func(n ast.Node) (w ast.Visitor)

func (v branchVisitor) Visit(n ast.Node) (w ast.Visitor) {
	return v(n)
}

func calcCycloComp(fd *ast.FuncDecl) uint {
	var comp uint = 1
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
