package inline

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestInlineData(t *testing.T) {
	type args struct {
		n ast.Node
	}
	testcases := []struct {
		src  string
		want float64
	}{
		{`package main

	import (
		"fmt"
	)
		

	func foo(t int) float64 {
		a := 123
		b := 1.23
		return float64(a+t) + b
	}
			
	func main(){
		fmt.Printf("Result: %v", foo(1))
	}`, 23},
		{`package main

	import (
		"fmt"
	)
	
	func foo(t int) complex128 {
		a := 'a'
		b := 32i
		result := (complex(float64(t), 0) + b)
		fmt.Printf("%s = %v\n", string(a), result)
		return result
	}
	
	func main() {
		fmt.Printf("Result: %v\n", foo(1))
	}`, 37},
		{`package main

	import (
		"fmt"
	)
	
	func foo() {
		a := struct{Num uint; Value string}{10, "val",};
		fmt.Printf("%v = {%v, %v}\n", "a", a.Num, a.Value)
	}
	
	func main() {
		foo()
	}`, 32},
	}
	cfg := Weights{Int: 2, Float: 2, Imag: 2, Char: 1, String: 1}
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
		}{fmt.Sprintf("Test%v", i), &Metric{Config: cfg}, args{file}, testcase.want}
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
