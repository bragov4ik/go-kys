// Package with metric which checks general code structure.
package codestruct

import "go/ast"

// Weights for metric
type Weights struct {
	// Function declaration weight
	Func float64 `xml:"func"`
	// Structure declaration weight
	Struct float64 `xml:"struct"`
	// Interface declaration weight
	Interface float64 `xml:"interface"`
}

// Intermidiate state for code structure metric
type Metric struct {
	// Config with weights
	Config Weights
	comp   float64
}

// Parses ast node and collects result of metric
func (m *Metric) ParseNode(n ast.Node) { m.comp += getCodeStructComp(n, &m.Config) }

// Returns final result of metric
func (m *Metric) Finish() float64 { return m.comp }

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
