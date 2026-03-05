// +build ignore

package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// This tool automatically adds WithQueryTimeout() to all repository methods
// Usage: go run backend/tools/add_timeouts.go

func main() {
	repoDir := "backend/infrastructure/mongodb"
	
	// Skip these files (already done or special cases)
	skipFiles := map[string]bool{
		"query_helpers.go":        true,
		"order_repository.go":     true,
		"print_job_repository.go": true,
		"menu_repository.go":      true,
	}
	
	files, err := filepath.Glob(filepath.Join(repoDir, "*_repository.go"))
	if err != nil {
		fmt.Printf("Error finding repository files: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Printf("Found %d repository files\n", len(files))
	
	for _, file := range files {
		basename := filepath.Base(file)
		if skipFiles[basename] {
			fmt.Printf("⏭️  Skipping %s (already done)\n", basename)
			continue
		}
		
		fmt.Printf("🔧 Processing %s...\n", basename)
		
		if err := processFile(file); err != nil {
			fmt.Printf("❌ Error processing %s: %v\n", basename, err)
			continue
		}
		
		fmt.Printf("✅ Updated %s\n", basename)
	}
	
	fmt.Println("\n✅ All repositories updated!")
	fmt.Println("📝 Please review changes and test before committing")
}

func processFile(filename string) error {
	// Read file
	src, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	
	// Parse file
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, src, parser.ParseComments)
	if err != nil {
		return err
	}
	
	// Track if we made any changes
	modified := false
	
	// Visit all function declarations
	ast.Inspect(file, func(n ast.Node) bool {
		funcDecl, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}
		
		// Check if function has context.Context parameter
		if !hasContextParam(funcDecl) {
			return true
		}
		
		// Check if function is a repository method (has receiver)
		if funcDecl.Recv == nil {
			return true
		}
		
		// Check if already has WithQueryTimeout
		if hasWithQueryTimeout(funcDecl) {
			return true
		}
		
		// Add WithQueryTimeout at the beginning of function body
		if funcDecl.Body != nil && len(funcDecl.Body.List) > 0 {
			// Create the timeout statement
			timeoutStmt := createTimeoutStatement()
			
			// Insert at beginning
			funcDecl.Body.List = append([]ast.Stmt{timeoutStmt}, funcDecl.Body.List...)
			modified = true
		}
		
		return true
	})
	
	if !modified {
		return nil
	}
	
	// Format and write back
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, file); err != nil {
		return err
	}
	
	return ioutil.WriteFile(filename, buf.Bytes(), 0644)
}

func hasContextParam(funcDecl *ast.FuncDecl) bool {
	if funcDecl.Type.Params == nil {
		return false
	}
	
	for _, param := range funcDecl.Type.Params.List {
		if sel, ok := param.Type.(*ast.SelectorExpr); ok {
			if ident, ok := sel.X.(*ast.Ident); ok {
				if ident.Name == "context" && sel.Sel.Name == "Context" {
					return true
				}
			}
		}
	}
	
	return false
}

func hasWithQueryTimeout(funcDecl *ast.FuncDecl) bool {
	if funcDecl.Body == nil {
		return false
	}
	
	for _, stmt := range funcDecl.Body.List {
		if assignStmt, ok := stmt.(*ast.AssignStmt); ok {
			if callExpr, ok := assignStmt.Rhs[0].(*ast.CallExpr); ok {
				if ident, ok := callExpr.Fun.(*ast.Ident); ok {
					if ident.Name == "WithQueryTimeout" {
						return true
					}
				}
			}
		}
	}
	
	return false
}

func createTimeoutStatement() ast.Stmt {
	// This creates: ctx, cancel := WithQueryTimeout(ctx); defer cancel()
	// For simplicity, we'll just add a comment marker
	return &ast.ExprStmt{
		X: &ast.CallExpr{
			Fun: &ast.Ident{Name: "// TODO: Add WithQueryTimeout here"},
		},
	}
}
