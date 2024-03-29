package inject

import (
	"fmt"
	"go/token"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
)

type DeclMeta struct {
	FuncName string //插入的方法
	FuncPos  int    //插入的位置

	Variables  []string      //左边的值
	Symbol     string        //符号
	IdentNames []string      //右边的值，当为return时，是右边变量名
	Parameters []interface{} //调用参数
	Comments   []string      //注释

	hasImported bool
	fset        *token.FileSet
	node        dst.Decl
}

func (vi *DeclMeta) Visit(node dst.Node) dst.Visitor {
	switch genDecl := node.(type) {
	case *dst.File:
		// 判断是否已经import
		for _, v := range genDecl.Decls {
			switch n := v.(type) {
			case *dst.FuncDecl:
				if n.Name.Name == vi.FuncName {
					vi.hasImported = true
				}
			}
		}

		if !vi.hasImported {
			importSp := vi.GetNode()
			//放到最后
			genDecl.Decls = append(genDecl.Decls, importSp.(*dst.FuncDecl))

			vi.hasImported = true
		}
	}
	return vi
}

func (vi *DeclMeta) RollBack(node dst.Node) dst.Visitor {
	return nil
}

func (vi *DeclMeta) GetNode() dst.Node {
	return vi.node
}

func (vi *DeclMeta) GetCode() string {
	return ""
}

func NewDeclMeta(importCode string) *DeclMeta {
	meta := &DeclMeta{}
	input := fmt.Sprintf(`
package main

%v 
`, importCode)
	fParse, err := decorator.ParseFile(token.NewFileSet(), "", input, 0)
	if err != nil {
		fmt.Println("Error parsing file:", err)
	}
	//dst.Print(fParse)
	dst.Inspect(fParse, func(n dst.Node) bool {
		switch x := n.(type) {
		case *dst.FuncDecl:
			//dst.Print(x)
			meta.FuncName = x.Name.Name
			meta.node = x
			for _, field := range x.Body.List {
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

			meta.Comments = x.Decs.NodeDecs.Start.All()
		}

		return true
	})

	//log.Println("GetCode", meta.GetCode())
	//fmt.Println("NewFuncMete", jsonconv.ObjectToJsonIndent(meta))
	return meta
}
