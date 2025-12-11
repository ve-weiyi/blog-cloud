package astx

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
)

type AstInjectionMeta struct {
	FilePath    string
	ImportMetas []*ImportMeta
	StructMetas []*StructMeta
	FuncMetas   []*FuncMeta
}

func (vi *AstInjectionMeta) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	switch n := node.(type) {
	case *ast.ImportSpec:
		if n.Path != nil {
			log.Printf("Import: %s\n", n.Path.Value)
		}
	case *ast.TypeSpec:
		if n.Name != nil {
			log.Printf("Type: %s\n", n.Name.Name)
		}
	case *ast.FuncDecl:
		if n.Name != nil {
			log.Printf("Function: %s\n", n.Name.Name)
		}
	}

	return vi
}

func (vi *AstInjectionMeta) Walk() error {
	fSet := token.NewFileSet()
	fParser, err := parser.ParseFile(fSet, vi.FilePath, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	ast.Walk(vi, fParser)
	return nil
}

func (vi *AstInjectionMeta) Inject() error {
	fSet := token.NewFileSet()
	fParser, err := parser.ParseFile(fSet, vi.FilePath, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	//当没有任何import时，需要 new 一个空的import列表
	hasImport := false
	for _, node := range fParser.Decls {
		if genDecl, ok := node.(*ast.GenDecl); ok {
			if genDecl.Tok == token.IMPORT {
				hasImport = true
			}
		}
	}
	if !hasImport {
		log.Println("new import")
		fParser.Decls = append([]ast.Decl{
			&ast.GenDecl{
				Doc:    nil,
				TokPos: 11,
				Tok:    token.IMPORT,
				Lparen: 0,
				Specs:  []ast.Spec{},
				Rparen: 0,
			},
		}, fParser.Decls...)
	}
	//log.Println("Execute start--")
	//ast.Println(token.NewFileSet(), fParser)

	for _, vi := range vi.ImportMetas {
		vi.fset = fSet
		ast.Inspect(fParser, func(node ast.Node) bool {
			vi.Visit(node)
			return true
		})
	}
	for _, vi := range vi.StructMetas {
		vi.fset = fSet
		ast.Inspect(fParser, func(node ast.Node) bool {
			vi.Visit(node)
			return true
		})
	}
	for _, vi := range vi.FuncMetas {
		vi.fset = fSet
		ast.Inspect(fParser, func(node ast.Node) bool {
			vi.Visit(node)
			return true
		})
	}
	//log.Println("Execute end--")
	//ast.Println(token.NewFileSet(), fParser)
	var output []byte
	buffer := bytes.NewBuffer(output)
	err = format.Node(buffer, fSet, fParser)
	if err != nil {
		log.Fatal(err)
	}
	// 写回数据
	return os.WriteFile(vi.FilePath, buffer.Bytes(), 0o600)
}

func (vi *AstInjectionMeta) RollBack() error {
	fSet := token.NewFileSet()
	fParser, err := parser.ParseFile(fSet, vi.FilePath, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	//log.Println("RollBack start--")
	//ast.Println(token.NewFileSet(), fParser)
	for _, vi := range vi.ImportMetas {
		ast.Inspect(fParser, func(node ast.Node) bool {
			vi.RollBack(node)
			return true
		})
		vi.RollBackImportSpec(fParser)
	}

	for _, vi := range vi.StructMetas {
		ast.Inspect(fParser, func(node ast.Node) bool {
			vi.RollBack(node)
			return true
		})
	}

	for _, vi := range vi.FuncMetas {
		ast.Inspect(fParser, func(node ast.Node) bool {
			vi.RollBack(node)
			return true
		})
	}
	//log.Println("RollBack end--")
	//ast.Println(token.NewFileSet(), fParser.Imports)
	var output []byte
	buffer := bytes.NewBuffer(output)
	err = format.Node(buffer, fSet, fParser)
	if err != nil {
		log.Fatal(err)
	}
	// 写回数据
	return os.WriteFile(vi.FilePath, buffer.Bytes(), 0o600)
}
