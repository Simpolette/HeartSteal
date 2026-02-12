package main

import (
	"fmt"
	"go/ast"
	"os"
	"strings"
	"unicode"

	"golang.org/x/tools/go/packages"
)

func main() {
	cfg := &packages.Config{
		Mode:  packages.NeedName | packages.NeedSyntax | packages.NeedFiles | packages.NeedTypes,
		Tests: true,
	}

	pkgs, err := packages.Load(cfg, "./internal/usecase/...")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading packages: %v\n", err)
		os.Exit(1)
	}

	// We gather ALL existing tests from ALL packages into one set.
	// This solves the issue where tests are in a sub-directory.
	existingTests := make(map[string]bool)

	for _, pkg := range pkgs {
		for _, file := range pkg.Syntax {
			// strictly check if it's a test file
			filename := pkg.Fset.File(file.Pos()).Name()
			if !strings.HasSuffix(filename, "_test.go") {
				continue
			}

			for _, decl := range file.Decls {
				if fn, ok := decl.(*ast.FuncDecl); ok {
					existingTests[fn.Name.Name] = true
				}
			}
		}
	}

	missingCount := 0
	

	for _, pkg := range pkgs {
		// Skip compiled test binaries (pkg.ID ends in .test])
		// But DO NOT skip packages named "test" or "usecase_test" if they contain source code
		if strings.HasSuffix(pkg.ID, ".test]") {
			continue
		}

		for _, file := range pkg.Syntax {
			filename := pkg.Fset.File(file.Pos()).Name()
			
			// Skip test files
			if strings.HasSuffix(filename, "_test.go") {
				continue
			}
			for _, decl := range file.Decls {
				fn, ok := decl.(*ast.FuncDecl)
				// Must be a function, and must be Exported
				if !ok || !fn.Name.IsExported() {
					continue
				}
				
				funcName := fn.Name.Name
				var expectedTestName string

				// Must be a Method (has a Receiver)
				if fn.Recv == nil || len(fn.Recv.List) < 1 {
					continue
				}

				// Attempt to get Struct Name
				typeExpr := fn.Recv.List[0].Type
				var structName string
				
				// Handle pointer receiver (*UserUsecase) vs value receiver (UserUsecase)
				if star, ok := typeExpr.(*ast.StarExpr); ok {
					if ident, ok := star.X.(*ast.Ident); ok {
						structName = Capitalize(ident.Name)
					}
				} else if ident, ok := typeExpr.(*ast.Ident); ok {
					structName = Capitalize(ident.Name)
				}

				// Naming convention for methods: TestStruct_Method
				expectedTestName = fmt.Sprintf("Test%s_%s", structName, funcName)

				if !existingTests[expectedTestName] {
					fmt.Printf("❌ Missing test for: %s \n   File: %s\n   Expected: %s\n", 
						funcName, filename, expectedTestName)
					missingCount++
				}
			}
		}
	}

	if missingCount > 0 {
		fmt.Printf("\nFinished: %d missing tests found.\n", missingCount)
		os.Exit(1)
	}

	fmt.Println("✅ All exported functions have matching test cases.")
}

func Capitalize(s string) string {
    if s == "" {
        return ""
    }
    r := []rune(s)
    
    // Uppercase the first character
    r[0] = unicode.ToUpper(r[0])
    
    return string(r)
}