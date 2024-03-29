package apiparser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"regexp"
	"sort"
	"strings"
)

type AstParserConfig struct {
	ApiBase string
}

type AstParser struct {
	AstParserConfig
}

func NewAstParser(cfg AstParserConfig) ApiParser {
	return &AstParser{
		AstParserConfig: cfg,
	}
}

func (s *AstParser) ParseApiDocsByRoots(root ...string) (out []*ApiDeclare, err error) {
	out = make([]*ApiDeclare, 0)
	for _, v := range root {
		apis, err := s.ParseApiDocsByRoot(v)
		if err != nil {
			return nil, err
		}

		out = append(out, apis...)
	}

	//sort.Slice(out, func(i, j int) bool {
	//	if out[i].Tag == out[j].Tag {
	//		if out[i].Router == out[j].Router {
	//			return out[i].Method < out[j].Method
	//		}
	//		return out[i].Router < out[j].Router
	//	}
	//	return out[i].Tag < out[j].Tag
	//})
	return out, nil
}

func (s *AstParser) ParseApiDocsByRoot(root string) (out []*ApiDeclare, err error) {
	apiDocs := make([]*ApiDeclare, 0)
	// 遍历目录下的所有文件
	VisitFile(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 是目录，则跳过
		if info.IsDir() {
			return nil
		}
		// 是文件，则判断是否是ctl.go文件
		if strings.HasSuffix(path, "ctl.go") {
			// 解析文件
			apiDocs = append(apiDocs, s.ParseApiDoc(path)...)
		}

		return nil
	})

	return apiDocs, nil
}

func (s *AstParser) ParseApiDoc(filepath string) []*ApiDeclare {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filepath, nil, parser.ParseComments)
	if err != nil {
		fmt.Println("解析代码时发生错误:", err)
		return nil
	}

	apiDocs := make([]*ApiDeclare, 0)

	for _, decl := range file.Decls {
		if f, ok := decl.(*ast.FuncDecl); ok {
			name := getFunctionNameFromAst(f)
			doc := ParseFunctionDoc(f)

			api := &ApiDeclare{}

			api.FunctionName = name
			for _, comment := range doc {
				// 按照空白分割
				content := strings.Fields(comment.Content)
				if len(content) == 0 {
					continue
				}

				switch comment.Tag {
				case "Tags":
					// 按空白分割
					api.Tag = content[0]

				case "Summary":
					api.Summary = content[0]

				case "Param":
					tp := content[1]
					field := &ApiParam{
						Name:        content[0],
						Type:        content[2],
						Description: strings.Trim(content[4], `"`),
					}

					switch tp {
					case "header":
						api.Header = append(api.Header, field)
					case "path":
						api.Path = append(api.Path, field)
					case "query":
						api.Query = append(api.Query, field)
					case "formData":
						api.Form = append(api.Form, field)
					case "body":
						api.Body = &ApiParam{
							Name:        content[0],
							Type:        content[2],
							Description: strings.Trim(content[4], `"`),
						}
					}

				case "Router":
					api.Router = s.ApiBase + content[0]
					api.Method = strings.TrimSuffix(strings.TrimPrefix(content[1], "["), "]")

				case "Success":
					api.Response = &ApiParam{
						Name:        content[0],
						Type:        content[2],
						Description: strings.Trim(content[3], `"`),
					}
				}
			}

			if api.Tag == "" {
				continue
			}
			apiDocs = append(apiDocs, api)
			// fmt.Println("函数注释:", jsonconv.ObjectToJsonIndent(api))
		}
	}

	return apiDocs
}

func (s *AstParser) ParseModelDocsByRoots(root ...string) (out []*ModelDeclare, err error) {
	out = make([]*ModelDeclare, 0)
	for _, v := range root {
		models, err := s.ParseModelDocsByRoot(v)
		if err != nil {
			return nil, err
		}

		out = append(out, models...)
	}

	sort.Slice(out, func(i, j int) bool {
		if out[i].Type == out[j].Type {
			return out[i].Type < out[j].Type
		}
		return out[i].Type < out[j].Type
	})
	return out, nil
}

func (s *AstParser) ParseModelDocsByRoot(root string) (out []*ModelDeclare, err error) {
	var models []*ModelDeclare
	// 遍历目录下的所有文件
	VisitFile(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 是目录，则跳过
		if info.IsDir() {
			return nil
		}
		// 是文件，则判断是否是.go文件
		if strings.HasSuffix(path, ".go") {
			// 解析文件
			model := s.ParseModelDoc(path)
			models = append(models, model...)
		}

		return nil
	})

	return models, nil
}

func (s *AstParser) ParseModelDoc(filepath string) []*ModelDeclare {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filepath, nil, parser.ParseComments)
	if err != nil {
		fmt.Println("解析代码时发生错误:", err)
		return nil
	}

	var models []*ModelDeclare
	for _, decl := range file.Decls {
		if f, ok := decl.(*ast.GenDecl); ok {
			for _, spec := range f.Specs {
				if t, ok := spec.(*ast.TypeSpec); ok {
					if s, ok := t.Type.(*ast.StructType); ok {
						var modelFields []*ModelField
						var extendFields []*ModelField
						for _, field := range s.Fields.List {
							if len(field.Names) > 0 {
								// 属性
								modelFields = append(modelFields, getFieldsFormNode(field))
							} else {
								// 继承
								extend := getFieldsFormNode(field)
								//if !strings.Contains(extend.Type, ".") {
								//	extend.Type = fmt.Sprintf("%v.%v", file.Name.Name, extend.Type)
								//}
								//fmt.Println("getFieldsFormNode", extend.Name, extend.Type, extend.Comment)
								extendFields = append(extendFields, extend)
							}
						}

						modelName := fmt.Sprintf("%s.%s", file.Name.Name, t.Name.Name)
						model := &ModelDeclare{
							Type:   modelName,
							Extend: extendFields,
							Fields: modelFields,
						}

						models = append(models, model)
					}
				}
			}
		}
	}

	return models
}

// 获取函数名称
func getFunctionNameFromAst(f *ast.FuncDecl) string {
	return f.Name.Name
}

func getFieldsFormNode(field ast.Node) *ModelField {

	switch node := field.(type) {
	case *ast.ArrayType:
		return getFieldsFormNode(node.Elt)
	case *ast.Field:
		if len(node.Names) > 0 {
			name := node.Names[0].Name
			tp := node.Type
			tag := getJsonTagFromField(node)
			fmt.Println("tag:", tag)
			elem := &ModelField{
				Name:    name,
				JsonTag: tag,
				Type:    getTypeNameFormExpr(tp),
			}

			// 读取字段的普通注释
			if node.Doc != nil {
				elem.Comment = strings.TrimSpace(node.Doc.Text())
			}

			// 读取字段的行尾注释
			if node.Comment != nil {
				elem.Comment = strings.TrimSpace(node.Comment.Text())
			}

			return elem
		}
		return getFieldsFormNode(node.Type)
	case *ast.StarExpr:
		return getFieldsFormNode(node.X)
	case *ast.SelectorExpr:
		if xIdent, ok := node.X.(*ast.Ident); ok {
			elem := &ModelField{
				Name: "",
				Type: fmt.Sprintf("%s.%s", xIdent.Name, node.Sel.Name),
			}

			return elem
		}
	case *ast.Ident:
		elem := &ModelField{
			Name: "",
			Type: node.Name,
		}

		return elem
	case *ast.InterfaceType:
		return &ModelField{
			Name: "",
			Type: "any",
		}

	default:
		ast.Print(nil, field)
	}

	return nil
}

// 从 ast.Expr 中取出名称。
func getTypeNameFormExpr(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		// 处理标识符，直接返回标识符的名称
		return t.Name
	case *ast.SelectorExpr:
		// 处理选择器表达式，格式化输出包名和选择器名称
		return fmt.Sprintf("%v.%v", getTypeNameFormExpr(t.X), t.Sel.Name)
	case *ast.ArrayType:
		// 处理数组类型，格式化输出数组的元素类型名称
		return fmt.Sprintf("[]%v", getTypeNameFormExpr(t.Elt))
	case *ast.StarExpr:
		// 处理指针类型，格式化输出指针指向的类型名称
		return fmt.Sprintf("*%v", getTypeNameFormExpr(t.X))
	case *ast.InterfaceType:
		// 处理接口类型，直接返回 "interface{}"
		return "interface{}"
	default:
		// 其他未处理的表达式类型，返回空字符串
		return ""
	}
}

func getJsonTagFromField(field *ast.Field) string {
	if field.Tag == nil {
		return ""
	}

	var jsonTag string
	// 定义匹配 json tag 的正则表达式
	jsonRegex := regexp.MustCompile(`json:"([^"]*)"`)

	// 提取 json tag 中的信息
	jsonTagMatches := jsonRegex.FindStringSubmatch(field.Tag.Value)
	if len(jsonTagMatches) > 1 {
		jsonTag = jsonTagMatches[1]
	}

	jsonTag = strings.Split(jsonTag, ",")[0]
	jsonTag = strings.TrimSpace(jsonTag)
	if jsonTag == "-" {
		jsonTag = ""
	}

	return jsonTag
}
