package inject

import (
	"fmt"
	"go/token"
	"log"
	"strconv"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
)

// 插入 import  ImportAlias ImportPackage
type ImportMeta struct {
	ImportAlias   string //导入别名
	ImportPackage string //导入包路径

	fset        *token.FileSet
	hasImported bool
	node        *dst.ImportSpec
}

func (vi *ImportMeta) Visit(node dst.Node) dst.Visitor {
	switch genDecl := node.(type) {
	case *dst.GenDecl:
		if genDecl.Tok == token.IMPORT {
			// 判断是否已经import
			for _, v := range genDecl.Specs {
				importSpec := v.(*dst.ImportSpec)
				// 如果已经包含
				if importSpec.Path.Value == strconv.Quote(vi.ImportPackage) {
					vi.hasImported = true
				}
			}
			if !vi.hasImported {
				importSp := vi.GetNode()
				//放到最后
				genDecl.Specs = append(genDecl.Specs, importSp.(*dst.ImportSpec))
				vi.hasImported = true
			}
		}
	}
	return vi
}

func (vi *ImportMeta) RollBack(node dst.Node) dst.Visitor {
	if genDecl, ok := node.(*dst.GenDecl); ok {
		if genDecl.Tok == token.IMPORT {
			//dst.Println(token.NewFileSet(),genDecl)
			// 是否已经import
			for k, v := range genDecl.Specs {
				importSpec := v.(*dst.ImportSpec)
				// 如果已经包含
				if importSpec.Path.Value == strconv.Quote(vi.ImportPackage) {
					//删除
					genDecl.Specs = append(append([]dst.Spec{}, genDecl.Specs[:k]...), genDecl.Specs[k+1:]...)
					log.Println("delete ", importSpec.Path.Value)
				}
			}

		}
	}
	return vi
}

func (vi *ImportMeta) GetNode() dst.Node {
	if vi.node != nil {
		return vi.node
	}

	importSpec := &dst.ImportSpec{
		Name: &dst.Ident{
			Name: vi.ImportAlias,
		},
		Path: &dst.BasicLit{
			Kind:  token.STRING,
			Value: strconv.Quote(vi.ImportPackage),
		},
	}
	return importSpec
}

func (vi *ImportMeta) GetCode() string {
	injectCode := fmt.Sprintf("import %v %v", vi.ImportAlias, vi.ImportPackage)

	return injectCode
}

func NewImportMete(importCode string) *ImportMeta {
	meta := &ImportMeta{}
	input := fmt.Sprintf(`
package main
import (%s)
`, importCode)
	fParse, err := decorator.ParseFile(token.NewFileSet(), "", input, 0)
	if err != nil {
		fmt.Println("Error parsing file:", err)
	}
	dst.Inspect(fParse, func(n dst.Node) bool {
		switch x := n.(type) {
		case *dst.ImportSpec:
			//dst.Print(x)
			if x.Name != nil {
				meta.ImportAlias = x.Name.Name
			}
			if x.Path != nil {
				meta.ImportPackage = x.Path.Value
			}
			meta.node = x
		}
		return true
	})

	//log.Println("GetCode", meta.GetCode())
	//log.Println("NewImportMete", jsonconv.AnyToJsonIndent(meta))
	return meta
}
