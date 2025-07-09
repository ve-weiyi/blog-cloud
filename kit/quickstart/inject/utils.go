package inject

import (
	"go/token"
	"log"
	"reflect"
	"strconv"

	"github.com/dave/dst"
)

func tokenToKind(t token.Token) reflect.Kind {
	switch t {
	case token.INT:
		return reflect.Int
	case token.FLOAT:
		return reflect.Float64
	case token.STRING:
		return reflect.String
	case token.CHAR:
		return reflect.Int32
	default:
		return reflect.Invalid
	}
}

func kindToToken(k reflect.Kind) token.Token {
	switch k {
	case reflect.Int:
		return token.INT
	case reflect.Float64:
		return token.FLOAT
	case reflect.String:
		return token.STRING
	case reflect.Int32:
		return token.CHAR
	default:
		return token.ILLEGAL
	}
}

func inferType(str string) (interface{}, error) {
	// 尝试将字符串解析为int
	i, err := strconv.Atoi(str)
	if err == nil {
		return i, nil
	}

	// 尝试将字符串解析为float
	f, err := strconv.ParseFloat(str, 64)
	if err == nil {
		return f, nil
	}

	// 尝试将字符串解析为带引号的string
	s, err := strconv.Unquote(str)
	if err == nil {
		return s, nil
	}

	// 如果都不匹配，则返回原始字符串
	return str, nil
}

func ExtractIdents(node dst.Node) []*dst.Ident {
	var idents []*dst.Ident

	switch n := node.(type) {
	case *dst.AssignStmt:
		return ExtractIdents(n.Rhs[0])

	case *dst.SelectorExpr:
		log.Println("SelectorExpr", n.Sel.Name)
		idents = append(idents, ExtractIdents(n.X)...)
		idents = append(idents, ExtractIdents(n.Sel)...)
	case *dst.Ident:
		log.Println("Ident", n.Name)
		idents = append(idents, n)
	case *dst.KeyValueExpr:
		log.Println("KeyValueExpr", n.Key)
		// 判断是否是复合字面值表达式的键值对
		idents = append(idents, ExtractIdents(n.Value)...)
	case *dst.CompositeLit:
		log.Println("CompositeLit", n.Type)
		idents = append(idents, ExtractIdents(n.Type)...)
		for _, elt := range n.Elts {
			idents = append(idents, ExtractIdents(elt)...)
		}
	default:
		log.Printf("default %T", n)
	}

	return idents
}

func insertStatements(stmts []dst.Stmt, pos int, toInsert ...dst.Stmt) []dst.Stmt {
	return append(stmts[:pos], append(toInsert, stmts[pos:]...)...)
}
