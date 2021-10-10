// Package with [wmfp](https://en.wikipedia.org/wiki/Weighted_Micro_Function_Points) metric
// calculator.
package wmfp

import (
	"go/ast"

	"github.com/bragov4ik/go-kys/pkg/arithmetic"
	codestruct "github.com/bragov4ik/go-kys/pkg/codestruct"
	comments "github.com/bragov4ik/go-kys/pkg/comments"
	cyclo "github.com/bragov4ik/go-kys/pkg/cyclocomp"
	halstead "github.com/bragov4ik/go-kys/pkg/halstead"
	inline "github.com/bragov4ik/go-kys/pkg/inline"
)

// State for WMFP metrics
type MeasurerWMFP struct {
	// State of comments metrics
	Comments *comments.Metric
	// State of cyclo complexity metrics
	Cyclo *cyclo.Metric
	// State of halstead metrics
	Halst *halstead.Metric
	// State of complexity of code structure
	Codestruct *codestruct.Metric
	// State of complexity of inline data constants
	InlineData *inline.Metric
	// State of complexity of arithmetic expressions
	ArithmeticComp *arithmetic.Metric

	halstWeight float64
}

// Interface for underlaying metrics
type Metric interface {
	// Parses ast node and collects all info for metric
	ParseNode(ast.Node)
	// Returns final score for metric
	Finish() float64
}

// Config with all weights for underlaying metrics
type Config struct {
	// Cyclo complexity weights
	CycloComp cyclo.Weights `xml:"cyclomatic"`
	// Comments complexity weights
	Comment comments.Weights `xml:"comment"`
	// Code structure complexity weights
	CodeStructComp codestruct.Weights `xml:"codestruct"`
	// Inline data complexity weights
	InlineData inline.Weights `xml:"inline"`
	// Arithmetic expression complexity weights
	ArithmeticComp arithmetic.Weights `xml:"arithmetic"`
	// Halstead metric weight
	Halstead float64 `xml:"halstead"`
}

// Constructor for WMFP metric
func NewMeasurerWMFP(config *Config) MeasurerWMFP {
	halst := halstead.NewMetric()
	return MeasurerWMFP{
		Comments: &comments.Metric{
			Config: config.Comment,
		},
		Cyclo: &cyclo.Metric{
			Config: config.CycloComp,
		},
		Halst: &halst,
		Codestruct: &codestruct.Metric{
			Config: config.CodeStructComp,
		},
		InlineData: &inline.Metric{
			Config: config.InlineData,
		},
		ArithmeticComp: &arithmetic.Metric{
			Config: config.ArithmeticComp,
		},
		halstWeight: config.Halstead,
	}
}

// Parses single file using WMFP metric
func (m *MeasurerWMFP) ParseFile(file *ast.File) {
	ast.Inspect(file, func(n ast.Node) bool {
		m.parseNode(n)
		return true
	})
}

// Returns final score of metric
func (m *MeasurerWMFP) Finish() (total float64) {
	total += m.Comments.Finish()
	total += m.Cyclo.Finish()
	total += m.Halst.Finish() * m.halstWeight
	total += m.Codestruct.Finish()
	total += m.InlineData.Finish()
	total += m.ArithmeticComp.Finish()
	return
}

func (m *MeasurerWMFP) metrics() []Metric {
	return []Metric{
		m.Comments,
		m.Cyclo,
		m.Halst,
		m.Codestruct,
		m.InlineData,
		m.ArithmeticComp,
	}
}

func (measurer *MeasurerWMFP) parseNode(n ast.Node) {
	for _, m := range measurer.metrics() {
		m.ParseNode(n)
	}
}
