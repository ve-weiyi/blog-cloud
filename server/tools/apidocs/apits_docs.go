package apidocs

import (
	"fmt"
	"path"
	"regexp"
	"sort"
	"strings"

	"github.com/ve-weiyi/ve-blog-golang/server/tools/apidocs/apiparser"
	"github.com/ve-weiyi/ve-blog-golang/server/tools/quickstart/invent"
	"github.com/ve-weiyi/ve-blog-golang/server/utils/jsonconv"
)

type Config struct {
	OutRoot   string
	ApiRoot   []string
	ModelRoot []string

	ApiBase        string
	ImportPkgPaths []string
	IgnoredModels  map[string]string
	ReplaceModels  map[string]string

	ApiFuncNameAs  func(api *apiparser.ApiDeclare) string
	ApiFieldNameAs func(model *apiparser.ModelField) string
	ApiFieldTypeAs func(name string) string
}

type AstApiDoc struct {
	Config

	Parser apiparser.ApiParser

	ApiDeclares  []*apiparser.ApiDeclare
	TypeDeclares []*apiparser.ModelDeclare
}

func NewAstApiDoc(config Config) *AstApiDoc {
	cfg := apiparser.AstParserConfig{
		ApiBase: "/api/v1",
	}
	return &AstApiDoc{
		Parser: apiparser.NewAstParser(cfg),
		Config: config,
	}
}

func (s *AstApiDoc) Parse() (err error) {
	// 解析api定义
	var apis []*apiparser.ApiDeclare
	apis, err = s.Parser.ParseApiDocsByRoots(s.ApiRoot...)
	if err != nil {
		return fmt.Errorf("解析api定义时发生错误:%v", err)
	}

	// 解析model定义
	var models []*apiparser.ModelDeclare
	models, err = s.Parser.ParseModelDocsByRoots(s.ModelRoot...)
	if err != nil {
		return fmt.Errorf("解析model定义时发生错误:%v", err)
	}

	// 根据tag对api分类
	s.ApiDeclares = apis

	for _, m := range models {
		// 过滤需要忽略的model
		_, ok := s.IgnoredModels[m.Type]
		if ok {
			continue
		}
		s.TypeDeclares = append(s.TypeDeclares, m)
	}

	//fmt.Println("ApiDeclares:", jsonconv.ObjectToJsonIndent(apis))
	//fmt.Println("TypeDeclares:", jsonconv.ObjectToJsonIndent(models))
	return nil
}

// 生成 TypeScript
func (s *AstApiDoc) GenerateTsTypeFile() {
	var tsModelDeclares []*TsModelDeclare
	for _, model := range s.TypeDeclares {
		// 过滤需要忽略的model
		_, ok := s.IgnoredModels[model.Type]
		if ok {
			continue
		}

		item := s.convertTsModelDeclare(model)
		if item != nil {
			tsModelDeclares = append(tsModelDeclares, item)
		}
	}

	//fmt.Println("tsModelDeclares:", jsonconv.ObjectToJsonIndent(tsModelDeclares))
	meta := invent.TemplateMeta{
		Key:            "",
		Mode:           invent.ModeCreateOrReplace,
		CodeOutPath:    path.Join(s.OutRoot, "types.ts"),
		TemplateString: ModelTypeScript,
		FunMap:         map[string]any{"joinArray": apiparser.JoinArray},
		Data:           tsModelDeclares,
	}

	err := meta.CreateTempFile()
	if err != nil {
		fmt.Println("生成 TypeScript 时发生错误:", err)
	}
}

func (s *AstApiDoc) GenerateTsApiFiles() {
	// 根据tag进行分组
	var apiGroups = map[string][]*apiparser.ApiDeclare{}
	for _, api := range s.ApiDeclares {
		apiGroups[api.Tag] = append(apiGroups[api.Tag], api)
	}

	tsApiDocs := s.convertTsApiDocs(apiGroups)
	//fmt.Println("tsApiDocs:", jsonconv.ObjectToJsonIndent(tsApiDocs))
	var metas []invent.TemplateMeta
	for _, apiDoc := range tsApiDocs {
		meta := invent.TemplateMeta{
			Key:            "",
			Mode:           invent.ModeCreateOrReplace,
			CodeOutPath:    path.Join(s.OutRoot, fmt.Sprintf("%s.ts", jsonconv.Camel2Case(apiDoc.Tag))),
			TemplateString: ApiTypeScript,
			FunMap:         map[string]any{"joinArray": joinArray},
			Data:           apiDoc,
		}
		//fmt.Println("apiDocs:", jsonconv.ObjectToJsonIndent(apiDoc))
		metas = append(metas, meta)
	}

	for _, meta := range metas {
		err := meta.CreateTempFile()
		if err != nil {
			fmt.Println("生成 TypeScript 时发生错误:", err)
		}
		fmt.Println("TypeScript 文件已生成：", meta.CodeOutPath)
	}
}

// 转换为分组为ts api文件
func (s *AstApiDoc) convertTsApiDocs(groups map[string][]*apiparser.ApiDeclare) []*TsApiDoc {

	var apiDocs []*TsApiDoc

	for tag, group := range groups {
		var params []*apiparser.ApiParam
		var tsApiDeclare []*TsApiDeclare
		for _, api := range group {
			tsApiDeclare = append(tsApiDeclare, s.convertTsApiDeclare(api))
			params = append(params, collectRequestParams(api)...)
			params = append(params, collectResponseParams(api)...)
		}

		// 移除重复元素
		//params = removeDuplicates(params)
		//fmt.Println("params:", params)

		tsModelsMap := make(map[string]*TsModelDeclare)
		for _, param := range params {

			//fmt.Println("tsModels:", param, jsonconv.ObjectToJsonIndent(s.findTsModelDeclareByName(param)))
			model := s.findModelDeclare(param.Type)
			tsModel := s.convertTsModelDeclare(model)
			if tsModel == nil {
				continue
			}

			tsModelsMap[tsModel.Name] = tsModel
		}

		var tsModels []*TsModelDeclare
		var tsModelsName []string

		for _, item := range tsModelsMap {
			tsModels = append(tsModels, item)
			tsModelsName = append(tsModelsName, item.Name)
		}

		sort.Slice(tsModelsName, func(i, j int) bool {
			return tsModelsName[i] < tsModelsName[j]
		})

		var apiDoc *TsApiDoc
		apiDoc = &TsApiDoc{
			Tag:              tag,
			ImportPkgPaths:   s.ImportPkgPaths,
			ImportModelTypes: tsModelsName,
			ModelDeclares:    tsModels,
			ApiDeclares:      tsApiDeclare,
		}

		apiDocs = append(apiDocs, apiDoc)
	}
	return apiDocs
}

func (s *AstApiDoc) convertTsModelDeclare(model *apiparser.ModelDeclare) *TsModelDeclare {
	if model == nil {
		return nil
	}

	tsFields := make([]*TsModelField, 0)
	tsExtends := make([]string, 0)

	// 属性
	for _, field := range model.Fields {
		tsField := &TsModelField{
			Name:    s.ApiFieldNameAs(field),
			Type:    s.ApiFieldTypeAs(field.Type),
			Comment: field.Comment,
		}

		tsFields = append(tsFields, tsField)
	}
	// 继承
	for _, extend := range model.Extend {
		tsExtends = append(tsExtends, s.ApiFieldTypeAs(extend.Type))
	}

	tsModel := &TsModelDeclare{
		Name:    s.ApiFieldTypeAs(model.Type),
		Extends: tsExtends,
		Fields:  tsFields,
	}

	return tsModel
}

func (s *AstApiDoc) convertTsApiDeclare(doc *apiparser.ApiDeclare) *TsApiDeclare {

	reqParams := collectRequestParams(doc)
	var reqStr string
	for i, param := range reqParams {
		if i > 0 {
			reqStr += ", "
		}
		ts := s.ApiFieldTypeAs(param.Type)
		// 需要替换的参数
		v, ok := s.ReplaceModels[ts]
		//fmt.Println("ReplaceModels", ts, v, ok)
		if ok {
			ts = v
		}
		reqStr += fmt.Sprintf("%s: %s", param.Name, ts)
	}

	// 提取参数
	respParams := collectResponseParams(doc)
	var respStr string
	for i, param := range respParams {
		ts := s.ApiFieldTypeAs(param.Type)

		// 需要替换的参数
		v, ok := s.ReplaceModels[ts]
		//fmt.Println("ReplaceModels", ts, v, ok)
		if ok {
			ts = v
		}

		if i > 0 {
			respStr += "<"
		}
		respStr += ts
	}
	for i := 0; i < len(respParams)-1; i++ {
		respStr += ">"
	}
	if respStr == "" {
		respStr = "any"
	}

	// {id}->${id}
	re := regexp.MustCompile(`\{(.+?)\}`)
	var tsDoc = &TsApiDeclare{
		Tag:          doc.Tag,
		FunctionName: s.ApiFuncNameAs(doc),
		Summary:      doc.Summary,
		Base:         s.ApiBase,
		Route:        re.ReplaceAllString(doc.Router, "${$1}"),
		Method:       doc.Method,
		Header:       convertTsParams(doc.Header),
		Path:         convertTsParams(doc.Path),
		Query:        convertTsParams(doc.Query),
		Form:         convertTsParams(doc.Form),
		Body:         convertTsParam(doc.Body),
		Request:      reqStr,
		Response:     respStr,
	}

	return tsDoc
}

// 查找接口model定义
func (s *AstApiDoc) findModelDeclare(name string) *apiparser.ModelDeclare {
	name = strings.Trim(name, "[]")

	for _, model := range s.TypeDeclares {
		if model.Type == name {
			return model
		}
		// package name 都相等的情况
		//if model.Pkg != "" {
		//	if fmt.Sprintf("%v.%v", model.Pkg, name) == model.Type {
		//		return model
		//	}
		//}
	}

	return nil
}

func convertTsParams(list []*apiparser.ApiParam) []*TsApiParam {
	if list == nil {
		return nil
	}
	var out []*TsApiParam
	for _, in := range list {
		out = append(out, convertTsParam(in))
	}

	return out
}

func convertTsParam(in *apiparser.ApiParam) *TsApiParam {
	if in == nil {
		return nil
	}

	out := &TsApiParam{
		Name: in.Name,
		Type: convertGoTypeToTsType(in.Type),
	}

	return out
}

func collectRequestParams(api *apiparser.ApiDeclare) []*apiparser.ApiParam {
	params := make([]*apiparser.ApiParam, 0)
	//if api.Header != nil {
	//	for _, param := range api.Header {
	//		params = append(params, param)
	//	}
	//}
	if api.Path != nil {
		for _, param := range api.Path {
			params = append(params, param)
		}
	}
	if api.Query != nil {
		for _, param := range api.Query {
			params = append(params, param)
		}
	}
	if api.Form != nil {
		for _, param := range api.Form {
			params = append(params, param)
		}
	}
	if api.Body != nil {
		params = append(params, api.Body)
	}

	return params
}

// response.Response{data=response.PageResult{list=[]entity.User}}-->Response、PageResult、[]entity.User
func collectResponseParams(api *apiparser.ApiDeclare) []*apiparser.ApiParam {
	if api.Response == nil {
		return nil
	}
	return extractResponseParams(api.Response.Type)
}

func extractResponseParams(response string) (fields []*apiparser.ApiParam) {
	left := strings.Index(response, "{")
	right := strings.LastIndex(response, "}")

	// data=entity.Api
	if left == -1 || right == -1 {
		flag := strings.Index(response, "=")
		return []*apiparser.ApiParam{{
			Name: response[:flag],
			Type: response[flag+1:],
		}}
	}

	// response.Response{data=entity.Api}
	if left != 1 {
		var name string
		flag := strings.Index(response, "=")
		if flag < left {
			name = response[:flag]
		} else {
			flag = -1
		}

		f := &apiparser.ApiParam{
			Name: name,
			Type: response[flag+1 : left],
		}
		fields = append(fields, f)
	}

	other := response[left+1 : right]
	return append(fields, extractResponseParams(other)...)
}
