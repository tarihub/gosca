package config

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
)

type AstFile struct {
	FilePath    string
	FileContent []byte
	FileAst     *ast.File
	Imports     []*ast.ImportSpec
}

func (a *AstFile) getFileContent() error {
	c, err := ioutil.ReadFile(a.FilePath)
	if err != nil {
		return err
	}
	a.FileContent = c
	return nil
}

func (a *AstFile) getFileAst() error {
	fSet := token.NewFileSet()
	err := error(nil)
	a.FileAst, err = parser.ParseFile(fSet, a.FilePath, a.FileContent, parser.ParseComments)
	if err != nil {
		return err
	}
	return nil
}

func (a *AstFile) getImports() {
	ast.Inspect(a.FileAst, func(n ast.Node) bool {
		if n != nil {
			switch n.(type) {
			case *ast.ImportSpec:
				a.Imports = append(a.Imports, n.(*ast.ImportSpec))
			}
		}
		return true
	})
}

func (a *AstFile) Load() error {
	err := a.getFileContent()
	if err != nil {
		return err
	}
	err = a.getFileAst()
	if err != nil {
		return err
	}
	a.getImports()
	return nil
}
