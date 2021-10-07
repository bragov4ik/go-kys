package wmfp

import (
	"go/ast"
	"reflect"

	comments "github.com/bragov4ik/go-kys/pkg/comments"
	cyclo "github.com/bragov4ik/go-kys/pkg/cyclocomp"
	halstead "github.com/bragov4ik/go-kys/pkg/halstead"
)

type MeasurerWMFP struct {
	Comments *comments.Metric
	Cyclo    *cyclo.Metric
	Halst    *halstead.Metric
}

type Metric interface {
	ParseNode(ast.Node)
	Finish() float64
}

type Config struct {
	CycloComp cyclo.Weights    `xml:"cyclomatic"`
	Comment   comments.Weights `xml:"comment"`
}

func NewMeasurerWMFP() MeasurerWMFP {
	halst := halstead.NewMetric()
	return MeasurerWMFP{
		Comments: &comments.Metric{},
		Cyclo:    &cyclo.Metric{},
		Halst:    &halst,
	}
}

func (m *MeasurerWMFP) Configure(config *Config) {
	m.Comments.Config = config.Comment
	m.Cyclo.Config = config.CycloComp
}

func (m *MeasurerWMFP) ParseFile(file *ast.File) {
	ast.Inspect(file, func(n ast.Node) bool {
		m.parseNode(n)
		return true
	})
}

func (m *MeasurerWMFP) Finish() (total float64) {
	v := reflect.ValueOf(*m)
	for i := 0; i < v.NumField(); i++ {
		m, ok := v.Field(i).Interface().(Metric)
		if ok {
			total += m.Finish()
		}
	}
	return
}

func (m *MeasurerWMFP) parseNode(n ast.Node) {
	v := reflect.ValueOf(*m)
	for i := 0; i < v.NumField(); i++ {
		m, ok := v.Field(i).Interface().(Metric)
		if ok {
			m.ParseNode(n)
		}
	}
}
