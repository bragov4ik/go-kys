// Package arithmetic calculates time spend on writing arithmetic
// expressions.
package arithmetic

import (
	"go/ast"
	"go/token"
)

// Weights is structure with weights for arithmetic metric calculator
type Weights struct {
	// Weight for addition
	Add float64 `xml:"add"`
	// Weight for subtraction
	Sub float64 `xml:"sub"`
	// Weight for multiplication
	Mul float64 `xml:"mul"`
	// Weight for division
	Quo float64 `xml:"quo"`
	// Weight for remainder of division
	Rem float64 `xml:"rem"`
	// Weight for addition assigned
	AddAssign float64 `xml:"add_assign"`
	// Weight for subtraction assigned
	SubAssign float64 `xml:"sub_assign"`
	// Weight for multiplication assigned
	MulAssign float64 `xml:"mul_assign"`
	// Weight for division assigned
	QuoAssign float64 `xml:"quo_assign"`
	// Weight for remainder assigned
	RemAssign float64 `xml:"rem_assign"`
	// Weight for increment by 1
	Inc float64 `xml:"inc"`
	// Weight for decrement by 1
	Dec float64 `xml:"dec"`
}

// Metric is the temporal state for calculations of metrics
type Metric struct {
	// Config with weights
	Config Weights
	comp   float64
}

// Parses node from ast
func (m *Metric) ParseNode(n ast.Node) { m.comp += getArithmeticComp(&n, &m.Config) }

// Finishes calculation and returns result
func (m Metric) Finish() float64 { return m.comp }

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
