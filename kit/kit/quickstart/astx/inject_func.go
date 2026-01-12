package astx

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
)

const (
	name = token.DEFINE
)

// Variables := IdentName.SelName(Parameters...)
// IdentName.SelName(Parameters...)
type FuncMeta struct {
	InjectCode string //注入的代码
	FuncName   string //插入的方法
	FuncPos    int    //插入的位置

	//LowerStartCamelName string //声明
	//IdentName string //引用名
	//SelName   string //元素struct名

	Variables  []string
	Symbol     string
	IdentNames []string
	Parameters []interface{}

	fset *token.FileSet
	node ast.Stmt
}

func (vi *FuncMeta) Visit(node ast.Node) ast.Visitor {
	switch fn := node.(type) {
	case *ast.FuncDecl:
		log.Printf("Func: %s %s\n", fn.Name.Name, vi.FuncName)
		if fn.Name.Name == vi.FuncName {

			hasVar := false
			for _, v := range fn.Body.List {
				switch varSpec := v.(type) {
				case *ast.AssignStmt: //赋值表达式
					for i := range varSpec.Lhs {
						switch nn := varSpec.Lhs[i].(type) {
						case *ast.Ident:
							if nn.Name == vi.Variables[i] {
								hasVar = true
							}
						}
					}

				case *ast.ExprStmt: //表达语句
					switch expr := varSpec.X.(type) {
					case *ast.CallExpr: //是调用语句
						if call, ok := expr.Fun.(*ast.SelectorExpr); ok {
							x, _ := call.X.(*ast.Ident)
							// 包名和方法名都相等
							if x.Name == vi.IdentNames[0] {
								hasVar = true
							}
						}
					}
				case *ast.IfStmt:
				case *ast.GoStmt:
				case *ast.ReturnStmt:
					break
				}
			}

			if hasVar {
				break
			}

			var insertStmt ast.Stmt
			insertStmt = vi.GetNode().(ast.Stmt)

			//放到第pos位置
			index := vi.FuncPos
			if vi.FuncPos >= len(fn.Body.List)-1 {
				index = len(fn.Body.List) - 1
			}
			// Find the line number after the third line
			for i, stmt := range fn.Body.List {
				ast.Print(vi.fset, fn)
				pos := index
				log.Println("line ", stmt)
				log.Println("line ", vi.fset.Position(stmt.Pos()).Line)
				log.Println("line ", pos)
				log.Println("line ", vi.fset.Position(fn.Body.Lbrace).Line)
				log.Println("line ", vi.fset.Position(fn.Body.Rbrace).Line)
				if vi.fset.Position(stmt.Pos()).Line > pos {
					fn.Body.List = insertStatements(fn.Body.List, i, insertStmt)
					break
				}
			}
			////放到第pos位置
			//fn.Body.List = append(fn.Body.List, fn.Body.List[index])
			//copy(fn.Body.List[index+1:], fn.Body.List[index:])
			//fn.Body.List[index] = assignStmt
		}
	}
	return vi
}

func (vi *FuncMeta) RollBack(node ast.Node) ast.Visitor {
	if funDecl, ok := node.(*ast.FuncDecl); ok {
		if funDecl.Name.Name == vi.FuncName {
			block := funDecl.Body
			vi.rollBack(block)
		}
	}
	return vi
}

func (vi *FuncMeta) GetNode() ast.Node {
	if vi.node != nil {
		return vi.node
	}

	return nil
}

func (vi *FuncMeta) GetCode() string {
	injectCode := fmt.Sprintf("import")

	return injectCode
}

func (vi *FuncMeta) rollBack(node ast.Node) ast.Visitor {
	if funDecl, ok := node.(*ast.FuncDecl); ok {
		if funDecl.Name.Name == vi.FuncName {
			for k, v := range funDecl.Body.List {
				switch varSpec := v.(type) {
				//赋值表达式
				case *ast.AssignStmt:
					for i := range varSpec.Lhs {
						switch nn := varSpec.Lhs[i].(type) {
						case *ast.Ident:
							if nn.Name == vi.IdentNames[0] {
								//删除
								block := funDecl.Body
								block.List = append(append([]ast.Stmt{}, block.List[:k]...), block.List[k+1:]...)
								return vi
							}
						}
					}
				case *ast.ExprStmt: //表达语句
					switch expr := varSpec.X.(type) {
					case *ast.CallExpr: //是调用语句
						if call, ok := expr.Fun.(*ast.SelectorExpr); ok {
							x, _ := call.X.(*ast.Ident)
							// 包名和方法名都相等
							if x.Name == vi.IdentNames[0] {
								//删除
								block := funDecl.Body
								block.List = append(append([]ast.Stmt{}, block.List[:k]...), block.List[k+1:]...)
								return vi
							}
						}
					}
				}
			}
		}
	}
	return vi
}

func NewFuncMete(findFuncName, importCode string) *FuncMeta {
	meta := &FuncMeta{}
	input := fmt.Sprintf(`
package main
func %v(){
   %v 
}`, findFuncName, importCode)
	fParse, err := parser.ParseFile(token.NewFileSet(), "", input, 0)
	if err != nil {
		fmt.Println("Error parsing file:", err, importCode)
	}
	ast.Inspect(fParse, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			if x.Name.Name != findFuncName {
				break
			}
			//dst.Print(x)
			meta.FuncName = findFuncName
			for _, field := range x.Body.List {
				meta.node = field
				switch s := n.(type) {
				case *ast.AssignStmt:
					ast.Print(token.NewFileSet(), s)

					variables := []string{}
					for _, lhs := range s.Lhs {
						variables = append(variables, lhs.(*ast.Ident).Name)
					}

					indentNames := []string{}
					parameters := []interface{}{}
					switch callExpr := s.Rhs[0].(type) {
					case *ast.CallExpr:
						idents := extractIdents(callExpr.Fun)
						for _, ident := range idents {
							indentNames = append(indentNames, ident.Name)
						}
						for _, arg := range callExpr.Args {
							switch arg.(type) {
							case *ast.BasicLit:
								value, _ := InferType(arg.(*ast.BasicLit).Value)
								parameters = append(parameters, value)
							case *ast.Ident:
								parameters = append(parameters, arg.(*ast.Ident).Name)
							}
						}
					case *ast.SelectorExpr:
						idents := extractIdents(callExpr)
						for _, ident := range idents {
							indentNames = append(indentNames, ident.Name)
						}
					}

					meta.Variables = variables
					meta.Symbol = s.Tok.String()
					meta.IdentNames = indentNames
					meta.Parameters = parameters

				case *ast.ExprStmt:
					ast.Print(token.NewFileSet(), s)

					indentNames := []string{}
					callExpr := s.X.(*ast.CallExpr)
					idents := extractIdents(callExpr.Fun)
					for _, ident := range idents {
						indentNames = append(indentNames, ident.Name)
					}

					parameters := []interface{}{}
					for _, arg := range callExpr.Args {
						switch arg.(type) {
						case *ast.BasicLit:
							value, _ := InferType(arg.(*ast.BasicLit).Value)
							parameters = append(parameters, value)
						case *ast.Ident:
							parameters = append(parameters, arg.(*ast.Ident).Name)
						}
					}

					meta.IdentNames = indentNames
					meta.Parameters = parameters
				}
			}
		}
		return true
	})

	//log.Println("GetCode", meta.GetCode())
	//log.Println("NewFuncMete", jsonconv.AnyToJsonIndent(meta))
	return meta
}
