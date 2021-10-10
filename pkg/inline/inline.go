// Metric for inline data complexity measurement
package inline

import (
	"go/ast"
	"go/token"
)

// Weight of different data complexity
type Weights struct {
	// Integer constants complexity
	Int float64 `xml:"int"`
	// Float constants complexity
	Float float64 `xml:"float"`
	// Imaginary numbers constants complexity
	Imag float64 `xml:"imag"`
	// Characters constants complexity
	Char float64 `xml:"char"`
	// Strings constants complexity
	String float64 `xml:"string"`
	// Composite literals (like structure initialization) constants complexity
	CompositeLit float64 `xml:"composite"`
}

// Intermidiate state of metric
type Metric struct {
	// Config with all metric's weights
	Config Weights
	comp   float64
}

// Parses ast node and collects all info about inlined constants in code
func (m *Metric) ParseNode(n ast.Node) {
	switch v := n.(type) {
	case *ast.BasicLit:
		m.comp += getBasicLitComp(v, &m.Config)
	case *ast.CompositeLit:
		m.comp += m.Config.CompositeLit
	}
}

// Returns final metric's result
func (m Metric) Finish() float64 { return m.comp }

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
