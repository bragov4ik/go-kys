package arithmetic

import (
	"go/ast"
	"go/token"
)

type Weights struct {
	Add       uint `xml:"add"`
	Sub       uint `xml:"sub"`
	Mul       uint `xml:"mul"`
	Quo       uint `xml:"quo"`
	Rem       uint `xml:"rem"`
	AddAssign uint `xml:"add_assign"`
	SubAssign uint `xml:"sub_assign"`
	MulAssign uint `xml:"mul_assign"`
	QuoAssign uint `xml:"quo_assign"`
	RemAssign uint `xml:"rem_assign"`
	Inc       uint `xml:"inc"`
	Dec       uint `xml:"dec"`
}

type Metric struct {
	Config Weights
	Comp   uint
}

func (m *Metric) ParseNode(n ast.Node) {
	m.Comp += getArithmeticComp(&n, &m.Config)
}

func (m Metric) Finish() float64 {
	return float64(m.Comp)
}

func getArithmeticComp(n *ast.Node, config *Weights) uint {
	var comp uint
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

func getBinaryComp(n *ast.BinaryExpr, config *Weights) uint {
	var comp uint
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

func getUnaryComp(n *ast.UnaryExpr, config *Weights) uint {
	var comp uint
	switch n.Op {
	case token.ADD:
		comp = config.Add
	case token.SUB:
		comp = config.Sub
	}
	return comp
}

func getIncDecComp(n *ast.IncDecStmt, config *Weights) uint {
	var comp uint
	switch n.Tok {
	case token.INC:
		comp = config.Inc
	case token.DEC:
		comp = config.Dec
	}
	return comp
}

func getAssignComp(n *ast.AssignStmt, config *Weights) uint {
	var comp uint
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
