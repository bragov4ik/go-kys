package halstead

import (
	"go/ast"
	"go/token"
	"math"
)

type Metric struct {
	operators map[token.Token]uint
	operands  map[string]uint
}

func NewMetric() Metric {
	return Metric{
		operators: make(map[token.Token]uint),
		operands:  make(map[string]uint),
	}
}

func (m *Metric) ParseNode(n ast.Node) {
	switch v := n.(type) {
	case *ast.AssignStmt:
		m.addAssignStmt(v)
	case *ast.BasicLit:
		m.addBasicLit(v)
	case *ast.BinaryExpr:
		m.addBinaryExpr(v)
	case *ast.BlockStmt:
		m.addBlockStmt(v)
	case *ast.BranchStmt:
		m.addBranchStmt(v)
	case *ast.CallExpr:
		m.addCallExpr(v)
	case *ast.CaseClause:
		m.addCaseClause(v)
	case *ast.ChanType:
		m.addChanType(v)
	case *ast.CommClause:
		m.addCommClause(v)
	case *ast.CompositeLit:
		m.addCompositeLit(v)
	case *ast.DeferStmt:
		m.addDeferStmt(v)
	case *ast.Ellipsis:
		m.addEllipsis(v)
	case *ast.EmptyStmt:
		m.addEmptyStmt(v)
	case *ast.FieldList:
		m.addFieldList(v)
	case *ast.File:
		m.addFile(v)
	case *ast.ForStmt:
		m.addForStmt(v)
	case *ast.FuncType:
		m.addFuncType(v)
	case *ast.GenDecl:
		m.addGenDecl(v)
	case *ast.GoStmt:
		m.addGoStmt(v)
	case *ast.Ident:
		m.addIdent(v)
	case *ast.IfStmt:
		m.addIfStmt(v)
	case *ast.IncDecStmt:
		m.addIncDecStmt(v)
	case *ast.IndexExpr:
		m.addIndexExpr(v)
	case *ast.InterfaceType:
		m.addInterfaceType(v)
	case *ast.KeyValueExpr:
		m.addKeyValueExpr(v)
	case *ast.LabeledStmt:
		m.addLabeledStmt(v)
	case *ast.MapType:
		m.addMapType(v)
	case *ast.ParenExpr:
		m.addParenExpr(v)
	case *ast.RangeStmt:
		m.addRangeStmt(v)
	case *ast.ReturnStmt:
		m.addReturnStmt(v)
	case *ast.SelectStmt:
		m.addSelectStmt(v)
	case *ast.SelectorExpr:
		m.addSelectorExpr(v)
	case *ast.SendStmt:
		m.addSendStmt(v)
	case *ast.SliceExpr:
		m.addSliceExpr(v)
	case *ast.StarExpr:
		m.addStarExpr(v)
	case *ast.StructType:
		m.addStructType(v)
	case *ast.SwitchStmt:
		m.addSwitchStmt(v)
	case *ast.TypeAssertExpr:
		m.addTypeAssertExpr(v)
	case *ast.TypeSpec:
		m.addTypeSpec(v)
	case *ast.TypeSwitchStmt:
		m.addTypeSwitchStmt(v)
	case *ast.UnaryExpr:
		m.addUnaryExpr(v)
	}
}

func (m Metric) Finish() float64 {
	n1 := float64(m.n1Distinct())
	n2 := float64(m.n2Distinct())
	N1 := float64(m.n1Total())
	N2 := float64(m.n2Total())
	return (N1 + N2) * math.Log2(n1+n2)
}

func (m *Metric) n1Distinct() uint { return uint(len(m.operators)) }
func (m *Metric) n2Distinct() uint { return uint(len(m.operands)) }

// Not universal key type because using interfaces is nasty
// and generics (with 1.18 version) are not released yet
func sumMap(targetMap map[string]uint) uint {
	var total uint = 0
	for _, count := range targetMap {
		total += count
	}
	return total
}

// Not universal key type because using interfaces is nasty
// and generics (with 1.18 version) are not released yet
func sumMapTok(targetMap map[token.Token]uint) uint {
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

func (m *Metric) n1Total() uint { return sumMap(m.operands) }
func (m *Metric) n2Total() uint { return sumMapTok(m.operators) }

func (m *Metric) addToken(nextToken token.Token) {
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
		m.operators[nextToken] += 1
	}
}

// TODO count commas

func (m *Metric) addAssignStmt(node *ast.AssignStmt) {
	// lhs and rhs should be visited in walk (expr contains node)
	m.addToken(node.Tok)
}

func (m *Metric) addBasicLit(node *ast.BasicLit) { m.operands[node.Value] += 1 }

func (m *Metric) addBinaryExpr(node *ast.BinaryExpr) {
	// x and y should be visited in walk (expr contains node)
	m.addToken(node.Op)
}

func (m *Metric) addBlockStmt(node *ast.BlockStmt) {
	// statements should be visited in walk
	m.addToken(token.LBRACE)
}

func (m *Metric) addBranchStmt(node *ast.BranchStmt) { m.addToken(node.Tok) }

func (m *Metric) addCallExpr(node *ast.CallExpr) {
	// leave expressions for further walk
	// add only one of the parentheses as they are in pairs
	m.addToken(token.LPAREN)
	// Ellipsis are handled separately
}

func (m *Metric) addCaseClause(node *ast.CaseClause) {
	// no way to distinguish case and default, so just always assume case for now
	m.addToken(token.CASE)
	m.addToken(token.COLON)
}

func (m *Metric) addChanType(node *ast.ChanType) {
	m.addToken(token.ARROW)
	if node.Arrow != token.NoPos {
		m.addToken(token.CHAN)
	}
}

func (m *Metric) addCommClause(node *ast.CommClause) {
	// no way to distinguish case and default, so just always assume case for now
	m.addToken(token.CASE)
	m.addToken(token.COLON)
}

func (m *Metric) addCompositeLit(node *ast.CompositeLit) {
	// Maybe also consider `Elts` (e.g. if it always results in commas added)
	m.addToken(token.LPAREN)
}

func (m *Metric) addDeferStmt(node *ast.DeferStmt) { m.addToken(token.DEFER) }
func (m *Metric) addEllipsis(node *ast.Ellipsis)   { m.addToken(token.ELLIPSIS) }

func (m *Metric) addEmptyStmt(node *ast.EmptyStmt) {
	if !node.Implicit {
		m.addToken(token.SEMICOLON)
	}
}

func (m *Metric) addFieldList(node *ast.FieldList) {
	if node.Opening.IsValid() {
		m.addToken(token.LPAREN)
	}
}

func (m *Metric) addFile(node *ast.File)         { m.addToken(token.PACKAGE) }
func (m *Metric) addForStmt(node *ast.ForStmt)   { m.addToken(token.FOR) }
func (m *Metric) addFuncType(node *ast.FuncType) { m.addToken(token.FUNC) }

func (m *Metric) addGenDecl(node *ast.GenDecl) {
	m.addToken(node.Tok)
	if node.Lparen.IsValid() {
		m.addToken(token.LPAREN)
	}
}

func (m *Metric) addGoStmt(node *ast.GoStmt)         { m.addToken(token.GO) }
func (m *Metric) addIdent(node *ast.Ident)           { m.operands[node.Name] += 1 }
func (m *Metric) addIfStmt(node *ast.IfStmt)         { m.addToken(token.IF) }
func (m *Metric) addIncDecStmt(node *ast.IncDecStmt) { m.addToken(node.Tok) }

func (m *Metric) addIndexExpr(node *ast.IndexExpr) {
	if node.Lbrack.IsValid() {
		m.addToken(token.LBRACK)
	}
}

func (m *Metric) addInterfaceType(node *ast.InterfaceType) { m.addToken(token.INTERFACE) }
func (m *Metric) addKeyValueExpr(node *ast.KeyValueExpr)   { m.addToken(token.COLON) }
func (m *Metric) addLabeledStmt(node *ast.LabeledStmt)     { m.addToken(token.COLON) }
func (m *Metric) addMapType(node *ast.MapType)             { m.addToken(token.MAP) }
func (m *Metric) addParenExpr(node *ast.ParenExpr)         { m.addToken(token.LPAREN) }

func (m *Metric) addRangeStmt(node *ast.RangeStmt) {
	// for should be already handled by `addForStmt`
	m.addToken(node.Tok)
	m.addToken(token.RANGE)
}

func (m *Metric) addReturnStmt(node *ast.ReturnStmt)         { m.addToken(token.RETURN) }
func (m *Metric) addSelectStmt(node *ast.SelectStmt)         { m.addToken(token.SELECT) }
func (m *Metric) addSelectorExpr(node *ast.SelectorExpr)     { m.addToken(token.PERIOD) }
func (m *Metric) addSendStmt(node *ast.SendStmt)             { m.addToken(token.ARROW) }
func (m *Metric) addSliceExpr(node *ast.SliceExpr)           { m.addToken(token.LBRACK) }
func (m *Metric) addStarExpr(node *ast.StarExpr)             { m.addToken(token.MUL) }
func (m *Metric) addStructType(node *ast.StructType)         { m.addToken(token.STRUCT) }
func (m *Metric) addSwitchStmt(node *ast.SwitchStmt)         { m.addToken(token.SWITCH) }
func (m *Metric) addTypeAssertExpr(node *ast.TypeAssertExpr) { m.addToken(token.LPAREN) }

func (m *Metric) addTypeSpec(node *ast.TypeSpec) {
	if node.Assign.IsValid() {
		m.addToken(token.ASSIGN)
	}
}

func (m *Metric) addTypeSwitchStmt(node *ast.TypeSwitchStmt) { m.addToken(token.SWITCH) }
func (m *Metric) addUnaryExpr(node *ast.UnaryExpr)           { m.addToken(node.Op) }
