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

type Config struct {
	cycloCompWeights CycloCompWeights
}
type CycloCompWeights struct {
	IF   uint
	FOR  uint
	RNG  uint
	CASE uint
	AND  uint
	OR   uint
}

func parseNode(n ast.Node, info *Info, config *Config) {
	switch v := n.(type) {
	case *ast.FuncDecl:
		info.CycloComp += calcCycloComp(v, config)
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

func GetInfo(file *ast.File, info *Info, config *Config) {
	ast.Inspect(file, func(n ast.Node) bool {
		parseNode(n, info, config)
		return true
	})
}

type branchVisitor func(n ast.Node) (w ast.Visitor)

func (v branchVisitor) Visit(n ast.Node) (w ast.Visitor) {
	return v(n)
}

func calcCycloComp(fd *ast.FuncDecl, config *Config) uint {
	var comp uint = 1
	var v ast.Visitor
	v = branchVisitor(func(n ast.Node) (w ast.Visitor) {
		switch n := n.(type) {
		case *ast.IfStmt:
			comp += config.cycloCompWeights.IF
		case *ast.ForStmt:
			comp += config.cycloCompWeights.FOR
		case *ast.RangeStmt:
			comp += config.cycloCompWeights.RNG
		case *ast.CaseClause:
			comp += config.cycloCompWeights.CASE
		case *ast.CommClause:
			comp += config.cycloCompWeights.CASE
		case *ast.BinaryExpr:
			if n.Op == token.LAND {
				comp += config.cycloCompWeights.AND
			} else if n.Op == token.LOR {
				comp += config.cycloCompWeights.OR
			}
		}
		return v
	})
	ast.Walk(v, fd)

	return comp
}
