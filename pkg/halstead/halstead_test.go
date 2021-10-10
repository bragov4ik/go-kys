package halstead

import (
	"go/ast"
	"go/parser"
	"go/token"
	"math"
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

type Result struct {
	n1     uint
	n2     uint
	N1     uint
	N2     uint
	finish uint
}

func getMetric(source string) Result {
	m := NewMetric()
	file, err := parser.ParseFile(token.NewFileSet(), "", "package main;"+source, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	ast.Inspect(file, func(n ast.Node) bool {
		m.ParseNode(n)
		return true
	})

	return Result{m.n1Distinct(), m.n2Distinct(), m.n1Total(), m.n2Total(), uint(m.Finish())}
}

func TestMetric(t *testing.T) {
	tourOfGo := `
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

	tests := []struct {
		Name   string
		Source string
		Want   Result
	}{
		{"tourOfGo", tourOfGo, Result{12, 16, 34, 29, uint(63 * math.Log2(28))}},
		{"addAssign", `func f() { x := 0; x += 2 }`, Result{6, 5, 6, 6, uint(12 * math.Log2(11))}},
		{"branchStmt", `func f() { for _, _ := range []int{} { break } }`, Result{7, 4, 5, 9, uint(14 * math.Log2(11))}},
		{"caseClause", `func f() { switch 10 { case 1: default: }}`, Result{7, 4, 4, 10, uint(14 * math.Log2(11))}},
		{"declStmt", `type X int; const I = 1; var a int`, Result{4, 6, 7, 4, uint(11 * math.Log2(10))}},
		{"ellipsis", `func f(aboba ...int) { f(aboba...) }`, Result{5, 4, 6, 6, uint(12 * math.Log2(9))}},
		{"emptyStmt", `func f() { ; }`, Result{5, 2, 2, 5, uint(7 * math.Log2(7))}},
		{"forStmt", `func f() { for i := 0; i < 0; i++ {} }`, Result{8, 4, 7, 9, uint(16 * math.Log2(12))}},
		{"channels", `func f(i chan int) { select { case a := <-i: i <- a } }`, Result{9, 5, 8, 12, uint(20 * math.Log2(14))}},
		{"defer", `func f() { defer f() }`, Result{5, 2, 3, 6, uint(9 * math.Log2(7))}},
		{"go", `func f() { go f() }`, Result{5, 2, 3, 6, uint(9 * math.Log2(7))}},
		{"if", `func f() bool { if f() { return false } else { return true } }`, Result{6, 5, 6, 10, uint(16 * math.Log2(11))}},
		{"index", `func f(a []int) int { return a[0] }`, Result{6, 5, 7, 6, uint(13 * math.Log2(11))}},
		{"interface", `type Aboba interface { Aboba() bool }`, Result{5, 3, 4, 6, uint(10 * math.Log2(8))}},
		{"keyvalue-expr", `var a = map[string]int{"a": 1}`, Result{5, 6, 6, 5, uint(11 * math.Log2(11))}},
		{"goto", `func x() { start: for { goto start }}`, Result{7, 3, 4, 8, uint(12 * math.Log2(10))}},
		{"misc", `
		type S struct {}
		type E = <-chan int

		func x() { if true || (false && true) {} }
		func y(a *int) { *a = 0 }
		func z(a []int) []int { return a[1:]}
		func a(b interface{}) int {
			switch b.(type) {
			case int:
				return b.(int)
			}
		}
		`, Result{19, 13, 26, 39, uint(65 * math.Log2(32))}},
	}

	for _, test := range tests {
		got := getMetric(test.Source)
		if test.Want != got {
			t.Errorf("%v: result{N1Distinct, N2Distinct, N1Total, N2Total} = %v, want %v", test.Name, got, test.Want)
		}
	}
}
