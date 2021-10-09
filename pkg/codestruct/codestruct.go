package codestruct

import (
	"go/ast"
)

type Weights struct {
	Func      uint `xml:"func"`
	Struct    uint `xml:"struct"`
	Interface uint `xml:"interface"`
}

type Metric struct {
	Config Weights
	Comp   uint
}

func (m *Metric) ParseNode(n ast.Node) {
	m.Comp += getCodeStructComp(n, &m.Config)
}

func getCodeStructComp(n ast.Node, cfg *Weights) uint {
	comp := uint(0)

	switch n.(type) {
	case *ast.StructType:
		comp += cfg.Struct
	case *ast.FuncDecl:
		comp += cfg.Func
	case *ast.InterfaceType:
		comp += cfg.Interface
	}

	return comp
}

func (m Metric) Finish() float64 {
	return float64(m.Comp)
}
