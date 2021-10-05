package kys

import (
	"fmt"
	"go/ast"

	comments "github.com/bragov4ik/go-kys/pkg/comments"
	cyclo "github.com/bragov4ik/go-kys/pkg/cyclocomp"
	"github.com/k0kubun/pp/v3"
)

type Info struct {
	CycloComp  uint
	FuncLit    uint
	ReturnStmt uint
	CallExpr   uint
	AssignStmt uint
	Comments   uint
}

type Config struct {
	CycloComp cyclo.Weights    `xml:"cyclomatic"`
	Comment   comments.Weights `xml:"comment"`
}

func parseNode(n ast.Node, info *Info, cfg *Config) {
	switch v := n.(type) {
	case *ast.FuncDecl:
		info.CycloComp += cyclo.GetCycloComp(v, &cfg.CycloComp)
	case *ast.FuncLit:
		info.FuncLit++
	case *ast.ReturnStmt:
		info.ReturnStmt++
	case *ast.CallExpr:
		info.CallExpr++
	case *ast.AssignStmt:
		info.AssignStmt++
	case *ast.Comment:
		info.Comments += comments.GetCommentComp(v, &cfg.Comment)

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
			fmt.Printf("Unhandled type: %T ", v)
			pp.Printf("%v\n", n)
		}
	}
}

func GetInfo(file *ast.File, info *Info, config *Config) {
	ast.Inspect(file, func(n ast.Node) bool {
		parseNode(n, info, config)
		return true
	})
}
