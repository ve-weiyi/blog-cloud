package astx

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"regexp"
	"strings"

	"github.com/ve-weiyi/ve-blog-golang/pkg/utils/jsonconv"
)

// 插入 FieldName  FieldType.SelName
type StructMeta struct {
	fset           *token.FileSet
	InjectCode     string
	FindStructName string //在 struct FindStructName{ //插入 FieldName  IdentName.SelName  }
	FieldName      string //元素别名
	FieldIdent     string //属性类型引用的包名
	FieldType      string //属性类型
	FieldComment   string //属性注释
}

func (vi *StructMeta) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.TypeSpec:
		//log.Printf("Type: %s %s\n", n.Group.Group, vi.FindStructName)
		// 找到结构体名称
		if n.Name.Name == vi.FindStructName {
			switch t := n.Type.(type) {
			case *ast.StructType:
				//数据已经存在，不重复添加
				for _, f := range t.Fields.List {
					if f.Names != nil {
						if f.Names[0].Name == vi.FieldName {
							return vi
						}
					}
				}
				f := &ast.Field{
					// 变量列表，可以有多个 var str1,str2 string
					Names: []*ast.Ident{
						{
							Name: vi.FieldName,
							Obj: &ast.Object{
								Kind: ast.Var,
								Name: vi.FieldName,
							},
						},
					},
					Type: &ast.SelectorExpr{
						X: &ast.Ident{
							Name: vi.FieldType,
						},
						Sel: &ast.Ident{
							Name: vi.FieldComment,
						},
					},
				}
				//放到最后
				t.Fields.List = append(t.Fields.List, f)
			}
		}
	}
	return vi
}

func (vi *StructMeta) RollBack(node ast.Node) ast.Visitor {
	if genDecl, ok := node.(*ast.GenDecl); ok {
		if genDecl.Tok == token.TYPE {
			//ast.Println(token.NewFileSet(),genDecl)
			for i := range genDecl.Specs {
				switch n := genDecl.Specs[i].(type) {
				case *ast.TypeSpec:
					// 找到结构体名称
					if n.Name.Name == vi.FindStructName {
						switch t := n.Type.(type) {
						case *ast.StructType:
							//数据已经存在，不重复添加
							for k, f := range t.Fields.List {
								if f.Names != nil {
									if f.Names[0].Name == vi.FieldName {
										//删除
										block := t.Fields
										block.List = append(append([]*ast.Field{}, block.List[:k]...), block.List[k+1:]...)
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

func (vi *StructMeta) ToInjectCode() string {
	vi.InjectCode = fmt.Sprintf("%s", vi.FieldType)

	if vi.FieldIdent != "" {
		vi.InjectCode = fmt.Sprintf("%s.%s", vi.FieldIdent, vi.InjectCode)
	}

	if vi.FieldName != "" {
		vi.InjectCode = fmt.Sprintf("%s %s", vi.FieldName, vi.InjectCode)
	}

	if vi.FieldComment != "" {
		vi.InjectCode = fmt.Sprintf("%s //%s", vi.InjectCode, vi.FieldComment)
	}

	return vi.InjectCode
}

func NewStructMete(findStructName, importCode string) {
	meta := &StructMeta{}
	input := fmt.Sprintf(`
package main

type %s struct {
	%s
}
`, findStructName, importCode)
	fParse, err := parser.ParseFile(token.NewFileSet(), "", input, parser.ParseComments)
	if err != nil {
		fmt.Println("Error parsing file:", err)
		return
	}
	ast.Print(token.NewFileSet(), fParse)

	ast.Inspect(fParse, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.TypeSpec:
			if x.Name.Name == findStructName {
				structType, ok := x.Type.(*ast.StructType)
				if !ok {
					return false
				}

				for _, field := range structType.Fields.List {
					fieldName := field.Names[0].Name
					// Extract field ident
					fieldIdent := ""
					fieldType := ""
					switch t := field.Type.(type) {
					case *ast.Ident:
						fieldIdent = t.Name
					case *ast.SelectorExpr:
						fieldIdent = t.X.(*ast.Ident).Name
						fieldType = t.Sel.Name
					}

					// Extract field comment
					fieldComment := ""
					if field.Comment != nil {
						re := regexp.MustCompile(`//\s*(.*)`)
						matches := re.FindStringSubmatch(field.Comment.List[0].Text)
						if len(matches) > 1 {
							fieldComment = strings.TrimSpace(matches[1])
						}
					}
					meta.FindStructName = findStructName
					meta.FieldName = fieldName
					meta.FieldIdent = fieldIdent
					meta.FieldType = fieldType
					meta.FieldComment = fieldComment
					log.Printf("Field Group: %s, Field Ident: %s, Field Type: %s, Field ColumnComment: %s\n", fieldName, fieldIdent, fieldType, fieldComment)
				}
			}
		}
		return true
	})

	log.Println("GetCode", meta.ToInjectCode())
	log.Println("NewStructMete", jsonconv.AnyToJsonIndent(meta))
}
