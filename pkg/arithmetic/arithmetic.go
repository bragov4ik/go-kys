package arithmetic

import (
	"go/ast"
	"go/token"
)

type Weights struct {
	Add       float64 `xml:"add"`
	Sub       float64 `xml:"sub"`
	Mul       float64 `xml:"mul"`
	Quo       float64 `xml:"quo"`
	Rem       float64 `xml:"rem"`
	AddAssign float64 `xml:"add_assign"`
	SubAssign float64 `xml:"sub_assign"`
	MulAssign float64 `xml:"mul_assign"`
	QuoAssign float64 `xml:"quo_assign"`
	RemAssign float64 `xml:"rem_assign"`
	Inc       float64 `xml:"inc"`
	Dec       float64 `xml:"dec"`
}

type Metric struct {
	Config Weights
	Comp   float64
}

func (m *Metric) ParseNode(n ast.Node) {
	m.Comp += getArithmeticComp(&n, &m.Config)
}

func (m Metric) Finish() float64 {
	return m.Comp
}

func getArithmeticComp(n *ast.Node, config *Weights) float64 {
	var comp float64
	switch v := (*n).(type) {
	case *ast.BinaryExpr:
		comp = getBinaryComp(v, config)
	case *ast.UnaryExpr:
		comp = getUnaryComp(v, config)
	case *ast.IncDecStmt:
		comp = getIncDecComp(v, config)
	case *ast.AssignStmt:
		comp = getAssignComp(v, config)
	}
	return comp
}

func getBinaryComp(n *ast.BinaryExpr, config *Weights) float64 {
	var comp float64
	switch n.Op {
	case token.ADD:
		comp = config.Add
	case token.SUB:
		comp = config.Sub
	case token.MUL:
		comp = config.Mul
	case token.QUO:
		comp = config.Quo
	case token.REM:
		comp = config.Rem
	}
	return comp
}

func getUnaryComp(n *ast.UnaryExpr, config *Weights) float64 {
	var comp float64
	switch n.Op {
	case token.ADD:
		comp = config.Add
	case token.SUB:
		comp = config.Sub
	}
	return comp
}

func getIncDecComp(n *ast.IncDecStmt, config *Weights) float64 {
	var comp float64
	switch n.Tok {
	case token.INC:
		comp = config.Inc
	case token.DEC:
		comp = config.Dec
	}
	return comp
}

func getAssignComp(n *ast.AssignStmt, config *Weights) float64 {
	var comp float64
	switch n.Tok {
	case token.ADD_ASSIGN:
		comp = config.AddAssign
	case token.SUB_ASSIGN:
		comp = config.SubAssign
	case token.MUL_ASSIGN:
		comp = config.MulAssign
	case token.QUO_ASSIGN:
		comp = config.QuoAssign
	case token.REM_ASSIGN:
		comp = config.RemAssign
	}
	return comp
}
