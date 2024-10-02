package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
)

var functions = make(map[string]*ast.FuncDecl)

func printFullFuncDecl(funcDecl *ast.FuncDecl, fset *token.FileSet) string {
	var buf bytes.Buffer
	// Use the printer package to format the function declaration
	err := printer.Fprint(&buf, fset, funcDecl)
	if err != nil {
		log.Printf("Failed to print function declaration: %v", err)
		return ""
	}
	fmt.Println(buf.String()) // Print the complete function declaration
	return buf.String()
}

func findFunctions(node ast.Node) {
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			functions[x.Name.Name] = x
		}
		return true
	})
}

func replaceGoRoutines() {
	// TODO replace goroutine call with a function call that:
	// registers an HTTP handler to listen for channel responses'
	// creates a pod, runs it
	// gets pod IP address
	// makes a POST request to the pod
}

func findGoroutines(node ast.Node, fset *token.FileSet) {
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.GoStmt:
			fmt.Printf("Found a goroutine at line %d\n", x.Pos())
			fmt.Printf("Goroutine: %v\n", x)

			if callExpr, ok := x.Call.Fun.(*ast.Ident); ok {
				if function, ok := functions[callExpr.Name]; ok {
					fmt.Printf("The goroutine is calling the function %s\n", function.Name.Name)
					fstring := printFullFuncDecl(function, fset)
					replaceGoRoutines()
					if err := generateServerFile(function.Name.Name, fstring); err != nil {
						log.Fatalf("Error generating server file: %s", err)
					}
				}
			}
		}
		return true
	})
}

func main() {
	filename := flag.String("file", "", "The Go source file to parse")
	fset := token.NewFileSet()

	flag.Parse()

	// parse the source file into an AST
	node, err := parser.ParseFile(fset, *filename, nil, parser.AllErrors)
	if err != nil {
		log.Fatalf("Error parsing file: %s", err)
		fmt.Println(err)
		return
	}

	findFunctions(node)
	findGoroutines(node, fset)
}
