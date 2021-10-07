package kys

import (
	"fmt"
	"go/ast"

	comments "github.com/bragov4ik/go-kys/pkg/comments"
	cyclo "github.com/bragov4ik/go-kys/pkg/cyclocomp"
	halstead "github.com/bragov4ik/go-kys/pkg/halstead"
	"github.com/k0kubun/pp/v3"
)

type Info struct {
	CycloComp  uint
	FuncLit    uint
	ReturnStmt uint
	CallExpr   uint
	AssignStmt uint
	Comments   uint
	Halstead   halstead.Info
}

func NewInfo() Info {
	return Info{Halstead: halstead.NewInfo()}
}

type Config struct {
	CycloComp cyclo.Weights    `xml:"cyclomatic"`
	Comment   comments.Weights `xml:"comment"`
}

func parseNode(n ast.Node, info *Info, cfg *Config) {
	switch v := n.(type) {
	case *ast.ArrayType:
		info.Halstead.AddArrayType(v)
	case *ast.AssignStmt:
		info.Halstead.AddAssignStmt(v)
		info.AssignStmt++
	case *ast.BadDecl:
		info.Halstead.AddBadDecl(v)
	case *ast.BadExpr:
		info.Halstead.AddBadExpr(v)
	case *ast.BadStmt:
		info.Halstead.AddBadStmt(v)
	case *ast.BasicLit:
		info.Halstead.AddBasicLit(v)
	case *ast.BinaryExpr:
		info.Halstead.AddBinaryExpr(v)
	case *ast.BlockStmt:
		// Should we parse blocks?
		info.Halstead.AddBlockStmt(v)
	case *ast.BranchStmt:
		info.Halstead.AddBranchStmt(v)
	case *ast.CallExpr:
		info.Halstead.AddCallExpr(v)
		info.CallExpr++
	case *ast.CaseClause:
		info.Halstead.AddCaseClause(v)
	case *ast.ChanType:
		info.Halstead.AddChanType(v)
	case *ast.CommClause:
		info.Halstead.AddCommClause(v)
	case *ast.Comment:
		info.Halstead.AddComment(v)
		info.Comments += comments.GetCommentComp(v, &cfg.Comment)
	case *ast.CommentGroup:
		info.Halstead.AddCommentGroup(v)
	case *ast.CompositeLit:
		info.Halstead.AddCompositeLit(v)
	case *ast.DeclStmt:
		info.Halstead.AddDeclStmt(v)
	case *ast.DeferStmt:
		info.Halstead.AddDeferStmt(v)
	case *ast.Ellipsis:
		info.Halstead.AddEllipsis(v)
	case *ast.EmptyStmt:
		info.Halstead.AddEmptyStmt(v)
	case *ast.ExprStmt:
		info.Halstead.AddExprStmt(v)
	case *ast.Field:
		info.Halstead.AddField(v)
	case *ast.FieldList:
		info.Halstead.AddFieldList(v)
	case *ast.File:
		info.Halstead.AddFile(v)
	case *ast.ForStmt:
		info.Halstead.AddForStmt(v)
	case *ast.FuncDecl:
		info.Halstead.AddFuncDecl(v)
		info.CycloComp += cyclo.GetCycloComp(v, &cfg.CycloComp)
	case *ast.FuncLit:
		info.Halstead.AddFuncLit(v)
		info.FuncLit++
	case *ast.FuncType:
		info.Halstead.AddFuncType(v)
	case *ast.GenDecl:
		info.Halstead.AddGenDecl(v)
	case *ast.GoStmt:
		info.Halstead.AddGoStmt(v)
	case *ast.Ident:
		info.Halstead.AddIdent(v)
	case *ast.IfStmt:
		info.Halstead.AddIfStmt(v)
	case *ast.ImportSpec:
		info.Halstead.AddImportSpec(v)
	case *ast.IncDecStmt:
		info.Halstead.AddIncDecStmt(v)
	case *ast.IndexExpr:
		info.Halstead.AddIndexExpr(v)
	case *ast.InterfaceType:
		info.Halstead.AddInterfaceType(v)
	case *ast.KeyValueExpr:
		info.Halstead.AddKeyValueExpr(v)
	case *ast.LabeledStmt:
		info.Halstead.AddLabeledStmt(v)
	case *ast.MapType:
		info.Halstead.AddMapType(v)
	case *ast.Package:
		info.Halstead.AddPackage(v)
	case *ast.ParenExpr:
		info.Halstead.AddParenExpr(v)
	case *ast.RangeStmt:
		info.Halstead.AddRangeStmt(v)
	case *ast.ReturnStmt:
		info.Halstead.AddReturnStmt(v)
		info.ReturnStmt++
	case *ast.SelectStmt:
		info.Halstead.AddSelectStmt(v)
	case *ast.SelectorExpr:
		info.Halstead.AddSelectorExpr(v)
	case *ast.SendStmt:
		info.Halstead.AddSendStmt(v)
	case *ast.SliceExpr:
		info.Halstead.AddSliceExpr(v)
	case *ast.StarExpr:
		info.Halstead.AddStarExpr(v)
	case *ast.StructType:
		info.Halstead.AddStructType(v)
	case *ast.SwitchStmt:
		info.Halstead.AddSwitchStmt(v)
	case *ast.TypeAssertExpr:
		info.Halstead.AddTypeAssertExpr(v)
	case *ast.TypeSpec:
		info.Halstead.AddTypeSpec(v)
	case *ast.TypeSwitchStmt:
		info.Halstead.AddTypeSwitchStmt(v)
	case *ast.UnaryExpr:
		info.Halstead.AddUnaryExpr(v)
	case *ast.ValueSpec:
		info.Halstead.AddValueSpec(v)

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
