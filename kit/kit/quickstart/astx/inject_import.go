package astx

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"strconv"

	"github.com/ve-weiyi/ve-blog-golang/kit/utils/jsonconv"
)

// 插入 import  ImportAlias ImportPackage
type ImportMeta struct {
	fset          *token.FileSet
	InjectCode    string
	ImportAlias   string //导入别名
	ImportPackage string //导入包路径
	HasDone       bool
}

func (vi *ImportMeta) Visit(node ast.Node) ast.Visitor {
	switch genDecl := node.(type) {
	case *ast.GenDecl:
		if genDecl.Tok == token.IMPORT {
			//ast.Println(token.NewFileSet(),genDecl)
			// 是否已经import
			hasImported := false
			for _, v := range genDecl.Specs {
				importSpec := v.(*ast.ImportSpec)
				// 如果已经包含
				if importSpec.Path.Value == strconv.Quote(vi.ImportPackage) {
					hasImported = true
				}
			}
			if !hasImported && !vi.HasDone {
				importSp := &ast.ImportSpec{
					Name: &ast.Ident{
						Name: vi.ImportAlias,
					},
					Path: &ast.BasicLit{
						Kind:  token.STRING,
						Value: strconv.Quote(vi.ImportPackage),
					},
				}
				//放到最后
				genDecl.Specs = append(genDecl.Specs, importSp)
				vi.HasDone = true
				//log.Println("ImportPackage end", vi.ImportPackage, vi.hasImported)
				//ast.Println(token.NewFileSet(), node)
			}
			//log.Println("--", jsonconv.AnyToJsonIndent(genDecl.Specs))
		} else {

		}
	}
	return vi
}

func (vi *ImportMeta) RollBack(node ast.Node) ast.Visitor {
	if genDecl, ok := node.(*ast.GenDecl); ok {
		if genDecl.Tok == token.IMPORT {
			//ast.Println(token.NewFileSet(),genDecl)
			// 是否已经import
			for k, v := range genDecl.Specs {
				importSpec := v.(*ast.ImportSpec)
				// 如果已经包含
				if importSpec.Path.Value == strconv.Quote(vi.ImportPackage) {
					//删除
					genDecl.Specs = append(append([]ast.Spec{}, genDecl.Specs[:k]...), genDecl.Specs[k+1:]...)
					genDecl.Specs = nil
				}
			}

		}
	}
	return vi
}

func (vi *ImportMeta) RollBackImportSpec(file *ast.File) ast.Visitor {
	var newImport []*ast.ImportSpec
	for _, imp := range file.Imports {
		if imp.Path.Value == strconv.Quote(vi.ImportPackage) {
			newImport = append(newImport, imp)
		}
	}
	file.Imports = newImport
	//log.Println("RollBackImportSpec", vi.ImportPackage)
	//ast.Println(token.NewFileSet(), file.Imports)
	return vi
}

func NewImportMete(importCode string) {
	meta := &ImportMeta{}
	input := fmt.Sprintf("package main\nimport (%s)", importCode)
	fParse, err := parser.ParseFile(token.NewFileSet(), "", input, 0)
	if err != nil {
		fmt.Println("Error parsing file:", err)
	}
	ast.Print(token.NewFileSet(), fParse)
	ast.Inspect(fParse, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.ImportSpec:
			meta.ImportAlias = x.Name.Name
			meta.ImportPackage = x.Path.Value
		}
		return true
	})
	log.Println("NewImportMete", jsonconv.AnyToJsonIndent(meta))
}
