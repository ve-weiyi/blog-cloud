package request

import (
	"encoding/json"

	"github.com/gin-gonic/gin"

	"github.com/ve-weiyi/ve-blog-golang/pkg/utils/jsonconv"
)

// 把请求参数转换为小写
func BindJSONIgnoreCase(c *gin.Context, req interface{}) (err error) {
	var tmp map[string]interface{}
	err = c.ShouldBindJSON(&tmp)
	if err != nil {
		return err
	}
	//如果obj已经是指针，则此处不需要指针
	js := jsonconv.AnyToJsonSnake(tmp)
	err = json.Unmarshal([]byte(js), req)
	if err != nil {
		return err
	}
	return nil
}
