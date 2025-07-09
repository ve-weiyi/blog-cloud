package inject

import (
	"bytes"
	"go/parser"
	"go/token"
	"log"
	"os"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
)

type Meta interface {
	Visit(node dst.Node) dst.Visitor
	RollBack(node dst.Node) dst.Visitor

	GetNode() dst.Node
	GetCode() string
}

type AstInjectMeta struct {
	Key         string
	FilePath    string
	ImportMetas []*ImportMeta // 导入包
	StructMetas []*StructMeta // 结构体中插入属性
	FuncMetas   []*FuncMeta   // 函数中插入代码
	DeclMetas   []*DeclMeta   // 文件中插入定义声明
}

func (vi *AstInjectMeta) Visit(node dst.Node) dst.Visitor {
	if node == nil {
		return nil
	}
	switch n := node.(type) {
	case *dst.ImportSpec:
		if n.Path != nil {
			log.Printf("Import: %s\n", n.Path.Value)
		}
	case *dst.TypeSpec:
		if n.Name != nil {
			log.Printf("Type: %s\n", n.Name.Name)
		}
	case *dst.FuncDecl:
		if n.Name != nil {
			log.Printf("Function: %s\n", n.Name.Name)
		}
	}

	return vi
}

func (vi *AstInjectMeta) Walk() error {
	fSet := token.NewFileSet()
	fParser, err := decorator.ParseFile(fSet, vi.FilePath, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	log.Println("Start...")
	dst.Print(fParser)
	dst.Walk(vi, fParser)
	return nil
}

func (vi *AstInjectMeta) Execute() error {
	fSet := token.NewFileSet()
	fParser, err := decorator.ParseFile(fSet, vi.FilePath, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	//当没有任何import时，需要 new 一个空的import列表
	hasImport := false
	for _, node := range fParser.Decls {
		if genDecl, ok := node.(*dst.GenDecl); ok {
			if genDecl.Tok == token.IMPORT {
				hasImport = true
			}
		}
	}
	if !hasImport {
		log.Println("new import")
		fParser.Decls = append([]dst.Decl{
			&dst.GenDecl{
				Tok:   token.IMPORT,
				Specs: []dst.Spec{},
			},
		}, fParser.Decls...)
	}
	//log.Println("Execute start--")
	//dst.Println(fParser)

	for _, vi := range vi.ImportMetas {
		vi.fset = fSet
		dst.Inspect(fParser, func(node dst.Node) bool {
			vi.Visit(node)
			return true
		})
	}
	for _, vi := range vi.StructMetas {
		vi.fset = fSet
		dst.Inspect(fParser, func(node dst.Node) bool {
			vi.Visit(node)
			return true
		})
	}
	for _, vi := range vi.FuncMetas {
		vi.fset = fSet
		dst.Inspect(fParser, func(node dst.Node) bool {
			vi.Visit(node)
			return true
		})
	}
	for _, vi := range vi.DeclMetas {
		vi.fset = fSet
		dst.Inspect(fParser, func(node dst.Node) bool {
			vi.Visit(node)
			return true
		})
	}
	//log.Println("Execute end--")
	//dst.Print(fParser)
	var output []byte
	buffer := bytes.NewBuffer(output)
	err = decorator.Fprint(buffer, fParser)
	if err != nil {
		panic(err)
	}

	//log.Println(buffer)
	// 写回数据
	return os.WriteFile(vi.FilePath, buffer.Bytes(), 0o600)
}

func (vi *AstInjectMeta) RollBack() error {
	fSet := token.NewFileSet()
	fParser, err := decorator.ParseFile(fSet, vi.FilePath, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	//log.Println("RollBack start--")
	//dst.Println(token.NewFileSet(), fParser)
	for _, vi := range vi.ImportMetas {
		dst.Inspect(fParser, func(node dst.Node) bool {
			vi.RollBack(node)
			return true
		})
	}

	for _, vi := range vi.StructMetas {
		dst.Inspect(fParser, func(node dst.Node) bool {
			vi.RollBack(node)
			return true
		})
	}

	for _, vi := range vi.FuncMetas {
		dst.Inspect(fParser, func(node dst.Node) bool {
			vi.RollBack(node)
			return true
		})
	}
	//log.Println("RollBack end--")
	//dst.Println(fParser)
	var output []byte
	buffer := bytes.NewBuffer(output)
	err = decorator.Fprint(buffer, fParser)
	if err != nil {
		panic(err)
	}

	//log.Println(buffer)
	// 写回数据
	return os.WriteFile(vi.FilePath, buffer.Bytes(), 0o600)
}
