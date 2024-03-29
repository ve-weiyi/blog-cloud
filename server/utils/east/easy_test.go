package east

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"regexp"
	"testing"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
}

func TestInject(t *testing.T) {
	inject := AstInjectionMeta{
		FilePath: "../ast/test/test.go",
		ImportMetas: []*ImportMeta{
			&ImportMeta{
				ImportAlias:   "jsoniter",
				ImportPackage: "github.com/json-iterator/go",
			},
			&ImportMeta{
				ImportAlias:   "",
				ImportPackage: "go/ast",
			},
		},
		StructMetas: []*StructMeta{
			&StructMeta{
				FindStructName: "ApiGroup",
				FieldName:      "visitor",
				FieldType:      "ast",
				FieldComment:   "Visitor",
			},
		},
		FuncMetas: []*FuncMeta{
			&FuncMeta{
				FuncName:   "GetName",
				FuncPos:    3,
				Variables:  []string{"value"},
				Symbol:     ":=",
				IdentNames: []string{"jsoniter", "AA", "ConfigCompatibleWithStandardLibrary"},
				//Parameters: []interface{}{"ss", 11},
			},
			//&FuncMeta{
			//	FuncName:   "GetName",
			//	FuncPos:    2,
			//	Variables:  []string{""},
			//	IdentNames: []string{"value", "UnmarshalFromString"},
			//	Parameters: []interface{}{"ss", 11},
			//},
		},
	}
	var err error
	//inject.Walk()
	//err = inject.RollBack()
	err = inject.Inject()
	log.Println("-->", err)
	if err != nil {
		return
	}
}
func TestNewAst(t *testing.T) {
	//NewImportMete(`jsoniter "github.com/json-iterator/go"`)
	//NewStructMete("ApiGroup", `Alias   ast.Visitor //元素别名`)
	NewFuncMete("main", `json := jsoniter.ConfigCompatibleWithStandardLibrary()`)
}

func TestParse(t *testing.T) {

	fSet := token.NewFileSet()
	fParser, err := parser.ParseFile(fSet, "../ast/test/test.go", nil, parser.ParseComments)
	if err != nil {
		return
	}
	log.Println("fun--")
	ast.Print(token.NewFileSet(), fParser)

	var output []byte
	buffer := bytes.NewBuffer(output)
	err = format.Node(buffer, fSet, fParser)
	if err != nil {
		log.Fatal(err)
	}
	// 写回数据
	os.WriteFile("../ast/test/test2.go", buffer.Bytes(), 0o600)
}

func TestValue(t *testing.T) {

	inputs := []string{
		`aa,bb :=        value.AA.BB.UnmarshalFromString(11.0     ,"22",anker)`,
		`ss,  zz     = value.AA.UnmarshalFromStringff(11,   ss22,true)`,
		`value.UnmarshalFromStringzz(11)`,
	}
	//1. `^([\w\s,]*)`：从字符串开始处匹配一个或多个字母、数字、下划线、空格或逗号。捕获到第一个分组。
	//2. `([:=]{1,2})`：匹配":="或者"="，捕获到第二个分组。
	//3. `\s*([\w.]+)`：匹配零个或多个空格，然后匹配一个或多个字母、数字、下划线或点。捕获到第三个分组。
	//4. `\(([^)]+)\)`：匹配一个左括号"("，然后匹配除右括号")"之外的一个或多个字符，最后匹配一个右括号")"。捕获到第四个分组。

	re := regexp.MustCompile(`^([\w\s,]*)?([:=]{1,2})?\s*([\w.]+)\(([^)]+)\)`)

	for _, input := range inputs {
		// 正则表达式匹配变量、调用方法和参数
		matches := re.FindStringSubmatch(input)
		if len(matches) == 0 {
			fmt.Println("No matches found")
		}
		// 提取变量
		variables := matches[1]
		fmt.Printf("Variables: %s\n", variables)
		// 提取符号
		symbol := matches[2]
		fmt.Printf("Symbol: %s\n", symbol)
		// 提取调用方法
		function := matches[3]
		fmt.Printf("Function: %s\n", function)
		// 提取参数
		args := matches[4]
		fmt.Printf("Arguments: %s\n", args)
	}

}
