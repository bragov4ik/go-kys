package halstead

import (
	"fmt"
	"go/ast"
	"go/token"
	"math"
)

type Info struct {
	operators map[string]uint
	operands  map[string]uint
}

func NewInfo() Info {
	return Info{
		operators: make(map[string]uint),
		operands:  make(map[string]uint),
	}
}

func (info *Info) getN1Distinct() uint {
	return uint(len(info.operators))
}

func (info *Info) getN2Distinct() uint {
	return uint(len(info.operands))
}

// Not universal key type because using interfaces is nasty
// and generics (with 1.18 version) are not released yet
func sumMap(targetMap map[string]uint) uint {
	var total uint = 0
	for _, count := range targetMap {
		total += count
	}
	return total
}

func tokenInArr(tok token.Token, arr []token.Token) bool {
	for _, op := range arr {
		if op == tok {
			return true
		}
	}
	return false
}

func (info *Info) getN1Total() uint {
	return sumMap(info.operands)
}

func (info *Info) getN2Total() uint {
	return sumMap(info.operators)
}

func (info *Info) Vocabulary() uint {
	return info.getN1Distinct() + info.getN2Distinct()
}

func (info *Info) Length() uint {
	return info.getN1Total() + info.getN2Total()
}

func (info *Info) Volume() float64 {
	nTot := info.Length()
	nDist := info.Vocabulary()
	return float64(nTot) * math.Log2(float64(nDist))
}

func (info *Info) Difficulty() float64 {
	n1Dist, n2Dist := info.getN1Distinct(), info.getN2Distinct()
	n2Tot := info.getN2Total()
	return (float64(n1Dist) * float64(n2Tot)) / (2 * float64(n2Dist))
}

func (info *Info) Effort() float64 {
	// Do not use other functions to avoid summing maps multiple times
	n1Dist, n2Dist := float64(info.getN1Distinct()), float64(info.getN2Distinct())
	n1Tot, n2Tot := float64(info.getN1Total()), float64(info.getN2Total())
	D := (n1Dist * n2Tot) / (2 * n2Dist)
	V := (n1Tot + n2Tot) * math.Log2(n1Dist+n2Dist)
	return D * V
}

func (info *Info) addToken(nextToken token.Token) {
	tokenName := fmt.Sprintf("token:%s", nextToken.String())

	NOT_OPERATORS := []token.Token{
		token.ILLEGAL,
		token.EOF,
		token.COMMENT,
		token.IDENT,
		token.INT,
		token.FLOAT,
		token.IMAG,
		token.CHAR,
		token.STRING,
	}

	if !tokenInArr(nextToken, NOT_OPERATORS[:]) {
		info.operators[tokenName] += 1
	}
}

// TODO count commas

func (info *Info) AddArrayType(node *ast.ArrayType) {
	// Nothing to count
}

func (info *Info) AddAssignStmt(node *ast.AssignStmt) {
	// lhs and rhs should be visited in walk (expr contains node)
	info.addToken(node.Tok)
}

func (info *Info) AddBadDecl(node *ast.BadDecl) {
	// Nothing to count
}

func (info *Info) AddBadExpr(node *ast.BadExpr) {
	// Nothing to count
}

func (info *Info) AddBadStmt(node *ast.BadStmt) {
	// Nothing to count
}

func (info *Info) AddBasicLit(node *ast.BasicLit) {
	info.operands[node.Value] += 1
}

func (info *Info) AddBinaryExpr(node *ast.BinaryExpr) {
	// x and y should be visited in walk (expr contains node)
	info.addToken(node.Op)
}

func (info *Info) AddBlockStmt(node *ast.BlockStmt) {
	// statements should be visited in walk
	info.addToken(token.LBRACE)
}

func (info *Info) AddBranchStmt(node *ast.BranchStmt) {
	info.addToken(node.Tok)
}

func (info *Info) AddCallExpr(node *ast.CallExpr) {
	// leave expressions for further walk
	// add only one of the parentheses as they are in pairs
	info.addToken(token.LPAREN)
	// Ellipsis are handled separately
}

func (info *Info) AddCaseClause(node *ast.CaseClause) {
	// no way to distinguish case and default, so just always assume case for now
	info.addToken(token.CASE)
	info.addToken(token.COLON)
}

func (info *Info) AddChanType(node *ast.ChanType) {
	info.addToken(token.ARROW)
	if node.Arrow != token.NoPos {
		info.addToken(token.CHAN)
	}
}

func (info *Info) AddCommClause(node *ast.CommClause) {
	// no way to distinguish case and default, so just always assume case for now
	info.addToken(token.CASE)
	info.addToken(token.COLON)
}

func (info *Info) AddComment(node *ast.Comment) {
	// Nothing to count
}

func (info *Info) AddCommentGroup(node *ast.CommentGroup) {
	// Nothing to count
}

func (info *Info) AddCompositeLit(node *ast.CompositeLit) {
	// Maybe also consider `Elts` (e.g. if it always results in commas added)
	info.addToken(token.LPAREN)
}

func (info *Info) AddDeclStmt(node *ast.DeclStmt) {
	// GenDecl is handled separately
}

func (info *Info) AddDeferStmt(node *ast.DeferStmt) {
	info.addToken(token.DEFER)
}

func (info *Info) AddEllipsis(node *ast.Ellipsis) {
	info.addToken(token.ELLIPSIS)
}

func (info *Info) AddEmptyStmt(node *ast.EmptyStmt) {
	if !node.Implicit {
		info.addToken(token.SEMICOLON)
	}
}

func (info *Info) AddExprStmt(node *ast.ExprStmt) {
	// Nothing to count
}

func (info *Info) AddField(node *ast.Field) {
	// Nothing to count
}

func (info *Info) AddFieldList(node *ast.FieldList) {
	if node.Opening.IsValid() {
		info.addToken(token.LPAREN)
	}
}

func (info *Info) AddFile(node *ast.File) {
	info.addToken(token.PACKAGE)
}

func (info *Info) AddForStmt(node *ast.ForStmt) {
	info.addToken(token.FOR)
}

func (info *Info) AddFuncDecl(node *ast.FuncDecl) {
	// Composite type only, nothing
}

func (info *Info) AddFuncLit(node *ast.FuncLit) {
	// Composite type only, nothing
}

func (info *Info) AddFuncType(node *ast.FuncType) {
	info.addToken(token.FUNC)
}

func (info *Info) AddGenDecl(node *ast.GenDecl) {
	info.addToken(node.Tok)
	if node.Lparen.IsValid() {
		info.addToken(token.LPAREN)
	}
}

func (info *Info) AddGoStmt(node *ast.GoStmt) {
	info.addToken(token.GO)
}

func (info *Info) AddIdent(node *ast.Ident) {
	info.operands[node.Name] += 1
}

func (info *Info) AddIfStmt(node *ast.IfStmt) {
	info.addToken(token.IF)
}

func (info *Info) AddImportSpec(node *ast.ImportSpec) {
	// Composite type only, nothing
}

func (info *Info) AddIncDecStmt(node *ast.IncDecStmt) {
	info.addToken(node.Tok)
}

func (info *Info) AddIndexExpr(node *ast.IndexExpr) {
	if node.Lbrack.IsValid() {
		info.addToken(token.LBRACK)
	}
}

func (info *Info) AddInterfaceType(node *ast.InterfaceType) {
	info.addToken(token.INTERFACE)
}

func (info *Info) AddKeyValueExpr(node *ast.KeyValueExpr) {
	info.addToken(token.COLON)
}

func (info *Info) AddLabeledStmt(node *ast.LabeledStmt) {
	info.addToken(token.COLON)
}

func (info *Info) AddMapType(node *ast.MapType) {
	info.addToken(token.MAP)
}

func (info *Info) AddPackage(node *ast.Package) {
	// name should be handled by ident, so ignore is as apparently it is
	// not a particular part of code but rather abstract entity (set of files?).
}

func (info *Info) AddParenExpr(node *ast.ParenExpr) {
	info.addToken(token.LPAREN)
}

func (info *Info) AddRangeStmt(node *ast.RangeStmt) {
	// for should be already handled by `AddForStmt`
	info.addToken(node.Tok)
	info.addToken(token.RANGE)
}

func (info *Info) AddReturnStmt(node *ast.ReturnStmt) {
	info.addToken(token.RETURN)
}

func (info *Info) AddSelectStmt(node *ast.SelectStmt) {
	info.addToken(token.SELECT)
}

func (info *Info) AddSelectorExpr(node *ast.SelectorExpr) {
	info.addToken(token.PERIOD)
}

func (info *Info) AddSendStmt(node *ast.SendStmt) {
	info.addToken(token.ARROW)
}

func (info *Info) AddSliceExpr(node *ast.SliceExpr) {
	info.addToken(token.LBRACK)
}

func (info *Info) AddStarExpr(node *ast.StarExpr) {
	info.addToken(token.MUL)
}

func (info *Info) AddStructType(node *ast.StructType) {
	info.addToken(token.STRUCT)
}

func (info *Info) AddSwitchStmt(node *ast.SwitchStmt) {
	info.addToken(token.SWITCH)
}

func (info *Info) AddTypeAssertExpr(node *ast.TypeAssertExpr) {
	info.addToken(token.LPAREN)
}

func (info *Info) AddTypeSpec(node *ast.TypeSpec) {
	if node.Assign.IsValid() {
		info.addToken(token.ASSIGN)
	}
}

func (info *Info) AddTypeSwitchStmt(node *ast.TypeSwitchStmt) {
	info.addToken(token.SWITCH)
}

func (info *Info) AddUnaryExpr(node *ast.UnaryExpr) {
	info.addToken(node.Op)
}

func (info *Info) AddValueSpec(node *ast.ValueSpec) {
	// Something composite only, ignore
}
