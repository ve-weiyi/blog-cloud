package notificationservicelogic

import (
	"strings"

	"github.com/ve-weiyi/blog-cloud/cmd/rpc/blog/internal/pb/notificationrpc"
	"github.com/ve-weiyi/blog-cloud/infra/domain/notify"
)

// pageOf 归一化分页参数。
func pageOf(q *notificationrpc.ListQuery) (page, size int) {
	page, size = 1, 20
	if q == nil {
		return page, size
	}
	if q.Page > 0 {
		page = int(q.Page)
	}
	if q.PageSize > 0 {
		size = int(q.PageSize)
	}
	return page, size
}

// where 拼装等值条件，空值跳过。
type where struct {
	clauses []string
	args    []interface{}
}

func (w *where) eq(column string, value *string) {
	if value == nil || *value == "" {
		return
	}
	w.clauses = append(w.clauses, column+" = ?")
	w.args = append(w.args, *value)
}

func (w *where) eqInt(column string, value *int64) {
	if value == nil || *value <= 0 {
		return
	}
	w.clauses = append(w.clauses, column+" = ?")
	w.args = append(w.args, *value)
}

func (w *where) build() (string, []interface{}) {
	return strings.Join(w.clauses, " AND "), w.args
}

// groupsOf 把外部传入的偏好分组名转成目录里的强类型值。
func groupsOf(in []string) []notify.Group {
	if len(in) == 0 {
		return nil
	}
	out := make([]notify.Group, 0, len(in))
	for _, g := range in {
		out = append(out, notify.Group(g))
	}
	return out
}

// strPtr 把空串转成 NULL，避免把「没填」写成空 JSON。
func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
