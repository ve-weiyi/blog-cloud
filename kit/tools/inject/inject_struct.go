package inject

import (
	"fmt"
	"go/parser"
	"go/token"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
)

// 插入 FieldName  FieldType.SelName
type StructMeta struct {
	FindStructName string //在 struct FindStructName{ //插入 FieldName  IdentName.SelName  }
	FieldName      string //元素别名
	FieldIdent     string //属性类型引用的包名
	FieldType      string //属性类型
	FieldComment   string //属性注释

	fset     *token.FileSet
	hasField bool
	node     *dst.Field
}

func (vi *StructMeta) Visit(node dst.Node) dst.Visitor {
	switch n := node.(type) {
	case *dst.TypeSpec:
		//log.Printf("Type: %s %s\n", n.Group.Group, vi.FindStructName)
		// 找到结构体名称
		if n.Name.Name == vi.FindStructName {
			switch t := n.Type.(type) {
			case *dst.StructType:
				//数据已经存在，不重复添加
				for _, f := range t.Fields.List {
					if f.Names != nil {
						if f.Names[0].Name == vi.FieldName {
							return vi
						}
					}
				}
				f := vi.GetNode()
				//放到最后
				t.Fields.List = append(t.Fields.List, f.(*dst.Field))
			}
		}
	}
	return vi
}

func (vi *StructMeta) RollBack(node dst.Node) dst.Visitor {
	if genDecl, ok := node.(*dst.GenDecl); ok {
		if genDecl.Tok == token.TYPE {
			//dst.Println(token.NewFileSet(),genDecl)
			for i := range genDecl.Specs {
				switch n := genDecl.Specs[i].(type) {
				case *dst.TypeSpec:
					// 找到结构体名称
					if n.Name.Name == vi.FindStructName {
						switch t := n.Type.(type) {
						case *dst.StructType:
							//数据已经存在，不重复添加
							for k, f := range t.Fields.List {
								if f.Names != nil {
									if f.Names[0].Name == vi.FieldName {
										//删除
										block := t.Fields
										block.List = append(append([]*dst.Field{}, block.List[:k]...), block.List[k+1:]...)
										return vi
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return vi
}

func (vi *StructMeta) GetNode() dst.Node {
	if vi.node != nil {
		return vi.node
	}
	f := &dst.Field{
		// 变量列表，可以有多个 var str1,str2 string
		Names: []*dst.Ident{
			{
				Name: vi.FieldName,
				Obj: &dst.Object{
					Kind: dst.Var,
					Name: vi.FieldName,
				},
			},
		},
		Type: &dst.SelectorExpr{
			X: &dst.Ident{
				Name: vi.FieldIdent,
			},
			Sel: &dst.Ident{
				Name: vi.FieldType,
			},
		},
		//添加注释
		Decs: dst.FieldDecorations{
			NodeDecs: dst.NodeDecs{
				Before: dst.NewLine,
				Start:  nil,
				End: dst.Decorations{
					fmt.Sprintf("//%s", vi.FieldComment),
				},
				After: dst.NewLine,
			},
			Type: nil,
		},
	}
	return f
}

func (vi *StructMeta) GetCode() string {
	injectCode := fmt.Sprintf("%s", vi.FieldType)

	if vi.FieldIdent != "" {
		injectCode = fmt.Sprintf("%s.%s", vi.FieldIdent, injectCode)
	}

	if vi.FieldName != "" {
		injectCode = fmt.Sprintf("%s %s", vi.FieldName, injectCode)
	}

	if vi.FieldComment != "" {
		injectCode = fmt.Sprintf("%s //%s", injectCode, vi.FieldComment)
	}

	return injectCode
}

func NewStructMete(findStructName, importCode string) *StructMeta {
	meta := &StructMeta{}
	input := fmt.Sprintf(`
package main

type %s struct {
	%s
}
`, findStructName, importCode)
	fParse, err := decorator.ParseFile(token.NewFileSet(), "", input, parser.ParseComments)
	if err != nil {
		fmt.Println("Error parsing file:", err)
		return nil
	}

	dst.Inspect(fParse, func(n dst.Node) bool {
		switch x := n.(type) {
		case *dst.TypeSpec:
			if x.Name.Name != findStructName {
				break
			}
			//dst.Print(x)
			structType, ok := x.Type.(*dst.StructType)
			if !ok {
				return false
			}

			for _, field := range structType.Fields.List {

				fieldName := field.Names[0].Name
				// Extract field ident
				fieldIdent := ""
				fieldType := ""
				switch t := field.Type.(type) {
				case *dst.Ident:
					fieldIdent = t.Name
				case *dst.SelectorExpr:
					fieldIdent = t.X.(*dst.Ident).Name
					fieldType = t.Sel.Name
				case *dst.StarExpr: //指针类型
					fieldIdent = "*" + t.X.(*dst.SelectorExpr).X.(*dst.Ident).Name
					fieldType = t.X.(*dst.SelectorExpr).Sel.Name
				}

				// Extract field comment
				fieldComment := ""
				if len(field.Decorations().End.All()) > 0 {
					fieldComment = field.Decorations().End.All()[0]
				}

				meta.FindStructName = findStructName
				meta.FieldName = fieldName
				meta.FieldIdent = fieldIdent
				meta.FieldType = fieldType
				meta.FieldComment = fieldComment
				meta.node = field
				//log.Printf("Field Group: %s, Field Ident: %s, Field Type: %s, Field ColumnComment: %s\n", fieldName, fieldIdent, fieldType, fieldComment)
			}
		}

		return true
	})

	//log.Println("GetCode", meta.GetCode())
	//log.Println("NewStructMete", jsonconv.AnyToJsonIndent(meta))
	return meta
}
