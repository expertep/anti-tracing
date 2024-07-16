package internal

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

func MigrateCode(pwd string) error {
	src, err := os.ReadFile(pwd)
	if err != nil {
		fmt.Println("ReadFile", err)
		return err
	}

	fset := token.NewFileSet()
	// node, err := parser.ParseFile(fset, "", string(src), parser.DeclarationErrors)
	node, err := parser.ParseFile(fset, "", string(src), parser.AllErrors)
	if err != nil {
		return err
	}

	ExamCode(node)
	WriteFile(pwd, node)

	return nil
}

type ProcessData struct {
	MapFileNames map[string]string
}

func processData() *ProcessData {
	return &ProcessData{
		MapFileNames: make(map[string]string),
	}

}
func ExamCode(node *ast.File) {
	// Traverse and modify the AST.
	processData := processData()

	AddImport(node)
	ast.Inspect(node, func(n ast.Node) bool {
		// Example transformation: rename function "main" to "Main"

		// check comment
		/* if fn, ok := n.(*ast.CommentGroup); ok {
			for _, comment := range fn.List {
				fmt.Printf("Comment: %s\n", comment.Text)
			}
		} else  */
		if fn, ok := n.(*ast.FuncDecl); ok {

			isMainController := isGinHandlerReturnType(fn)

			if isMainController {
				/* if hasDoComment(fn) {
				} */
				AddCommandStartParentTrace(fn)
			} else {

				processData.AddCommandStartTrace(fn)

			}

			// AddCommandErrorSpan(fn)

		}

		return true
	})
	fmt.Println("MapFileNames", processData.MapFileNames)
	ast.Inspect(node, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok {
			processData.AddParam(fn)
		}
		return true
	})
}

func AddImport(node *ast.File) {
	addFmtImport := true
	for _, imp := range node.Imports {
		if imp.Path.Value == "\"fmt\"" {
			addFmtImport = false
			break
		}
	}

	if addFmtImport {
		// Add `import "fmt"` to the file's import section.
		// alias fmt
		newImport := &ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: "\"fmt\"",
			},
			Name: &ast.Ident{
				Name: "fmt",
			},
		}
		node.Imports = append(node.Imports, newImport)

		// Ensure `import` declaration exists.
		if node.Decls == nil || len(node.Decls) == 0 {
			node.Decls = append(node.Decls, &ast.GenDecl{
				Tok:    token.IMPORT,
				Lparen: token.NoPos,
				Specs:  []ast.Spec{newImport},
			})
		} else {
			// Insert into existing import declarations or create a new one.
			importFound := false
			for _, decl := range node.Decls {
				if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.IMPORT {
					genDecl.Specs = append(genDecl.Specs, newImport)
					importFound = true
					break
				}
			}
			if !importFound {
				node.Decls = append([]ast.Decl{
					&ast.GenDecl{
						Tok:    token.IMPORT,
						Lparen: token.NoPos,
						Specs:  []ast.Spec{newImport},
					},
				}, node.Decls...)
			}
		}
	}

}
func (processData *ProcessData) AddParam(fn *ast.FuncDecl) {
	// fmt.Println("AddParam", fn.Name.Name)
	for _, v := range fn.Body.List {
		if stmt, ok := v.(*ast.ExprStmt); ok {
			if callExpr, ok := stmt.X.(*ast.CallExpr); ok {
				if _, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
					// fmt.Println("ExprStmt", selExpr.Sel, selExpr.Sel.Name)
					/* if _, ok := processData.MapFileNames[selExpr.Sel.Name]; ok {

						callExpr.Args = append(callExpr.Args, &ast.BasicLit{
							Value: "ot",
						})
					} */
				}
			}
		} else if stmt, ok := v.(*ast.AssignStmt); ok {
			//add param ot
			for _, rhs := range stmt.Rhs {
				// show code at this line
				fmt.Println("AssignStmt", rhs)
				if callExpr, ok := rhs.(*ast.CallExpr); ok {
					// fmt.Println("AssignStmt", callExpr.)
					fmt.Println("AssignStmt", callExpr.Args)
					if ident, ok := callExpr.Fun.(*ast.Ident); ok {
						if _, ok := processData.MapFileNames[ident.Name]; ok {
							// Add two new arguments to the function call in the assignment.
							callExpr.Args = append(callExpr.Args,
								&ast.BasicLit{
									Kind:  token.INT,
									Value: "20", // Example argument for `int`
								},
								&ast.BasicLit{
									Kind:  token.STRING,
									Value: "\"extra\"", // Example argument for `string`
								},
							)
						}
					}
				}
				/* rhs.(*ast.CallExpr).Args = append(rhs.(*ast.CallExpr).Args, &ast.BasicLit{

				} */
			}
		}

	}
}
func (processData *ProcessData) AddCommandStartTrace(fn *ast.FuncDecl) {
	/* defer func(name string, value reflect.Value) bool {

	}()  */
	cmdCreateTrace := &ast.ExprStmt{
		X: &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   ast.NewIdent("fmt"),
				Sel: ast.NewIdent("Println"),
			},
			Args: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf("\"%s\"", fn.Name.Name),
				},
			},
		},
	}
	cmdCloseTrace := &ast.DeferStmt{
		Call: &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   ast.NewIdent("fmt"),
				Sel: ast.NewIdent("Println"),
			},
			Args: []ast.Expr{
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf("\"close %s\"", fn.Name.Name),
				},
			},
		},
	}

	fn.Body.List = InsertSlice(fn.Body.List, []ast.Stmt{
		cmdCreateTrace,
		cmdCloseTrace,
	}, 0)

	// modify add 1 arg for func
	if len(fn.Type.Params.List) >= 1 {
		// Modify the function to add two more parameters.
		param1 := fn.Type.Params.List[0] // Existing parameter
		param2 := &ast.Field{
			Names: []*ast.Ident{
				ast.NewIdent("ot"),
			},
			Type: ast.NewIdent("otelCus.IOtelPattern"),
		}
		fn.Type.Params.List = []*ast.Field{param1, param2}
		processData.MapFileNames[fn.Name.Name] = ""
	}

}
func AddCommandStartParentTrace(fn *ast.FuncDecl) {
	for i, stmt := range fn.Body.List {
		if retStmt, ok := stmt.(*ast.ReturnStmt); ok {
			if len(retStmt.Results) == 1 {
				if funcLit, ok := retStmt.Results[0].(*ast.FuncLit); ok {

					//  init variable
					cmdInitCtx := &ast.AssignStmt{
						Lhs: []ast.Expr{
							ast.NewIdent("ctx"),
						},
						Tok: token.DEFINE,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: ast.NewIdent("context.Background"),
							},
						},
					}

					cmdCreateTrace := &ast.AssignStmt{
						Lhs: []ast.Expr{
							ast.NewIdent("ot"),
						},
						Tok: token.DEFINE,
						Rhs: []ast.Expr{
							&ast.CallExpr{
								Fun: ast.NewIdent("otelCus.CreateOtel"),
								Args: []ast.Expr{
									&ast.BasicLit{
										// Kind:  token.IDENT,
										Value: "ctx",
									},
									&ast.BasicLit{
										Kind:  token.STRING,
										Value: fmt.Sprintf("\"%s\"", fn.Name.Name),
									},
									&ast.BasicLit{
										Kind:  token.STRING,
										Value: fmt.Sprintf("\"%s\"", fn.Name.Name),
									},
								},
							},
						},
					}

					printStmt := &ast.ExprStmt{
						X: &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("fmt"),
								Sel: ast.NewIdent("Println"),
							},
							Args: []ast.Expr{
								&ast.BasicLit{
									Kind:  token.STRING,
									Value: fmt.Sprintf("\"%s\"", fn.Name.Name),
								},
							},
						},
					}
					cmdCloseTrace := &ast.DeferStmt{
						Call: &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("fmt"),
								Sel: ast.NewIdent("Println"),
							},
							Args: []ast.Expr{
								&ast.BasicLit{
									Kind:  token.STRING,
									Value: fmt.Sprintf("\"close %s\"", fn.Name.Name),
								},
							},
						},
					}
					// funcLit.Body.List = append([]ast.Stmt{printStmt}, funcLit.Body.List...)
					funcLit.Body.List = InsertSlice(funcLit.Body.List, []ast.Stmt{
						cmdInitCtx,
						cmdCreateTrace,
						printStmt,
						cmdCloseTrace,
					}, 0)
					fn.Body.List[i] = &ast.ReturnStmt{Results: []ast.Expr{funcLit}}

				}
			}
		}
	}
	/* fn.Body.List = InsertSlice(fn.Body.List, []ast.Stmt{
		cmdCreateTrace,
		cmdCloseTrace,
	}, 0) */

}
func AddCommandErrorSpan(fn *ast.FuncDecl) {
	for index, v := range fn.Body.List {
		if _, ok := v.(*ast.ReturnStmt); ok {
			// fmt.Println("return", v.Results)
			cmdErrorTrace := []ast.Stmt{
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   ast.NewIdent("fmt"),
							Sel: ast.NewIdent("Println"),
						},
						Args: []ast.Expr{
							&ast.BasicLit{
								Kind:  token.STRING,
								Value: `"error trace"`,
							},
						},
					},
				},
			}
			fn.Body.List = InsertSlice(fn.Body.List, cmdErrorTrace, index)
		}
	}
}
func InsertSlice(target []ast.Stmt, insert []ast.Stmt, index int) []ast.Stmt {
	if index < 0 || index > len(target) {
		// Handling invalid index, returning the target slice unchanged
		return target
	}

	// Create a new slice to hold the result
	result := make([]ast.Stmt, 0, len(target)+len(insert))

	// Add elements from the target slice before the index
	result = append(result, target[:index]...)

	// Add elements from the insert slice
	result = append(result, insert...)

	// Add remaining elements from the target slice after the index
	result = append(result, target[index:]...)

	return result
}

func isGinHandlerReturnType(fn *ast.FuncDecl) bool {
	if fn.Type.Results == nil || len(fn.Type.Results.List) != 1 {
		return false
	}

	// Check if the return type is a function.
	if funcType, ok := fn.Type.Results.List[0].Type.(*ast.FuncType); ok {
		// Check if the function has exactly one parameter.
		if len(funcType.Params.List) == 1 {
			// Check if the parameter is `*gin.Context`.
			if starExpr, ok := funcType.Params.List[0].Type.(*ast.StarExpr); ok {
				if selExpr, ok := starExpr.X.(*ast.SelectorExpr); ok {
					if selExpr.Sel.Name == "Context" {
						if id, ok := selExpr.X.(*ast.Ident); ok && id.Name == "gin" {
							return true
						}
					}
				}
			}
		}
	}
	return false
}
func hasDoComment(fn *ast.FuncDecl) bool {
	if fn.Doc != nil {
		for _, comment := range fn.Doc.List {

			if comment.Text == "// otel" {
				return true
			}
		}
	}
	return false
}
