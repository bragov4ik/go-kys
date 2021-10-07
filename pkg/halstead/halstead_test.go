package halstead

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestSumMap(t *testing.T) {
	type args struct {
		targetMap map[string]uint
	}
	tests := []struct {
		name string
		args args
		want uint
	}{
		{"EmptyMap", args{map[string]uint{}}, 0},
		{"Map1", args{map[string]uint{"a": 1}}, 1},
		{"Map2", args{map[string]uint{"a": 1, "b": 2}}, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sumMap(tt.args.targetMap); got != tt.want {
				t.Errorf("sumMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTokenInArr(t *testing.T) {
	type args struct {
		tok token.Token
		arr []token.Token
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"Test empty", args{tok: token.ADD, arr: []token.Token{}}, false},
		{"Test 1", args{tok: token.ADD, arr: []token.Token{token.ADD}}, true},
		{"Test 2", args{tok: token.ADD, arr: []token.Token{token.SUB}}, false},
		{"Test 3", args{tok: token.ADD, arr: []token.Token{token.SUB, token.XOR, token.ADD}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tokenInArr(tt.args.tok, tt.args.arr); got != tt.want {
				t.Errorf("tokenInArr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMetric(t *testing.T) {
	metric := NewMetric()
	// Code taken from A Tour of Go
	source := `package main

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
	}
	`
	type result struct {
		n1 uint
		n2 uint
		N1 uint
		N2 uint
	}

	notEqual := func(r1 result, r2 result) bool {
		return r1.n1 != r2.n1 || r1.n2 != r2.n2 || r1.N1 != r2.N1 || r1.N2 != r2.N2
	}

	want := result{12, 16, 34, 29}
	file, err := parser.ParseFile(token.NewFileSet(), "", source, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	ast.Inspect(file, func(n ast.Node) bool {
		metric.ParseNode(n)
		return true
	})

	got := result{
		metric.getN1Distinct(),
		metric.getN2Distinct(),
		metric.getN1Total(),
		metric.getN2Total(),
	}
	if notEqual(want, got) {
		t.Errorf("result{metric.getN1Distinct(), metric.getN2Distinct(), metric.getN1Total(), metric.getN2Total(),} = %v, want %v", got, want)
	}
}
