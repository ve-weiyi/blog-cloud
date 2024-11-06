package inject

import (
	"fmt"
	"go/token"
	"log"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
)

const (
	name = token.DEFINE
)

// Variables := IdentName.SelName(Parameters...)
// IdentName.SelName(Parameters...)
type FuncMeta struct {
	FuncName string //插入的方法
	FuncPos  int    //插入的位置

	Variables  []string      //左边的值
	Symbol     string        //符号
	IdentNames []string      //右边的值，当为return时，是右边变量名
	Parameters []interface{} //调用参数
	Comments   string        //注释

	fset *token.FileSet
	node dst.Stmt
}

func (vi *FuncMeta) Visit(node dst.Node) dst.Visitor {
	switch fn := node.(type) {
	case *dst.FuncDecl:
		if fn.Name.Name == vi.FuncName {
			//log.Printf("Func: %s %s\n", fn.Group.Group, vi.FuncName)
			for _, v := range fn.Body.List {
				switch stmt := v.(type) {
				case *dst.AssignStmt: //赋值表达式
					if vi.Symbol != ":=" {
						break
					}
					if len(vi.Variables) == len(stmt.Lhs) {
						for i, lhs := range stmt.Lhs {
							switch nn := lhs.(type) {
							case *dst.Ident:
								if nn.Name == vi.Variables[i] {
									log.Println("hasVar", nn.Name, vi.Variables[i])
									return vi
								}
							}
						}
					}

				case *dst.ExprStmt: //表达语句
					if vi.Symbol != "" {
						break
					}
					switch expr := stmt.X.(type) {
					case *dst.CallExpr: //是调用语句，调用语句不查重
						if call, ok := expr.Fun.(*dst.SelectorExpr); ok {
							_, _ = call.X.(*dst.Ident)
							// 包名xd和方法名都相等
						}
					}
				case *dst.BlockStmt: //空块，注释的效果等于方法名
					if vi.Symbol != "{}" {
						break
					}
					if len(stmt.Decs.Start.All()) != 0 {

					}
				case *dst.IfStmt:
				case *dst.GoStmt:
				case *dst.ReturnStmt:
					if vi.Symbol != "return" {
						break
					}
					break
				}
			}

			insertStmt := vi.GetNode()

			//放到第pos位置
			index := vi.FuncPos
			if vi.FuncPos >= len(fn.Body.List)-1 {
				//除了return语句
				index = len(fn.Body.List) - 1
			}
			for pos, stmt := range fn.Body.List {
				//log.Println("line ", pos, index, vi.IdentNames)
				if pos >= index {
					//dst.Print(stmt)
					switch stmt.(type) {
					case *dst.BlockStmt:
						if vi.Symbol != "{}" {
							break
						}
						//log.Println("insert BlockStmt", vi.IdentNames)
						stmt.(*dst.BlockStmt).List = insertStatements(stmt.(*dst.BlockStmt).List, pos, insertStmt.(dst.Stmt))
					case *dst.ReturnStmt:
						if vi.Symbol != token.RETURN.String() {
							break
						}
						//log.Println("insert ReturnStmt", vi.IdentNames)
						results := stmt.(*dst.ReturnStmt).Results
						if len(results) != len(vi.IdentNames) {
							break
						}
						for i, ret := range results {
							_, ok := ret.(*dst.UnaryExpr)
							if !ok {
								continue
							}
							name1 := ret.(*dst.UnaryExpr).X.(*dst.CompositeLit).Type.(*dst.Ident).Name
							//name2 := insertStmt.(*dst.ReturnStmt).Results[0].(*dst.UnaryExpr).X.(*dst.CompositeLit).Type.(*dst.Ident).Group
							if name1 == vi.IdentNames[i] {
								kv := ret.(*dst.UnaryExpr).X.(*dst.CompositeLit).Elts
								addKv := insertStmt.(*dst.ReturnStmt).Results[0].(*dst.UnaryExpr).X.(*dst.CompositeLit).Elts
								//未导入的才添加
								for _, add := range addKv {
									has := false
									for _, exist := range kv {
										if add.(*dst.KeyValueExpr).Key.(*dst.Ident).Name == exist.(*dst.KeyValueExpr).Key.(*dst.Ident).Name {
											has = true
											break
										}
									}
									if !has {
										kv = append(kv, add)
									}
								}
								ret.(*dst.UnaryExpr).X.(*dst.CompositeLit).Elts = kv
							}
						}
					case *dst.AssignStmt:
						if vi.Symbol != token.DEFINE.String() {
							break
						}
						//log.Println("insert AssignStmt", vi.IdentNames)
						fn.Body.List = insertStatements(fn.Body.List, pos, insertStmt.(dst.Stmt))
					default:
						//log.Println("insert default", vi.IdentNames)
					}
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

func (vi *FuncMeta) RollBack(node dst.Node) dst.Visitor {
	if funDecl, ok := node.(*dst.FuncDecl); ok {
		if funDecl.Name.Name == vi.FuncName {
			block := funDecl.Body
			vi.rollBack(block)
		}
	}
	return vi
}

func (vi *FuncMeta) GetNode() dst.Node {
	if vi.node != nil {
		return vi.node
	}

	return nil
}

func (vi *FuncMeta) GetCode() string {
	injectCode := fmt.Sprintf("import")

	return injectCode
}

func NewFuncMete(findFuncName, importCode string) *FuncMeta {
	meta := &FuncMeta{}
	input := fmt.Sprintf(`
package main
func %v(){
   %v 
}`, findFuncName, importCode)
	fParse, err := decorator.ParseFile(token.NewFileSet(), "", input, 0)
	if err != nil {
		fmt.Println("Error parsing file:", err)
	}
	dst.Inspect(fParse, func(n dst.Node) bool {
		switch x := n.(type) {
		case *dst.FuncDecl:
			if x.Name.Name != findFuncName {
				break
			}
			//dst.Print(x)
			meta.FuncName = findFuncName
			for _, field := range x.Body.List {
				meta.node = field
				switch s := field.(type) {
				case *dst.AssignStmt:

					variables := []string{}
					for _, lhs := range s.Lhs {
						variables = append(variables, lhs.(*dst.Ident).Name)
					}

					indentNames := []string{}
					parameters := []interface{}{}
					switch callExpr := s.Rhs[0].(type) {
					case *dst.CallExpr:
						idents := ExtractIdents(callExpr.Fun)
						for _, ident := range idents {
							indentNames = append(indentNames, ident.Name)
						}
						for _, arg := range callExpr.Args {
							switch arg.(type) {
							case *dst.BasicLit:
								value, _ := inferType(arg.(*dst.BasicLit).Value)
								parameters = append(parameters, value)
							case *dst.Ident:
								parameters = append(parameters, arg.(*dst.Ident).Name)
							}
						}
					case *dst.SelectorExpr:
						idents := ExtractIdents(callExpr)
						for _, ident := range idents {
							indentNames = append(indentNames, ident.Name)
						}
					}

					meta.Variables = variables
					meta.Symbol = s.Tok.String()
					meta.IdentNames = indentNames
					meta.Parameters = parameters

				case *dst.ExprStmt:

					indentNames := []string{}
					callExpr := s.X.(*dst.CallExpr)
					idents := ExtractIdents(callExpr.Fun)
					for _, ident := range idents {
						indentNames = append(indentNames, ident.Name)
					}

					parameters := []interface{}{}
					for _, arg := range callExpr.Args {
						switch arg.(type) {
						case *dst.BasicLit:
							value, _ := inferType(arg.(*dst.BasicLit).Value)
							parameters = append(parameters, value)
						case *dst.Ident:
							parameters = append(parameters, arg.(*dst.Ident).Name)
						}
					}

					meta.IdentNames = indentNames
					meta.Parameters = parameters

				case *dst.ReturnStmt:
					meta.Symbol = token.RETURN.String()
					//op := s.Results[0].(*dst.UnaryExpr).Op.String()
					ret := s.Results[0].(*dst.UnaryExpr).X.(*dst.CompositeLit)
					ident := ret.Type.(*dst.Ident).Name

					meta.IdentNames = []string{ident}
				}
			}
		}

		return true
	})

	//log.Println("GetCode", meta.GetCode())
	//log.Println("NewFuncMete", jsonconv.AnyToJsonIndent(meta))
	return meta
}

func (vi *FuncMeta) rollBack(block *dst.BlockStmt) dst.Visitor {
	for k, v := range block.List {
		switch stmt := v.(type) {
		//赋值表达式
		case *dst.AssignStmt:
			if vi.Symbol != ":=" {
				break
			}
			if len(vi.Variables) == len(stmt.Lhs) {
				for i, lhs := range stmt.Lhs {
					switch nn := lhs.(type) {
					case *dst.Ident:
						if nn.Name == vi.Variables[i] {
							//删除
							block.List = append(append([]dst.Stmt{}, block.List[:k]...), block.List[k+1:]...)
							return vi
						}
					}
				}
			}
		case *dst.ExprStmt: //表达语句
			if vi.Symbol != "" {
				break
			}
			switch expr := stmt.X.(type) {
			case *dst.CallExpr: //是调用语句
				if len(vi.IdentNames) == 0 {
					break
				}
				if call, ok := expr.Fun.(*dst.SelectorExpr); ok {
					x, _ := call.X.(*dst.Ident)
					// 包名和方法名都相等
					if x.Name == vi.IdentNames[0] {
						//删除
						block.List = append(append([]dst.Stmt{}, block.List[:k]...), block.List[k+1:]...)
						return vi
					}
				}
			}
		case *dst.BlockStmt: //{}
			if vi.Symbol != "{}" {
				break
			}
			vi.rollBack(stmt)
			if len(stmt.List) == 0 {
				//删除空块
				block.List = append(append([]dst.Stmt{}, block.List[:k]...), block.List[k+1:]...)
				return vi
			}
		case *dst.ReturnStmt:
			if vi.Symbol != "return" {
				break
			}
			results := stmt.Results
			insertStmt := vi.GetNode()
			if len(results) != len(vi.IdentNames) {
				break
			}

			for i, ret := range results {
				_, ok := ret.(*dst.UnaryExpr)
				if !ok {
					continue
				}
				name1 := ret.(*dst.UnaryExpr).X.(*dst.CompositeLit).Type.(*dst.Ident).Name
				//找到需要加入的变量名称
				if name1 == vi.IdentNames[i] {
					kv := ret.(*dst.UnaryExpr).X.(*dst.CompositeLit).Elts
					addKv := insertStmt.(*dst.ReturnStmt).Results[0].(*dst.UnaryExpr).X.(*dst.CompositeLit).Elts
					//未导入的才添加
					for _, add := range addKv {
						for j, exist := range kv {
							//如果已经包含
							if add.(*dst.KeyValueExpr).Key.(*dst.Ident).Name == exist.(*dst.KeyValueExpr).Key.(*dst.Ident).Name {
								kv = append(kv[:j], kv[j+1:]...)
								//break
							}
						}

					}
					ret.(*dst.UnaryExpr).X.(*dst.CompositeLit).Elts = kv
				}
			}

		}
	}

	return vi
}
