package arithmetic

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestArithmeticComp(t *testing.T) {
	type args struct {
		n ast.Node
	}
	testcases := []struct {
		src  string
		want float64
	}{
		{`package main; func main(){ a := 1+2-3*4/5%6 }`, 5},
		{`package main; func main(){ a := -0; a+=1; a-=2; a*=3; a%=4; a/=5 }`, 6},
		{`package main; func main(){ a := +0; a++; a++; a--; a++; a--; a++ }`, 7},
	}
	cfg := Weights{Add: 1, Sub: 1, Mul: 1, Quo: 1, Rem: 1, AddAssign: 1, SubAssign: 1, MulAssign: 1, QuoAssign: 1, RemAssign: 1, Inc: 1, Dec: 1}
	tests := make([]struct {
		name string
		m    *Metric
		args args
		want float64
	}, len(testcases))
	for i, testcase := range testcases {
		file, err := parser.ParseFile(token.NewFileSet(), "", testcase.src, parser.ParseComments)
		if err != nil {
			t.Fatal(err)
		}
		tests[i] = struct {
			name string
			m    *Metric
			args args
			want float64
		}{fmt.Sprintf("Test%v", i+1), &Metric{Config: cfg}, args{file}, testcase.want}
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast.Inspect(tt.args.n, func(n ast.Node) bool {
				tt.m.ParseNode(n)
				return true
			})
			if got := tt.m.Finish(); got != tt.want {
				t.Errorf("Metric.Finish() = %v, want %v", got, tt.want)
			}
		})
	}
}
