package codestruct

import (
	"go/ast"
)

type Weights struct {
	Func      float64 `xml:"func"`
	Struct    float64 `xml:"struct"`
	Interface float64 `xml:"interface"`
}

type Metric struct {
	Config Weights
	Comp   float64
}

func (m *Metric) ParseNode(n ast.Node) {
	m.Comp += getCodeStructComp(n, &m.Config)
}

func getCodeStructComp(n ast.Node, cfg *Weights) float64 {
	var comp float64

	switch n.(type) {
	case *ast.StructType:
		comp = cfg.Struct
	case *ast.FuncDecl:
		comp = cfg.Func
	case *ast.InterfaceType:
		comp = cfg.Interface
	}

	return comp
}

func (m Metric) Finish() float64 {
	return float64(m.Comp)
}
