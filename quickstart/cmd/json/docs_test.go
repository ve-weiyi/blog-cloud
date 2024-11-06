package swagger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/zeromicro/go-zero/tools/goctl/api/parser"
)

func TestDocs(t *testing.T) {
	genSwagJson()
}

func genSwagJson() error {
	var host, basePath, schemes string

	filename := "/Users/weiyi/Github/ve-blog-golang/zero/service/api/blog/proto/blog.api"
	sp, err := parser.Parse(filename)
	if err != nil {
		fmt.Println(err)
	}

	swagger, err := applyGenerate(sp, host, basePath, schemes)
	if err != nil {
		fmt.Println(err)
	}

	var formatted bytes.Buffer
	enc := json.NewEncoder(&formatted)
	enc.SetIndent("", "  ")

	if err := enc.Encode(swagger); err != nil {
		fmt.Println(err)
	}

	output := "./test.json"

	err = os.WriteFile(output, formatted.Bytes(), 0666)
	if err != nil {
		fmt.Println(err)
	}
	return err
}
