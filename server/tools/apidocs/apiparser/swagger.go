package apiparser

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

type SwaggerParser struct {
}

func NewSwaggerParser() ApiParser {
	return &SwaggerParser{}
}

func (s *SwaggerParser) ReadSwagJSON(filepath string) (out *SwaggerDefinition, err error) {

	// 读取 JSON 文件内容
	jsonData, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read JSON file:%v", err)
	}

	// 解析 JSON 数据到结构体
	var swagger SwaggerDefinition
	err = json.Unmarshal(jsonData, &swagger)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON:%v", err)
	}

	//fmt.Println("swagger:", jsonconv.ObjectToJsonIndent(swagger))

	return &swagger, nil
}

func (s *SwaggerParser) ParseApiDocsByRoots(root ...string) (out []*ApiDeclare, err error) {
	out = make([]*ApiDeclare, 0)
	for _, v := range root {
		apis, err := s.ParseApiDocsByRoot(v)
		if err != nil {
			return nil, err
		}

		out = append(out, apis...)
	}

	sort.Slice(out, func(i, j int) bool {
		if out[i].Tag == out[j].Tag {
			if out[i].Router == out[j].Router {
				return out[i].Method < out[j].Method
			}
			return out[i].Router < out[j].Router
		}
		return out[i].Tag < out[j].Tag
	})
	return out, nil
}

func (s *SwaggerParser) ParseApiDocsByRoot(root string) (out []*ApiDeclare, err error) {
	out = make([]*ApiDeclare, 0)
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
		if strings.HasSuffix(path, ".json") {
			// 解析文件
			swagger, err := s.ReadSwagJSON(path)
			if err != nil {
				return err
			}

			out = append(out, s.ParseApiDoc(swagger)...)
		}
		return nil
	})

	return out, nil
}

func (s *SwaggerParser) ParseApiDoc(swagger *SwaggerDefinition) []*ApiDeclare {
	// 生成 TypeScript 代码
	apis := make([]*ApiDeclare, 0)
	for path, paths := range swagger.Paths {

		//fmt.Println("path:", path)

		for method, apiMethod := range paths {

			//fmt.Println("method:", method, "apiMethod:", jsonconv.ObjectToJsonIndent(apiMethod))

			tag := JoinArray(apiMethod.Tags)

			// 未创建tag
			apiPath := swagger.BasePath + path

			var header []*ApiParam
			var query []*ApiParam
			var params []*ApiParam
			var form []*ApiParam
			var body *ApiParam
			for _, param := range apiMethod.Parameters {

				//fmt.Println("param:", jsonconv.ObjectToJsonIndent(param))
				//t := getTypeNameFormSchema(param.SchemaObject)

				if param.In == "body" {
					body = &ApiParam{
						Name:        param.Name,
						Type:        getTypeNameFormSchema(param.Schema),
						Description: param.Description,
					}

				} else {
					p := &ApiParam{
						Name:        param.Name,
						Type:        param.Type,
						Description: param.Description,
					}
					switch param.In {
					case "header":
						header = append(header, p)
					case "query":
						query = append(query, p)
					case "path":
						params = append(params, p)
					case "form":
						form = append(form, p)
					default:
					}
				}
			}

			var response *ApiParam
			for k, v := range apiMethod.Responses {
				// 200
				//fmt.Println("in:", jsonconv.ObjectToJsonIndent(v.Schema))
				response = &ApiParam{
					Name:        k,
					Type:        getTypeNameFormSchema(v.Schema),
					Description: v.Description,
				}
				//fmt.Println("out:", jsonconv.ObjectToJsonIndent(response))
				break
			}

			apiTs := &ApiDeclare{
				Tag:          tag,
				FunctionName: getFunctionNameFormPath(path),
				Summary:      apiMethod.Summary,
				Router:       apiPath,
				Method:       method,
				Header:       header,
				Path:         params,
				Query:        query,
				Form:         form,
				Body:         body,
				Response:     response,
			}

			apis = append(apis, apiTs)
		}

	}

	return apis
}

func (s *SwaggerParser) ParseModelDocsByRoots(root ...string) (out []*ModelDeclare, err error) {
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

func (s *SwaggerParser) ParseModelDocsByRoot(root string) (out []*ModelDeclare, err error) {
	out = make([]*ModelDeclare, 0)
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
		if strings.HasSuffix(path, ".json") {
			// 解析文件
			swagger, err := s.ReadSwagJSON(path)
			if err != nil {
				return err
			}

			out = append(out, s.ParseModelDoc(swagger)...)
		}
		return nil
	})

	return out, nil
}

func (s *SwaggerParser) ParseModelDoc(swagger *SwaggerDefinition) []*ModelDeclare {
	out := make([]*ModelDeclare, 0)

	for k, v := range swagger.Definitions {
		//var pkg string
		//if strings.Contains(k, ".") {
		//	pkg = strings.Split(k, ".")[0]
		//}

		model := &ModelDeclare{
			//Pkg:    pkg,
			Type:   k,
			Extend: nil,
			Fields: getFieldsFormSchema(v.Properties),
		}

		out = append(out, model)
	}

	return out
}

func getFunctionNameFormPath(path string) string {
	if path == "/" || path == "" {
		return ""
	}

	var name string
	name = strings.ReplaceAll(path, ":", "_")
	name = strings.ReplaceAll(path, "{", "_")
	name = strings.ReplaceAll(path, "}", "_")
	name = strings.ReplaceAll(path, "/", "_")

	var newName string
	var flag bool = false
	for i := 0; i < len(name); i++ {
		if name[i] == '_' {
			flag = true
			continue
		}

		if flag {
			newName = newName + strings.ToUpper(string(byte(rune(name[i]))))
			flag = false
		} else {
			newName = newName + string(byte(rune(name[i])))
		}
	}

	return strings.ToLower(string(byte(rune(newName[0])))) + newName[1:]
}

func getFieldsFormSchema(in map[string]SchemaObject) []*ModelField {
	if in == nil {
		return nil
	}

	out := make([]*ModelField, 0)
	for k, v := range in {

		var comment string
		if v.Items != nil {
			comment = v.Items.Description
		}

		field := &ModelField{
			Name:    k,
			Type:    getTypeNameFormSchema(&v),
			Comment: comment,
		}
		out = append(out, field)
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].Name < out[j].Name
	})
	return out
}

func getTypeNameFormSchema(in *SchemaObject) string {
	if in == nil {
		return ""
	}

	//out := &ApiObject{
	//Properties: make(map[string]*ApiObject),
	//}
	var out string
	if in.Items != nil {
		//t += getTypeNameFormItems(in.Items)

		switch in.Type {
		case "array":
			out += "[]"
			out += strings.TrimPrefix(in.Items.Items.Ref, "#/definitions/")
			out += strings.TrimPrefix(in.Items.Items.Type, "#/definitions/")
		case "string":
			out += "string"
		case "integer":
			out += "int"
		case "":
			if in.Ref != "" {
				//out.Type = "object"
				out += strings.TrimPrefix(in.Ref, "#/definitions/")
			}
		default:
			out += "any"
		}

	}

	if in.AllOf != nil {
		for _, v := range in.AllOf {
			// 引用类型
			if v.Ref != "" {
				out += strings.TrimPrefix(v.Ref, "#/definitions/")
			}
			// 类型
			//if v.Type != "" {
			//	out += v.Type
			//}

			for k, p := range v.Properties {
				//t += getTypeNameFormSchema(&p)
				out += "{"
				out += k + "=" + getTypeNameFormSchema(&p)
				out += "}"
			}
		}
	}

	return out
}
