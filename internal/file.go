package internal

import (
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"os"
	"strings"
)

func ClearDir() {
	os.RemoveAll("dist")
	os.Mkdir("dist", os.ModePerm)
}

func GetFiles(path string) error {
	dir, err := os.ReadDir(path)
	if err != nil {
		fmt.Println(err)
		return err
	}

	for _, file := range dir {
		fileName := file.Name()

		pwd := path + "/" + fileName

		if file.IsDir() {
			GetFiles(pwd)
		} else if strings.HasSuffix(fileName, ".go") {
			// to string

			if err := MigrateCode(pwd); err != nil {
				fmt.Println("MigrateCode", err)
				return err
			}
		}
	}
	return nil
}

func WriteFile(pwd string, node *ast.File) error {
	// Print the modified AST.

	pwds := strings.Split(pwd, "/")
	listDir := pwds[0 : len(pwds)-1]

	for i := 1; i < len(listDir); i++ {
		// fmt.Println("dist/" + strings.Join(listDir[0:i], "/"))
	}
	os.MkdirAll("dist/"+strings.Join(listDir, "/"), os.ModePerm)

	f, err := os.Create("dist/" + pwd)
	if err != nil {
		return err
	}
	defer f.Close()
	fset := token.NewFileSet()
	printer.Fprint(f, fset, node)
	return nil
}
