package apiparser

import (
	"go/ast"
	"regexp"
	"strings"
)

type ApiCommentLine struct {
	Tag     string
	Content string
}

// 获取函数注释
func ParseFunctionDoc(f *ast.FuncDecl) []ApiCommentLine {
	if f.Doc != nil {
		// return f.Doc.Text()
		// 定义一个map来存储注释
		comments := make([]ApiCommentLine, 0)
		// 定义正则表达式，匹配形式为 @标签名 内容
		re := regexp.MustCompile(`@(\w+)\s+(.+?)\s*$`)
		// 获取函数节点上方的注释
		for _, comment := range f.Doc.List {
			commentText := strings.TrimSpace(comment.Text)
			// 使用正则表达式提取注释中的标签和内容
			matches := re.FindAllStringSubmatch(commentText, -1)
			for _, match := range matches {
				tag := strings.TrimSpace(match[1])
				comments = append(comments, ApiCommentLine{Tag: tag, Content: strings.TrimSpace(match[2])})
			}
		}
		return comments
	}

	return nil
}
