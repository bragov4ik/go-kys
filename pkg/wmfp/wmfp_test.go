package wmfp

import (
	"go/parser"
	"go/token"
	"testing"

	"github.com/bragov4ik/go-kys/pkg/arithmetic"
	"github.com/bragov4ik/go-kys/pkg/codestruct"
	"github.com/bragov4ik/go-kys/pkg/comments"
	cyclo "github.com/bragov4ik/go-kys/pkg/cyclocomp"
	"github.com/bragov4ik/go-kys/pkg/inline"
)

func TestWMFP(t *testing.T) {
	tourOfGo := `
	package main

	import "fmt"
	
	const (
		// Create a huge number by shifting a 1 bit left 100 places.
		// In other words, the binary number that is 1 followed by 100 zeroes.
		Big = 1 << 100
		// Shift it right again 99 places, so we end up with 1<<1, or 2.
		Small = Big >> 99
	)
	
	func needInt(x int) int { return x*10 + 1 }
	func needFloat(x float64) float64 {
		return x * 0.1
	}

	func main() {
		fmt.Println(needInt(Small))
		// fmt.Println(needInt(Big))
		fmt.Println(needFloat(Small))
		fmt.Println(needFloat(Big))
	}`

	cfg := Config{
		CycloComp: cyclo.Weights{
			If:   1,
			For:  1,
			Rng:  1,
			Case: 1,
			And:  1,
			Or:   1,
		},
		Comment: comments.Weights{Word: 1},
		CodeStructComp: codestruct.Weights{
			Func:      1,
			Struct:    1,
			Interface: 1,
		},
		InlineData: inline.Weights{
			Int:          1,
			Float:        1,
			Imag:         1,
			Char:         1,
			String:       1,
			CompositeLit: 1,
		},
		ArithmeticComp: arithmetic.Weights{
			Add:       1,
			Sub:       1,
			Mul:       1,
			Quo:       1,
			Rem:       1,
			AddAssign: 1,
			SubAssign: 1,
			MulAssign: 1,
			QuoAssign: 1,
			RemAssign: 1,
			Inc:       1,
			Dec:       1,
		},
		Halstead: 1,
	}

	m := NewMeasurerWMFP(&cfg)
	file, err := parser.ParseFile(token.NewFileSet(), "", tourOfGo, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	m.ParseFile(file)
	var expect uint = 62

	if got := m.Finish(); uint(got) != expect {
		t.Errorf("Finish = %v, want %v", got, expect)
	}
}
