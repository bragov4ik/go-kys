package wmfp

import (
	"go/ast"
	"reflect"

	"github.com/bragov4ik/go-kys/pkg/arithmetic"
	codestruct "github.com/bragov4ik/go-kys/pkg/codestruct"
	comments "github.com/bragov4ik/go-kys/pkg/comments"
	cyclo "github.com/bragov4ik/go-kys/pkg/cyclocomp"
	halstead "github.com/bragov4ik/go-kys/pkg/halstead"
	inline "github.com/bragov4ik/go-kys/pkg/inline"
)

type MeasurerWMFP struct {
	Comments       *comments.Metric
	Cyclo          *cyclo.Metric
	Halst          *halstead.Metric
	HalstWeight    float64
	Codestruct     *codestruct.Metric
	InlineData     *inline.Metric
	ArithmeticComp *arithmetic.Metric
}

type Metric interface {
	ParseNode(ast.Node)
	Finish() float64
}

type Config struct {
	CycloComp      cyclo.Weights      `xml:"cyclomatic"`
	Comment        comments.Weights   `xml:"comment"`
	CodeStructComp codestruct.Weights `xml:"codestruct"`
	InlineData     inline.Weights     `xml:"inline"`
	ArithmeticComp arithmetic.Weights `xml:"arithmetic"`
	Halstead       float64            `xml:"halstead"`
}

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
	}
}

func (m *MeasurerWMFP) ParseFile(file *ast.File) {
	ast.Inspect(file, func(n ast.Node) bool {
		m.parseNode(n)
		return true
	})
}

func (m *MeasurerWMFP) Finish() (total float64) {
	total += m.Comments.Finish()
	total += m.Cyclo.Finish()
	total += m.Halst.Finish() * m.HalstWeight
	total += m.Codestruct.Finish()
	total += m.InlineData.Finish()
	total += m.ArithmeticComp.Finish()
	return
}

func (measurer *MeasurerWMFP) parseNode(n ast.Node) {
	v := reflect.ValueOf(*measurer)
	for i := 0; i < v.NumField(); i++ {
		if m, ok := v.Field(i).Interface().(Metric); ok {
			m.ParseNode(n)
		}
	}
}
