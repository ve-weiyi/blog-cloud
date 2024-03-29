package initest

import (
	"testing"

	"github.com/ve-weiyi/ve-blog-golang/server/global"
	"github.com/ve-weiyi/ve-blog-golang/server/infra/rbac"
	"github.com/ve-weiyi/ve-blog-golang/server/utils/jsonconv"
)

func TestRbacHolder(t *testing.T) {
	Init()
	r := rbac.NewPermissionHolder(global.DB, global.LOG)
	permission, err := r.FindApiPermission("/api/v1/user", "GET")
	t.Log(err)

	t.Log(jsonconv.ObjectToJsonIndent(permission))
}
