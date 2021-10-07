package halstead

import (
	"fmt"
	"go/ast"
	"go/token"
	"math"
)

type HalsteadInfo struct {
	operators map[string]uint
	operands  map[string]uint
}

func NewHalsteadInfo() *HalsteadInfo {
	return &HalsteadInfo{
		operators: make(map[string]uint),
		operands: make(map[string]uint),
	}
}

func (info *HalsteadInfo) String() string {
	return fmt.Sprintf("%v, %v", info.operators, info.operands)
}

func (info *HalsteadInfo) getN1Distinct() uint {
	return uint(len(info.operators))
}

func (info *HalsteadInfo) getN2Distinct() uint {
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

func (info *HalsteadInfo) getN1Total() uint {
	sum := sumMap(info.operands)
	return sum
}

func (info *HalsteadInfo) getN2Total() uint {
	sum := sumMap(info.operators)
	return sum
}

func (info *HalsteadInfo) Vocabuary() uint {
	n1Dist, n2Dist := info.getN1Distinct(), info.getN2Distinct()
	return n1Dist + n2Dist
}

func (info *HalsteadInfo) Length() uint {
	n1Tot, n2Tot := info.getN1Total(), info.getN2Total()
	return n1Tot + n2Tot
}

func (info *HalsteadInfo) Volume() float64 {
	nTot := info.Length()
	nDist := info.Vocabuary()
	return float64(nTot) * math.Log2(float64(nDist))
}

func (info *HalsteadInfo) Difficulty() float64 {
	n1Dist, n2Dist := info.getN1Distinct(), info.getN2Distinct()
	n2Tot := info.getN2Total()
	return (float64(n1Dist) * float64(n2Tot)) / (2 * float64(n2Dist))
}

func (info *HalsteadInfo) Effort() float64 {
	// Do not use other functions to avoid summing maps multiple times
	n1Dist, n2Dist := float64(info.getN1Distinct()), float64(info.getN2Distinct())
	n1Tot, n2Tot := float64(info.getN1Total()), float64(info.getN2Total())
	D := (n1Dist * n2Tot) / (2 * n2Dist)
	V := (n1Tot + n2Tot) * math.Log2(n1Dist+n2Dist)
	return D * V
}

func (info *HalsteadInfo) addToken(nextToken token.Token) {
	tokenName := fmt.Sprintf("token:%s", nextToken.String())

	NOT_OPERATORS := []token.Token{token.ILLEGAL, token.EOF, token.COMMENT, token.IDENT, token.INT, token.FLOAT, token.IMAG, token.CHAR, token.STRING}
	
	if !tokenInArr(nextToken, NOT_OPERATORS[:]) {
		info.operators[tokenName] += 1
	}
}

// TODO count commas

func (info *HalsteadInfo) AddArrayType(node *ast.ArrayType) {
	// Nothing to count
	return
}

func (info *HalsteadInfo) AddAssignStmt(node *ast.AssignStmt) {
	// lhs and rhs should be visited in walk (expr contains node)
	info.addToken(node.Tok)
}

func (info *HalsteadInfo) AddBadDecl(node *ast.BadDecl) {
	// Nothing to count
	return
}

func (info *HalsteadInfo) AddBadExpr(node *ast.BadExpr) {
	// Nothing to count
	return
}

func (info *HalsteadInfo) AddBadStmt(node *ast.BadStmt) {
	// Nothing to count
	return
}

func (info *HalsteadInfo) AddBasicLit(node *ast.BasicLit) {
	info.operands[node.Value] += 1
}

func (info *HalsteadInfo) AddBinaryExpr(node *ast.BinaryExpr) {
	// x and y should be visited in walk (expr contains node)
	info.addToken(node.Op)
}

func (info *HalsteadInfo) AddBlockStmt(node *ast.BlockStmt) {
	// statements should be visited in walk
	info.addToken(token.LBRACE)
}

func (info *HalsteadInfo) AddBranchStmt(node *ast.BranchStmt) {
	info.addToken(node.Tok)
}

func (info *HalsteadInfo) AddCallExpr(node *ast.CallExpr) {
	// leave expressions for further walk
	// add only one of the parentheses as they are in pairs
	info.addToken(token.LPAREN)
	// Ellipsis are handled separately
}

func (info *HalsteadInfo) AddCaseClause(node *ast.CaseClause) {
	// no way to distinguish case and default, so just always assume case for now
	info.addToken(token.CASE)
	info.addToken(token.COLON)
}

func (info *HalsteadInfo) AddChanType(node *ast.ChanType) {
	info.addToken(token.ARROW)
	if node.Arrow != token.NoPos {
		info.addToken(token.CHAN)
	}
}

func (info *HalsteadInfo) AddCommClause(node *ast.CommClause) {
	// no way to distinguish case and default, so just always assume case for now
	info.addToken(token.CASE)
	info.addToken(token.COLON)
}

func (info *HalsteadInfo) AddComment(node *ast.Comment) {
	// Nothing to count
	return
}

func (info *HalsteadInfo) AddCommentGroup(node *ast.CommentGroup) {
	// Nothing to count
	return
}

func (info *HalsteadInfo) AddCompositeLit(node *ast.CompositeLit) {
	// Maybe also consider `Elts` (e.g. if it always results in commas added)
	info.addToken(token.LPAREN)
}

func (info *HalsteadInfo) AddDeclStmt(node *ast.DeclStmt) {
	// GenDecl is handled separately
	return
}

func (info *HalsteadInfo) AddDeferStmt(node *ast.DeferStmt) {
	info.addToken(token.DEFER)
}

func (info *HalsteadInfo) AddEllipsis(node *ast.Ellipsis) {
	info.addToken(token.ELLIPSIS)
}

func (info *HalsteadInfo) AddEmptyStmt(node *ast.EmptyStmt) {
	if !node.Implicit{
		info.addToken(token.SEMICOLON)
	}
}

func (info *HalsteadInfo) AddExprStmt(node *ast.ExprStmt) {
	// Nothing to count
	return
}

func (info *HalsteadInfo) AddField(node *ast.Field) {
	// Nothing to count
	return
}

func (info *HalsteadInfo) AddFieldList(node *ast.FieldList) {
	if node.Opening.IsValid() {
		info.addToken(token.LPAREN)
	}
}

func (info *HalsteadInfo) AddFile(node *ast.File) {
	info.addToken(token.PACKAGE)
}

func (info *HalsteadInfo) AddForStmt(node *ast.ForStmt) {
	info.addToken(token.FOR)
}

func (info *HalsteadInfo) AddFuncDecl(node *ast.FuncDecl) {
	// Composite type only, nothing
	return
}

func (info *HalsteadInfo) AddFuncLit(node *ast.FuncLit) {
	// Composite type only, nothing
	return
}

func (info *HalsteadInfo) AddFuncType(node *ast.FuncType) {
	info.addToken(token.FUNC)
}

func (info *HalsteadInfo) AddGenDecl(node *ast.GenDecl) {
	info.addToken(node.Tok)
	if node.Lparen.IsValid() {
		info.addToken(token.LPAREN)
	}
}

func (info *HalsteadInfo) AddGoStmt(node *ast.GoStmt) {
	info.addToken(token.GO)
}

func (info *HalsteadInfo) AddIdent(node *ast.Ident) {
	info.operands[node.Name] += 1
}

func (info *HalsteadInfo) AddIfStmt(node *ast.IfStmt) {
	info.addToken(token.IF)
}

func (info *HalsteadInfo) AddImportSpec(node *ast.ImportSpec) {
	// Composite type only, nothing
	return
}

func (info *HalsteadInfo) AddIncDecStmt(node *ast.IncDecStmt) {
	info.addToken(node.Tok)
}

func (info *HalsteadInfo) AddIndexExpr(node *ast.IndexExpr) {
	if node.Lbrack.IsValid() {
		info.addToken(token.LBRACK)
	}
}

func (info *HalsteadInfo) AddInterfaceType(node *ast.InterfaceType) {
	info.addToken(token.INTERFACE)
}

func (info *HalsteadInfo) AddKeyValueExpr(node *ast.KeyValueExpr) {
	info.addToken(token.COLON)
}

func (info *HalsteadInfo) AddLabeledStmt(node *ast.LabeledStmt) {
	info.addToken(token.COLON)
}

func (info *HalsteadInfo) AddMapType(node *ast.MapType) {
	info.addToken(token.MAP)
}

func (info *HalsteadInfo) AddPackage(node *ast.Package) {
	// name should be handled by ident, so ignore is as apparently it is 
	// not a particular part of code but rather abstract entity (set of files?).
	return
}

func (info *HalsteadInfo) AddParenExpr(node *ast.ParenExpr) {
	info.addToken(token.LPAREN)
}

func (info *HalsteadInfo) AddRangeStmt(node *ast.RangeStmt) {
	// for should be already handled by `AddForStmt`
	info.addToken(node.Tok)
	info.addToken(token.RANGE)
}

func (info *HalsteadInfo) AddReturnStmt(node *ast.ReturnStmt) {
	info.addToken(token.RETURN)
}

func (info *HalsteadInfo) AddSelectStmt(node *ast.SelectStmt) {
	info.addToken(token.SELECT)
}

func (info *HalsteadInfo) AddSelectorExpr(node *ast.SelectorExpr) {
	info.addToken(token.PERIOD)
	return
}

func (info *HalsteadInfo) AddSendStmt(node *ast.SendStmt) {
	info.addToken(token.ARROW)
}

func (info *HalsteadInfo) AddSliceExpr(node *ast.SliceExpr) {
	info.addToken(token.LBRACK)
}

func (info *HalsteadInfo) AddStarExpr(node *ast.StarExpr) {
	info.addToken(token.MUL)
}

func (info *HalsteadInfo) AddStructType(node *ast.StructType) {
	info.addToken(token.STRUCT)
}

func (info *HalsteadInfo) AddSwitchStmt(node *ast.SwitchStmt) {
	info.addToken(token.SWITCH)
}

func (info *HalsteadInfo) AddTypeAssertExpr(node *ast.TypeAssertExpr) {
	info.addToken(token.LPAREN)
}

func (info *HalsteadInfo) AddTypeSpec(node *ast.TypeSpec) {
	if node.Assign.IsValid() {
		info.addToken(token.ASSIGN)
	}
}

func (info *HalsteadInfo) AddTypeSwitchStmt(node *ast.TypeSwitchStmt) {
	info.addToken(token.SWITCH)
}

func (info *HalsteadInfo) AddUnaryExpr(node *ast.UnaryExpr) {
	info.addToken(node.Op)
}

func (info *HalsteadInfo) AddValueSpec(node *ast.ValueSpec) {
	// Something composite only, ignore
	return
}
