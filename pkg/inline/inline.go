package inline

import (
	"go/ast"
	"go/token"
)

type Weights struct {
	Int          float64 `xml:"int"`
	Float        float64 `xml:"float"`
	Imag         float64 `xml:"imag"`
	Char         float64 `xml:"char"`
	String       float64 `xml:"string"`
	CompositeLit float64 `xml:"composite"`
}

type Metric struct {
	Config Weights
	Comp   float64
}

func (m *Metric) ParseNode(n ast.Node) {
	switch v := n.(type) {
	case *ast.BasicLit:
		m.Comp += getBasicLitComp(v, &m.Config)
	case *ast.CompositeLit:
		m.Comp += m.Config.CompositeLit
	}
}

func (m Metric) Finish() float64 {
	return m.Comp
}

func getBasicLitComp(literal *ast.BasicLit, config *Weights) float64 {
	var comp float64
	len := float64(len(literal.Value))
	switch literal.Kind {
	case token.INT:
		comp = config.Int
	case token.FLOAT:
		comp = config.Float
	case token.IMAG:
		comp = config.Imag
	case token.CHAR:
		comp = config.Char
	case token.STRING:
		comp = config.String * len
	}
	return comp
}
