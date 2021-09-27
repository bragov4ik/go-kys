package kys

import "go/ast"

func WMFPcalc(file *ast.File, score map[string]int) {
	ast.Inspect(file, func(n ast.Node) bool {
		// Find function declarations
		_, ok := n.(*ast.FuncDecl)
		if ok {
			score["funcDecl"]++
			return true
		}
		// Find return statements
		_, ok = n.(*ast.ReturnStmt)
		if ok {
			score["returnStmt"]++
			return true
		}
		// Find function calls
		_, ok = n.(*ast.CallExpr)
		if ok {
			score["callExpr"]++
			return true
		}
		// Find assignment statements
		_, ok = n.(*ast.AssignStmt)
		if ok {
			score["assignStmt"]++
		}
		return true
	})
}
