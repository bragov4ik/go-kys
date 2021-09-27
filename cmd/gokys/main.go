package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
)

func main() {
	flag.Parse()
	args := flag.Args()

	var files []string

	for _, v := range args {
		info, err := os.Stat(v)

		if os.IsNotExist(err) {
			log.Fatal("File does not exist.")
		}

		if info.IsDir() {
			files, err = WalkMatch(v, "*.go")
		} else {
			files = append(files, v)
		}

		if err != nil {
			log.Fatal(err)
		}
	}

	fset := token.NewFileSet()
	scores := make(map[string]int)

	for _, file := range files {
		node, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
		if err != nil {
			log.Fatal(err)
		}
		WMFPcalc(node, scores)
	}

	fmt.Println(scores)

}

func WalkMatch(root, pattern string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}

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
