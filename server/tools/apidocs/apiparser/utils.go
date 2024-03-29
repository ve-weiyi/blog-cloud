package apiparser

import (
	"fmt"
	"os"
	"path/filepath"
)

// 深度遍历目录下的所有文件，包括目录和文件
func VisitFile(root string, visitFile func(path string, f os.FileInfo, err error) error) {
	err := filepath.Walk(root, visitFile)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}

func JoinArray(arr []string) string {
	var result string
	for i, v := range arr {
		result += v
		if i < len(arr)-1 {
			result += ", "
		}
	}
	return result
}
